package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	PlanID     uuid.UUID `gorm:"type:uuid;index" json:"plan_id"`
	ActivityID uuid.UUID `gorm:"type:uuid;index" json:"activity_id"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`
	Resource   string    `gorm:"type:varchar(50);not null" json:"resource"`
	ResourceID uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	OldValue   string    `gorm:"type:text" json:"old_value"`
	NewValue   string    `gorm:"type:text" json:"new_value"`
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}
