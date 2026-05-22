package model

import (
	"time"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusRejected UserStatus = "rejected"
)

type User struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Username         string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password         string         `gorm:"size:255;not null" json:"-"`
	Email            string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone            string         `gorm:"size:20" json:"phone"`
	RealName         string         `gorm:"size:50" json:"real_name"`
	IDCard           string         `gorm:"size:20" json:"id_card"`
	LicenseNumber    string         `gorm:"size:50" json:"license_number"`
	LicenseImage     string         `gorm:"size:255" json:"license_image"`
	IDCardFront      string         `gorm:"size:255" json:"id_card_front"`
	IDCardBack       string         `gorm:"size:255" json:"id_card_back"`
	AuthStatus       UserStatus     `gorm:"size:20;default:pending" json:"auth_status"`
	Status           UserStatus     `gorm:"size:20;default:pending" json:"status"`
	RoleID           uint           `json:"role_id"`
	Role             *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Avatar           string         `gorm:"size:255" json:"avatar"`
	LastLoginAt      *time.Time     `json:"last_login_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        *time.Time     `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type Role struct {
	ID          uint         `gorm:"primarykey" json:"id"`
	Name        string       `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Resource    string    `gorm:"size:50;not null" json:"resource"`
	Action      string    `gorm:"size:50;not null" json:"action"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}