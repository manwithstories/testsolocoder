package models

import (
	"time"

	"gorm.io/gorm"
)

type CertStatus string

const (
	CertStatusPending  CertStatus = "pending"
	CertStatusApproved CertStatus = "approved"
	CertStatusRejected CertStatus = "rejected"
)

type RoasterCertification struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	User           *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CertName       string         `json:"cert_name" gorm:"size:100"`
	CertNumber     string         `json:"cert_number" gorm:"size:100"`
	OrgName        string         `json:"org_name" gorm:"size:200"`
	CertFile       string         `json:"cert_file" gorm:"size:500"`
	Experience     string         `json:"experience" gorm:"type:text"`
	Specialty      string         `json:"specialty" gorm:"size:200"`
	Status         CertStatus     `json:"status" gorm:"size:20;default:pending"`
	ReviewerID     *uint          `json:"reviewer_id"`
	Reviewer       *User          `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
	ReviewComment  string         `json:"review_comment" gorm:"type:text"`
	ReviewedAt     *time.Time     `json:"reviewed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type ApplyCertificationRequest struct {
	CertName   string `json:"cert_name" binding:"required"`
	CertNumber string `json:"cert_number" binding:"required"`
	OrgName    string `json:"org_name" binding:"required"`
	CertFile   string `json:"cert_file"`
	Experience string `json:"experience" binding:"required"`
	Specialty  string `json:"specialty"`
}

type ReviewCertificationRequest struct {
	Status        CertStatus `json:"status" binding:"required,oneof=approved rejected"`
	ReviewComment string     `json:"review_comment"`
}

type RoasterProfile struct {
	User          *User                  `json:"user"`
	Certification *RoasterCertification `json:"certification"`
	Products      []Product              `json:"products"`
	RoastingRecords []RoastingRecord     `json:"roasting_records"`
	TotalProducts int64                  `json:"total_products"`
	TotalRoasts   int64                  `json:"total_roasts"`
	AvgScore      float64                `json:"avg_score"`
}
