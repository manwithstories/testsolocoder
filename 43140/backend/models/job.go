package models

import (
	"time"

	"gorm.io/gorm"
)

type JobType string

const (
	JobTypeFullTime   JobType = "full-time"
	JobTypePartTime   JobType = "part-time"
	JobTypeContract   JobType = "contract"
	JobTypeInternship JobType = "internship"
	JobTypeRemote     JobType = "remote"
)

type JobStatus string

const (
	JobStatusOpen     JobStatus = "open"
	JobStatusPaused   JobStatus = "paused"
	JobStatusClosed   JobStatus = "closed"
)

type Job struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CompanyID       uint           `gorm:"not null;index" json:"company_id"`
	DepartmentID    *uint          `json:"department_id"`
	Title           string         `gorm:"not null" json:"title"`
	Location        string         `json:"location"`
	SalaryMin       float64        `json:"salary_min"`
	SalaryMax       float64        `json:"salary_max"`
	SalaryType      string         `gorm:"type:varchar(20);default:monthly" json:"salary_type"`
	Description     string         `gorm:"type:text" json:"description"`
	Requirements    string         `gorm:"type:text" json:"requirements"`
	Skills          string         `gorm:"type:text" json:"skills"`
	JobType         JobType        `gorm:"type:varchar(20);not null" json:"job_type"`
	Status          JobStatus      `gorm:"type:varchar(20);default:open" json:"status"`
	Views           int            `gorm:"default:0" json:"views"`
	Applications    []Application  `gorm:"foreignKey:JobID" json:"applications,omitempty"`
	Company         Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Department      *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
