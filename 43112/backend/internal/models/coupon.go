package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CouponType string

const (
	CouponFixed   CouponType = "fixed"
	CouponPercent CouponType = "percent"
)

type Coupon struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Type        CouponType     `gorm:"size:20;not null" json:"type"`
	Value       float64        `gorm:"not null" json:"value"`
	MinAmount   float64        `gorm:"default:0" json:"min_amount"`
	MaxDiscount float64        `gorm:"default:0" json:"max_discount"`
	TotalCount  int            `gorm:"default:0" json:"total_count"`
	UsedCount   int            `gorm:"default:0" json:"used_count"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CouponUsed struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CouponID  uuid.UUID `gorm:"type:uuid;index;not null" json:"coupon_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	OrderID   uuid.UUID `gorm:"type:uuid;index;not null" json:"order_id"`
	UsedAt    time.Time `json:"used_at"`
}

func (c *Coupon) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (c *CouponUsed) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
