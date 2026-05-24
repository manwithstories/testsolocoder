package models

import (
	"time"

	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	ApplicationStatusPending   ApplicationStatus = "pending"
	ApplicationStatusReviewed  ApplicationStatus = "reviewed"
	ApplicationStatusInterview ApplicationStatus = "interview"
	ApplicationStatusAccepted  ApplicationStatus = "accepted"
	ApplicationStatusRejected  ApplicationStatus = "rejected"
	ApplicationStatusHold      ApplicationStatus = "hold"
)

type Application struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	JobID       uint              `gorm:"not null;index" json:"job_id"`
	JobSeekerID uint              `gorm:"not null;index" json:"jobseeker_id"`
	ResumeID    uint              `gorm:"not null" json:"resume_id"`
	Status      ApplicationStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	CoverLetter string            `gorm:"type:text" json:"cover_letter"`
	Job         Job               `gorm:"foreignKey:JobID" json:"job,omitempty"`
	JobSeeker   JobSeeker         `gorm:"foreignKey:JobSeekerID" json:"jobseeker,omitempty"`
	Resume      Resume            `gorm:"foreignKey:ResumeID" json:"resume,omitempty"`
	Interviews  []Interview       `gorm:"foreignKey:ApplicationID" json:"interviews,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}
