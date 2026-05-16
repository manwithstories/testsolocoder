from __future__ import annotations

import json
import os
import uuid
from datetime import datetime, timezone
from pathlib import Path
from typing import Dict, List, Optional

from .models import Deck, FlashCard


class Storage:
    def __init__(self, data_dir: Optional[str] = None):
        if data_dir is None:
            home = Path.home()
            data_dir = os.path.join(home, ".flashcards")
        self.data_dir = Path(data_dir)
        self.decks_dir = self.data_dir / "decks"
        self.data_file = self.data_dir / "decks.json"
        self._ensure_dirs()

    def _ensure_dirs(self) -> None:
        self.data_dir.mkdir(parents=True, exist_ok=True)
        self.decks_dir.mkdir(parents=True, exist_ok=True)

    def _read_all(self) -> Dict[str, Deck]:
        if not self.data_file.exists():
            return {}
        try:
            with open(self.data_file, "r", encoding="utf-8") as f:
                data = json.load(f)
            decks = {}
            for deck_data in data.get("decks", []):
                deck = Deck.from_dict(deck_data)
                decks[deck.id] = deck
            return decks
        except json.JSONDecodeError:
            return {}

    def _write_all(self, decks: Dict[str, Deck]) -> None:
        data = {"decks": [deck.to_dict() for deck in decks.values()]}
        with open(self.data_file, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=2)

    def list_decks(self) -> List[Deck]:
        decks = self._read_all()
        return sorted(decks.values(), key=lambda d: d.updated_at, reverse=True)

    def get_deck(self, deck_id: str) -> Optional[Deck]:
        decks = self._read_all()
        return decks.get(deck_id)

    def get_deck_by_name(self, name: str) -> Optional[Deck]:
        decks = self._read_all()
        for deck in decks.values():
            if deck.name.lower() == name.lower():
                return deck
        return None

    def create_deck(self, name: str) -> Deck:
        existing = self.get_deck_by_name(name)
        if existing:
            raise ValueError(f"牌组 '{name}' 已存在")
        deck = Deck(name=name)
        decks = self._read_all()
        decks[deck.id] = deck
        self._write_all(decks)
        return deck

    def rename_deck(self, deck_id: str, new_name: str) -> Optional[Deck]:
        decks = self._read_all()
        deck = decks.get(deck_id)
        if not deck:
            return None
        existing = self.get_deck_by_name(new_name)
        if existing and existing.id != deck_id:
            raise ValueError(f"牌组 '{new_name}' 已存在")
        deck.name = new_name
        deck.touch()
        self._write_all(decks)
        return deck

    def delete_deck(self, deck_id: str) -> bool:
        decks = self._read_all()
        if deck_id not in decks:
            return False
        del decks[deck_id]
        self._write_all(decks)
        return True

    def add_card(self, deck_id: str, card: FlashCard) -> Optional[FlashCard]:
        decks = self._read_all()
        deck = decks.get(deck_id)
        if not deck:
            return None
        deck.add_card(card)
        self._write_all(decks)
        return card

    def update_card(self, deck_id: str, card: FlashCard) -> Optional[FlashCard]:
        decks = self._read_all()
        deck = decks.get(deck_id)
        if not deck:
            return None
        for i, c in enumerate(deck.cards):
            if c.id == card.id:
                card.touch()
                deck.cards[i] = card
                deck.touch()
                self._write_all(decks)
                return card
        return None

    def delete_card(self, deck_id: str, card_id: str) -> bool:
        decks = self._read_all()
        deck = decks.get(deck_id)
        if not deck:
            return False
        if deck.remove_card(card_id):
            self._write_all(decks)
            return True
        return False

    def save_deck(self, deck: Deck) -> None:
        decks = self._read_all()
        deck.touch()
        decks[deck.id] = deck
        self._write_all(decks)

    def export_deck_to_json(self, deck_id: str, filepath: str) -> bool:
        deck = self.get_deck(deck_id)
        if not deck:
            return False
        cards_data = [card.to_dict() for card in deck.cards]
        with open(filepath, "w", encoding="utf-8") as f:
            json.dump(cards_data, f, ensure_ascii=False, indent=2)
        return True

    def import_cards_from_json(self, deck_id: str, filepath: str, overwrite: bool = False) -> Dict[str, int]:
        from .models import Difficulty

        deck = self.get_deck(deck_id)
        if not deck:
            return {"added": 0, "updated": 0, "skipped": 0}
        try:
            with open(filepath, "r", encoding="utf-8") as f:
                cards_data = json.load(f)
        except (json.JSONDecodeError, FileNotFoundError):
            return {"added": 0, "updated": 0, "skipped": 0}

        added = 0
        updated = 0
        skipped = 0

        existing_fronts = {card.front: card for card in deck.cards}

        for card_data in cards_data:
            front = card_data.get("front", "").strip()
            back = card_data.get("back", "").strip()
            if not front or not back:
                skipped += 1
                continue
            if front in existing_fronts:
                if overwrite:
                    existing_card = existing_fronts[front]
                    existing_card.back = back
                    existing_card.tags = card_data.get("tags", existing_card.tags)
                    diff_value = card_data.get("difficulty")
                    if diff_value:
                        existing_card.difficulty = Difficulty(diff_value)
                    existing_card.touch()
                    updated += 1
                else:
                    skipped += 1
            else:
                new_card = FlashCard.from_dict(card_data)
                deck.add_card(new_card)
                added += 1

        self.save_deck(deck)
        return {"added": added, "updated": updated, "skipped": skipped}

    def export_deck_to_csv(self, deck_id: str, filepath: str) -> bool:
        import csv

        deck = self.get_deck(deck_id)
        if not deck:
            return False
        with open(filepath, "w", encoding="utf-8", newline="") as f:
            writer = csv.writer(f)
            writer.writerow(["front", "back", "tags", "difficulty"])
            for card in deck.cards:
                tags_str = ",".join(card.tags)
                writer.writerow([card.front, card.back, tags_str, card.difficulty.value])
        return True

    def import_cards_from_csv(self, deck_id: str, filepath: str, overwrite: bool = False) -> Dict[str, int]:
        import csv

        deck = self.get_deck(deck_id)
        if not deck:
            return {"added": 0, "updated": 0, "skipped": 0}
        try:
            with open(filepath, "r", encoding="utf-8") as f:
                reader = csv.DictReader(f)
                cards_data = []
                for row in reader:
                    tags = []
                    if row.get("tags"):
                        tags = [t.strip() for t in row["tags"].split(",") if t.strip()]
                    cards_data.append({
                        "front": row.get("front", ""),
                        "back": row.get("back", ""),
                        "tags": tags,
                        "difficulty": row.get("difficulty", "medium"),
                    })
        except FileNotFoundError:
            return {"added": 0, "updated": 0, "skipped": 0}

        added = 0
        updated = 0
        skipped = 0

        existing_fronts = {card.front: card for card in deck.cards}

        for card_data in cards_data:
            front = card_data.get("front", "").strip()
            back = card_data.get("back", "").strip()
            if not front or not back:
                skipped += 1
                continue
            if front in existing_fronts:
                if overwrite:
                    from .models import Difficulty
                    existing_card = existing_fronts[front]
                    existing_card.back = back
                    existing_card.tags = card_data.get("tags", existing_card.tags)
                    diff_value = card_data.get("difficulty")
                    if diff_value:
                        existing_card.difficulty = Difficulty(diff_value)
                    existing_card.touch()
                    updated += 1
                else:
                    skipped += 1
            else:
                from .models import Difficulty
                difficulty = Difficulty(card_data.get("difficulty", "medium"))
                new_card = FlashCard(
                    front=front,
                    back=back,
                    tags=card_data.get("tags", []),
                    difficulty=difficulty,
                )
                deck.add_card(new_card)
                added += 1

        self.save_deck(deck)
        return {"added": added, "updated": updated, "skipped": skipped}
