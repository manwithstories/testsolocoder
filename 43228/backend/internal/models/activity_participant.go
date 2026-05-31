package models

import (
	"time"
)

type ParticipantStatus string

const (
	ParticipantStatusPending   ParticipantStatus = "pending"
	ParticipantStatusConfirmed ParticipantStatus = "confirmed"
	ParticipantStatusCancelled ParticipantStatus = "cancelled"
)

type ActivityParticipant struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	ActivityID uint              `gorm:"index;not null" json:"activity_id"`
	UserID     uint              `gorm:"index;not null" json:"user_id"`
	Status     ParticipantStatus `gorm:"size:32;not null;default:pending;index" json:"status"`
	JoinedAt   time.Time         `json:"joined_at"`

	Activity *TastingActivity `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
	User     *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ActivityParticipant) TableName() string {
	return "activity_participants"
}
