package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Progress struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	CourseID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"course_id"`
	LessonID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"lesson_id"`
	LastPosition   float64        `gorm:"default:0" json:"last_position"`
	TotalDuration  float64        `gorm:"default:0" json:"total_duration"`
	IsCompleted    bool           `gorm:"default:false" json:"is_completed"`
	CompletedAt    *time.Time     `json:"completed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Note struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	LessonID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"lesson_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Timestamp  float64        `json:"timestamp"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User       *User          `json:"user,omitempty"`
}

type InstructorApplication struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	RealName     string         `gorm:"size:50;not null" json:"real_name"`
	IDCard       string         `gorm:"size:50;not null" json:"-"`
	Qualification string        `gorm:"type:text;not null" json:"qualification"`
	Experience   string         `gorm:"type:text" json:"experience"`
	Certificates string         `gorm:"size:500" json:"certificates"`
	Status       InstructorStatus `gorm:"size:20;default:pending;index" json:"status"`
	ReviewerID   *uuid.UUID     `gorm:"type:uuid" json:"reviewer_id,omitempty"`
	ReviewRemark string         `gorm:"size:500" json:"review_remark"`
	ReviewedAt   *time.Time     `json:"reviewed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User         *User          `json:"user,omitempty"`
}

type FileType string

const (
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeDocument FileType = "document"
	FileTypeOther    FileType = "other"
)

type File struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	FileType     FileType       `gorm:"size:20;not null" json:"file_type"`
	FileName     string         `gorm:"size:200;not null" json:"file_name"`
	FilePath     string         `gorm:"size:500;not null" json:"file_path"`
	FileURL      string         `gorm:"size:500" json:"file_url"`
	FileSize     int64          `gorm:"default:0" json:"file_size"`
	MimeType     string         `gorm:"size:100" json:"mime_type"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Progress) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (n *Note) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

func (a *InstructorApplication) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
