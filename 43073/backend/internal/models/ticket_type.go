package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TicketTypeNormal  = "normal"
	TicketTypeVIP     = "vip"
	TicketTypeEarly   = "early_bird"
	TicketStatusOnSale  = "on_sale"
	TicketStatusSoldOut = "sold_out"
	TicketStatusOffSale = "off_sale"
)

type TicketType struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ActivityID uint           `json:"activityId" validate:"required"`
	Name       string         `gorm:"size:100;not null" json:"name" validate:"required,max=100"`
	Type       string         `gorm:"size:20;not null" json:"type" validate:"required,oneof=normal vip early_bird"`
	Price      float64        `json:"price" validate:"required,min=0"`
	Stock      int            `json:"stock" validate:"required,min=0"`
	SoldCount  int            `json:"soldCount"`
	Status     string         `gorm:"size:20;default:'on_sale'" json:"status"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Activity   *Activity      `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
}
