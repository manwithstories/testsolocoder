package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
	MessageTypeSystem MessageType = "system"
)

type Message struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SenderID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"sender_id"`
	Sender     User           `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	ReceiverID uuid.UUID      `gorm:"type:uuid;index;not null" json:"receiver_id"`
	Receiver   User           `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Type       MessageType    `gorm:"type:varchar(20);default:'text'" json:"type"`
	Content    string         `gorm:"type:text" json:"content"`
	Encrypted  bool           `gorm:"default:false" json:"encrypted"`
	FileURL    string         `gorm:"size:500" json:"file_url"`
	FileName   string         `gorm:"size:200" json:"file_name"`
	FileSize   int64          `json:"file_size"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	ReadAt     *time.Time     `json:"read_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type Conversation struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Participant1 uuid.UUID      `gorm:"type:uuid;index;not null" json:"participant1_id"`
	Participant2 uuid.UUID      `gorm:"type:uuid;index;not null" json:"participant2_id"`
	LastMessage  string         `gorm:"type:text" json:"last_message"`
	LastMessageAt *time.Time   `json:"last_message_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
