import signal
import sys
import time
from typing import Optional

import click

from .models import init_db, get_session, Feed, Group, Tag, Article, FilterRule
from .scheduler import FeedScheduler
from .report import ReportGenerator
from .emailer import EmailNotifier
from .logger import logger
from .config import settings


@click.group()
@click.version_option()
def cli():
    """RSS Reader - 自动化RSS订阅与内容聚合工具"""
    init_db()


@cli.group()
def feed():
    """订阅源管理"""
    pass


@feed.command("add")
@click.argument("url")
@click.option("--title", "-t", help="订阅源标题")
@click.option("--group", "-g", help="分组名称")
@click.option("--interval", "-i", default=30, help="抓取间隔(分钟)")
@click.option("--timeout", default=30, help="请求超时(秒)")
@click.option("--retries", default=3, help="最大重试次数")
def feed_add(url, title, group, interval, timeout, retries):
    """添加RSS订阅源"""
    session = get_session()
    try:
        existing = session.query(Feed).filter_by(url=url).first()
        if existing:
            click.echo(f"⚠️  订阅源已存在: {existing.title}")
            return

        group_id = None
        if group:
            grp = session.query(Group).filter_by(name=group).first()
            if not grp:
                grp = Group(name=group)
                session.add(grp)
                session.flush()
            group_id = grp.id

        feed_obj = Feed(
            title=title or url,
            url=url,
            group_id=group_id,
            fetch_interval=interval,
            request_timeout=timeout,
            max_retries=retries,
        )
        session.add(feed_obj)
        session.commit()
        click.echo(f"✅ 添加成功: {feed_obj.title} (ID: {feed_obj.id})")
    except Exception as e:
        click.echo(f"❌ 添加失败: {e}", err=True)
    finally:
        session.close()


@feed.command("list")
@click.option("--group", "-g", help="按分组筛选")
def feed_list(group):
    """列出所有订阅源"""
    session = get_session()
    try:
        query = session.query(Feed)
        if group:
            grp = session.query(Group).filter_by(name=group).first()
            if grp:
                query = query.filter_by(group_id=grp.id)

        feeds = query.order_by(Feed.id).all()
        if not feeds:
            click.echo("暂无订阅源")
            return

        click.echo(f"{'ID':<5} {'状态':<6} {'分组':<10} {'间隔':<6} {'标题'}")
        click.echo("-" * 60)
        for f in feeds:
            status = "✅" if f.enabled else "❌"
            group_name = f.group.name if f.group else "-"
            click.echo(f"{f.id:<5} {status:<6} {group_name:<10} {f.fetch_interval:<6} {f.title}")
    finally:
        session.close()


@feed.command("remove")
@click.argument("feed_id", type=int)
def feed_remove(feed_id):
    """删除订阅源"""
    session = get_session()
    try:
        feed_obj = session.get(Feed, feed_id)
        if not feed_obj:
            click.echo(f"❌ 订阅源不存在: {feed_id}")
            return
        session.delete(feed_obj)
        session.commit()
        click.echo(f"✅ 已删除: {feed_obj.title}")
    except Exception as e:
        click.echo(f"❌ 删除失败: {e}", err=True)
    finally:
        session.close()


@feed.command("enable")
@click.argument("feed_id", type=int)
def feed_enable(feed_id):
    """启用订阅源"""
    session = get_session()
    try:
        feed_obj = session.get(Feed, feed_id)
        if feed_obj:
            feed_obj.enabled = True
            session.commit()
            click.echo(f"✅ 已启用: {feed_obj.title}")
    finally:
        session.close()


@feed.command("disable")
@click.argument("feed_id", type=int)
def feed_disable(feed_id):
    """禁用订阅源"""
    session = get_session()
    try:
        feed_obj = session.get(Feed, feed_id)
        if feed_obj:
            feed_obj.enabled = False
            session.commit()
            click.echo(f"✅ 已禁用: {feed_obj.title}")
    finally:
        session.close()


@cli.group()
def group():
    """分组管理"""
    pass


@group.command("add")
@click.argument("name")
@click.option("--description", "-d", help="分组描述")
def group_add(name, description):
    """添加分组"""
    session = get_session()
    try:
        existing = session.query(Group).filter_by(name=name).first()
        if existing:
            click.echo(f"⚠️  分组已存在: {name}")
            return
        grp = Group(name=name, description=description)
        session.add(grp)
        session.commit()
        click.echo(f"✅ 添加分组: {name} (ID: {grp.id})")
    finally:
        session.close()


@group.command("list")
def group_list():
    """列出所有分组"""
    session = get_session()
    try:
        groups = session.query(Group).all()
        if not groups:
            click.echo("暂无分组")
            return
        for g in groups:
            count = len(g.feeds)
            click.echo(f"[{g.id}] {g.name} - {count} 个订阅源")
    finally:
        session.close()


@cli.group()
def filter():
    """过滤规则管理"""
    pass


@filter.command("add")
@click.option("--name", "-n", required=True, help="规则名称")
@click.option("--type", "rule_type", required=True, type=click.Choice(["whitelist", "blacklist", "tag"]), help="规则类型")
@click.option("--pattern", "-p", required=True, help="匹配模式(逗号分隔关键词或正则)")
@click.option("--action", "-a", required=True, type=click.Choice(["ignore", "tag"]), help="匹配后的动作")
@click.option("--regex/--no-regex", default=False, help="是否使用正则")
@click.option("--apply-to", default="title_content", type=click.Choice(["title", "summary", "content", "title_content"]))
@click.option("--tag", "tag_name", help="要添加的标签名(action=tag时)")
@click.option("--priority", default=0, help="优先级")
def filter_add(name, rule_type, pattern, action, regex, apply_to, tag_name, priority):
    """添加过滤规则"""
    from .filters import FilterEngine
    engine = FilterEngine()
    rule_id = engine.add_rule(
        name=name,
        rule_type=rule_type,
        pattern=pattern,
        action=action,
        is_regex=regex,
        apply_to=apply_to,
        tag_name=tag_name,
        priority=priority,
    )
    if rule_id:
        click.echo(f"✅ 规则已添加 (ID: {rule_id})")
    else:
        click.echo("❌ 添加失败", err=True)


@filter.command("list")
def filter_list():
    """列出所有过滤规则"""
    session = get_session()
    try:
        rules = session.query(FilterRule).order_by(FilterRule.priority.desc()).all()
        if not rules:
            click.echo("暂无过滤规则")
            return
        for r in rules:
            tag_name = r.tag.name if r.tag else "-"
            status = "✅" if r.enabled else "❌"
            click.echo(f"[{r.id}] {status} {r.name} | {r.rule_type} -> {r.action} | 标签: {tag_name}")
    finally:
        session.close()


@cli.command("fetch")
@click.option("--feed-id", type=int, help="指定抓取某个源")
@click.option("--all", "fetch_all", is_flag=True, help="抓取所有源")
def cmd_fetch(feed_id, fetch_all):
    """抓取RSS内容"""
    scheduler = FeedScheduler()
    if feed_id:
        session = get_session()
        try:
            feed_obj = session.get(Feed, feed_id)
            if feed_obj:
                scheduler._process_feed(feed_id)
                click.echo("✅ 抓取完成")
        finally:
            session.close()
    elif fetch_all:
        new_articles = scheduler.run_once()
        click.echo(f"✅ 抓取完成, 新增 {len(new_articles)} 篇文章")
    else:
        click.echo("请指定 --feed-id 或 --all")


@cli.command("run")
@click.option("--with-report", is_flag=True, help="同时启用每日报告")
def cmd_run(with_report):
    """启动定时调度服务"""
    scheduler = FeedScheduler()

    def signal_handler(signum, frame):
        click.echo("\n正在停止服务...")
        scheduler.stop()
        sys.exit(0)

    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    scheduler.start()

    if with_report:
        emailer = EmailNotifier()
        scheduler.add_report_job(lambda: emailer.send_daily_report(period="daily"))
        click.echo("📧 已启用每日报告推送")

    click.echo(f"🚀 RSS Reader 已启动, 数据目录: {settings.data_dir}")
    click.echo("按 Ctrl+C 停止服务")

    while True:
        time.sleep(1)


@cli.command("report")
@click.option("--period", default="daily", type=click.Choice(["daily", "weekly"]))
@click.option("--format", "fmt", default="markdown", type=click.Choice(["markdown", "html"]))
@click.option("--output", "-o", help="输出目录")
@click.option("--send-email", is_flag=True, help="生成后发送邮件")
@click.option("--group-id", type=int, help="指定分组ID")
def cmd_report(period, fmt, output, send_email, group_id):
    """生成摘要报告"""
    reporter = ReportGenerator()
    path = reporter.save_report(period=period, fmt=fmt, output_dir=output)
    click.echo(f"✅ 报告已生成: {path}")

    if send_email:
        emailer = EmailNotifier()
        ok = emailer.send_daily_report(period=period, group_id=group_id)
        if ok:
            click.echo("📧 邮件已发送")
        else:
            click.echo("⚠️  邮件发送失败(可能未配置SMTP)")


@cli.command("stats")
def cmd_stats():
    """查看统计信息"""
    session = get_session()
    try:
        total_feeds = session.query(Feed).count()
        enabled_feeds = session.query(Feed).filter_by(enabled=True).count()
        total_articles = session.query(Article).count()
        dup_articles = session.query(Article).filter_by(is_duplicate=True).count()
        ignored_articles = session.query(Article).filter_by(is_ignored=True).count()
        total_tags = session.query(Tag).count()

        click.echo("📊 RSS Reader 统计")
        click.echo("=" * 30)
        click.echo(f"订阅源: {enabled_feeds}/{total_feeds} 已启用")
        click.echo(f"文章总数: {total_articles}")
        click.echo(f"重复文章: {dup_articles}")
        click.echo(f"忽略文章: {ignored_articles}")
        click.echo(f"标签数量: {total_tags}")
    finally:
        session.close()


@cli.command("init-demo")
def cmd_init_demo():
    """初始化示例数据"""
    session = get_session()
    try:
        groups = [
            Group(name="科技新闻", description="科技行业动态"),
            Group(name="技术博客", description="程序员技术博客"),
            Group(name="开源项目", description="开源项目更新"),
        ]
        for g in groups:
            if not session.query(Group).filter_by(name=g.name).first():
                session.add(g)
        session.flush()

        tech_group = session.query(Group).filter_by(name="科技新闻").first()
        blog_group = session.query(Group).filter_by(name="技术博客").first()

        demo_feeds = [
            Feed(title="Solidot", url="https://www.solidot.org/index.rss", group_id=tech_group.id, fetch_interval=30),
            Feed(title="Hacker News", url="https://news.ycombinator.com/rss", group_id=tech_group.id, fetch_interval=30),
            Feed(title="阮一峰的网络日志", url="https://www.ruanyifeng.com/blog/atom.xml", group_id=blog_group.id, fetch_interval=1440),
            Feed(title="酷壳", url="https://coolshell.cn/feed", group_id=blog_group.id, fetch_interval=1440),
        ]

        for f in demo_feeds:
            if not session.query(Feed).filter_by(url=f.url).first():
                session.add(f)

        session.commit()
        click.echo("✅ 示例数据已初始化")
        click.echo("已添加: 3个分组, 4个订阅源")
    except Exception as e:
        click.echo(f"❌ 初始化失败: {e}", err=True)
        session.rollback()
    finally:
        session.close()


if __name__ == "__main__":
    cli()
