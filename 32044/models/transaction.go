package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TransactionTypeIncome  = "income"
	TransactionTypeExpense = "expense"
)

type Transaction struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	AccountID  uint           `gorm:"not null;index" json:"account_id" binding:"required"`
	CategoryID uint           `gorm:"not null;index" json:"category_id" binding:"required"`
	Type       string         `gorm:"size:10;not null;index" json:"type" binding:"required,oneof=income expense"`
	Amount     float64        `gorm:"type:decimal(15,2);not null" json:"amount" binding:"required,gt=0"`
	Remark     string         `gorm:"size:255" json:"remark" binding:"max=255"`
	Date       time.Time      `gorm:"not null;index" json:"date"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Account  *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
