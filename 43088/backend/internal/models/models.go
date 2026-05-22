package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

type Podcast struct {
	BaseModel
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	FeedURL     string    `gorm:"uniqueIndex;not null" json:"feed_url"`
	Website     string    `json:"website"`
	Author      string    `json:"author"`
	CoverImage  string    `json:"cover_image"`
	Language    string    `json:"language"`
	Category    string    `json:"category"`
	LastChecked time.Time `json:"last_checked"`
	LastUpdated time.Time `json:"last_updated"`
	Episodes    []Episode `gorm:"foreignKey:PodcastID" json:"episodes,omitempty"`
}

type Episode struct {
	BaseModel
	PodcastID    uuid.UUID `gorm:"type:uuid;index" json:"podcast_id"`
	Title        string    `gorm:"not null" json:"title"`
	Description  string    `json:"description"`
	GUID         string    `gorm:"uniqueIndex;not null" json:"guid"`
	AudioURL     string    `json:"audio_url"`
	AudioType    string    `json:"audio_type"`
	Duration     int       `json:"duration"`
	PubDate      time.Time `json:"pub_date"`
	EpisodeType    string    `json:"episode_type"`
	SeasonNumber int      `json:"season_number"`
	EpisodeNumber int     `json:"episode_number"`
	IsNew        bool      `gorm:"default:true" json:"is_new"`
	Podcast      Podcast   `gorm:"foreignKey:PodcastID" json:"podcast,omitempty"`
}

type PlaybackProgress struct {
	BaseModel
	EpisodeID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_episode" json:"episode_id"`
	CurrentTime  float64   `json:"current_time"`
	Completed    bool      `gorm:"default:false" json:"completed"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	PlayCount    int       `gorm:"default:0" json:"play_count"`
	LastPlayedAt time.Time `json:"last_played_at"`
	Episode      Episode   `gorm:"foreignKey:EpisodeID" json:"episode,omitempty"`
}

type Note struct {
	BaseModel
	EpisodeID uuid.UUID `gorm:"type:uuid;index" json:"episode_id"`
	Timestamp  float64   `json:"timestamp"`
	Content    string    `gorm:"type:text" json:"content"`
	Tags       []string  `gorm:"type:text[]" json:"tags"`
	Episode    Episode   `gorm:"foreignKey:EpisodeID" json:"episode,omitempty"`
}

type Playlist struct {
	BaseModel
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
	Items       []PlaylistItem `gorm:"foreignKey:PlaylistID" json:"items,omitempty"`
}

type PlaylistItem struct {
	BaseModel
	PlaylistID uuid.UUID `gorm:"type:uuid;index" json:"playlist_id"`
	EpisodeID  uuid.UUID `gorm:"type:uuid;index" json:"episode_id"`
	Position   int       `json:"position"`
	Playlist   Playlist  `gorm:"foreignKey:PlaylistID" json:"playlist,omitempty"`
	Episode    Episode   `gorm:"foreignKey:EpisodeID" json:"episode,omitempty"`
}

type ListeningHistory struct {
	BaseModel
	EpisodeID    uuid.UUID `gorm:"type:uuid;index" json:"episode_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     int       `json:"duration"`
	Completion   float64   `json:"completion"`
	Episode      Episode   `gorm:"foreignKey:EpisodeID" json:"episode,omitempty"`
}

type PodcastStats struct {
	PodcastID    uuid.UUID `json:"podcast_id"`
	TotalEpisodes int64     `json:"total_episodes"`
	UnplayedCount int64     `json:"unplayed_count"`
	CompletedCount int64     `json:"completed_count"`
	TotalListened int64     `json:"total_listened_seconds"`
}

type ListeningStats struct {
	Date             string  `json:"date"`
	TotalDuration    int64   `json:"total_duration"`
	EpisodeCount     int64   `json:"episode_count"`
	CompletionRate   float64 `json:"completion_rate"`
}
