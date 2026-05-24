package models

import (
	"time"

	"gorm.io/gorm"
)

type ReviewType string

const (
	ReviewTypeCourier   ReviewType = "courier"
	ReviewTypePublisher ReviewType = "publisher"
)

type Review struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderID      uint           `gorm:"uniqueIndex" json:"order_id"`
	Order        Order          `json:"order,omitempty"`
	ReviewerID   uint           `gorm:"index" json:"reviewer_id"`
	Reviewer     User           `json:"reviewer,omitempty"`
	RevieweeID   uint           `gorm:"index" json:"reviewee_id"`
	Reviewee     User           `json:"reviewee,omitempty"`
	ReviewType   ReviewType     `gorm:"size:20" json:"review_type"`
	Rating       int            `gorm:"default:5" json:"rating"`
	Content      string         `gorm:"size:1000" json:"content"`
	Tags         string         `gorm:"size:255" json:"tags"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type ChatMessage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   uint           `gorm:"index" json:"order_id"`
	SenderID  uint           `gorm:"index" json:"sender_id"`
	Sender    User           `json:"sender,omitempty"`
	Content   string         `gorm:"size:1000" json:"content"`
	MsgType   string         `gorm:"size:20;default:text" json:"msg_type"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	User      User           `json:"user,omitempty"`
	Type      string         `gorm:"size:20" json:"type"`
	Title     string         `gorm:"size:100" json:"title"`
	Content   string         `gorm:"size:500" json:"content"`
	RelatedID *uint          `json:"related_id,omitempty"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
