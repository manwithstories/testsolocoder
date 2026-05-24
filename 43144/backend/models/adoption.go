package models

import (
	"time"

	"gorm.io/gorm"
)

type AdoptionStatus string

const (
	AdoptionStatusPending   AdoptionStatus = "pending"
	AdoptionStatusApproved  AdoptionStatus = "approved"
	AdoptionStatusRejected  AdoptionStatus = "rejected"
	AdoptionStatusSigned    AdoptionStatus = "signed"
	AdoptionStatusCompleted AdoptionStatus = "completed"
	AdoptionStatusCancelled AdoptionStatus = "cancelled"
)

type AdoptionApplication struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	PetID         uint           `json:"pet_id" gorm:"index;not null"`
	Pet           *Pet           `json:"pet,omitempty" gorm:"foreignKey:PetID"`
	AdopterID     uint           `json:"adopter_id" gorm:"index;not null"`
	Adopter       *User          `json:"adopter,omitempty" gorm:"foreignKey:AdopterID"`
	RescueID      uint           `json:"rescue_id" gorm:"index;not null"`
	Status        AdoptionStatus `json:"status" gorm:"size:20;default:pending;index;not null"`
	Reason        string         `json:"reason" gorm:"type:text"`
	LivingSituation string       `json:"living_situation" gorm:"type:text"`
	PetExperience string         `json:"pet_experience" gorm:"type:text"`
	FamilyMembers int            `json:"family_members"`
	HasChildren   bool           `json:"has_children" gorm:"default:false"`
	HasOtherPets  bool           `json:"has_other_pets" gorm:"default:false"`
	OtherPetsDesc string         `json:"other_pets_desc" gorm:"type:text"`
	HousingType   string         `json:"housing_type" gorm:"size:50"`
	IncomeLevel   string         `json:"income_level" gorm:"size:50"`
	CanAffordVet  bool           `json:"can_afford_vet" gorm:"default:true"`
	AgreeToVisit  bool           `json:"agree_to_visit" gorm:"default:true"`
	ReviewedBy    *uint          `json:"reviewed_by,omitempty" gorm:"index"`
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	RejectReason  string         `json:"reject_reason" gorm:"type:text"`
	SignedAt      *time.Time     `json:"signed_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type AdoptionAgreement struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ApplicationID   uint           `json:"application_id" gorm:"uniqueIndex;not null"`
	Application     *AdoptionApplication `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
	AdopterSign     bool           `json:"adopter_sign" gorm:"default:false"`
	AdopterSignedAt *time.Time     `json:"adopter_signed_at"`
	RescueSign      bool           `json:"rescue_sign" gorm:"default:false"`
	RescueSignedAt  *time.Time     `json:"rescue_signed_at"`
	AgreementTerms  string         `json:"agreement_terms" gorm:"type:text;not null"`
	AgreementFile   string         `json:"agreement_file" gorm:"size:500"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type FollowUpRecord struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ApplicationID  uint           `json:"application_id" gorm:"index;not null"`
	Application    *AdoptionApplication `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
	PetID          uint           `json:"pet_id" gorm:"index;not null"`
	Pet            *Pet           `json:"pet,omitempty" gorm:"foreignKey:PetID"`
	AdopterID      uint           `json:"adopter_id" gorm:"index;not null"`
	RescueID       uint           `json:"rescue_id" gorm:"index;not null"`
	FollowUpDate   *time.Time     `json:"follow_up_date"`
	HealthStatus   string         `json:"health_status" gorm:"type:text"`
	LivingCondition string        `json:"living_condition" gorm:"type:text"`
	Notes          string         `json:"notes" gorm:"type:text"`
	PhotoEvidence  string         `json:"photo_evidence" gorm:"type:text"`
	RecordedBy     uint           `json:"recorded_by" gorm:"index;not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateAdoptionRequest struct {
	PetID            uint   `json:"pet_id" binding:"required"`
	Reason           string `json:"reason" binding:"required"`
	LivingSituation  string `json:"living_situation" binding:"required"`
	PetExperience    string `json:"pet_experience"`
	FamilyMembers    int    `json:"family_members"`
	HasChildren      bool   `json:"has_children"`
	HasOtherPets     bool   `json:"has_other_pets"`
	OtherPetsDesc    string `json:"other_pets_desc"`
	HousingType      string `json:"housing_type" binding:"required"`
	IncomeLevel      string `json:"income_level"`
	CanAffordVet     bool   `json:"can_afford_vet"`
	AgreeToVisit     bool   `json:"agree_to_visit" binding:"required"`
}

type ReviewAdoptionRequest struct {
	Action       string `json:"action" binding:"required,oneof=approve reject"`
	RejectReason string `json:"reject_reason"`
}

type CreateFollowUpRequest struct {
	ApplicationID  uint   `json:"application_id" binding:"required"`
	FollowUpDate   string `json:"follow_up_date" binding:"required"`
	HealthStatus   string `json:"health_status"`
	LivingCondition string `json:"living_condition"`
	Notes          string `json:"notes"`
}

type AdoptionListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status"`
	PetID    uint   `form:"pet_id"`
	RescueID uint   `form:"rescue_id"`
}
