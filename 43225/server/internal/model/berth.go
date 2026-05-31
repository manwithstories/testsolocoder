package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BerthType string

const (
	BerthTypeSmall  BerthType = "small"
	BerthTypeMedium BerthType = "medium"
	BerthTypeLarge  BerthType = "large"
)

type BerthStatus string

const (
	BerthStatusAvailable BerthStatus = "available"
	BerthStatusOccupied  BerthStatus = "occupied"
	BerthStatusReserved  BerthStatus = "reserved"
	BerthStatusMaintenance BerthStatus = "maintenance"
)

type Berth struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	DockID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"dock_id"`
	Number        string         `gorm:"size:50;not null" json:"number" binding:"required"`
	BerthType     BerthType      `gorm:"type:varchar(20);not null" json:"berth_type" binding:"required,oneof=small medium large"`
	MaxLength     float64        `gorm:"type:decimal(10,2)" json:"max_length"`
	MaxWidth      float64        `gorm:"type:decimal(10,2)" json:"max_width"`
	HourlyRate    float64        `gorm:"type:decimal(12,2);not null;default:0" json:"hourly_rate"`
	DailyRate     float64        `gorm:"type:decimal(12,2);not null;default:0" json:"daily_rate"`
	HasWater      bool           `gorm:"default:true" json:"has_water"`
	HasElectric   bool           `gorm:"default:true" json:"has_electric"`
	HasInternet   bool           `gorm:"default:true" json:"has_internet"`
	Description   string         `gorm:"type:text" json:"description"`
	Status        BerthStatus    `gorm:"type:varchar(20);default:available" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Reservations  []BerthReservation `gorm:"foreignKey:BerthID" json:"reservations,omitempty"`
}

type BerthReservation struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BerthID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"berth_id"`
	Berth       Berth          `gorm:"foreignKey:BerthID" json:"berth,omitempty"`
	ShipID      *uuid.UUID     `gorm:"type:uuid;index" json:"ship_id"`
	Ship        *Ship          `gorm:"foreignKey:ShipID" json:"ship,omitempty"`
	RentalID    *uuid.UUID     `gorm:"type:uuid;index" json:"rental_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	StartTime   time.Time      `gorm:"not null" json:"start_time"`
	EndTime     time.Time      `gorm:"not null" json:"end_time"`
	TotalAmount float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status      string         `gorm:"type:varchar(20);default:confirmed" json:"status"`
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WaterLevel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DockID    uuid.UUID `gorm:"type:uuid;not null;index" json:"dock_id"`
	Height    float64   `gorm:"type:decimal(10,2);not null" json:"height"`
	Unit      string    `gorm:"size:20;default:meters" json:"unit"`
	RecordedAt time.Time `gorm:"not null" json:"recorded_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Dock struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name" binding:"required"`
	Address     string         `gorm:"size:500;not null" json:"address" binding:"required"`
	City        string         `gorm:"size:100" json:"city"`
	Country     string         `gorm:"size:100" json:"country"`
	Latitude    float64        `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude   float64        `gorm:"type:decimal(10,7)" json:"longitude"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"size:500" json:"image_url"`
	Amenities   string         `gorm:"type:text" json:"amenities"`
	OpenTime    string         `gorm:"size:50" json:"open_time"`
	CloseTime   string         `gorm:"size:50" json:"close_time"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	AverageRating float64      `gorm:"type:decimal(3,2);default:0" json:"average_rating"`
	ReviewCount int            `gorm:"default:0" json:"review_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Berths      []Berth        `gorm:"foreignKey:DockID" json:"berths,omitempty"`
}

func (b *Berth) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (br *BerthReservation) BeforeCreate(tx *gorm.DB) error {
	if br.ID == uuid.Nil {
		br.ID = uuid.New()
	}
	return nil
}

func (w *WaterLevel) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (d *Dock) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type CreateBerthRequest struct {
	DockID      string    `json:"dock_id" binding:"required,uuid"`
	Number      string    `json:"number" binding:"required"`
	BerthType   BerthType `json:"berth_type" binding:"required,oneof=small medium large"`
	MaxLength   float64   `json:"max_length"`
	MaxWidth    float64   `json:"max_width"`
	HourlyRate  float64   `json:"hourly_rate" binding:"required,min=0"`
	DailyRate   float64   `json:"daily_rate" binding:"required,min=0"`
	HasWater    bool      `json:"has_water"`
	HasElectric bool      `json:"has_electric"`
	HasInternet bool      `json:"has_internet"`
	Description string    `json:"description"`
}

type CreateReservationRequest struct {
	BerthID   string    `json:"berth_id" binding:"required,uuid"`
	ShipID    string    `json:"ship_id"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Notes     string    `json:"notes"`
}

type CheckAvailabilityRequest struct {
	BerthID   string    `form:"berth_id" binding:"required,uuid"`
	StartTime time.Time `form:"start_time" binding:"required"`
	EndTime   time.Time `form:"end_time" binding:"required"`
}

type CreateDockRequest struct {
	Name        string  `json:"name" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Amenities   string  `json:"amenities"`
	OpenTime    string  `json:"open_time"`
	CloseTime   string  `json:"close_time"`
}

type RecordWaterLevelRequest struct {
	DockID     string    `json:"dock_id" binding:"required,uuid"`
	Height     float64   `json:"height" binding:"required"`
	Unit       string    `json:"unit"`
	RecordedAt time.Time `json:"recorded_at" binding:"required"`
}
