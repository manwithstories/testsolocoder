package model

import (
	"time"
)

type OrderType string

const (
	OrderTypeVenue  OrderType = "venue"
	OrderTypeDevice OrderType = "device"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	OrderNo      string      `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID       uint        `gorm:"index;not null" json:"user_id"`
	Type         OrderType   `gorm:"size:20;not null" json:"type"`
	ItemID       uint        `gorm:"index;not null" json:"item_id"`
	ItemName     string      `gorm:"size:100" json:"item_name"`
	StartTime    time.Time   `gorm:"not null" json:"start_time"`
	EndTime      time.Time   `gorm:"not null" json:"end_time"`
	TotalHours   float64     `gorm:"default:0" json:"total_hours"`
	Quantity     int         `gorm:"default:1" json:"quantity"`
	TotalAmount  float64     `gorm:"default:0" json:"total_amount"`
	DepositAmount float64    `gorm:"default:0" json:"deposit_amount"`
	Status       OrderStatus `gorm:"size:20;default:'pending'" json:"status"`
	Purpose      string      `gorm:"type:text" json:"purpose"`
	ContactName  string      `gorm:"size:50" json:"contact_name"`
	ContactPhone string      `gorm:"size:20" json:"contact_phone"`
	CancelReason string      `gorm:"type:text" json:"cancel_reason"`
	CancelledAt  *time.Time  `json:"cancelled_at"`
	ReviewedBy   *uint       `json:"reviewed_by"`
	ReviewNote   string      `gorm:"type:text" json:"review_note"`
	ReviewedAt   *time.Time  `json:"reviewed_at"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	User         *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
