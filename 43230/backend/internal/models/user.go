package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleDesigner UserRole = "designer"
	RolePrinter  UserRole = "printer"
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Email           string     `json:"email" gorm:"uniqueIndex;not null"`
	Username        string     `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash    string     `json:"-" gorm:"not null"`
	Role            UserRole   `json:"role" gorm:"not null"`
	Status          UserStatus `json:"status" gorm:"default:active"`
	Phone           string     `json:"phone"`
	Avatar          string     `json:"avatar"`
	CreditScore     float64    `json:"credit_score" gorm:"default:5.0"`
	RealName        string     `json:"real_name"`
	IDCardNumber    string     `json:"-"`
	IDCardVerified  bool       `json:"id_card_verified" gorm:"default:false"`
	Balance         float64    `json:"balance" gorm:"default:0"`
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LastLoginIP     string     `json:"last_login_ip"`
	DesignerProfile *DesignerProfile `json:"designer_profile,omitempty" gorm:"foreignKey:UserID"`
	PrinterProfile  *PrinterProfile  `json:"printer_profile,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type DesignerProfile struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	Nickname        string    `json:"nickname"`
	Bio             string    `json:"bio"`
	PortfolioURL    string    `json:"portfolio_url"`
	TotalModels     int       `json:"total_models" gorm:"default:0"`
	TotalSales      float64   `json:"total_sales" gorm:"default:0"`
	Rating          float64   `json:"rating" gorm:"default:5.0"`
	RatingCount     int       `json:"rating_count" gorm:"default:0"`
	Specialties     []string  `json:"specialties" gorm:"type:text[]"`
	ExperienceYears int       `json:"experience_years"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PrinterProfile struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	CompanyName   string    `json:"company_name"`
	Address       string    `json:"address"`
	BusinessLicense string  `json:"business_license"`
	LicenseVerified bool    `json:"license_verified" gorm:"default:false"`
	MaxPrintSize  string    `json:"max_print_size"`
	SupportedMaterials []string `json:"supported_materials" gorm:"type:text[]"`
	Rating        float64   `json:"rating" gorm:"default:5.0"`
	RatingCount   int       `json:"rating_count" gorm:"default:0"`
	TotalOrders   int       `json:"total_orders" gorm:"default:0"`
	TotalRevenue  float64   `json:"total_revenue" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
