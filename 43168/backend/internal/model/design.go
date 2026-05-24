package model

import (
	"time"
)

const (
	ProjectStatusDraft     = "draft"
	ProjectStatusSubmitted = "submitted"
	ProjectStatusApproved  = "approved"
	ProjectStatusRejected  = "rejected"
)

const (
	CommentTypeText   = "text"
	CommentTypeMarker = "marker"
)

type DesignProject struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DesignerID  uint      `gorm:"index;not null" json:"designer_id"`
	OwnerID     uint      `gorm:"index;not null" json:"owner_id"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:32;not null;default:draft;index" json:"status"`
	CoverImage  string    `gorm:"size:255" json:"cover_image"`
	RoomType    string    `gorm:"size:64" json:"room_type"`
	Area        float64   `gorm:"default:0" json:"area"`
	Budget      float64   `gorm:"default:0" json:"budget"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DesignProject) TableName() string {
	return "design_projects"
}

func ValidProjectStatus(status string) bool {
	switch status {
	case ProjectStatusDraft, ProjectStatusSubmitted, ProjectStatusApproved, ProjectStatusRejected:
		return true
	default:
		return false
	}
}

type DesignImage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	ImageURL    string    `gorm:"size:512;not null" json:"image_url"`
	Description string    `gorm:"size:255" json:"description"`
	Sort        int       `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
}

func (DesignImage) TableName() string {
	return "design_images"
}

type DesignComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index;not null" json:"project_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	UserRole  string    `gorm:"size:32;not null" json:"user_role"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Type      string    `gorm:"size:32;not null;default:text" json:"type"`
	PositionX float64   `gorm:"default:0" json:"position_x"`
	PositionY float64   `gorm:"default:0" json:"position_y"`
	ParentID  uint      `gorm:"default:0" json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (DesignComment) TableName() string {
	return "design_comments"
}

func ValidCommentType(typ string) bool {
	switch typ {
	case CommentTypeText, CommentTypeMarker:
		return true
	default:
		return false
	}
}
