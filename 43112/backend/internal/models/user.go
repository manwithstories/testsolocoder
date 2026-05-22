package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleStudent    UserRole = "student"
	RoleInstructor UserRole = "instructor"
	RoleAdmin      UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type InstructorStatus string

const (
	InstructorPending   InstructorStatus = "pending"
	InstructorApproved  InstructorStatus = "approved"
	InstructorRejected  InstructorStatus = "rejected"
)

type User struct {
	ID                uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	Username          string           `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email             string           `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password          string           `gorm:"size:255;not null" json:"-"`
	Nickname          string           `gorm:"size:50" json:"nickname"`
	Avatar            string           `gorm:"size:500" json:"avatar"`
	Role              UserRole         `gorm:"size:20;default:student" json:"role"`
	Status            UserStatus       `gorm:"size:20;default:active" json:"status"`
	InstructorStatus  InstructorStatus `gorm:"size:20;default:pending" json:"instructor_status"`
	Phone             string           `gorm:"size:20" json:"phone"`
	Bio               string           `gorm:"type:text" json:"bio"`
	Credentials       string           `gorm:"type:text" json:"-"`
	EmailVerified     bool             `gorm:"default:false" json:"email_verified"`
	LastLoginAt       *time.Time       `json:"last_login_at"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `gorm:"index" json:"-"`

	Courses           []Course         `gorm:"foreignKey:InstructorID" json:"-"`
	Orders            []Order          `gorm:"foreignKey:UserID" json:"-"`
	Reviews           []Review         `gorm:"foreignKey:UserID" json:"-"`
	Questions         []Question       `gorm:"foreignKey:UserID" json:"-"`
	Answers           []Answer         `gorm:"foreignKey:UserID" json:"-"`
	Progresses        []Progress       `gorm:"foreignKey:UserID" json:"-"`
	Notes             []Note           `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
