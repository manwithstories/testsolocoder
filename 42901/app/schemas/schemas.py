from pydantic import BaseModel, Field, field_validator
from datetime import date, datetime
from typing import Optional, List
from enum import Enum


class RecordType(str, Enum):
    INCOME = "income"
    EXPENSE = "expense"


class FamilyCreate(BaseModel):
    name: str = Field(..., min_length=1, max_length=50)


class FamilyResponse(BaseModel):
    id: int
    name: str
    invite_code: str
    created_at: datetime

    class Config:
        from_attributes = True


class MemberCreate(BaseModel):
    name: str = Field(..., min_length=1, max_length=50)
    invite_code: str = Field(..., min_length=1)


class MemberResponse(BaseModel):
    id: int
    name: str
    family_id: int
    created_at: datetime

    class Config:
        from_attributes = True


class MemberWithFamily(BaseModel):
    id: int
    name: str
    family_id: int
    family_name: str
    created_at: datetime


class RecordCreate(BaseModel):
    member_id: int
    type: RecordType
    category: str = Field(..., min_length=1, max_length=30)
    amount: float = Field(..., gt=0)
    date: date
    description: Optional[str] = None

    @field_validator('date')
    def date_not_in_future(cls, v):
        if v > date.today():
            raise ValueError('日期不能选择未来的日期')
        return v


class RecordResponse(BaseModel):
    id: int
    member_id: int
    member_name: str
    type: str
    category: str
    amount: float
    date: date
    description: Optional[str]
    created_at: datetime


class BudgetCreate(BaseModel):
    family_id: int
    category: str = Field(..., min_length=1, max_length=30)
    amount: float = Field(..., gt=0)


class BudgetUpdate(BaseModel):
    amount: float = Field(..., gt=0)


class BudgetResponse(BaseModel):
    id: int
    family_id: int
    category: str
    amount: float
    created_at: datetime

    class Config:
        from_attributes = True


class BudgetWithWarning(BudgetResponse):
    current_spent: float
    warning: bool


class CategorySummary(BaseModel):
    category: str
    total_amount: float
    percentage: float


class MemberStatistics(BaseModel):
    member_id: int
    member_name: str
    total_income: float
    total_expense: float
    net_balance: float
    income_by_category: List[CategorySummary]
    expense_by_category: List[CategorySummary]


class StatisticsResponse(BaseModel):
    period: str
    start_date: date
    end_date: date
    family_id: int
    family_name: str
    member_statistics: List[MemberStatistics]
