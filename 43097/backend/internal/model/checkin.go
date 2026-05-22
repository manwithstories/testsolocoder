package model

import (
	"time"

	"gorm.io/gorm"
)

type CheckInStatus string

const (
	CheckInStatusActive    CheckInStatus = "active"
	CheckInStatusCheckedOut CheckInStatus = "checked_out"
)

type CheckIn struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CheckInNo        string         `gorm:"size:32;not null;uniqueIndex" json:"check_in_no"`
	BookingID        *uint          `gorm:"index" json:"booking_id"`
	RoomID           uint           `gorm:"not null;index" json:"room_id"`
	GuestName        string         `gorm:"size:50;not null" json:"guest_name"`
	GuestPhone       string         `gorm:"size:20;not null;index" json:"guest_phone"`
	GuestIDCard      string         `gorm:"size:20;index" json:"guest_id_card"`
	CheckInTime      time.Time      `gorm:"not null;index" json:"check_in_time"`
	ExpectedCheckOut time.Time      `gorm:"not null" json:"expected_check_out"`
	ActualCheckOut   *time.Time     `json:"actual_check_out"`
	Status           CheckInStatus  `gorm:"size:20;not null;default:'active'" json:"status"`
	Deposit          float64        `gorm:"type:decimal(10,2);not null;default:0" json:"deposit"`
	ExtraCharges     float64        `gorm:"type:decimal(10,2);not null;default:0" json:"extra_charges"`
	TotalAmount      float64        `gorm:"type:decimal(10,2);not null;default:0" json:"total_amount"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	Booking          *Booking       `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Room             *Room          `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}
