from __future__ import annotations

import json
import uuid
from dataclasses import dataclass, field, asdict
from datetime import datetime, timezone
from enum import Enum
from typing import List, Optional


class Difficulty(str, Enum):
    EASY = "easy"
    MEDIUM = "medium"
    HARD = "hard"


@dataclass
class FlashCard:
    front: str
    back: str
    id: str = field(default_factory=lambda: uuid.uuid4().hex)
    tags: List[str] = field(default_factory=list)
    difficulty: Difficulty = Difficulty.MEDIUM
    created_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    updated_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    last_reviewed_at: Optional[str] = None
    next_review_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    interval: float = 0.0
    ease_factor: float = 2.5
    repetitions: int = 0
    review_count: int = 0

    def to_dict(self) -> dict:
        data = asdict(self)
        data["difficulty"] = self.difficulty.value
        return data

    @classmethod
    def from_dict(cls, data: dict) -> "FlashCard":
        difficulty = Difficulty(data.get("difficulty", Difficulty.MEDIUM.value))
        return cls(
            id=data.get("id", uuid.uuid4().hex),
            front=data["front"],
            back=data["back"],
            tags=data.get("tags", []),
            difficulty=difficulty,
            created_at=data.get("created_at", datetime.now(timezone.utc).isoformat()),
            updated_at=data.get("updated_at", datetime.now(timezone.utc).isoformat()),
            last_reviewed_at=data.get("last_reviewed_at"),
            next_review_at=data.get("next_review_at", datetime.now(timezone.utc).isoformat()),
            interval=data.get("interval", 0.0),
            ease_factor=data.get("ease_factor", 2.5),
            repetitions=data.get("repetitions", 0),
            review_count=data.get("review_count", 0),
        )

    def touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc).isoformat()

    def is_due(self) -> bool:
        now = datetime.now(timezone.utc)
        next_review = datetime.fromisoformat(self.next_review_at)
        return now >= next_review


@dataclass
class StudyRecord:
    card_id: str
    quality: int
    reviewed_at: str
    session_id: str

    def to_dict(self) -> dict:
        return asdict(self)

    @classmethod
    def from_dict(cls, data: dict) -> "StudyRecord":
        return cls(**data)


@dataclass
class Deck:
    name: str
    id: str = field(default_factory=lambda: uuid.uuid4().hex)
    cards: List[FlashCard] = field(default_factory=list)
    study_records: List[StudyRecord] = field(default_factory=list)
    created_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    updated_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "name": self.name,
            "cards": [card.to_dict() for card in self.cards],
            "study_records": [record.to_dict() for record in self.study_records],
            "created_at": self.created_at,
            "updated_at": self.updated_at,
        }

    @classmethod
    def from_dict(cls, data: dict) -> "Deck":
        return cls(
            id=data["id"],
            name=data["name"],
            cards=[FlashCard.from_dict(c) for c in data.get("cards", [])],
            study_records=[StudyRecord.from_dict(r) for r in data.get("study_records", [])],
            created_at=data.get("created_at", datetime.now(timezone.utc).isoformat()),
            updated_at=data.get("updated_at", datetime.now(timezone.utc).isoformat()),
        )

    def touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc).isoformat()

    def add_card(self, card: FlashCard) -> None:
        self.cards.append(card)
        self.touch()

    def remove_card(self, card_id: str) -> Optional[FlashCard]:
        for i, card in enumerate(self.cards):
            if card.id == card_id:
                self.touch()
                return self.cards.pop(i)
        return None

    def get_card(self, card_id: str) -> Optional[FlashCard]:
        for card in self.cards:
            if card.id == card_id:
                return card
        return None

    def add_study_record(self, record: StudyRecord) -> None:
        self.study_records.append(record)

    def get_due_cards(self) -> List[FlashCard]:
        return [card for card in self.cards if card.is_due()]

    def get_all_tags(self) -> List[str]:
        tags = set()
        for card in self.cards:
            tags.update(card.tags)
        return sorted(tags)
