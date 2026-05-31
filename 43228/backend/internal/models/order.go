package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusShipped        OrderStatus = "shipped"
	OrderStatusCompleted      OrderStatus = "completed"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusRefunding      OrderStatus = "refunding"
	OrderStatusRefunded       OrderStatus = "refunded"
)

type Order struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderNo         string         `gorm:"size:64;uniqueIndex;not null" json:"order_no"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	TeaID           uint           `gorm:"index;not null" json:"tea_id"`
	Quantity        int            `gorm:"not null;default:1" json:"quantity"`
	UnitPrice       float64        `gorm:"type:decimal(12,2);not null" json:"unit_price"`
	TotalPrice      float64        `gorm:"type:decimal(12,2);not null" json:"total_price"`
	Status          OrderStatus    `gorm:"size:32;not null;default:pending_payment;index" json:"status"`
	PaymentMethod   string         `gorm:"size:64" json:"payment_method"`
	PaymentTime     *time.Time     `json:"payment_time"`
	ShippingAddress string         `gorm:"size:512" json:"shipping_address"`
	TrackingNumber  string         `gorm:"size:64" json:"tracking_number"`
	ShippedTime     *time.Time     `json:"shipped_time"`
	CompletedTime   *time.Time     `json:"completed_time"`
	CancelReason    string         `gorm:"size:255" json:"cancel_reason"`
	RefundReason    string         `gorm:"size:255" json:"refund_reason"`
	Remark          string         `gorm:"size:512" json:"remark"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User      *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tea       *Tea        `gorm:"foreignKey:TeaID" json:"tea,omitempty"`
	OrderLogs []OrderLog  `gorm:"foreignKey:OrderID" json:"order_logs,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
