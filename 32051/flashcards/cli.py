from __future__ import annotations

import os
import sys
from datetime import datetime
from typing import List, Optional

import typer
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.text import Text
from rich.prompt import Prompt, Confirm
from rich.progress import Progress, SpinnerColumn, TextColumn

from .models import Difficulty, FlashCard
from .session import StudySession
from .sm2 import format_interval_days, get_quality_description
from .stats import StatsCalculator
from .storage import Storage


app = typer.Typer(
    help="间隔重复闪卡学习工具 - 基于SM-2算法",
    add_completion=False,
    rich_markup_mode="rich",
)
console = Console()
storage = Storage()


def _get_deck_by_identifier(identifier: str):
    deck = storage.get_deck(identifier)
    if not deck:
        deck = storage.get_deck_by_name(identifier)
    return deck


def _print_error(message: str) -> None:
    console.print(f"[bold red]错误:[/bold red] {message}")


def _print_success(message: str) -> None:
    console.print(f"[bold green]成功:[/bold green] {message}")


def _print_warning(message: str) -> None:
    console.print(f"[bold yellow]警告:[/bold yellow] {message}")


def _print_info(message: str) -> None:
    console.print(f"[bold blue]信息:[/bold blue] {message}")


@app.command("list", help="列出所有牌组")
def list_decks() -> None:
    decks = storage.list_decks()
    if not decks:
        _print_info("还没有创建任何牌组，使用 'flashcards create <名称>' 创建第一个牌组")
        return

    table = Table(title="牌组列表", show_lines=True)
    table.add_column("ID", style="cyan", no_wrap=True)
    table.add_column("名称", style="magenta")
    table.add_column("卡片数", justify="right")
    table.add_column("待复习", justify="right", style="yellow")
    table.add_column("创建时间", style="dim")
    table.add_column("更新时间", style="dim")

    for deck in decks:
        due_count = len(deck.get_due_cards())
        created = datetime.fromisoformat(deck.created_at).strftime("%Y-%m-%d %H:%M")
        updated = datetime.fromisoformat(deck.updated_at).strftime("%Y-%m-%d %H:%M")
        table.add_row(
            deck.id[:8],
            deck.name,
            str(len(deck.cards)),
            str(due_count),
            created,
            updated,
        )

    console.print(table)


@app.command("create", help="创建新牌组")
def create_deck(name: str) -> None:
    try:
        deck = storage.create_deck(name)
        _print_success(f"已创建牌组 '{deck.name}' (ID: {deck.id[:8]})")
    except ValueError as e:
        _print_error(str(e))
        raise typer.Exit(code=1)


@app.command("rename", help="重命名牌组")
def rename_deck(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    new_name: str = typer.Argument(..., help="新名称"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    try:
        updated = storage.rename_deck(deck.id, new_name)
        if updated:
            _print_success(f"牌组已重命名为 '{new_name}'")
    except ValueError as e:
        _print_error(str(e))
        raise typer.Exit(code=1)


@app.command("delete", help="删除牌组")
def delete_deck(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    force: bool = typer.Option(False, "--force", "-f", help="不确认直接删除"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    if not force:
        console.print(Panel(
            f"[bold red]即将删除牌组:[/bold red] {deck.name}\n"
            f"包含 {len(deck.cards)} 张卡片和 {len(deck.study_records)} 条学习记录\n"
            f"[bold]此操作不可恢复！[/bold]",
            title="删除确认",
            border_style="red",
        ))
        confirm = Confirm.ask("确认要删除这个牌组吗?", default=False)
        if not confirm:
            _print_info("已取消删除")
            raise typer.Exit()

    if storage.delete_deck(deck.id):
        _print_success(f"已删除牌组 '{deck.name}'")
    else:
        _print_error("删除失败")


@app.command("add-card", help="添加卡片到牌组")
def add_card(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    front: str = typer.Option(..., "--front", "-q", prompt="正面内容", help="卡片正面（问题）"),
    back: str = typer.Option(..., "--back", "-a", prompt="背面内容", help="卡片背面（答案）"),
    tags: Optional[List[str]] = typer.Option(None, "--tag", "-t", help="标签，可多次指定"),
    difficulty: str = typer.Option("medium", "--difficulty", "-d", help="难度: easy/medium/hard"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    try:
        diff = Difficulty(difficulty.lower())
    except ValueError:
        _print_error(f"无效的难度等级 '{difficulty}'，必须是: easy, medium, hard")
        raise typer.Exit(code=1)

    if not front.strip():
        _print_error("正面内容不能为空")
        raise typer.Exit(code=1)
    if not back.strip():
        _print_error("背面内容不能为空")
        raise typer.Exit(code=1)

    for existing in deck.cards:
        if existing.front == front.strip():
            _print_warning(f"已存在正面相同的卡片: '{front}'")
            overwrite = Confirm.ask("是否覆盖?", default=False)
            if not overwrite:
                _print_info("已取消添加")
                raise typer.Exit()
            existing.back = back.strip()
            existing.tags = tags or []
            existing.difficulty = diff
            existing.touch()
            storage.save_deck(deck)
            _print_success("已更新卡片")
            return

    card = FlashCard(
        front=front.strip(),
        back=back.strip(),
        tags=tags or [],
        difficulty=diff,
    )
    storage.add_card(deck.id, card)
    _print_success(f"已添加卡片到 '{deck.name}'")


@app.command("list-cards", help="列出牌组中的卡片")
def list_cards(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    keyword: Optional[str] = typer.Option(None, "--search", "-s", help="按关键词搜索"),
    tag: Optional[List[str]] = typer.Option(None, "--tag", "-t", help="按标签筛选，可多次指定"),
    difficulty: Optional[str] = typer.Option(None, "--difficulty", "-d", help="按难度筛选: easy/medium/hard"),
    due_only: bool = typer.Option(False, "--due", help="只显示待复习卡片"),
    limit: int = typer.Option(50, "--limit", "-n", help="显示数量限制"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    cards = deck.cards

    if due_only:
        cards = [c for c in cards if c.is_due()]

    if keyword:
        kw = keyword.lower()
        cards = [
            c for c in cards
            if kw in c.front.lower() or kw in c.back.lower()
        ]

    if tag:
        cards = [
            c for c in cards
            if any(t in c.tags for t in tag)
        ]

    if difficulty:
        try:
            diff = Difficulty(difficulty.lower())
            cards = [c for c in cards if c.difficulty == diff]
        except ValueError:
            _print_error(f"无效的难度等级 '{difficulty}'")
            raise typer.Exit(code=1)

    if not cards:
        _print_info("没有符合条件的卡片")
        return

    cards = cards[:limit]

    table = Table(title=f"卡片列表 - {deck.name}", show_lines=True)
    table.add_column("ID", style="cyan", no_wrap=True)
    table.add_column("正面", style="white", overflow="fold")
    table.add_column("标签", style="green")
    table.add_column("难度", style="yellow")
    table.add_column("状态", style="blue")
    table.add_column("下次复习", style="magenta")

    for card in cards:
        status = "待复习" if card.is_due() else f"已掌握({card.repetitions}次)"
        next_review = format_interval_days(card.interval) if card.interval > 0 else "今天"
        tags_str = ", ".join(card.tags) if card.tags else "-"
        table.add_row(
            card.id[:8],
            card.front[:60] + "..." if len(card.front) > 60 else card.front,
            tags_str,
            card.difficulty.value,
            status,
            next_review,
        )

    console.print(table)
    _print_info(f"显示 {len(cards)} 张卡片，总共 {len(deck.cards)} 张")


@app.command("show-card", help="查看卡片详情")
def show_card(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    card_id: str = typer.Argument(..., help="卡片ID"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    card = deck.get_card(card_id)
    if not card:
        for c in deck.cards:
            if c.id.startswith(card_id):
                card = c
                break

    if not card:
        _print_error(f"找不到卡片 '{card_id}'")
        raise typer.Exit(code=1)

    console.print(Panel(f"[bold]{card.front}[/bold]", title="正面", border_style="cyan"))
    console.print(Panel(f"{card.back}", title="背面", border_style="green"))

    info_table = Table(show_header=False, show_lines=False)
    info_table.add_column("属性", style="bold")
    info_table.add_column("值")

    info_table.add_row("ID", card.id)
    info_table.add_row("标签", ", ".join(card.tags) if card.tags else "-")
    info_table.add_row("难度", card.difficulty.value)
    info_table.add_row("复习次数", str(card.review_count))
    info_table.add_row("连续正确", str(card.repetitions))
    info_table.add_row("难度系数", str(card.ease_factor))
    info_table.add_row("当前间隔", format_interval_days(card.interval))
    info_table.add_row("上次复习", datetime.fromisoformat(card.last_reviewed_at).strftime("%Y-%m-%d %H:%M") if card.last_reviewed_at else "从未复习")
    info_table.add_row("下次复习", datetime.fromisoformat(card.next_review_at).strftime("%Y-%m-%d %H:%M"))
    info_table.add_row("创建时间", datetime.fromisoformat(card.created_at).strftime("%Y-%m-%d %H:%M"))

    console.print(info_table)


@app.command("edit-card", help="编辑卡片")
def edit_card(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    card_id: str = typer.Argument(..., help="卡片ID"),
    front: Optional[str] = typer.Option(None, "--front", "-q", help="新的正面内容"),
    back: Optional[str] = typer.Option(None, "--back", "-a", help="新的背面内容"),
    tags: Optional[List[str]] = typer.Option(None, "--tag", "-t", help="新标签，可多次指定"),
    difficulty: Optional[str] = typer.Option(None, "--difficulty", "-d", help="新难度: easy/medium/hard"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    card = deck.get_card(card_id)
    if not card:
        for c in deck.cards:
            if c.id.startswith(card_id):
                card = c
                break

    if not card:
        _print_error(f"找不到卡片 '{card_id}'")
        raise typer.Exit(code=1)

    if front is not None:
        card.front = front
    if back is not None:
        card.back = back
    if tags is not None:
        card.tags = tags
    if difficulty is not None:
        try:
            card.difficulty = Difficulty(difficulty.lower())
        except ValueError:
            _print_error(f"无效的难度等级 '{difficulty}'")
            raise typer.Exit(code=1)

    card.touch()
    storage.save_deck(deck)
    _print_success(f"已更新卡片")


@app.command("delete-card", help="删除卡片")
def delete_card(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    card_id: str = typer.Argument(..., help="卡片ID"),
    force: bool = typer.Option(False, "--force", "-f", help="不确认直接删除"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    card = deck.get_card(card_id)
    if not card:
        for c in deck.cards:
            if c.id.startswith(card_id):
                card = c
                break

    if not card:
        _print_error(f"找不到卡片 '{card_id}'")
        raise typer.Exit(code=1)

    if not force:
        console.print(Panel(
            f"[bold]正面:[/bold] {card.front}\n"
            f"[bold]背面:[/bold] {card.back[:100]}...",
            title="删除卡片确认",
            border_style="yellow",
        ))
        confirm = Confirm.ask("确认要删除这张卡片吗?", default=False)
        if not confirm:
            _print_info("已取消删除")
            raise typer.Exit()

    if storage.delete_card(deck.id, card.id):
        _print_success("已删除卡片")
    else:
        _print_error("删除失败")


@app.command("study", help="开始学习会话")
def study(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    tags: Optional[List[str]] = typer.Option(None, "--tag", "-t", help="按标签筛选"),
    difficulty: Optional[str] = typer.Option(None, "--difficulty", "-d", help="按难度筛选"),
    max_cards: Optional[int] = typer.Option(None, "--max", "-n", help="最大复习数量"),
    no_shuffle: bool = typer.Option(False, "--no-shuffle", help="不随机排序"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    session = StudySession(
        deck=deck,
        tags=tags,
        difficulty=difficulty,
        max_cards=max_cards,
        shuffle=not no_shuffle,
    )

    if not session.review_queue:
        _print_info("没有待复习的卡片，继续保持！")
        raise typer.Exit()

    console.print(Panel(
        f"牌组: [bold]{deck.name}[/bold]\n"
        f"待复习: [bold yellow]{len(session.review_queue)}[/bold yellow] 张卡片\n"
        f"输入数字 0-5 评分，或输入 'q' 退出",
        title="学习会话开始",
        border_style="green",
    ))

    try:
        while session.has_next():
            current, total = session.get_progress()
            card = session.next_card()
            if not card:
                break

            console.print()
            console.print(f"[dim]进度: {current}/{total}[/dim]")
            console.print(Panel(f"[bold]{card.front}[/bold]", title=f"卡片 {current + 1}", border_style="cyan"))

            action = Prompt.ask("按 Enter 显示答案，或输入 'q' 退出", default="")
            if action.lower() == "q":
                break

            console.print(Panel(f"{card.back}", title="答案", border_style="green"))

            console.print()
            console.print("[bold]请评分:[/bold]")
            console.print("  [0] 完全忘记  [1] 毫无印象  [2] 似曾相识")
            console.print("  [3] 勉强记住  [4] 轻松记住  [5] 牢记在心")

            while True:
                quality_input = Prompt.ask("输入评分 (0-5)，或输入 's' 跳过")
                if quality_input.lower() == "s":
                    session.skip_card(card)
                    _print_info("已跳过，稍后再复习")
                    break
                if quality_input.lower() == "q":
                    session.has_next = lambda: False
                    break
                try:
                    quality = int(quality_input)
                    if 0 <= quality <= 5:
                        session.record_answer(card, quality)
                        next_interval = format_interval_days(card.interval)
                        quality_desc = get_quality_description(quality)
                        color = "green" if quality >= 3 else "red"
                        console.print(f"[{color}]{quality_desc} - 下次复习: {next_interval}[/{color}]")
                        break
                    else:
                        _print_warning("请输入 0-5 之间的数字")
                except ValueError:
                    _print_warning("请输入有效的数字")
    except KeyboardInterrupt:
        console.print()
        _print_info("学习被中断")

    result = session.finish()
    storage.save_deck(deck)

    if result["total_cards"] == 0:
        _print_info("本次没有复习任何卡片")
        return

    result_panel = Panel(
        f"复习卡片: [bold]{result['total_cards']}[/bold] 张\n"
        f"正确: [bold green]{result['correct_count']}[/bold green]  "
        f"错误: [bold red]{result['incorrect_count']}[/bold red]\n"
        f"正确率: [bold]{result['accuracy']}%[/bold]\n"
        f"用时: [bold]{result['duration_formatted']}[/bold]\n"
        f"平均每张: [bold]{result['avg_time_per_card']}[/bold] 秒",
        title="学习会话结束",
        border_style="blue",
    )
    console.print()
    console.print(result_panel)


@app.command("export", help="导出牌组卡片")
def export_deck(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    output: str = typer.Argument(..., help="输出文件路径 (.json 或 .csv)"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    ext = os.path.splitext(output)[1].lower()

    success = False
    if ext == ".json":
        success = storage.export_deck_to_json(deck.id, output)
    elif ext == ".csv":
        success = storage.export_deck_to_csv(deck.id, output)
    else:
        _print_error(f"不支持的文件格式 '{ext}'，请使用 .json 或 .csv")
        raise typer.Exit(code=1)

    if success:
        _print_success(f"已导出 {len(deck.cards)} 张卡片到 {output}")
    else:
        _print_error("导出失败")


@app.command("import", help="从文件导入卡片")
def import_cards(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    filepath: str = typer.Argument(..., help="导入文件路径 (.json 或 .csv)"),
    overwrite: bool = typer.Option(False, "--overwrite", "-o", help="覆盖已存在的卡片"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    if not os.path.exists(filepath):
        _print_error(f"文件不存在: {filepath}")
        raise typer.Exit(code=1)

    ext = os.path.splitext(filepath)[1].lower()

    has_duplicates = False
    try:
        if ext == ".json":
            import json
            with open(filepath, "r", encoding="utf-8") as f:
                cards_data = json.load(f)
            for card_data in cards_data:
                front = card_data.get("front", "").strip()
                if any(c.front == front for c in deck.cards):
                    has_duplicates = True
                    break
        elif ext == ".csv":
            import csv
            with open(filepath, "r", encoding="utf-8") as f:
                reader = csv.DictReader(f)
                for row in reader:
                    front = row.get("front", "").strip()
                    if any(c.front == front for c in deck.cards):
                        has_duplicates = True
                        break
    except Exception as e:
        _print_error(f"读取文件失败: {e}")
        raise typer.Exit(code=1)

    if has_duplicates and not overwrite:
        _print_warning("检测到重复卡片（正面内容相同）")
        overwrite = Confirm.ask("是否覆盖已存在的卡片?", default=False)

    result = {"added": 0, "updated": 0, "skipped": 0}
    if ext == ".json":
        result = storage.import_cards_from_json(deck.id, filepath, overwrite=overwrite)
    elif ext == ".csv":
        result = storage.import_cards_from_csv(deck.id, filepath, overwrite=overwrite)
    else:
        _print_error(f"不支持的文件格式 '{ext}'")
        raise typer.Exit(code=1)

    if result["added"] == 0 and result["updated"] == 0 and result["skipped"] == 0:
        _print_error("导入失败，请检查文件格式")
        raise typer.Exit(code=1)

    _print_success(
        f"导入完成: 新增 {result['added']} 张, "
        f"更新 {result['updated']} 张, "
        f"跳过 {result['skipped']} 张"
    )


@app.command("stats", help="查看学习统计")
def show_stats(
    deck_id: str = typer.Argument(..., help="牌组ID或名称"),
    days: int = typer.Option(30, "--days", "-n", help="统计天数范围"),
    tag: Optional[str] = typer.Option(None, "--tag", "-t", help="查看指定标签的统计"),
) -> None:
    deck = _get_deck_by_identifier(deck_id)
    if not deck:
        _print_error(f"找不到牌组 '{deck_id}'")
        raise typer.Exit(code=1)

    stats = StatsCalculator(deck)

    if tag:
        tag_stats = stats.get_tag_stats(tag)
        if not tag_stats:
            _print_error(f"找不到标签 '{tag}'")
            raise typer.Exit(code=1)

        table = Table(title=f"标签统计 - {tag}", show_header=False)
        table.add_column("指标", style="bold")
        table.add_column("值")
        table.add_row("卡片总数", str(tag_stats["total_cards"]))
        table.add_row("待复习", str(tag_stats["due_cards"]))
        table.add_row("总复习次数", str(tag_stats["total_reviews"]))
        table.add_row("正确率", f"{tag_stats['accuracy']}%")
        console.print(table)
        return

    overall = stats.get_overall_stats()
    streak = stats.get_streak_days()
    difficulty_bd = stats.get_difficulty_breakdown()
    tag_bd = stats.get_tag_breakdown()
    retention = stats.get_retention_rate()

    grid = Table.grid(expand=True)
    grid.add_column()
    grid.add_column()

    overview_table = Table(title="总体概览", show_header=False)
    overview_table.add_column("指标", style="bold")
    overview_table.add_column("值")
    overview_table.add_row("卡片总数", str(overall["total_cards"]))
    overview_table.add_row("待复习", f"[yellow]{overall['due_cards']}[/yellow]")
    overview_table.add_row("已学习", str(overall["reviewed_cards"]))
    overview_table.add_row("未学习", str(overall["unreviewed_cards"]))
    overview_table.add_row("总复习次数", str(overall["total_reviews"]))
    overview_table.add_row("总正确率", f"[green]{overall['accuracy']}%[/green]")
    overview_table.add_row("连续学习天数", f"[bold magenta]{streak}[/bold magenta] 天")

    diff_table = Table(title="难度分布", show_header=True)
    diff_table.add_column("难度", style="bold")
    diff_table.add_column("数量", justify="right")
    for diff in ["easy", "medium", "hard"]:
        count = difficulty_bd.get(diff, 0)
        diff_table.add_row(diff, str(count))

    grid.add_row(overview_table, diff_table)
    console.print(grid)

    if tag_bd:
        tag_table = Table(title="标签分布 (前10)", show_lines=False)
        tag_table.add_column("标签", style="green")
        tag_table.add_column("卡片数", justify="right")
        for i, (tag_name, count) in enumerate(list(tag_bd.items())[:10]):
            tag_table.add_row(tag_name, str(count))
        console.print(tag_table)

    retention_table = Table(title="记忆保持率", show_header=True)
    retention_table.add_column("时间段", style="bold")
    retention_table.add_column("保持率", justify="right")
    for period, rate in retention.items():
        color = "green" if rate >= 80 else "yellow" if rate >= 60 else "red"
        retention_table.add_row(period, f"[{color}]{rate}%[/{color}]")
    console.print(retention_table)

    daily_activity = stats.get_daily_activity(days=days)
    active_days = sum(1 for d in daily_activity if d["reviews"] > 0)
    total_reviews = sum(d["reviews"] for d in daily_activity)

    summary_table = Table(title=f"最近 {days} 天学习汇总", show_lines=False)
    summary_table.add_column("活跃天数", style="bold")
    summary_table.add_column("总复习次数", justify="right")
    summary_table.add_column("日均复习", justify="right")
    summary_table.add_row(
        f"{active_days}/{days}",
        str(total_reviews),
        f"{total_reviews / max(active_days, 1):.1f}"
    )
    console.print(summary_table)

    if active_days > 0:
        detail_table = Table(title=f"最近 {days} 天学习明细", show_lines=False)
        detail_table.add_column("日期", style="cyan", justify="center")
        detail_table.add_column("复习次数", justify="right")
        detail_table.add_column("正确数", justify="right")
        detail_table.add_column("正确率", justify="right")
        detail_table.add_column("趋势", justify="left")

        max_reviews = max(d["reviews"] for d in daily_activity) if active_days > 0 else 1

        for day_data in daily_activity:
            reviews = day_data["reviews"]
            correct = day_data["correct"]
            accuracy = day_data["accuracy"]

            if reviews == 0:
                bar = "  "
                accuracy_str = "-"
                correct_str = "-"
            else:
                bar_length = int(reviews / max(max_reviews, 1) * 10)
                bar = "█" * bar_length + "░" * (10 - bar_length)
                accuracy_str = f"{accuracy}%"
                correct_str = str(correct)

            color = "green" if accuracy >= 80 else "yellow" if accuracy >= 60 else "red"
            if reviews == 0:
                color = "dim"

            detail_table.add_row(
                day_data["date"],
                str(reviews),
                correct_str,
                f"[{color}]{accuracy_str}[/{color}]",
                f"[{color}]{bar}[/{color}]",
            )

        console.print(detail_table)


@app.callback()
def main() -> None:
    """间隔重复闪卡学习工具 - 基于SM-2算法的命令行闪卡应用"""
    pass


def run() -> None:
    try:
        app()
    except typer.Exit:
        raise
    except Exception as e:
        _print_error(f"发生错误: {e}")
        console.print("[dim]如需帮助，请使用 --help 参数[/dim]")
        sys.exit(1)


if __name__ == "__main__":
    run()
