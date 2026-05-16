from rich.console import Console
from rich.text import Text

console = Console()


def print_bar_chart(data: dict[str, float], title: str, max_width: int = 30):
    if not data:
        console.print(f"[yellow]{title}: 没有数据[/yellow]")
        return

    max_value = max(data.values()) if data else 1
    total = sum(data.values())

    console.print(f"\n[bold]{title}[/bold]")
    console.print(f"[dim]总计: ¥{total:.2f}[/dim]\n")

    for category, amount in sorted(data.items(), key=lambda x: x[1], reverse=True):
        bar_length = int((amount / max_value) * max_width) if max_value > 0 else 0
        bar = "█" * bar_length
        percentage = (amount / total * 100) if total > 0 else 0

        line = Text()
        line.append(f"{category:<10}", style="cyan")
        line.append(f" ¥{amount:>8.2f} ", style="yellow")
        line.append(f"{bar}", style="green")
        line.append(f" {percentage:5.1f}%", style="dim")

        console.print(line)

    console.print()
