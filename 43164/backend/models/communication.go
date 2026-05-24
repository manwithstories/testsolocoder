package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"bookingId"`
	TeacherID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacherId"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	ReviewerID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"reviewerId"`
	RevieweeID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"revieweeId"`
	Rating      int            `gorm:"not null" json:"rating" binding:"required,min=1,max=5"`
	Content     string         `gorm:"type:text;not null" json:"content" binding:"required"`
	Tags        string         `json:"tags"`
	IsAnonymous bool           `gorm:"default:false" json:"isAnonymous"`
	TeacherReply string        `gorm:"type:text" json:"teacherReply"`
	TeacherRepliedAt *time.Time `json:"teacherRepliedAt"`
	IsHidden    bool           `gorm:"default:false" json:"isHidden"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Booking  *Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Teacher  *User    `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Student  *User    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Reviewer *User    `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Reviewee *User    `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type MessageType string

const (
	MessageTypeText    MessageType = "text"
	MessageTypeFile    MessageType = "file"
	MessageTypeImage   MessageType = "image"
	MessageTypeSystem  MessageType = "system"
	MessageTypeBooking MessageType = "booking"
)

type Message struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SenderID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"senderId"`
	ReceiverID uuid.UUID      `gorm:"type:uuid;not null;index" json:"receiverId"`
	BookingID  *uuid.UUID     `gorm:"type:uuid;index" json:"bookingId"`
	Type       MessageType    `gorm:"type:varchar(20);not null;default:text" json:"type"`
	Content    string         `gorm:"type:text;not null" json:"content" binding:"required"`
	IsRead     bool           `gorm:"default:false" json:"isRead"`
	ReadAt     *time.Time     `json:"readAt"`
	IsDeleted  bool           `gorm:"default:false" json:"isDeleted"`
	DeletedBy  *uuid.UUID     `gorm:"type:uuid" json:"deletedBy"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Sender   *User          `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver *User          `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Booking  *Booking       `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Files    []MessageFile  `gorm:"foreignKey:MessageID" json:"files,omitempty"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type MessageFile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	MessageID uuid.UUID      `gorm:"type:uuid;not null;index" json:"messageId"`
	FileName  string         `gorm:"not null" json:"fileName"`
	FileURL   string         `gorm:"not null" json:"fileUrl"`
	FileSize  int64          `json:"fileSize"`
	FileType  string         `json:"fileType"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Message *Message `gorm:"foreignKey:MessageID" json:"-"`
}

func (mf *MessageFile) BeforeCreate(tx *gorm.DB) error {
	if mf.ID == uuid.Nil {
		mf.ID = uuid.New()
	}
	return nil
}

type NotificationType string

const (
	NotificationTypeBookingCreated   NotificationType = "booking_created"
	NotificationTypeBookingUpdated   NotificationType = "booking_updated"
	NotificationTypeBookingCancelled NotificationType = "booking_cancelled"
	NotificationTypeBookingReminder  NotificationType = "booking_reminder"
	NotificationTypeLessonStart      NotificationType = "lesson_start"
	NotificationTypeLessonEnd        NotificationType = "lesson_end"
	NotificationTypeNewReview        NotificationType = "new_review"
	NotificationTypePaymentReceived   NotificationType = "payment_received"
	NotificationTypePaymentFailed    NotificationType = "payment_failed"
	NotificationTypeWithdrawApproved NotificationType = "withdraw_approved"
	NotificationTypeWithdrawRejected NotificationType = "withdraw_rejected"
	NotificationTypeNewMessage       NotificationType = "new_message"
	NotificationTypeHomeworkAssigned NotificationType = "homework_assigned"
	NotificationTypeHomeworkGraded   NotificationType = "homework_graded"
	NotificationTypeMilestone        NotificationType = "milestone"
	NotificationTypeSystem           NotificationType = "system"
)

type Notification struct {
	ID       uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Type     NotificationType `gorm:"type:varchar(30);not null" json:"type"`
	Title    string         `gorm:"not null" json:"title"`
	Content  string         `gorm:"type:text;not null" json:"content"`
	Data     string         `gorm:"type:text" json:"data"`
	IsRead   bool           `gorm:"default:false" json:"isRead"`
	ReadAt   *time.Time     `json:"readAt"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type AdminAction struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AdminID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"adminId"`
	ActionType  string         `gorm:"not null" json:"actionType"`
	TargetType  string         `gorm:"not null" json:"targetType"`
	TargetID    uuid.UUID      `gorm:"type:uuid;not null" json:"targetId"`
	Details     string         `gorm:"type:text" json:"details"`
	IPAddress   string         `json:"ipAddress"`
	UserAgent   string         `json:"userAgent"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Admin *User `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
}

func (aa *AdminAction) BeforeCreate(tx *gorm.DB) error {
	if aa.ID == uuid.Nil {
		aa.ID = uuid.New()
	}
	return nil
}

type SystemLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Level     string         `gorm:"type:varchar(20);not null" json:"level"`
	Module    string         `gorm:"type:varchar(50);not null" json:"module"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Context   string         `gorm:"type:text" json:"context"`
	UserID    *uuid.UUID     `gorm:"type:uuid;index" json:"userId"`
	IPAddress string         `json:"ipAddress"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (sl *SystemLog) BeforeCreate(tx *gorm.DB) error {
	if sl.ID == uuid.Nil {
		sl.ID = uuid.New()
	}
	return nil
}
