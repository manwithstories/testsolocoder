package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusPicked    OrderStatus = "picked"
	OrderStatusReturned  OrderStatus = "returned"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type RentalOrder struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderNo      string         `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DroneID      uint           `gorm:"index;not null" json:"drone_id"`
	Drone        *Drone         `gorm:"foreignKey:DroneID" json:"drone,omitempty"`
	StartDate    time.Time      `json:"start_date"`
	EndDate      time.Time      `json:"end_date"`
	ReturnDate   *time.Time     `json:"return_date"`
	Region       string         `gorm:"size:64" json:"region"`
	Address      string         `gorm:"size:256" json:"address"`
	ContactName  string         `gorm:"size:64" json:"contact_name"`
	ContactPhone string         `gorm:"size:32" json:"contact_phone"`
	TotalDays    int            `json:"total_days"`
	PricePerDay  float64        `json:"price_per_day"`
	RentalFee    float64        `json:"rental_fee"`
	Deposit      float64        `json:"deposit"`
	InsuranceFee float64        `json:"insurance_fee"`
	LateFee      float64        `gorm:"default:0" json:"late_fee"`
	TotalAmount  float64        `json:"total_amount"`
	PaidAmount   float64        `gorm:"default:0" json:"paid_amount"`
	RefundAmount float64        `gorm:"default:0" json:"refund_amount"`
	Status       OrderStatus    `gorm:"size:16;default:pending;index" json:"status"`
	Remark       string         `gorm:"size:500" json:"remark"`
	CancelReason string         `gorm:"size:256" json:"cancel_reason"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusRefund  PaymentStatus = "refund"
)

type Payment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PaymentNo   string         `gorm:"size:32;uniqueIndex;not null" json:"payment_no"`
	OrderID     uint           `gorm:"index;not null" json:"order_id"`
	Order       *RentalOrder   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	Amount      float64        `json:"amount"`
	PayType     string         `gorm:"size:32" json:"pay_type"`
	Status      PaymentStatus  `gorm:"size:16;default:pending;index" json:"status"`
	TradeNo     string         `gorm:"size:64" json:"trade_no"`
	PaidAt      *time.Time     `json:"paid_at"`
	RefundAt    *time.Time     `json:"refund_at"`
	RefundReason string        `gorm:"size:256" json:"refund_reason"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
