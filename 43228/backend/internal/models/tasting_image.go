package models

import (
	"time"

	"gorm.io/gorm"
)

type TastingImage struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TastingRecordID uint           `gorm:"index;not null" json:"tasting_record_id"`
	ImageURL        string         `gorm:"size:512;not null" json:"image_url"`
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	TastingRecord *TastingRecord `gorm:"foreignKey:TastingRecordID" json:"tasting_record,omitempty"`
}

func (TastingImage) TableName() string {
	return "tasting_images"
}
