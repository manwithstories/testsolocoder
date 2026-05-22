package model

import (
	"time"

	"gorm.io/gorm"
)

type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"
	MemberStatusInactive MemberStatus = "inactive"
)

type Member struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MemberNo  string         `gorm:"size:32;not null;uniqueIndex" json:"member_no"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Phone     string         `gorm:"size:20;not null;uniqueIndex" json:"phone"`
	Email     string         `gorm:"size:100;index" json:"email"`
	IDCard    string         `gorm:"size:20;uniqueIndex" json:"id_card"`
	LevelID   uint           `gorm:"not null;index" json:"level_id"`
	Points    int            `gorm:"not null;default:0" json:"points"`
	Balance   float64        `gorm:"type:decimal(10,2);not null;default:0" json:"balance"`
	Status    MemberStatus   `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Level     *MemberLevel   `gorm:"foreignKey:LevelID" json:"level,omitempty"`
	Bookings  []Booking      `gorm:"foreignKey:MemberID" json:"bookings,omitempty"`
}
