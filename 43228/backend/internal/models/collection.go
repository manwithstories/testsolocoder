package models

import (
	"time"

	"gorm.io/gorm"
)

type Collection struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"index;not null" json:"user_id"`
	TeaID             uint           `gorm:"index;not null" json:"tea_id"`
	Quantity          int            `gorm:"not null;default:0" json:"quantity"`
	StorageLocation   string         `gorm:"size:255" json:"storage_location"`
	EnvironmentParams string         `gorm:"type:json" json:"environment_params"`
	PurchaseDate      *time.Time     `json:"purchase_date"`
	ExpiryDate        *time.Time     `json:"expiry_date"`
	MinStock          int            `gorm:"not null;default:0" json:"min_stock"`
	AlertEnabled      bool           `gorm:"not null;default:false" json:"alert_enabled"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tea  *Tea  `gorm:"foreignKey:TeaID" json:"tea,omitempty"`
}

func (Collection) TableName() string {
	return "collections"
}
