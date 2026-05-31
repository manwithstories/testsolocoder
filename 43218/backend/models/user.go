package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Username        string         `gorm:"uniqueIndex;size:50;not null" json:"username" binding:"required,min=3,max=50"`
	Password        string         `gorm:"size:255;not null" json:"-" binding:"required,min=6"`
	Email           string         `gorm:"uniqueIndex;size:100" json:"email" binding:"omitempty,email"`
	Phone           string         `gorm:"uniqueIndex;size:20" json:"phone" binding:"omitempty,len=11"`
	Nickname        string         `gorm:"size:50" json:"nickname"`
	Avatar          string         `gorm:"size:255" json:"avatar"`
	Role            string         `gorm:"size:20;not null;default:buyer" json:"role"`
	RealName        string         `gorm:"size:50" json:"real_name"`
	IDCard          string         `gorm:"size:20" json:"id_card"`
	IsAuthenticated bool           `gorm:"default:false" json:"is_authenticated"`
	CreditScore     int            `gorm:"default:100" json:"credit_score"`
	WalletBalance   float64        `gorm:"default:0" json:"wallet_balance"`
	Status          int            `gorm:"default:1" json:"status"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Products    []Product    `gorm:"foreignKey:SellerID" json:"-"`
	Orders      []Order      `gorm:"foreignKey:BuyerID" json:"-"`
	SellOrders  []Order      `gorm:"foreignKey:SellerID" json:"-"`
	RepairServices []RepairService `gorm:"foreignKey:TechnicianID" json:"-"`
	RepairOrders []RepairOrder `gorm:"foreignKey:TechnicianID" json:"-"`
	Reviews     []Review     `gorm:"foreignKey:ReviewerID" json:"-"`
}

func (User) TableName() string {
	return "users"
}

const (
	RoleSeller     = "seller"
	RoleBuyer      = "buyer"
	RoleTechnician = "technician"
	RoleAdmin      = "admin"
)

const (
	UserStatusNormal  = 1
	UserStatusFrozen  = 2
	UserStatusBanned  = 3
)

type TechnicianCert struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	CertType      string         `gorm:"size:50;not null" json:"cert_type"`
	CertNumber    string         `gorm:"size:50;not null" json:"cert_number"`
	CertImage     string         `gorm:"size:255" json:"cert_image"`
	Status        int            `gorm:"default:0" json:"status"`
	RejectReason  string         `gorm:"size:255" json:"reject_reason"`
	ReviewedBy    *uint          `json:"reviewed_by"`
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TechnicianCert) TableName() string {
	return "technician_certs"
}

const (
	CertStatusPending  = 0
	CertStatusApproved = 1
	CertStatusRejected = 2
)
