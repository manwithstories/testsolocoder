from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from typing import Optional
from datetime import date
from app.database import get_db
from app.schemas.schemas import StatisticsResponse
from app.services import statistics_service, family_service

router = APIRouter(prefix="/api/statistics", tags=["统计分析"])


@router.get("/family/{family_id}", response_model=StatisticsResponse, summary="获取家庭收支统计")
def get_family_statistics(
    family_id: int,
    period: str = Query(..., description="统计周期: week(按周) 或 month(按月)"),
    reference_date: Optional[date] = Query(None, description="参考日期，默认为今天"),
    db: Session = Depends(get_db)
):
    family = family_service.get_family_by_id(db, family_id)
    if not family:
        raise HTTPException(status_code=404, detail="家庭不存在")

    try:
        return statistics_service.get_family_statistics(db, family_id, period, reference_date)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
