package models

import (
	"time"

	"gorm.io/gorm"
)

type Traceability struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TeaID       uint           `gorm:"index;not null" json:"tea_id"`
	Stage       string         `gorm:"size:64;not null" json:"stage"`
	Description string         `gorm:"type:text" json:"description"`
	Location    string         `gorm:"size:255" json:"location"`
	Operator    string         `gorm:"size:128" json:"operator"`
	RecordTime  time.Time      `json:"record_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Tea *Tea `gorm:"foreignKey:TeaID" json:"tea,omitempty"`
}

func (Traceability) TableName() string {
	return "traceability"
}
