package model

import (
	"time"

	"gorm.io/gorm"
)

type MemberLevel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;uniqueIndex" json:"name"`
	DiscountRate float64       `gorm:"type:decimal(5,2);not null;default:1.0" json:"discount_rate"`
	PointsRate  float64        `gorm:"type:decimal(5,2);not null;default:1.0" json:"points_rate"`
	MinPoints   int            `gorm:"not null;default:0;index" json:"min_points"`
	MaxPoints   int            `gorm:"not null;default:0;index" json:"max_points"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Members     []Member       `gorm:"foreignKey:LevelID" json:"members,omitempty"`
}
