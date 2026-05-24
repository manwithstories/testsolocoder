package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	RoleID    uint           `json:"role_id" gorm:"default:3"`
	Role      *Role          `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Status    int            `json:"status" gorm:"default:1;comment:1=active,2=disabled"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Role struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"size:20;uniqueIndex;not null;comment:admin,editor,viewer"`
	Description string `json:"description" gorm:"size:100"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
