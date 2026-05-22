package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quiz struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	LessonID   uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"lesson_id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	PassScore  int            `gorm:"default:60" json:"pass_score"`
	TimeLimit  int            `gorm:"default:0" json:"time_limit"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Questions  []QuizQuestion `json:"questions,omitempty"`
}

type QuizQuestion struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	QuizID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"quiz_id"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Type        string         `gorm:"size:20;default:single" json:"type"`
	Score       int            `gorm:"default:10" json:"score"`
	Position    int            `gorm:"default:0" json:"position"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Options     []QuizOption   `json:"options,omitempty"`
}

type QuizOption struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	QuestionID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"question_id"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	IsCorrect   bool           `gorm:"default:false" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserQuiz struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	QuizID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"quiz_id"`
	Score      int            `gorm:"default:0" json:"score"`
	IsPassed   bool           `gorm:"default:false" json:"is_passed"`
	CompletedAt *time.Time    `json:"completed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

func (q *QuizQuestion) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

func (q *QuizOption) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

func (q *UserQuiz) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
