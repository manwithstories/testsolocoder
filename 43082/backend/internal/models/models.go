package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleBuyer  = "buyer"
	RoleSeller = "seller"
	RoleAdmin  = "admin"

	ShopStatusPending  = "pending"
	ShopStatusApproved = "approved"
	ShopStatusRejected = "rejected"
	ShopStatusSuspended = "suspended"

	ProductStatusDraft     = "draft"
	ProductStatusOnSale    = "on_sale"
	ProductStatusOffSale   = "off_sale"
	ProductStatusRejected  = "rejected"

	OrderStatusPendingPayment = "pending_payment"
	OrderStatusPendingShip    = "pending_ship"
	OrderStatusShipped        = "shipped"
	OrderStatusCompleted      = "completed"
	OrderStatusRefunded       = "refunded"
	OrderStatusCancelled      = "cancelled"

	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusFailed  = "failed"

	RefundStatusPending  = "pending"
	RefundStatusApproved = "approved"
	RefundStatusRejected = "rejected"
	RefundStatusCompleted = "completed"

	NotificationTypeOrder       = "order"
	NotificationTypeShipment    = "shipment"
	NotificationTypeRefund      = "refund"
	NotificationTypeShopReview  = "shop_review"

	DisputeStatusOpen     = "open"
	DisputeStatusProcessing = "processing"
	DisputeStatusResolved = "resolved"
	DisputeStatusClosed   = "closed"
)

type User struct {
	gorm.Model
	Username string  `gorm:"size:50;uniqueIndex;not null"`
	Email    string  `gorm:"size:100;uniqueIndex;not null"`
	Phone    string  `gorm:"size:20;uniqueIndex"`
	Password string  `gorm:"size:255;not null"`
	Role     string  `gorm:"size:20;not null;default:'buyer'"`
	Avatar   string  `gorm:"size:255"`
	Nickname string  `gorm:"size:50"`
	Status   string  `gorm:"size:20;default:'active'"`
}

type Shop struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index"`
	Name         string `gorm:"size:100;uniqueIndex;not null"`
	Description  string `gorm:"type:text"`
	Logo         string `gorm:"size:255"`
	ContactName  string `gorm:"size:50"`
	ContactPhone string `gorm:"size:20"`
	Address      string `gorm:"size:255"`
	IDCardFront  string `gorm:"size:255"`
	IDCardBack   string `gorm:"size:255"`
	BusinessLicense string `gorm:"size:255"`
	Status       string `gorm:"size:20;default:'pending'"`
	RejectReason string `gorm:"size:255"`
	ApprovedAt   *time.Time
	Rating       float64 `gorm:"default:5"`
	User         User
}

type Category struct {
	gorm.Model
	Name       string `gorm:"size:50;not null"`
	Icon       string `gorm:"size:255"`
	ParentID   *uint
	Level      int    `gorm:"default:1"`
	Sort       int    `gorm:"default:0"`
	Status     string `gorm:"size:20;default:'active'"`
	Products   []Product
}

type Product struct {
	gorm.Model
	ShopID      uint   `gorm:"not null;index"`
	CategoryID  uint   `gorm:"not null;index"`
	Name        string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`
	MainImage   string `gorm:"size:255"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Stock       int    `gorm:"not null;default:0"`
	Sales       int    `gorm:"default:0"`
	Status      string `gorm:"size:20;default:'draft'"`
	Weight      float64 `gorm:"type:decimal(10,2)"`
	IsHot       bool   `gorm:"default:false"`
	IsRecommend bool   `gorm:"default:false"`
	Shop        Shop
	Category    Category
	Images      []ProductImage
	Specs       []ProductSpec
	Skus        []SKU
}

type ProductImage struct {
	gorm.Model
	ProductID uint   `gorm:"not null;index"`
	URL       string `gorm:"size:255;not null"`
	Sort      int    `gorm:"default:0"`
}

type ProductSpec struct {
	gorm.Model
	ProductID uint   `gorm:"not null;index"`
	Name      string `gorm:"size:50;not null"`
	Values    string `gorm:"type:jsonb;not null"`
}

type SKU struct {
	gorm.Model
	ProductID uint    `gorm:"not null;index"`
	Specs     string  `gorm:"type:jsonb;not null"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
	Stock     int     `gorm:"not null;default:0"`
	SKUCode   string  `gorm:"size:50;uniqueIndex"`
}

type CartItem struct {
	gorm.Model
	UserID    uint `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`
	SKU_ID    *uint
	Quantity  int `gorm:"not null;default:1"`
	Product   Product
}

type Order struct {
	gorm.Model
	OrderNo     string  `gorm:"size:50;uniqueIndex;not null"`
	UserID      uint    `gorm:"not null;index"`
	ShopID      uint    `gorm:"not null;index"`
	TotalAmount float64 `gorm:"type:decimal(10,2);not null"`
	ShippingFee float64 `gorm:"type:decimal(10,2);default:0"`
	Status      string  `gorm:"size:30;not null;default:'pending_payment'"`
	ReceiverName string `gorm:"size:50;not null"`
	ReceiverPhone string `gorm:"size:20;not null"`
	ReceiverAddress string `gorm:"size:255;not null"`
	Remark      string `gorm:"size:500"`
	PaidAt      *time.Time
	ShippedAt   *time.Time
	CompletedAt *time.Time
	TrackingNo  string `gorm:"size:50"`
	TrackingCompany string `gorm:"size:50"`
	User        User
	Shop        Shop
	Items       []OrderItem
	Payment     *Payment
	Refunds     []Refund
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null;index"`
	SKU_ID    *uint
	ProductName string `gorm:"size:200;not null"`
	ProductImage string `gorm:"size:255"`
	Specs     string  `gorm:"type:jsonb"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
	Quantity  int     `gorm:"not null"`
	Subtotal  float64 `gorm:"type:decimal(10,2);not null"`
	Reviewed  bool    `gorm:"default:false"`
}

type Payment struct {
	gorm.Model
	OrderID     uint    `gorm:"not null;index"`
	PaymentNo   string  `gorm:"size:50;uniqueIndex;not null"`
	Amount      float64 `gorm:"type:decimal(10,2);not null"`
	Method      string  `gorm:"size:20"`
	Status      string  `gorm:"size:20;default:'pending'"`
	PaidAt      *time.Time
	TransactionID string `gorm:"size:100"`
}

type Refund struct {
	gorm.Model
	OrderID     uint    `gorm:"not null;index"`
	OrderItemID *uint
	UserID      uint    `gorm:"not null;index"`
	ShopID      uint    `gorm:"not null;index"`
	RefundNo    string  `gorm:"size:50;uniqueIndex;not null"`
	Amount      float64 `gorm:"type:decimal(10,2);not null"`
	Reason      string  `gorm:"size:500;not null"`
	Status      string  `gorm:"size:20;default:'pending'"`
	Type        string  `gorm:"size:20"`
	RejectReason string `gorm:"size:255"`
	ApprovedAt  *time.Time
	CompletedAt *time.Time
}

type Review struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`
	ProductID uint   `gorm:"not null;index"`
	OrderID   uint   `gorm:"not null;index"`
	OrderItemID uint `gorm:"not null;index"`
	ShopID    uint   `gorm:"not null;index"`
	Rating    int    `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	Images    string `gorm:"type:jsonb"`
	Reply     string `gorm:"type:text"`
	ReplyAt   *time.Time
	User      User
	Product   Product
}

type Favorite struct {
	gorm.Model
	UserID    uint  `gorm:"not null;index"`
	ShopID    *uint `gorm:"index"`
	ProductID *uint `gorm:"index"`
}

type Notification struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index"`
	Type    string `gorm:"size:20;not null"`
	Title   string `gorm:"size:100;not null"`
	Content string `gorm:"type:text;not null"`
	Data    string `gorm:"type:jsonb"`
	IsRead  bool   `gorm:"default:false"`
}

type Dispute struct {
	gorm.Model
	OrderID    uint   `gorm:"not null;index"`
	UserID     uint   `gorm:"not null;index"`
	ShopID     uint   `gorm:"not null;index"`
	Type       string `gorm:"size:20;not null"`
	Reason     string `gorm:"size:500;not null"`
	Evidence   string `gorm:"type:jsonb"`
	Status     string `gorm:"size:20;default:'open'"`
	Result     string `gorm:"type:text"`
	AdminID    *uint
	ResolvedAt *time.Time
	User       User
	Shop       Shop
}
