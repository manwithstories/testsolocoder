package models

import (
	"time"

	"gorm.io/gorm"
)

type Beehive struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null;index" json:"user_id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Code           string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Latitude       float64        `gorm:"type:decimal(10,7);not null" json:"latitude"`
	Longitude      float64        `gorm:"type:decimal(10,7);not null" json:"longitude"`
	Region         string         `gorm:"size:100;index" json:"region"`
	BeeSpecies     string         `gorm:"size:50;index" json:"bee_species"`
	GroupName      string         `gorm:"size:50;index" json:"group_name"`
	Status         string         `gorm:"size:20;default:active" json:"status"`
	HealthStatus   string         `gorm:"size:20;default:healthy" json:"health_status"`
	QueenStatus    string         `gorm:"size:20;default:normal" json:"queen_status"`
	WorkerCount    int            `gorm:"default:0" json:"worker_count"`
	LastInspection *time.Time     `gorm:"type:date" json:"last_inspection"`
	Notes          string         `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Beehive) TableName() string {
	return "beehives"
}

type CreateBeehiveRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Code        string  `json:"code" binding:"required,max=50"`
	Latitude    float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude   float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
	Region      string  `json:"region" binding:"max=100"`
	BeeSpecies  string  `json:"bee_species" binding:"max=50"`
	GroupName   string  `json:"group_name" binding:"max=50"`
	QueenStatus string  `json:"queen_status"`
	WorkerCount int     `json:"worker_count" binding:"gte=0"`
	Notes       string  `json:"notes"`
}

type UpdateBeehiveRequest struct {
	Name         *string  `json:"name" binding:"omitempty,max=100"`
	Latitude     *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude    *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	Region       *string  `json:"region" binding:"omitempty,max=100"`
	BeeSpecies   *string  `json:"bee_species" binding:"omitempty,max=50"`
	GroupName    *string  `json:"group_name" binding:"omitempty,max=50"`
	Status       *string  `json:"status"`
	HealthStatus *string  `json:"health_status"`
	QueenStatus  *string  `json:"queen_status"`
	WorkerCount  *int     `json:"worker_count" binding:"omitempty,gte=0"`
	Notes        *string  `json:"notes"`
}
