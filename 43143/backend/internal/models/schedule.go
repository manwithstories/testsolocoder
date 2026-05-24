package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
	Saturday  DayOfWeek = "saturday"
	Sunday    DayOfWeek = "sunday"
)

type ScheduleType string

const (
	ScheduleTypeAvailability ScheduleType = "availability"
	ScheduleTypeBusy         ScheduleType = "busy"
	ScheduleTypeBooking      ScheduleType = "booking"
)

type Schedule struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        ScheduleType   `gorm:"type:varchar(20);default:'availability'" json:"type"`
	DayOfWeek   DayOfWeek      `gorm:"type:varchar(20)" json:"day_of_week"`
	SpecificDate *time.Time   `json:"specific_date"`
	StartTime   string         `gorm:"size:10" json:"start_time"`
	EndTime     string         `gorm:"size:10" json:"end_time"`
	IsRecurring bool           `gorm:"default:true" json:"is_recurring"`
	BookingID   *uuid.UUID     `gorm:"type:uuid;index" json:"booking_id"`
	Booking     *Booking       `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Title       string         `gorm:"size:200" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
