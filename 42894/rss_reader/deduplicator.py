import re
from datetime import datetime, timedelta
from typing import Optional, List, Tuple
from simhash import Simhash

from .config import settings
from .logger import logger
from .models import Article, get_session


class SimHashDeduplicator:
    def __init__(self, threshold: int = 3):
        self.threshold = threshold or settings.dedup_threshold

    def _tokenize(self, text: str) -> List[str]:
        if not text:
            return []
        text = text.lower()
        text = re.sub(r'[^\w\u4e00-\u9fff\s]', ' ', text)
        words = re.findall(r'[\w\u4e00-\u9fff]+', text)
        return [w for w in words if len(w) > 1]

    def compute_simhash(self, article: Article) -> str:
        text = f"{article.title}\n{article.summary or ''}\n{article.content or ''}"
        tokens = self._tokenize(text)
        if not tokens:
            return ""
        sh = Simhash(tokens)
        return str(sh.value)

    def hamming_distance(self, hash1: str, hash2: str) -> int:
        try:
            return bin(int(hash1) ^ int(hash2)).count("1")
        except (ValueError, TypeError):
            return 999

    def find_duplicate(
        self,
        article: Article,
        existing_hashes: List[Tuple[int, str]],
    ) -> Optional[int]:
        if not article.simhash:
            return None

        for art_id, art_hash in existing_hashes:
            if art_id == article.id or not art_hash:
                continue
            distance = self.hamming_distance(article.simhash, art_hash)
            if distance <= self.threshold:
                return art_id
        return None

    def process_article(self, article: Article) -> Tuple[bool, Optional[int]]:
        if not article.simhash:
            article.simhash = self.compute_simhash(article)
            if not article.simhash:
                return False, None

        session = get_session()
        try:
            since_date = article.published_at or article.fetched_at
            search_window = since_date - timedelta(days=30)

            existing = session.query(Article.id, Article.simhash).filter(
                Article.id != article.id,
                Article.simhash.isnot(None),
                Article.simhash != "",
                Article.fetched_at >= search_window,
            ).all()

            duplicate_of = self.find_duplicate(article, existing)
            if duplicate_of:
                article.is_duplicate = True
                article.duplicate_of = duplicate_of
                logger.info(f"Article {article.id} marked as duplicate of {duplicate_of}")
                return True, duplicate_of
            return False, None
        finally:
            session.close()

    def process_articles(self, article_ids: List[int]) -> int:
        session = get_session()
        try:
            articles = session.query(Article).filter(Article.id.in_(article_ids)).all()
            dup_count = 0

            for article in articles:
                if not article.simhash:
                    article.simhash = self.compute_simhash(article)

            existing_hashes = []
            since_date = datetime.utcnow() - timedelta(days=30)
            existing = session.query(Article.id, Article.simhash).filter(
                Article.simhash.isnot(None),
                Article.simhash != "",
                Article.fetched_at >= since_date,
            ).all()
            existing_hashes = [(a.id, a.simhash) for a in existing]

            for article in articles:
                if article.simhash:
                    duplicate_of = self.find_duplicate(article, existing_hashes)
                    if duplicate_of:
                        article.is_duplicate = True
                        article.duplicate_of = duplicate_of
                        dup_count += 1
                    else:
                        existing_hashes.append((article.id, article.simhash))

            session.commit()
            logger.info(f"Deduplication complete: {dup_count}/{len(articles)} articles marked as duplicates")
            return dup_count
        except Exception as e:
            logger.exception(f"Deduplication error: {e}")
            session.rollback()
            return 0
        finally:
            session.close()
