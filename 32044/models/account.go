package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	Name          string         `gorm:"size:100;not null" json:"name" binding:"required,max=100"`
	Currency      string         `gorm:"size:10;not null;default:'CNY'" json:"currency" binding:"required,max=10"`
	InitialBalance float64       `gorm:"type:decimal(15,2);not null;default:0" json:"initial_balance"`
	Balance       float64        `gorm:"type:decimal(15,2);not null;default:0" json:"balance"`
	Remark        string         `gorm:"size:255" json:"remark" binding:"max=255"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"-"`
}
