import re
from typing import List, Dict, Any
from bs4 import BeautifulSoup, Comment
import html2text

from .logger import logger
from .models import Article, get_session


AD_SELECTORS = [
    "script",
    "style",
    "iframe",
    "noscript",
    "aside",
    "nav",
    "footer",
    "header",
    ".ads",
    ".ad",
    ".advertisement",
    ".advertising",
    ".banner",
    ".sponsored",
    ".promo",
    "#ads",
    "#ad",
    ".social-share",
    ".share-buttons",
    ".related-posts",
    ".related-articles",
    "[class*='ad-']",
    "[class*='ads-']",
    "[id*='ad-']",
    "[id*='ads-']",
]


class ContentCleaner:
    def __init__(self):
        self.html2md = html2text.HTML2Text()
        self.html2md.ignore_links = False
        self.html2md.ignore_images = False
        self.html2md.body_width = 0
        self.html2md.ignore_emphasis = False
        self.html2md.mark_code = True

    def _remove_ads(self, soup: BeautifulSoup) -> BeautifulSoup:
        for selector in AD_SELECTORS:
            try:
                for element in soup.select(selector):
                    element.decompose()
            except Exception:
                continue

        for comment in soup.find_all(string=lambda text: isinstance(text, Comment)):
            comment.extract()

        return soup

    def _clean_html(self, html: str) -> str:
        if not html:
            return ""
        soup = BeautifulSoup(html, "lxml")
        soup = self._remove_ads(soup)

        for tag in soup.find_all(True):
            if tag.name in ["script", "style", "noscript"]:
                continue
            attrs_to_remove = []
            for attr in tag.attrs:
                if attr in ["onclick", "onload", "onerror", "onmouseover", "style", "class", "id"]:
                    attrs_to_remove.append(attr)
            for attr in attrs_to_remove:
                del tag[attr]

        return str(soup)

    def _html_to_markdown(self, html: str) -> str:
        if not html:
            return ""
        try:
            return self.html2md.handle(html).strip()
        except Exception as e:
            logger.warning(f"HTML to Markdown conversion failed: {e}")
            return html.strip()

    def _extract_plain_text(self, html: str) -> str:
        if not html:
            return ""
        soup = BeautifulSoup(html, "lxml")
        for tag in soup(["script", "style"]):
            tag.decompose()
        return soup.get_text(separator=" ", strip=True)

    def _generate_summary(self, content: str, max_len: int = 300) -> str:
        plain = self._extract_plain_text(content) if content else ""
        if len(plain) <= max_len:
            return plain
        return plain[:max_len].rsplit(" ", 1)[0] + "..."

    def clean_article_data(self, article_data: Dict[str, Any]) -> Dict[str, Any]:
        cleaned_html = self._clean_html(article_data.get("content_html", ""))
        markdown = self._html_to_markdown(cleaned_html)

        summary = article_data.get("summary", "")
        if not summary or len(summary) < 50:
            summary = self._generate_summary(cleaned_html)

        return {
            **article_data,
            "content_html": cleaned_html,
            "content": markdown,
            "summary": summary,
        }

    def clean_and_save(self, articles_data: List[Dict[str, Any]]) -> List[int]:
        session = get_session()
        saved_ids = []
        try:
            for data in articles_data:
                cleaned = self.clean_article_data(data)
                article = Article(**cleaned)
                session.add(article)
                session.flush()
                saved_ids.append(article.id)
            session.commit()
            logger.info(f"Cleaned and saved {len(saved_ids)} articles")
            return saved_ids
        except Exception as e:
            logger.exception(f"Failed to save articles: {e}")
            session.rollback()
            return []
        finally:
            session.close()
