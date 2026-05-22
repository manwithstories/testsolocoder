package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending     OrderStatus = "pending"
	OrderStatusConfirmed   OrderStatus = "confirmed"
	OrderStatusInService   OrderStatus = "in_service"
	OrderStatusCompleted   OrderStatus = "completed"
	OrderStatusCancelled   OrderStatus = "cancelled"
	OrderStatusDisputed    OrderStatus = "disputed"
)

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusRejected InvitationStatus = "rejected"
	InvitationStatusExpired  InvitationStatus = "expired"
)

type Order struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	OrderNo          string         `json:"order_no" gorm:"uniqueIndex;size:32;not null"`
	CustomerID       uint           `json:"customer_id" gorm:"not null;index"`
	ProviderID       uint           `json:"provider_id" gorm:"index"`
	ServiceItemID    uint           `json:"service_item_id" gorm:"not null;index"`
	Status           OrderStatus    `json:"status" gorm:"size:20;not null;default:pending"`
	AddressID        uint           `json:"address_id" gorm:"not null"`
	ServiceAddress   string         `json:"service_address" gorm:"size:255;not null"`
	ContactName      string         `json:"contact_name" gorm:"size:50;not null"`
	ContactPhone     string         `json:"contact_phone" gorm:"size:20;not null"`
	Longitude        float64        `json:"longitude" gorm:"default:0"`
	Latitude         float64        `json:"latitude" gorm:"default:0"`
	AppointmentTime  time.Time      `json:"appointment_time" gorm:"not null;index"`
	ActualStartTime  *time.Time     `json:"actual_start_time"`
	ActualEndTime    *time.Time     `json:"actual_end_time"`
	StartLocation    string         `json:"start_location" gorm:"size:255"`
	EndLocation      string         `json:"end_location" gorm:"size:255"`
	Duration         int            `json:"duration" gorm:"default:0"`
	BasePrice        float64        `json:"base_price" gorm:"not null"`
	TotalAmount      float64        `json:"total_amount" gorm:"not null"`
	PlatformFee      float64        `json:"platform_fee" gorm:"default:0"`
	ProviderIncome   float64        `json:"provider_income" gorm:"default:0"`
	PenaltyAmount    float64        `json:"penalty_amount" gorm:"default:0"`
	Remark           string         `json:"remark" gorm:"size:500"`
	CancelReason     string         `json:"cancel_reason" gorm:"size:500"`
	CancelledBy      string         `json:"cancelled_by" gorm:"size:20"`
	CancelledAt      *time.Time     `json:"cancelled_at"`
	CompletedAt      *time.Time     `json:"completed_at"`
	CustomerRated    bool           `json:"customer_rated" gorm:"default:false"`
	ProviderRated    bool           `json:"provider_rated" gorm:"default:false"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	Customer         *User             `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Provider         *User             `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	ServiceItem      *ServiceItem      `json:"service_item,omitempty" gorm:"foreignKey:ServiceItemID"`
	Invitations      []OrderInvitation `json:"invitations,omitempty" gorm:"foreignKey:OrderID"`
	Review           *Review           `json:"review,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderInvitation struct {
	ID              uint              `json:"id" gorm:"primaryKey"`
	OrderID         uint              `json:"order_id" gorm:"not null;index"`
	ProviderID      uint              `json:"provider_id" gorm:"not null;index"`
	Status          InvitationStatus  `json:"status" gorm:"size:20;not null;default:pending"`
	RespondedAt     *time.Time        `json:"responded_at"`
	RejectReason    string            `json:"reject_reason" gorm:"size:500"`
	ExpiredAt       time.Time         `json:"expired_at"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"-" gorm:"index"`

	Provider        *User             `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
}
