import time
from datetime import datetime
from typing import List, Dict, Optional, Any
from urllib.parse import urlparse

import feedparser
import requests
from dateutil import parser as date_parser

from .config import settings
from .logger import logger
from .models import Feed, Article, FetchLog, get_session


class RSSFetcher:
    def __init__(self):
        self.session = requests.Session()
        self.session.headers.update({
            "User-Agent": settings.user_agent,
            "Accept": "application/rss+xml, application/rdf+xml, application/atom+xml, application/xml;q=0.9, text/xml;q=0.8"
        })

    def _fetch_with_retry(self, url: str, timeout: int = 30, max_retries: int = 3) -> Optional[str]:
        last_error = None
        for attempt in range(max_retries):
            try:
                response = self.session.get(url, timeout=timeout)
                response.raise_for_status()
                if not response.content:
                    raise ValueError("Empty response")
                return response.text
            except Exception as e:
                last_error = e
                logger.warning(f"Fetch attempt {attempt + 1}/{max_retries} failed for {url}: {e}")
                if attempt < max_retries - 1:
                    time.sleep(settings.retry_interval * (attempt + 1))
        logger.error(f"Failed to fetch {url} after {max_retries} attempts: {last_error}")
        return None

    def _parse_feed(self, feed_content: str, feed_url: str) -> feedparser.FeedParserDict:
        parsed = feedparser.parse(feed_content)
        if parsed.bozo and not parsed.entries:
            logger.warning(f"Feed parse error for {feed_url}: {parsed.bozo_exception}")
        return parsed

    def _parse_datetime(self, entry: feedparser.FeedParserDict) -> Optional[datetime]:
        for key in ["published_parsed", "updated_parsed", "created_parsed"]:
            val = getattr(entry, key, None) or entry.get(key)
            if val:
                try:
                    return datetime(*val[:6])
                except (TypeError, ValueError):
                    pass
        for key in ["published", "updated", "created"]:
            val = entry.get(key)
            if val:
                try:
                    return date_parser.parse(val)
                except (ValueError, TypeError):
                    pass
        return None

    def _extract_entry_data(self, entry: feedparser.FeedParserDict, feed_id: int) -> Dict[str, Any]:
        guid = entry.get("id") or entry.get("guid") or entry.get("link")
        if not guid:
            return {}

        title = entry.get("title", "")
        link = entry.get("link", "")
        author = entry.get("author", "")

        summary = entry.get("summary", "") or entry.get("description", "")
        content_html = ""
        if "content" in entry:
            for c in entry["content"]:
                if c.get("type", "").startswith("text/html") or not content_html:
                    content_html = c.get("value", "")
                    break
        if not content_html:
            content_html = summary

        return {
            "feed_id": feed_id,
            "guid": guid,
            "title": title,
            "summary": summary,
            "content_html": content_html,
            "url": link,
            "author": author,
            "published_at": self._parse_datetime(entry),
        }

    def fetch_feed(self, feed: Feed) -> List[Dict[str, Any]]:
        feed_id = feed.id
        log = FetchLog(feed_id=feed_id, status="running")
        session = get_session()
        try:
            session.add(log)
            session.commit()

            feed_db = session.get(Feed, feed_id)
            if not feed_db:
                logger.error(f"Feed {feed_id} not found in database")
                return []

            logger.info(f"Fetching feed: {feed_db.title} ({feed_db.url})")

            content = self._fetch_with_retry(
                feed_db.url,
                timeout=feed_db.request_timeout,
                max_retries=feed_db.max_retries
            )

            if not content:
                log.status = "failed"
                log.error_message = "Failed to fetch feed content"
                feed_db.last_status = "failed"
                feed_db.error_count += 1
                session.commit()
                return []

            parsed = self._parse_feed(content, feed_db.url)
            entries_data = []
            new_count = 0

            for entry in parsed.entries:
                data = self._extract_entry_data(entry, feed_id)
                if not data:
                    continue

                existing = session.query(Article).filter_by(
                    feed_id=feed_id,
                    guid=data["guid"]
                ).first()

                if not existing:
                    entries_data.append(data)
                    new_count += 1

            if hasattr(parsed.feed, "title") and parsed.feed.title != feed_db.title:
                feed_db.title = parsed.feed.title

            feed_db.last_fetched = datetime.utcnow()
            feed_db.last_status = "success"
            feed_db.error_count = 0
            feed_db.updated_at = datetime.utcnow()

            log.status = "success"
            log.new_articles = new_count
            log.finished_at = datetime.utcnow()
            session.commit()

            logger.info(f"Fetched {feed_db.title}: {new_count} new articles, {len(parsed.entries)} total")
            return entries_data

        except Exception as e:
            logger.exception(f"Error fetching feed {feed_id}: {e}")
            log.status = "failed"
            log.error_message = str(e)
            try:
                feed_db = session.get(Feed, feed_id)
                if feed_db:
                    feed_db.last_status = "error"
                    feed_db.error_count += 1
                    feed_db.updated_at = datetime.utcnow()
                session.commit()
            except Exception:
                pass
            return []
        finally:
            session.close()

    def fetch_all(self) -> List[Dict[str, Any]]:
        session = get_session()
        try:
            feeds = session.query(Feed).filter_by(enabled=True).all()
            all_entries = []
            for feed in feeds:
                entries = self.fetch_feed(feed)
                all_entries.extend(entries)
            return all_entries
        finally:
            session.close()
