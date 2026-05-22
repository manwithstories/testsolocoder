package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type RoomStatus int

const (
	RoomStatusActive   RoomStatus = 1
	RoomStatusInactive RoomStatus = 0
)

type Room struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Floor       string         `json:"floor" gorm:"size:50;not null;index"`
	Capacity    int            `json:"capacity" gorm:"not null"`
	Location    string         `json:"location" gorm:"size:200"`
	PricePerHour float64       `json:"price_per_hour" gorm:"type:decimal(10,2);not null"`
	Equipment   string         `json:"equipment" gorm:"type:text"`
	Description string         `json:"description" gorm:"type:text"`
	AvailableStart string      `json:"available_start" gorm:"size:10;default:08:00"`
	AvailableEnd   string      `json:"available_end" gorm:"size:10;default:22:00"`
	Status      RoomStatus     `json:"status" gorm:"default:1"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Photos      []RoomPhoto    `json:"photos" gorm:"foreignKey:RoomID"`
}

func (Room) TableName() string {
	return "rooms"
}

func (r *Room) GetEquipmentList() []string {
	if r.Equipment == "" {
		return []string{}
	}
	var list []string
	json.Unmarshal([]byte(r.Equipment), &list)
	return list
}

func (r *Room) SetEquipmentList(items []string) {
	data, _ := json.Marshal(items)
	r.Equipment = string(data)
}

type RoomPhoto struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	RoomID    uint           `json:"room_id" gorm:"index;not null"`
	URL       string         `json:"url" gorm:"size:500;not null"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (RoomPhoto) TableName() string {
	return "room_photos"
}
