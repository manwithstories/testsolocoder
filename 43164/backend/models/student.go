package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentProfile struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	GradeLevel       string         `json:"gradeLevel"`
	School           string         `json:"school"`
	LearningStyle    string         `json:"learningStyle"`
	PreferredTime    string         `json:"preferredTime"`
	Notes            string         `gorm:"type:text" json:"notes"`
	ParentName       string         `json:"parentName"`
	ParentPhone      string         `json:"parentPhone"`
	ParentEmail      string         `json:"parentEmail"`
	AssessmentStatus string         `gorm:"type:varchar(20);default:not_started" json:"assessmentStatus"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	User           *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LearningGoals  []LearningGoal   `gorm:"foreignKey:StudentID" json:"learningGoals,omitempty"`
	Assessments    []AssessmentAnswer `gorm:"foreignKey:StudentID" json:"assessments,omitempty"`
}

func (sp *StudentProfile) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

type LearningGoal struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	SubjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"subjectId"`
	Title       string         `gorm:"not null" json:"title" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	TargetScore float64        `json:"targetScore"`
	CurrentScore float64       `json:"currentScore"`
	Deadline    *time.Time     `json:"deadline"`
	IsAchieved  bool           `gorm:"default:false" json:"isAchieved"`
	AchievedAt  *time.Time     `json:"achievedAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Student *StudentProfile `gorm:"foreignKey:StudentID" json:"-"`
	Subject *Subject        `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

func (lg *LearningGoal) BeforeCreate(tx *gorm.DB) error {
	if lg.ID == uuid.Nil {
		lg.ID = uuid.New()
	}
	return nil
}

type AssessmentQuestion struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SubjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"subjectId"`
	Question    string         `gorm:"type:text;not null" json:"question" binding:"required"`
	QuestionType string        `gorm:"type:varchar(20);not null;default:text" json:"questionType"`
	Options     string         `gorm:"type:text" json:"options"`
	Difficulty  string         `gorm:"type:varchar(20);default:medium" json:"difficulty"`
	SortOrder   int            `gorm:"default:0" json:"sortOrder"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Subject *Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

func (aq *AssessmentQuestion) BeforeCreate(tx *gorm.DB) error {
	if aq.ID == uuid.Nil {
		aq.ID = uuid.New()
	}
	return nil
}

type AssessmentAnswer struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"studentId"`
	QuestionID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"questionId"`
	Answer        string         `gorm:"type:text" json:"answer"`
	Score         float64        `json:"score"`
	CompletedAt   *time.Time     `json:"completedAt"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Student  *StudentProfile      `gorm:"foreignKey:StudentID" json:"-"`
	Question *AssessmentQuestion  `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (aa *AssessmentAnswer) BeforeCreate(tx *gorm.DB) error {
	if aa.ID == uuid.Nil {
		aa.ID = uuid.New()
	}
	return nil
}
