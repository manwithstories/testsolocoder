package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonType string

const (
	LessonTypeVideo    LessonType = "video"
	LessonTypeDocument LessonType = "document"
	LessonTypeQuiz     LessonType = "quiz"
)

type Lesson struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ChapterID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"chapter_id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Type        LessonType     `gorm:"size:20;default:video" json:"type"`
	Content     string         `gorm:"type:text" json:"content"`
	VideoURL    string         `gorm:"size:500" json:"video_url"`
	VideoLength int            `gorm:"default:0" json:"video_length"`
	DocURL      string         `gorm:"size:500" json:"doc_url"`
	DocName     string         `gorm:"size:200" json:"doc_name"`
	Position    int            `gorm:"default:0" json:"position"`
	IsFree      bool           `gorm:"default:false" json:"is_free"`
	IsPublished bool           `gorm:"default:true" json:"is_published"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Quiz        *Quiz          `json:"quiz,omitempty"`
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
