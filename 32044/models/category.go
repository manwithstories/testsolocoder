package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	CategoryTypeIncome  = "income"
	CategoryTypeExpense = "expense"
)

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"size:50;not null" json:"name" binding:"required,max=50"`
	Type      string         `gorm:"size:10;not null;index" json:"type" binding:"required,oneof=income expense"`
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Parent     *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children   []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID" json:"-"`
}
