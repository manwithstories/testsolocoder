package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleRoaster  UserRole = "roaster"
	RoleUser     UserRole = "user"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusPending  UserStatus = "pending"
)

type User struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Username        string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email           string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Phone           string         `json:"phone" gorm:"size:20"`
	Password        string         `json:"-" gorm:"size:255;not null"`
	Nickname        string         `json:"nickname" gorm:"size:50"`
	Avatar          string         `json:"avatar" gorm:"size:255"`
	Role            UserRole       `json:"role" gorm:"size:20;default:user"`
	Status          UserStatus     `json:"status" gorm:"size:20;default:active"`
	Bio             string         `json:"bio" gorm:"type:text"`
	Address         string         `json:"address" gorm:"size:500"`
	IsCertified     bool           `json:"is_certified" gorm:"default:false"`
	CertificationID *uint          `json:"certification_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Address  string `json:"address"`
}
