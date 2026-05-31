package models

import (
	"time"

	"gorm.io/gorm"
)

type TeaImageType string

const (
	TeaImageTypeMain      TeaImageType = "main"
	TeaImageTypeDetail    TeaImageType = "detail"
	TeaImageTypePackaging TeaImageType = "packaging"
)

type TeaImage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeaID     uint           `gorm:"index;not null" json:"tea_id"`
	ImageURL  string         `gorm:"size:512;not null" json:"image_url"`
	ImageType TeaImageType   `gorm:"size:32;not null;default:detail" json:"image_type"`
	Sort      int            `gorm:"not null;default:0" json:"sort"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Tea *Tea `gorm:"foreignKey:TeaID" json:"tea,omitempty"`
}

func (TeaImage) TableName() string {
	return "tea_images"
}
