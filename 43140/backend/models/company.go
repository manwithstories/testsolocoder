package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	CompanyName     string         `gorm:"not null" json:"company_name"`
	Industry        string         `json:"industry"`
	CompanySize     string         `json:"company_size"`
	CompanyType     string         `json:"company_type"`
	Address         string         `json:"address"`
	Website         string         `json:"website"`
	Logo            string         `json:"logo"`
	Description     string         `gorm:"type:text" json:"description"`
	Departments     []Department   `gorm:"foreignKey:CompanyID" json:"departments,omitempty"`
	Jobs            []Job          `gorm:"foreignKey:CompanyID" json:"jobs,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Department struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CompanyID uint           `gorm:"not null" json:"company_id"`
	Name      string         `gorm:"not null" json:"name"`
	Jobs      []Job          `gorm:"foreignKey:DepartmentID" json:"jobs,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
