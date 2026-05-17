from typing import List, Optional, Dict
from sqlalchemy.orm import Session
from datetime import datetime
from .database import TrainingSession, TrainingSet, Exercise


def start_training(db: Session, notes: Optional[str] = None) -> TrainingSession:
    active_session = db.query(TrainingSession).filter(TrainingSession.status == "in_progress").first()
    if active_session:
        raise ValueError("已有进行中的训练，请先结束或取消当前训练")

    session = TrainingSession(status="in_progress", notes=notes)
    db.add(session)
    db.commit()
    db.refresh(session)
    return session


def get_active_session(db: Session) -> Optional[TrainingSession]:
    return db.query(TrainingSession).filter(TrainingSession.status == "in_progress").first()


def get_session(db: Session, session_id: int) -> Optional[TrainingSession]:
    return db.query(TrainingSession).filter(TrainingSession.id == session_id).first()


def add_set(
    db: Session,
    session_id: int,
    exercise_id: int,
    reps: int,
    weight: float = 0.0,
    duration: int = 0,
    notes: Optional[str] = None,
) -> TrainingSet:
    session = get_session(db, session_id)
    if not session:
        raise ValueError(f"训练会话 {session_id} 不存在")
    if session.status != "in_progress":
        raise ValueError("只能向进行中的训练添加组")

    exercise = db.query(Exercise).filter(Exercise.id == exercise_id).first()
    if not exercise:
        raise ValueError(f"动作 {exercise_id} 不存在")
    if exercise.is_archived:
        raise ValueError(f"动作 '{exercise.name}' 已归档，无法添加到训练")

    existing_sets = db.query(TrainingSet).filter(
        TrainingSet.session_id == session_id,
        TrainingSet.exercise_id == exercise_id
    ).count()
    set_number = existing_sets + 1

    training_set = TrainingSet(
        session_id=session_id,
        exercise_id=exercise_id,
        set_number=set_number,
        reps=reps,
        weight=weight,
        duration=duration,
        notes=notes,
    )
    db.add(training_set)
    db.commit()
    db.refresh(training_set)

    _update_session_volume(db, session)

    return training_set


def _update_session_volume(db: Session, session: TrainingSession) -> None:
    total_volume = 0.0
    for s in session.sets:
        if s.weight > 0 and s.reps > 0:
            total_volume += s.weight * s.reps
        elif s.duration > 0:
            total_volume += s.duration
    session.total_volume = total_volume
    db.commit()


def remove_set(db: Session, set_id: int) -> bool:
    training_set = db.query(TrainingSet).filter(TrainingSet.id == set_id).first()
    if not training_set:
        return False

    session = training_set.session
    if session.status != "in_progress":
        raise ValueError("只能修改进行中的训练")

    exercise_id = training_set.exercise_id
    session_id = training_set.session_id
    set_number = training_set.set_number

    db.delete(training_set)
    db.commit()

    remaining_sets = db.query(TrainingSet).filter(
        TrainingSet.session_id == session_id,
        TrainingSet.exercise_id == exercise_id,
        TrainingSet.set_number > set_number
    ).all()
    for s in remaining_sets:
        s.set_number -= 1
    db.commit()

    _update_session_volume(db, session)

    return True


def list_sets(db: Session, session_id: int) -> List[TrainingSet]:
    return db.query(TrainingSet).filter(TrainingSet.session_id == session_id).order_by(TrainingSet.id).all()


def get_session_exercises(db: Session, session_id: int) -> Dict[int, List[TrainingSet]]:
    sets = list_sets(db, session_id)
    exercises: Dict[int, List[TrainingSet]] = {}
    for s in sets:
        if s.exercise_id not in exercises:
            exercises[s.exercise_id] = []
        exercises[s.exercise_id].append(s)
    return exercises


def finish_training(db: Session, session_id: Optional[int] = None) -> Optional[TrainingSession]:
    if session_id:
        session = get_session(db, session_id)
    else:
        session = get_active_session(db)

    if not session:
        return None

    session.status = "completed"
    session.end_time = datetime.utcnow()
    db.commit()
    db.refresh(session)
    return session


def cancel_training(db: Session, session_id: Optional[int] = None) -> bool:
    if session_id:
        session = get_session(db, session_id)
    else:
        session = get_active_session(db)

    if not session:
        return False

    for s in session.sets:
        db.delete(s)
    db.delete(session)
    db.commit()
    return True


def list_sessions(db: Session, limit: int = 20, status: Optional[str] = None) -> List[TrainingSession]:
    query = db.query(TrainingSession)
    if status:
        query = query.filter(TrainingSession.status == status)
    return query.order_by(TrainingSession.start_time.desc()).limit(limit).all()
