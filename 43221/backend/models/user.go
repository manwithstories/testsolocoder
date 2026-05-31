package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleClient      Role = "client"
	RoleProfessional Role = "professional"
	RoleAdmin       Role = "admin"
)

type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationApproved VerificationStatus = "approved"
	VerificationRejected VerificationStatus = "rejected"
)

type User struct {
	ID               uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey"`
	Username         string             `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email            string             `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password         string             `json:"-" gorm:"not null"`
	Role             Role               `json:"role" gorm:"type:varchar(20);not null;default:client"`
	FullName         string             `json:"full_name" gorm:"size:100"`
	Avatar           string             `json:"avatar"`
	Phone            string             `json:"phone" gorm:"size:20"`
	VerificationStatus VerificationStatus `json:"verification_status" gorm:"type:varchar(20);default:pending"`
	VerificationDocs string             `json:"-" gorm:"type:text"`
	VerificationNote string             `json:"verification_note"`
	IsActive         bool               `json:"is_active" gorm:"default:true"`
	LastLoginAt      *time.Time         `json:"last_login_at"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	DeletedAt        gorm.DeletedAt     `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
