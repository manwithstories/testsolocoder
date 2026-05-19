from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List
from app.database import get_db
from app.schemas.schemas import BudgetCreate, BudgetUpdate, BudgetResponse, BudgetWithWarning
from app.services import budget_service, family_service

router = APIRouter(prefix="/api/budgets", tags=["预算管理"])


@router.post("", response_model=BudgetResponse, summary="设置分类预算")
def create_budget(budget_data: BudgetCreate, db: Session = Depends(get_db)):
    try:
        return budget_service.create_budget(db, budget_data)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.get("/family/{family_id}", response_model=List[BudgetWithWarning], summary="获取家庭所有预算及预警状态")
def get_family_budgets(family_id: int, db: Session = Depends(get_db)):
    family = family_service.get_family_by_id(db, family_id)
    if not family:
        raise HTTPException(status_code=404, detail="家庭不存在")
    try:
        return budget_service.get_budgets_by_family(db, family_id)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.get("/{budget_id}", response_model=BudgetWithWarning, summary="获取单条预算详情")
def get_budget(budget_id: int, db: Session = Depends(get_db)):
    budget = budget_service.get_budget_by_id(db, budget_id)
    if not budget:
        raise HTTPException(status_code=404, detail="预算不存在")

    current_spent = budget_service.calculate_category_spent(db, budget.family_id, budget.category)
    warning = current_spent >= budget.amount * 0.8
    return {
        "id": budget.id,
        "family_id": budget.family_id,
        "category": budget.category,
        "amount": budget.amount,
        "created_at": budget.created_at,
        "current_spent": round(current_spent, 2),
        "warning": warning
    }


@router.put("/{budget_id}", response_model=BudgetResponse, summary="更新预算金额")
def update_budget(budget_id: int, budget_data: BudgetUpdate, db: Session = Depends(get_db)):
    try:
        return budget_service.update_budget(db, budget_id, budget_data)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.delete("/{budget_id}", summary="删除预算")
def delete_budget(budget_id: int, db: Session = Depends(get_db)):
    try:
        budget_service.delete_budget(db, budget_id)
        return {"message": "预算已删除"}
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))
