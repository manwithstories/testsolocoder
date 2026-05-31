package models

import (
	"time"

	"gorm.io/gorm"
)

type TastingRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	TeaID           uint           `gorm:"index;not null" json:"tea_id"`
	BrewMethod      string         `gorm:"size:128" json:"brew_method"`
	WaterTemp       float64        `json:"water_temp"`
	BrewTime        int            `json:"brew_time"`
	WaterQuality    string         `gorm:"size:128" json:"water_quality"`
	TeaAmount       float64        `json:"tea_amount"`
	AromaScore      float64        `json:"aroma_score"`
	TasteScore      float64        `json:"taste_score"`
	AftertasteScore float64        `json:"aftertaste_score"`
	OverallScore    float64        `json:"overall_score"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User   *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tea    *Tea            `gorm:"foreignKey:TeaID" json:"tea,omitempty"`
	Images []TastingImage  `gorm:"foreignKey:TastingRecordID" json:"images,omitempty"`
}

func (TastingRecord) TableName() string {
	return "tasting_records"
}
