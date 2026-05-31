package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityStatus string

const (
	ActivityStatusRecruiting ActivityStatus = "recruiting"
	ActivityStatusOngoing    ActivityStatus = "ongoing"
	ActivityStatusCompleted  ActivityStatus = "completed"
	ActivityStatusCancelled  ActivityStatus = "cancelled"
)

type TastingActivity struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"index;not null" json:"user_id"`
	Title             string         `gorm:"size:255;not null;index" json:"title"`
	Description       string         `gorm:"type:text" json:"description"`
	Location          string         `gorm:"size:512" json:"location"`
	ActivityTime      *time.Time     `json:"activity_time"`
	MaxParticipants   int            `gorm:"not null;default:0" json:"max_participants"`
	CurrentParticipants int         `gorm:"not null;default:0" json:"current_participants"`
	Status            ActivityStatus `gorm:"size:32;not null;default:recruiting;index" json:"status"`
	CoverImage        string         `gorm:"size:512" json:"cover_image"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	Organizer    *User                 `gorm:"foreignKey:UserID" json:"organizer,omitempty"`
	Participants []ActivityParticipant `gorm:"foreignKey:ActivityID" json:"participants,omitempty"`
}

func (TastingActivity) TableName() string {
	return "tasting_activities"
}
