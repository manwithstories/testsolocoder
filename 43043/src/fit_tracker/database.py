import os
from pathlib import Path
from sqlalchemy import create_engine, Column, Integer, String, Float, DateTime, Boolean, ForeignKey, Text
from sqlalchemy.orm import declarative_base, sessionmaker, relationship
from datetime import datetime

DB_PATH = os.path.expanduser("~/.fit_tracker/fit_tracker.db")
Path(os.path.dirname(DB_PATH)).mkdir(parents=True, exist_ok=True)

engine = create_engine(f"sqlite:///{DB_PATH}")
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()


class Exercise(Base):
    __tablename__ = "exercises"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(100), nullable=False, unique=True)
    category = Column(String(50), nullable=False)
    default_sets = Column(Integer, default=3)
    default_reps = Column(Integer, default=10)
    default_duration = Column(Integer, default=0)
    is_archived = Column(Boolean, default=False)
    created_at = Column(DateTime, default=datetime.utcnow)

    training_sets = relationship("TrainingSet", back_populates="exercise")
    plan_exercises = relationship("PlanExercise", back_populates="exercise")


class TrainingSession(Base):
    __tablename__ = "training_sessions"

    id = Column(Integer, primary_key=True, index=True)
    start_time = Column(DateTime, default=datetime.utcnow)
    end_time = Column(DateTime, nullable=True)
    status = Column(String(20), default="in_progress")
    total_volume = Column(Float, default=0.0)
    notes = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    sets = relationship("TrainingSet", back_populates="session", cascade="all, delete-orphan")


class TrainingSet(Base):
    __tablename__ = "training_sets"

    id = Column(Integer, primary_key=True, index=True)
    session_id = Column(Integer, ForeignKey("training_sessions.id"), nullable=False)
    exercise_id = Column(Integer, ForeignKey("exercises.id"), nullable=False)
    set_number = Column(Integer, nullable=False)
    reps = Column(Integer, default=0)
    weight = Column(Float, default=0.0)
    duration = Column(Integer, default=0)
    notes = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    session = relationship("TrainingSession", back_populates="sets")
    exercise = relationship("Exercise", back_populates="training_sets")


class PlanTemplate(Base):
    __tablename__ = "plan_templates"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(100), nullable=False)
    description = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    exercises = relationship("PlanExercise", back_populates="plan", cascade="all, delete-orphan")
    schedules = relationship("ScheduledPlan", back_populates="plan", cascade="all, delete-orphan")


class PlanExercise(Base):
    __tablename__ = "plan_exercises"

    id = Column(Integer, primary_key=True, index=True)
    plan_id = Column(Integer, ForeignKey("plan_templates.id"), nullable=False)
    exercise_id = Column(Integer, ForeignKey("exercises.id"), nullable=False)
    order_index = Column(Integer, default=0)
    sets = Column(Integer, default=3)
    reps = Column(Integer, default=10)
    duration = Column(Integer, default=0)
    notes = Column(Text, nullable=True)

    plan = relationship("PlanTemplate", back_populates="exercises")
    exercise = relationship("Exercise", back_populates="plan_exercises")


class ScheduledPlan(Base):
    __tablename__ = "scheduled_plans"

    id = Column(Integer, primary_key=True, index=True)
    plan_id = Column(Integer, ForeignKey("plan_templates.id"), nullable=False)
    day_of_week = Column(Integer, nullable=False)
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    plan = relationship("PlanTemplate", back_populates="schedules")


def init_db():
    Base.metadata.create_all(bind=engine)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
