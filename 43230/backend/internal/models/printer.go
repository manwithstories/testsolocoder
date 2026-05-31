package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrinterStatus string

const (
	PrinterStatusIdle      PrinterStatus = "idle"
	PrinterStatusPrinting  PrinterStatus = "printing"
	PrinterStatusMaintenance PrinterStatus = "maintenance"
	PrinterStatusOffline   PrinterStatus = "offline"
)

type PrinterDevice struct {
	ID               uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey"`
	PrinterID        uuid.UUID     `json:"printer_id" gorm:"type:uuid;not null;index"`
	Printer          *User         `json:"printer,omitempty" gorm:"foreignKey:PrinterID"`
	Name             string        `json:"name" gorm:"not null"`
	Model            string        `json:"model"`
	Manufacturer     string        `json:"manufacturer"`
	MaxPrintSize     string        `json:"max_print_size"`
	MaxPrintVolume   float64       `json:"max_print_volume"`
	SupportedMaterials []string    `json:"supported_materials" gorm:"type:text[]"`
	SupportedQualities []PrintQuality `json:"supported_qualities" gorm:"type:text[]"`
	Status           PrinterStatus `json:"status" gorm:"default:idle"`
	CurrentOrderID   *uuid.UUID    `json:"current_order_id" gorm:"type:uuid"`
	TotalPrintHours  float64       `json:"total_print_hours" gorm:"default:0"`
	TotalPrintJobs   int           `json:"total_print_jobs" gorm:"default:0"`
	LastMaintenance  *time.Time    `json:"last_maintenance"`
	NextMaintenance  *time.Time    `json:"next_maintenance"`
	IPAddress        string        `json:"ip_address"`
	FirmwareVersion  string        `json:"firmware_version"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type MaterialInventory struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	PrinterID    uuid.UUID    `json:"printer_id" gorm:"type:uuid;not null;index"`
	Printer      *User        `json:"printer,omitempty" gorm:"foreignKey:PrinterID"`
	MaterialID   uuid.UUID    `json:"material_id" gorm:"type:uuid;not null;index"`
	Material     *Material    `json:"material,omitempty" gorm:"foreignKey:MaterialID"`
	Color        string       `json:"color" gorm:"not null"`
	QuantityGrams float64     `json:"quantity_grams" gorm:"not null"`
	ReorderLevel  float64     `json:"reorder_level" gorm:"default:500"`
	LastUpdated   time.Time    `json:"last_updated"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type PrintSchedule struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	PrinterID     uuid.UUID   `json:"printer_id" gorm:"type:uuid;not null;index"`
	DeviceID      uuid.UUID   `json:"device_id" gorm:"type:uuid;not null;index"`
	Device        *PrinterDevice `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	OrderID       uuid.UUID   `json:"order_id" gorm:"type:uuid;not null;index"`
	Order         *PrintOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ScheduledStart *time.Time `json:"scheduled_start"`
	ScheduledEnd   *time.Time `json:"scheduled_end"`
	ActualStart    *time.Time `json:"actual_start"`
	ActualEnd      *time.Time `json:"actual_end"`
	Status         string      `json:"status" gorm:"default:scheduled"`
	Priority       int         `json:"priority" gorm:"default:0"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type Review struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	OrderID       uuid.UUID   `json:"order_id" gorm:"type:uuid;uniqueIndex;not null"`
	Order         *PrintOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	CustomerID    uuid.UUID   `json:"customer_id" gorm:"type:uuid;not null;index"`
	Customer      *User       `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	ModelID       uuid.UUID   `json:"model_id" gorm:"type:uuid;not null;index"`
	Model         *Model3D    `json:"model,omitempty" gorm:"foreignKey:ModelID"`
	PrinterID     uuid.UUID   `json:"printer_id" gorm:"type:uuid;not null;index"`
	Printer       *User       `json:"printer,omitempty" gorm:"foreignKey:PrinterID"`
	ModelRating   int         `json:"model_rating" gorm:"check:model_rating >= 1 AND model_rating <= 5"`
	PrintRating   int         `json:"print_rating" gorm:"check:print_rating >= 1 AND print_rating <= 5"`
	ModelComment  string      `json:"model_comment"`
	PrintComment  string      `json:"print_comment"`
	Images        []string    `json:"images" gorm:"type:text[]"`
	IsAnonymous   bool        `json:"is_anonymous" gorm:"default:false"`
	DesignerReply string      `json:"designer_reply"`
	PrinterReply  string      `json:"printer_reply"`
	HelpfulCount  int         `json:"helpful_count" gorm:"default:0"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
