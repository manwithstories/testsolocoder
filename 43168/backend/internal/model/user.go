package model

import (
	"time"
)

// 角色常量
const (
	RoleManufacturer = "manufacturer" // 厂商
	RoleDesigner     = "designer"     // 设计师
	RoleOwner        = "owner"        // 业主
)

// 状态常量
const (
	StatusActive   = 1 // 启用
	StatusInactive = 0 // 禁用
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      string    `gorm:"size:32;not null;default:owner" json:"role"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Email     string    `gorm:"size:128" json:"email"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsActive 判断用户是否启用
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsManufacturer 判断是否为厂商角色
func (u *User) IsManufacturer() bool {
	return u.Role == RoleManufacturer
}

// IsDesigner 判断是否为设计师角色
func (u *User) IsDesigner() bool {
	return u.Role == RoleDesigner
}

// IsOwner 判断是否为业主角色
func (u *User) IsOwner() bool {
	return u.Role == RoleOwner
}

// ValidRole 校验角色是否合法
func ValidRole(role string) bool {
	switch role {
	case RoleManufacturer, RoleDesigner, RoleOwner:
		return true
	default:
		return false
	}
}
