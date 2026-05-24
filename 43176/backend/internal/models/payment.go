package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TxTypeDeposit      TransactionType = "deposit"
	TxTypeWithdraw     TransactionType = "withdraw"
	TxTypePayment      TransactionType = "payment"
	TxTypeRefund       TransactionType = "refund"
	TxTypeSettlement   TransactionType = "settlement"
	TxTypeServiceFee   TransactionType = "service_fee"
)

type TransactionStatus string

const (
	TxStatusPending    TransactionStatus = "pending"
	TxStatusCompleted  TransactionStatus = "completed"
	TxStatusFailed     TransactionStatus = "failed"
	TxStatusCancelled  TransactionStatus = "cancelled"
)

type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	OrderID         *uint             `gorm:"index" json:"order_id,omitempty"`
	UserID          uint              `gorm:"index" json:"user_id"`
	User            User              `json:"user,omitempty"`
	Type            TransactionType   `gorm:"size:20" json:"type"`
	Amount          float64           `json:"amount"`
	BalanceBefore   float64           `json:"balance_before"`
	BalanceAfter    float64           `json:"balance_after"`
	Status          TransactionStatus `gorm:"size:20;default:pending" json:"status"`
	Description     string            `gorm:"size:500" json:"description"`
	PaymentMethod   string            `gorm:"size:20" json:"payment_method"`
	TransactionNo   string            `gorm:"uniqueIndex;size:50" json:"transaction_no"`
	FailureReason   string            `gorm:"size:500" json:"failure_reason"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
}

type WithdrawRequest struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index" json:"user_id"`
	User          User           `json:"user,omitempty"`
	Amount        float64        `json:"amount"`
	AccountType   string         `gorm:"size:20" json:"account_type"`
	AccountNo     string         `gorm:"size:50" json:"account_no"`
	AccountName   string         `gorm:"size:50" json:"account_name"`
	Status        string         `gorm:"size:20;default:pending" json:"status"`
	Reason        string         `gorm:"size:500" json:"reason"`
	TransactionID *uint          `gorm:"index" json:"transaction_id,omitempty"`
	ReviewedAt    *time.Time     `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type RefundRequest struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderID       uint           `gorm:"index" json:"order_id"`
	Order         Order          `json:"order,omitempty"`
	UserID        uint           `gorm:"index" json:"user_id"`
	User          User           `json:"user,omitempty"`
	Amount        float64        `json:"amount"`
	Reason        string         `gorm:"size:1000" json:"reason"`
	Status        string         `gorm:"size:20;default:pending" json:"status"`
	Decision      string         `gorm:"size:20" json:"decision"`
	DecisionNote  string         `gorm:"size:500" json:"decision_note"`
	TransactionID *uint          `gorm:"index" json:"transaction_id,omitempty"`
	ReviewedAt    *time.Time     `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
