from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List
from app.database import get_db
from app.schemas.schemas import FamilyCreate, FamilyResponse, MemberCreate, MemberResponse
from app.services import family_service

router = APIRouter(prefix="/api/family", tags=["家庭与成员管理"])


@router.post("", response_model=FamilyResponse, summary="创建家庭")
def create_family(family_data: FamilyCreate, db: Session = Depends(get_db)):
    return family_service.create_family(db, family_data)


@router.get("/{family_id}", response_model=FamilyResponse, summary="获取家庭信息")
def get_family(family_id: int, db: Session = Depends(get_db)):
    family = family_service.get_family_by_id(db, family_id)
    if not family:
        raise HTTPException(status_code=404, detail="家庭不存在")
    return family


@router.post("/members", response_model=MemberResponse, summary="成员加入家庭")
def add_member(member_data: MemberCreate, db: Session = Depends(get_db)):
    try:
        return family_service.add_member(db, member_data)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.get("/{family_id}/members", response_model=List[MemberResponse], summary="获取家庭所有成员")
def get_family_members(family_id: int, db: Session = Depends(get_db)):
    family = family_service.get_family_by_id(db, family_id)
    if not family:
        raise HTTPException(status_code=404, detail="家庭不存在")
    return family_service.get_family_members(db, family_id)


@router.get("/members/{member_id}", response_model=MemberResponse, summary="获取成员信息")
def get_member(member_id: int, db: Session = Depends(get_db)):
    member = family_service.get_member_by_id(db, member_id)
    if not member:
        raise HTTPException(status_code=404, detail="成员不存在")
    return member


@router.delete("/members/{member_id}", summary="移除成员")
def remove_member(member_id: int, db: Session = Depends(get_db)):
    try:
        family_service.remove_member(db, member_id)
        return {"message": "成员已移除"}
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))
