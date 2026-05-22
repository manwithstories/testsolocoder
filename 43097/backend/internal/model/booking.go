package model

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	BookingNo      string         `gorm:"size:32;not null;uniqueIndex" json:"booking_no"`
	RoomID         uint           `gorm:"not null;index" json:"room_id"`
	MemberID       *uint          `gorm:"index" json:"member_id"`
	GuestName      string         `gorm:"size:50;not null" json:"guest_name"`
	GuestPhone     string         `gorm:"size:20;not null;index" json:"guest_phone"`
	GuestIDCard    string         `gorm:"size:20;index" json:"guest_id_card"`
	CheckInDate    time.Time      `gorm:"not null;index" json:"check_in_date"`
	CheckOutDate   time.Time      `gorm:"not null;index" json:"check_out_date"`
	Days           int            `gorm:"not null" json:"days"`
	TotalPrice     float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status         BookingStatus  `gorm:"size:20;not null;default:'pending'" json:"status"`
	PaidAmount     float64        `gorm:"type:decimal(10,2);not null;default:0" json:"paid_amount"`
	Remarks        string         `gorm:"type:text" json:"remarks"`
	CancelDeadline *time.Time     `json:"cancel_deadline"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Room           *Room          `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Member         *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	CheckIns       []CheckIn      `gorm:"foreignKey:BookingID" json:"check_ins,omitempty"`
}
