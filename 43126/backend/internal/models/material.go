package models

import (
	"time"

	"gorm.io/gorm"
)

type MeetingMaterial struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	BookingID   uint           `json:"booking_id" gorm:"index;not null"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	FileName    string         `json:"file_name" gorm:"size:255;not null"`
	FilePath    string         `json:"file_path" gorm:"size:500;not null"`
	FileSize    int64          `json:"file_size"`
	FileType    string         `json:"file_type" gorm:"size:50"`
	MeetingDate time.Time      `json:"meeting_date" gorm:"index"`
	ExpireAt    *time.Time     `json:"expire_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (MeetingMaterial) TableName() string {
	return "meeting_materials"
}

type UsageStats struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoomID       uint      `json:"room_id" gorm:"index;not null"`
	Date         time.Time `json:"date" gorm:"index;not null"`
	TotalHours   float64   `json:"total_hours" gorm:"type:decimal(10,2)"`
	Bookings     int       `json:"bookings"`
	Revenue      float64   `json:"revenue" gorm:"type:decimal(10,2)"`
	Department   string    `json:"department" gorm:"size:50;index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (UsageStats) TableName() string {
	return "usage_stats"
}
