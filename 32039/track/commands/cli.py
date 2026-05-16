from datetime import datetime

import typer
from rich.console import Console

from track.models.transaction import Transaction, TransactionType
from track.storage.database import Database
from track.utils.chart import print_bar_chart

console = Console()
app = typer.Typer()


@app.command()
def add(
    category: str = typer.Option(..., "-c", "--category", help="分类，如 餐饮"),
    amount: float = typer.Option(..., "-a", "--amount", help="金额"),
    note: str = typer.Option("", "-n", "--note", help="备注"),
    income: bool = typer.Option(False, "-i", "--income", help="标记为收入"),
):
    """添加一笔支出或收入"""
    db = Database()
    t_type = TransactionType.INCOME if income else TransactionType.EXPENSE
    transaction = Transaction.create(t_type, category, amount, note)
    transaction_id = db.add_transaction(transaction)

    type_label = "收入" if income else "支出"
    type_color = "green" if income else "red"

    console.print(f"[{type_color}]✓ {type_label}已添加[/{type_color}]")
    console.print(f"  ID: {transaction_id}")
    console.print(f"  分类: {category}")
    console.print(f"  金额: ¥{amount:.2f}")
    if note:
        console.print(f"  备注: {note}")

    current_month = datetime.now().strftime("%Y-%m")
    budget = db.get_budget(category, current_month)
    if budget and not income:
        summary = db.get_monthly_summary(current_month)
        category_spent = summary["expenses"].get(category, 0)
        remaining = budget.amount - category_spent

        if remaining < 0:
            console.print(f"\n[bold red]⚠ 预算超支! 该分类本月已超出 ¥{abs(remaining):.2f}[/bold red]")
        elif remaining < budget.amount * 0.2:
            console.print(f"\n[bold yellow]⚠ 预算警告: 该分类本月剩余 ¥{remaining:.2f}[/bold yellow]")
        else:
            console.print(f"\n[dim]该分类本月剩余预算: ¥{remaining:.2f}[/dim]")


@app.command("add-income")
def add_income(
    category: str = typer.Option(..., "-c", "--category", help="分类，如 工资"),
    amount: float = typer.Option(..., "-a", "--amount", help="金额"),
    note: str = typer.Option("", "-n", "--note", help="备注"),
):
    """快速添加一笔收入"""
    add(category=category, amount=amount, note=note, income=True)


@app.command("list")
def list_transactions(
    month: str = typer.Option(None, "-m", "--month", help="月份，格式 YYYY-MM"),
    category: str = typer.Option(None, "-c", "--category", help="分类筛选"),
    limit: int = typer.Option(20, "-l", "--limit", help="显示条数"),
):
    """列出交易记录"""
    from track.utils.formatting import print_transactions

    db = Database()
    if not month:
        month = datetime.now().strftime("%Y-%m")

    transactions = db.get_transactions(month=month, category=category)
    print_transactions(transactions[:limit])


@app.command()
def summary(
    month: str = typer.Option(None, "-m", "--month", help="月份，格式 YYYY-MM"),
):
    """按月汇总，显示分类柱状图"""
    db = Database()
    if not month:
        month = datetime.now().strftime("%Y-%m")

    summary_data = db.get_monthly_summary(month)
    expenses = summary_data["expenses"]
    incomes = summary_data["incomes"]

    total_expense = sum(expenses.values())
    total_income = sum(incomes.values())
    net = total_income - total_expense

    console.print(f"\n[bold cyan]===== {month} 月度汇总 =====[/bold cyan]\n")
    console.print(f"总收入: [green]¥{total_income:.2f}[/green]")
    console.print(f"总支出: [red]¥{total_expense:.2f}[/red]")

    net_color = "green" if net >= 0 else "red"
    console.print(f"结余: [{net_color}]¥{net:.2f}[/{net_color}]\n")

    if expenses:
        print_bar_chart(expenses, "支出分布")

    if incomes:
        print_bar_chart(incomes, "收入分布")

    budgets = db.get_budgets(month)
    if budgets:
        console.print("\n[bold yellow]预算执行情况:[/bold yellow]")
        for budget in budgets:
            spent = expenses.get(budget.category, 0)
            remaining = budget.amount - spent
            percentage = (spent / budget.amount * 100) if budget.amount > 0 else 0

            if remaining < 0:
                status = f"[bold red]超支 ¥{abs(remaining):.2f} ({percentage:.0f}%)[/bold red]"
            elif percentage >= 80:
                status = f"[bold yellow]剩余 ¥{remaining:.2f} ({percentage:.0f}%)[/bold yellow]"
            else:
                status = f"[green]剩余 ¥{remaining:.2f} ({percentage:.0f}%)[/green]"

            console.print(f"  {budget.category}: ¥{spent:.2f} / ¥{budget.amount:.2f} - {status}")


@app.command()
def export(
    output: str = typer.Option("transactions.csv", "-o", "--output", help="输出文件路径"),
    month: str = typer.Option(None, "-m", "--month", help="月份，格式 YYYY-MM"),
    category: str = typer.Option(None, "-c", "--category", help="分类筛选"),
):
    """导出交易记录为CSV"""
    from pathlib import Path

    from track.utils.exporter import export_transactions_to_csv

    db = Database()
    transactions = db.get_transactions(month=month, category=category)

    output_path = Path(output)
    export_transactions_to_csv(transactions, output_path)

    console.print(f"[green]✓ 已导出 {len(transactions)} 条记录到 {output_path}[/green]")


@app.command()
def budget(
    category: str = typer.Option(..., "-c", "--category", help="分类"),
    amount: float = typer.Option(..., "-a", "--amount", help="预算金额"),
    month: str = typer.Option(None, "-m", "--month", help="月份，格式 YYYY-MM，默认为当月"),
):
    """设置分类月度预算"""
    from track.models.transaction import Budget

    db = Database()
    if not month:
        month = datetime.now().strftime("%Y-%m")

    budget = Budget.create(category, amount, month)
    db.set_budget(budget)

    console.print(f"[green]✓ 已设置 {month} {category} 预算: ¥{amount:.2f}[/green]")


@app.command("budget-list")
def budget_list(
    month: str = typer.Option(None, "-m", "--month", help="月份，格式 YYYY-MM，默认为当月"),
):
    """查看预算设置"""
    db = Database()
    if not month:
        month = datetime.now().strftime("%Y-%m")

    budgets = db.get_budgets(month)
    if not budgets:
        console.print(f"[yellow]{month} 暂无预算设置[/yellow]")
        return

    from rich.table import Table

    table = Table(title=f"{month} 预算设置")
    table.add_column("分类", style="cyan")
    table.add_column("预算金额", style="yellow", justify="right")

    for budget in budgets:
        table.add_row(budget.category, f"¥{budget.amount:.2f}")

    console.print(table)


@app.command()
def delete(
    transaction_id: int = typer.Argument(..., help="交易记录ID"),
):
    """删除交易记录"""
    db = Database()
    success = db.delete_transaction(transaction_id)

    if success:
        console.print(f"[green]✓ 已删除记录 ID: {transaction_id}[/green]")
    else:
        console.print(f"[red]✗ 未找到记录 ID: {transaction_id}[/red]")
