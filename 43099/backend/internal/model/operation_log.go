package model

import (
	"time"
)

type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	Module    string    `gorm:"size:50" json:"module"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Details   string    `gorm:"type:json" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
