package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal StringArray value")
	}
	return json.Unmarshal(bytes, a)
}

type RoomType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	BasePrice   float64        `gorm:"type:decimal(10,2);not null" json:"base_price"`
	BedCount    int            `gorm:"not null;default:1" json:"bed_count"`
	MaxGuests   int            `gorm:"not null;default:1" json:"max_guests"`
	Facilities  StringArray    `gorm:"type:jsonb" json:"facilities"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Rooms       []Room         `gorm:"foreignKey:RoomTypeID" json:"rooms,omitempty"`
}
