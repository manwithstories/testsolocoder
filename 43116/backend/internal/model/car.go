package model

import (
	"time"
)

type CarStatus string

const (
	CarStatusAvailable CarStatus = "available"
	CarStatusRented    CarStatus = "rented"
	CarStatusMaintenance CarStatus = "maintenance"
	CarStatusDisabled  CarStatus = "disabled"
)

type Car struct {
	ID              uint        `gorm:"primarykey" json:"id"`
	Brand           string      `gorm:"size:50;not null" json:"brand"`
	Model           string      `gorm:"size:100;not null" json:"model"`
	Year            int         `gorm:"not null" json:"year"`
	Seats           int         `gorm:"not null" json:"seats"`
	Transmission    string      `gorm:"size:20;not null" json:"transmission"`
	FuelType        string      `gorm:"size:20" json:"fuel_type"`
	DailyRent       float64     `gorm:"type:decimal(10,2);not null" json:"daily_rent"`
	Deposit         float64     `gorm:"type:decimal(10,2);default:0" json:"deposit"`
	Status          CarStatus   `gorm:"size:20;default:available" json:"status"`
	LicensePlate    string      `gorm:"uniqueIndex;size:20" json:"license_plate"`
	Color           string      `gorm:"size:20" json:"color"`
	Mileage         int         `gorm:"default:0" json:"mileage"`
	Rating          float64     `gorm:"type:decimal(2,1);default:0" json:"rating"`
	ReviewCount     int         `gorm:"default:0" json:"review_count"`
	StoreID         uint        `json:"store_id"`
	Store           *Store      `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Images          []CarImage  `gorm:"foreignKey:CarID" json:"images,omitempty"`
	Features        string      `gorm:"type:text" json:"features"`
	Description     string      `gorm:"type:text" json:"description"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	DeletedAt       *time.Time  `gorm:"index" json:"-"`
}

func (Car) TableName() string {
	return "cars"
}

type CarImage struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CarID     uint       `gorm:"index;not null" json:"car_id"`
	URL       string     `gorm:"size:500;not null" json:"url"`
	IsCover   bool       `gorm:"default:false" json:"is_cover"`
	SortOrder int        `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
}

func (CarImage) TableName() string {
	return "car_images"
}

type City struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Name      string     `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code      string     `gorm:"uniqueIndex;size:10" json:"code"`
	Province  string     `gorm:"size:50" json:"province"`
	Stores    []Store    `gorm:"foreignKey:CityID" json:"stores,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (City) TableName() string {
	return "cities"
}

type Store struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	CityID      uint       `gorm:"index;not null" json:"city_id"`
	City        *City      `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Address     string     `gorm:"size:255;not null" json:"address"`
	Phone       string     `gorm:"size:20" json:"phone"`
	BusinessHours string   `gorm:"size:100" json:"business_hours"`
	Latitude    float64    `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude   float64    `gorm:"type:decimal(10,7)" json:"longitude"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Store) TableName() string {
	return "stores"
}