package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeSystem        MessageType = "system"
	MessageTypeOrder         MessageType = "order"
	MessageTypeInvitation    MessageType = "invitation"
	MessageTypeReview        MessageType = "review"
	MessageTypeComplaint     MessageType = "complaint"
	MessageTypeWithdraw      MessageType = "withdraw"
)

type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Type      MessageType    `json:"type" gorm:"size:30;not null"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	RelatedID *uint          `json:"related_id"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type OperationLog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OperatorID  uint           `json:"operator_id" gorm:"index"`
	OperatorRole string        `json:"operator_role" gorm:"size:20"`
	Module      string         `json:"module" gorm:"size:50"`
	Action      string         `json:"action" gorm:"size:50"`
	TargetID    *uint          `json:"target_id"`
	TargetType  string         `json:"target_type" gorm:"size:50"`
	Content     string         `json:"content" gorm:"type:text"`
	IP          string         `json:"ip" gorm:"size:50"`
	UserAgent   string         `json:"user_agent" gorm:"size:255"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
