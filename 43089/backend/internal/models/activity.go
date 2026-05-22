package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"plan_id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Description string         `gorm:"type:text" json:"description"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type" validate:"required,oneof=sightseeing transport accommodation food other"`
	Date        time.Time      `gorm:"not null" json:"date"`
	StartTime   string         `gorm:"type:varchar(10)" json:"start_time"`
	EndTime     string         `gorm:"type:varchar(10)" json:"end_time"`
	Location    string         `gorm:"type:varchar(500)" json:"location"`
	Latitude    float64        `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude   float64        `gorm:"type:decimal(10,7)" json:"longitude"`
	Cost        float64        `gorm:"type:decimal(12,2);default:0" json:"cost"`
	Currency    string         `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	Notes       string         `gorm:"type:text" json:"notes"`
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Booked      bool           `gorm:"default:false" json:"booked"`
	Confirmation string        `gorm:"type:varchar(200)" json:"confirmation"`
	ContactInfo string         `gorm:"type:text" json:"contact_info"`
	OrderIndex  int            `gorm:"default:0" json:"order_index"`
	Version     int            `gorm:"default:1" json:"version"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Plan *TravelPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Files []File     `gorm:"many2many:activity_files;" json:"files,omitempty"`
}

type ActivityFile struct {
	ActivityID uuid.UUID `gorm:"type:uuid;primaryKey" json:"activity_id"`
	FileID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"file_id"`
	CreatedAt  time.Time `json:"created_at"`
}
