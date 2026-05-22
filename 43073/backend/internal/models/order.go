package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusCanceled  = "canceled"
	OrderStatusRefunded  = "refunded"
)

type Order struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderNo      string         `gorm:"size:50;uniqueIndex;not null" json:"orderNo"`
	UserID       uint           `json:"userId"`
	ActivityID   uint           `json:"activityId"`
	CouponID     *uint          `json:"couponId"`
	TotalAmount  float64        `json:"totalAmount"`
	Discount     float64        `json:"discount"`
	PayAmount    float64        `json:"payAmount"`
	Status       string         `gorm:"size:20;default:'pending'" json:"status"`
	Remark       string         `gorm:"size:500" json:"remark"`
	PaidAt       *time.Time     `json:"paidAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Activity     *Activity      `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
	Coupon       *Coupon        `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
	OrderItems   []OrderItem    `gorm:"foreignKey:OrderID" json:"orderItems,omitempty"`
}

type OrderItem struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderID      uint           `json:"orderId"`
	TicketTypeID uint           `json:"ticketTypeId"`
	Quantity     int            `json:"quantity"`
	UnitPrice    float64        `json:"unitPrice"`
	Subtotal     float64        `json:"subtotal"`
	CreatedAt    time.Time      `json:"createdAt"`
	TicketType   *TicketType    `gorm:"foreignKey:TicketTypeID" json:"ticketType,omitempty"`
	CheckIns     []CheckIn      `gorm:"foreignKey:OrderItemID" json:"checkIns,omitempty"`
}
