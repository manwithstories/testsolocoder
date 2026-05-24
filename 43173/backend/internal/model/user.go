package model

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleArtist    UserRole = "artist"
	RoleFan       UserRole = "fan"
	RoleLabel     UserRole = "label"
)

type UserStatus int

const (
	UserStatusInactive UserStatus = 0
	UserStatusActive   UserStatus = 1
	UserStatusBanned   UserStatus = 2
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:100"`
	Phone        string         `json:"phone" gorm:"size:20"`
	Password     string         `json:"-" gorm:"size:255;not null"`
	Nickname     string         `json:"nickname" gorm:"size:50"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	Bio          string         `json:"bio" gorm:"type:text"`
	Role         UserRole       `json:"role" gorm:"size:20;default:fan"`
	Status       UserStatus     `json:"status" gorm:"default:1"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `json:"last_login_ip" gorm:"size:45"`
	ArtistInfo   *ArtistInfo    `json:"artist_info,omitempty" gorm:"foreignKey:UserID"`
	Followers    []Follow       `json:"-" gorm:"foreignKey:FollowingID"`
	Followings   []Follow       `json:"-" gorm:"foreignKey:FollowerID"`
	PlayRecords  []PlayRecord   `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type ArtistInfo struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	RealName        string         `json:"real_name" gorm:"size:50"`
	IDCard          string         `json:"id_card" gorm:"size:20"`
	ArtistName      string         `json:"artist_name" gorm:"size:100"`
	Genre           string         `json:"genre" gorm:"size:100"`
	Label           string         `json:"label" gorm:"size:100"`
	Website         string         `json:"website" gorm:"size:255"`
	Facebook        string         `json:"facebook" gorm:"size:255"`
	Instagram       string         `json:"instagram" gorm:"size:255"`
	Twitter         string         `json:"twitter" gorm:"size:255"`
	YouTube         string         `json:"youtube" gorm:"size:255"`
	Spotify         string         `json:"spotify" gorm:"size:255"`
	AppleMusic      string         `json:"apple_music" gorm:"size:255"`
	IsVerified      bool           `json:"is_verified" gorm:"default:false"`
	VerifiedAt      *time.Time     `json:"verified_at"`
	TotalPlays      int64          `json:"total_plays" gorm:"default:0"`
	TotalFollowers  int64          `json:"total_followers" gorm:"default:0"`
	TotalWorks      int64          `json:"total_works" gorm:"default:0"`
	Balance         float64        `json:"balance" gorm:"default:0"`
	FrozenBalance   float64        `json:"frozen_balance" gorm:"default:0"`
	BankAccount     string         `json:"bank_account" gorm:"size:50"`
	BankName        string         `json:"bank_name" gorm:"size:50"`
	BankHolder      string         `json:"bank_holder" gorm:"size:50"`
	AlipayAccount   string         `json:"alipay_account" gorm:"size:100"`
	WechatAccount   string         `json:"wechat_account" gorm:"size:100"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ArtistInfo) TableName() string {
	return "artist_infos"
}

type WorkStatus int

const (
	WorkStatusDraft    WorkStatus = 0
	WorkStatusReview   WorkStatus = 1
	WorkStatusPublished WorkStatus = 2
	WorkStatusRejected WorkStatus = 3
	WorkStatusOffline  WorkStatus = 4
)

type WorkType string

const (
	WorkTypeSingle WorkType = "single"
	WorkTypeAlbum  WorkType = "album"
	WorkTypeEP     WorkType = "ep"
)

type Work struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"index;not null"`
	ArtistID      uint           `json:"artist_id" gorm:"index;not null"`
	Title         string         `json:"title" gorm:"size:200;not null"`
	ArtistName    string         `json:"artist_name" gorm:"size:100"`
	AlbumID       *uint          `json:"album_id" gorm:"index"`
	Type          WorkType       `json:"type" gorm:"size:20;default:single"`
	Genre         string         `json:"genre" gorm:"size:100"`
	SubGenre      string         `json:"sub_genre" gorm:"size:100"`
	Language      string         `json:"language" gorm:"size:50"`
	Duration      int            `json:"duration"`
	Description   string         `json:"description" gorm:"type:text"`
	Lyrics        string         `json:"lyrics" gorm:"type:text"`
	Composer      string         `json:"composer" gorm:"size:100"`
	Lyricist      string         `json:"lyricist" gorm:"size:100"`
	Arranger      string         `json:"arranger" gorm:"size:100"`
	Producer      string         `json:"producer" gorm:"size:100"`
	CoverURL      string         `json:"cover_url" gorm:"size:255"`
	AudioURL      string         `json:"audio_url" gorm:"size:255"`
	AudioFormat   string         `json:"audio_format" gorm:"size:20"`
	AudioQuality  string         `json:"audio_quality" gorm:"size:20"`
	Bitrate       int            `json:"bitrate"`
	SampleRate    int            `json:"sample_rate"`
	FileSize      int64          `json:"file_size"`
	MD5           string         `json:"md5" gorm:"size:32"`
	Status        WorkStatus     `json:"status" gorm:"default:0"`
	PublishedAt   *time.Time     `json:"published_at"`
	PlayCount     int64          `json:"play_count" gorm:"default:0"`
	LikeCount     int64          `json:"like_count" gorm:"default:0"`
	CommentCount  int64          `json:"comment_count" gorm:"default:0"`
	ShareCount    int64          `json:"share_count" gorm:"default:0"`
	IsPublic      bool           `json:"is_public" gorm:"default:true"`
	Explicit      bool           `json:"explicit" gorm:"default:false"`
	Tags          []Tag          `json:"tags,omitempty" gorm:"many2many:work_tags"`
	Copyright     *Copyright     `json:"copyright,omitempty" gorm:"foreignKey:WorkID"`
	Comments      []Comment      `json:"-" gorm:"foreignKey:WorkID"`
	PlayRecords   []PlayRecord   `json:"-" gorm:"foreignKey:WorkID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Work) TableName() string {
	return "works"
}

type Album struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id" gorm:"index;not null"`
	ArtistID       uint           `json:"artist_id" gorm:"index;not null"`
	Title          string         `json:"title" gorm:"size:200;not null"`
	ArtistName     string         `json:"artist_name" gorm:"size:100"`
	Description    string         `json:"description" gorm:"type:text"`
	CoverURL       string         `json:"cover_url" gorm:"size:255"`
	ReleaseDate    time.Time      `json:"release_date"`
	Type           WorkType       `json:"type" gorm:"size:20;default:album"`
	Genre          string         `json:"genre" gorm:"size:100"`
	RecordLabel    string         `json:"record_label" gorm:"size:100"`
	Upc            string         `json:"upc" gorm:"size:50"`
	Status         WorkStatus     `json:"status" gorm:"default:0"`
	PublishedAt    *time.Time     `json:"published_at"`
	PlayCount      int64          `json:"play_count" gorm:"default:0"`
	LikeCount      int64          `json:"like_count" gorm:"default:0"`
	CommentCount   int64          `json:"comment_count" gorm:"default:0"`
	WorkCount      int            `json:"work_count" gorm:"default:0"`
	Works          []Work         `json:"works,omitempty" gorm:"foreignKey:AlbumID"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Album) TableName() string {
	return "albums"
}

type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Type      string         `json:"type" gorm:"size:50"`
	UsageCount int64         `json:"usage_count" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Tag) TableName() string {
	return "tags"
}

type WorkTag struct {
	WorkID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

func (WorkTag) TableName() string {
	return "work_tags"
}

type Copyright struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	WorkID        uint           `json:"work_id" gorm:"uniqueIndex;not null"`
	CopyrightType string         `json:"copyright_type" gorm:"size:50"`
	Owner         string         `json:"owner" gorm:"size:100"`
	LicenseType   string         `json:"license_type" gorm:"size:50"`
	RoyaltiesRate float64        `json:"royalties_rate" gorm:"default:0"`
	StartDate     time.Time      `json:"start_date"`
	EndDate       *time.Time     `json:"end_date"`
	IsExclusive   bool           `json:"is_exclusive" gorm:"default:false"`
	Territory     string         `json:"territory" gorm:"size:255"`
	ContractNo    string         `json:"contract_no" gorm:"size:50"`
	ContractFile  string         `json:"contract_file" gorm:"size:255"`
	Notes         string         `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (Copyright) TableName() string {
	return "copyrights"
}
