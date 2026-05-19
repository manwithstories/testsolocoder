from sqlalchemy.orm import Session
from sqlalchemy import and_
from app.models.models import Budget, Family, Member, Record
from app.schemas.schemas import BudgetCreate, BudgetUpdate
from datetime import date


def get_month_start_end(reference_date: date = None):
    if reference_date is None:
        reference_date = date.today()
    start = reference_date.replace(day=1)
    if start.month == 12:
        end = date(start.year + 1, 1, 1)
    else:
        end = date(start.year, start.month + 1, 1)
    return start, end


def calculate_category_spent(db: Session, family_id: int, category: str, reference_date: date = None):
    start_date, end_date = get_month_start_end(reference_date)

    spent = db.query(Record).join(Member).filter(
        Member.family_id == family_id,
        Record.type == "expense",
        Record.category == category,
        Record.date >= start_date,
        Record.date < end_date
    ).all()

    return sum(record.amount for record in spent)


def create_budget(db: Session, budget_data: BudgetCreate) -> Budget:
    family = db.query(Family).filter(Family.id == budget_data.family_id).first()
    if not family:
        raise ValueError("家庭不存在")

    existing = db.query(Budget).filter(
        and_(
            Budget.family_id == budget_data.family_id,
            Budget.category == budget_data.category
        )
    ).first()

    if existing:
        raise ValueError("该分类的预算已存在")

    db_budget = Budget(
        family_id=budget_data.family_id,
        category=budget_data.category,
        amount=budget_data.amount
    )
    db.add(db_budget)
    db.commit()
    db.refresh(db_budget)
    return db_budget


def get_budgets_by_family(db: Session, family_id: int, include_warning: bool = True) -> list:
    family = db.query(Family).filter(Family.id == family_id).first()
    if not family:
        raise ValueError("家庭不存在")

    budgets = db.query(Budget).filter(Budget.family_id == family_id).all()

    if not include_warning:
        return budgets

    result = []
    for budget in budgets:
        current_spent = calculate_category_spent(db, family_id, budget.category)
        warning = current_spent >= budget.amount * 0.8
        result.append({
            "id": budget.id,
            "family_id": budget.family_id,
            "category": budget.category,
            "amount": budget.amount,
            "created_at": budget.created_at,
            "current_spent": round(current_spent, 2),
            "warning": warning
        })
    return result


def get_budget_by_id(db: Session, budget_id: int) -> Budget:
    return db.query(Budget).filter(Budget.id == budget_id).first()


def update_budget(db: Session, budget_id: int, budget_data: BudgetUpdate) -> Budget:
    budget = get_budget_by_id(db, budget_id)
    if not budget:
        raise ValueError("预算不存在")

    budget.amount = budget_data.amount
    db.commit()
    db.refresh(budget)
    return budget


def delete_budget(db: Session, budget_id: int) -> bool:
    budget = get_budget_by_id(db, budget_id)
    if not budget:
        raise ValueError("预算不存在")
    db.delete(budget)
    db.commit()
    return True
