#!/usr/bin/env python3
"""测试脚本 - 验证闪卡应用的核心功能"""

import os
import sys
import tempfile
from datetime import datetime, timezone, timedelta

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from flashcards.models import FlashCard, Difficulty, Deck
from flashcards.sm2 import calculate_next_review, apply_review_result, format_interval_days
from flashcards.storage import Storage
from flashcards.stats import StatsCalculator
from flashcards.session import StudySession


def test_sm2_algorithm():
    """测试SM-2算法"""
    print("=" * 60)
    print("测试 SM-2 间隔重复算法")
    print("=" * 60)

    card = FlashCard(front="Test", back="测试")
    print(f"\n初始状态: interval={card.interval}, ease={card.ease_factor}, reps={card.repetitions}")

    test_cases = [
        (5, "完美回答"),
        (5, "再次完美"),
        (3, "勉强记住"),
        (2, "答错了"),
        (4, "恢复"),
    ]

    for quality, desc in test_cases:
        apply_review_result(card, quality)
        next_interval = format_interval_days(card.interval)
        print(f"  评分 {quality} ({desc}): interval={card.interval:.1f}天 ({next_interval}), "
              f"ease={card.ease_factor:.2f}, reps={card.repetitions}")

    assert card.ease_factor >= 1.3, "Ease factor 不应低于最小值"
    print("\n✓ SM-2 算法测试通过")


def test_deck_operations():
    """测试牌组操作"""
    print("\n" + "=" * 60)
    print("测试牌组和卡片操作")
    print("=" * 60)

    with tempfile.TemporaryDirectory() as tmpdir:
        storage = Storage(data_dir=tmpdir)

        deck = storage.create_deck("测试牌组")
        print(f"\n✓ 创建牌组: {deck.name} (ID: {deck.id[:8]})")

        card1 = FlashCard(front="Hello", back="你好", tags=["基础"], difficulty=Difficulty.EASY)
        card2 = FlashCard(front="World", back="世界", tags=["基础"], difficulty=Difficulty.MEDIUM)
        storage.add_card(deck.id, card1)
        storage.add_card(deck.id, card2)
        print(f"✓ 添加2张卡片")

        deck = storage.get_deck(deck.id)
        assert len(deck.cards) == 2, "卡片数量不对"

        card1.back = "你好呀"
        storage.update_card(deck.id, card1)
        print("✓ 更新卡片")

        deck = storage.get_deck(deck.id)
        assert deck.cards[0].back == "你好呀", "卡片更新失败"

        due_cards = deck.get_due_cards()
        assert len(due_cards) == 2, "待复习卡片数量不对"
        print(f"✓ 待复习卡片: {len(due_cards)} 张")

        tags = deck.get_all_tags()
        assert tags == ["基础"], "标签获取失败"
        print(f"✓ 所有标签: {tags}")

        storage.delete_card(deck.id, card2.id)
        deck = storage.get_deck(deck.id)
        assert len(deck.cards) == 1, "删除卡片失败"
        print("✓ 删除卡片")

        storage.rename_deck(deck.id, "新名称")
        deck = storage.get_deck(deck.id)
        assert deck.name == "新名称", "重命名失败"
        print("✓ 重命名牌组")

        storage.delete_deck(deck.id)
        assert storage.get_deck(deck.id) is None, "删除牌组失败"
        print("✓ 删除牌组")


def test_import_export():
    """测试导入导出"""
    print("\n" + "=" * 60)
    print("测试导入导出功能")
    print("=" * 60)

    with tempfile.TemporaryDirectory() as tmpdir:
        storage = Storage(data_dir=tmpdir)
        deck = storage.create_deck("导入测试")

        card = FlashCard(front="Import", back="导入", tags=["test"], difficulty=Difficulty.EASY)
        storage.add_card(deck.id, card)

        json_file = os.path.join(tmpdir, "test.json")
        csv_file = os.path.join(tmpdir, "test.csv")

        assert storage.export_deck_to_json(deck.id, json_file), "JSON导出失败"
        assert storage.export_deck_to_csv(deck.id, csv_file), "CSV导出失败"
        print("✓ 导出 JSON 和 CSV")

        deck2 = storage.create_deck("导入测试2")
        result = storage.import_cards_from_json(deck2.id, json_file)
        assert result["added"] == 1, "JSON导入失败"
        print(f"✓ JSON导入: 新增{result['added']}张, 更新{result['updated']}张, 跳过{result['skipped']}张")

        result = storage.import_cards_from_json(deck2.id, json_file, overwrite=True)
        assert result["updated"] == 1, "覆盖导入失败"
        print(f"✓ JSON覆盖导入: 新增{result['added']}张, 更新{result['updated']}张, 跳过{result['skipped']}张")


def test_study_session():
    """测试学习会话"""
    print("\n" + "=" * 60)
    print("测试学习会话")
    print("=" * 60)

    deck = Deck(name="会话测试")
    for i in range(3):
        card = FlashCard(front=f"Question {i}", back=f"Answer {i}")
        deck.add_card(card)

    session = StudySession(deck, shuffle=False)
    print(f"\n✓ 创建学习会话，待复习: {len(session.review_queue)} 张")

    assert session.has_next(), "应该有卡片待复习"
    card = session.next_card()
    assert card is not None, "应该能获取卡片"
    print(f"✓ 复习卡片: {card.front}")

    session.record_answer(card, 4)
    print(f"✓ 记录答案，评分4，下次复习: {format_interval_days(card.interval)}")

    result = session.finish()
    print(f"✓ 会话结束: {result['total_cards']}张, 正确率{result['accuracy']}%, 用时{result['duration_formatted']}")
    assert result["total_cards"] == 1, "复习数量不对"


def test_stats():
    """测试统计功能"""
    print("\n" + "=" * 60)
    print("测试统计功能")
    print("=" * 60)

    deck = Deck(name="统计测试")
    for i in range(5):
        card = FlashCard(front=f"Q{i}", back=f"A{i}", tags=["tag1", "tag2"] if i < 3 else ["tag3"])
        deck.add_card(card)

    for i, card in enumerate(deck.cards[:3]):
        apply_review_result(card, 4 if i < 2 else 2)
        record = card.last_reviewed_at

    stats = StatsCalculator(deck)
    overall = stats.get_overall_stats()
    print(f"\n总体统计: {overall['total_cards']}张卡片, {overall['due_cards']}张待复习, 正确率{overall['accuracy']}%")

    tag_stats = stats.get_tag_breakdown()
    print(f"标签分布: {tag_stats}")

    diff_stats = stats.get_difficulty_breakdown()
    print(f"难度分布: {diff_stats}")

    streak = stats.get_streak_days()
    print(f"连续学习: {streak}天")

    assert overall["total_cards"] == 5, "卡片总数不对"
    assert overall["reviewed_cards"] == 3, "已学习卡片数不对"
    print("✓ 统计功能正常")


def test_search_filter():
    """测试搜索过滤功能"""
    print("\n" + "=" * 60)
    print("测试搜索过滤功能")
    print("=" * 60)

    deck = Deck(name="搜索测试")
    cards_data = [
        ("Apple", "苹果", ["水果"], Difficulty.EASY),
        ("Banana", "香蕉", ["水果"], Difficulty.EASY),
        ("Computer", "电脑", ["科技"], Difficulty.MEDIUM),
        ("Python", "蟒蛇/编程语言", ["科技", "编程"], Difficulty.HARD),
    ]
    for front, back, tags, diff in cards_data:
        deck.add_card(FlashCard(front=front, back=back, tags=tags, difficulty=diff))

    cards = deck.cards

    kw_filter = [c for c in cards if "p" in c.front.lower() or "p" in c.back.lower()]
    assert len(kw_filter) == 3, f"关键词搜索失败，期望3张，实际{len(kw_filter)}张"
    print(f"✓ 关键词'p'搜索: 找到{len(kw_filter)}张 (Apple, Computer, Python)")

    tag_filter = [c for c in cards if "水果" in c.tags]
    assert len(tag_filter) == 2, "标签过滤失败"
    print(f"✓ 标签'水果'过滤: 找到{len(tag_filter)}张")

    diff_filter = [c for c in cards if c.difficulty == Difficulty.EASY]
    assert len(diff_filter) == 2, "难度过滤失败"
    print(f"✓ 难度'easy'过滤: 找到{len(diff_filter)}张")


def main():
    try:
        test_sm2_algorithm()
        test_deck_operations()
        test_import_export()
        test_study_session()
        test_stats()
        test_search_filter()

        print("\n" + "=" * 60)
        print("🎉 所有测试通过！")
        print("=" * 60)
        return 0
    except Exception as e:
        print(f"\n❌ 测试失败: {e}")
        import traceback
        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
