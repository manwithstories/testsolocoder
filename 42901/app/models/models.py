from sqlalchemy import Column, Integer, String, Float, Date, ForeignKey, DateTime, UniqueConstraint
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime

Base = declarative_base()


class Family(Base):
    __tablename__ = "families"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, index=True)
    invite_code = Column(String, unique=True, index=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    members = relationship("Member", back_populates="family", cascade="all, delete-orphan")
    budgets = relationship("Budget", back_populates="family", cascade="all, delete-orphan")


class Member(Base):
    __tablename__ = "members"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, index=True)
    family_id = Column(Integer, ForeignKey("families.id"))
    created_at = Column(DateTime, default=datetime.utcnow)

    family = relationship("Family", back_populates="members")
    records = relationship("Record", back_populates="member", cascade="all, delete-orphan")


class Record(Base):
    __tablename__ = "records"

    id = Column(Integer, primary_key=True, index=True)
    member_id = Column(Integer, ForeignKey("members.id"))
    type = Column(String)  # income or expense
    category = Column(String, index=True)
    amount = Column(Float)
    date = Column(Date)
    description = Column(String, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    member = relationship("Member", back_populates="records")


class Budget(Base):
    __tablename__ = "budgets"
    __table_args__ = (UniqueConstraint('family_id', 'category', name='_family_category_uc'),)

    id = Column(Integer, primary_key=True, index=True)
    family_id = Column(Integer, ForeignKey("families.id"))
    category = Column(String)
    amount = Column(Float)
    created_at = Column(DateTime, default=datetime.utcnow)

    family = relationship("Family", back_populates="budgets")
