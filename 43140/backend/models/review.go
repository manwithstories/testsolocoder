package models

import (
	"time"

	"gorm.io/gorm"
)

type ReviewStatus string

const (
	ReviewStatusOffer     ReviewStatus = "offer"
	ReviewStatusPass      ReviewStatus = "pass"
	ReviewStatusReject    ReviewStatus = "reject"
	ReviewStatusPending   ReviewStatus = "pending"
)

type Review struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	InterviewID uint           `gorm:"uniqueIndex;not null" json:"interview_id"`
	Rating      int            `gorm:"default:0" json:"rating"`
	Feedback    string         `gorm:"type:text" json:"feedback"`
	Strengths   string         `gorm:"type:text" json:"strengths"`
	Weaknesses  string         `gorm:"type:text" json:"weaknesses"`
	Status      ReviewStatus   `gorm:"type:varchar(20);default:pending" json:"status"`
	Interview   Interview      `gorm:"foreignKey:InterviewID" json:"interview,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
