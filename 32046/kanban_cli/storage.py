import json
import os
import shutil
from datetime import datetime, timedelta
from pathlib import Path
from typing import List, Optional, Dict, Any

from .models import Board, Card


class Storage:
    def __init__(self, data_dir: Optional[str] = None):
        if data_dir is None:
            data_dir = os.path.join(Path.home(), ".kanban-cli")
        self.data_dir = data_dir
        self.boards_file = os.path.join(self.data_dir, "boards.json")
        self.trash_dir = os.path.join(self.data_dir, "trash")
        self.card_trash_dir = os.path.join(self.data_dir, "card_trash")
        self._ensure_dirs()

    def _ensure_dirs(self) -> None:
        os.makedirs(self.data_dir, exist_ok=True)
        os.makedirs(self.trash_dir, exist_ok=True)
        os.makedirs(self.card_trash_dir, exist_ok=True)
        if not os.path.exists(self.boards_file):
            self._save_boards([])

    def _load_boards(self) -> List[Dict[str, Any]]:
        try:
            with open(self.boards_file, "r", encoding="utf-8") as f:
                return json.load(f)
        except (json.JSONDecodeError, FileNotFoundError):
            return []

    def _save_boards(self, boards_data: List[Dict[str, Any]]) -> None:
        with open(self.boards_file, "w", encoding="utf-8") as f:
            json.dump(boards_data, f, indent=2, ensure_ascii=False)

    def load_all_boards(self) -> List[Board]:
        data = self._load_boards()
        return [Board.from_dict(b) for b in data]

    def save_all_boards(self, boards: List[Board]) -> None:
        data = [b.to_dict() for b in boards]
        self._save_boards(data)

    def load_board(self, board_id: str) -> Optional[Board]:
        boards = self.load_all_boards()
        for board in boards:
            if board.id == board_id:
                return board
        return None

    def save_board(self, board: Board) -> None:
        boards = self.load_all_boards()
        for i, b in enumerate(boards):
            if b.id == board.id:
                boards[i] = board
                break
        else:
            boards.append(board)
        self.save_all_boards(boards)

    def delete_board(self, board_id: str) -> Optional[Board]:
        boards = self.load_all_boards()
        deleted = None
        for i, b in enumerate(boards):
            if b.id == board_id:
                deleted = boards.pop(i)
                break
        if deleted:
            self.save_all_boards(boards)
            self._move_to_trash(deleted)
        return deleted

    def _move_to_trash(self, board: Board) -> None:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        trash_file = os.path.join(self.trash_dir, f"{board.id}_{timestamp}.json")
        with open(trash_file, "w", encoding="utf-8") as f:
            json.dump(board.to_dict(), f, indent=2, ensure_ascii=False)

    def restore_from_trash(self, board_id: str) -> Optional[Board]:
        trash_files = sorted(os.listdir(self.trash_dir), reverse=True)
        for filename in trash_files:
            if filename.startswith(board_id):
                filepath = os.path.join(self.trash_dir, filename)
                try:
                    with open(filepath, "r", encoding="utf-8") as f:
                        data = json.load(f)
                    board = Board.from_dict(data)
                    boards = self.load_all_boards()
                    boards.append(board)
                    self.save_all_boards(boards)
                    os.remove(filepath)
                    return board
                except (json.JSONDecodeError, FileNotFoundError):
                    continue
        return None

    def cleanup_old_trash(self, days: int = 30) -> int:
        cutoff = datetime.now() - timedelta(days=days)
        count = 0
        for filename in os.listdir(self.trash_dir):
            filepath = os.path.join(self.trash_dir, filename)
            if os.path.getmtime(filepath) < cutoff.timestamp():
                os.remove(filepath)
                count += 1
        return count

    def export_board(self, board_id: str, export_path: str) -> bool:
        board = self.load_board(board_id)
        if not board:
            return False
        with open(export_path, "w", encoding="utf-8") as f:
            json.dump(board.to_dict(), f, indent=2, ensure_ascii=False)
        return True

    def import_board(self, import_path: str) -> Optional[Board]:
        try:
            with open(import_path, "r", encoding="utf-8") as f:
                data = json.load(f)
            return Board.from_dict(data)
        except (json.JSONDecodeError, FileNotFoundError, KeyError):
            return None

    def board_exists(self, name: str) -> bool:
        boards = self.load_all_boards()
        return any(b.name.lower() == name.lower() for b in boards)

    def _move_card_to_trash(self, card: Card, board_id: str) -> None:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        trash_file = os.path.join(self.card_trash_dir, f"{card.id}_{timestamp}.json")
        data = {
            "card": card.to_dict(),
            "board_id": board_id,
            "deleted_at": datetime.now().isoformat(),
        }
        with open(trash_file, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2, ensure_ascii=False)

    def restore_card_from_trash(self, card_id: str) -> Optional[Card]:
        trash_files = sorted(os.listdir(self.card_trash_dir), reverse=True)
        for filename in trash_files:
            if filename.startswith(card_id):
                filepath = os.path.join(self.card_trash_dir, filename)
                try:
                    with open(filepath, "r", encoding="utf-8") as f:
                        data = json.load(f)
                    card = Card.from_dict(data["card"])
                    board_id = data["board_id"]
                    board = self.load_board(board_id)
                    if not board:
                        return None
                    existing_ids = {c.id for c in board.cards}
                    if card.id not in existing_ids:
                        board.cards.append(card)
                        self.save_board(board)
                    os.remove(filepath)
                    return card
                except (json.JSONDecodeError, FileNotFoundError, KeyError):
                    continue
        return None

    def delete_card_permanently(self, board_id: str, card_id: str) -> Optional[Card]:
        board = self.load_board(board_id)
        if not board:
            return None
        deleted = None
        for i, c in enumerate(board.cards):
            if c.id == card_id:
                deleted = board.cards.pop(i)
                break
        if deleted:
            self._move_card_to_trash(deleted, board_id)
            self.save_board(board)
        return deleted

    def cleanup_old_card_trash(self, days: int = 30) -> int:
        cutoff = datetime.now() - timedelta(days=days)
        count = 0
        for filename in os.listdir(self.card_trash_dir):
            filepath = os.path.join(self.card_trash_dir, filename)
            if os.path.getmtime(filepath) < cutoff.timestamp():
                os.remove(filepath)
                count += 1
        return count
