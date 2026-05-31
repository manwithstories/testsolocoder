package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShipType string

const (
	ShipTypeSailboat ShipType = "sailboat"
	ShipTypeMotorboat ShipType = "motorboat"
	ShipTypeYacht    ShipType = "yacht"
	ShipTypeFishing  ShipType = "fishing"
	ShipTypeCargo    ShipType = "cargo"
)

type ShipStatus string

const (
	ShipStatusAvailable  ShipStatus = "available"
	ShipStatusRented     ShipStatus = "rented"
	ShipStatusMaintenance ShipStatus = "maintenance"
	ShipStatusInactive   ShipStatus = "inactive"
)

type Ship struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	Owner           User           `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Name            string         `gorm:"size:200;not null" json:"name" binding:"required,min=2,max=200"`
	Description     string         `gorm:"type:text" json:"description"`
	ShipType        ShipType       `gorm:"type:varchar(50);not null" json:"ship_type" binding:"required,oneof=sailboat motorboat yacht fishing cargo"`
	Capacity        int            `gorm:"not null;default:0" json:"capacity" binding:"required,min=1"`
	CabinCount      int            `gorm:"default:0" json:"cabin_count"`
	BathroomCount   int            `gorm:"default:0" json:"bathroom_count"`
	Length          float64        `gorm:"type:decimal(10,2)" json:"length"`
	Width           float64        `gorm:"type:decimal(10,2)" json:"width"`
	YearBuilt       int            `json:"year_built"`
	Equipment       string         `gorm:"type:text" json:"equipment"`
	Features        string         `gorm:"type:text" json:"features"`
	SailingArea     string         `gorm:"size:500" json:"sailing_area"`
	HomePort        string         `gorm:"size:200" json:"home_port"`
	LicenseNumber   string         `gorm:"size:100" json:"license_number"`
	HourlyRate      float64        `gorm:"type:decimal(12,2);not null;default:0" json:"hourly_rate"`
	DailyRate       float64        `gorm:"type:decimal(12,2);not null;default:0" json:"daily_rate"`
	DepositAmount   float64        `gorm:"type:decimal(12,2);not null;default:0" json:"deposit_amount"`
	InsuranceRequired bool         `gorm:"default:true" json:"insurance_required"`
	CancellationPolicy string      `gorm:"type:text" json:"cancellation_policy"`
	Status          ShipStatus     `gorm:"type:varchar(20);default:available" json:"status"`
	AverageRating   float64        `gorm:"type:decimal(3,2);default:0" json:"average_rating"`
	ReviewCount     int            `gorm:"default:0" json:"review_count"`
	ViewCount       int            `gorm:"default:0" json:"view_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Images          []ShipImage    `gorm:"foreignKey:ShipID" json:"images,omitempty"`
}

type ShipImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ShipID    uuid.UUID `gorm:"type:uuid;not null;index" json:"ship_id"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Ship) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (si *ShipImage) BeforeCreate(tx *gorm.DB) error {
	if si.ID == uuid.Nil {
		si.ID = uuid.New()
	}
	return nil
}

type CreateShipRequest struct {
	Name              string   `json:"name" binding:"required,min=2,max=200"`
	Description       string   `json:"description"`
	ShipType          ShipType `json:"ship_type" binding:"required,oneof=sailboat motorboat yacht fishing cargo"`
	Capacity          int      `json:"capacity" binding:"required,min=1"`
	CabinCount        int      `json:"cabin_count"`
	BathroomCount     int      `json:"bathroom_count"`
	Length            float64  `json:"length"`
	Width             float64  `json:"width"`
	YearBuilt         int      `json:"year_built"`
	Equipment         string   `json:"equipment"`
	Features          string   `json:"features"`
	SailingArea       string   `json:"sailing_area"`
	HomePort          string   `json:"home_port"`
	LicenseNumber     string   `json:"license_number"`
	HourlyRate        float64  `json:"hourly_rate" binding:"required,min=0"`
	DailyRate         float64  `json:"daily_rate" binding:"required,min=0"`
	DepositAmount     float64  `json:"deposit_amount"`
	InsuranceRequired bool     `json:"insurance_required"`
	CancellationPolicy string  `json:"cancellation_policy"`
}

type UpdateShipRequest struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	ShipType          *ShipType  `json:"ship_type"`
	Capacity          *int       `json:"capacity"`
	CabinCount        *int       `json:"cabin_count"`
	BathroomCount     *int       `json:"bathroom_count"`
	Length            *float64   `json:"length"`
	Width             *float64   `json:"width"`
	YearBuilt         *int       `json:"year_built"`
	Equipment         string     `json:"equipment"`
	Features          string     `json:"features"`
	SailingArea       string     `json:"sailing_area"`
	HomePort          string     `json:"home_port"`
	LicenseNumber     string     `json:"license_number"`
	HourlyRate        *float64   `json:"hourly_rate"`
	DailyRate         *float64   `json:"daily_rate"`
	DepositAmount     *float64   `json:"deposit_amount"`
	InsuranceRequired *bool      `json:"insurance_required"`
	CancellationPolicy string    `json:"cancellation_policy"`
	Status            *ShipStatus `json:"status"`
}

type UploadImageRequest struct {
	IsPrimary bool `json:"is_primary"`
	SortOrder int  `json:"sort_order"`
}

type SearchShipRequest struct {
	ShipType    *ShipType `form:"ship_type"`
	MinCapacity *int      `form:"min_capacity"`
	MaxCapacity *int      `form:"max_capacity"`
	MinPrice    *float64  `form:"min_price"`
	MaxPrice    *float64  `form:"max_price"`
	Location    string    `form:"location"`
	StartDate   string    `form:"start_date"`
	EndDate     string    `form:"end_date"`
	Page        int       `form:"page,default=1"`
	PageSize    int       `form:"page_size,default=10"`
	SortBy      string    `form:"sort_by,default=rating"`
	SortOrder   string    `form:"sort_order,default=desc"`
}
