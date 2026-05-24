package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperationType string

const (
	OperationCreate OperationType = "create"
	OperationUpdate OperationType = "update"
	OperationDelete OperationType = "delete"
	OperationLogin  OperationType = "login"
	OperationPay    OperationType = "pay"
	OperationReview OperationType = "review"
	OperationCancel OperationType = "cancel"
)

type OperationLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      *uuid.UUID     `gorm:"type:uuid;index" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TargetType  string         `gorm:"size:50" json:"target_type"`
	TargetID    *uuid.UUID     `gorm:"type:uuid" json:"target_id"`
	Operation   OperationType  `gorm:"type:varchar(20)" json:"operation"`
	Detail      string         `gorm:"type:text" json:"detail"`
	IP          string         `gorm:"size:50" json:"ip"`
	UserAgent   string         `gorm:"size:500" json:"user_agent"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (o *OperationLog) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type ComplaintType string

const (
	ComplaintReview ComplaintType = "review"
	ComplaintBooking ComplaintType = "booking"
	ComplaintPayment ComplaintType = "payment"
	ComplaintUser   ComplaintType = "user"
)

type ComplaintStatus string

const (
	ComplaintStatusPending   ComplaintStatus = "pending"
	ComplaintStatusProcessing ComplaintStatus = "processing"
	ComplaintStatusResolved  ComplaintStatus = "resolved"
	ComplaintStatusRejected  ComplaintStatus = "rejected"
)

type Complaint struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ReporterID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"reporter_id"`
	Reporter    User           `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	RespondentID uuid.UUID     `gorm:"type:uuid;index;not null" json:"respondent_id"`
	Respondent  User           `gorm:"foreignKey:RespondentID" json:"respondent,omitempty"`
	Type        ComplaintType  `gorm:"type:varchar(20)" json:"type"`
	TargetID    uuid.UUID      `gorm:"type:uuid;index" json:"target_id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Evidence    string         `gorm:"type:text" json:"evidence"`
	Status      ComplaintStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	HandlerID   *uuid.UUID     `gorm:"type:uuid;index" json:"handler_id"`
	Handler     *User          `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
	Result      string         `gorm:"type:text" json:"result"`
	HandledAt   *time.Time     `json:"handled_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Complaint) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
