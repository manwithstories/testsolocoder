from __future__ import annotations

import time
import uuid
from datetime import datetime, timezone
from typing import List, Optional, Tuple

from .models import Deck, FlashCard, StudyRecord
from .sm2 import apply_review_result, format_interval_days, get_quality_description


class StudySession:
    def __init__(
        self,
        deck: Deck,
        tags: Optional[List[str]] = None,
        difficulty: Optional[str] = None,
        max_cards: Optional[int] = None,
        shuffle: bool = True,
    ):
        self.deck = deck
        self.session_id = uuid.uuid4().hex
        self.start_time = time.time()
        self.tags = tags or []
        self.difficulty = difficulty
        self.max_cards = max_cards
        self.shuffle = shuffle
        self.review_queue = self._build_queue()
        self.current_index = 0
        self.results: List[Tuple[FlashCard, int]] = []

    def _build_queue(self) -> List[FlashCard]:
        cards = self.deck.get_due_cards()

        if self.tags:
            cards = [
                c for c in cards if any(tag in c.tags for tag in self.tags)
            ]

        if self.difficulty:
            cards = [c for c in cards if c.difficulty.value == self.difficulty]

        if self.shuffle:
            import random
            random.shuffle(cards)

        if self.max_cards:
            cards = cards[: self.max_cards]

        return cards

    def has_next(self) -> bool:
        return self.current_index < len(self.review_queue)

    def next_card(self) -> Optional[FlashCard]:
        if not self.has_next():
            return None
        card = self.review_queue[self.current_index]
        self.current_index += 1
        return card

    def record_answer(self, card: FlashCard, quality: int) -> None:
        apply_review_result(card, quality)
        record = StudyRecord(
            card_id=card.id,
            quality=quality,
            reviewed_at=datetime.now(timezone.utc).isoformat(),
            session_id=self.session_id,
        )
        self.deck.add_study_record(record)
        self.results.append((card, quality))

    def skip_card(self, card: FlashCard) -> None:
        self.review_queue.append(card)

    def get_progress(self) -> Tuple[int, int]:
        return self.current_index, len(self.review_queue)

    def finish(self) -> dict:
        end_time = time.time()
        duration_seconds = end_time - self.start_time

        total_reviewed = len(self.results)
        if total_reviewed == 0:
            return {
                "total_cards": 0,
                "correct_count": 0,
                "incorrect_count": 0,
                "accuracy": 0.0,
                "duration_seconds": 0,
                "duration_formatted": "0秒",
                "avg_time_per_card": 0.0,
                "session_id": self.session_id,
            }

        correct_count = sum(1 for _, q in self.results if q >= 3)
        incorrect_count = total_reviewed - correct_count
        accuracy = round(correct_count / total_reviewed * 100, 1)
        avg_time = round(duration_seconds / total_reviewed, 1)

        minutes, seconds = divmod(int(duration_seconds), 60)
        if minutes > 0:
            duration_formatted = f"{minutes}分{seconds}秒"
        else:
            duration_formatted = f"{seconds}秒"

        return {
            "total_cards": total_reviewed,
            "correct_count": correct_count,
            "incorrect_count": incorrect_count,
            "accuracy": accuracy,
            "duration_seconds": duration_seconds,
            "duration_formatted": duration_formatted,
            "avg_time_per_card": avg_time,
            "session_id": self.session_id,
        }


class CardReviewResult:
    def __init__(self, card: FlashCard, quality: int):
        self.card = card
        self.quality = quality

    def get_quality_text(self) -> str:
        return get_quality_description(self.quality)

    def get_next_interval_text(self) -> str:
        return format_interval_days(self.card.interval)
