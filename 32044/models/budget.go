package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index:idx_user_category_month,unique" json:"user_id"`
	CategoryID uint           `gorm:"not null;index:idx_user_category_month,unique" json:"category_id" binding:"required"`
	Month      string         `gorm:"size:7;not null;index:idx_user_category_month,unique" json:"month" binding:"required,len=7"`
	Limit      float64        `gorm:"type:decimal(15,2);not null" json:"limit" binding:"required,gt=0"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
