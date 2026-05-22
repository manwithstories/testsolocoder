package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID       uuid.UUID      `gorm:"type:uuid;index" json:"plan_id"`
	ActivityID   uuid.UUID      `gorm:"type:uuid;index" json:"activity_id"`
	UploadedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"uploaded_by"`
	FileName     string         `gorm:"type:varchar(255);not null" json:"file_name"`
	OriginalName string         `gorm:"type:varchar(255);not null" json:"original_name"`
	FileType     string         `gorm:"type:varchar(100)" json:"file_type"`
	FileSize     int64          `gorm:"default:0" json:"file_size"`
	FileURL      string         `gorm:"type:varchar(500)" json:"file_url"`
	Category     string         `gorm:"type:varchar(50)" json:"category"`
	Description  string         `gorm:"type:text" json:"description"`
	Tags         []string       `gorm:"type:varchar(255)[]" json:"tags"`
	IsPublic     bool           `gorm:"default:false" json:"is_public"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Plan     *TravelPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Uploader *User       `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}
