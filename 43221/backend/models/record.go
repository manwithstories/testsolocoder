package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConsultRecord struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	AppointmentID   uuid.UUID      `json:"appointment_id" gorm:"type:uuid;not null;uniqueIndex"`
	ClientID        uuid.UUID      `json:"client_id" gorm:"type:uuid;not null;index"`
	ProfessionalID  uuid.UUID      `json:"professional_id" gorm:"type:uuid;not null;index"`
	Summary         string         `json:"summary" gorm:"type:text"`
	Advice          string         `json:"advice" gorm:"type:text"`
	FollowUpDate    *time.Time     `json:"follow_up_date"`
	IsConfidential  bool           `json:"is_confidential" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Appointment     Appointment    `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	Client          User           `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	Professional    User           `json:"professional,omitempty" gorm:"foreignKey:ProfessionalID"`
}

type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "pending"
	ReviewApproved ReviewStatus = "approved"
	ReviewRejected ReviewStatus = "rejected"
)

type Review struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	AppointmentID   uuid.UUID      `json:"appointment_id" gorm:"type:uuid;not null;uniqueIndex"`
	ClientID        uuid.UUID      `json:"client_id" gorm:"type:uuid;not null;index"`
	ProfessionalID  uuid.UUID      `json:"professional_id" gorm:"type:uuid;not null;index"`
	ServiceID       uuid.UUID      `json:"service_id" gorm:"type:uuid;not null;index"`
	Rating          int            `json:"rating" gorm:"not null"`
	Content         string         `json:"content" gorm:"type:text"`
	Status          ReviewStatus   `json:"status" gorm:"type:varchar(20);default:pending"`
	RejectReason    string         `json:"reject_reason"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Appointment     Appointment    `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	Client          User           `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	Professional    User           `json:"professional,omitempty" gorm:"foreignKey:ProfessionalID"`
	Service         Service        `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

func (c *ConsultRecord) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
