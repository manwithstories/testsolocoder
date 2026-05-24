package models

import (
	"time"
)

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Username      string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password      string    `json:"-" gorm:"size:255;not null"`
	RealName      string    `json:"realName" gorm:"size:50"`
	IDCard        string    `json:"-" gorm:"size:18"`
	Role          string    `json:"role" gorm:"size:20;default:'renter"`
	Verified      bool      `json:"verified" gorm:"default:false"`
	Phone         string    `json:"phone" gorm:"size:20"`
	Avatar        string    `json:"avatar" gorm:"size:255"`
	DepositBalance float64   `json:"depositBalance" gorm:"default:0"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Equipment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OwnerID     uint      `json:"ownerId" gorm:"index;not null"`
	Owner       User      `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Category    string    `json:"category" gorm:"size:50;not null"`
	Brand       string    `json:"brand" gorm:"size:50;not null"`
	Model       string    `json:"model" gorm:"size:100;not null"`
	PurchaseDate string   `json:"purchaseDate" gorm:"size:20"`
	Status      string    `json:"status" gorm:"size:20;default:'available'"`
	Deposit     float64   `json:"deposit" gorm:"not null"`
	DailyRent   float64   `json:"dailyRent" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Rating      float64   `json:"rating" gorm:"default:0"`
	ReviewCount int       `json:"reviewCount" gorm:"default:0"`
	Images      []EquipmentImage `json:"images,omitempty" gorm:"foreignKey:EquipmentID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EquipmentImage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EquipmentID uint      `json:"equipmentId" gorm:"index;not null"`
	ImageURL    string    `json:"imageUrl" gorm:"size:255;not null"`
	SortOrder   int       `json:"sortOrder" gorm:"default:0"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Order struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	EquipmentID   uint      `json:"equipmentId" gorm:"index;not null"`
	Equipment     Equipment `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	RenterID      uint      `json:"renterId" gorm:"index;not null"`
	Renter        User      `json:"renter,omitempty" gorm:"foreignKey:RenterID"`
	OwnerID       uint      `json:"ownerId" gorm:"index;not null"`
	Owner         User      `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	StartDate     string    `json:"startDate" gorm:"size:20;not null"`
	EndDate       string    `json:"endDate" gorm:"size:20;not null"`
	TotalRent     float64   `json:"totalRent" gorm:"not null"`
	Deposit        float64   `json:"deposit" gorm:"not null"`
	Status         string    `json:"status" gorm:"size:20;default:'pending'"`
	DeliveryMethod string   `json:"deliveryMethod" gorm:"size:50"`
	DeliveryAddress string  `json:"deliveryAddress" gorm:"size:255"`
	RejectReason  string  `json:"rejectReason,omitempty" gorm:"size:255"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"orderId" gorm:"index;not null"`
	FromUserID uint    `json:"fromUserId" gorm:"index;not null"`
	FromUser  *User     `json:"fromUser,omitempty" gorm:"foreignKey:FromUserID"`
	ToUserID   uint    `json:"toUserId" gorm:"index;not null"`
	ToUser    *User     `json:"toUser,omitempty" gorm:"foreignKey:ToUserID"`
	EquipmentID  uint   `json:"equipmentId" gorm:"index;not null"`
	Equipment *Equipment `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	Rating     int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Content    string    `json:"content" gorm:"type:text"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Settlement struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OrderID       uint      `json:"orderId" gorm:"uniqueIndex;not null"`
	Order          Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	TotalRent      float64   `json:"totalRent" gorm:"not null"`
	Deposit        float64   `json:"deposit" gorm:"not null"`
	RefundDeposit float64   `json:"refundDeposit" gorm:"default:0"`
	FinalAmount  float64   `json:"finalAmount" gorm:"not null"`
	DamageFee     float64   `json:"damageFee" gorm:"default:0"`
	Status        string    `json:"status" gorm:"size:20;default:'completed'"`
	Remark        string    `json:"remark,omitempty" gorm:"size:255"`
	SettledAt    time.Time `json:"settledAt"`
	CreatedAt    time.Time `json:"createdAt"`
}
