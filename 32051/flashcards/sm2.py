from __future__ import annotations

from datetime import datetime, timedelta, timezone
from typing import Tuple

from .models import FlashCard


MIN_EASE_FACTOR = 1.3


def calculate_next_review(
    card: FlashCard, quality: int
) -> Tuple[float, float, int, datetime]:
    """
    根据SM-2算法计算下次复习时间

    quality 评分说明:
    0 - 完全不知道，完全忘记
    1 - 看到答案完全没有印象
    2 - 看到答案有点印象但想不起来
    3 - 回答正确但需要努力回忆
    4 - 回答正确但有点迟疑
    5 - 回答完全正确，毫不犹豫

    返回: (new_interval, new_ease_factor, new_repetitions, next_review_at)
    """
    quality = max(0, min(5, quality))

    interval = card.interval
    ease_factor = card.ease_factor
    repetitions = card.repetitions

    if quality < 3:
        repetitions = 0
        interval = 0.0
    else:
        if repetitions == 0:
            interval = 1.0
        elif repetitions == 1:
            interval = 6.0
        else:
            interval = round(interval * ease_factor, 1)

        repetitions += 1

    ease_factor = ease_factor + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02))
    ease_factor = max(MIN_EASE_FACTOR, round(ease_factor, 2))

    if quality < 3:
        interval = 0.0
        next_review_at = datetime.now(timezone.utc)
    else:
        next_review_at = datetime.now(timezone.utc) + timedelta(days=interval)

    return interval, ease_factor, repetitions, next_review_at


def apply_review_result(card: FlashCard, quality: int) -> None:
    """将复习结果应用到卡片上"""
    interval, ease_factor, repetitions, next_review_at = calculate_next_review(
        card, quality
    )

    card.interval = interval
    card.ease_factor = ease_factor
    card.repetitions = repetitions
    card.last_reviewed_at = datetime.now(timezone.utc).isoformat()
    card.next_review_at = next_review_at.isoformat()
    card.review_count += 1
    card.touch()


def get_quality_description(quality: int) -> str:
    """获取评分的文字描述"""
    descriptions = {
        0: "完全忘记 - 完全不知道答案",
        1: "毫无印象 - 看到答案也想不起来",
        2: "似曾相识 - 有点印象但答不出来",
        3: "勉强记住 - 需要努力回忆才答对",
        4: "轻松记住 - 稍微迟疑就答对了",
        5: "牢记在心 - 毫不犹豫就答对了",
    }
    return descriptions.get(quality, "未知")


def format_interval_days(days: float) -> str:
    """格式化间隔天数为友好字符串"""
    if days < 1:
        return "今天"
    elif days < 30:
        return f"{int(days)}天"
    elif days < 365:
        months = round(days / 30, 1)
        return f"{months}个月"
    else:
        years = round(days / 365, 1)
        return f"{years}年"
