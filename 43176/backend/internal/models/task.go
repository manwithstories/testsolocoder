package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskType string

const (
	TaskTypeBuy      TaskType = "buy"
	TaskTypePickup   TaskType = "pickup"
	TaskTypeDeliver  TaskType = "deliver"
	TaskTypeQueue    TaskType = "queue"
	TaskTypeErrand   TaskType = "errand"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAccepted   TaskStatus = "accepted"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusTimeout    TaskStatus = "timeout"
)

type Task struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PublisherID   uint           `gorm:"index" json:"publisher_id"`
	Publisher     User           `json:"publisher,omitempty"`
	CourierID     *uint          `gorm:"index" json:"courier_id,omitempty"`
	Courier       *User          `json:"courier,omitempty"`
	Type          TaskType       `gorm:"size:20" json:"type"`
	Title         string         `gorm:"size:100" json:"title"`
	Description   string         `gorm:"size:1000" json:"description"`
	StartAddr     string         `gorm:"size:255" json:"start_addr"`
	StartLat      float64        `json:"start_lat"`
	StartLng      float64        `json:"start_lng"`
	EndAddr       string         `gorm:"size:255" json:"end_addr"`
	EndLat        float64        `json:"end_lat"`
	EndLng        float64        `json:"end_lng"`
	PublishTime   time.Time      `json:"publish_time"`
	Deadline      time.Time      `json:"deadline"`
	Reward        float64        `json:"reward"`
	Status        TaskStatus     `gorm:"size:20;default:pending" json:"status"`
	Images        []TaskImage    `gorm:"foreignKey:TaskID" json:"images,omitempty"`
	Order         *Order         `gorm:"foreignKey:TaskID" json:"order,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type TaskImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"index" json:"task_id"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskAcceptLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"index" json:"task_id"`
	CourierID uint      `gorm:"index" json:"courier_id"`
	Accepted  bool      `gorm:"default:false" json:"accepted"`
	CreatedAt time.Time `json:"created_at"`
}
