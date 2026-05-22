package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleSpaceAdmin UserRole = "space_admin"
	RoleUser       UserRole = "user"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password     string         `json:"-" gorm:"size:255;not null"`
	RealName     string         `json:"real_name" gorm:"size:50"`
	Phone        string         `json:"phone" gorm:"size:20"`
	Department   string         `json:"department" gorm:"size:50"`
	Role         UserRole       `json:"role" gorm:"size:20;default:user"`
	Floor        string         `json:"floor" gorm:"size:50"`
	Status       int            `json:"status" gorm:"default:1"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
