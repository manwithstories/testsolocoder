package models

import (
	"time"

	"gorm.io/gorm"
)

type AppointmentType string

const (
	AppointmentVisit   AppointmentType = "visit"
	AppointmentCheckup AppointmentType = "checkup"
)

type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "pending"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

type Appointment struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	UserID        uint              `json:"user_id" gorm:"index;not null"`
	User          *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PetID         uint              `json:"pet_id" gorm:"index;not null"`
	Pet           *Pet              `json:"pet,omitempty" gorm:"foreignKey:PetID"`
	RescueID      uint              `json:"rescue_id" gorm:"index;not null"`
	AppointmentType AppointmentType  `json:"appointment_type" gorm:"size:20;not null"`
	AppointmentDate time.Time       `json:"appointment_date"`
	StartTime     string            `json:"start_time" gorm:"size:10;not null"`
	EndTime       string            `json:"end_time" gorm:"size:10;not null"`
	Status        AppointmentStatus `json:"status" gorm:"size:20;default:pending;index;not null"`
	Location      string            `json:"location" gorm:"size:255"`
	Notes         string            `json:"notes" gorm:"type:text"`
	CancelReason  string            `json:"cancel_reason" gorm:"type:text"`
	OriginalID    *uint             `json:"original_id,omitempty" gorm:"index"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
}

type CreateAppointmentRequest struct {
	PetID           uint   `json:"pet_id" binding:"required"`
	AppointmentType string `json:"appointment_type" binding:"required,oneof=visit checkup"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	StartTime       string `json:"start_time" binding:"required"`
	EndTime         string `json:"end_time" binding:"required"`
	Location        string `json:"location"`
	Notes           string `json:"notes"`
}

type UpdateAppointmentRequest struct {
	AppointmentDate string `json:"appointment_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	Location        string `json:"location"`
	Notes           string `json:"notes"`
}

type AppointmentListQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	Status     string `form:"status"`
	PetID      uint   `form:"pet_id"`
	RescueID   uint   `form:"rescue_id"`
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
}
