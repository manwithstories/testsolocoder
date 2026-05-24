package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	RealName     string    `gorm:"size:64" json:"real_name"`
	Email        string    `gorm:"size:128" json:"email"`
	Phone        string    `gorm:"size:32" json:"phone"`
	RoleID       uint      `gorm:"not null" json:"role_id"`
	Role         Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	CustomerID   *uint     `json:"customer_id,omitempty"`
	Customer     *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Active       bool      `gorm:"default:true" json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"uniqueIndex;size:32;not null" json:"name"`
	Description string `gorm:"size:256" json:"description"`
}

type Customer struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Contact     string    `gorm:"size:64" json:"contact"`
	Phone       string    `gorm:"size:32" json:"phone"`
	Email       string    `gorm:"size:128" json:"email"`
	Address     string    `gorm:"size:256" json:"address"`
	Level       string    `gorm:"size:16;default:normal" json:"level"`
	CreditLimit float64   `gorm:"default:0" json:"credit_limit"`
	Balance     float64   `gorm:"default:0" json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Template struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	Name        string             `gorm:"size:128;not null" json:"name"`
	Category    string             `gorm:"size:32;not null;index" json:"category"`
	WidthMM     float64            `gorm:"not null" json:"width_mm"`
	HeightMM    float64            `gorm:"not null" json:"height_mm"`
	Description string             `gorm:"size:512" json:"description"`
	Thumbnail   string             `gorm:"size:512" json:"thumbnail"`
	Materials   []TemplateMaterial `gorm:"foreignKey:TemplateID" json:"materials,omitempty"`
	Processes   []TemplateProcess  `gorm:"foreignKey:TemplateID" json:"processes,omitempty"`
	Options     []TemplateOption   `gorm:"foreignKey:TemplateID" json:"options,omitempty"`
	Active      bool               `gorm:"default:true" json:"active"`
	CreatedBy   uint               `json:"created_by"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type TemplateMaterial struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	TemplateID  uint    `gorm:"not null;index" json:"template_id"`
	Name        string  `gorm:"size:64;not null" json:"name"`
	Description string  `gorm:"size:256" json:"description"`
	BasePrice   float64 `gorm:"not null" json:"base_price"`
	Unit        string  `gorm:"size:16;default:sheet" json:"unit"`
}

type TemplateProcess struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	TemplateID  uint    `gorm:"not null;index" json:"template_id"`
	Name        string  `gorm:"size:64;not null" json:"name"`
	Description string  `gorm:"size:256" json:"description"`
	ExtraPrice  float64 `gorm:"default:0" json:"extra_price"`
}

type TemplateOption struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TemplateID uint   `gorm:"not null;index" json:"template_id"`
	Type       string `gorm:"size:32;not null" json:"type"`
	Name       string `gorm:"size:64;not null" json:"name"`
	Value      string `gorm:"size:256" json:"value"`
	Required   bool   `gorm:"default:false" json:"required"`
}

type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	OrderNo       string      `gorm:"uniqueIndex;size:32;not null" json:"order_no"`
	CustomerID    uint        `gorm:"not null;index" json:"customer_id"`
	Customer      Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Status        string      `gorm:"size:32;not null;default:created;index" json:"status"`
	Urgent        bool        `gorm:"default:false" json:"urgent"`
	TotalPrice    float64     `gorm:"default:0" json:"total_price"`
	Discount      float64     `gorm:"default:0" json:"discount"`
	FinalPrice    float64     `gorm:"default:0" json:"final_price"`
	Remark        string      `gorm:"size:512" json:"remark"`
	CreatedBy     uint        `json:"created_by"`
	Items         []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Schedules     []ProductionSchedule `gorm:"foreignKey:OrderID" json:"schedules,omitempty"`
	ReviewedBy    *uint       `json:"reviewed_by,omitempty"`
	ReviewedAt    *time.Time  `json:"reviewed_at,omitempty"`
	ProducedAt    *time.Time  `json:"produced_at,omitempty"`
	ShippedAt     *time.Time  `json:"shipped_at,omitempty"`
	CancelledAt   *time.Time  `json:"cancelled_at,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"not null;index" json:"order_id"`
	TemplateID    uint      `gorm:"not null" json:"template_id"`
	Template      Template  `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	MaterialID    uint      `gorm:"not null" json:"material_id"`
	ProcessIDs    string    `gorm:"size:256" json:"process_ids"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	FileAssetID   *uint     `json:"file_asset_id,omitempty"`
	FileAsset     *FileAsset `gorm:"foreignKey:FileAssetID" json:"file_asset,omitempty"`
	UnitPrice     float64   `gorm:"not null" json:"unit_price"`
	SubTotal      float64   `gorm:"not null" json:"sub_total"`
	Specification string    `gorm:"size:512" json:"specification"`
}

type PriceRule struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Category    string  `gorm:"size:32;not null;index" json:"category"`
	MaterialID  *uint   `gorm:"index" json:"material_id,omitempty"`
	MinQty      int     `gorm:"not null" json:"min_qty"`
	MaxQty      int     `gorm:"not null" json:"max_qty"`
	UnitPrice   float64 `gorm:"not null" json:"unit_price"`
	Discount    float64 `gorm:"default:0" json:"discount"`
	CustomerLevel string `gorm:"size:16;default:all" json:"customer_level"`
	Description string  `gorm:"size:256" json:"description"`
}

type ProductionLine struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:64;not null" json:"name"`
	Code       string    `gorm:"uniqueIndex;size:32;not null" json:"code"`
	Capacity   int       `gorm:"default:1000" json:"capacity"`
	Workload   int       `gorm:"default:0" json:"workload"`
	Active     bool      `gorm:"default:true" json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductionSchedule struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	OrderID    uint           `gorm:"not null;index" json:"order_id"`
	Order      Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	LineID     uint           `gorm:"not null;index" json:"line_id"`
	Line       ProductionLine `gorm:"foreignKey:LineID" json:"line,omitempty"`
	PlannedQty int            `gorm:"not null" json:"planned_qty"`
	ProducedQty int           `gorm:"default:0" json:"produced_qty"`
	StartDate  time.Time      `json:"start_date"`
	EndDate    time.Time      `json:"end_date"`
	Status     string         `gorm:"size:32;default:scheduled" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type FileAsset struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FileName    string    `gorm:"size:256;not null" json:"file_name"`
	FilePath    string    `gorm:"size:512;not null" json:"file_path"`
	FileSize    int64     `json:"file_size"`
	FileType    string    `gorm:"size:64" json:"file_type"`
	FileHash    string    `gorm:"size:128" json:"file_hash"`
	UploaderID  uint      `json:"uploader_id"`
	Status      string    `gorm:"size:32;default:uploaded" json:"status"`
	PreviewURL  string    `gorm:"size:512" json:"preview_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type Invoice struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	InvoiceNo    string        `gorm:"uniqueIndex;size:32;not null" json:"invoice_no"`
	CustomerID   uint          `gorm:"not null;index" json:"customer_id"`
	Customer     Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	PeriodStart  time.Time     `json:"period_start"`
	PeriodEnd    time.Time     `json:"period_end"`
	TotalAmount  float64       `gorm:"default:0" json:"total_amount"`
	PaidAmount   float64       `gorm:"default:0" json:"paid_amount"`
	Status       string        `gorm:"size:32;default:unpaid" json:"status"`
	Items        []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type InvoiceItem struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	InvoiceID  uint    `gorm:"not null;index" json:"invoice_id"`
	OrderID    uint    `json:"order_id"`
	Amount     float64 `gorm:"not null" json:"amount"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Action    string    `gorm:"size:64;index" json:"action"`
	Resource  string    `gorm:"size:64" json:"resource"`
	ResourceID uint     `json:"resource_id"`
	Detail    string    `gorm:"size:2048" json:"detail"`
	IP        string    `gorm:"size:64" json:"ip"`
	UserAgent string    `gorm:"size:512" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
