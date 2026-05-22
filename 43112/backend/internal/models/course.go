package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseStatus string

const (
	CourseDraft     CourseStatus = "draft"
	CoursePublished CourseStatus = "published"
	CourseOffline   CourseStatus = "offline"
	CourseRejected  CourseStatus = "rejected"
)

type CourseLevel string

const (
	LevelBeginner     CourseLevel = "beginner"
	LevelIntermediate CourseLevel = "intermediate"
	LevelAdvanced     CourseLevel = "advanced"
)

type Course struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	InstructorID  uuid.UUID    `gorm:"type:uuid;index;not null" json:"instructor_id"`
	Title         string       `gorm:"size:200;not null" json:"title"`
	Subtitle      string       `gorm:"size:500" json:"subtitle"`
	Description   string       `gorm:"type:text;not null" json:"description"`
	Cover         string       `gorm:"size:500" json:"cover"`
	Category      string       `gorm:"size:100;index" json:"category"`
	Level         CourseLevel  `gorm:"size:20;default:beginner" json:"level"`
	Price         float64      `gorm:"default:0" json:"price"`
	OriginalPrice float64      `json:"original_price"`
	Status        CourseStatus `gorm:"size:20;default:draft;index" json:"status"`
	IsFree        bool         `gorm:"default:false" json:"is_free"`
	Tags          string       `gorm:"size:500" json:"tags"`
	TotalHours    float64      `gorm:"default:0" json:"total_hours"`
	StudentCount  int          `gorm:"default:0" json:"student_count"`
	AvgRating     float64      `gorm:"default:0" json:"avg_rating"`
	ReviewCount   int          `gorm:"default:0" json:"review_count"`
	PublishedAt   *time.Time   `json:"published_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Instructor    *User        `gorm:"foreignKey:InstructorID" json:"instructor,omitempty"`
	Chapters      []Chapter    `json:"chapters,omitempty"`
	Reviews       []Review     `gorm:"foreignKey:CourseID" json:"-"`
	Orders        []Order      `gorm:"foreignKey:CourseID" json:"-"`
	Questions     []Question   `gorm:"foreignKey:CourseID" json:"-"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
