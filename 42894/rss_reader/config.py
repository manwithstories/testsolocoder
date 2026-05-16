import os
from pathlib import Path
from dataclasses import dataclass, field
from typing import List


def _getenv(name: str, default: str = "") -> str:
    val = os.getenv(name, default)
    if isinstance(val, str):
        return val.strip().strip('"').strip("'")
    return val


def _getenv_int(name: str, default: int = 0) -> int:
    val = _getenv(name, str(default))
    try:
        return int(val)
    except (ValueError, TypeError):
        return default


def _getenv_bool(name: str, default: bool = False) -> bool:
    val = _getenv(name, str(default)).lower()
    return val in ("true", "1", "yes", "y", "on")


def _getenv_list(name: str, default: List[str] = None) -> List[str]:
    default = default or []
    val = _getenv(name, "")
    if val:
        return [s.strip() for s in val.split(",") if s.strip()]
    return default


@dataclass
class SMTPConfig:
    host: str = field(default_factory=lambda: _getenv("RSS_SMTP_HOST", ""))
    port: int = field(default_factory=lambda: _getenv_int("RSS_SMTP_PORT", 587))
    username: str = field(default_factory=lambda: _getenv("RSS_SMTP_USERNAME", ""))
    password: str = field(default_factory=lambda: _getenv("RSS_SMTP_PASSWORD", ""))
    use_tls: bool = field(default_factory=lambda: _getenv_bool("RSS_SMTP_USE_TLS", True))
    from_addr: str = field(default_factory=lambda: _getenv("RSS_SMTP_FROM_ADDR", ""))
    to_addrs: List[str] = field(default_factory=lambda: _getenv_list("RSS_SMTP_TO_ADDRS"))


@dataclass
class Settings:
    data_dir: Path = field(default_factory=lambda: Path(_getenv("RSS_DATA_DIR", str(Path.home() / ".rss_reader"))).expanduser())
    db_url: str = field(default_factory=lambda: _getenv("RSS_DB_URL", ""))
    log_level: str = field(default_factory=lambda: _getenv("RSS_LOG_LEVEL", "INFO"))
    user_agent: str = field(default_factory=lambda: _getenv("RSS_USER_AGENT", "Mozilla/5.0 (compatible; RSSReader/0.1)"))
    request_timeout: int = field(default_factory=lambda: _getenv_int("RSS_REQUEST_TIMEOUT", 30))
    max_retries: int = field(default_factory=lambda: _getenv_int("RSS_MAX_RETRIES", 3))
    retry_interval: int = field(default_factory=lambda: _getenv_int("RSS_RETRY_INTERVAL", 5))
    dedup_threshold: int = field(default_factory=lambda: _getenv_int("RSS_DEDUP_THRESHOLD", 3))
    report_time: str = field(default_factory=lambda: _getenv("RSS_REPORT_TIME", "08:00"))
    report_format: str = field(default_factory=lambda: _getenv("RSS_REPORT_FORMAT", "markdown"))
    smtp: SMTPConfig = field(default_factory=SMTPConfig)

    def __post_init__(self):
        self.data_dir.mkdir(parents=True, exist_ok=True)
        if not self.db_url:
            self.db_url = f"sqlite:///{self.data_dir}/rss.db"


settings = Settings()
