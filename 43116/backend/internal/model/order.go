package model

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusRefunded  OrderStatus = "refunded"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Order struct {
	ID            uint          `gorm:"primarykey" json:"id"`
	OrderNo       string        `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	UserID        uint          `gorm:"index;not null" json:"user_id"`
	User          *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CarID         uint          `gorm:"index;not null" json:"car_id"`
	Car           *Car          `gorm:"foreignKey:CarID" json:"car,omitempty"`
	BookingID     uint          `gorm:"index;not null" json:"booking_id"`
	Booking       *Booking      `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	TotalAmount   float64       `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Discount      float64       `gorm:"type:decimal(10,2);default:0" json:"discount"`
	FinalAmount   float64       `gorm:"type:decimal(10,2);not null" json:"final_amount"`
	PaidAmount    float64       `gorm:"type:decimal(10,2);default:0" json:"paid_amount"`
	PaymentStatus PaymentStatus `gorm:"size:20;default:pending" json:"payment_status"`
	PaymentMethod string        `gorm:"size:20" json:"payment_method"`
	PaidAt        *time.Time    `json:"paid_at"`
	Status        OrderStatus   `gorm:"size:20;default:pending" json:"status"`
	RefundReason  string        `gorm:"size:500" json:"refund_reason"`
	RefundedAt    *time.Time    `json:"refunded_at"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     *time.Time    `gorm:"index" json:"-"`
}

func (Order) TableName() string {
	return "orders"
}

type Review struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CarID       uint       `gorm:"index;not null" json:"car_id"`
	Car         *Car       `gorm:"foreignKey:CarID" json:"car,omitempty"`
	BookingID   uint       `gorm:"index;not null" json:"booking_id"`
	Booking     *Booking   `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Rating      int        `gorm:"not null" json:"rating"`
	Content     string     `gorm:"type:text" json:"content"`
	Images      string     `gorm:"type:text" json:"images"`
	IsAnonymous bool       `gorm:"default:false" json:"is_anonymous"`
	Likes       int        `gorm:"default:0" json:"likes"`
	IsHidden    bool       `gorm:"default:false" json:"is_hidden"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Review) TableName() string {
	return "reviews"
}

type MaintenanceStatus string

const (
	MaintenanceStatusScheduled MaintenanceStatus = "scheduled"
	MaintenanceStatusInProgress MaintenanceStatus = "in_progress"
	MaintenanceStatusCompleted  MaintenanceStatus = "completed"
	MaintenanceStatusCancelled  MaintenanceStatus = "cancelled"
)

type MaintenancePlan struct {
	ID           uint             `gorm:"primarykey" json:"id"`
	CarID        uint             `gorm:"index;not null" json:"car_id"`
	Car          *Car             `gorm:"foreignKey:CarID" json:"car,omitempty"`
	Title        string           `gorm:"size:200;not null" json:"title"`
	Description  string           `gorm:"type:text" json:"description"`
	StartDate    time.Time        `gorm:"not null" json:"start_date"`
	EndDate      time.Time        `gorm:"not null" json:"end_date"`
	ActualStart  *time.Time       `json:"actual_start"`
	ActualEnd    *time.Time       `json:"actual_end"`
	Cost         float64          `gorm:"type:decimal(10,2);default:0" json:"cost"`
	Status       MaintenanceStatus `gorm:"size:20;default:scheduled" json:"status"`
	Notes        string           `gorm:"type:text" json:"notes"`
	CreatedBy    uint             `json:"created_by"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func (MaintenancePlan) TableName() string {
	return "maintenance_plans"
}

type MessageType string

const (
	MessageTypeBooking     MessageType = "booking"
	MessageTypePickup      MessageType = "pickup"
	MessageTypeReturn      MessageType = "return"
	MessageTypeReview      MessageType = "review"
	MessageTypeMaintenance MessageType = "maintenance"
	MessageTypeSystem      MessageType = "system"
)

type Message struct {
	ID        uint        `gorm:"primarykey" json:"id"`
	UserID    uint        `gorm:"index;not null" json:"user_id"`
	User      *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      MessageType `gorm:"size:20;not null" json:"type"`
	Title     string      `gorm:"size:200;not null" json:"title"`
	Content   string      `gorm:"type:text;not null" json:"content"`
	RelatedID *uint       `json:"related_id"`
	IsRead    bool        `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time   `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

type OperationLog struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Module      string     `gorm:"size:50;not null" json:"module"`
	Action      string     `gorm:"size:50;not null" json:"action"`
	TargetType  string     `gorm:"size:50" json:"target_type"`
	TargetID    *uint      `json:"target_id"`
	Description string     `gorm:"type:text" json:"description"`
	IP          string     `gorm:"size:50" json:"ip"`
	UserAgent   string     `gorm:"size:500" json:"user_agent"`
	Method      string     `gorm:"size:10" json:"method"`
	Path        string     `gorm:"size:500" json:"path"`
	Status      int        `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Channel   string    `gorm:"size:20;not null" json:"channel"`
	IsSent    bool      `gorm:"default:false" json:"is_sent"`
	SentAt    *time.Time `json:"sent_at"`
	ErrorMsg  string    `gorm:"size:500" json:"error_msg"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}