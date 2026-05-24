package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAccepted   OrderStatus = "accepted"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TaskID        uint           `gorm:"uniqueIndex" json:"task_id"`
	Task          Task           `json:"task,omitempty"`
	PublisherID   uint           `gorm:"index" json:"publisher_id"`
	Publisher     User           `json:"publisher,omitempty"`
	CourierID     uint           `gorm:"index" json:"courier_id"`
	Courier       User           `json:"courier,omitempty"`
	Status        OrderStatus    `gorm:"size:20;default:pending" json:"status"`
	Reward        float64        `json:"reward"`
	ServiceFee    float64        `json:"service_fee"`
	ActualPayment float64        `json:"actual_payment"`
	StartTime     *time.Time     `json:"start_time,omitempty"`
	EndTime       *time.Time     `json:"end_time,omitempty"`
	Tracks        []OrderTrack   `gorm:"foreignKey:OrderID" json:"tracks,omitempty"`
	ProofImages   []OrderProof   `gorm:"foreignKey:OrderID" json:"proof_images,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderTrack struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index" json:"order_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Address   string    `gorm:"size:255" json:"address"`
	Message   string    `gorm:"size:500" json:"message"`
	EventType string    `gorm:"size:20" json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderProof struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index" json:"order_id"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	ProofType string    `gorm:"size:20" json:"proof_type"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderStatusLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index" json:"order_id"`
	OldStatus string    `gorm:"size:20" json:"old_status"`
	NewStatus string    `gorm:"size:20" json:"new_status"`
	Reason    string    `gorm:"size:500" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
