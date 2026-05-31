package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      *uuid.UUID     `gorm:"type:uuid;index" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action      string         `gorm:"size:100;not null;index" json:"action"`
	EntityType  string         `gorm:"size:50;not null;index" json:"entity_type"`
	EntityID    *uuid.UUID     `gorm:"type:uuid;index" json:"entity_id"`
	OldValue    string         `gorm:"type:text" json:"old_value"`
	NewValue    string         `gorm:"type:text" json:"new_value"`
	IPAddress   string         `gorm:"size:50" json:"ip_address"`
	UserAgent   string         `gorm:"size:500" json:"user_agent"`
	RequestID   string         `gorm:"size:100" json:"request_id"`
	Status      string         `gorm:"type:varchar(20);default:success" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type CreateAuditLogRequest struct {
	UserID     *uuid.UUID `json:"user_id"`
	Action     string     `json:"action" binding:"required"`
	EntityType string     `json:"entity_type" binding:"required"`
	EntityID   *uuid.UUID `json:"entity_id"`
	OldValue   string     `json:"old_value"`
	NewValue   string     `json:"new_value"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	RequestID  string     `json:"request_id"`
	Status     string     `json:"status"`
}
