package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chapter struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"course_id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Position   int            `gorm:"default:0" json:"position"`
	IsFree     bool           `gorm:"default:false" json:"is_free"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Lessons    []Lesson       `json:"lessons,omitempty"`
}

func (c *Chapter) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
