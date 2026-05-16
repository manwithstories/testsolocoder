from rich.console import Console
from rich.table import Table

from track.models.transaction import Transaction, TransactionType

console = Console()


def print_transactions(transactions: list[Transaction]):
    if not transactions:
        console.print("[yellow]没有找到记录[/yellow]")
        return

    table = Table(title="交易记录")
    table.add_column("ID", style="cyan", justify="right")
    table.add_column("类型", style="magenta")
    table.add_column("分类", style="green")
    table.add_column("金额", style="yellow", justify="right")
    table.add_column("备注", style="blue")
    table.add_column("时间", style="dim")

    for t in transactions:
        type_color = "red" if t.type == TransactionType.EXPENSE else "green"
        type_label = "支出" if t.type == TransactionType.EXPENSE else "收入"
        amount_sign = "-" if t.type == TransactionType.EXPENSE else "+"
        table.add_row(
            str(t.id),
            f"[{type_color}]{type_label}[/{type_color}]",
            t.category,
            f"{amount_sign}¥{t.amount:.2f}",
            t.note,
            t.created_at.strftime("%Y-%m-%d %H:%M"),
        )

    console.print(table)
