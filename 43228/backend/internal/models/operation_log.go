package models

import (
	"time"

	"gorm.io/gorm"
)

type OperationLog struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	Username    string         `gorm:"size:64" json:"username"`
	Operation   string         `gorm:"size:128;not null" json:"operation"`
	Method      string         `gorm:"size:16" json:"method"`
	Path        string         `gorm:"size:512" json:"path"`
	IP          string         `gorm:"size:64" json:"ip"`
	UserAgent   string         `gorm:"size:512" json:"user_agent"`
	Status      int            `json:"status"`
	Detail      string         `gorm:"type:text" json:"detail"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
