from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from typing import List, Optional
from datetime import date
from app.database import get_db
from app.schemas.schemas import RecordCreate, RecordResponse
from app.services import record_service, family_service

router = APIRouter(prefix="/api/records", tags=["记账管理"])


def format_record_response(record):
    return {
        "id": record.id,
        "member_id": record.member_id,
        "member_name": record.member.name,
        "type": record.type,
        "category": record.category,
        "amount": record.amount,
        "date": record.date,
        "description": record.description,
        "created_at": record.created_at
    }


@router.post("", response_model=RecordResponse, summary="添加记账记录")
def create_record(record_data: RecordCreate, db: Session = Depends(get_db)):
    try:
        record = record_service.create_record(db, record_data)
        return format_record_response(record)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.get("/{record_id}", response_model=RecordResponse, summary="获取单条记录")
def get_record(record_id: int, db: Session = Depends(get_db)):
    record = record_service.get_record_by_id(db, record_id)
    if not record:
        raise HTTPException(status_code=404, detail="记录不存在")
    return format_record_response(record)


@router.get("/member/{member_id}", response_model=List[RecordResponse], summary="获取成员的记账记录")
def get_member_records(
    member_id: int,
    start_date: Optional[date] = Query(None, description="开始日期"),
    end_date: Optional[date] = Query(None, description="结束日期"),
    db: Session = Depends(get_db)
):
    member = family_service.get_member_by_id(db, member_id)
    if not member:
        raise HTTPException(status_code=404, detail="成员不存在")
    records = record_service.get_records_by_member(db, member_id, start_date, end_date)
    return [format_record_response(r) for r in records]


@router.get("/family/{family_id}", response_model=List[RecordResponse], summary="获取家庭的记账记录")
def get_family_records(
    family_id: int,
    start_date: Optional[date] = Query(None, description="开始日期"),
    end_date: Optional[date] = Query(None, description="结束日期"),
    db: Session = Depends(get_db)
):
    family = family_service.get_family_by_id(db, family_id)
    if not family:
        raise HTTPException(status_code=404, detail="家庭不存在")
    records = record_service.get_records_by_family(db, family_id, start_date, end_date)
    return [format_record_response(r) for r in records]


@router.put("/{record_id}", response_model=RecordResponse, summary="更新记账记录")
def update_record(record_id: int, record_data: RecordCreate, db: Session = Depends(get_db)):
    try:
        record = record_service.update_record(db, record_id, record_data)
        return format_record_response(record)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.delete("/{record_id}", summary="删除记账记录")
def delete_record(record_id: int, db: Session = Depends(get_db)):
    try:
        record_service.delete_record(db, record_id)
        return {"message": "记录已删除"}
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))
