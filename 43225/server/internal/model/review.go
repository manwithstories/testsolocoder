package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewType string

const (
	ReviewTypeShip ReviewType = "ship"
	ReviewTypeDock ReviewType = "dock"
)

type Review struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RentalID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"rental_id"`
	Rental        Rental         `gorm:"foreignKey:RentalID" json:"rental,omitempty"`
	ReviewerID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"reviewer_id"`
	Reviewer      User           `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	TargetType    ReviewType     `gorm:"type:varchar(20);not null" json:"target_type"`
	TargetID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"target_id"`
	Rating        int            `gorm:"not null;default:5" json:"rating" binding:"required,min=1,max=5"`
	Content       string         `gorm:"type:text" json:"content"`
	Response      string         `gorm:"type:text" json:"response"`
	ResponseBy    *uuid.UUID     `gorm:"type:uuid" json:"response_by"`
	Responder     *User          `gorm:"foreignKey:ResponseBy" json:"responder,omitempty"`
	ResponseAt    *time.Time     `json:"response_at"`
	IsRecommended bool           `gorm:"default:true" json:"is_recommended"`
	HelpfulCount  int            `gorm:"default:0" json:"helpful_count"`
	IsDeleted     bool           `gorm:"default:false" json:"is_deleted"`
	DeletedAt     *time.Time     `json:"deleted_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type CreateReviewRequest struct {
	RentalID      string     `json:"rental_id" binding:"required,uuid"`
	TargetType    ReviewType `json:"target_type" binding:"required,oneof=ship dock"`
	TargetID      string     `json:"target_id" binding:"required,uuid"`
	Rating        int        `json:"rating" binding:"required,min=1,max=5"`
	Content       string     `json:"content"`
	IsRecommended bool       `json:"is_recommended"`
}

type RespondToReviewRequest struct {
	Response string `json:"response" binding:"required,min=1"`
}

type SearchReviewRequest struct {
	TargetType ReviewType `form:"target_type" binding:"required,oneof=ship dock"`
	TargetID   string     `form:"target_id" binding:"required,uuid"`
	MinRating  int        `form:"min_rating"`
	MaxRating  int        `form:"max_rating"`
	Page       int        `form:"page,default=1"`
	PageSize   int        `form:"page_size,default=10"`
	SortBy     string     `form:"sort_by,default=created_at"`
	SortOrder  string     `form:"sort_order,default=desc"`
}
