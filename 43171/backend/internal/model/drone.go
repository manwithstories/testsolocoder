package model

import (
	"time"

	"gorm.io/gorm"
)

type DroneStatus string

const (
	DroneStatusOffline  DroneStatus = "offline"
	DroneStatusOnline   DroneStatus = "online"
	DroneStatusRented  DroneStatus = "rented"
	DroneStatusMaint   DroneStatus = "maintenance"
)

type Drone struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OwnerID       uint           `gorm:"index;not null" json:"owner_id"`
	Owner         *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Name          string         `gorm:"size:128;not null" json:"name"`
	Brand         string         `gorm:"size:64" json:"brand"`
	Model         string         `gorm:"size:64" json:"model"`
	SerialNo      string         `gorm:"size:64;uniqueIndex" json:"serial_no"`
	Weight        float64        `json:"weight"`
	BatteryLife   int            `json:"battery_life"`
	GimbalSpec    string         `gorm:"size:128" json:"gimbal_spec"`
	CameraSpec    string         `gorm:"size:128" json:"camera_spec"`
	MaxSpeed      float64        `json:"max_speed"`
	MaxAltitude   float64        `json:"max_altitude"`
	PurchaseDate  *time.Time     `json:"purchase_date"`
	Region        string         `gorm:"size:64" json:"region"`
	Description   string         `gorm:"size:500" json:"description"`
	Images        string         `gorm:"size:1024" json:"images"`
	PricePerDay   float64        `gorm:"not null" json:"price_per_day"`
	Deposit       float64        `json:"deposit"`
	Status        DroneStatus    `gorm:"size:16;default:offline;index" json:"status"`
	Rating        float64        `gorm:"default:5.0" json:"rating"`
	RatingCount   int            `gorm:"default:0" json:"rating_count"`
	AvailableFrom *time.Time     `json:"available_from"`
	AvailableTo   *time.Time     `json:"available_to"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
