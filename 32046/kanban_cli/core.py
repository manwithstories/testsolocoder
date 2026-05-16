from datetime import datetime
from typing import List, Optional, Dict, Any, Tuple
from enum import Enum

from .models import Board, Column, Card, Priority
from .storage import Storage


class DuplicateResolution(str, Enum):
    SKIP = "skip"
    REPLACE = "replace"
    RENAME = "rename"


class KanbanManager:
    def __init__(self, storage: Optional[Storage] = None):
        self.storage = storage or Storage()

    def create_board(self, name: str, description: str = "",
                     default_columns: bool = True) -> Board:
        board = Board.create(name, description)
        if default_columns:
            board.columns = [
                Column.create("待办", 0),
                Column.create("进行中", 1),
                Column.create("已完成", 2, is_completed=True),
            ]
        self.storage.save_board(board)
        return board

    def list_boards(self, include_archived: bool = False) -> List[Board]:
        boards = self.storage.load_all_boards()
        if not include_archived:
            boards = [b for b in boards if not b.archived]
        return sorted(boards, key=lambda b: b.created_at, reverse=True)

    def get_board(self, board_id: str) -> Optional[Board]:
        return self.storage.load_board(board_id)

    def get_board_by_name(self, name: str) -> Optional[Board]:
        boards = self.storage.load_all_boards()
        for board in boards:
            if board.name.lower() == name.lower():
                return board
        return None

    def update_board(self, board: Board) -> None:
        self.storage.save_board(board)

    def archive_board(self, board_id: str) -> Optional[Board]:
        board = self.storage.load_board(board_id)
        if board:
            board.archived = True
            self.storage.save_board(board)
        return board

    def unarchive_board(self, board_id: str) -> Optional[Board]:
        board = self.storage.load_board(board_id)
        if board:
            board.archived = False
            self.storage.save_board(board)
        return board

    def delete_board(self, board_id: str) -> Optional[Board]:
        return self.storage.delete_board(board_id)

    def restore_board(self, board_id: str) -> Optional[Board]:
        return self.storage.restore_from_trash(board_id)

    def add_column(self, board_id: str, name: str, is_completed: bool = False) -> Optional[Column]:
        board = self.storage.load_board(board_id)
        if not board:
            return None
        max_order = max((c.order for c in board.columns), default=-1)
        column = Column.create(name, max_order + 1, is_completed)
        board.columns.append(column)
        self.storage.save_board(board)
        return column

    def rename_column(self, board_id: str, column_id: str, new_name: str) -> Optional[Column]:
        board = self.storage.load_board(board_id)
        if not board:
            return None
        column = board.get_column(column_id)
        if column:
            column.name = new_name
            self.storage.save_board(board)
        return column

    def delete_column(self, board_id: str, column_id: str, delete_cards: bool = False) -> bool:
        board = self.storage.load_board(board_id)
        if not board:
            return False
        column = board.get_column(column_id)
        if not column:
            return False
        board.columns = [c for c in board.columns if c.id != column_id]
        if delete_cards:
            board.cards = [c for c in board.cards if c.column_id != column_id]
        else:
            for card in board.cards:
                if card.column_id == column_id:
                    card.column_id = ""
        self.storage.save_board(board)
        return True

    def reorder_columns(self, board_id: str, column_order: List[str]) -> bool:
        board = self.storage.load_board(board_id)
        if not board:
            return False
        for i, col_id in enumerate(column_order):
            column = board.get_column(col_id)
            if column:
                column.order = i
        self.storage.save_board(board)
        return True

    def add_card(self, board_id: str, column_id: str, title: str,
                 description: str = "", priority: Priority = Priority.MEDIUM,
                 tags: Optional[List[str]] = None, due_date: Optional[str] = None) -> Optional[Card]:
        board = self.storage.load_board(board_id)
        if not board:
            return None
        if column_id and not board.get_column(column_id):
            return None
        card = Card.create(title, description, priority, tags, due_date)
        card.column_id = column_id
        board.cards.append(card)
        self.storage.save_board(board)
        return card

    def update_card(self, board_id: str, card_id: str, **kwargs) -> Optional[Card]:
        board = self.storage.load_board(board_id)
        if not board:
            return None
        card = board.get_card(card_id)
        if not card:
            return None
        for key, value in kwargs.items():
            if hasattr(card, key) and value is not None:
                setattr(card, key, value)
        self.storage.save_board(board)
        return card

    def move_card(self, board_id: str, card_id: str, new_column_id: str) -> Optional[Card]:
        board = self.storage.load_board(board_id)
        if not board:
            return None
        card = board.get_card(card_id)
        if not card:
            return None
        new_column = board.get_column(new_column_id)
        if not new_column:
            return None
        old_column = board.get_column(card.column_id)
        if new_column.is_completed and (old_column is None or not old_column.is_completed):
            card.completed_at = datetime.now().isoformat()
        elif old_column and old_column.is_completed and not new_column.is_completed:
            card.completed_at = None
        card.column_id = new_column_id
        self.storage.save_board(board)
        return card

    def delete_card(self, board_id: str, card_id: str) -> Optional[Card]:
        return self.storage.delete_card_permanently(board_id, card_id)

    def restore_card(self, card_id: str) -> Optional[Card]:
        return self.storage.restore_card_from_trash(card_id)

    def add_tag(self, board_id: str, card_id: str, tag: str) -> Optional[Card]:
        board = self.storage.load_board(board_id)
        if not board:
            return None
        card = board.get_card(card_id)
        if not card:
            return None
        if tag not in card.tags:
            card.tags.append(tag)
            self.storage.save_board(board)
        return card

    def remove_tag(self, board_id: str, card_id: str, tag: str) -> Optional[Card]:
        board = self.storage.load_board(board_id)
        if not board:
            return None
        card = board.get_card(card_id)
        if not card:
            return None
        if tag in card.tags:
            card.tags.remove(tag)
            self.storage.save_board(board)
        return card

    def search_cards(self, board_id: str, query: str = "",
                     tags: Optional[List[str]] = None,
                     priority: Optional[Priority] = None,
                     date_from: Optional[str] = None,
                     date_to: Optional[str] = None,
                     column_id: Optional[str] = None,
                     overdue_only: bool = False) -> List[Card]:
        board = self.storage.load_board(board_id)
        if not board:
            return []
        cards = board.cards
        if column_id:
            cards = [c for c in cards if c.column_id == column_id]
        if query:
            query_lower = query.lower()
            cards = [c for c in cards if query_lower in c.title.lower()
                     or query_lower in c.description.lower()]
        if tags:
            cards = [c for c in cards if any(tag in c.tags for tag in tags)]
        if priority:
            cards = [c for c in cards if c.priority == priority]
        if date_from:
            cards = [c for c in cards if c.due_date and c.due_date >= date_from]
        if date_to:
            cards = [c for c in cards if c.due_date and c.due_date <= date_to]
        if overdue_only:
            cards = [c for c in cards if c.is_overdue()]
        return cards

    def get_all_tags(self, board_id: str) -> List[str]:
        board = self.storage.load_board(board_id)
        if not board:
            return []
        all_tags = set()
        for card in board.cards:
            all_tags.update(card.tags)
        return sorted(all_tags)

    def export_board(self, board_id: str, export_path: str) -> bool:
        return self.storage.export_board(board_id, export_path)

    def import_board(self, import_path: str,
                     resolution: DuplicateResolution = DuplicateResolution.SKIP) -> Tuple[Optional[Board], str]:
        imported_board = self.storage.import_board(import_path)
        if not imported_board:
            return None, "导入失败：文件格式错误或不存在"
        existing = self.get_board_by_name(imported_board.name)
        if existing:
            if resolution == DuplicateResolution.SKIP:
                return None, f"导入取消：已存在同名看板 '{imported_board.name}'"
            elif resolution == DuplicateResolution.REPLACE:
                self.storage.delete_board(existing.id)
            elif resolution == DuplicateResolution.RENAME:
                counter = 1
                while True:
                    new_name = f"{imported_board.name} ({counter})"
                    if not self.storage.board_exists(new_name):
                        imported_board.name = new_name
                        break
                    counter += 1
        import uuid
        imported_board.id = str(uuid.uuid4())
        self.storage.save_board(imported_board)
        return imported_board, f"成功导入看板 '{imported_board.name}'"

    def get_stats(self, board_id: str) -> Dict[str, Any]:
        board = self.storage.load_board(board_id)
        if not board:
            return {}
        total_cards = len(board.cards)
        completed_cards = [c for c in board.cards if c.completed_at]
        completed_count = len(completed_cards)
        completion_rate = (completed_count / total_cards * 100) if total_cards > 0 else 0
        avg_days = 0
        if completed_count > 0:
            total_days = 0
            for card in completed_cards:
                try:
                    created = datetime.fromisoformat(card.created_at)
                    completed = datetime.fromisoformat(card.completed_at)
                    total_days += (completed - created).days
                except (ValueError, TypeError):
                    continue
            avg_days = total_days / completed_count
        column_distribution = {}
        for col in board.sorted_columns():
            col_cards = board.get_cards_in_column(col.id)
            column_distribution[col.name] = len(col_cards)
        unassigned = [c for c in board.cards if not c.column_id]
        if unassigned:
            column_distribution["未分配"] = len(unassigned)
        priority_distribution = {}
        for priority in Priority:
            count = len([c for c in board.cards if c.priority == priority])
            if count > 0:
                priority_distribution[priority.value] = count
        overdue_count = len([c for c in board.cards if c.is_overdue()])
        return {
            "total_cards": total_cards,
            "completed_cards": completed_count,
            "completion_rate": round(completion_rate, 2),
            "avg_completion_days": round(avg_days, 1),
            "column_distribution": column_distribution,
            "priority_distribution": priority_distribution,
            "overdue_cards": overdue_count,
        }
