package model

import (
	"time"

	"gorm.io/gorm"
)

type FlightRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OrderID     *uint          `gorm:"index" json:"order_id"`
	Order       *RentalOrder   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ServiceID   *uint          `gorm:"index" json:"service_id"`
	Service     *AerialService `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	DroneID     uint           `gorm:"index;not null" json:"drone_id"`
	Drone       *Drone         `gorm:"foreignKey:DroneID" json:"drone,omitempty"`
	PilotID     uint           `gorm:"index;not null" json:"pilot_id"`
	Pilot       *User          `gorm:"foreignKey:PilotID" json:"pilot,omitempty"`
	StartPoint  string         `gorm:"size:128" json:"start_point"`
	EndPoint    string         `gorm:"size:128" json:"end_point"`
	Route       string         `gorm:"size:2048" json:"route"`
	AltitudeMax float64        `json:"altitude_max"`
	AltitudeAvg float64        `json:"altitude_avg"`
	Duration    int            `json:"duration"`
	Distance    float64        `json:"distance"`
	FlightDate  time.Time      `json:"flight_date"`
	FlightLog   string         `gorm:"size:2048" json:"flight_log"`
	Images      string         `gorm:"size:1024" json:"images"`
	Remark      string         `gorm:"size:500" json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
