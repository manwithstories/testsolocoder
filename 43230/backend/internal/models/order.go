package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending      OrderStatus = "pending"
	OrderStatusPaid         OrderStatus = "paid"
	OrderStatusPrinting     OrderStatus = "printing"
	OrderStatusQualityCheck OrderStatus = "quality_check"
	OrderStatusShipped      OrderStatus = "shipped"
	OrderStatusDelivered    OrderStatus = "delivered"
	OrderStatusCompleted    OrderStatus = "completed"
	OrderStatusCancelled    OrderStatus = "cancelled"
	OrderStatusRefunded     OrderStatus = "refunded"
)

type PrintQuality string

const (
	QualityDraft   PrintQuality = "draft"
	QualityStandard PrintQuality = "standard"
	QualityHigh    PrintQuality = "high"
	QualityUltra   PrintQuality = "ultra"
)

type MaterialType string

const (
	MaterialPLA   MaterialType = "pla"
	MaterialABS   MaterialType = "abs"
	MaterialPETG  MaterialType = "petg"
	MaterialTPU   MaterialType = "tpu"
	MaterialResin MaterialType = "resin"
	MaterialNylon MaterialType = "nylon"
	MaterialPC    MaterialType = "pc"
)

type ColorOption struct {
	ID   uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Name string     `json:"name"`
	Code string     `json:"code"`
}

type Material struct {
	ID             uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	Type           MaterialType `json:"type" gorm:"not null"`
	Name           string       `json:"name" gorm:"not null"`
	Description    string       `json:"description"`
	ColorOptions   []ColorOption `json:"color_options" gorm:"many2many:material_colors"`
	PricePerGram   float64      `json:"price_per_gram" gorm:"not null"`
	PrintSpeed     float64      `json:"print_speed"`
	Density        float64      `json:"density"`
	Strength       string       `json:"strength"`
	TemperatureResistance string `json:"temperature_resistance"`
	IsAvailable    bool         `json:"is_available" gorm:"default:true"`
	ImageURL       string       `json:"image_url"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type PrintOrder struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	OrderNo         string         `json:"order_no" gorm:"uniqueIndex;not null"`
	CustomerID      uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null;index"`
	PrinterID       uuid.UUID      `json:"printer_id" gorm:"type:uuid;index"`
	ModelID         uuid.UUID      `json:"model_id" gorm:"type:uuid;not null;index"`
	Customer        *User          `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Printer         *User          `json:"printer,omitempty" gorm:"foreignKey:PrinterID"`
	Model           *Model3D       `json:"model,omitempty" gorm:"foreignKey:ModelID"`
	Quantity        int            `json:"quantity" gorm:"not null;default:1"`
	MaterialID      uuid.UUID      `json:"material_id" gorm:"type:uuid;not null"`
	Material        *Material      `json:"material,omitempty" gorm:"foreignKey:MaterialID"`
	Color           string         `json:"color" gorm:"not null"`
	Quality         PrintQuality   `json:"quality" gorm:"not null"`
	LayerHeight     float64        `json:"layer_height"`
	InfillPercent   int            `json:"infill_percent" gorm:"default:20"`
	Supports        bool           `json:"supports" gorm:"default:false"`
	EstimatedVolume float64        `json:"estimated_volume"`
	EstimatedWeight float64        `json:"estimated_weight"`
	EstimatedPrintTime float64     `json:"estimated_print_time"`
	BasePrice       float64        `json:"base_price"`
	MaterialCost    float64        `json:"material_cost"`
	ServiceFee      float64        `json:"service_fee"`
	ShippingFee     float64        `json:"shipping_fee"`
	TotalAmount     float64        `json:"total_amount"`
	Status          OrderStatus    `json:"status" gorm:"default:pending"`
	ShippingAddress string         `json:"shipping_address"`
	TrackingNumber  string         `json:"tracking_number"`
	Notes           string         `json:"notes"`
	CancelledReason string         `json:"cancelled_reason"`
	PrintStartedAt  *time.Time     `json:"print_started_at"`
	PrintFinishedAt *time.Time     `json:"print_finished_at"`
	ShippedAt       *time.Time     `json:"shipped_at"`
	DeliveredAt     *time.Time     `json:"delivered_at"`
	CompletedAt     *time.Time     `json:"completed_at"`
	CancelledAt     *time.Time     `json:"cancelled_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type OrderHistory struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	OrderID     uuid.UUID   `json:"order_id" gorm:"type:uuid;not null;index"`
	Order       *PrintOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Status      OrderStatus `json:"status"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
}

type Settlement struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	OrderID       uuid.UUID   `json:"order_id" gorm:"type:uuid;uniqueIndex;not null"`
	Order         *PrintOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	TotalAmount   float64     `json:"total_amount"`
	PlatformFee   float64     `json:"platform_fee"`
	DesignerShare float64     `json:"designer_share"`
	PrinterShare  float64     `json:"printer_share"`
	Status        string      `json:"status" gorm:"default:pending"`
	SettledAt     *time.Time  `json:"settled_at"`
	CreatedAt     time.Time   `json:"created_at"`
}
