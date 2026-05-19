package model

import (
	"time"
)

type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Username      string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email         string     `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash  string     `gorm:"size:255;not null" json:"-"`
	RealName      string     `gorm:"size:50" json:"real_name"`
	Phone         string     `gorm:"size:20" json:"phone"`
	Avatar        string     `gorm:"size:255" json:"avatar"`
	Role          UserRole   `gorm:"size:20;default:'user'" json:"role"`
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`
	Status        UserStatus `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
