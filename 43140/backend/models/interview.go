package models

import (
	"time"

	"gorm.io/gorm"
)

type InterviewStatus string

const (
	InterviewStatusScheduled InterviewStatus = "scheduled"
	InterviewStatusConfirmed InterviewStatus = "confirmed"
	InterviewStatusDeclined  InterviewStatus = "declined"
	InterviewStatusCompleted InterviewStatus = "completed"
	InterviewStatusCancelled InterviewStatus = "cancelled"
)

type Interview struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	ApplicationID   uint            `gorm:"not null;index" json:"application_id"`
	ScheduledAt     time.Time       `gorm:"not null" json:"scheduled_at"`
	Duration        int             `gorm:"default:60" json:"duration"`
	Location        string          `json:"location"`
	Interviewer     string          `json:"interviewer"`
	InterviewerEmail string         `json:"interviewer_email"`
	Notes           string          `gorm:"type:text" json:"notes"`
	Status          InterviewStatus `gorm:"type:varchar(20);default:scheduled" json:"status"`
	Application     Application     `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	Review          *Review         `gorm:"foreignKey:InterviewID" json:"review,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
}
