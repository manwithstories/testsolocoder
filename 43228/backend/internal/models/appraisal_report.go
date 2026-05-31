package models

import (
	"time"

	"gorm.io/gorm"
)

type AppraisalStatus string

const (
	AppraisalStatusPending   AppraisalStatus = "pending"
	AppraisalStatusCompleted AppraisalStatus = "completed"
)

type AppraisalReport struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	AppraiserID        uint            `gorm:"index;not null" json:"appraiser_id"`
	TeaID              uint            `gorm:"index;not null" json:"tea_id"`
	ReportNo           string          `gorm:"size:64;uniqueIndex;not null" json:"report_no"`
	IsAuthentic        bool            `gorm:"not null;default:false" json:"is_authentic"`
	GradeAssessment    string          `gorm:"size:64" json:"grade_assessment"`
	ValueEstimate      float64         `gorm:"type:decimal(12,2)" json:"value_estimate"`
	CollectionSuggestion string        `gorm:"type:text" json:"collection_suggestion"`
	ReportContent      string          `gorm:"type:text" json:"report_content"`
	Status             AppraisalStatus `gorm:"size:32;not null;default:pending;index" json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"-"`

	Appraiser *User `gorm:"foreignKey:AppraiserID" json:"appraiser,omitempty"`
	Tea       *Tea  `gorm:"foreignKey:TeaID" json:"tea,omitempty"`
}

func (AppraisalReport) TableName() string {
	return "appraisal_reports"
}
