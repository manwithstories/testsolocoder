package model

import (
	"time"

	"gorm.io/gorm"
)

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusOccupied    RoomStatus = "occupied"
	RoomStatusReserved    RoomStatus = "reserved"
	RoomStatusMaintenance RoomStatus = "maintenance"
)

type Room struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	RoomNo     string         `gorm:"size:20;not null;uniqueIndex" json:"room_no"`
	Floor      int            `gorm:"not null;index" json:"floor"`
	RoomTypeID uint           `gorm:"not null;index" json:"room_type_id"`
	Status     RoomStatus     `gorm:"size:20;not null;default:'available'" json:"status"`
	Price      float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Facilities StringArray    `gorm:"type:jsonb" json:"facilities"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	RoomType   *RoomType      `gorm:"foreignKey:RoomTypeID" json:"room_type,omitempty"`
}
