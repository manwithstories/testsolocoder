package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reminder struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"plan_id"`
	ActivityID    uuid.UUID      `gorm:"type:uuid;index" json:"activity_id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title         string         `gorm:"type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Description   string         `gorm:"type:text" json:"description"`
	ReminderTime  time.Time      `gorm:"not null" json:"reminder_time" validate:"required"`
	IsSent        bool           `gorm:"default:false" json:"is_sent"`
	SentAt        *time.Time     `json:"sent_at"`
	Channel       string         `gorm:"type:varchar(20);default:'email'" json:"channel"`
	RepeatPattern string         `gorm:"type:varchar(50)" json:"repeat_pattern"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Plan     *TravelPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Activity *Activity   `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
	User     *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
