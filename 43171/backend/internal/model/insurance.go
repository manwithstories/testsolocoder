package model

import (
	"time"

	"gorm.io/gorm"
)

type ClaimStatus string

const (
	ClaimStatusPending  ClaimStatus = "pending"
	ClaimStatusReviewing ClaimStatus = "reviewing"
	ClaimStatusApproved ClaimStatus = "approved"
	ClaimStatusRejected ClaimStatus = "rejected"
)

type InsuranceClaim struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ClaimNo       string         `gorm:"size:32;uniqueIndex;not null" json:"claim_no"`
	OrderID         uint           `gorm:"index;not null" json:"order_id"`
	Order           *RentalOrder   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	User            *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DroneID        *uint          `gorm:"index" json:"drone_id"`
	Drone           *Drone         `gorm:"foreignKey:DroneID" json:"drone,omitempty"`
	DamageDesc     string         `gorm:"size:1000" json:"damage_desc"`
	DamageImages     string         `gorm:"size:1024" json:"damage_images"`
	EstimatedCost  float64        `json:"estimated_cost"`
	ActualCost     float64        `json:"actual_cost"`
	Status         ClaimStatus    `gorm:"size:16;default:pending;index" json:"status"`
	ReviewerID     *uint          `gorm:"index" json:"reviewer_id"`
	Reviewer       *User          `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	ReviewRemark  string         `gorm:"size:500" json:"review_remark"`
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	DeductedAmount float64      `json:"deducted_amount"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
