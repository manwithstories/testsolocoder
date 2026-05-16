from dataclasses import dataclass
from datetime import datetime
from enum import Enum
from typing import Optional


class TransactionType(str, Enum):
    EXPENSE = "expense"
    INCOME = "income"


@dataclass
class Transaction:
    id: Optional[int]
    type: TransactionType
    category: str
    amount: float
    note: str
    created_at: datetime

    @classmethod
    def create(
        cls,
        type: TransactionType,
        category: str,
        amount: float,
        note: str = "",
    ) -> "Transaction":
        return cls(
            id=None,
            type=type,
            category=category,
            amount=amount,
            note=note,
            created_at=datetime.now(),
        )


@dataclass
class Budget:
    id: Optional[int]
    category: str
    amount: float
    month: str

    @classmethod
    def create(cls, category: str, amount: float, month: str) -> "Budget":
        return cls(id=None, category=category, amount=amount, month=month)
