package model

import (
	"time"

	"gorm.io/gorm"
)

type PointsLogType string

const (
	PointsLogTypeEarn    PointsLogType = "earn"
	PointsLogTypeUse     PointsLogType = "use"
	PointsLogTypeRecharge PointsLogType = "recharge"
	PointsLogTypeRefund  PointsLogType = "refund"
)

type MemberPointsLog struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	MemberID    uint           `gorm:"not null;index" json:"member_id"`
	Points      int            `gorm:"not null" json:"points"`
	Balance     int            `gorm:"not null" json:"balance"`
	Type        PointsLogType  `gorm:"size:20;not null;index" json:"type"`
	Description string         `gorm:"size:255" json:"description"`
	OrderNo     string         `gorm:"size:64;index" json:"order_no"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Member      *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}
