package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"course_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Rating     int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Content    string         `gorm:"type:text" json:"content"`
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	IsAnonymous bool          `gorm:"default:false" json:"is_anonymous"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User       *User          `json:"user,omitempty"`
	Course     *Course        `json:"-"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
