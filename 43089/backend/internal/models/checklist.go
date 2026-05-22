package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Checklist struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"plan_id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Type      string         `gorm:"type:varchar(20);default:'packing'" json:"type"`
	CreatedBy uuid.UUID      `gorm:"type:uuid" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Plan  *TravelPlan      `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Items []ChecklistItem `gorm:"foreignKey:ChecklistID" json:"items,omitempty"`
}

type ChecklistItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ChecklistID uuid.UUID      `gorm:"type:uuid;not null;index" json:"checklist_id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"type:varchar(50)" json:"category"`
	Quantity    int            `gorm:"default:1" json:"quantity"`
	IsCompleted bool           `gorm:"default:false" json:"is_completed"`
	CompletedBy uuid.UUID      `gorm:"type:uuid" json:"completed_by"`
	CompletedAt *time.Time     `json:"completed_at"`
	OrderIndex  int            `gorm:"default:0" json:"order_index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Checklist *Checklist `gorm:"foreignKey:ChecklistID" json:"checklist,omitempty"`
}
