package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus int

const (
	BookingStatusPending   BookingStatus = 0
	BookingStatusConfirmed BookingStatus = 1
	BookingStatusCancelled BookingStatus = 2
	BookingStatusCompleted BookingStatus = 3
)

type RecurrenceType string

const (
	RecurrenceNone     RecurrenceType = "none"
	RecurrenceDaily    RecurrenceType = "daily"
	RecurrenceWeekly   RecurrenceType = "weekly"
	RecurrenceBiweekly RecurrenceType = "biweekly"
	RecurrenceMonthly  RecurrenceType = "monthly"
)

type Booking struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	RoomID          uint           `json:"room_id" gorm:"index;not null"`
	UserID          uint           `json:"user_id" gorm:"index;not null"`
	Title           string         `json:"title" gorm:"size:200;not null"`
	Description     string         `json:"description" gorm:"type:text"`
	StartTime       time.Time      `json:"start_time" gorm:"index;not null"`
	EndTime         time.Time      `json:"end_time" gorm:"index;not null"`
	RecurrenceType  RecurrenceType `json:"recurrence_type" gorm:"size:20;default:none"`
	RecurrenceEnd   *time.Time     `json:"recurrence_end"`
	ParentBookingID *uint          `json:"parent_booking_id" gorm:"index"`
	Status          BookingStatus  `json:"status" gorm:"default:0;index"`
	Attendees       string         `json:"attendees" gorm:"type:text"`
	TotalPrice      float64        `json:"total_price" gorm:"type:decimal(10,2)"`
	Reminded        bool           `json:"reminded" gorm:"default:false"`
	CancelledBy     uint           `json:"cancelled_by"`
	CancelledAt     *time.Time     `json:"cancelled_at"`
	CancelReason    string         `json:"cancel_reason" gorm:"size:500"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	Room            *Room          `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	User            *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Booking) TableName() string {
	return "bookings"
}

func (b *Booking) CalculateTotalPrice(pricePerHour float64) float64 {
	duration := b.EndTime.Sub(b.StartTime).Hours()
	return float64(int(duration*100)) / 100 * pricePerHour
}
