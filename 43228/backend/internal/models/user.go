package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleBuyer     UserRole = "buyer"
	UserRoleSeller    UserRole = "seller"
	UserRoleAppraiser UserRole = "appraiser"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Email     string         `gorm:"size:128;uniqueIndex" json:"email"`
	Phone     string         `gorm:"size:32;uniqueIndex" json:"phone"`
	Avatar    string         `gorm:"size:512" json:"avatar"`
	Role      UserRole       `gorm:"size:32;not null;default:buyer" json:"role"`
	Status    UserStatus     `gorm:"size:32;not null;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Collections     []Collection     `gorm:"foreignKey:UserID" json:"collections,omitempty"`
	Orders          []Order          `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	TastingRecords  []TastingRecord  `gorm:"foreignKey:UserID" json:"tasting_records,omitempty"`
	Posts           []Post           `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments        []Comment        `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Likes           []Like           `gorm:"foreignKey:UserID" json:"likes,omitempty"`
	Activities      []TastingActivity `gorm:"foreignKey:UserID" json:"activities,omitempty"`
	ActivityJoins   []ActivityParticipant `gorm:"foreignKey:UserID" json:"activity_joins,omitempty"`
	AppraisalReports []AppraisalReport `gorm:"foreignKey:AppraiserID" json:"appraisal_reports,omitempty"`
	OperationLogs   []OperationLog   `gorm:"foreignKey:UserID" json:"operation_logs,omitempty"`
}

func (User) TableName() string {
	return "users"
}
