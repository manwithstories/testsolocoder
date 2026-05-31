package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleHR       UserRole = "hr"
	RoleEmployee UserRole = "employee"
	RoleAgency   UserRole = "agency"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	RealName     string         `gorm:"size:50" json:"real_name"`
	Phone        string         `gorm:"size:20" json:"phone"`
	Email        string         `gorm:"size:100" json:"email"`
	Role         UserRole       `gorm:"size:20;not null;default:employee" json:"role"`
	Status       int            `gorm:"default:1" json:"status"`
	CompanyID    *uint          `gorm:"index" json:"company_id"`
	AgencyID     *uint          `gorm:"index" json:"agency_id"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Company      *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Agency       *Agency        `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
}
