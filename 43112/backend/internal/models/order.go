package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderRefunding OrderStatus = "refunding"
	OrderRefunded  OrderStatus = "refunded"
	OrderCancelled OrderStatus = "cancelled"
	OrderFailed    OrderStatus = "failed"
)

type PayMethod string

const (
	PayMethodAlipay  PayMethod = "alipay"
	PayMethodWechat  PayMethod = "wechat"
	PayMethodBalance PayMethod = "balance"
)

type Order struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID    `gorm:"type:uuid;index;not null" json:"user_id"`
	CourseID      uuid.UUID    `gorm:"type:uuid;index;not null" json:"course_id"`
	OrderNo       string       `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	Amount        float64      `gorm:"not null" json:"amount"`
	OriginalPrice float64      `json:"original_price"`
	Discount      float64      `gorm:"default:0" json:"discount"`
	CouponID      *uuid.UUID   `gorm:"type:uuid" json:"coupon_id,omitempty"`
	Status        OrderStatus  `gorm:"size:20;default:pending;index" json:"status"`
	PayMethod     PayMethod    `gorm:"size:20" json:"pay_method"`
	TransactionNo string       `gorm:"size:100" json:"transaction_no"`
	PaidAt        *time.Time   `json:"paid_at"`
	RefundReason  string       `gorm:"type:text" json:"refund_reason"`
	RefundedAt    *time.Time   `json:"refunded_at"`
	Remark        string       `gorm:"size:500" json:"remark"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Course        *Course      `json:"course,omitempty"`
	User          *User        `json:"-"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
