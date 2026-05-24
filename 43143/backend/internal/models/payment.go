package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentMethodAlipay PaymentMethod = "alipay"
	PaymentMethodWechat PaymentMethod = "wechat"
	PaymentMethodBank   PaymentMethod = "bank_transfer"
	PaymentMethodCash   PaymentMethod = "cash"
)

type TransactionType string

const (
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeWithdraw   TransactionType = "withdraw"
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeFee        TransactionType = "fee"
)

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusProcessing TransactionStatus = "processing"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
)

type Payment struct {
	ID              uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID       uuid.UUID        `gorm:"type:uuid;index" json:"booking_id"`
	Booking         Booking          `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	PayerID         uuid.UUID        `gorm:"type:uuid;index;not null" json:"payer_id"`
	Payer           User             `gorm:"foreignKey:PayerID" json:"payer,omitempty"`
	PayeeID         uuid.UUID        `gorm:"type:uuid;index;not null" json:"payee_id"`
	Payee           User             `gorm:"foreignKey:PayeeID" json:"payee,omitempty"`
	Amount          float64          `gorm:"not null" json:"amount"`
	Currency        string           `gorm:"size:10;default:'CNY'" json:"currency"`
	Method          PaymentMethod    `gorm:"type:varchar(20);not null" json:"method"`
	Type            TransactionType  `gorm:"type:varchar(20);not null" json:"type"`
	Status          TransactionStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PlatformFee     float64          `gorm:"default:0" json:"platform_fee"`
	NetAmount       float64          `gorm:"default:0" json:"net_amount"`
	EscrowReleaseAt *time.Time       `json:"escrow_release_at"`
	Released        bool             `gorm:"default:false" json:"released"`
	ReleasedAt      *time.Time       `json:"released_at"`
	TransactionID   string           `gorm:"size:100" json:"transaction_id"`
	ThirdPartyID    string           `gorm:"size:100" json:"third_party_id"`
	FailureReason   string           `gorm:"size:500" json:"failure_reason"`
	Remark          string           `gorm:"size:500" json:"remark"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Wallet struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Balance    float64        `gorm:"default:0" json:"balance"`
	Frozen     float64        `gorm:"default:0" json:"frozen"`
	Currency   string         `gorm:"size:10;default:'CNY'" json:"currency"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
