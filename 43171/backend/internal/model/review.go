package model

import (
	"time"

	"gorm.io/gorm"
)

type ReviewType string

const (
	ReviewTypeRental  ReviewType = "rental"
	ReviewTypeService   ReviewType = "service"
)

type Review struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Type        ReviewType     `gorm:"size:16;not null" json:"type"`
	OrderID     *uint          `gorm:"index" json:"order_id"`
	Order       *RentalOrder   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ServiceID   *uint          `gorm:"index" json:"service_id"`
	Service     *AerialService `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	ReviewerID  uint           `gorm:"index;not null" json:"reviewer_id"`
	Reviewer   *User          `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	RevieweeID uint           `gorm:"index;not null" json:"reviewee_id"`
	Reviewee   *User          `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
	DroneID   *uint          `gorm:"index" json:"drone_id"`
	Drone     *Drone         `gorm:"foreignKey:DroneID" json:"drone,omitempty"`
	Rating      int            `gorm:"type:decimal(2,1)" json:"rating"`
	Content     string         `gorm:"size:1000" json:"content"`
	Images       string         `gorm:"size:1024" json:"images"`
	Reply        string         `gorm:"size:500" json:"reply"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
