package models

import (
	"time"

	"gorm.io/gorm"
)

type BillType string

const (
	BillTypeIncome      BillType = "income"
	BillTypeWithdraw    BillType = "withdraw"
	BillTypePenalty     BillType = "penalty"
	BillTypeRefund      BillType = "refund"
	BillTypeCommission  BillType = "commission"
)

type BillStatus string

const (
	BillStatusPending   BillStatus = "pending"
	BillStatusCompleted BillStatus = "completed"
	BillStatusFailed    BillStatus = "failed"
)

type Bill struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ProviderID    uint           `json:"provider_id" gorm:"not null;index"`
	OrderID       *uint          `json:"order_id" gorm:"index"`
	BillType      BillType       `json:"bill_type" gorm:"size:20;not null"`
	Amount        float64        `json:"amount" gorm:"not null"`
	Balance       float64        `json:"balance" gorm:"not null"`
	Status        BillStatus     `json:"status" gorm:"size:20;not null;default:completed"`
	Description   string         `json:"description" gorm:"size:500"`
	TransactionNo string         `json:"transaction_no" gorm:"size:64"`
	SettledAt     *time.Time     `json:"settled_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	Provider      *User          `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	Order         *Order         `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

type WithdrawStatus string

const (
	WithdrawStatusPending   WithdrawStatus = "pending"
	WithdrawStatusApproved  WithdrawStatus = "approved"
	WithdrawStatusRejected  WithdrawStatus = "rejected"
	WithdrawStatusCompleted WithdrawStatus = "completed"
	WithdrawStatusFailed    WithdrawStatus = "failed"
)

type WithdrawRequest struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	ProviderID    uint             `json:"provider_id" gorm:"not null;index"`
	Amount        float64          `json:"amount" gorm:"not null"`
	BankName      string           `json:"bank_name" gorm:"size:50"`
	BankAccount   string           `json:"bank_account" gorm:"size:50"`
	AccountHolder string           `json:"account_holder" gorm:"size:50"`
	Status        WithdrawStatus   `json:"status" gorm:"size:20;not null;default:pending"`
	HandlerID     *uint            `json:"handler_id" gorm:"index"`
	HandleRemark  string           `json:"handle_remark" gorm:"size:500"`
	TransferNo    string           `json:"transfer_no" gorm:"size:64"`
	HandledAt     *time.Time       `json:"handled_at"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `json:"-" gorm:"index"`

	Provider      *User            `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
}
