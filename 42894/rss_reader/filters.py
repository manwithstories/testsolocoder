import re
from typing import List, Tuple, Optional
from sqlalchemy import and_

from .logger import logger
from .models import Article, FilterRule, Tag, get_session


class FilterEngine:
    def __init__(self):
        pass

    def _get_match_text(self, article: Article, apply_to: str) -> str:
        parts = []
        if apply_to in ["title", "title_content"]:
            parts.append(article.title or "")
        if apply_to in ["summary", "title_content"]:
            parts.append(article.summary or "")
        if apply_to in ["content", "title_content"]:
            parts.append(article.content or "")
        return "\n".join(parts)

    def _match_rule(self, text: str, rule: FilterRule) -> bool:
        if rule.is_regex:
            try:
                pattern = re.compile(rule.pattern, re.IGNORECASE | re.MULTILINE)
                return bool(pattern.search(text))
            except re.error as e:
                logger.warning(f"Invalid regex pattern in rule {rule.id}: {e}")
                return False
        else:
            keywords = [k.strip() for k in rule.pattern.split(",") if k.strip()]
            text_lower = text.lower()
            if rule.rule_type == "whitelist":
                return any(kw.lower() in text_lower for kw in keywords)
            elif rule.rule_type == "blacklist":
                return any(kw.lower() in text_lower for kw in keywords)
            elif rule.rule_type == "tag":
                return any(kw.lower() in text_lower for kw in keywords)
        return False

    def _process_article_with_session(self, article: Article, session) -> Tuple[bool, List[int]]:
        rules = session.query(FilterRule).filter_by(enabled=True).order_by(FilterRule.priority.desc()).all()

        should_ignore = False
        tag_ids_to_add: List[int] = []
        article_text = ""

        for rule in rules:
            if not article_text:
                article_text = self._get_match_text(article, rule.apply_to)

            if self._match_rule(article_text, rule):
                if rule.action == "ignore":
                    should_ignore = True
                    logger.info(f"Article {article.id} matched ignore rule: {rule.name}")
                elif rule.action == "tag" and rule.tag_id and rule.tag_id not in tag_ids_to_add:
                    tag_ids_to_add.append(rule.tag_id)
                    tag = session.get(Tag, rule.tag_id)
                    if tag:
                        logger.info(f"Article {article.id} tagged with '{tag.name}' by rule: {rule.name}")

        return should_ignore, tag_ids_to_add

    def process_articles(self, article_ids: List[int]) -> None:
        session = get_session()
        try:
            articles = session.query(Article).filter(Article.id.in_(article_ids)).all()
            for article in articles:
                should_ignore, tag_ids = self._process_article_with_session(article, session)
                if should_ignore:
                    article.is_ignored = True
                for tag_id in tag_ids:
                    tag = session.get(Tag, tag_id)
                    if tag and tag not in article.tags:
                        article.tags.append(tag)
            session.commit()
            logger.info(f"Processed {len(articles)} articles through filter rules")
        except Exception as e:
            logger.exception(f"Error processing articles: {e}")
            session.rollback()
        finally:
            session.close()

    def add_rule(
        self,
        name: str,
        rule_type: str,
        pattern: str,
        action: str,
        is_regex: bool = False,
        apply_to: str = "title_content",
        tag_name: Optional[str] = None,
        priority: int = 0,
    ) -> Optional[int]:
        session = get_session()
        try:
            tag_id = None
            if action == "tag" and tag_name:
                tag = session.query(Tag).filter_by(name=tag_name).first()
                if not tag:
                    tag = Tag(name=tag_name, is_auto=True)
                    session.add(tag)
                    session.flush()
                tag_id = tag.id

            rule = FilterRule(
                name=name,
                rule_type=rule_type,
                pattern=pattern,
                is_regex=is_regex,
                apply_to=apply_to,
                action=action,
                tag_id=tag_id,
                priority=priority,
            )
            session.add(rule)
            session.commit()
            logger.info(f"Added filter rule: {name}")
            return rule.id
        except Exception as e:
            logger.exception(f"Failed to add filter rule: {e}")
            session.rollback()
            return None
        finally:
            session.close()
