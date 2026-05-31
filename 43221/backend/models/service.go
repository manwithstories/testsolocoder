package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceType string

const (
	ServiceTypeLegal       ServiceType = "legal"
	ServiceTypeCounseling  ServiceType = "counseling"
	ServiceTypeFinancial   ServiceType = "financial"
	ServiceTypeOther       ServiceType = "other"
)

type ServiceStatus string

const (
	ServiceStatusActive   ServiceStatus = "active"
	ServiceStatusInactive ServiceStatus = "inactive"
)

type Service struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	ProfessionalID  uuid.UUID      `json:"professional_id" gorm:"type:uuid;not null;index"`
	Title           string         `json:"title" gorm:"not null;size:200"`
	Description     string         `json:"description" gorm:"type:text"`
	ServiceType     ServiceType    `json:"service_type" gorm:"type:varchar(30);not null"`
	Price           float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	DurationMinutes int            `json:"duration_minutes" gorm:"not null;default:60"`
	Status          ServiceStatus  `json:"status" gorm:"type:varchar(20);default:active"`
	AverageRating   float64        `json:"average_rating" gorm:"type:decimal(3,2);default:0"`
	ReviewCount     int            `json:"review_count" gorm:"default:0"`
	Tags            string         `json:"tags" gorm:"size:500"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Professional    User           `json:"professional,omitempty" gorm:"foreignKey:ProfessionalID"`
	Schedules       []Schedule     `json:"schedules,omitempty" gorm:"foreignKey:ServiceID"`
}

type Schedule struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	ServiceID    uuid.UUID      `json:"service_id" gorm:"type:uuid;not null;index"`
	Date         time.Time      `json:"date" gorm:"type:date;not null"`
	StartTime    string         `json:"start_time" gorm:"type:varchar(5);not null"`
	EndTime      string         `json:"end_time" gorm:"type:varchar(5);not null"`
	IsBooked     bool           `json:"is_booked" gorm:"default:false"`
	IsAvailable  bool           `json:"is_available" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	Appointment  *Appointment   `json:"appointment,omitempty" gorm:"foreignKey:ScheduleID"`
}

func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
