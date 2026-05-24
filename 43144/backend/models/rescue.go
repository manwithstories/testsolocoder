package models

import (
	"time"

	"gorm.io/gorm"
)

type RescueStatus string

const (
	RescueStatusPending  RescueStatus = "pending"
	RescueStatusApproved RescueStatus = "approved"
	RescueStatusRejected RescueStatus = "rejected"
)

type RescueStation struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"size:255;not null"`
	Address       string         `json:"address" gorm:"size:500"`
	ContactPerson string         `json:"contact_person" gorm:"size:100"`
	ContactPhone  string         `json:"contact_phone" gorm:"size:20"`
	ContactEmail  string         `json:"contact_email" gorm:"size:255"`
	LicenseNumber string         `json:"license_number" gorm:"size:100"`
	LicenseFile   string         `json:"license_file" gorm:"size:500"`
	Description   string         `json:"description" gorm:"type:text"`
	Status        RescueStatus   `json:"status" gorm:"size:20;default:pending;index;not null"`
	VerifiedBy    *uint          `json:"verified_by,omitempty" gorm:"index"`
	VerifiedAt    *time.Time     `json:"verified_at"`
	RejectReason  string         `json:"reject_reason" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type ReviewRescueRequest struct {
	Action       string `json:"action" binding:"required,oneof=approve reject"`
	RejectReason string `json:"reject_reason"`
}

type RescueStats struct {
	TotalPets          int64 `json:"total_pets"`
	AdoptablePets      int64 `json:"adoptable_pets"`
	AdoptedPets        int64 `json:"adopted_pets"`
	TreatmentPets      int64 `json:"treatment_pets"`
	DeceasedPets       int64 `json:"deceased_pets"`
	TotalAdoptions     int64 `json:"total_adoptions"`
	PendingApplications int64 `json:"pending_applications"`
	CompletedFollowUps int64 `json:"completed_follow_ups"`
	TotalFollowUps     int64 `json:"total_follow_ups"`
	AdoptionRate       float64 `json:"adoption_rate"`
	FollowUpRate       float64 `json:"follow_up_rate"`
}

type RescueListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status"`
	Search   string `form:"search"`
}
