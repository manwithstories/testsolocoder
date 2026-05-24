package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Balance     float64        `gorm:"not null;default:0" json:"balance"`
	Currency    string         `gorm:"default:USD" json:"currency"`
	TotalIncome float64        `gorm:"default:0" json:"totalIncome"`
	TotalSpent  float64        `gorm:"default:0" json:"totalSpent"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction  `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeCommission TransactionType = "commission"
	TransactionTypeTransfer   TransactionType = "transfer"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

type Transaction struct {
	ID            uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	WalletID      uuid.UUID         `gorm:"type:uuid;not null;index" json:"walletId"`
	UserID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"userId"`
	Type          TransactionType   `gorm:"type:varchar(20);not null" json:"type"`
	Amount        float64           `gorm:"not null" json:"amount"`
	Currency      string            `gorm:"default:USD" json:"currency"`
	BalanceAfter  float64           `json:"balanceAfter"`
	Status        TransactionStatus `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	Description   string            `json:"description"`
	ReferenceID   string            `json:"referenceId"`
	ReferenceType string            `json:"referenceType"`
	BookingID     *uuid.UUID        `gorm:"type:uuid" json:"bookingId"`
	PaymentMethod string            `json:"paymentMethod"`
	TransactionID string            `json:"transactionId"`
	GatewayResponse string          `gorm:"type:text" json:"gatewayResponse"`
	FailureReason string            `json:"failureReason"`
	CompletedAt   *time.Time        `json:"completedAt"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`

	Wallet  *Wallet  `gorm:"foreignKey:WalletID" json:"-"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Booking *Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
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
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"userId"`
	WalletID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"walletId"`
	Amount        float64         `gorm:"not null" json:"amount"`
	Currency      string          `gorm:"default:USD" json:"currency"`
	BankName      string          `json:"bankName"`
	BankAccount   string          `json:"bankAccount"`
	AccountHolder string          `json:"accountHolder"`
	Status        WithdrawStatus  `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	AdminNotes    string          `gorm:"type:text" json:"adminNotes"`
	ProcessedAt   *time.Time      `json:"processedAt"`
	ProcessedBy   *uuid.UUID      `gorm:"type:uuid" json:"processedBy"`
	TransactionID *uuid.UUID      `gorm:"type:uuid" json:"transactionId"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Wallet *Wallet `gorm:"foreignKey:WalletID" json:"-"`
}

func (wr *WithdrawRequest) BeforeCreate(tx *gorm.DB) error {
	if wr.ID == uuid.Nil {
		wr.ID = uuid.New()
	}
	return nil
}

type PaymentConfig struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Provider       string         `gorm:"not null" json:"provider"`
	APIKey         string         `json:"apiKey"`
	APISecret      string         `json:"apiSecret"`
	WebhookURL     string         `json:"webhookUrl"`
	IsActive       bool           `gorm:"default:true" json:"isActive"`
	ExtraConfig    string         `gorm:"type:text" json:"extraConfig"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (pc *PaymentConfig) BeforeCreate(tx *gorm.DB) error {
	if pc.ID == uuid.Nil {
		pc.ID = uuid.New()
	}
	return nil
}
