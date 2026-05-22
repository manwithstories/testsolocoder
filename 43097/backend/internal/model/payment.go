package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderType string

const (
	OrderTypeBooking OrderType = "booking"
	OrderTypeCheckIn OrderType = "checkin"
)

type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodWeChat PaymentMethod = "wechat"
	PaymentMethodAlipay PaymentMethod = "alipay"
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodTransfer PaymentMethod = "transfer"
)

type PaymentType string

const (
	PaymentTypePrepaid PaymentType = "prepaid"
	PaymentTypeExtra   PaymentType = "extra"
	PaymentTypeRefund  PaymentType = "refund"
	PaymentTypeDeposit PaymentType = "deposit"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PaymentNo     string         `gorm:"size:32;not null;uniqueIndex" json:"payment_no"`
	OrderType     OrderType      `gorm:"size:20;not null;index" json:"order_type"`
	OrderID       uint           `gorm:"not null;index" json:"order_id"`
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentMethod PaymentMethod  `gorm:"size:20;not null" json:"payment_method"`
	PaymentType   PaymentType    `gorm:"size:20;not null" json:"payment_type"`
	Status        PaymentStatus  `gorm:"size:20;not null;default:'pending'" json:"status"`
	TransactionID string         `gorm:"size:100;index" json:"transaction_id"`
	Remark        string         `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
