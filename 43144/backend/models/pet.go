package models

import (
	"time"

	"gorm.io/gorm"
)

type PetStatus string

const (
	PetStatusAdoptable PetStatus = "adoptable"
	PetStatusAdopted   PetStatus = "adopted"
	PetStatusTreatment PetStatus = "treatment"
	PetStatusDeceased  PetStatus = "deceased"
)

type PetGender string

const (
	GenderMale   PetGender = "male"
	GenderFemale PetGender = "female"
	GenderUnknown PetGender = "unknown"
)

type Pet struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ArchiveNumber   string         `json:"archive_number" gorm:"uniqueIndex;size:50;not null"`
	Name            string         `json:"name" gorm:"size:100;not null"`
	Species         string         `json:"species" gorm:"size:50;not null"`
	Breed           string         `json:"breed" gorm:"size:100"`
	Age             string         `json:"age" gorm:"size:50"`
	Gender          PetGender      `json:"gender" gorm:"size:20"`
	Weight          float64        `json:"weight"`
	Color           string         `json:"color" gorm:"size:50"`
	Description     string         `json:"description" gorm:"type:text"`
	Status          PetStatus      `json:"status" gorm:"size:20;default:adoptable;index;not null"`
	Photos          string         `json:"photos" gorm:"type:text"`
	Videos          string         `json:"videos" gorm:"type:text"`
	HealthStatus    string         `json:"health_status" gorm:"size:255"`
	Vaccinated      bool           `json:"vaccinated" gorm:"default:false"`
	Neutered        bool           `json:"neutered" gorm:"default:false"`
	RescueID        uint           `json:"rescue_id" gorm:"index;not null"`
	Rescue          *RescueStation `json:"rescue,omitempty" gorm:"foreignKey:RescueID"`
	AdopterID       *uint          `json:"adopter_id,omitempty" gorm:"index"`
	Adopter         *User          `json:"adopter,omitempty" gorm:"foreignKey:AdopterID"`
	FoundLocation   string         `json:"found_location" gorm:"size:255"`
	FoundDate       *time.Time     `json:"found_date"`
	AdoptedDate     *time.Time     `json:"adopted_date"`
	MedicalHistory  string         `json:"medical_history" gorm:"type:text"`
	Personality     string         `json:"personality" gorm:"type:text"`
	SpecialNeeds    string         `json:"special_needs" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreatePetRequest struct {
	Name           string    `json:"name" binding:"required"`
	Species        string    `json:"species" binding:"required"`
	Breed          string    `json:"breed"`
	Age            string    `json:"age"`
	Gender         string    `json:"gender" binding:"required"`
	Weight         float64   `json:"weight"`
	Color          string    `json:"color"`
	Description    string    `json:"description"`
	HealthStatus   string    `json:"health_status"`
	Vaccinated     bool      `json:"vaccinated"`
	Neutered       bool      `json:"neutered"`
	FoundLocation  string    `json:"found_location"`
	FoundDate      *time.Time `json:"found_date"`
	MedicalHistory string    `json:"medical_history"`
	Personality    string    `json:"personality"`
	SpecialNeeds   string    `json:"special_needs"`
}

type UpdatePetRequest struct {
	Name           string     `json:"name"`
	Species        string     `json:"species"`
	Breed          string     `json:"breed"`
	Age            string     `json:"age"`
	Gender         string     `json:"gender"`
	Weight         *float64   `json:"weight"`
	Color          string     `json:"color"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	HealthStatus   string     `json:"health_status"`
	Vaccinated     *bool      `json:"vaccinated"`
	Neutered       *bool      `json:"neutered"`
	FoundLocation  string     `json:"found_location"`
	FoundDate      *time.Time `json:"found_date"`
	MedicalHistory string     `json:"medical_history"`
	Personality    string     `json:"personality"`
	SpecialNeeds   string     `json:"special_needs"`
}

type PetListQuery struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	Status    string `form:"status"`
	Species   string `form:"species"`
	Gender    string `form:"gender"`
	Search    string `form:"search"`
	RescueID  uint   `form:"rescue_id"`
	Breed     string `form:"breed"`
}
