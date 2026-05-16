import sqlite3
from datetime import datetime
from pathlib import Path
from typing import List, Optional

from track.models.transaction import Budget, Transaction, TransactionType


def get_db_path() -> Path:
    home = Path.home()
    track_dir = home / ".track"
    track_dir.mkdir(exist_ok=True)
    return track_dir / "track.db"


class Database:
    def __init__(self, db_path: Optional[Path] = None):
        self.db_path = db_path or get_db_path()
        self._init_db()

    def _get_connection(self) -> sqlite3.Connection:
        return sqlite3.connect(self.db_path)

    def _init_db(self):
        conn = self._get_connection()
        cursor = conn.cursor()
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS transactions (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                type TEXT NOT NULL,
                category TEXT NOT NULL,
                amount REAL NOT NULL,
                note TEXT,
                created_at TEXT NOT NULL
            )
            """
        )
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS budgets (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                category TEXT NOT NULL,
                amount REAL NOT NULL,
                month TEXT NOT NULL,
                UNIQUE(category, month)
            )
            """
        )
        conn.commit()
        conn.close()

    def add_transaction(self, transaction: Transaction) -> int:
        conn = self._get_connection()
        cursor = conn.cursor()
        cursor.execute(
            """
            INSERT INTO transactions (type, category, amount, note, created_at)
            VALUES (?, ?, ?, ?, ?)
            """,
            (
                transaction.type.value,
                transaction.category,
                transaction.amount,
                transaction.note,
                transaction.created_at.isoformat(),
            ),
        )
        conn.commit()
        transaction_id = cursor.lastrowid
        conn.close()
        return transaction_id

    def get_transactions(
        self,
        month: Optional[str] = None,
        category: Optional[str] = None,
        type: Optional[TransactionType] = None,
    ) -> List[Transaction]:
        conn = self._get_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()

        query = "SELECT * FROM transactions WHERE 1=1"
        params = []

        if month:
            query += " AND strftime('%Y-%m', created_at) = ?"
            params.append(month)
        if category:
            query += " AND category = ?"
            params.append(category)
        if type:
            query += " AND type = ?"
            params.append(type.value)

        query += " ORDER BY created_at DESC"

        cursor.execute(query, params)
        rows = cursor.fetchall()
        conn.close()

        return [
            Transaction(
                id=row["id"],
                type=TransactionType(row["type"]),
                category=row["category"],
                amount=row["amount"],
                note=row["note"],
                created_at=datetime.fromisoformat(row["created_at"]),
            )
            for row in rows
        ]

    def get_monthly_summary(self, month: str) -> dict:
        conn = self._get_connection()
        cursor = conn.cursor()
        cursor.execute(
            """
            SELECT category, type, SUM(amount) as total
            FROM transactions
            WHERE strftime('%Y-%m', created_at) = ?
            GROUP BY category, type
            """,
            (month,),
        )
        rows = cursor.fetchall()
        conn.close()

        summary = {"expenses": {}, "incomes": {}}
        for category, type_, total in rows:
            if type_ == TransactionType.EXPENSE.value:
                summary["expenses"][category] = total
            else:
                summary["incomes"][category] = total

        return summary

    def set_budget(self, budget: Budget) -> int:
        conn = self._get_connection()
        cursor = conn.cursor()
        cursor.execute(
            """
            INSERT INTO budgets (category, amount, month)
            VALUES (?, ?, ?)
            ON CONFLICT(category, month) DO UPDATE SET amount = excluded.amount
            """,
            (budget.category, budget.amount, budget.month),
        )
        conn.commit()
        budget_id = cursor.lastrowid
        conn.close()
        return budget_id

    def get_budgets(self, month: Optional[str] = None) -> List[Budget]:
        conn = self._get_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()

        if month:
            cursor.execute("SELECT * FROM budgets WHERE month = ?", (month,))
        else:
            cursor.execute("SELECT * FROM budgets")

        rows = cursor.fetchall()
        conn.close()

        return [
            Budget(
                id=row["id"],
                category=row["category"],
                amount=row["amount"],
                month=row["month"],
            )
            for row in rows
        ]

    def get_budget(self, category: str, month: str) -> Optional[Budget]:
        conn = self._get_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute(
            "SELECT * FROM budgets WHERE category = ? AND month = ?",
            (category, month),
        )
        row = cursor.fetchone()
        conn.close()

        if row:
            return Budget(
                id=row["id"],
                category=row["category"],
                amount=row["amount"],
                month=row["month"],
            )
        return None

    def delete_transaction(self, transaction_id: int) -> bool:
        conn = self._get_connection()
        cursor = conn.cursor()
        cursor.execute("DELETE FROM transactions WHERE id = ?", (transaction_id,))
        conn.commit()
        affected = cursor.rowcount > 0
        conn.close()
        return affected
