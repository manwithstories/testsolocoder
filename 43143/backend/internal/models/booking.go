package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusRejected  BookingStatus = "rejected"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusNoShow    BookingStatus = "no_show"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusPaid       PaymentStatus = "paid"
	PaymentStatusHeld       PaymentStatus = "held"
	PaymentStatusReleased   PaymentStatus = "released"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusFailed     PaymentStatus = "failed"
)

type Booking struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PostingID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"posting_id"`
	Posting         SkillPosting   `gorm:"foreignKey:PostingID" json:"posting,omitempty"`
	StudentID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"student_id"`
	Student         User           `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	TeacherID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"teacher_id"`
	Teacher         User           `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	ScheduledStart  time.Time      `json:"scheduled_start"`
	ScheduledEnd    time.Time      `json:"scheduled_end"`
	ActualStart     *time.Time     `json:"actual_start"`
	ActualEnd       *time.Time     `json:"actual_end"`
	Status          BookingStatus  `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Price           float64        `gorm:"not null" json:"price"`
	PlatformFee     float64        `gorm:"default:0" json:"platform_fee"`
	TeacherEarnings float64        `gorm:"default:0" json:"teacher_earnings"`
	Note            string         `gorm:"type:text" json:"note"`
	RejectReason    string         `gorm:"size:500" json:"reject_reason"`
	CancelReason    string         `gorm:"size:500" json:"cancel_reason"`
	CancelledBy     *uuid.UUID     `gorm:"type:uuid" json:"cancelled_by"`
	ReviewedByStudent bool         `gorm:"default:false" json:"reviewed_by_student"`
	ReviewedByTeacher bool         `gorm:"default:false" json:"reviewed_by_teacher"`
	PaymentStatus   PaymentStatus  `gorm:"type:varchar(20);default:'pending'" json:"payment_status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type Review struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"booking_id"`
	Booking    Booking        `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	ReviewerID uuid.UUID      `gorm:"type:uuid;index;not null" json:"reviewer_id"`
	Reviewer   User           `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	RevieweeID uuid.UUID      `gorm:"type:uuid;index;not null" json:"reviewee_id"`
	Reviewee   User           `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
	PostingID  uuid.UUID      `gorm:"type:uuid;index" json:"posting_id"`
	Rating     int            `gorm:"not null" json:"rating"`
	Content    string         `gorm:"type:text" json:"content"`
	IsPublic   bool           `gorm:"default:true" json:"is_public"`
	IsFeatured bool           `gorm:"default:false" json:"is_featured"`
	HelpfulCount int          `gorm:"default:0" json:"helpful_count"`
	Status     string         `gorm:"type:varchar(20);default:'approved'" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type ReviewReply struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ReviewID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"review_id"`
	Review    Review         `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
	UserID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string         `gorm:"type:text" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *ReviewReply) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
