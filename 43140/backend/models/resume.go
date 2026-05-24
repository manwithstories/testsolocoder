package models

import (
	"time"

	"gorm.io/gorm"
)

type Resume struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	JobSeekerID     uint             `gorm:"not null;index" json:"jobseeker_id"`
	Title           string           `gorm:"not null" json:"title"`
	FullName        string           `gorm:"not null" json:"full_name"`
	Email           string           `json:"email"`
	Phone           string           `json:"phone"`
	Location        string           `json:"location"`
	Summary         string           `gorm:"type:text" json:"summary"`
	FilePath        string           `json:"file_path"`
	FileName        string           `json:"file_name"`
	IsDefault       bool             `gorm:"default:false" json:"is_default"`
	EducationList   []Education      `gorm:"foreignKey:ResumeID" json:"education_list,omitempty"`
	WorkExperiences []WorkExperience `gorm:"foreignKey:ResumeID" json:"work_experiences,omitempty"`
	Skills          []Skill          `gorm:"foreignKey:ResumeID" json:"skills,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
}

type Education struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ResumeID    uint           `gorm:"not null" json:"resume_id"`
	School      string         `gorm:"not null" json:"school"`
	Degree      string         `json:"degree"`
	Major       string         `json:"major"`
	StartDate   string         `json:"start_date"`
	EndDate     string         `json:"end_date"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WorkExperience struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ResumeID    uint           `gorm:"not null" json:"resume_id"`
	Company     string         `gorm:"not null" json:"company"`
	Position    string         `gorm:"not null" json:"position"`
	StartDate   string         `json:"start_date"`
	EndDate     string         `json:"end_date"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Skill struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ResumeID  uint           `gorm:"not null" json:"resume_id"`
	Name      string         `gorm:"not null" json:"name"`
	Level     string         `gorm:"type:varchar(20)" json:"level"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
