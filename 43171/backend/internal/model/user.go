package model

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleClient    Role = "client"
	RolePilot     Role = "pilot"
	RoleOwner     Role = "owner"
)

type VerifyStatus string

const (
	VerifyPending  VerifyStatus = "pending"
	VerifyApproved VerifyStatus = "approved"
	VerifyRejected VerifyStatus = "rejected"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password     string         `gorm:"size:256;not null" json:"-"`
	Nickname     string         `gorm:"size:64" json:"nickname"`
	Phone        string         `gorm:"size:32" json:"phone"`
	Email        string         `gorm:"size:128" json:"email"`
	Avatar       string         `gorm:"size:256" json:"avatar"`
	Role         Role           `gorm:"size:16;not null;index" json:"role"`
	RealName     string         `gorm:"size:64" json:"real_name"`
	IDCardNo     string         `gorm:"size:32" json:"id_card_no"`
	LicenseNo    string         `gorm:"size:64" json:"license_no"`
	LicenseImage string         `gorm:"size:256" json:"license_image"`
	VerifyStatus VerifyStatus   `gorm:"size:16;default:pending" json:"verify_status"`
	VerifyRemark string         `gorm:"size:256" json:"verify_remark"`
	Rating       float64        `gorm:"default:5.0" json:"rating"`
	RatingCount  int            `gorm:"default:0" json:"rating_count"`
	Balance      float64        `gorm:"default:0" json:"balance"`
	Deposit      float64        `gorm:"default:0" json:"deposit"`
	Status       int            `gorm:"default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
