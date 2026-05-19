package model

import (
	"time"
)

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

type Review struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	OrderID    uint         `gorm:"index;not null" json:"order_id"`
	UserID     uint         `gorm:"index;not null" json:"user_id"`
	Rating     int          `gorm:"not null" json:"rating"`
	Content    string       `gorm:"type:text" json:"content"`
	Images     string       `gorm:"type:json" json:"-"`
	Status     ReviewStatus `gorm:"size:20;default:'pending'" json:"status"`
	ReviewNote string       `gorm:"type:text" json:"review_note"`
	ReviewedBy *uint        `json:"reviewed_by"`
	ReviewedAt *time.Time   `json:"reviewed_at"`
	CreatedAt  time.Time    `json:"created_at"`
	Order      *Order       `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User       *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
