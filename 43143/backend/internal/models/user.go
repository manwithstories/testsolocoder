package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleLearner UserRole = "learner"
	RoleTeacher UserRole = "teacher"
	RoleBoth    UserRole = "both"
	RoleAdmin   UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
)

type AuthType string

const (
	AuthTypeEmail    AuthType = "email"
	AuthTypePhone    AuthType = "phone"
	AuthTypeWechat   AuthType = "wechat"
	AuthTypeGoogle   AuthType = "google"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email           string         `gorm:"uniqueIndex;size:100" json:"email"`
	Phone           string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Password        string         `gorm:"size:255" json:"-"`
	Nickname        string         `gorm:"size:50;not null" json:"nickname"`
	Avatar          string         `gorm:"size:500" json:"avatar"`
	Bio             string         `gorm:"size:500" json:"bio"`
	Gender          string         `gorm:"size:10" json:"gender"`
	Birthday        *time.Time     `json:"birthday"`
	Location        string         `gorm:"size:200" json:"location"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	Role            UserRole       `gorm:"type:varchar(20);default:'learner'" json:"role"`
	Status          UserStatus     `gorm:"type:varchar(20);default:'active'" json:"status"`
	AuthType        AuthType       `gorm:"type:varchar(20)" json:"auth_type"`
	EmailVerified   bool           `gorm:"default:false" json:"email_verified"`
	PhoneVerified   bool           `gorm:"default:false" json:"phone_verified"`
	SkillTags       []SkillTag     `gorm:"many2many:user_skills;constraint:OnDelete:CASCADE" json:"skill_tags,omitempty"`
	LearnTags       []SkillTag     `gorm:"many2many:user_learn_skills;constraint:OnDelete:CASCADE" json:"learn_tags,omitempty"`
	TeachTags       []SkillTag     `gorm:"many2many:user_teach_skills;constraint:OnDelete:CASCADE" json:"teach_tags,omitempty"`
	TeachingHours   float64        `gorm:"default:0" json:"teaching_hours"`
	StudentCount    int            `gorm:"default:0" json:"student_count"`
	Rating          float64        `gorm:"default:0" json:"rating"`
	ReviewCount     int            `gorm:"default:0" json:"review_count"`
	Balance         float64        `gorm:"default:0" json:"balance"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	LastLoginIP     string         `gorm:"size:50" json:"last_login_ip"`
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

type Certification struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SkillTagID uuid.UUID `gorm:"type:uuid;index" json:"skill_tag_id"`
	SkillTag   SkillTag  `gorm:"foreignKey:SkillTagID" json:"skill_tag,omitempty"`
	CertType   string    `gorm:"type:varchar(50);not null" json:"cert_type"`
	CertName   string    `gorm:"size:200;not null" json:"cert_name"`
	CertFile   string    `gorm:"size:500" json:"cert_file"`
	Status     string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	VerifiedAt *time.Time `json:"verified_at"`
	Remark     string    `gorm:"size:500" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Certification) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
