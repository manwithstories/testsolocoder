package models

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleCaptain Role = "captain"
	RolePlayer Role = "player"
	RoleReferee Role = "referee"
)

type User struct {
	gorm.Model
	Email      string  `gorm:"uniqueIndex;size:100" json:"email"`
	Password   string  `json:"-"`
	FullName   string  `gorm:"size:100" json:"full_name"`
	Role       Role    `gorm:"size:20;default:player" json:"role"`
	Phone      string  `gorm:"size:20" json:"phone"`
	TeamID     *uint   `json:"team_id,omitempty"`
	Team       *Team   `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	IsActive   bool    `gorm:"default:true" json:"is_active"`
}
