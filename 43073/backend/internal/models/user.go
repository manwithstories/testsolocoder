package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Password  string         `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email" validate:"omitempty,email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Role      string         `gorm:"size:20;default:'user'" json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
