package model

import (
	"time"
)

type PaymentMethod string

const (
	PaymentMethodWeChat PaymentMethod = "wechat"
	PaymentMethodAlipay PaymentMethod = "alipay"
	PaymentMethodCash   PaymentMethod = "cash"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	OrderID       uint          `gorm:"index;not null" json:"order_id"`
	TransactionNo string        `gorm:"size:64" json:"transaction_no"`
	Amount        float64       `gorm:"default:0" json:"amount"`
	PaymentMethod PaymentMethod `gorm:"size:20" json:"payment_method"`
	Status        PaymentStatus `gorm:"size:20;default:'pending'" json:"status"`
	PaidAt        *time.Time    `json:"paid_at"`
	CreatedAt     time.Time     `json:"created_at"`
	Order         *Order        `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}
