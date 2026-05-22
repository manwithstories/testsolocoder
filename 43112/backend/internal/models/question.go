package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"course_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	LessonID   *uuid.UUID     `gorm:"type:uuid;index" json:"lesson_id,omitempty"`
	Title      string         `gorm:"size:500;not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	ReplyCount int            `gorm:"default:0" json:"reply_count"`
	ViewCount  int            `gorm:"default:0" json:"view_count"`
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	IsResolved bool           `gorm:"default:false" json:"is_resolved"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User       *User          `json:"user,omitempty"`
	Course     *Course        `json:"-"`
	Answers    []Answer       `json:"answers,omitempty"`
}

type Answer struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	QuestionID uuid.UUID      `gorm:"type:uuid;index;not null" json:"question_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	IsBest     bool           `gorm:"default:false" json:"is_best"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User       *User          `json:"user,omitempty"`
	Question   *Question      `json:"-"`
}

type AnswerLike struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AnswerID   uuid.UUID `gorm:"type:uuid;index;not null" json:"answer_id"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type QuestionLike struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;index;not null" json:"question_id"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

func (a *Answer) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (a *AnswerLike) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (q *QuestionLike) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
