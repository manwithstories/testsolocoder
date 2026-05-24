package model

import (
	"encoding/json"
	"time"
)

// Review 售后评价模型
type Review struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"index;not null" json:"order_id"`
	ProductID     uint      `gorm:"index;not null" json:"product_id"`
	OwnerID       uint      `gorm:"index;not null" json:"owner_id"`
	ProductRating int       `gorm:"not null;default:5" json:"product_rating"`
	ServiceRating int       `gorm:"not null;default:5" json:"service_rating"`
	Content       string    `gorm:"type:text" json:"content"`
	Images        string    `gorm:"type:text" json:"images"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 指定表名
func (Review) TableName() string {
	return "reviews"
}

// ParseImages 解析图片 JSON 数组
func (r *Review) ParseImages() ([]string, error) {
	if r.Images == "" {
		return nil, nil
	}
	var list []string
	if err := json.Unmarshal([]byte(r.Images), &list); err != nil {
		return nil, err
	}
	return list, nil
}

// SetImages 将图片列表序列化为 JSON
func (r *Review) SetImages(images []string) error {
	if len(images) == 0 {
		r.Images = ""
		return nil
	}
	data, err := json.Marshal(images)
	if err != nil {
		return err
	}
	r.Images = string(data)
	return nil
}

// ValidRating 校验评分范围（1-5）
func ValidRating(rating int) bool {
	return rating >= 1 && rating <= 5
}
