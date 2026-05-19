from typing import List, Optional
from sqlalchemy.orm import Session
from .database import Exercise, TrainingSet, PlanExercise

CATEGORIES = ["力量", "有氧", "柔韧性", "其他"]


def add_exercise(
    db: Session,
    name: str,
    category: str,
    default_sets: int = 3,
    default_reps: int = 10,
    default_duration: int = 0,
) -> Exercise:
    if category not in CATEGORIES:
        raise ValueError(f"分类必须是以下之一: {', '.join(CATEGORIES)}")

    existing = db.query(Exercise).filter(Exercise.name == name).first()
    if existing:
        if existing.is_archived:
            existing.is_archived = False
            existing.category = category
            existing.default_sets = default_sets
            existing.default_reps = default_reps
            existing.default_duration = default_duration
            db.commit()
            db.refresh(existing)
            return existing
        raise ValueError(f"动作 '{name}' 已存在")

    exercise = Exercise(
        name=name,
        category=category,
        default_sets=default_sets,
        default_reps=default_reps,
        default_duration=default_duration,
    )
    db.add(exercise)
    db.commit()
    db.refresh(exercise)
    return exercise


def get_exercise(db: Session, exercise_id: int) -> Optional[Exercise]:
    return db.query(Exercise).filter(Exercise.id == exercise_id).first()


def get_exercise_by_name(db: Session, name: str) -> Optional[Exercise]:
    return db.query(Exercise).filter(Exercise.name == name).first()


def list_exercises(db: Session, category: Optional[str] = None, include_archived: bool = False) -> List[Exercise]:
    query = db.query(Exercise)
    if category:
        query = query.filter(Exercise.category == category)
    if not include_archived:
        query = query.filter(Exercise.is_archived == False)
    return query.order_by(Exercise.category, Exercise.name).all()


def update_exercise(
    db: Session,
    exercise_id: int,
    name: Optional[str] = None,
    category: Optional[str] = None,
    default_sets: Optional[int] = None,
    default_reps: Optional[int] = None,
    default_duration: Optional[int] = None,
) -> Optional[Exercise]:
    exercise = get_exercise(db, exercise_id)
    if not exercise:
        return None

    if category and category not in CATEGORIES:
        raise ValueError(f"分类必须是以下之一: {', '.join(CATEGORIES)}")

    if name:
        existing = db.query(Exercise).filter(Exercise.name == name, Exercise.id != exercise_id).first()
        if existing:
            raise ValueError(f"动作名称 '{name}' 已被使用")
        exercise.name = name

    if category:
        exercise.category = category
    if default_sets is not None:
        exercise.default_sets = default_sets
    if default_reps is not None:
        exercise.default_reps = default_reps
    if default_duration is not None:
        exercise.default_duration = default_duration

    db.commit()
    db.refresh(exercise)
    return exercise


def delete_exercise(db: Session, exercise_id: int) -> bool:
    exercise = get_exercise(db, exercise_id)
    if not exercise:
        return False

    has_training_records = db.query(TrainingSet).filter(TrainingSet.exercise_id == exercise_id).first() is not None
    has_plan_references = db.query(PlanExercise).filter(PlanExercise.exercise_id == exercise_id).first() is not None

    if has_training_records or has_plan_references:
        exercise.is_archived = True
        db.commit()
        return True
    else:
        db.delete(exercise)
        db.commit()
        return True


def get_categories() -> List[str]:
    return CATEGORIES
