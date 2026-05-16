from typing import List, Dict, Any
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich.text import Text
from rich import box
from datetime import datetime

from .models import Board, Card, Column, Priority

console = Console()


def priority_color(priority: Priority) -> str:
    colors = {
        Priority.LOW: "green",
        Priority.MEDIUM: "yellow",
        Priority.HIGH: "orange",
        Priority.URGENT: "red",
    }
    return colors.get(priority, "white")


def priority_icon(priority: Priority) -> str:
    icons = {
        Priority.LOW: "↓",
        Priority.MEDIUM: "→",
        Priority.HIGH: "↑",
        Priority.URGENT: "!!",
    }
    return icons.get(priority, "")


def format_due_date(card: Card) -> Text:
    if not card.due_date:
        return Text("无", style="dim")
    try:
        due = datetime.fromisoformat(card.due_date).date()
        days = card.days_until_due()
        if card.is_overdue():
            return Text(f"已过期 {abs(days)} 天", style="bold red")
        elif days == 0:
            return Text("今天截止", style="bold yellow")
        elif days <= 3:
            return Text(f"{due} ({days}天后)", style="yellow")
        else:
            return Text(f"{due}", style="green")
    except (ValueError, TypeError):
        return Text(card.due_date, style="dim")


def format_tags(tags: List[str]) -> Text:
    if not tags:
        return Text("")
    tag_text = Text()
    for tag in tags:
        tag_text.append(f"[{tag}] ", style="cyan")
    return tag_text


def format_card_title(card: Card, show_id: bool = False) -> Text:
    title = Text()
    title.append(f"{priority_icon(card.priority)} ", style=priority_color(card.priority))
    title.append(card.title)
    if show_id:
        title.append(f" ({card.id[:8]})", style="dim")
    if card.is_overdue():
        title.append(" ⚠️", style="red")
    return title


def display_boards(boards: List[Board]) -> None:
    if not boards:
        console.print("[yellow]暂无看板[/yellow]")
        return
    table = Table(show_header=True, header_style="bold magenta", box=box.ROUNDED)
    table.add_column("ID", style="dim", width=36)
    table.add_column("名称")
    table.add_column("描述", style="dim")
    table.add_column("列数", justify="right")
    table.add_column("卡片数", justify="right")
    table.add_column("状态")
    table.add_column("创建时间", style="dim")
    for board in boards:
        status = "已归档" if board.archived else "正常"
        status_style = "yellow" if board.archived else "green"
        table.add_row(
            board.id,
            board.name,
            board.description or "-",
            str(len(board.columns)),
            str(len(board.cards)),
            Text(status, style=status_style),
            board.created_at[:19].replace("T", " "),
        )
    console.print(table)


def display_board(board: Board, show_details: bool = False) -> None:
    console.print(f"\n[bold magenta]📋 {board.name}[/bold magenta]")
    if board.description:
        console.print(f"[dim]{board.description}[/dim]")
    if board.archived:
        console.print("[yellow]⚠️  此看板已归档[/yellow]")
    console.print()
    if not board.columns:
        console.print("[yellow]暂无列，请先添加列[/yellow]")
        return
    for column in board.sorted_columns():
        cards = board.get_cards_in_column(column.id)
        title = f"[bold blue]{column.name}[/bold blue] ({len(cards)})"
        if column.is_completed:
            title += " [green]✓[/green]"
        content = []
        for card in cards:
            line = format_card_title(card, show_id=show_details)
            if show_details:
                if card.tags:
                    line.append(" ")
                    line.append(format_tags(card.tags))
                line.append("\n  ")
                line.append(format_due_date(card))
            content.append(line)
        if not content:
            content_text = Text("空", style="dim")
        else:
            content_text = Text("\n").join(content)
        panel = Panel(content_text, title=title, border_style="blue", box=box.ROUNDED)
        console.print(panel)
    unassigned = [c for c in board.cards if not c.column_id]
    if unassigned:
        title = f"[bold yellow]未分配[/bold yellow] ({len(unassigned)})"
        content = Text("\n").join([format_card_title(c, show_id=show_details) for c in unassigned])
        panel = Panel(content, title=title, border_style="yellow", box=box.ROUNDED)
        console.print(panel)


def display_card(card: Card, column: Column = None) -> None:
    title = format_card_title(card)
    console.print(f"\n[bold magenta]🃏 {title}[/bold magenta]")
    console.print(f"[dim]ID: {card.id}[/dim]")
    if column:
        console.print(f"列: [blue]{column.name}[/blue]")
    if card.description:
        console.print(f"\n[bold]描述:[/bold]\n{card.description}")
    console.print(f"\n优先级: [{priority_color(card.priority)}]{card.priority.value}[/]")
    if card.tags:
        console.print(f"标签: {format_tags(card.tags)}")
    console.print(f"截止日期: {format_due_date(card)}")
    console.print(f"创建时间: {card.created_at[:19].replace('T', ' ')}")
    if card.completed_at:
        console.print(f"完成时间: {card.completed_at[:19].replace('T', ' ')} [green]✓[/green]")


def display_cards(cards: List[Card], board: Board = None, show_details: bool = False) -> None:
    if not cards:
        console.print("[yellow]未找到匹配的卡片[/yellow]")
        return
    for i, card in enumerate(cards, 1):
        line = Text(f"{i}. ")
        line.append(format_card_title(card, show_id=True))
        if show_details:
            if board:
                column = board.get_column(card.column_id)
                if column:
                    line.append(f" [dim][{column.name}][/dim]")
            if card.tags:
                line.append(" ")
                line.append(format_tags(card.tags))
            line.append(" ")
            line.append(format_due_date(card))
        console.print(line)


def display_stats(stats: Dict[str, Any]) -> None:
    console.print("\n[bold magenta]📊 统计信息[/bold magenta]\n")
    table = Table(show_header=False, box=box.ROUNDED, border_style="dim")
    table.add_column(style="bold")
    table.add_column()
    table.add_row("总卡片数", str(stats.get("total_cards", 0)))
    table.add_row("已完成", f"[green]{stats.get('completed_cards', 0)}[/green]")
    table.add_row("完成率", f"[blue]{stats.get('completion_rate', 0)}%[/blue]")
    table.add_row("平均完成天数", str(stats.get("avg_completion_days", 0)))
    table.add_row("过期卡片", f"[red]{stats.get('overdue_cards', 0)}[/red]")
    console.print(table)
    col_dist = stats.get("column_distribution", {})
    if col_dist:
        console.print("\n[bold]列分布:[/bold]")
        for col, count in col_dist.items():
            bar = "█" * min(count, 20)
            console.print(f"  {col:<10} {count:>4} {bar}")
    prio_dist = stats.get("priority_distribution", {})
    if prio_dist:
        console.print("\n[bold]优先级分布:[/bold]")
        for prio, count in prio_dist.items():
            color = priority_color(Priority(prio))
            bar = "█" * min(count, 20)
            console.print(f"  [{color}]{prio:<10}[/{color}] {count:>4} {bar}")


def display_columns(columns: List[Column]) -> None:
    if not columns:
        console.print("[yellow]暂无列[/yellow]")
        return
    table = Table(show_header=True, header_style="bold magenta", box=box.ROUNDED)
    table.add_column("ID", style="dim", width=36)
    table.add_column("名称")
    table.add_column("顺序", justify="right")
    table.add_column("完成列")
    for col in sorted(columns, key=lambda c: c.order):
        is_completed = "✓" if col.is_completed else "-"
        completed_style = "green" if col.is_completed else "dim"
        table.add_row(
            col.id,
            col.name,
            str(col.order),
            Text(is_completed, style=completed_style),
        )
    console.print(table)


def success(message: str) -> None:
    console.print(f"[green]✓ {message}[/green]")


def error(message: str) -> None:
    console.print(f"[red]✗ {message}[/red]")


def warning(message: str) -> None:
    console.print(f"[yellow]⚠ {message}[/yellow]")


def info(message: str) -> None:
    console.print(f"[blue]ℹ {message}[/blue]")
