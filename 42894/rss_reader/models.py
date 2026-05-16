from datetime import datetime
from typing import List, Optional
from sqlalchemy import (
    Column, Integer, String, Text, DateTime, Boolean,
    ForeignKey, Table, create_engine
)
from sqlalchemy.orm import (
    relationship, declarative_base, sessionmaker, Mapped, mapped_column
)

from .config import settings

Base = declarative_base()


article_tags = Table(
    "article_tags",
    Base.metadata,
    Column("article_id", Integer, ForeignKey("articles.id"), primary_key=True),
    Column("tag_id", Integer, ForeignKey("tags.id"), primary_key=True),
)


class Group(Base):
    __tablename__ = "groups"

    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(100), unique=True, index=True)
    description: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)

    feeds: Mapped[List["Feed"]] = relationship(back_populates="group")


class Feed(Base):
    __tablename__ = "feeds"

    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(255))
    url: Mapped[str] = mapped_column(String(500), unique=True, index=True)
    group_id: Mapped[Optional[int]] = mapped_column(ForeignKey("groups.id"), nullable=True)
    fetch_interval: Mapped[int] = mapped_column(default=30)
    request_timeout: Mapped[int] = mapped_column(default=30)
    max_retries: Mapped[int] = mapped_column(default=3)
    enabled: Mapped[bool] = mapped_column(default=True)
    last_fetched: Mapped[Optional[datetime]] = mapped_column(nullable=True)
    last_status: Mapped[Optional[str]] = mapped_column(String(50), nullable=True)
    error_count: Mapped[int] = mapped_column(default=0)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    updated_at: Mapped[datetime] = mapped_column(default=datetime.utcnow, onupdate=datetime.utcnow)

    group: Mapped[Optional[Group]] = relationship(back_populates="feeds")
    articles: Mapped[List["Article"]] = relationship(back_populates="feed", cascade="all, delete-orphan")


class Tag(Base):
    __tablename__ = "tags"

    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(50), unique=True, index=True)
    color: Mapped[str] = mapped_column(String(7), default="#3b82f6")
    is_auto: Mapped[bool] = mapped_column(default=False)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)

    articles: Mapped[List["Article"]] = relationship(
        secondary=article_tags, back_populates="tags"
    )


class FilterRule(Base):
    __tablename__ = "filter_rules"

    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(100))
    rule_type: Mapped[str] = mapped_column(String(20))
    pattern: Mapped[str] = mapped_column(Text)
    is_regex: Mapped[bool] = mapped_column(default=False)
    apply_to: Mapped[str] = mapped_column(String(20), default="title_content")
    action: Mapped[str] = mapped_column(String(20))
    tag_id: Mapped[Optional[int]] = mapped_column(ForeignKey("tags.id"), nullable=True)
    enabled: Mapped[bool] = mapped_column(default=True)
    priority: Mapped[int] = mapped_column(default=0)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)

    tag: Mapped[Optional[Tag]] = relationship()


class Article(Base):
    __tablename__ = "articles"

    id: Mapped[int] = mapped_column(primary_key=True)
    feed_id: Mapped[int] = mapped_column(ForeignKey("feeds.id"))
    guid: Mapped[str] = mapped_column(String(500), index=True)
    title: Mapped[str] = mapped_column(String(500))
    summary: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    content: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    content_html: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    url: Mapped[str] = mapped_column(String(500))
    author: Mapped[Optional[str]] = mapped_column(String(200), nullable=True)
    published_at: Mapped[Optional[datetime]] = mapped_column(nullable=True)
    fetched_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    simhash: Mapped[Optional[str]] = mapped_column(String(64), nullable=True, index=True)
    is_duplicate: Mapped[bool] = mapped_column(default=False)
    duplicate_of: Mapped[Optional[int]] = mapped_column(nullable=True)
    is_read: Mapped[bool] = mapped_column(default=False)
    is_ignored: Mapped[bool] = mapped_column(default=False)

    feed: Mapped[Feed] = relationship(back_populates="articles")
    tags: Mapped[List[Tag]] = relationship(
        secondary=article_tags, back_populates="articles"
    )


class FetchLog(Base):
    __tablename__ = "fetch_logs"

    id: Mapped[int] = mapped_column(primary_key=True)
    feed_id: Mapped[int] = mapped_column(ForeignKey("feeds.id"), index=True)
    started_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    finished_at: Mapped[Optional[datetime]] = mapped_column(nullable=True)
    status: Mapped[str] = mapped_column(String(20))
    new_articles: Mapped[int] = mapped_column(default=0)
    error_message: Mapped[Optional[str]] = mapped_column(Text, nullable=True)


engine = create_engine(settings.db_url, echo=False)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def init_db():
    Base.metadata.create_all(bind=engine)


def get_session():
    return SessionLocal()
