package models

import (
	"time"

	"gorm.io/gorm"
)

type JobSeeker struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	User           *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BirthDate      *time.Time     `json:"birth_date"`
	Gender         string         `json:"gender"`
	Location       string         `json:"location"`
	CurrentTitle   string         `json:"current_title"`
	CurrentCompany string         `json:"current_company"`
	Experience     string         `json:"experience"`
	Education      string         `json:"education"`
	Resumes        []Resume       `gorm:"foreignKey:JobSeekerID" json:"resumes,omitempty"`
	Applications   []Application  `gorm:"foreignKey:JobSeekerID" json:"applications,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
