package model

import (
	"time"

	"gorm.io/gorm"
)

type Survey struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Title           string         `json:"title" gorm:"size:200;not null"`
	Description     string         `json:"description" gorm:"type:text"`
	CoverImage      string         `json:"cover_image" gorm:"size:500"`
	UserID          uint           `json:"user_id" gorm:"index;not null"`
	User            *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Status          int            `json:"status" gorm:"default:1;comment:1=draft,2=published,3=closed"`
	StartTime       *time.Time     `json:"start_time"`
	EndTime         *time.Time     `json:"end_time"`
	Anonymous       bool           `json:"anonymous" gorm:"default:false"`
	Password        string         `json:"-" gorm:"size:100"`
	MaxResponses    int            `json:"max_responses" gorm:"default:0;comment:0=unlimited"`
	MaxPerUser      int            `json:"max_per_user" gorm:"default:1;comment:0=unlimited per user"`
	RequiresLogin   bool           `json:"requires_login" gorm:"default:false"`
	AllowResume     bool           `json:"allow_resume" gorm:"default:true"`
	ResponseCount   int            `json:"response_count" gorm:"default:0"`
	Category        string         `json:"category" gorm:"size:50"`
	Tags            string         `json:"tags" gorm:"size:500"`
	Questions       []Question     `json:"questions,omitempty" gorm:"foreignKey:SurveyID"`
	DistributionLinks []DistributionLink `json:"distribution_links,omitempty" gorm:"foreignKey:SurveyID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (s *Survey) IsClosed() bool {
	if s.Status == 3 {
		return true
	}
	if s.EndTime != nil && time.Now().After(*s.EndTime) {
		return true
	}
	if s.MaxResponses > 0 && s.ResponseCount >= s.MaxResponses {
		return true
	}
	return false
}

func (s *Survey) IsStarted() bool {
	if s.Status == 1 {
		return false
	}
	if s.StartTime != nil && time.Now().Before(*s.StartTime) {
		return false
	}
	return true
}
