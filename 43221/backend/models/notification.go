package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationAppointmentSuccess NotificationType = "appointment_success"
	NotificationAppointmentCancel  NotificationType = "appointment_cancel"
	NotificationAppointmentRemind  NotificationType = "appointment_remind"
	NotificationPaymentSuccess     NotificationType = "payment_success"
	NotificationPaymentRefund      NotificationType = "payment_refund"
	NotificationReviewReply        NotificationType = "review_reply"
	NotificationSystem             NotificationType = "system"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID        `json:"user_id" gorm:"type:uuid;not null;index"`
	Type      NotificationType `json:"type" gorm:"type:varchar(30);not null"`
	Title     string           `json:"title" gorm:"not null;size:200"`
	Content   string           `json:"content" gorm:"type:text"`
	Data      string           `json:"data" gorm:"type:text"`
	IsRead    bool             `json:"is_read" gorm:"default:false"`
	ReadAt    *time.Time       `json:"read_at"`
	CreatedAt time.Time        `json:"created_at"`
	DeletedAt gorm.DeletedAt   `json:"-" gorm:"index"`

	User      User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type NotificationTemplate struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Type      string         `json:"type" gorm:"uniqueIndex;not null;size:50"`
	Title     string         `json:"title" gorm:"not null;size:200"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Variables string         `json:"variables" gorm:"type:text"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

func (n *NotificationTemplate) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
