from sqlalchemy.orm import Session
from app.models.models import Member, Record
from datetime import date, timedelta
from collections import defaultdict


def get_week_range(reference_date: date = None):
    if reference_date is None:
        reference_date = date.today()
    start = reference_date - timedelta(days=reference_date.weekday())
    end = start + timedelta(days=6)
    return start, end


def get_month_range(reference_date: date = None):
    if reference_date is None:
        reference_date = date.today()
    start = reference_date.replace(day=1)
    if start.month == 12:
        end = date(start.year + 1, 1, 1) - timedelta(days=1)
    else:
        end = date(start.year, start.month + 1, 1) - timedelta(days=1)
    return start, end


def calculate_category_summary(records, total_amount):
    category_totals = defaultdict(float)
    for record in records:
        category_totals[record.category] += record.amount

    summaries = []
    for category, amount in category_totals.items():
        percentage = (amount / total_amount * 100) if total_amount > 0 else 0
        summaries.append({
            "category": category,
            "total_amount": round(amount, 2),
            "percentage": round(percentage, 2)
        })
    summaries.sort(key=lambda x: x["total_amount"], reverse=True)
    return summaries


def get_family_statistics(db: Session, family_id: int, period: str, reference_date: date = None):
    if period == "week":
        start_date, end_date = get_week_range(reference_date)
    elif period == "month":
        start_date, end_date = get_month_range(reference_date)
    else:
        raise ValueError("无效的周期类型，只能是week或month")

    members = db.query(Member).filter(Member.family_id == family_id).all()
    if not members:
        raise ValueError("家庭不存在或没有成员")

    family = members[0].family

    member_statistics = []
    for member in members:
        records = db.query(Record).filter(
            Record.member_id == member.id,
            Record.date >= start_date,
            Record.date <= end_date
        ).all()

        income_records = [r for r in records if r.type == "income"]
        expense_records = [r for r in records if r.type == "expense"]

        total_income = sum(r.amount for r in income_records)
        total_expense = sum(r.amount for r in expense_records)

        income_by_category = calculate_category_summary(income_records, total_income)
        expense_by_category = calculate_category_summary(expense_records, total_expense)

        member_statistics.append({
            "member_id": member.id,
            "member_name": member.name,
            "total_income": round(total_income, 2),
            "total_expense": round(total_expense, 2),
            "net_balance": round(total_income - total_expense, 2),
            "income_by_category": income_by_category,
            "expense_by_category": expense_by_category
        })

    return {
        "period": period,
        "start_date": start_date,
        "end_date": end_date,
        "family_id": family.id,
        "family_name": family.name,
        "member_statistics": member_statistics
    }
