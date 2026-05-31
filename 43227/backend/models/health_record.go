package models

import (
	"time"

	"gorm.io/gorm"
)

type HealthRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	BeehiveID       uint           `gorm:"not null;index" json:"beehive_id"`
	RecordDate      time.Time      `gorm:"type:date;not null" json:"record_date"`
	QueenStatus     string         `gorm:"size:20" json:"queen_status"`
	WorkerCount     int            `json:"worker_count"`
	HasDisease      bool           `gorm:"default:false" json:"has_disease"`
	DiseaseType     string         `gorm:"size:50" json:"disease_type"`
	DiseaseSeverity string         `gorm:"size:20" json:"disease_severity"`
	Treatment       string         `gorm:"type:text" json:"treatment"`
	Season          string         `gorm:"size:20" json:"season"`
	Recommendations string         `gorm:"type:text" json:"recommendations"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	Beehive         Beehive        `gorm:"foreignKey:BeehiveID" json:"beehive,omitempty"`
}

func (HealthRecord) TableName() string {
	return "health_records"
}

type CreateHealthRecordRequest struct {
	BeehiveID       uint   `json:"beehive_id" binding:"required"`
	RecordDate      string `json:"record_date" binding:"required"`
	QueenStatus     string `json:"queen_status"`
	WorkerCount     int    `json:"worker_count"`
	HasDisease      bool   `json:"has_disease"`
	DiseaseType     string `json:"disease_type"`
	DiseaseSeverity string `json:"disease_severity"`
	Treatment       string `json:"treatment"`
	Season          string `json:"season"`
	Recommendations string `json:"recommendations"`
	Notes           string `json:"notes"`
}

type DiseaseWarning struct {
	ID          uint   `json:"id"`
	BeehiveID   uint   `json:"beehive_id"`
	BeehiveName string `json:"beehive_name"`
	BeehiveCode string `json:"beehive_code"`
	DiseaseType string `json:"disease_type"`
	Severity    string `json:"severity"`
	WarningTime string `json:"warning_time"`
}
