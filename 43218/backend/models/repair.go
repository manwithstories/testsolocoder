package models

import (
	"time"

	"gorm.io/gorm"
)

type RepairService struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TechnicianID  uint           `gorm:"index;not null" json:"technician_id"`
	ServiceType   string         `gorm:"size:50;not null;index" json:"service_type" binding:"required"`
	Title         string         `gorm:"size:200;not null" json:"title" binding:"required,min=2,max=200"`
	Description   string         `gorm:"type:text" json:"description" binding:"required"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price" binding:"required,min=0"`
	MinPrice      float64        `gorm:"type:decimal(10,2)" json:"min_price"`
	MaxPrice      float64        `gorm:"type:decimal(10,2)" json:"max_price"`
	EstimatedDays int            `gorm:"default:3" json:"estimated_days"`
	Images        string         `gorm:"type:text" json:"images"`
	Status        int            `gorm:"default:1;index" json:"status"`
	OrderCount    int            `gorm:"default:0" json:"order_count"`
	Rating        float64        `gorm:"default:0" json:"rating"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Technician User `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
}

func (RepairService) TableName() string {
	return "repair_services"
}

var ServiceTypes = []string{"屏幕更换", "电池更换", "主板维修", "外壳更换", "摄像头维修", "充电接口维修", "扬声器维修", "系统维护", "其他"}

const (
	ServiceStatusActive   = 1
	ServiceStatusInactive = 2
)

type RepairOrder struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrderNo        string         `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	BuyerID        uint           `gorm:"index;not null" json:"buyer_id"`
	TechnicianID   uint           `gorm:"index;not null" json:"technician_id"`
	ServiceID      uint           `gorm:"index;not null" json:"service_id"`
	DeviceType     string         `gorm:"size:50;not null" json:"device_type"`
	DeviceBrand    string         `gorm:"size:100;not null" json:"device_brand"`
	DeviceModel    string         `gorm:"size:100;not null" json:"device_model"`
	FaultDescription string       `gorm:"type:text;not null" json:"fault_description"`
	ContactName    string         `gorm:"size:50;not null" json:"contact_name"`
	ContactPhone   string         `gorm:"size:20;not null" json:"contact_phone"`
	Address        string         `gorm:"size:255" json:"address"`
	ServicePrice   float64        `gorm:"type:decimal(10,2);not null" json:"service_price"`
	FinalPrice     float64        `gorm:"type:decimal(10,2)" json:"final_price"`
	Status         int            `gorm:"default:1;index" json:"status"`
	PaymentMethod  string         `gorm:"size:20" json:"payment_method"`
	PaidAt         *time.Time     `json:"paid_at"`
	AcceptedAt     *time.Time     `json:"accepted_at"`
	CompletedAt    *time.Time     `json:"completed_at"`
	PickedUpAt     *time.Time     `json:"picked_up_at"`
	WarrantyDays   int            `gorm:"default:90" json:"warranty_days"`
	WarrantyUntil  *time.Time     `json:"warranty_until"`
	Remark         string         `gorm:"size:255" json:"remark"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Buyer      User          `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Technician User          `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	Service    RepairService `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	Reviews    []Review      `gorm:"foreignKey:RepairOrderID" json:"reviews,omitempty"`
}

func (RepairOrder) TableName() string {
	return "repair_orders"
}

const (
	RepairStatusPending    = 1
	RepairStatusAccepted   = 2
	RepairStatusRepairing  = 3
	RepairStatusCompleted  = 4
	RepairStatusPickedUp   = 5
	RepairStatusCancelled  = 6
	RepairStatusRefunding  = 7
	RepairStatusRefunded   = 8
)

var RepairStatusText = map[int]string{
	1: "待接单",
	2: "已接单",
	3: "维修中",
	4: "待取件",
	5: "已取件",
	6: "已取消",
	7: "退款中",
	8: "已退款",
}
