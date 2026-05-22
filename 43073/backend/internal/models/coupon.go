package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	CouponTypeFixed    = "fixed"
	CouponTypeDiscount = "discount"
	CouponStatusActive = "active"
	CouponStatusUsed   = "used"
	CouponStatusExpired = "expired"
)

type Coupon struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Code       string         `gorm:"size:50;uniqueIndex;not null" json:"code" validate:"required,max=50"`
	Type       string         `gorm:"size:20;not null" json:"type" validate:"required,oneof=fixed discount"`
	Value      float64        `json:"value" validate:"required,min=0"`
	MinAmount  float64        `json:"minAmount" validate:"min=0"`
	TotalCount int            `json:"totalCount" validate:"min=1"`
	UsedCount  int            `json:"usedCount"`
	StartTime  time.Time      `json:"startTime"`
	EndTime    time.Time      `json:"endTime"`
	Status     string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
