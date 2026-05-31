package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	SellerID      uint           `gorm:"index;not null" json:"seller_id"`
	Title         string         `gorm:"size:200;not null" json:"title" binding:"required,min=2,max=200"`
	Category      string         `gorm:"size:50;not null;index" json:"category" binding:"required"`
	Brand         string         `gorm:"size:100;not null" json:"brand" binding:"required,max=100"`
	Model         string         `gorm:"size:100;not null" json:"model" binding:"required,max=100"`
	Condition     string         `gorm:"size:20;not null;index" json:"condition" binding:"required,oneof=全新 95新 9成新 8成新 7成新及以下"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price" binding:"required,min=0"`
	OriginalPrice float64        `gorm:"type:decimal(10,2)" json:"original_price"`
	Description   string         `gorm:"type:text" json:"description" binding:"required"`
	WarrantyDays  int            `gorm:"default:30" json:"warranty_days"`
	Images        string         `gorm:"type:text" json:"images"`
	Status        int            `gorm:"default:0;index" json:"status"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	FavoriteCount int            `gorm:"default:0" json:"favorite_count"`
	SoldCount     int            `gorm:"default:0" json:"sold_count"`
	RejectReason  string         `gorm:"size:255" json:"reject_reason"`
	ReviewedBy    *uint          `json:"reviewed_by"`
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Seller  User   `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Orders  []Order `gorm:"foreignKey:ProductID" json:"-"`
}

func (Product) TableName() string {
	return "products"
}

const (
	ProductStatusPending  = 0
	ProductStatusApproved = 1
	ProductStatusRejected = 2
	ProductStatusOnSale   = 3
	ProductStatusSoldOut  = 4
	ProductStatusOffShelf = 5
)

var ProductCategories = []string{"手机", "电脑", "相机", "耳机", "平板", "智能手表", "游戏机", "其他数码"}
var ProductConditions = []string{"全新", "95新", "9成新", "8成新", "7成新及以下"}

type ProductImage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	ImageURL  string         `gorm:"size:255;not null" json:"image_url"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

type Favorite struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Favorite) TableName() string {
	return "favorites"
}
