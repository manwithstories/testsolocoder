package models

import (
	"time"

	"gorm.io/gorm"
)

type Inspection struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	InspectorID     *uint          `gorm:"index" json:"inspector_id"`
	InventoryID     uint           `gorm:"not null;index" json:"inventory_id"`
	BatchCode       string         `gorm:"size:50;not null" json:"batch_code"`
	AppointmentDate time.Time      `gorm:"type:date;not null" json:"appointment_date"`
	InspectionDate  *time.Time     `gorm:"type:date" json:"inspection_date"`
	Status          string         `gorm:"size:20;default:pending" json:"status"`
	ReportURL       string         `gorm:"size:255" json:"report_url"`
	Result          string         `gorm:"size:20" json:"result"`
	Grade           string         `gorm:"size:20" json:"grade"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Inspector       *User          `gorm:"foreignKey:InspectorID" json:"inspector,omitempty"`
	Inventory       Inventory      `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
}

func (Inspection) TableName() string {
	return "inspections"
}

type CreateInspectionRequest struct {
	InventoryID     uint   `json:"inventory_id" binding:"required"`
	BatchCode       string `json:"batch_code" binding:"required,max=50"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	Notes           string `json:"notes"`
}

type UpdateInspectionRequest struct {
	Status         *string `json:"status"`
	ReportURL      *string `json:"report_url"`
	Result         *string `json:"result"`
	Grade          *string `json:"grade"`
	InspectionDate *string `json:"inspection_date"`
	Notes          *string `json:"notes"`
}

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	InventoryID uint           `gorm:"not null;index" json:"inventory_id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	HoneyType   string         `gorm:"size:50;not null" json:"honey_type"`
	BatchCode   string         `gorm:"size:50" json:"batch_code"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       float64        `gorm:"type:decimal(10,2);not null" json:"stock"`
	Unit        string         `gorm:"size:10;default:kg" json:"unit"`
	Images      []string       `gorm:"serializer:json" json:"images"`
	Grade       string         `gorm:"size:20" json:"grade"`
	Status      string         `gorm:"size:20;default:on_sale;index" json:"status"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Inventory   Inventory      `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

type CreateProductRequest struct {
	InventoryID uint     `json:"inventory_id" binding:"required"`
	Title       string   `json:"title" binding:"required,max=200"`
	Description string   `json:"description"`
	HoneyType   string   `json:"honey_type" binding:"required,max=50"`
	BatchCode   string   `json:"batch_code" binding:"max=50"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Stock       float64  `json:"stock" binding:"required,gte=0"`
	Unit        string   `json:"unit" binding:"max=10"`
	Images      []string `json:"images"`
	Grade       string   `json:"grade"`
}

type UpdateProductRequest struct {
	Title       *string   `json:"title" binding:"omitempty,max=200"`
	Description *string   `json:"description"`
	Price       *float64  `json:"price" binding:"omitempty,gt=0"`
	Stock       *float64  `json:"stock" binding:"omitempty,gte=0"`
	Images      *[]string `json:"images"`
	Status      *string   `json:"status"`
}
