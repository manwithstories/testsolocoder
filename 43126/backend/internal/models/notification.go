package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeBookingConfirmation NotificationType = "booking_confirmation"
	NotificationTypeBookingReminder     NotificationType = "booking_reminder"
	NotificationTypeBookingCancellation NotificationType = "booking_cancellation"
	NotificationTypeBookingModification NotificationType = "booking_modification"
)

type NotificationStatus int

const (
	NotificationStatusPending NotificationStatus = 0
	NotificationStatusSent    NotificationStatus = 1
	NotificationStatusFailed  NotificationStatus = 2
)

type Notification struct {
	ID        uint              `json:"id" gorm:"primaryKey"`
	UserID    uint              `json:"user_id" gorm:"index;not null"`
	BookingID *uint             `json:"booking_id" gorm:"index"`
	Type      NotificationType  `json:"type" gorm:"size:50;not null"`
	Subject   string            `json:"subject" gorm:"size:200"`
	Content   string            `json:"content" gorm:"type:text"`
	Status    NotificationStatus `json:"status" gorm:"default:0"`
	RetryCount int              `json:"retry_count" gorm:"default:0"`
	SentAt    *time.Time        `json:"sent_at"`
	ErrorMsg  string            `json:"error_msg" gorm:"size:500"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"-" gorm:"index"`
}

func (Notification) TableName() string {
	return "notifications"
}
