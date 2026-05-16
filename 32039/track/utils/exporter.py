import csv
from pathlib import Path
from typing import List

from track.models.transaction import Transaction


def export_transactions_to_csv(transactions: List[Transaction], output_path: Path):
    with open(output_path, "w", newline="", encoding="utf-8-sig") as f:
        writer = csv.writer(f)
        writer.writerow(["ID", "类型", "分类", "金额", "备注", "时间"])

        for t in transactions:
            type_label = "支出" if t.type.value == "expense" else "收入"
            writer.writerow([
                t.id,
                type_label,
                t.category,
                t.amount,
                t.note,
                t.created_at.strftime("%Y-%m-%d %H:%M:%S"),
            ])
