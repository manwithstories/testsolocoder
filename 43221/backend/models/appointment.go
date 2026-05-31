package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentStatus string

const (
	AppointmentPending   AppointmentStatus = "pending"
	AppointmentConfirmed AppointmentStatus = "confirmed"
	AppointmentCompleted AppointmentStatus = "completed"
	AppointmentCancelled AppointmentStatus = "cancelled"
	AppointmentRefunded  AppointmentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentOnline  PaymentMethod = "online"
	PaymentOffline PaymentMethod = "offline"
)

type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentPaid       PaymentStatus = "paid"
	PaymentRefunded   PaymentStatus = "refunded"
	PaymentFailed     PaymentStatus = "failed"
	PaymentCancelled  PaymentStatus = "cancelled"
)

type Appointment struct {
	ID              uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey"`
	ClientID        uuid.UUID         `json:"client_id" gorm:"type:uuid;not null;index"`
	ProfessionalID  uuid.UUID         `json:"professional_id" gorm:"type:uuid;not null;index"`
	ServiceID       uuid.UUID         `json:"service_id" gorm:"type:uuid;not null;index"`
	ScheduleID      uuid.UUID         `json:"schedule_id" gorm:"type:uuid;uniqueIndex"`
	Status          AppointmentStatus `json:"status" gorm:"type:varchar(20);default:pending"`
	Notes           string            `json:"notes" gorm:"type:text"`
	CancelReason    string            `json:"cancel_reason"`
	CancelledAt     *time.Time        `json:"cancelled_at"`
	CompletedAt     *time.Time        `json:"completed_at"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"-" gorm:"index"`

	Client          User              `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	Professional    User              `json:"professional,omitempty" gorm:"foreignKey:ProfessionalID"`
	Service         Service           `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	Schedule        Schedule          `json:"schedule,omitempty" gorm:"foreignKey:ScheduleID"`
	Payment         *Payment          `json:"payment,omitempty" gorm:"foreignKey:AppointmentID"`
	ConsultRecord   *ConsultRecord    `json:"consult_record,omitempty" gorm:"foreignKey:AppointmentID"`
	Review          *Review           `json:"review,omitempty" gorm:"foreignKey:AppointmentID"`
}

type Payment struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey"`
	AppointmentID   uuid.UUID     `json:"appointment_id" gorm:"type:uuid;not null;uniqueIndex"`
	ClientID        uuid.UUID     `json:"client_id" gorm:"type:uuid;not null;index"`
	ProfessionalID  uuid.UUID     `json:"professional_id" gorm:"type:uuid;not null;index"`
	Amount          float64       `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"type:varchar(20);not null"`
	Status          PaymentStatus `json:"status" gorm:"type:varchar(20);default:pending"`
	TransactionID   string        `json:"transaction_id" gorm:"size:100"`
	RefundReason    string        `json:"refund_reason"`
	RefundStatus    string        `json:"refund_status" gorm:"type:varchar(20)"`
	RefundedAt      *time.Time    `json:"refunded_at"`
	PaidAt          *time.Time    `json:"paid_at"`
	ExpiresAt       time.Time     `json:"expires_at"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Appointment     Appointment   `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
