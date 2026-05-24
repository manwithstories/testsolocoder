package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending    BookingStatus = "pending"
	BookingStatusConfirmed  BookingStatus = "confirmed"
	BookingStatusCompleted  BookingStatus = "completed"
	BookingStatusCancelled  BookingStatus = "cancelled"
	BookingStatusRescheduled BookingStatus = "rescheduled"
)

type Booking struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	TeacherID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacherId"`
	SubjectID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"subjectId"`
	StartTime      time.Time      `gorm:"not null" json:"startTime"`
	EndTime        time.Time      `gorm:"not null" json:"endTime"`
	Duration       int            `gorm:"not null" json:"duration"`
	HourlyRate     float64        `gorm:"not null" json:"hourlyRate"`
	TotalAmount    float64        `gorm:"not null" json:"totalAmount"`
	Status         BookingStatus  `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	Notes          string         `gorm:"type:text" json:"notes"`
	StudentNotes   string         `gorm:"type:text" json:"studentNotes"`
	TeacherNotes   string         `gorm:"type:text" json:"teacherNotes"`
	CancelledBy    *uuid.UUID     `gorm:"type:uuid" json:"cancelledBy"`
	CancelledAt    *time.Time     `json:"cancelledAt"`
	CancelReason   string         `json:"cancelReason"`
	OriginalBookingID *uuid.UUID  `gorm:"type:uuid" json:"originalBookingId"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Student *User           `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher *User           `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Subject *Subject        `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	VideoSession *VideoSession `gorm:"foreignKey:BookingID" json:"videoSession,omitempty"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BookingRescheduleHistory struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"bookingId"`
	OldStartTime    time.Time      `json:"oldStartTime"`
	OldEndTime      time.Time      `json:"oldEndTime"`
	NewStartTime    time.Time      `json:"newStartTime"`
	NewEndTime      time.Time      `json:"newEndTime"`
	RescheduledBy   uuid.UUID      `gorm:"type:uuid;not null" json:"rescheduledBy"`
	Reason          string         `json:"reason"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Booking *Booking `gorm:"foreignKey:BookingID" json:"-"`
}

func (brh *BookingRescheduleHistory) BeforeCreate(tx *gorm.DB) error {
	if brh.ID == uuid.Nil {
		brh.ID = uuid.New()
	}
	return nil
}

type VideoSession struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID      uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"bookingId"`
	SessionID      string         `gorm:"uniqueIndex;not null" json:"sessionId"`
	RoomName       string         `json:"roomName"`
	Token          string         `json:"token"`
	JoinURL        string         `json:"joinUrl"`
	Status         string         `gorm:"type:varchar(20);not null;default:scheduled" json:"status"`
	ActualStartAt  *time.Time     `json:"actualStartAt"`
	ActualEndAt    *time.Time     `json:"actualEndAt"`
	ActualDuration int            `json:"actualDuration"`
	QualityScore   float64        `json:"qualityScore"`
	ConnectionLogs string         `gorm:"type:text" json:"connectionLogs"`
	RecordingURL   string         `json:"recordingUrl"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Booking *Booking `gorm:"foreignKey:BookingID" json:"-"`
	Logs    []VideoSessionLog `gorm:"foreignKey:SessionID" json:"logs,omitempty"`
}

func (vs *VideoSession) BeforeCreate(tx *gorm.DB) error {
	if vs.ID == uuid.Nil {
		vs.ID = uuid.New()
	}
	return nil
}

type VideoSessionLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"sessionId"`
	EventType string         `gorm:"not null" json:"eventType"`
	Timestamp time.Time      `json:"timestamp"`
	Details   string         `gorm:"type:text" json:"details"`
	UserID    *uuid.UUID     `gorm:"type:uuid" json:"userId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Session *VideoSession `gorm:"foreignKey:SessionID" json:"-"`
}

func (vsl *VideoSessionLog) BeforeCreate(tx *gorm.DB) error {
	if vsl.ID == uuid.Nil {
		vsl.ID = uuid.New()
	}
	return nil
}
