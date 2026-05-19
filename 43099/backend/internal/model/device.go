package model

import (
	"time"
)

type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
)

type DeviceCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Device struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	CategoryID      uint         `gorm:"index;not null" json:"category_id"`
	Name            string       `gorm:"size:100;not null" json:"name"`
	Description     string       `gorm:"type:text" json:"description"`
	Specification   string       `gorm:"type:text" json:"specification"`
	StockQuantity   int          `gorm:"default:0" json:"stock_quantity"`
	AvailableQuantity int        `gorm:"default:0" json:"available_quantity"`
	RentalPrice     float64      `gorm:"default:0" json:"rental_price"`
	DepositAmount   float64      `gorm:"default:0" json:"deposit_amount"`
	Status          DeviceStatus `gorm:"size:20;default:'online'" json:"status"`
	Images          string       `gorm:"type:json" json:"-"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Category        *DeviceCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
