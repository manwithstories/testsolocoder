import typer
from typing import Optional, List
from datetime import date, timedelta

from .models import Priority
from .core import KanbanManager, DuplicateResolution
from . import display

app = typer.Typer(help="命令行看板管理工具", no_args_is_help=True)
manager = KanbanManager()


def _resolve_board(board_id: Optional[str], board_name: Optional[str]):
    if board_id:
        board = manager.get_board(board_id)
        if not board:
            display.error(f"未找到ID为 {board_id} 的看板")
            return None
        return board
    if board_name:
        board = manager.get_board_by_name(board_name)
        if not board:
            display.error(f"未找到名为 '{board_name}' 的看板")
            return None
        return board
    display.error("请提供看板ID (--board-id) 或名称 (--board-name)")
    return None


def _confirm(prompt: str) -> bool:
    return typer.confirm(prompt, default=False)


@app.command()
def init():
    """初始化数据目录"""
    display.success("数据目录已初始化")
    display.info(f"数据存储位置: {manager.storage.data_dir}")


@app.command("list")
def list_boards(
    include_archived: bool = typer.Option(False, "--archived", help="包含已归档的看板"),
):
    """列出所有看板"""
    boards = manager.list_boards(include_archived=include_archived)
    display.display_boards(boards)


@app.command()
def create(
    name: str = typer.Argument(..., help="看板名称"),
    description: str = typer.Option("", "--desc", help="看板描述"),
    no_default_columns: bool = typer.Option(False, "--no-default", help="不创建默认列"),
):
    """创建新看板"""
    board = manager.create_board(name, description, default_columns=not no_default_columns)
    display.success(f"已创建看板 '{board.name}'")
    display.info(f"ID: {board.id}")


@app.command()
def show(
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
    details: bool = typer.Option(False, "--details", help="显示详细信息"),
):
    """显示看板详情"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    display.display_board(board, show_details=details)


@app.command()
def archive(
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """归档看板"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    if _confirm(f"确定要归档看板 '{board.name}' 吗？"):
        manager.archive_board(board.id)
        display.success(f"已归档看板 '{board.name}'")


@app.command()
def unarchive(
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """取消归档看板"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    manager.unarchive_board(board.id)
    display.success(f"已取消归档看板 '{board.name}'")


@app.command()
def delete(
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """删除看板（可恢复）"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    if _confirm(f"确定要删除看板 '{board.name}' 吗？此操作可在30天内撤销。"):
        deleted = manager.delete_board(board.id)
        if deleted:
            display.success(f"已删除看板 '{deleted.name}'")
            display.info("可使用 'kanban restore --id <ID>' 恢复")


@app.command()
def restore(
    board_id: str = typer.Option(..., "--id", help="要恢复的看板ID"),
):
    """恢复已删除的看板"""
    board = manager.restore_board(board_id)
    if board:
        display.success(f"已恢复看板 '{board.name}'")
    else:
        display.error(f"无法恢复ID为 {board_id} 的看板，可能已超过30天或不存在")


board_app = typer.Typer(help="看板管理")
app.add_typer(board_app, name="board")


@board_app.command("add-column")
def add_column(
    name: str = typer.Argument(..., help="列名称"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
    is_completed: bool = typer.Option(False, "--completed", help="标记为完成列"),
):
    """添加列"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    column = manager.add_column(board.id, name, is_completed=is_completed)
    if column:
        display.success(f"已添加列 '{column.name}'")
        display.info(f"ID: {column.id}")


@board_app.command("rename-column")
def rename_column(
    column_id: str = typer.Argument(..., help="列ID"),
    new_name: str = typer.Argument(..., help="新名称"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """重命名列"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    column = manager.rename_column(board.id, column_id, new_name)
    if column:
        display.success(f"已将列重命名为 '{column.name}'")
    else:
        display.error(f"未找到ID为 {column_id} 的列")


@board_app.command("delete-column")
def delete_column(
    column_id: str = typer.Argument(..., help="列ID"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
    delete_cards: bool = typer.Option(False, "--delete-cards", help="同时删除列中的卡片"),
):
    """删除列"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    column = board.get_column(column_id)
    if not column:
        display.error(f"未找到ID为 {column_id} 的列")
        raise typer.Exit(code=1)
    card_count = len(board.get_cards_in_column(column_id))
    if delete_cards and card_count > 0:
        if not _confirm(f"确定要删除列 '{column.name}' 及其 {card_count} 张卡片吗？"):
            raise typer.Exit(code=0)
    elif card_count > 0:
        if not _confirm(f"列 '{column.name}' 中有 {card_count} 张卡片，这些卡片将变为未分配状态。确定删除吗？"):
            raise typer.Exit(code=0)
    if manager.delete_column(board.id, column_id, delete_cards=delete_cards):
        display.success(f"已删除列 '{column.name}'")


@board_app.command("reorder-columns")
def reorder_columns(
    column_ids: List[str] = typer.Argument(..., help="按顺序排列的列ID列表"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """重新排序列"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    if manager.reorder_columns(board.id, column_ids):
        display.success("已更新列顺序")
        updated_board = manager.get_board(board.id)
        if updated_board:
            display.display_columns(updated_board.columns)


@board_app.command("columns")
def list_columns(
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """列出所有列"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    display.display_columns(board.columns)


card_app = typer.Typer(help="卡片管理")
app.add_typer(card_app, name="card")


@card_app.command("add")
def add_card(
    title: str = typer.Argument(..., help="卡片标题"),
    column_name: str = typer.Argument(..., help="列名称"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
    description: str = typer.Option("", "--desc", help="卡片描述"),
    priority: Priority = typer.Option(Priority.MEDIUM, "--priority", help="优先级"),
    tags: Optional[List[str]] = typer.Option(None, "--tag", help="标签，可多次指定"),
    due_days: Optional[int] = typer.Option(None, "--due-days", help="多少天后截止"),
    due_date: Optional[str] = typer.Option(None, "--due-date", help="截止日期 (YYYY-MM-DD)"),
):
    """添加卡片"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    column = board.get_column_by_name(column_name)
    if not column:
        display.error(f"未找到列 '{column_name}'")
        raise typer.Exit(code=1)
    due = due_date
    if due_days is not None:
        due = (date.today() + timedelta(days=due_days)).isoformat()
    card = manager.add_card(
        board.id, column.id, title,
        description=description,
        priority=priority,
        tags=tags,
        due_date=due,
    )
    if card:
        display.success(f"已添加卡片 '{card.title}'")
        display.info(f"ID: {card.id}")


@card_app.command("show")
def show_card(
    card_id: str = typer.Argument(..., help="卡片ID"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """显示卡片详情"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    card = board.get_card(card_id)
    if not card:
        display.error(f"未找到ID为 {card_id} 的卡片")
        raise typer.Exit(code=1)
    column = board.get_column(card.column_id)
    display.display_card(card, column)


@card_app.command("move")
def move_card(
    card_id: str = typer.Argument(..., help="卡片ID"),
    to_column: str = typer.Argument(..., help="目标列名称"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """移动卡片到另一列"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    card = board.get_card(card_id)
    if not card:
        display.error(f"未找到ID为 {card_id} 的卡片")
        raise typer.Exit(code=1)
    new_column = board.get_column_by_name(to_column)
    if not new_column:
        display.error(f"未找到列 '{to_column}'")
        raise typer.Exit(code=1)
    old_column = board.get_column(card.column_id)
    old_name = old_column.name if old_column else "未分配"
    moved = manager.move_card(board.id, card_id, new_column.id)
    if moved:
        display.success(f"已将卡片 '{card.title}' 从 '{old_name}' 移动到 '{new_column.name}'")


@card_app.command("update")
def update_card(
    card_id: str = typer.Argument(..., help="卡片ID"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
    title: Optional[str] = typer.Option(None, "--title", help="新标题"),
    description: Optional[str] = typer.Option(None, "--desc", help="新描述"),
    priority: Optional[Priority] = typer.Option(None, "--priority", help="新优先级"),
    due_days: Optional[int] = typer.Option(None, "--due-days", help="多少天后截止"),
    due_date: Optional[str] = typer.Option(None, "--due-date", help="截止日期 (YYYY-MM-DD)"),
):
    """更新卡片信息"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    card = board.get_card(card_id)
    if not card:
        display.error(f"未找到ID为 {card_id} 的卡片")
        raise typer.Exit(code=1)
    updates = {}
    if title is not None:
        updates["title"] = title
    if description is not None:
        updates["description"] = description
    if priority is not None:
        updates["priority"] = priority
    due = due_date
    if due_days is not None:
        due = (date.today() + timedelta(days=due_days)).isoformat()
    if due is not None:
        updates["due_date"] = due
    if not updates:
        display.warning("没有提供要更新的字段")
        raise typer.Exit(code=0)
    updated = manager.update_card(board.id, card_id, **updates)
    if updated:
        display.success(f"已更新卡片 '{updated.title}'")


@card_app.command("delete")
def delete_card(
    card_id: str = typer.Argument(..., help="卡片ID"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """删除卡片"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    card = board.get_card(card_id)
    if not card:
        display.error(f"未找到ID为 {card_id} 的卡片")
        raise typer.Exit(code=1)
    if _confirm(f"确定要删除卡片 '{card.title}' 吗？此操作可在30天内撤销。"):
        deleted = manager.delete_card(board.id, card_id)
        if deleted:
            display.success(f"已删除卡片 '{deleted.title}'")
            display.info("可使用 'kanban card restore --id <ID>' 恢复")


@card_app.command("restore")
def restore_card(
    card_id: str = typer.Option(..., "--id", help="要恢复的卡片ID"),
):
    """恢复已删除的卡片"""
    card = manager.restore_card(card_id)
    if card:
        display.success(f"已恢复卡片 '{card.title}'")
    else:
        display.error(f"无法恢复ID为 {card_id} 的卡片，可能已超过30天或不存在")


@card_app.command("tag")
def add_tag(
    card_id: str = typer.Argument(..., help="卡片ID"),
    tag: str = typer.Argument(..., help="标签"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """为卡片添加标签"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    card = manager.add_tag(board.id, card_id, tag)
    if card:
        display.success(f"已为卡片 '{card.title}' 添加标签 [{tag}]")


@card_app.command("untag")
def remove_tag(
    card_id: str = typer.Argument(..., help="卡片ID"),
    tag: str = typer.Argument(..., help="标签"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """移除卡片标签"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    card = manager.remove_tag(board.id, card_id, tag)
    if card:
        display.success(f"已从卡片 '{card.title}' 移除标签 [{tag}]")


search_app = typer.Typer(help="搜索和过滤")
app.add_typer(search_app, name="search")


@search_app.command("cards")
def search_cards(
    query: Optional[str] = typer.Option("", "--query", "-q", help="搜索关键词"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
    tags: Optional[List[str]] = typer.Option(None, "--tag", help="按标签过滤"),
    priority: Optional[Priority] = typer.Option(None, "--priority", help="按优先级过滤"),
    date_from: Optional[str] = typer.Option(None, "--from", help="起始日期 (YYYY-MM-DD)"),
    date_to: Optional[str] = typer.Option(None, "--to", help="结束日期 (YYYY-MM-DD)"),
    column_name: Optional[str] = typer.Option(None, "--column", help="按列过滤"),
    overdue: bool = typer.Option(False, "--overdue", help="只显示过期卡片"),
    details: bool = typer.Option(False, "--details", help="显示详细信息"),
):
    """搜索和过滤卡片"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    column_id = None
    if column_name:
        column = board.get_column_by_name(column_name)
        if column:
            column_id = column.id
    cards = manager.search_cards(
        board.id,
        query=query,
        tags=tags,
        priority=priority,
        date_from=date_from,
        date_to=date_to,
        column_id=column_id,
        overdue_only=overdue,
    )
    display.display_cards(cards, board=board, show_details=details)


@search_app.command("tags")
def list_tags(
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """列出所有标签"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    tags = manager.get_all_tags(board.id)
    if not tags:
        display.info("暂无标签")
        return
    display.info(f"共 {len(tags)} 个标签:")
    for tag in tags:
        typer.echo(f"  [{tag}]")


io_app = typer.Typer(help="导入导出")
app.add_typer(io_app, name="io")


@io_app.command("export")
def export_board(
    export_path: str = typer.Argument(..., help="导出文件路径"),
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """导出版到JSON文件"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    if manager.export_board(board.id, export_path):
        display.success(f"已导出版 '{board.name}' 到 {export_path}")
    else:
        display.error("导出失败")


@io_app.command("import")
def import_board(
    import_path: str = typer.Argument(..., help="导入文件路径"),
    resolution: DuplicateResolution = typer.Option(
        DuplicateResolution.SKIP,
        "--resolve",
        help="重复处理方式: skip(跳过), replace(替换), rename(重命名)",
    ),
):
    """从JSON文件导入看板"""
    board, message = manager.import_board(import_path, resolution=resolution)
    if board:
        display.success(message)
        display.info(f"ID: {board.id}")
    else:
        display.error(message)


stats_app = typer.Typer(help="统计")
app.add_typer(stats_app, name="stats")


@stats_app.command("show")
def show_stats(
    board_id: Optional[str] = typer.Option(None, "--id", help="看板ID"),
    board_name: Optional[str] = typer.Option(None, "--name", help="看板名称"),
):
    """显示看板统计信息"""
    board = _resolve_board(board_id, board_name)
    if not board:
        raise typer.Exit(code=1)
    stats = manager.get_stats(board.id)
    display.display_stats(stats)


def main():
    app()


if __name__ == "__main__":
    main()
