from datetime import datetime, timedelta
from typing import List, Dict, Optional
from collections import defaultdict

from .config import settings
from .logger import logger
from .models import Article, Feed, Group, Tag, get_session


REPORT_TEMPLATES = {
    "html": """<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RSS 摘要报告 - {{ date_range }}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; max-width: 900px; margin: 0 auto; padding: 20px; color: #333; }
        h1 { color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 10px; }
        h2 { color: #34495e; margin-top: 30px; background: #ecf0f1; padding: 8px 15px; border-radius: 4px; }
        h3 { color: #2980b9; margin-top: 20px; }
        .article { border-left: 4px solid #3498db; padding: 15px; margin: 15px 0; background: #f8f9fa; border-radius: 0 6px 6px 0; }
        .article-title { font-size: 18px; margin: 0 0 8px 0; }
        .article-title a { color: #2c3e50; text-decoration: none; }
        .article-title a:hover { color: #3498db; }
        .article-meta { color: #7f8c8d; font-size: 13px; margin-bottom: 8px; }
        .article-summary { color: #555; line-height: 1.6; }
        .tag { display: inline-block; background: #3498db; color: white; padding: 2px 8px; border-radius: 10px; font-size: 12px; margin-right: 5px; }
        .stats { background: #e8f4fd; padding: 15px; border-radius: 6px; margin-bottom: 20px; }
        .stats-item { display: inline-block; margin-right: 30px; font-size: 14px; }
        .stats-value { font-weight: bold; color: #2980b9; }
        .group-section { margin-top: 30px; }
        .duplicate { opacity: 0.6; }
    </style>
</head>
<body>
    <h1>📰 RSS 摘要报告</h1>
    <p><strong>时间范围:</strong> {{ date_range }}</p>

    <div class="stats">
        <div class="stats-item">📝 总文章数: <span class="stats-value">{{ stats.total }}</span></div>
        <div class="stats-item">✅ 新文章: <span class="stats-value">{{ stats.new }}</span></div>
        <div class="stats-item">🔄 重复文章: <span class="stats-value">{{ stats.duplicates }}</span></div>
        <div class="stats-item">🏷️ 标签数: <span class="stats-value">{{ stats.tags }}</span></div>
    </div>

    {% for group_name, group_data in groups.items() %}
    <div class="group-section">
        <h2>📂 {{ group_name }}</h2>
        <p>共 <strong>{{ group_data.count }}</strong> 篇文章</p>

        {% for tag_name, articles in group_data.by_tag.items() %}
        <h3>🏷️ {{ tag_name }} ({{ articles|length }})</h3>
        {% for article in articles %}
        <div class="article {% if article.is_duplicate %}duplicate{% endif %}">
            <h4 class="article-title">
                <a href="{{ article.url }}" target="_blank">{{ article.title }}</a>
                {% if article.is_duplicate %}<span style="color:#e74c3c;font-size:12px;">[重复]</span>{% endif %}
            </h4>
            <div class="article-meta">
                📰 {{ article.feed_title }} |
                📅 {{ article.published_at }} |
                ✍️ {{ article.author or '未知' }}
            </div>
            <div class="article-summary">{{ article.summary }}</div>
            {% if article.tags %}
            <div style="margin-top:8px;">
                {% for t in article.tags %}
                <span class="tag" style="background:{{ t.color }};">{{ t.name }}</span>
                {% endfor %}
            </div>
            {% endif %}
        </div>
        {% endfor %}
        {% endfor %}
    </div>
    {% endfor %}

    <p style="margin-top:40px;color:#95a5a6;font-size:12px;text-align:center;">
        由 RSS Reader 自动生成 | {{ generated_at }}
    </p>
</body>
</html>
""",
    "markdown": """# 📰 RSS 摘要报告

**时间范围:** {{ date_range }}

## 📊 统计概览

- 📝 总文章数: **{{ stats.total }}**
- ✅ 新文章: **{{ stats.new }}**
- 🔄 重复文章: **{{ stats.duplicates }}**
- 🏷️ 标签数: **{{ stats.tags }}**

{% for group_name, group_data in groups.items() %}
---

## 📂 {{ group_name }}

共 **{{ group_data.count }}** 篇文章

{% for tag_name, articles in group_data.by_tag.items() %}
### 🏷️ {{ tag_name }} ({{ articles|length }})

{% for article in articles %}
#### {{ '🔄 ' if article.is_duplicate else '' }}[{{ article.title }}]({{ article.url }})

> 📰 {{ article.feed_title }} | 📅 {{ article.published_at }} | ✍️ {{ article.author or '未知' }}

{{ article.summary }}

{% if article.tags %}标签: {% for t in article.tags %} `{{ t.name }}` {% endfor %}{% endif %}

{% endfor %}
{% endfor %}
{% endfor %}

---

*由 RSS Reader 自动生成 | {{ generated_at }}*
""",
}


class ReportGenerator:
    def __init__(self):
        try:
            from jinja2 import Template
            self._template_cls = Template
        except ImportError:
            logger.warning("Jinja2 not available, using basic string formatting")
            self._template_cls = None

    def _get_date_range(self, period: str = "daily") -> str:
        now = datetime.now()
        if period == "daily":
            start = now - timedelta(days=1)
            return f"{start.strftime('%Y-%m-%d')} 至 {now.strftime('%Y-%m-%d')}"
        elif period == "weekly":
            start = now - timedelta(days=7)
            return f"{start.strftime('%Y-%m-%d')} 至 {now.strftime('%Y-%m-%d')}"
        else:
            return f"截至 {now.strftime('%Y-%m-%d %H:%M')}"

    def _prepare_template_data(self, period: str, group_id: Optional[int] = None):
        groups = defaultdict(lambda: {"count": 0, "by_tag": defaultdict(list)})
        tag_count = set()

        session = get_session()
        try:
            from sqlalchemy.orm import joinedload

            now = datetime.utcnow()
            if period == "daily":
                since = now - timedelta(days=1)
            elif period == "weekly":
                since = now - timedelta(days=7)
            else:
                since = now - timedelta(days=30)

            query = session.query(Article).options(
                joinedload(Article.tags),
                joinedload(Article.feed).joinedload(Feed.group)
            ).filter(
                Article.fetched_at >= since,
                Article.is_ignored == False,
            )

            if group_id is not None:
                query = query.join(Feed).filter(Feed.group_id == group_id)

            articles = query.order_by(Article.published_at.desc()).all()

            for article in articles:
                feed = article.feed
                group_name = feed.group.name if feed and feed.group else "未分组"
                feed_title = feed.title if feed else "未知来源"

                tags = [{"name": t.name, "color": t.color} for t in article.tags]
                for t in article.tags:
                    tag_count.add(t.id)

                article_data = {
                    "id": article.id,
                    "title": article.title,
                    "url": article.url,
                    "summary": article.summary or "",
                    "published_at": article.published_at.strftime("%Y-%m-%d %H:%M") if article.published_at else "未知",
                    "author": article.author,
                    "feed_title": feed_title,
                    "is_duplicate": article.is_duplicate,
                    "tags": tags,
                }

                groups[group_name]["count"] += 1
                if article.tags:
                    for tag in article.tags:
                        groups[group_name]["by_tag"][tag.name].append(article_data)
                else:
                    groups[group_name]["by_tag"]["未分类"].append(article_data)

            total = len(articles)
            duplicates = sum(1 for a in articles if a.is_duplicate)

            return {
                "date_range": self._get_date_range(period),
                "generated_at": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                "stats": {
                    "total": total,
                    "new": total - duplicates,
                    "duplicates": duplicates,
                    "tags": len(tag_count),
                },
                "groups": dict(groups),
            }
        finally:
            session.close()

    def generate(self, period: str = "daily", fmt: str = "markdown", group_id: Optional[int] = None) -> str:
        data = self._prepare_template_data(period, group_id)
        template_str = REPORT_TEMPLATES.get(fmt, REPORT_TEMPLATES["markdown"])

        if self._template_cls:
            template = self._template_cls(template_str)
            return template.render(**data)
        else:
            return self._simple_format(template_str, data)

    def _simple_format(self, template: str, data: dict) -> str:
        result = template
        for key, value in data.items():
            if isinstance(value, (str, int)):
                result = result.replace("{{ " + key + " }}", str(value))
            elif isinstance(value, dict) and key == "stats":
                for k, v in value.items():
                    result = result.replace("{{ stats." + k + " }}", str(v))
        return result

    def save_report(self, period: str = "daily", fmt: str = "markdown", output_dir: Optional[str] = None) -> str:
        content = self.generate(period, fmt)
        ext = "html" if fmt == "html" else "md"
        timestamp = datetime.now().strftime("%Y%m%d")

        if output_dir:
            out_path = f"{output_dir}/rss_report_{period}_{timestamp}.{ext}"
        else:
            out_path = f"{settings.data_dir}/reports/rss_report_{period}_{timestamp}.{ext}"

        import os
        os.makedirs(os.path.dirname(out_path), exist_ok=True)

        with open(out_path, "w", encoding="utf-8") as f:
            f.write(content)

        logger.info(f"Report saved to {out_path}")
        return out_path
