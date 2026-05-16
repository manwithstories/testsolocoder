from dataclasses import dataclass, field, asdict
from typing import List, Optional, Dict, Any
from datetime import datetime, date
from enum import Enum
import uuid


class Priority(str, Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    URGENT = "urgent"


@dataclass
class Card:
    id: str
    title: str
    description: str = ""
    priority: Priority = Priority.MEDIUM
    tags: List[str] = field(default_factory=list)
    due_date: Optional[str] = None
    column_id: str = ""
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    completed_at: Optional[str] = None

    @classmethod
    def create(cls, title: str, description: str = "", priority: Priority = Priority.MEDIUM,
               tags: Optional[List[str]] = None, due_date: Optional[str] = None) -> "Card":
        return cls(
            id=str(uuid.uuid4()),
            title=title,
            description=description,
            priority=priority,
            tags=tags or [],
            due_date=due_date,
        )

    def is_overdue(self) -> bool:
        if not self.due_date:
            return False
        try:
            due = datetime.fromisoformat(self.due_date).date()
            return date.today() > due
        except (ValueError, TypeError):
            return False

    def days_until_due(self) -> Optional[int]:
        if not self.due_date:
            return None
        try:
            due = datetime.fromisoformat(self.due_date).date()
            return (due - date.today()).days
        except (ValueError, TypeError):
            return None

    def to_dict(self) -> Dict[str, Any]:
        data = asdict(self)
        data["priority"] = self.priority.value if isinstance(self.priority, Priority) else self.priority
        return data

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Card":
        priority = data.get("priority", Priority.MEDIUM)
        if isinstance(priority, str):
            priority = Priority(priority)
        return cls(
            id=data["id"],
            title=data["title"],
            description=data.get("description", ""),
            priority=priority,
            tags=data.get("tags", []),
            due_date=data.get("due_date"),
            column_id=data.get("column_id", ""),
            created_at=data.get("created_at", datetime.now().isoformat()),
            completed_at=data.get("completed_at"),
        )


@dataclass
class Column:
    id: str
    name: str
    order: int = 0
    is_completed: bool = False

    @classmethod
    def create(cls, name: str, order: int = 0, is_completed: bool = False) -> "Column":
        return cls(
            id=str(uuid.uuid4()),
            name=name,
            order=order,
            is_completed=is_completed,
        )

    def to_dict(self) -> Dict[str, Any]:
        return asdict(self)

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Column":
        return cls(
            id=data["id"],
            name=data["name"],
            order=data.get("order", 0),
            is_completed=data.get("is_completed", False),
        )


@dataclass
class Board:
    id: str
    name: str
    description: str = ""
    columns: List[Column] = field(default_factory=list)
    cards: List[Card] = field(default_factory=list)
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    archived: bool = False

    @classmethod
    def create(cls, name: str, description: str = "") -> "Board":
        return cls(
            id=str(uuid.uuid4()),
            name=name,
            description=description,
        )

    def get_column(self, column_id: str) -> Optional[Column]:
        for col in self.columns:
            if col.id == column_id:
                return col
        return None

    def get_column_by_name(self, name: str) -> Optional[Column]:
        for col in self.columns:
            if col.name.lower() == name.lower():
                return col
        return None

    def get_card(self, card_id: str) -> Optional[Card]:
        for card in self.cards:
            if card.id == card_id:
                return card
        return None

    def get_cards_in_column(self, column_id: str) -> List[Card]:
        return [card for card in self.cards if card.column_id == column_id]

    def sorted_columns(self) -> List[Column]:
        return sorted(self.columns, key=lambda c: c.order)

    def to_dict(self) -> Dict[str, Any]:
        return {
            "id": self.id,
            "name": self.name,
            "description": self.description,
            "columns": [col.to_dict() for col in self.columns],
            "cards": [card.to_dict() for card in self.cards],
            "created_at": self.created_at,
            "archived": self.archived,
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Board":
        return cls(
            id=data["id"],
            name=data["name"],
            description=data.get("description", ""),
            columns=[Column.from_dict(col) for col in data.get("columns", [])],
            cards=[Card.from_dict(card) for card in data.get("cards", [])],
            created_at=data.get("created_at", datetime.now().isoformat()),
            archived=data.get("archived", False),
        )
