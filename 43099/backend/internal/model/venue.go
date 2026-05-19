package model

import (
	"time"
)

type VenueStatus string

const (
	VenueStatusOnline  VenueStatus = "online"
	VenueStatusOffline VenueStatus = "offline"
)

type Venue struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `gorm:"size:100;not null" json:"name"`
	Location    string      `gorm:"size:255" json:"location"`
	Capacity    int         `gorm:"default:0" json:"capacity"`
	Facilities  string      `gorm:"type:text" json:"facilities"`
	Description string      `gorm:"type:text" json:"description"`
	CoverImage  string      `gorm:"size:255" json:"cover_image"`
	Status      VenueStatus `gorm:"size:20;default:'online'" json:"status"`
	CreatedBy   uint        `json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type TimeSlot struct {
	Start string  `json:"start"`
	End   string  `json:"end"`
	Price float64 `json:"price"`
}

type VenuePrice struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	VenueID   uint       `gorm:"index;not null" json:"venue_id"`
	DayOfWeek int        `gorm:"not null" json:"day_of_week"`
	TimeSlots string     `gorm:"type:json" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
