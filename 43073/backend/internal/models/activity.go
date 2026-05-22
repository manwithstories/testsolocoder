package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ActivityStatusDraft     = "draft"
	ActivityStatusPublished = "published"
	ActivityStatusCanceled  = "canceled"
)

type Activity struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title" validate:"required,max=200"`
	Description string         `gorm:"type:text" json:"description"`
	StartTime   time.Time      `json:"startTime" validate:"required"`
	EndTime     time.Time      `json:"endTime" validate:"required,gtfield=StartTime"`
	Location    string         `gorm:"size:500" json:"location" validate:"required,max=500"`
	Capacity    int            `json:"capacity" validate:"required,min=1"`
	Poster      string         `gorm:"size:500" json:"poster"`
	Status      string         `gorm:"size:20;default:'draft'" json:"status"`
	CreatedBy   uint           `json:"createdBy"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	TicketTypes []TicketType   `gorm:"foreignKey:ActivityID" json:"ticketTypes,omitempty"`
}
