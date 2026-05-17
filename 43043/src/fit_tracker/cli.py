import click
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from datetime import datetime
from .database import init_db, SessionLocal
from . import exercises as exercise_module
from . import training as training_module
from . import plans as plan_module
from . import statistics as stats_module
from . import export as export_module

console = Console()


def get_db():
    db = SessionLocal()
    try:
        return db
    finally:
        pass


@click.group()
@click.version_option()
def cli():
    """健身训练记录工具"""
    init_db()


@cli.group()
def exercise():
    """动作库管理"""
    pass


@exercise.command("add")
@click.option("--name", "-n", required=True, help="动作名称")
@click.option("--category", "-c", required=True, help="分类：力量/有氧/柔韧性/其他")
@click.option("--sets", "-s", default=3, help="默认组数")
@click.option("--reps", "-r", default=10, help="默认次数")
@click.option("--duration", "-d", default=0, help="默认时长(秒)")
def add_exercise(name, category, sets, reps, duration):
    """添加新动作"""
    db = get_db()
    try:
        ex = exercise_module.add_exercise(db, name, category, sets, reps, duration)
        console.print(f"[green]✓ 动作 '{ex.name}' 添加成功！[/green]")
        console.print(f"  ID: {ex.id}")
        console.print(f"  分类: {ex.category}")
        console.print(f"  默认: {ex.default_sets}组 × {ex.default_reps}次" + (f" / {ex.default_duration}秒" if ex.default_duration > 0 else ""))
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@exercise.command("list")
@click.option("--category", "-c", help="按分类筛选")
@click.option("--all", "-a", is_flag=True, help="显示已归档的动作")
def list_exercises(category, all):
    """列出所有动作"""
    db = get_db()
    exercises = exercise_module.list_exercises(db, category, include_archived=all)

    if not exercises:
        console.print("[yellow]暂无动作数据[/yellow]")
        return

    table = Table(title="动作库")
    table.add_column("ID", style="cyan", no_wrap=True)
    table.add_column("名称", style="magenta")
    table.add_column("分类", style="green")
    table.add_column("默认组数", justify="right")
    table.add_column("默认次数", justify="right")
    table.add_column("默认时长", justify="right")
    table.add_column("状态", style="yellow")

    for ex in exercises:
        status = "[dim]已归档[/dim]" if ex.is_archived else "[green]正常[/green]"
        table.add_row(
            str(ex.id),
            ex.name,
            ex.category,
            str(ex.default_sets),
            str(ex.default_reps),
            str(ex.default_duration) + "s" if ex.default_duration > 0 else "-",
            status
        )

    console.print(table)


@exercise.command("update")
@click.argument("exercise_id", type=int)
@click.option("--name", "-n", help="新名称")
@click.option("--category", "-c", help="新分类")
@click.option("--sets", "-s", type=int, help="新默认组数")
@click.option("--reps", "-r", type=int, help="新默认次数")
@click.option("--duration", "-d", type=int, help="新默认时长(秒)")
def update_exercise(exercise_id, name, category, sets, reps, duration):
    """更新动作信息"""
    db = get_db()
    try:
        ex = exercise_module.update_exercise(db, exercise_id, name, category, sets, reps, duration)
        if ex:
            console.print(f"[green]✓ 动作 '{ex.name}' 更新成功！[/green]")
        else:
            console.print(f"[red]✗ 动作 ID {exercise_id} 不存在[/red]")
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@exercise.command("delete")
@click.argument("exercise_id", type=int)
def delete_exercise(exercise_id):
    """删除动作（有关联记录则归档）"""
    db = get_db()
    ex = exercise_module.get_exercise(db, exercise_id)
    if not ex:
        console.print(f"[red]✗ 动作 ID {exercise_id} 不存在[/red]")
        return

    if click.confirm(f"确定要删除动作 '{ex.name}' 吗？"):
        exercise_module.delete_exercise(db, exercise_id)
        ex = exercise_module.get_exercise(db, exercise_id)
        if ex and ex.is_archived:
            console.print(f"[yellow]⚠ 动作 '{ex.name}' 有关联记录，已标记为归档[/yellow]")
        else:
            console.print(f"[green]✓ 动作已删除[/green]")


@exercise.command("categories")
def list_categories():
    """显示可用分类"""
    categories = exercise_module.get_categories()
    console.print("可用分类：")
    for c in categories:
        console.print(f"  • {c}")


@cli.group()
def train():
    """训练记录"""
    pass


@train.command("start")
@click.option("--notes", "-n", help="训练备注")
def start_training(notes):
    """开始新训练"""
    db = get_db()
    try:
        session = training_module.start_training(db, notes)
        console.print(f"[green]✓ 训练开始！[/green] 训练ID: {session.id}")
        if notes:
            console.print(f"  备注: {notes}")
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@train.command("add-set")
@click.option("--exercise", "-e", required=True, help="动作名称或ID")
@click.option("--reps", "-r", type=int, required=True, help="次数")
@click.option("--weight", "-w", type=float, default=0.0, help="重量(kg)")
@click.option("--duration", "-d", type=int, default=0, help="时长(秒)")
@click.option("--notes", "-n", help="组备注")
def add_set(exercise, reps, weight, duration, notes):
    """向当前训练添加一组"""
    db = get_db()
    active_session = training_module.get_active_session(db)
    if not active_session:
        console.print("[red]✗ 没有进行中的训练，请先使用 'fit train start' 开始训练[/red]")
        return

    ex = None
    if exercise.isdigit():
        ex = exercise_module.get_exercise(db, int(exercise))
    else:
        ex = exercise_module.get_exercise_by_name(db, exercise)

    if not ex:
        console.print(f"[red]✗ 找不到动作 '{exercise}'[/red]")
        return

    try:
        training_set = training_module.add_set(db, active_session.id, ex.id, reps, weight, duration, notes)
        volume = weight * reps if weight > 0 and reps > 0 else duration
        console.print(f"[green]✓ 已添加:[/green] {ex.name} - {training_set.set_number}组: {reps}次" +
                      (f" × {weight}kg" if weight > 0 else "") +
                      (f" / {duration}秒" if duration > 0 else "") +
                      f" [cyan](容量: {volume})[/cyan]")
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@train.command("remove-set")
@click.argument("set_id", type=int)
def remove_set(set_id):
    """删除一组"""
    db = get_db()
    try:
        if training_module.remove_set(db, set_id):
            console.print(f"[green]✓ 组 {set_id} 已删除[/green]")
        else:
            console.print(f"[red]✗ 组 {set_id} 不存在[/red]")
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@train.command("status")
def training_status():
    """查看当前训练状态"""
    db = get_db()
    active_session = training_module.get_active_session(db)
    if not active_session:
        console.print("[yellow]当前没有进行中的训练[/yellow]")
        return

    exercises = training_module.get_session_exercises(db, active_session.id)

    console.print(Panel.fit(
        f"[bold]训练ID:[/bold] {active_session.id}\n"
        f"[bold]开始时间:[/bold] {active_session.start_time.strftime('%Y-%m-%d %H:%M:%S')}\n"
        f"[bold]当前容量:[/bold] [cyan]{active_session.total_volume}[/cyan]\n"
        f"[bold]已记录组数:[/bold] {len(active_session.sets)}",
        title="训练进行中",
        border_style="green"
    ))

    if exercises:
        table = Table(title="已记录动作")
        table.add_column("动作", style="magenta")
        table.add_column("组", justify="right")
        table.add_column("详情", style="cyan")

        for ex_id, sets in exercises.items():
            ex = exercise_module.get_exercise(db, ex_id)
            ex_name = ex.name if ex else f"ID:{ex_id}"
            set_details = []
            for s in sets:
                detail = f"{s.reps}次"
                if s.weight > 0:
                    detail += f"×{s.weight}kg"
                if s.duration > 0:
                    detail += f"/{s.duration}s"
                set_details.append(detail)
            table.add_row(ex_name, str(len(sets)), ", ".join(set_details))

        console.print(table)


@train.command("finish")
def finish_training():
    """结束当前训练"""
    db = get_db()
    active_session = training_module.get_active_session(db)
    if not active_session:
        console.print("[yellow]当前没有进行中的训练[/yellow]")
        return

    if click.confirm(f"确定要结束训练吗？总容量: {active_session.total_volume}"):
        session = training_module.finish_training(db)
        duration = (session.end_time - session.start_time).total_seconds() / 60 if session.end_time else 0
        console.print(f"[green]✓ 训练完成！[/green]")
        console.print(f"  总容量: {session.total_volume}")
        console.print(f"  总时长: {duration:.1f}分钟")
        console.print(f"  总组数: {len(session.sets)}")


@train.command("cancel")
def cancel_training():
    """取消当前训练（丢弃所有记录）"""
    db = get_db()
    active_session = training_module.get_active_session(db)
    if not active_session:
        console.print("[yellow]当前没有进行中的训练[/yellow]")
        return

    if click.confirm("[red]确定要取消本次训练吗？所有记录将被丢弃！[/red]", abort=True):
        training_module.cancel_training(db)
        console.print("[yellow]训练已取消，记录已丢弃[/yellow]")


@train.command("history")
@click.option("--limit", "-l", default=10, help="显示最近N次训练")
def training_history(limit):
    """查看训练历史"""
    db = get_db()
    sessions = training_module.list_sessions(db, limit=limit, status="completed")

    if not sessions:
        console.print("[yellow]暂无训练历史[/yellow]")
        return

    table = Table(title=f"最近 {len(sessions)} 次训练")
    table.add_column("ID", style="cyan")
    table.add_column("开始时间", style="blue")
    table.add_column("时长(分)", justify="right")
    table.add_column("总容量", justify="right", style="green")
    table.add_column("组数", justify="right")

    for s in sessions:
        duration = (s.end_time - s.start_time).total_seconds() / 60 if s.end_time else 0
        table.add_row(
            str(s.id),
            s.start_time.strftime("%Y-%m-%d %H:%M"),
            f"{duration:.1f}",
            str(s.total_volume),
            str(len(s.sets))
        )

    console.print(table)


@cli.group()
def plan():
    """训练计划管理"""
    pass


@plan.command("create")
@click.option("--name", "-n", required=True, help="计划名称")
@click.option("--description", "-d", help="计划描述")
def create_plan(name, description):
    """创建训练计划模板"""
    db = get_db()
    try:
        p = plan_module.create_plan(db, name, description)
        console.print(f"[green]✓ 计划 '{p.name}' 创建成功！[/green] ID: {p.id}")
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@plan.command("list")
def list_plans():
    """列出所有计划"""
    db = get_db()
    plans = plan_module.list_plans(db)

    if not plans:
        console.print("[yellow]暂无计划模板[/yellow]")
        return

    table = Table(title="训练计划模板")
    table.add_column("ID", style="cyan")
    table.add_column("名称", style="magenta")
    table.add_column("动作数", justify="right")
    table.add_column("描述", style="blue")

    for p in plans:
        table.add_row(
            str(p.id),
            p.name,
            str(len(p.exercises)),
            p.description or "-"
        )

    console.print(table)


@plan.command("add-exercise")
@click.option("--plan-id", "-p", type=int, required=True, help="计划ID")
@click.option("--exercise", "-e", required=True, help="动作名称或ID")
@click.option("--sets", "-s", type=int, help="组数（使用动作默认值）")
@click.option("--reps", "-r", type=int, help="次数（使用动作默认值）")
@click.option("--duration", "-d", type=int, help="时长(秒)")
@click.option("--notes", "-n", help="备注")
def add_exercise_to_plan(plan_id, exercise, sets, reps, duration, notes):
    """向计划添加动作"""
    db = get_db()

    ex = None
    if exercise.isdigit():
        ex = exercise_module.get_exercise(db, int(exercise))
    else:
        ex = exercise_module.get_exercise_by_name(db, exercise)

    if not ex:
        console.print(f"[red]✗ 找不到动作 '{exercise}'[/red]")
        return

    if sets is None:
        sets = ex.default_sets
    if reps is None:
        reps = ex.default_reps
    if duration is None:
        duration = ex.default_duration

    try:
        pe = plan_module.add_exercise_to_plan(db, plan_id, ex.id, sets, reps, duration, notes)
        console.print(f"[green]✓ 已添加:[/green] {ex.name} - {sets}组 × {reps}次" +
                      (f" / {duration}秒" if duration > 0 else ""))
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@plan.command("remove-exercise")
@click.option("--plan-id", "-p", type=int, required=True, help="计划ID")
@click.option("--exercise-id", "-e", type=int, required=True, help="动作ID")
def remove_exercise_from_plan(plan_id, exercise_id):
    """从计划移除动作"""
    db = get_db()
    if plan_module.remove_exercise_from_plan(db, plan_id, exercise_id):
        console.print("[green]✓ 动作已从计划中移除[/green]")
    else:
        console.print("[red]✗ 未找到该计划或动作[/red]")


@plan.command("show")
@click.argument("plan_id", type=int)
def show_plan(plan_id):
    """显示计划详情"""
    db = get_db()
    p = plan_module.get_plan(db, plan_id)
    if not p:
        console.print(f"[red]✗ 计划 ID {plan_id} 不存在[/red]")
        return

    exercises_with_status = plan_module.get_plan_exercises_with_status(db, plan_id)
    has_archived = any(archived for _, archived in exercises_with_status)

    console.print(Panel.fit(
        f"[bold]名称:[/bold] {p.name}\n"
        f"[bold]描述:[/bold] {p.description or '-'}",
        title=f"计划详情 (ID: {p.id})",
        border_style="blue"
    ))

    if has_archived:
        console.print("[yellow]⚠ 该计划包含已归档的动作[/yellow]")

    if exercises_with_status:
        table = Table(title="计划动作")
        table.add_column("顺序", justify="right")
        table.add_column("动作名称", style="magenta")
        table.add_column("组数", justify="right")
        table.add_column("次数", justify="right")
        table.add_column("时长", justify="right")
        table.add_column("状态", style="yellow")

        for idx, (pe, archived) in enumerate(exercises_with_status, 1):
            ex = exercise_module.get_exercise(db, pe.exercise_id)
            ex_name = ex.name if ex else f"ID:{pe.exercise_id}"
            status = "[red]已归档[/red]" if archived else "[green]正常[/green]"
            table.add_row(
                str(idx),
                ex_name,
                str(pe.sets),
                str(pe.reps),
                str(pe.duration) + "s" if pe.duration > 0 else "-",
                status
            )

        console.print(table)
    else:
        console.print("[yellow]该计划暂无动作[/yellow]")


@plan.command("delete")
@click.argument("plan_id", type=int)
def delete_plan(plan_id):
    """删除计划"""
    db = get_db()
    p = plan_module.get_plan(db, plan_id)
    if not p:
        console.print(f"[red]✗ 计划 ID {plan_id} 不存在[/red]")
        return

    if click.confirm(f"确定要删除计划 '{p.name}' 吗？"):
        plan_module.delete_plan(db, plan_id)
        console.print("[green]✓ 计划已删除[/green]")


@plan.command("schedule")
@click.option("--plan-id", "-p", type=int, required=True, help="计划ID")
@click.option("--day", "-d", type=int, required=True, help="周几：0=周一, 6=周日")
def schedule_plan(plan_id, day):
    """安排计划到每周具体日期"""
    db = get_db()
    try:
        sp = plan_module.schedule_plan(db, plan_id, day)
        console.print(f"[green]✓ 计划已安排在{plan_module.get_weekday_name(day)}[/green]")
    except ValueError as e:
        console.print(f"[red]✗ {e}[/red]")


@plan.command("unschedule")
@click.argument("scheduled_id", type=int)
def unschedule_plan(scheduled_id):
    """取消计划安排"""
    db = get_db()
    if plan_module.unschedule_plan(db, scheduled_id):
        console.print("[green]✓ 已取消安排[/green]")
    else:
        console.print(f"[red]✗ 安排 ID {scheduled_id} 不存在[/red]")


@plan.command("week")
def show_week_schedule():
    """查看本周训练计划"""
    db = get_db()
    scheduled = plan_module.get_scheduled_plans(db)

    if not scheduled:
        console.print("[yellow]本周暂无安排的训练计划[/yellow]")
        return

    table = Table(title="本周训练计划")
    table.add_column("安排ID", style="cyan")
    table.add_column("周几", style="green")
    table.add_column("计划名称", style="magenta")
    table.add_column("动作数", justify="right")

    for sp in scheduled:
        table.add_row(
            str(sp.id),
            plan_module.get_weekday_name(sp.day_of_week),
            sp.plan.name,
            str(len(sp.plan.exercises))
        )

    console.print(table)


@cli.group()
def stats():
    """统计进度"""
    pass


@stats.command("overview")
def stats_overview():
    """总体统计概览"""
    db = get_db()
    total = stats_module.get_total_stats(db)

    console.print(Panel.fit(
        f"[bold]总训练次数:[/bold] {total['total_sessions']}\n"
        f"[bold]总训练容量:[/bold] [cyan]{total['total_volume']}[/cyan]\n"
        f"[bold]总记录组数:[/bold] {total['total_sets']}\n"
        f"[bold]连续训练天数:[/bold] [green]{total['streak_days']} 天[/green]",
        title="训练统计概览",
        border_style="green"
    ))


@stats.command("weekly")
@click.option("--weeks", "-w", type=int, default=4, help="显示最近N周")
def stats_weekly(weeks):
    """每周训练统计"""
    db = get_db()
    weekly = stats_module.get_weekly_stats(db, weeks)

    table = Table(title=f"最近 {weeks} 周训练统计")
    table.add_column("周", style="blue")
    table.add_column("训练次数", justify="right")
    table.add_column("总容量", justify="right", style="green")

    max_volume = max((w["total_volume"] for w in weekly), default=1)

    for w in weekly:
        bar_length = int(w["total_volume"] / max_volume * 20) if max_volume > 0 else 0
        bar = "█" * bar_length + "░" * (20 - bar_length)
        table.add_row(
            w["label"],
            str(w["session_count"]),
            f"{w['total_volume']} {bar}"
        )

    console.print(table)


@stats.command("monthly")
@click.option("--months", "-m", type=int, default=6, help="显示最近N月")
def stats_monthly(months):
    """每月训练统计"""
    db = get_db()
    monthly = stats_module.get_monthly_stats(db, months)

    table = Table(title=f"最近 {months} 月训练统计")
    table.add_column("月份", style="blue")
    table.add_column("训练次数", justify="right")
    table.add_column("总容量", justify="right", style="green")

    max_volume = max((m["total_volume"] for m in monthly), default=1)

    for m in monthly:
        bar_length = int(m["total_volume"] / max_volume * 20) if max_volume > 0 else 0
        bar = "█" * bar_length + "░" * (20 - bar_length)
        table.add_row(
            m["label"],
            str(m["session_count"]),
            f"{m['total_volume']} {bar}"
        )

    console.print(table)


@stats.command("pr")
@click.option("--exercise", "-e", help="动作名称或ID（不指定则显示所有）")
def stats_pr(exercise):
    """个人最佳记录"""
    db = get_db()

    if exercise:
        ex = None
        if exercise.isdigit():
            ex = exercise_module.get_exercise(db, int(exercise))
        else:
            ex = exercise_module.get_exercise_by_name(db, exercise)

        if not ex:
            console.print(f"[red]✗ 找不到动作 '{exercise}'[/red]")
            return

        pr = stats_module.get_exercise_pr(db, ex.id)
        if pr:
            console.print(Panel.fit(
                f"[bold]动作:[/bold] {ex.name}\n"
                f"[bold]最大重量:[/bold] [green]{pr['max_weight']} kg[/green]\n"
                f"[bold]最大容量:[/bold] [cyan]{pr['max_volume']}[/cyan]\n"
                f"[bold]总训练组数:[/bold] {pr['total_sets']}",
                title=f"个人最佳 - {ex.name}",
                border_style="yellow"
            ))
        else:
            console.print(f"[yellow]动作 '{ex.name}' 暂无训练记录[/yellow]")
    else:
        prs = stats_module.get_all_prs(db)
        if not prs:
            console.print("[yellow]暂无个人最佳记录[/yellow]")
            return

        table = Table(title="个人最佳记录")
        table.add_column("动作", style="magenta")
        table.add_column("分类", style="green")
        table.add_column("最大重量(kg)", justify="right")
        table.add_column("最大容量", justify="right", style="cyan")

        for pr in prs:
            table.add_row(
                pr["exercise"].name,
                pr["exercise"].category,
                str(pr["max_weight"]),
                str(pr["max_volume"])
            )

        console.print(table)


@stats.command("streak")
def stats_streak():
    """连续训练天数"""
    db = get_db()
    streak = stats_module.get_streak_days(db)
    console.print(f"[bold green]🔥 连续训练 {streak} 天[/bold green]")


@cli.group()
def export():
    """数据导出"""
    pass


@export.command("training")
@click.option("--output", "-o", default="training_history.csv", help="输出文件路径")
def export_training(output):
    """导出训练历史为CSV"""
    db = get_db()
    if export_module.export_training_history(db, output):
        console.print(f"[green]✓ 训练历史已导出到 {output}[/green]")
    else:
        console.print(f"[red]✗ 导出失败[/red]")


@export.command("exercises")
@click.option("--output", "-o", default="exercises.csv", help="输出文件路径")
def export_exercises(output):
    """导出动作库为CSV"""
    db = get_db()
    if export_module.export_exercises(db, output):
        console.print(f"[green]✓ 动作库已导出到 {output}[/green]")
    else:
        console.print(f"[red]✗ 导出失败[/red]")


@export.command("plans")
@click.option("--output", "-o", default="plans.csv", help="输出文件路径")
def export_plans(output):
    """导出训练计划为CSV"""
    db = get_db()
    if export_module.export_plans(db, output):
        console.print(f"[green]✓ 训练计划已导出到 {output}[/green]")
    else:
        console.print(f"[red]✗ 导出失败[/red]")


@export.command("all")
@click.option("--output-dir", "-o", default="fit_export", help="输出目录")
def export_all(output_dir):
    """导出所有数据"""
    db = get_db()
    results = export_module.export_all(db, output_dir)

    for name, success in results.items():
        if success:
            console.print(f"[green]✓ {name} 已导出[/green]")
        else:
            console.print(f"[red]✗ {name} 导出失败[/red]")


if __name__ == "__main__":
    cli()
