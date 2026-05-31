package models

import (
	"time"

	"gorm.io/gorm"
)

type Harvest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	BeehiveID   uint           `gorm:"not null;index" json:"beehive_id"`
	HarvestDate time.Time      `gorm:"type:date;not null" json:"harvest_date"`
	HoneyType   string         `gorm:"size:50;not null" json:"honey_type"`
	Quantity    float64        `gorm:"type:decimal(10,2);not null" json:"quantity"`
	Unit        string         `gorm:"size:10;default:kg" json:"unit"`
	Quality     string         `gorm:"size:20;default:normal" json:"quality"`
	BatchCode   string         `gorm:"uniqueIndex;size:50;not null" json:"batch_code"`
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	Beehive     Beehive        `gorm:"foreignKey:BeehiveID" json:"beehive,omitempty"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Harvest) TableName() string {
	return "harvests"
}

type CreateHarvestRequest struct {
	BeehiveID   uint    `json:"beehive_id" binding:"required"`
	HarvestDate string  `json:"harvest_date" binding:"required"`
	HoneyType   string  `json:"honey_type" binding:"required,max=50"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Unit        string  `json:"unit" binding:"max=10"`
	Quality     string  `json:"quality"`
	BatchCode   string  `json:"batch_code" binding:"required,max=50"`
	Notes       string  `json:"notes"`
}

type Inventory struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"not null;index" json:"user_id"`
	HarvestID        uint           `gorm:"not null;index" json:"harvest_id"`
	HoneyType        string         `gorm:"size:50;not null" json:"honey_type"`
	BatchCode        string         `gorm:"size:50;not null;index" json:"batch_code"`
	Quantity         float64        `gorm:"type:decimal(10,2);not null" json:"quantity"`
	Unit             string         `gorm:"size:10;default:kg" json:"unit"`
	ExpiryDate       time.Time      `gorm:"type:date;not null" json:"expiry_date"`
	InspectionReport string         `gorm:"size:255" json:"inspection_report"`
	Grade            string         `gorm:"size:20;default:ungraded" json:"grade"`
	Status           string         `gorm:"size:20;default:in_stock;index" json:"status"`
	Threshold        float64        `gorm:"type:decimal(10,2);default:10" json:"threshold"`
	Price            float64        `gorm:"type:decimal(10,2)" json:"price"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	Harvest          Harvest        `gorm:"foreignKey:HarvestID" json:"harvest,omitempty"`
	User             User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Inventory) TableName() string {
	return "inventory"
}

type UpdateInventoryRequest struct {
	Quantity         *float64 `json:"quantity" binding:"omitempty,gte=0"`
	ExpiryDate       *string  `json:"expiry_date"`
	InspectionReport *string  `json:"inspection_report"`
	Grade            *string  `json:"grade"`
	Threshold        *float64 `json:"threshold" binding:"omitempty,gte=0"`
	Price            *float64 `json:"price" binding:"omitempty,gt=0"`
}
