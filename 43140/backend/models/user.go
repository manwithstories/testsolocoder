package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleCompany   Role = "company"
	RoleJobSeeker Role = "jobseeker"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	Password     string         `gorm:"not null" json:"-"`
	Role         Role           `gorm:"type:varchar(20);not null" json:"role"`
	Name         string         `gorm:"not null" json:"name"`
	Phone        string         `json:"phone"`
	Avatar       string         `json:"avatar"`
	Status       string         `gorm:"type:varchar(20);default:active" json:"status"`
	Company      *Company       `gorm:"foreignKey:UserID" json:"company,omitempty"`
	JobSeeker    *JobSeeker     `gorm:"foreignKey:UserID" json:"jobseeker,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
