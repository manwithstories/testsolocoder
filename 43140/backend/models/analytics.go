package models

import (
	"time"

	"gorm.io/gorm"
)

type JobView struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	JobID     uint           `gorm:"not null;index" json:"job_id"`
	UserID    *uint          `json:"user_id"`
	IP        string         `json:"ip"`
	ViewedAt  time.Time      `json:"viewed_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type NotificationType string

const (
	NotificationTypeApplication NotificationType = "application"
	NotificationTypeInterview   NotificationType = "interview"
	NotificationTypeReview      NotificationType = "review"
	NotificationTypeSystem      NotificationType = "system"
)

type Notification struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	UserID    uint              `gorm:"not null;index" json:"user_id"`
	Type      NotificationType  `gorm:"type:varchar(20)" json:"type"`
	Title     string            `json:"title"`
	Message   string            `gorm:"type:text" json:"message"`
	Read      bool              `gorm:"default:false" json:"read"`
	RelatedID *uint             `json:"related_id"`
	CreatedAt time.Time         `json:"created_at"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`
}
