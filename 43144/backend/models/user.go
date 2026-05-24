package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleRescue   UserRole = "rescue"
	RoleAdopter  UserRole = "adopter"
)

type User struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Email           string         `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password        string         `json:"-" gorm:"not null"`
	Name            string         `json:"name" gorm:"size:100;not null"`
	Phone           string         `json:"phone" gorm:"size:20"`
	Role            UserRole       `json:"role" gorm:"size:20;default:adopter;not null"`
	RescueID        *uint          `json:"rescue_id,omitempty" gorm:"index"`
	Rescue          *RescueStation `json:"rescue,omitempty" gorm:"foreignKey:RescueID"`
	Address         string         `json:"address" gorm:"size:500"`
	IsVerified      bool           `json:"is_verified" gorm:"default:false"`
	Avatar          string         `json:"avatar" gorm:"size:500"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role" binding:"required,oneof=adopter rescue"`
	RescueName string `json:"rescue_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	RescueID *uint  `json:"rescue_id,omitempty"`
}
