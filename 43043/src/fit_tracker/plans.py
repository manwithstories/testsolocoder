from typing import List, Optional, Tuple
from sqlalchemy.orm import Session
from .database import PlanTemplate, PlanExercise, ScheduledPlan, Exercise


WEEKDAYS = ["周一", "周二", "周三", "周四", "周五", "周六", "周日"]


def create_plan(db: Session, name: str, description: Optional[str] = None) -> PlanTemplate:
    existing = db.query(PlanTemplate).filter(PlanTemplate.name == name).first()
    if existing:
        raise ValueError(f"计划模板 '{name}' 已存在")

    plan = PlanTemplate(name=name, description=description)
    db.add(plan)
    db.commit()
    db.refresh(plan)
    return plan


def get_plan(db: Session, plan_id: int) -> Optional[PlanTemplate]:
    return db.query(PlanTemplate).filter(PlanTemplate.id == plan_id).first()


def get_plan_by_name(db: Session, name: str) -> Optional[PlanTemplate]:
    return db.query(PlanTemplate).filter(PlanTemplate.name == name).first()


def list_plans(db: Session) -> List[PlanTemplate]:
    return db.query(PlanTemplate).order_by(PlanTemplate.name).all()


def add_exercise_to_plan(
    db: Session,
    plan_id: int,
    exercise_id: int,
    sets: int = 3,
    reps: int = 10,
    duration: int = 0,
    notes: Optional[str] = None,
) -> PlanExercise:
    plan = get_plan(db, plan_id)
    if not plan:
        raise ValueError(f"计划模板 {plan_id} 不存在")

    exercise = db.query(Exercise).filter(Exercise.id == exercise_id).first()
    if not exercise:
        raise ValueError(f"动作 {exercise_id} 不存在")

    existing = db.query(PlanExercise).filter(
        PlanExercise.plan_id == plan_id,
        PlanExercise.exercise_id == exercise_id
    ).first()
    if existing:
        raise ValueError(f"动作 '{exercise.name}' 已在计划中")

    max_order = db.query(PlanExercise).filter(PlanExercise.plan_id == plan_id).count()

    plan_exercise = PlanExercise(
        plan_id=plan_id,
        exercise_id=exercise_id,
        order_index=max_order,
        sets=sets,
        reps=reps,
        duration=duration,
        notes=notes,
    )
    db.add(plan_exercise)
    db.commit()
    db.refresh(plan_exercise)
    return plan_exercise


def remove_exercise_from_plan(db: Session, plan_id: int, exercise_id: int) -> bool:
    plan_exercise = db.query(PlanExercise).filter(
        PlanExercise.plan_id == plan_id,
        PlanExercise.exercise_id == exercise_id
    ).first()
    if not plan_exercise:
        return False

    order_index = plan_exercise.order_index
    db.delete(plan_exercise)
    db.commit()

    remaining = db.query(PlanExercise).filter(
        PlanExercise.plan_id == plan_id,
        PlanExercise.order_index > order_index
    ).all()
    for pe in remaining:
        pe.order_index -= 1
    db.commit()

    return True


def update_plan_exercise(
    db: Session,
    plan_id: int,
    exercise_id: int,
    sets: Optional[int] = None,
    reps: Optional[int] = None,
    duration: Optional[int] = None,
    notes: Optional[str] = None,
) -> Optional[PlanExercise]:
    plan_exercise = db.query(PlanExercise).filter(
        PlanExercise.plan_id == plan_id,
        PlanExercise.exercise_id == exercise_id
    ).first()
    if not plan_exercise:
        return None

    if sets is not None:
        plan_exercise.sets = sets
    if reps is not None:
        plan_exercise.reps = reps
    if duration is not None:
        plan_exercise.duration = duration
    if notes is not None:
        plan_exercise.notes = notes

    db.commit()
    db.refresh(plan_exercise)
    return plan_exercise


def delete_plan(db: Session, plan_id: int) -> bool:
    plan = get_plan(db, plan_id)
    if not plan:
        return False

    db.delete(plan)
    db.commit()
    return True


def get_plan_exercises_with_status(db: Session, plan_id: int) -> List[Tuple[PlanExercise, bool]]:
    plan_exercises = db.query(PlanExercise).filter(
        PlanExercise.plan_id == plan_id
    ).order_by(PlanExercise.order_index).all()

    result = []
    for pe in plan_exercises:
        exercise = db.query(Exercise).filter(Exercise.id == pe.exercise_id).first()
        is_archived = exercise.is_archived if exercise else True
        result.append((pe, is_archived))
    return result


def schedule_plan(db: Session, plan_id: int, day_of_week: int) -> ScheduledPlan:
    if day_of_week < 0 or day_of_week > 6:
        raise ValueError("day_of_week 必须在 0-6 之间 (0=周一, 6=周日)")

    plan = get_plan(db, plan_id)
    if not plan:
        raise ValueError(f"计划模板 {plan_id} 不存在")

    existing = db.query(ScheduledPlan).filter(
        ScheduledPlan.plan_id == plan_id,
        ScheduledPlan.day_of_week == day_of_week
    ).first()
    if existing:
        if not existing.is_active:
            existing.is_active = True
            db.commit()
            db.refresh(existing)
            return existing
        raise ValueError(f"该计划已安排在{WEEKDAYS[day_of_week]}")

    scheduled = ScheduledPlan(plan_id=plan_id, day_of_week=day_of_week)
    db.add(scheduled)
    db.commit()
    db.refresh(scheduled)
    return scheduled


def unschedule_plan(db: Session, scheduled_id: int) -> bool:
    scheduled = db.query(ScheduledPlan).filter(ScheduledPlan.id == scheduled_id).first()
    if not scheduled:
        return False

    scheduled.is_active = False
    db.commit()
    return True


def get_scheduled_plans(db: Session, day_of_week: Optional[int] = None, active_only: bool = True) -> List[ScheduledPlan]:
    query = db.query(ScheduledPlan)
    if day_of_week is not None:
        query = query.filter(ScheduledPlan.day_of_week == day_of_week)
    if active_only:
        query = query.filter(ScheduledPlan.is_active == True)
    return query.order_by(ScheduledPlan.day_of_week).all()


def get_weekday_name(day_of_week: int) -> str:
    return WEEKDAYS[day_of_week] if 0 <= day_of_week <= 6 else "未知"
