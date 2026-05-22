package models

import (
	"time"

	"gorm.io/gorm"
)

type ServiceCategory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null;uniqueIndex"`
	Description string         `json:"description" gorm:"size:500"`
	Icon        string         `json:"icon" gorm:"size:255"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	ServiceItems []ServiceItem `json:"service_items,omitempty" gorm:"foreignKey:CategoryID"`
}

type ServiceItem struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	CategoryID    uint           `json:"category_id" gorm:"not null;index"`
	ProviderID    uint           `json:"provider_id" gorm:"not null;index"`
	Name          string         `json:"name" gorm:"size:100;not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Images        string         `json:"images" gorm:"type:text"`
	BasePrice     float64        `json:"base_price" gorm:"not null"`
	PriceUnit     string         `json:"price_unit" gorm:"size:20;default:hour"`
	MinDuration   int            `json:"min_duration" gorm:"default:60"`
	MaxDuration   int            `json:"max_duration" gorm:"default:480"`
	Rating        float64        `json:"rating" gorm:"default:5.0"`
	ReviewCount   int            `json:"review_count" gorm:"default:0"`
	OrderCount    int            `json:"order_count" gorm:"default:0"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	Category    *ServiceCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Provider    *User            `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	ServiceAreas []ServiceArea   `json:"service_areas,omitempty" gorm:"many2many:service_item_areas;"`
}

type ServiceArea struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Province  string         `json:"province" gorm:"size:50;not null"`
	City      string         `json:"city" gorm:"size:50;not null"`
	District  string         `json:"district" gorm:"size:50;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ServiceItemArea struct {
	ServiceItemID uint `json:"service_item_id" gorm:"primaryKey"`
	ServiceAreaID uint `json:"service_area_id" gorm:"primaryKey"`
}
