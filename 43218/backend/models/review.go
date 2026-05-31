package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderID       *uint          `gorm:"index" json:"order_id"`
	RepairOrderID *uint          `gorm:"index" json:"repair_order_id"`
	ReviewerID    uint           `gorm:"index;not null" json:"reviewer_id"`
	RevieweeID    uint           `gorm:"index;not null" json:"reviewee_id"`
	ReviewType    string         `gorm:"size:20;not null" json:"review_type"`
	Rating        int            `gorm:"not null" json:"rating" binding:"required,min=1,max=5"`
	Content       string         `gorm:"type:text" json:"content" binding:"max=500"`
	Images        string         `gorm:"type:text" json:"images"`
	QualityScore  *int           `json:"quality_score"`
	ServiceScore  *int           `json:"service_score"`
	Anonymous     bool           `gorm:"default:false" json:"anonymous"`
	Status        int            `gorm:"default:1" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Reviewer User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Reviewee User `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}

const (
	ReviewTypeProduct     = "product"
	ReviewTypeSeller      = "seller"
	ReviewTypeRepair      = "repair"
	ReviewTypeTechnician  = "technician"
)

type Report struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ReporterID    uint           `gorm:"index;not null" json:"reporter_id"`
	TargetType    string         `gorm:"size:20;not null" json:"target_type"`
	TargetID      uint           `gorm:"index;not null" json:"target_id"`
	Reason        string         `gorm:"size:100;not null" json:"reason"`
	Description   string         `gorm:"type:text" json:"description"`
	Images        string         `gorm:"type:text" json:"images"`
	Status        int            `gorm:"default:0" json:"status"`
	HandleResult  string         `gorm:"type:text" json:"handle_result"`
	HandledBy     *uint          `json:"handled_by"`
	HandledAt     *time.Time     `json:"handled_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Reporter User `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
}

func (Report) TableName() string {
	return "reports"
}

const (
	ReportStatusPending = 0
	ReportStatusApproved = 1
	ReportStatusRejected = 2
)

const (
	ReportTargetProduct   = "product"
	ReportTargetUser      = "user"
	ReportTargetService   = "service"
	ReportTargetReview    = "review"
)

var ReportReasons = []string{
	"虚假宣传",
	"商品质量问题",
	"欺诈行为",
	"恶意差评",
	"违规发布",
	"骚扰行为",
	"其他违规",
}

type Warranty struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderID       *uint          `gorm:"index" json:"order_id"`
	RepairOrderID *uint          `gorm:"index" json:"repair_order_id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	Type          string         `gorm:"size:20;not null" json:"type"`
	Description   string         `gorm:"type:text" json:"description"`
	Images        string         `gorm:"type:text" json:"images"`
	Status        int            `gorm:"default:0" json:"status"`
	HandleResult  string         `gorm:"type:text" json:"handle_result"`
	HandledBy     *uint          `json:"handled_by"`
	HandledAt     *time.Time     `json:"handled_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Warranty) TableName() string {
	return "warranties"
}

const (
	WarrantyStatusPending  = 0
	WarrantyStatusApproved = 1
	WarrantyStatusRejected = 2
	WarrantyStatusProcessing = 3
	WarrantyStatusCompleted = 4
)

const (
	WarrantyTypeProduct = "product"
	WarrantyTypeRepair  = "repair"
)
