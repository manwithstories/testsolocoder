package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PostID    uint           `gorm:"index;not null" json:"post_id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	LikeCount int            `gorm:"not null;default:0" json:"like_count"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Post   *Post     `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User   *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}
