package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Order struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	OrderNo       string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID        uint           `json:"user_id" gorm:"index;not null"`
	User          *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	TotalAmount   float64        `json:"total_amount" gorm:"type:decimal(12,2);not null"`
	Status        OrderStatus    `json:"status" gorm:"size:20;default:pending"`
	PaymentStatus PaymentStatus  `json:"payment_status" gorm:"size:20;default:pending"`
	ReceiverName  string         `json:"receiver_name" gorm:"size:50"`
	ReceiverPhone string         `json:"receiver_phone" gorm:"size:20"`
	Address       string         `json:"address" gorm:"size:500"`
	Remark        string         `json:"remark" gorm:"size:500"`
	Items         []OrderItem    `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	PaidAt        *time.Time     `json:"paid_at"`
	ShippedAt     *time.Time     `json:"shipped_at"`
	DeliveredAt   *time.Time     `json:"delivered_at"`
	CancelledAt   *time.Time     `json:"cancelled_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type OrderItem struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	OrderID    uint           `json:"order_id" gorm:"index;not null"`
	ProductID  uint           `json:"product_id" gorm:"not null"`
	Product    *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ProductName string        `json:"product_name" gorm:"size:200"`
	Price      float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity   int            `json:"quantity" gorm:"not null"`
	Subtotal   float64        `json:"subtotal" gorm:"type:decimal(12,2);not null"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type CartItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	ProductID uint           `json:"product_id" gorm:"index;not null"`
	Product   *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int            `json:"quantity" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateOrderRequest struct {
	ReceiverName  string         `json:"receiver_name" binding:"required"`
	ReceiverPhone string         `json:"receiver_phone" binding:"required"`
	Address       string         `json:"address" binding:"required"`
	Remark        string         `json:"remark"`
	Items         []CartOrderItem `json:"items" binding:"required"`
}

type CartOrderItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending paid processing shipped delivered cancelled refunded"`
}

type PaymentRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Method  string `json:"method" binding:"required"`
}
