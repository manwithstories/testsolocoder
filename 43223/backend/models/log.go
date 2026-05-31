package models

import (
	"time"

	"gorm.io/gorm"
)

type OperationLog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index"`
	Username    string         `json:"username" gorm:"size:50"`
	Action      string         `json:"action" gorm:"size:100;index"`
	Resource    string         `json:"resource" gorm:"size:100"`
	ResourceID  string         `json:"resource_id" gorm:"size:100"`
	Method      string         `json:"method" gorm:"size:20"`
	Path        string         `json:"path" gorm:"size:500"`
	IP          string         `json:"ip" gorm:"size:50"`
	UserAgent   string         `json:"user_agent" gorm:"size:500"`
	StatusCode  int            `json:"status_code"`
	RequestData string         `json:"request_data" gorm:"type:text"`
	ResponseData string        `json:"response_data" gorm:"type:text"`
	ErrorMsg    string         `json:"error_msg" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
