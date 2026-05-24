package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonNote struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"bookingId"`
	TeacherID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacherId"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	SubjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"subjectId"`
	Title       string         `gorm:"not null" json:"title" binding:"required"`
	Content     string         `gorm:"type:text;not null" json:"content" binding:"required"`
	Tags        string         `json:"tags"`
	IsPrivate   bool           `gorm:"default:false" json:"isPrivate"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Booking *Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Teacher *User    `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Student *User    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Subject *Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

func (ln *LessonNote) BeforeCreate(tx *gorm.DB) error {
	if ln.ID == uuid.Nil {
		ln.ID = uuid.New()
	}
	return nil
}

type HomeworkStatus string

const (
	HomeworkStatusPending    HomeworkStatus = "pending"
	HomeworkStatusSubmitted  HomeworkStatus = "submitted"
	HomeworkStatusGraded     HomeworkStatus = "graded"
)

type Homework struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"bookingId"`
	TeacherID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacherId"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	SubjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"subjectId"`
	Title       string         `gorm:"not null" json:"title" binding:"required"`
	Description string         `gorm:"type:text;not null" json:"description" binding:"required"`
	DueDate     time.Time      `gorm:"not null" json:"dueDate"`
	Status      HomeworkStatus `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	MaxScore    float64        `json:"maxScore"`
	Attachments string         `json:"attachments"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Booking    *Booking              `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Teacher    *User                 `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Student    *User                 `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Subject    *Subject              `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Submission *HomeworkSubmission   `gorm:"foreignKey:HomeworkID" json:"submission,omitempty"`
}

func (h *Homework) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

type HomeworkSubmission struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	HomeworkID  uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"homeworkId"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	Content     string         `gorm:"type:text" json:"content"`
	Attachments string         `json:"attachments"`
	SubmittedAt time.Time      `json:"submittedAt"`
	Score       float64        `json:"score"`
	Feedback    string         `gorm:"type:text" json:"feedback"`
	GradedAt    *time.Time     `json:"gradedAt"`
	GradedBy    *uuid.UUID     `gorm:"type:uuid" json:"gradedBy"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Homework *Homework `gorm:"foreignKey:HomeworkID" json:"-"`
	Student  *User     `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

func (hs *HomeworkSubmission) BeforeCreate(tx *gorm.DB) error {
	if hs.ID == uuid.Nil {
		hs.ID = uuid.New()
	}
	return nil
}

type Feedback struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"bookingId"`
	TeacherID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacherId"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	SubjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"subjectId"`
	Content     string         `gorm:"type:text;not null" json:"content" binding:"required"`
	Type        string         `gorm:"type:varchar(20);not null;default:general" json:"type"`
	Progress    string         `json:"progress"`
	Suggestions string         `gorm:"type:text" json:"suggestions"`
	NextSteps   string         `gorm:"type:text" json:"nextSteps"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Booking *Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Teacher *User    `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Student *User    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Subject *Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

func (f *Feedback) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type MilestoneType string

const (
	MilestoneTypeSessionComplete MilestoneType = "session_complete"
	MilestoneTypeGoalAchieved    MilestoneType = "goal_achieved"
	MilestoneTypeScoreImprove    MilestoneType = "score_improve"
	MilestoneTypeStreak          MilestoneType = "streak"
	MilestoneTypeCustom          MilestoneType = "custom"
)

type Milestone struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	SubjectID   *uuid.UUID     `gorm:"type:uuid;index" json:"subjectId"`
	Title       string         `gorm:"not null" json:"title" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	Type        MilestoneType  `gorm:"type:varchar(20);not null;default:custom" json:"type"`
	Icon        string         `json:"icon"`
	Color       string         `json:"color"`
	IsAchieved  bool           `gorm:"default:false" json:"isAchieved"`
	AchievedAt  *time.Time     `json:"achievedAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Student *StudentProfile `gorm:"foreignKey:StudentID" json:"-"`
	Subject *Subject        `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

func (m *Milestone) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
