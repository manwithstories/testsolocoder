package models

import (
	"time"

	"gorm.io/gorm"
)

type PostCategory string

const (
	PostCategoryKnowledge PostCategory = "knowledge"
	PostCategorySharing   PostCategory = "sharing"
	PostCategoryActivity  PostCategory = "activity"
)

type PostStatus string

const (
	PostStatusPublished PostStatus = "published"
	PostStatusHidden    PostStatus = "hidden"
	PostStatusDeleted   PostStatus = "deleted"
)

type Post struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Title        string         `gorm:"size:255;not null;index" json:"title"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	Category     PostCategory   `gorm:"size:32;not null;index" json:"category"`
	CoverImage   string         `gorm:"size:512" json:"cover_image"`
	ViewCount    int            `gorm:"not null;default:0" json:"view_count"`
	LikeCount    int            `gorm:"not null;default:0" json:"like_count"`
	CommentCount int            `gorm:"not null;default:0" json:"comment_count"`
	Status       PostStatus     `gorm:"size:32;not null;default:published;index" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Likes    []Like    `gorm:"foreignKey:PostID" json:"likes,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}
