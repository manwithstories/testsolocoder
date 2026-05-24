package model

import (
	"time"

	"gorm.io/gorm"
)

type Response struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	SurveyID        uint           `json:"survey_id" gorm:"index;not null"`
	Survey          *Survey        `json:"survey,omitempty" gorm:"foreignKey:SurveyID"`
	UserID          *uint          `json:"user_id" gorm:"index"`
	User            *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DistributionID  *uint          `json:"distribution_id" gorm:"index"`
	Distribution    *DistributionLink `json:"distribution,omitempty" gorm:"foreignKey:DistributionID"`
	SessionID       string         `json:"session_id" gorm:"size:100;index"`
	IPAddress       string         `json:"ip_address" gorm:"size:50"`
	UserAgent       string         `json:"user_agent" gorm:"size:500"`
	Status          int            `json:"status" gorm:"default:1;comment:1=in_progress,2=completed,3=abandoned"`
	StartTime       *time.Time     `json:"start_time"`
	CompleteTime    *time.Time     `json:"complete_time"`
	Duration        int            `json:"duration" gorm:"comment:seconds"`
	Answers         []Answer       `json:"answers,omitempty" gorm:"foreignKey:ResponseID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type Answer struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ResponseID   uint           `json:"response_id" gorm:"index;not null"`
	Response     *Response      `json:"response,omitempty" gorm:"foreignKey:ResponseID"`
	QuestionID   uint           `json:"question_id" gorm:"index;not null"`
	Question     *Question      `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	OptionID     *uint          `json:"option_id" gorm:"index"`
	Option       *Option        `json:"option,omitempty" gorm:"foreignKey:OptionID"`
	TextValue    string         `json:"text_value" gorm:"type:text"`
	NumericValue float64        `json:"numeric_value" gorm:"default:0"`
	MatrixValues string         `json:"matrix_values" gorm:"type:text;comment:JSON for matrix questions"`
	RankingOrder string         `json:"ranking_order" gorm:"size:200;comment:JSON array"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type DistributionLink struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	SurveyID     uint           `json:"survey_id" gorm:"index;not null"`
	Survey       *Survey        `json:"survey,omitempty" gorm:"foreignKey:SurveyID"`
	LinkToken    string         `json:"link_token" gorm:"size:100;uniqueIndex;not null"`
	Channel      string         `json:"channel" gorm:"size:20;comment:email,qrcode,wechat,dm,api"`
	MaxUses      int            `json:"max_uses" gorm:"default:0;comment:0=unlimited"`
	UseCount     int            `json:"use_count" gorm:"default:0"`
	ExpiresAt    *time.Time     `json:"expires_at"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Invitation struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	SurveyID       uint           `json:"survey_id" gorm:"index;not null"`
	Survey         *Survey        `json:"survey,omitempty" gorm:"foreignKey:SurveyID"`
	Email          string         `json:"email" gorm:"size:100;not null"`
	LinkToken      string         `json:"link_token" gorm:"size:100;uniqueIndex;not null"`
	Status         int            `json:"status" gorm:"default:1;comment:1=pending,2=sent,3=opened,4=responded,5=failed"`
	SentAt         *time.Time     `json:"sent_at"`
	OpenedAt       *time.Time     `json:"opened_at"`
	RespondedAt    *time.Time     `json:"responded_at"`
	RetryCount     int            `json:"retry_count" gorm:"default:0"`
	ErrorMessage   string         `json:"error_message" gorm:"size:500"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
