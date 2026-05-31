package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusRefunded  TransactionStatus = "refunded"
)

type Transaction struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RentalID        *uuid.UUID     `gorm:"type:uuid;index" json:"rental_id"`
	Rental          *Rental        `gorm:"foreignKey:RentalID" json:"rental,omitempty"`
	PayerID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"payer_id"`
	Payer           User           `gorm:"foreignKey:PayerID" json:"payer,omitempty"`
	PayeeID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"payee_id"`
	Payee           User           `gorm:"foreignKey:PayeeID" json:"payee,omitempty"`
	Amount          float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string         `gorm:"size:10;default:USD" json:"currency"`
	TransactionType TransactionType `gorm:"type:varchar(20);not null" json:"transaction_type"`
	Description     string         `gorm:"size:500" json:"description"`
	Status          TransactionStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	PaymentMethod   string         `gorm:"type:varchar(50)" json:"payment_method"`
	TransactionRef  string         `gorm:"size:100" json:"transaction_ref"`
	ExchangeRate    float64        `gorm:"type:decimal(10,6);default:1" json:"exchange_rate"`
	OriginalAmount  float64        `gorm:"type:decimal(12,2)" json:"original_amount"`
	OriginalCurrency string        `gorm:"size:10" json:"original_currency"`
	PlatformFee     float64        `gorm:"type:decimal(12,2);default:0" json:"platform_fee"`
	DockFee         float64        `gorm:"type:decimal(12,2);default:0" json:"dock_fee"`
	NetAmount       float64        `gorm:"type:decimal(12,2)" json:"net_amount"`
	PaidAt          *time.Time     `json:"paid_at"`
	FailedReason    string         `gorm:"type:text" json:"failed_reason"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Settlement struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PeriodStart   time.Time      `gorm:"not null" json:"period_start"`
	PeriodEnd     time.Time      `gorm:"not null" json:"period_end"`
	TotalIncome   float64        `gorm:"type:decimal(12,2);default:0" json:"total_income"`
	TotalExpense  float64        `gorm:"type:decimal(12,2);default:0" json:"total_expense"`
	Currency      string         `gorm:"size:10;default:USD" json:"currency"`
	Status        string         `gorm:"type:varchar(20);default:pending" json:"status"`
	PaidAt        *time.Time     `json:"paid_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (s *Settlement) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type CreateTransactionRequest struct {
	RentalID        string          `json:"rental_id"`
	PayerID         string          `json:"payer_id" binding:"required,uuid"`
	PayeeID         string          `json:"payee_id" binding:"required,uuid"`
	Amount          float64         `json:"amount" binding:"required,gt=0"`
	Currency        string          `json:"currency"`
	TransactionType TransactionType `json:"transaction_type" binding:"required,oneof=income expense"`
	Description     string          `json:"description"`
	PaymentMethod   string          `json:"payment_method"`
}

type FinancialReportRequest struct {
	UserID      string `form:"user_id"`
	StartDate   string `form:"start_date" binding:"required"`
	EndDate     string `form:"end_date" binding:"required"`
	Currency    string `form:"currency"`
	TransactionType string `form:"transaction_type"`
	Format      string `form:"format" binding:"required,oneof=pdf csv"`
}

type MonthlyReportRequest struct {
	Year   int    `form:"year" binding:"required"`
	Month  int    `form:"month" binding:"required,min=1,max=12"`
	UserID string `form:"user_id"`
	Format string `form:"format" binding:"required,oneof=pdf csv"`
}
