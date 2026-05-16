from datetime import datetime
from typing import Dict, Optional
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.interval import IntervalTrigger
from apscheduler.triggers.cron import CronTrigger

from .config import settings
from .logger import logger
from .models import Feed, get_session
from .fetcher import RSSFetcher
from .content_cleaner import ContentCleaner
from .filters import FilterEngine
from .deduplicator import SimHashDeduplicator


class FeedScheduler:
    def __init__(self):
        self.scheduler = BackgroundScheduler(timezone="UTC")
        self.fetcher = RSSFetcher()
        self.cleaner = ContentCleaner()
        self.filter_engine = FilterEngine()
        self.deduplicator = SimHashDeduplicator()
        self._job_ids: Dict[int, str] = {}

    def _process_feed(self, feed_id: int):
        session = get_session()
        try:
            feed = session.get(Feed, feed_id)
            if not feed or not feed.enabled:
                return

            logger.info(f"Scheduled fetch for feed: {feed.title}")

            raw_entries = self.fetcher.fetch_feed(feed)
            if not raw_entries:
                return

            saved_ids = self.cleaner.clean_and_save(raw_entries)
            if not saved_ids:
                return

            self.filter_engine.process_articles(saved_ids)
            self.deduplicator.process_articles(saved_ids)

            logger.info(f"Completed processing feed {feed.title}: {len(saved_ids)} articles")

        except Exception as e:
            logger.exception(f"Error in scheduled feed processing {feed_id}: {e}")
        finally:
            session.close()

    def add_feed_job(self, feed: Feed) -> Optional[str]:
        job_id = f"feed_{feed.id}"
        if job_id in self._job_ids:
            self.scheduler.remove_job(job_id)

        try:
            trigger = IntervalTrigger(minutes=feed.fetch_interval)
            job = self.scheduler.add_job(
                self._process_feed,
                trigger=trigger,
                args=[feed.id],
                id=job_id,
                name=f"Fetch {feed.title}",
                replace_existing=True,
                next_run_time=datetime.utcnow(),
            )
            self._job_ids[feed.id] = job_id
            logger.info(f"Scheduled feed {feed.title} every {feed.fetch_interval} minutes")
            return job.id
        except Exception as e:
            logger.exception(f"Failed to schedule feed {feed.id}: {e}")
            return None

    def remove_feed_job(self, feed_id: int):
        job_id = self._job_ids.pop(feed_id, None)
        if job_id and self.scheduler.get_job(job_id):
            self.scheduler.remove_job(job_id)
            logger.info(f"Removed scheduled job for feed {feed_id}")

    def load_all_feeds(self):
        session = get_session()
        try:
            feeds = session.query(Feed).filter_by(enabled=True).all()
            for feed in feeds:
                self.add_feed_job(feed)
            logger.info(f"Loaded {len(feeds)} feeds into scheduler")
        finally:
            session.close()

    def add_report_job(self, report_func, cron_expr: str = None):
        if not cron_expr:
            hour, minute = settings.report_time.split(":")
            cron_expr = f"{minute} {hour} * * *"

        try:
            trigger = CronTrigger.from_crontab(cron_expr)
            job = self.scheduler.add_job(
                report_func,
                trigger=trigger,
                id="daily_report",
                name="Daily Digest Report",
                replace_existing=True,
            )
            logger.info(f"Scheduled daily report at {settings.report_time}")
            return job.id
        except Exception as e:
            logger.exception(f"Failed to schedule report: {e}")
            return None

    def start(self):
        if not self.scheduler.running:
            self.load_all_feeds()
            self.scheduler.start()
            logger.info("Scheduler started")

    def stop(self):
        if self.scheduler.running:
            self.scheduler.shutdown()
            logger.info("Scheduler stopped")

    def run_once(self):
        session = get_session()
        try:
            feeds = session.query(Feed).filter_by(enabled=True).all()
            all_saved = []
            for feed in feeds:
                raw_entries = self.fetcher.fetch_feed(feed)
                if raw_entries:
                    saved_ids = self.cleaner.clean_and_save(raw_entries)
                    if saved_ids:
                        self.filter_engine.process_articles(saved_ids)
                        self.deduplicator.process_articles(saved_ids)
                        all_saved.extend(saved_ids)
            logger.info(f"One-time run complete: {len(all_saved)} new articles")
            return all_saved
        finally:
            session.close()
