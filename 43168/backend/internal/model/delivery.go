package model

import (
	"time"
)

// 配送安装类型常量
const (
	DeliveryTypeDelivery = "delivery" // 仅配送
	DeliveryTypeInstall  = "install"  // 仅安装
	DeliveryTypeBoth     = "both"     // 配送并安装
)

// ValidDeliveryType 校验配送安装类型是否合法
func ValidDeliveryType(t string) bool {
	switch t {
	case DeliveryTypeDelivery, DeliveryTypeInstall, DeliveryTypeBoth:
		return true
	default:
		return false
	}
}

// 配送安装状态常量
const (
	DeliveryStatusPending   = "pending"   // 待确认
	DeliveryStatusConfirmed = "confirmed" // 已确认
	DeliveryStatusCompleted = "completed" // 已完成
	DeliveryStatusCancelled = "cancelled" // 已取消
)

// ValidDeliveryStatus 校验状态是否合法
func ValidDeliveryStatus(s string) bool {
	switch s {
	case DeliveryStatusPending, DeliveryStatusConfirmed,
		DeliveryStatusCompleted, DeliveryStatusCancelled:
		return true
	default:
		return false
	}
}

// Delivery 配送安装模型
type Delivery struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"index;not null" json:"order_id"`
	OwnerID      uint      `gorm:"index;not null" json:"owner_id"`
	Type         string    `gorm:"size:32;not null;index" json:"type"`
	TimeSlot     string    `gorm:"size:128" json:"time_slot"`
	Address      string    `gorm:"size:255" json:"address"`
	ContactName  string    `gorm:"size:64" json:"contact_name"`
	ContactPhone string    `gorm:"size:20" json:"contact_phone"`
	Status       string    `gorm:"size:32;not null;default:pending;index" json:"status"`
	InstallerID  uint      `gorm:"index" json:"installer_id"`
	Remark       string    `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Delivery) TableName() string {
	return "deliveries"
}
