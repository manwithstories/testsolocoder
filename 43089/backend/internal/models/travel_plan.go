package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TravelPlan struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Description string         `gorm:"type:text" json:"description"`
	Destination string         `gorm:"type:varchar(200);not null" json:"destination" validate:"required"`
	StartDate   time.Time      `gorm:"not null" json:"start_date" validate:"required"`
	EndDate     time.Time      `gorm:"not null" json:"end_date" validate:"required,gtefield=StartDate"`
	Budget      float64        `gorm:"type:decimal(12,2);default:0" json:"budget"`
	Currency    string         `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	Status      string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	Version     int            `gorm:"default:1" json:"version"`
	CoverImage  string         `gorm:"type:varchar(255)" json:"cover_image"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Owner        *User             `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Participants []PlanParticipant `gorm:"foreignKey:PlanID" json:"participants,omitempty"`
	Activities   []Activity        `gorm:"foreignKey:PlanID" json:"activities,omitempty"`
	Files        []File            `gorm:"foreignKey:PlanID" json:"files,omitempty"`
	Checklists   []Checklist       `gorm:"foreignKey:PlanID" json:"checklists,omitempty"`
	Reminders    []Reminder        `gorm:"foreignKey:PlanID" json:"reminders,omitempty"`
}

type PlanParticipant struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null;index" json:"plan_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Role      string    `gorm:"type:varchar(20);default:'viewer'" json:"role"`
	CanEdit   bool      `gorm:"default:false" json:"can_edit"`
	CanDelete bool      `gorm:"default:false" json:"can_delete"`
	JoinedAt  time.Time `json:"joined_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Plan *TravelPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	User *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (tp *TravelPlan) BeforeUpdate(tx *gorm.DB) error {
	tp.Version++
	return nil
}
