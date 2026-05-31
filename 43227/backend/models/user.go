package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Username   string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password   string         `gorm:"size:255;not null" json:"-"`
	Email      string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone      string         `gorm:"size:20" json:"phone"`
	Role       string         `gorm:"size:20;not null" json:"role"`
	Avatar     string         `gorm:"size:255" json:"avatar"`
	Reputation float64        `gorm:"type:decimal(3,2);default:5.00" json:"reputation"`
	Address    string         `gorm:"type:text" json:"address"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
