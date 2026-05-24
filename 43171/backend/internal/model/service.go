package model

import (
	"time"

	"gorm.io/gorm"
)

type ServiceStatus string

const (
	ServiceStatusOpen     ServiceStatus = "open"
	ServiceStatusAssigned ServiceStatus = "assigned"
	ServiceStatusProgress ServiceStatus = "progress"
	ServiceStatusCompleted ServiceStatus = "completed"
	ServiceStatusCancelled ServiceStatus = "cancelled"
)

type BidStatus string

const (
	BidStatusPending  BidStatus = "pending"
	BidStatusAccepted BidStatus = "accepted"
	BidStatusRejected BidStatus = "rejected"
)

type AerialService struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ServiceNo   string         `gorm:"size:32;uniqueIndex;not null" json:"service_no"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PilotID     *uint          `gorm:"index" json:"pilot_id"`
	Pilot       *User          `gorm:"foreignKey:PilotID" json:"pilot,omitempty"`
	Title       string         `gorm:"size:128;not null" json:"title"`
	Description string         `gorm:"size:1000" json:"description"`
	Region      string         `gorm:"size:64;index" json:"region"`
	Address     string         `gorm:"size:256" json:"address"`
	ServiceDate *time.Time     `json:"service_date"`
	ServiceTime string         `gorm:"size:32" json:"service_time"`
	Duration    int            `json:"duration"`
	BudgetMin   float64        `json:"budget_min"`
	BudgetMax   float64        `json:"budget_max"`
	FinalPrice  float64        `json:"final_price"`
	Status      ServiceStatus  `gorm:"size:16;default:open;index" json:"status"`
	Images      string         `gorm:"size:1024" json:"images"`
	Remark      string         `gorm:"size:500" json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ServiceBid struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ServiceID   uint           `gorm:"index;not null" json:"service_id"`
	Service     *AerialService `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	PilotID     uint           `gorm:"index;not null" json:"pilot_id"`
	Pilot       *User          `gorm:"foreignKey:PilotID" json:"pilot,omitempty"`
	Price       float64        `json:"price"`
	Message     string         `gorm:"size:500" json:"message"`
	Status      BidStatus      `gorm:"size:16;default:pending" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
