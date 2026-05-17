import csv
from typing import Optional
from pathlib import Path
from sqlalchemy.orm import Session
from .database import TrainingSession, TrainingSet, Exercise, PlanTemplate


def export_training_history(db: Session, output_path: str) -> bool:
    try:
        Path(output_path).parent.mkdir(parents=True, exist_ok=True)

        with open(output_path, "w", newline="", encoding="utf-8-sig") as f:
            writer = csv.writer(f)
            writer.writerow([
                "训练ID", "开始时间", "结束时间", "状态", "总容量", "备注",
                "动作名称", "组号", "次数", "重量(kg)", "时长(秒)", "组备注"
            ])

            sessions = db.query(TrainingSession).order_by(TrainingSession.start_time).all()

            for session in sessions:
                for s in session.sets:
                    exercise = db.query(Exercise).filter(Exercise.id == s.exercise_id).first()
                    exercise_name = exercise.name if exercise else "未知动作"

                    writer.writerow([
                        session.id,
                        session.start_time.strftime("%Y-%m-%d %H:%M:%S") if session.start_time else "",
                        session.end_time.strftime("%Y-%m-%d %H:%M:%S") if session.end_time else "",
                        session.status,
                        session.total_volume,
                        session.notes or "",
                        exercise_name,
                        s.set_number,
                        s.reps,
                        s.weight,
                        s.duration,
                        s.notes or "",
                    ])

        return True
    except Exception:
        return False


def export_exercises(db: Session, output_path: str) -> bool:
    try:
        Path(output_path).parent.mkdir(parents=True, exist_ok=True)

        with open(output_path, "w", newline="", encoding="utf-8-sig") as f:
            writer = csv.writer(f)
            writer.writerow([
                "动作ID", "名称", "分类", "默认组数", "默认次数", "默认时长(秒)", "是否归档", "创建时间"
            ])

            exercises = db.query(Exercise).order_by(Exercise.category, Exercise.name).all()

            for ex in exercises:
                writer.writerow([
                    ex.id,
                    ex.name,
                    ex.category,
                    ex.default_sets,
                    ex.default_reps,
                    ex.default_duration,
                    "是" if ex.is_archived else "否",
                    ex.created_at.strftime("%Y-%m-%d %H:%M:%S") if ex.created_at else "",
                ])

        return True
    except Exception:
        return False


def export_plans(db: Session, output_path: str) -> bool:
    try:
        Path(output_path).parent.mkdir(parents=True, exist_ok=True)

        with open(output_path, "w", newline="", encoding="utf-8-sig") as f:
            writer = csv.writer(f)
            writer.writerow([
                "计划ID", "计划名称", "描述", "动作顺序", "动作名称", "组数", "次数", "时长(秒)", "备注", "动作是否归档"
            ])

            plans = db.query(PlanTemplate).order_by(PlanTemplate.name).all()

            for plan in plans:
                for pe in sorted(plan.exercises, key=lambda x: x.order_index):
                    exercise = db.query(Exercise).filter(Exercise.id == pe.exercise_id).first()
                    exercise_name = exercise.name if exercise else "未知动作"
                    is_archived = exercise.is_archived if exercise else True

                    writer.writerow([
                        plan.id,
                        plan.name,
                        plan.description or "",
                        pe.order_index + 1,
                        exercise_name,
                        pe.sets,
                        pe.reps,
                        pe.duration,
                        pe.notes or "",
                        "是" if is_archived else "否",
                    ])

        return True
    except Exception:
        return False


def export_all(db: Session, output_dir: str) -> dict:
    dir_path = Path(output_dir)
    dir_path.mkdir(parents=True, exist_ok=True)

    results = {
        "training_history": export_training_history(db, str(dir_path / "training_history.csv")),
        "exercises": export_exercises(db, str(dir_path / "exercises.csv")),
        "plans": export_plans(db, str(dir_path / "plans.csv")),
    }

    return results
