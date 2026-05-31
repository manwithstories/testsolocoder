package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleClient      UserRole = "client"
	RoleTranslator  UserRole = "translator"
	RolePM          UserRole = "pm"
	RoleAdmin       UserRole = "admin"
)

type ProjectStatus string

const (
	ProjectStatusPending   ProjectStatus = "pending"
	ProjectStatusApproved  ProjectStatus = "approved"
	ProjectStatusAssigned  ProjectStatus = "assigned"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusReview    ProjectStatus = "review"
	ProjectStatusRevision  ProjectStatus = "revision"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"
)

type DocumentType string

const (
	DocTypeWord    DocumentType = "word"
	DocTypeExcel   DocumentType = "excel"
	DocTypePDF     DocumentType = "pdf"
	DocTypePPT     DocumentType = "ppt"
	DocTypeTXT     DocumentType = "txt"
	DocTypeOther   DocumentType = "other"
)

type ReviewStatus string

const (
	ReviewStatusPending    ReviewStatus = "pending"
	ReviewStatusApproved   ReviewStatus = "approved"
	ReviewStatusRejected   ReviewStatus = "rejected"
)

type UrgencyLevel string

const (
	UrgencyNormal  UrgencyLevel = "normal"
	UrgencyFast    UrgencyLevel = "fast"
	UrgencyUrgent  UrgencyLevel = "urgent"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Username        string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password        string         `gorm:"size:255;not null" json:"-"`
	Email           string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone           string         `gorm:"size:20" json:"phone"`
	RealName        string         `gorm:"size:50" json:"real_name"`
	Avatar          string         `gorm:"size:255" json:"avatar"`
	Role            UserRole       `gorm:"size:20;not null;default:client" json:"role"`
	LanguagePairs   []LanguagePair `gorm:"many2many:user_language_pairs;" json:"language_pairs,omitempty"`
	ExpertiseTags   []ExpertiseTag `gorm:"many2many:user_expertise_tags;" json:"expertise_tags,omitempty"`
	Rating          float64        `gorm:"default:5.0" json:"rating"`
	CompletedCount  int            `gorm:"default:0" json:"completed_count"`
	CurrentWorkload int            `gorm:"default:0" json:"current_workload"`
	DailyCapacity   int            `gorm:"default:2000" json:"daily_capacity"`
	Status          string         `gorm:"size:20;default:active" json:"status"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type LanguagePair struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SourceLang   string    `gorm:"size:10;not null" json:"source_lang"`
	TargetLang   string    `gorm:"size:10;not null" json:"target_lang"`
	DisplayName  string    `gorm:"size:50" json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type ExpertiseTag struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Project struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	ClientID      uint           `gorm:"not null;index" json:"client_id"`
	Client        User           `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	SourceLang    string         `gorm:"size:10;not null" json:"source_lang"`
	TargetLang    string         `gorm:"size:10;not null" json:"target_lang"`
	ExpertiseTags []ExpertiseTag `gorm:"many2many:project_expertise_tags;" json:"expertise_tags,omitempty"`
	WordCount     int            `gorm:"default:0" json:"word_count"`
	Urgency       UrgencyLevel   `gorm:"size:20;default:normal" json:"urgency"`
	Difficulty    float64        `gorm:"default:1.0" json:"difficulty"`
	UnitPrice     float64        `gorm:"default:0" json:"unit_price"`
	TotalAmount   float64        `gorm:"default:0" json:"total_amount"`
	Deadline      time.Time      `json:"deadline"`
	PMID          *uint          `gorm:"index" json:"pm_id,omitempty"`
	PM            *User          `gorm:"foreignKey:PMID" json:"pm,omitempty"`
	TranslatorID  *uint          `gorm:"index" json:"translator_id,omitempty"`
	Translator    *User          `gorm:"foreignKey:TranslatorID" json:"translator,omitempty"`
	ReviewerID    *uint          `gorm:"index" json:"reviewer_id,omitempty"`
	Reviewer      *User          `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Status        ProjectStatus  `gorm:"size:20;default:pending;index" json:"status"`
	ReviewStatus  ReviewStatus   `gorm:"size:20;default:pending" json:"review_status"`
	Documents     []Document     `json:"documents,omitempty"`
	Comments      []ProjectComment `json:"comments,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type ProjectComment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProjectID uint           `gorm:"not null;index" json:"project_id"`
	Project   Project        `gorm:"foreignKey:ProjectID" json:"-"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

type Document struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProjectID   uint           `gorm:"not null;index" json:"project_id"`
	Project     Project        `gorm:"foreignKey:ProjectID" json:"-"`
	FileName    string         `gorm:"size:255;not null" json:"file_name"`
	FilePath    string         `gorm:"size:500;not null" json:"file_path"`
	FileType    DocumentType   `gorm:"size:20" json:"file_type"`
	FileSize    int64          `gorm:"default:0" json:"file_size"`
	WordCount   int            `gorm:"default:0" json:"word_count"`
	Version     int            `gorm:"default:1" json:"version"`
	IsSource    bool           `gorm:"default:true" json:"is_source"`
	UploadedBy  uint           `gorm:"not null" json:"uploaded_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
}

type TranslationSegment struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ProjectID      uint           `gorm:"not null;index" json:"project_id"`
	DocumentID     uint           `gorm:"index" json:"document_id"`
	SourceText     string         `gorm:"type:text;not null" json:"source_text"`
	TranslatedText string         `gorm:"type:text" json:"translated_text"`
	IsMemory       bool           `gorm:"default:false" json:"is_memory"`
	MemoryMatch    float64        `gorm:"default:0" json:"memory_match"`
	Status         string         `gorm:"size:20;default:pending" json:"status"`
	ReviewComment  string         `gorm:"type:text" json:"review_comment"`
	ReviewStatus   ReviewStatus   `gorm:"size:20;default:pending" json:"review_status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type TranslationMemory struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SourceText     string         `gorm:"type:text;not null" json:"source_text"`
	TranslatedText string         `gorm:"type:text;not null" json:"translated_text"`
	SourceLang     string         `gorm:"size:10;not null;index" json:"source_lang"`
	TargetLang     string         `gorm:"size:10;not null;index" json:"target_lang"`
	ProjectID      *uint          `gorm:"index" json:"project_id,omitempty"`
	ExpertiseTagID *uint          `gorm:"index" json:"expertise_tag_id,omitempty"`
	UsageCount     int            `gorm:"default:0" json:"usage_count"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt      time.Time      `json:"created_at"`
}

type GlossaryTerm struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SourceTerm     string         `gorm:"size:200;not null" json:"source_term"`
	TargetTerm     string         `gorm:"size:200;not null" json:"target_term"`
	SourceLang     string         `gorm:"size:10;not null;index" json:"source_lang"`
	TargetLang     string         `gorm:"size:10;not null;index" json:"target_lang"`
	Domain         string         `gorm:"size:50" json:"domain"`
	Definition     string         `gorm:"type:text" json:"definition"`
	PartOfSpeech   string         `gorm:"size:20" json:"part_of_speech"`
	CreatedBy      uint           `gorm:"not null" json:"created_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt      time.Time      `json:"created_at"`
}

type ReviewTask struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ProjectID    uint           `gorm:"not null;index" json:"project_id"`
	Project      Project        `gorm:"foreignKey:ProjectID" json:"-"`
	ReviewerID   uint           `gorm:"not null;index" json:"reviewer_id"`
	Reviewer     User           `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	SegmentID    uint           `gorm:"index" json:"segment_id"`
	Segment      TranslationSegment `gorm:"foreignKey:SegmentID" json:"segment,omitempty"`
	Round        int            `gorm:"default:1" json:"round"`
	Comment      string         `gorm:"type:text" json:"comment"`
	Suggestion   string         `gorm:"type:text" json:"suggestion"`
	Status       ReviewStatus   `gorm:"size:20;default:pending" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
}

type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProjectID     uint           `gorm:"not null;index" json:"project_id"`
	Project       Project        `gorm:"foreignKey:ProjectID" json:"-"`
	ClientID      uint           `gorm:"not null;index" json:"client_id"`
	Client        User           `gorm:"foreignKey:ClientID" json:"-"`
	TranslatorID  uint           `gorm:"index" json:"translator_id"`
	Translator    *User          `gorm:"foreignKey:TranslatorID" json:"translator,omitempty"`
	Amount        float64        `gorm:"not null" json:"amount"`
	BaseAmount    float64        `gorm:"not null" json:"base_amount"`
	UrgencyFee    float64        `gorm:"default:0" json:"urgency_fee"`
	DifficultyFee float64        `gorm:"default:0" json:"difficulty_fee"`
	Status        string         `gorm:"size:20;default:pending" json:"status"`
	PaidAt        *time.Time     `json:"paid_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
}

type OperationLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    *uint          `gorm:"index" json:"user_id,omitempty"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action    string         `gorm:"size:100;not null;index" json:"action"`
	Module    string         `gorm:"size:50;index" json:"module"`
	Detail    string         `gorm:"type:text" json:"detail"`
	IPAddress string         `gorm:"size:50" json:"ip_address"`
	UserAgent string         `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
}

type Statistics struct {
	Date              time.Time `gorm:"uniqueIndex" json:"date"`
	NewProjects       int       `gorm:"default:0" json:"new_projects"`
	CompletedProjects int       `gorm:"default:0" json:"completed_projects"`
	TotalWords        int       `gorm:"default:0" json:"total_words"`
	Revenue           float64   `gorm:"default:0" json:"revenue"`
	CreatedAt         time.Time `json:"created_at"`
}
