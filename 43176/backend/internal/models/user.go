package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RolePublisher UserRole = "publisher"
	RoleCourier   UserRole = "courier"
	RoleAdmin     UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusFrozen   UserStatus = "frozen"
	UserStatusVerified UserStatus = "verified"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Phone        string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Password     string         `json:"-"`
	Nickname     string         `gorm:"size:50" json:"nickname"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Role         UserRole       `gorm:"size:20;default:publisher" json:"role"`
	Status       UserStatus     `gorm:"size:20;default:active" json:"status"`
	Balance      float64        `gorm:"default:0" json:"balance"`
	Rating       float64        `gorm:"default:5" json:"rating"`
	OrderCount   int            `gorm:"default:0" json:"order_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Verification struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	User       User           `json:"user,omitempty"`
	IDCardNo   string         `gorm:"size:20" json:"id_card_no"`
	IDCardName string         `gorm:"size:50" json:"id_card_name"`
	IDCardFront string        `gorm:"size:255" json:"id_card_front"`
	IDCardBack  string        `gorm:"size:255" json:"id_card_back"`
	Status     string         `gorm:"size:20;default:pending" json:"status"`
	Reason     string         `gorm:"size:500" json:"reason"`
	ReviewedAt *time.Time     `json:"reviewed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type CourierProfile struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"uniqueIndex" json:"user_id"`
	User             User           `json:"user,omitempty"`
	Status           string         `gorm:"size:20;default:pending" json:"status"`
	Level            int            `gorm:"default:1" json:"level"`
	TotalOrders      int            `gorm:"default:0" json:"total_orders"`
	CompletedOrders  int            `gorm:"default:0" json:"completed_orders"`
	CancelledOrders  int            `gorm:"default:0" json:"cancelled_orders"`
	Rating           float64        `gorm:"default:5" json:"rating"`
	CurrentTaskID    *uint          `json:"current_task_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
