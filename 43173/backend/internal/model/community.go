package model

import (
	"time"

	"gorm.io/gorm"
)

type Follow struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FollowerID  uint           `json:"follower_id" gorm:"index;not null"`
	FollowingID uint           `json:"following_id" gorm:"index;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Follow) TableName() string {
	return "follows"
}

type Comment struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	WorkID     uint           `json:"work_id" gorm:"index"`
	AlbumID    *uint          `json:"album_id" gorm:"index"`
	PlaylistID *uint          `json:"playlist_id" gorm:"index"`
	ParentID   *uint          `json:"parent_id" gorm:"index"`
	Content    string         `json:"content" gorm:"type:text;not null"`
	LikeCount  int64          `json:"like_count" gorm:"default:0"`
	ReplyCount int64          `json:"reply_count" gorm:"default:0"`
	IsPinned   bool           `json:"is_pinned" gorm:"default:false"`
	Status     int            `json:"status" gorm:"default:1"`
	User       *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Replies    []Comment      `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Comment) TableName() string {
	return "comments"
}

type Playlist struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"index;not null"`
	Title        string         `json:"title" gorm:"size:200;not null"`
	Description  string         `json:"description" gorm:"type:text"`
	CoverURL     string         `json:"cover_url" gorm:"size:255"`
	IsPublic     bool           `json:"is_public" gorm:"default:true"`
	PlayCount    int64          `json:"play_count" gorm:"default:0"`
	LikeCount    int64          `json:"like_count" gorm:"default:0"`
	FollowCount  int64          `json:"follow_count" gorm:"default:0"`
	WorkCount    int            `json:"work_count" gorm:"default:0"`
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Works        []Work         `json:"works,omitempty" gorm:"many2many:playlist_works"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Playlist) TableName() string {
	return "playlists"
}

type PlaylistWork struct {
	PlaylistID uint `gorm:"primaryKey"`
	WorkID     uint `gorm:"primaryKey"`
	AddedAt    time.Time `json:"added_at"`
}

func (PlaylistWork) TableName() string {
	return "playlist_works"
}

type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	WorkID    *uint          `json:"work_id" gorm:"index"`
	AlbumID   *uint          `json:"album_id" gorm:"index"`
	CommentID *uint          `json:"comment_id" gorm:"index"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Like) TableName() string {
	return "likes"
}

type PlayRecord struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	WorkID    uint           `json:"work_id" gorm:"index;not null"`
	ArtistID  uint           `json:"artist_id" gorm:"index;not null"`
	Duration  int            `json:"duration"`
	IP        string         `json:"ip" gorm:"size:45"`
	Device    string         `json:"device" gorm:"size:100"`
	Source    string         `json:"source" gorm:"size:50"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (PlayRecord) TableName() string {
	return "play_records"
}

type Share struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	WorkID     *uint          `json:"work_id" gorm:"index"`
	AlbumID    *uint          `json:"album_id" gorm:"index"`
	PlaylistID *uint          `json:"playlist_id" gorm:"index"`
	Platform   string         `json:"platform" gorm:"size:50"`
	IP         string         `json:"ip" gorm:"size:45"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Share) TableName() string {
	return "shares"
}

type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	Type      string         `json:"type" gorm:"size:50"`
	Title     string         `json:"title" gorm:"size:200"`
	Content   string         `json:"content" gorm:"type:text"`
	Data      string         `json:"data" gorm:"type:text"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Notification) TableName() string {
	return "notifications"
}
