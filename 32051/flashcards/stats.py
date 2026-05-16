from __future__ import annotations

from collections import defaultdict
from datetime import datetime, timedelta, timezone
from typing import Dict, List, Optional

from .models import Deck, FlashCard


class StatsCalculator:
    def __init__(self, deck: Deck):
        self.deck = deck

    def get_overall_stats(self) -> Dict:
        total_cards = len(self.deck.cards)
        due_cards = len(self.deck.get_due_cards())
        reviewed_cards = len([c for c in self.deck.cards if c.review_count > 0])
        total_reviews = sum(c.review_count for c in self.deck.cards)

        if total_reviews > 0:
            correct_reviews = sum(
                1 for r in self.deck.study_records if r.quality >= 3
            )
            accuracy = round(correct_reviews / total_reviews * 100, 1)
        else:
            accuracy = 0.0

        return {
            "total_cards": total_cards,
            "due_cards": due_cards,
            "reviewed_cards": reviewed_cards,
            "unreviewed_cards": total_cards - reviewed_cards,
            "total_reviews": total_reviews,
            "accuracy": accuracy,
        }

    def get_difficulty_breakdown(self) -> Dict[str, int]:
        breakdown = defaultdict(int)
        for card in self.deck.cards:
            breakdown[card.difficulty.value] += 1
        return dict(breakdown)

    def get_tag_breakdown(self) -> Dict[str, int]:
        breakdown = defaultdict(int)
        for card in self.deck.cards:
            for tag in card.tags:
                breakdown[tag] += 1
        return dict(sorted(breakdown.items(), key=lambda x: (-x[1], x[0])))

    def get_daily_activity(self, days: int = 30) -> List[Dict]:
        end_date = datetime.now(timezone.utc).date()
        start_date = end_date - timedelta(days=days - 1)

        activity = defaultdict(lambda: {"reviews": 0, "correct": 0})

        for record in self.deck.study_records:
            review_date = datetime.fromisoformat(record.reviewed_at).date()
            if start_date <= review_date <= end_date:
                date_str = review_date.isoformat()
                activity[date_str]["reviews"] += 1
                if record.quality >= 3:
                    activity[date_str]["correct"] += 1

        result = []
        current_date = start_date
        while current_date <= end_date:
            date_str = current_date.isoformat()
            day_data = activity.get(date_str, {"reviews": 0, "correct": 0})
            accuracy = (
                round(day_data["correct"] / day_data["reviews"] * 100, 1)
                if day_data["reviews"] > 0
                else 0.0
            )
            result.append({
                "date": date_str,
                "reviews": day_data["reviews"],
                "correct": day_data["correct"],
                "accuracy": accuracy,
            })
            current_date += timedelta(days=1)

        return result

    def get_streak_days(self) -> int:
        if not self.deck.study_records:
            return 0

        review_dates = set()
        for record in self.deck.study_records:
            review_date = datetime.fromisoformat(record.reviewed_at).date()
            review_dates.add(review_date)

        today = datetime.now(timezone.utc).date()
        streak = 0
        current_date = today

        while current_date in review_dates:
            streak += 1
            current_date -= timedelta(days=1)

        return streak

    def get_retention_rate(self) -> Dict[str, float]:
        intervals = {
            "1天": 1,
            "1周": 7,
            "1个月": 30,
            "3个月": 90,
        }

        retention = {}
        now = datetime.now(timezone.utc)

        for label, days in intervals.items():
            cutoff = now - timedelta(days=days)
            recent_cards = [
                c
                for c in self.deck.cards
                if c.last_reviewed_at
                and datetime.fromisoformat(c.last_reviewed_at) >= cutoff
            ]

            if recent_cards:
                retained = sum(1 for c in recent_cards if c.repetitions > 0)
                retention[label] = round(retained / len(recent_cards) * 100, 1)
            else:
                retention[label] = 0.0

        return retention

    def get_ease_factor_distribution(self) -> Dict[str, int]:
        ranges = [
            ("< 1.5", lambda x: x < 1.5),
            ("1.5-2.0", lambda x: 1.5 <= x < 2.0),
            ("2.0-2.5", lambda x: 2.0 <= x < 2.5),
            ("2.5-3.0", lambda x: 2.5 <= x < 3.0),
            (">= 3.0", lambda x: x >= 3.0),
        ]

        dist = defaultdict(int)
        for card in self.deck.cards:
            for label, check in ranges:
                if check(card.ease_factor):
                    dist[label] += 1
                    break

        return dict(dist)

    def get_tag_stats(self, tag: str) -> Optional[Dict]:
        tag_cards = [c for c in self.deck.cards if tag in c.tags]
        if not tag_cards:
            return None

        total_reviews = sum(c.review_count for c in tag_cards)
        tag_records = [
            r for r in self.deck.study_records
            if any(c.id == r.card_id for c in tag_cards)
        ]

        if tag_records:
            correct = sum(1 for r in tag_records if r.quality >= 3)
            accuracy = round(correct / len(tag_records) * 100, 1)
        else:
            accuracy = 0.0

        due_count = len([c for c in tag_cards if c.is_due()])

        return {
            "total_cards": len(tag_cards),
            "due_cards": due_count,
            "total_reviews": total_reviews,
            "accuracy": accuracy,
        }
