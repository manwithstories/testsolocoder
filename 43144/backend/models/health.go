package models

import (
	"time"

	"gorm.io/gorm"
)

type HealthRecordType string

const (
	HealthRecordVaccine  HealthRecordType = "vaccine"
	HealthRecordDeworm   HealthRecordType = "deworm"
	HealthRecordCheckup  HealthRecordType = "checkup"
	HealthRecordDisease  HealthRecordType = "disease"
	HealthRecordSurgery  HealthRecordType = "surgery"
	HealthRecordOther    HealthRecordType = "other"
)

type HealthRecord struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	PetID       uint              `json:"pet_id" gorm:"index;not null"`
	Pet         *Pet              `json:"pet,omitempty" gorm:"foreignKey:PetID"`
	RecordType  HealthRecordType  `json:"record_type" gorm:"size:20;not null"`
	Title       string            `json:"title" gorm:"size:255;not null"`
	Description string            `json:"description" gorm:"type:text"`
	VaccineName string            `json:"vaccine_name" gorm:"size:100"`
	RecordDate  time.Time         `json:"record_date"`
	NextDate    *time.Time        `json:"next_date"`
	Weight      float64           `json:"weight"`
	Temperature float64           `json:"temperature"`
	VetName     string            `json:"vet_name" gorm:"size:100"`
	Hospital    string            `json:"hospital" gorm:"size:200"`
	ReportFile  string            `json:"report_file" gorm:"size:500"`
	Notes       string            `json:"notes" gorm:"type:text"`
	RecordedBy  uint              `json:"recorded_by" gorm:"index;not null"`
	RescueID    uint              `json:"rescue_id" gorm:"index;not null"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `json:"-" gorm:"index"`
}

type CreateHealthRecordRequest struct {
	PetID       uint    `json:"pet_id" binding:"required"`
	RecordType  string  `json:"record_type" binding:"required,oneof=vaccine deworm checkup disease surgery other"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	VaccineName string  `json:"vaccine_name"`
	RecordDate  string  `json:"record_date" binding:"required"`
	NextDate    string  `json:"next_date"`
	Weight      float64 `json:"weight"`
	Temperature float64 `json:"temperature"`
	VetName     string  `json:"vet_name"`
	Hospital    string  `json:"hospital"`
	Notes       string  `json:"notes"`
}

type HealthRecordListQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	PetID      uint   `form:"pet_id"`
	RecordType string `form:"record_type"`
}

type HealthReminder struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PetID        uint      `json:"pet_id" gorm:"index;not null"`
	RecordID     *uint     `json:"record_id,omitempty" gorm:"index"`
	Title        string    `json:"title" gorm:"size:255;not null"`
	ReminderDate time.Time `json:"reminder_date"`
	IsCompleted  bool      `json:"is_completed" gorm:"default:false"`
	Notes        string    `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
