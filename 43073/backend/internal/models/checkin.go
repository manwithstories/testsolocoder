package models

import (
	"time"

	"gorm.io/gorm"
)

type CheckIn struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderID      uint           `json:"orderId"`
	OrderItemID  uint           `json:"orderItemId"`
	UserID       uint           `json:"userId"`
	ActivityID   uint           `json:"activityId"`
	TicketTypeID uint           `json:"ticketTypeId"`
	QrCode       string         `gorm:"size:100;uniqueIndex;not null" json:"qrCode"`
	CheckedIn    bool           `json:"checkedIn"`
	CheckedAt    *time.Time     `json:"checkedAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Order        *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Activity     *Activity      `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
}
