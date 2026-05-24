package model

import (
	"time"
)

// 产品状态常量
const (
	ProductStatusOnSale   = 1 // 上架
	ProductStatusOffSale  = 2 // 下架
	ProductStatusSoldOut  = 3 // 售罄
)

// 选项类型常量
const (
	OptionTypeSize     = "size"     // 尺寸
	OptionTypeMaterial = "material" // 材质
	OptionTypeColor    = "color"    // 颜色
)

// ValidProductStatus 校验产品状态是否合法
func ValidProductStatus(status int) bool {
	switch status {
	case ProductStatusOnSale, ProductStatusOffSale, ProductStatusSoldOut:
		return true
	default:
		return false
	}
}

// ValidOptionType 校验选项类型是否合法
func ValidOptionType(t string) bool {
	switch t {
	case OptionTypeSize, OptionTypeMaterial, OptionTypeColor:
		return true
	default:
		return false
	}
}

// Product 产品模型
type Product struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ManufacturerID uint      `gorm:"index;not null" json:"manufacturer_id"` // 所属厂商用户 ID
	Name           string    `gorm:"size:128;not null" json:"name"`
	Category       string    `gorm:"size:64;index" json:"category"`
	Description    string    `gorm:"type:text" json:"description"`
	BasePrice      float64   `gorm:"type:decimal(12,2);not null" json:"base_price"`
	Stock          int       `gorm:"not null;default:0" json:"stock"`
	Status         int       `gorm:"not null;default:1;index" json:"status"`
	IsHot          bool      `gorm:"not null;default:false;index" json:"is_hot"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// IsOnSale 是否在售
func (p *Product) IsOnSale() bool {
	return p.Status == ProductStatusOnSale
}

// ProductOption 产品选项（尺寸/材质/颜色等）
type ProductOption struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProductID    uint      `gorm:"index;not null" json:"product_id"`
	OptionType   string    `gorm:"size:32;not null" json:"option_type"`
	OptionValue  string    `gorm:"size:128;not null" json:"option_value"`
	PriceAdjust  float64   `gorm:"type:decimal(12,2);not null;default:0" json:"price_adjust"`
	Sort         int       `gorm:"not null;default:0" json:"sort"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ProductOption) TableName() string {
	return "product_options"
}

// ProductImage 产品图片
type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	ImageURL  string    `gorm:"size:255;not null" json:"image_url"`
	Sort      int       `gorm:"not null;default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (ProductImage) TableName() string {
	return "product_images"
}
