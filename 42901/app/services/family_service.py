import uuid
from sqlalchemy.orm import Session
from app.models.models import Family, Member
from app.schemas.schemas import FamilyCreate, MemberCreate


def generate_invite_code() -> str:
    return str(uuid.uuid4())[:8].upper()


def create_family(db: Session, family_data: FamilyCreate) -> Family:
    invite_code = generate_invite_code()
    while db.query(Family).filter(Family.invite_code == invite_code).first():
        invite_code = generate_invite_code()

    db_family = Family(name=family_data.name, invite_code=invite_code)
    db.add(db_family)
    db.commit()
    db.refresh(db_family)
    return db_family


def get_family_by_id(db: Session, family_id: int) -> Family:
    return db.query(Family).filter(Family.id == family_id).first()


def get_family_by_invite_code(db: Session, invite_code: str) -> Family:
    return db.query(Family).filter(Family.invite_code == invite_code.upper()).first()


def add_member(db: Session, member_data: MemberCreate) -> Member:
    family = get_family_by_invite_code(db, member_data.invite_code)
    if not family:
        raise ValueError("邀请码无效")

    member_count = db.query(Member).filter(Member.family_id == family.id).count()
    if member_count >= 8:
        raise ValueError("家庭成员数量已达上限（最多8人）")

    db_member = Member(name=member_data.name, family_id=family.id)
    db.add(db_member)
    db.commit()
    db.refresh(db_member)
    return db_member


def get_member_by_id(db: Session, member_id: int) -> Member:
    return db.query(Member).filter(Member.id == member_id).first()


def get_family_members(db: Session, family_id: int) -> list[Member]:
    return db.query(Member).filter(Member.family_id == family_id).all()


def remove_member(db: Session, member_id: int) -> bool:
    member = get_member_by_id(db, member_id)
    if not member:
        raise ValueError("成员不存在")
    db.delete(member)
    db.commit()
    return True
