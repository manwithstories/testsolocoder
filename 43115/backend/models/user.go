package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleCustomer         UserRole = "customer"
	RoleServiceProvider  UserRole = "service_provider"
	RoleAdmin            UserRole = "admin"
)

type ProviderStatus string

const (
	ProviderStatusPending  ProviderStatus = "pending"
	ProviderStatusApproved ProviderStatus = "approved"
	ProviderStatusRejected ProviderStatus = "rejected"
	ProviderStatusSuspended ProviderStatus = "suspended"
)

type User struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Phone           string         `json:"phone" gorm:"uniqueIndex;size:20;not null"`
	Password        string         `json:"-" gorm:"not null"`
	Nickname        string         `json:"nickname" gorm:"size:50"`
	Avatar          string         `json:"avatar" gorm:"size:255"`
	Role            UserRole       `json:"role" gorm:"size:20;not null;default:customer"`
	RealName        string         `json:"real_name" gorm:"size:50"`
	IDCard          string         `json:"id_card" gorm:"size:20"`
	Gender          string         `json:"gender" gorm:"size:10"`
	Age             int            `json:"age" gorm:"default:0"`
	ProviderStatus  ProviderStatus `json:"provider_status" gorm:"size:20;default:pending"`
	CertificationDesc string       `json:"certification_desc" gorm:"type:text"`
	RejectReason    string         `json:"reject_reason" gorm:"size:500"`
	Rating          float64        `json:"rating" gorm:"default:5.0"`
	OrderCount      int            `json:"order_count" gorm:"default:0"`
	Balance         float64        `json:"balance" gorm:"default:0"`
	TotalIncome     float64        `json:"total_income" gorm:"default:0"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	WechatOpenID    string         `json:"-" gorm:"size:100"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Addresses       []Address             `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	Certifications  []ServiceProviderCert `json:"certifications,omitempty" gorm:"foreignKey:ProviderID"`
}

type Address struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	ContactName string         `json:"contact_name" gorm:"size:50;not null"`
	ContactPhone string        `json:"contact_phone" gorm:"size:20;not null"`
	Province    string         `json:"province" gorm:"size:50;not null"`
	City        string         `json:"city" gorm:"size:50;not null"`
	District    string         `json:"district" gorm:"size:50;not null"`
	Address     string         `json:"address" gorm:"size:255;not null"`
	Longitude   float64        `json:"longitude" gorm:"default:0"`
	Latitude    float64        `json:"latitude" gorm:"default:0"`
	IsDefault   bool           `json:"is_default" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ServiceProviderCert struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ProviderID   uint           `json:"provider_id" gorm:"not null;index"`
	CertType     string         `json:"cert_type" gorm:"size:50;not null"`
	CertName     string         `json:"cert_name" gorm:"size:100;not null"`
	CertNumber   string         `json:"cert_number" gorm:"size:50"`
	CertImage    string         `json:"cert_image" gorm:"size:255"`
	IssuedAt     *time.Time     `json:"issued_at"`
	ExpiredAt    *time.Time     `json:"expired_at"`
	IsVerified   bool           `json:"is_verified" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
