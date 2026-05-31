package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOwner    Role = "owner"
	RoleTenant   Role = "tenant"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Username        string         `gorm:"uniqueIndex;size:100;not null" json:"username" binding:"required,min=3,max=100"`
	Email           string         `gorm:"uniqueIndex;size:200;not null" json:"email" binding:"required,email"`
	Password        string         `gorm:"size:255;not null" json:"-"`
	Phone           string         `gorm:"size:20" json:"phone"`
	Role            Role           `gorm:"type:varchar(20);default:tenant;not null" json:"role"`
	FullName        string         `gorm:"size:200" json:"full_name"`
	AvatarURL       string         `gorm:"size:500" json:"avatar_url"`
	Company         string         `gorm:"size:200" json:"company"`
	Address         string         `gorm:"size:500" json:"address"`
	City            string         `gorm:"size:100" json:"city"`
	Country         string         `gorm:"size:100" json:"country"`
	Timezone        string         `gorm:"size:50;default:UTC" json:"timezone"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Role     Role   `json:"role" binding:"required,oneof=admin owner tenant"`
	FullName string `json:"full_name" binding:"required,min=2,max=200"`
	Phone    string `json:"phone"`
	Company  string `json:"company"`
}

type UpdateProfileRequest struct {
	FullName   string `json:"full_name"`
	Phone      string `json:"phone"`
	Company    string `json:"company"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Timezone   string `json:"timezone"`
	AvatarURL  string `json:"avatar_url"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
