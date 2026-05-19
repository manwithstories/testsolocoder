from typing import List, Dict, Optional, Tuple
from sqlalchemy.orm import Session
from sqlalchemy import func, extract
from datetime import datetime, timedelta, date
from dateutil.relativedelta import relativedelta
from .database import TrainingSession, TrainingSet, Exercise


def get_weekly_stats(db: Session, weeks: int = 4) -> List[Dict]:
    results = []
    today = date.today()

    for i in range(weeks - 1, -1, -1):
        week_start = today - timedelta(days=today.weekday() + i * 7)
        week_end = week_start + timedelta(days=6)

        sessions = db.query(TrainingSession).filter(
            TrainingSession.status == "completed",
            func.date(TrainingSession.start_time) >= week_start,
            func.date(TrainingSession.start_time) <= week_end
        ).all()

        total_volume = sum(s.total_volume for s in sessions)
        session_count = len(sessions)

        results.append({
            "week_start": week_start,
            "week_end": week_end,
            "session_count": session_count,
            "total_volume": round(total_volume, 2),
            "label": f"{week_start.month}/{week_start.day}-{week_end.month}/{week_end.day}"
        })

    return results


def get_monthly_stats(db: Session, months: int = 6) -> List[Dict]:
    results = []
    today = date.today()

    for i in range(months - 1, -1, -1):
        current_month = today - relativedelta(months=i)
        month_start = date(current_month.year, current_month.month, 1)
        if current_month.month == 12:
            month_end = date(current_month.year + 1, 1, 1) - timedelta(days=1)
        else:
            month_end = date(current_month.year, current_month.month + 1, 1) - timedelta(days=1)

        sessions = db.query(TrainingSession).filter(
            TrainingSession.status == "completed",
            func.date(TrainingSession.start_time) >= month_start,
            func.date(TrainingSession.start_time) <= month_end
        ).all()

        total_volume = sum(s.total_volume for s in sessions)
        session_count = len(sessions)

        results.append({
            "month": month_start,
            "session_count": session_count,
            "total_volume": round(total_volume, 2),
            "label": f"{month_start.year}年{month_start.month}月"
        })

    return results


def get_exercise_pr(db: Session, exercise_id: int) -> Optional[Dict]:
    exercise = db.query(Exercise).filter(Exercise.id == exercise_id).first()
    if not exercise:
        return None

    sets = db.query(TrainingSet).filter(
        TrainingSet.exercise_id == exercise_id,
        TrainingSession.status == "completed"
    ).join(TrainingSession).all()

    if not sets:
        return None

    max_weight = 0.0
    max_volume = 0.0
    max_weight_set = None
    max_volume_set = None

    for s in sets:
        volume = s.weight * s.reps if s.weight > 0 and s.reps > 0 else s.duration
        if s.weight > max_weight:
            max_weight = s.weight
            max_weight_set = s
        if volume > max_volume:
            max_volume = volume
            max_volume_set = s

    return {
        "exercise": exercise,
        "max_weight": round(max_weight, 2),
        "max_volume": round(max_volume, 2),
        "max_weight_set": max_weight_set,
        "max_volume_set": max_volume_set,
        "total_sets": len(sets),
    }


def get_all_prs(db: Session) -> List[Dict]:
    exercises = db.query(Exercise).filter(Exercise.is_archived == False).all()
    prs = []
    for ex in exercises:
        pr = get_exercise_pr(db, ex.id)
        if pr:
            prs.append(pr)
    return prs


def get_streak_days(db: Session) -> int:
    completed_sessions = db.query(TrainingSession).filter(
        TrainingSession.status == "completed"
    ).order_by(TrainingSession.start_time.desc()).all()

    if not completed_sessions:
        return 0

    today = date.today()
    training_dates = set()
    for s in completed_sessions:
        d = s.start_time.date()
        if d <= today:
            training_dates.add(d)

    streak = 0
    current_date = today
    while current_date in training_dates:
        streak += 1
        current_date -= timedelta(days=1)

    return streak


def get_total_stats(db: Session) -> Dict:
    total_sessions = db.query(TrainingSession).filter(TrainingSession.status == "completed").count()

    total_volume = db.query(func.sum(TrainingSession.total_volume)).filter(
        TrainingSession.status == "completed"
    ).scalar() or 0.0

    total_sets = db.query(TrainingSet).count()

    return {
        "total_sessions": total_sessions,
        "total_volume": round(total_volume, 2),
        "total_sets": total_sets,
        "streak_days": get_streak_days(db),
    }


def get_exercise_history(db: Session, exercise_id: int, limit: int = 20) -> List[Dict]:
    sets = db.query(TrainingSet).filter(
        TrainingSet.exercise_id == exercise_id,
        TrainingSession.status == "completed"
    ).join(TrainingSession).order_by(TrainingSet.created_at.desc()).limit(limit).all()

    history = []
    for s in sets:
        volume = s.weight * s.reps if s.weight > 0 and s.reps > 0 else s.duration
        history.append({
            "date": s.created_at,
            "set_number": s.set_number,
            "reps": s.reps,
            "weight": s.weight,
            "duration": s.duration,
            "volume": round(volume, 2),
            "session_id": s.session_id,
        })

    return history
