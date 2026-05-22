package models

import (
	"time"

	"gorm.io/gorm"
)

type ReviewType string

const (
	ReviewTypeCustomer ReviewType = "customer_to_provider"
	ReviewTypeProvider ReviewType = "provider_to_customer"
)

type Review struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OrderID        uint           `json:"order_id" gorm:"not null;uniqueIndex"`
	ReviewerID     uint           `json:"reviewer_id" gorm:"not null;index"`
	RevieweeID     uint           `json:"reviewee_id" gorm:"not null;index"`
	ReviewType     ReviewType     `json:"review_type" gorm:"size:30;not null"`
	AttitudeRating int            `json:"attitude_rating" gorm:"default:5"`
	PunctualRating int            `json:"punctual_rating" gorm:"default:5"`
	ProfessionalRating int        `json:"professional_rating" gorm:"default:5"`
	OverallRating  float64        `json:"overall_rating" gorm:"not null"`
	Content        string         `json:"content" gorm:"size:1000"`
	Images         string         `json:"images" gorm:"type:text"`
	IsAnonymous    bool           `json:"is_anonymous" gorm:"default:false"`
	ReplyContent   string         `json:"reply_content" gorm:"size:500"`
	RepliedAt      *time.Time     `json:"replied_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	Reviewer       *User          `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
	Reviewee       *User          `json:"reviewee,omitempty" gorm:"foreignKey:RevieweeID"`
}

type ComplaintStatus string

const (
	ComplaintStatusPending   ComplaintStatus = "pending"
	ComplaintStatusProcessing ComplaintStatus = "processing"
	ComplaintStatusResolved  ComplaintStatus = "resolved"
	ComplaintStatusRejected  ComplaintStatus = "rejected"
)

type Complaint struct {
	ID             uint             `json:"id" gorm:"primaryKey"`
	OrderID        uint             `json:"order_id" gorm:"not null;index"`
	ComplainantID  uint             `json:"complainant_id" gorm:"not null;index"`
	RespondentID   uint             `json:"respondent_id" gorm:"not null;index"`
	Status         ComplaintStatus  `json:"status" gorm:"size:20;not null;default:pending"`
	Title          string           `json:"title" gorm:"size:200;not null"`
	Content        string           `json:"content" gorm:"type:text;not null"`
	Images         string           `json:"images" gorm:"type:text"`
	HandlerID      *uint            `json:"handler_id" gorm:"index"`
	HandleResult   string           `json:"handle_result" gorm:"size:500"`
	HandledAt      *time.Time       `json:"handled_at"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `json:"-" gorm:"index"`

	Complainant    *User            `json:"complainant,omitempty" gorm:"foreignKey:ComplainantID"`
	Respondent     *User            `json:"respondent,omitempty" gorm:"foreignKey:RespondentID"`
	Order          *Order           `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}
