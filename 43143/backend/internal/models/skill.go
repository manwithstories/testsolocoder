package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillCategory struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Icon      string         `gorm:"size:500" json:"icon"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Skills    []Skill        `gorm:"foreignKey:CategoryID" json:"skills,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *SkillCategory) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type SkillTag struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	CategoryID uuid.UUID      `gorm:"type:uuid;index" json:"category_id"`
	Category   SkillCategory  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	UsageCount int            `gorm:"default:0" json:"usage_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *SkillTag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "beginner"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyAdvanced     DifficultyLevel = "advanced"
	DifficultyExpert       DifficultyLevel = "expert"
)

type Skill struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string          `gorm:"size:200;not null" json:"title"`
	Description   string          `gorm:"type:text" json:"description"`
	CategoryID    uuid.UUID       `gorm:"type:uuid;index;not null" json:"category_id"`
	Category      SkillCategory   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags          []SkillTag      `gorm:"many2many:skill_tags;constraint:OnDelete:CASCADE" json:"tags,omitempty"`
	Difficulty    DifficultyLevel `gorm:"type:varchar(20);default:'beginner'" json:"difficulty"`
	CoverImage    string          `gorm:"size:500" json:"cover_image"`
	VideoURL      string          `gorm:"size:500" json:"video_url"`
	Prerequisites string          `gorm:"type:text" json:"prerequisites"`
	Outcomes      string          `gorm:"type:text" json:"outcomes"`
	IsActive      bool            `gorm:"default:true" json:"is_active"`
	PostingCount  int             `gorm:"default:0" json:"posting_count"`
	Rating        float64         `gorm:"default:0" json:"rating"`
	ReviewCount   int             `gorm:"default:0" json:"review_count"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (s *Skill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type TeachingMethod string

const (
	TeachingMethodOnline  TeachingMethod = "online"
	TeachingMethodOffline TeachingMethod = "offline"
	TeachingMethodBoth    TeachingMethod = "both"
)

type TeachingMode string

const (
	TeachingModeOneToOne  TeachingMode = "one_to_one"
	TeachingModeSmallClass TeachingMode = "small_class"
	TeachingModeGroup     TeachingMode = "group"
)

type SkillPosting struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TeacherID        uuid.UUID      `gorm:"type:uuid;index;not null" json:"teacher_id"`
	Teacher          User           `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	SkillID          uuid.UUID      `gorm:"type:uuid;index;not null" json:"skill_id"`
	Skill            Skill          `gorm:"foreignKey:SkillID" json:"skill,omitempty"`
	Title            string         `gorm:"size:200;not null" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	TeachingMethod   TeachingMethod `gorm:"type:varchar(20);default:'online'" json:"teaching_method"`
	TeachingMode     TeachingMode   `gorm:"type:varchar(20);default:'one_to_one'" json:"teaching_mode"`
	MaxStudents      int            `gorm:"default:1" json:"max_students"`
	PricePerHour     float64        `gorm:"not null" json:"price_per_hour"`
	Currency         string         `gorm:"size:10;default:'CNY'" json:"currency"`
	SessionDuration  int            `gorm:"default:60" json:"session_duration"`
	Location         string         `gorm:"size:200" json:"location"`
	Latitude         float64        `json:"latitude"`
	Longitude        float64        `json:"longitude"`
	Availability     string         `gorm:"type:jsonb" json:"availability"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	Rating           float64        `gorm:"default:0" json:"rating"`
	ReviewCount      int            `gorm:"default:0" json:"review_count"`
	BookingCount     int            `gorm:"default:0" json:"booking_count"`
	TotalHours       float64        `gorm:"default:0" json:"total_hours"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *SkillPosting) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
