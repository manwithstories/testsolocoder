from sqlalchemy.orm import Session
from app.models.models import Record, Member
from app.schemas.schemas import RecordCreate
from datetime import date


def create_record(db: Session, record_data: RecordCreate) -> Record:
    member = db.query(Member).filter(Member.id == record_data.member_id).first()
    if not member:
        raise ValueError("成员不存在")

    db_record = Record(
        member_id=record_data.member_id,
        type=record_data.type,
        category=record_data.category,
        amount=record_data.amount,
        date=record_data.date,
        description=record_data.description
    )
    db.add(db_record)
    db.commit()
    db.refresh(db_record)
    return db_record


def get_record_by_id(db: Session, record_id: int) -> Record:
    return db.query(Record).filter(Record.id == record_id).first()


def get_records_by_member(db: Session, member_id: int, start_date: date = None, end_date: date = None) -> list[Record]:
    query = db.query(Record).filter(Record.member_id == member_id)
    if start_date:
        query = query.filter(Record.date >= start_date)
    if end_date:
        query = query.filter(Record.date <= end_date)
    return query.order_by(Record.date.desc(), Record.created_at.desc()).all()


def get_records_by_family(db: Session, family_id: int, start_date: date = None, end_date: date = None) -> list[Record]:
    query = db.query(Record).join(Member).filter(Member.family_id == family_id)
    if start_date:
        query = query.filter(Record.date >= start_date)
    if end_date:
        query = query.filter(Record.date <= end_date)
    return query.order_by(Record.date.desc(), Record.created_at.desc()).all()


def update_record(db: Session, record_id: int, record_data: RecordCreate) -> Record:
    record = get_record_by_id(db, record_id)
    if not record:
        raise ValueError("记录不存在")

    member = db.query(Member).filter(Member.id == record_data.member_id).first()
    if not member:
        raise ValueError("成员不存在")

    record.member_id = record_data.member_id
    record.type = record_data.type
    record.category = record_data.category
    record.amount = record_data.amount
    record.date = record_data.date
    record.description = record_data.description

    db.commit()
    db.refresh(record)
    return record


def delete_record(db: Session, record_id: int) -> bool:
    record = get_record_by_id(db, record_id)
    if not record:
        raise ValueError("记录不存在")
    db.delete(record)
    db.commit()
    return True
