package models

import (
	"time"
)

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"index:idx_post_user,unique;not null" json:"post_id"`
	UserID    uint      `gorm:"index:idx_post_user,unique;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	Post *Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Like) TableName() string {
	return "likes"
}
