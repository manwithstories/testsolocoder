package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WeatherCondition string

const (
	WeatherSunny   WeatherCondition = "sunny"
	WeatherCloudy  WeatherCondition = "cloudy"
	WeatherRainy   WeatherCondition = "rainy"
	WeatherStormy  WeatherCondition = "stormy"
	WeatherFoggy   WeatherCondition = "foggy"
)

type VoyageLog struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RentalID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"rental_id"`
	Rental          Rental         `gorm:"foreignKey:RentalID" json:"rental,omitempty"`
	ShipID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"ship_id"`
	Ship            Ship           `gorm:"foreignKey:ShipID" json:"ship,omitempty"`
	CaptainID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"captain_id"`
	Captain         User           `gorm:"foreignKey:CaptainID" json:"captain,omitempty"`
	LogDate         time.Time      `gorm:"not null" json:"log_date"`
	DeparturePort   string         `gorm:"size:200" json:"departure_port"`
	ArrivalPort     string         `gorm:"size:200" json:"arrival_port"`
	DepartureTime   *time.Time     `json:"departure_time"`
	ArrivalTime     *time.Time     `json:"arrival_time"`
	Distance        float64        `gorm:"type:decimal(10,2)" json:"distance"`
	EngineHours     float64        `gorm:"type:decimal(10,2)" json:"engine_hours"`
	FuelConsumed    float64        `gorm:"type:decimal(10,2)" json:"fuel_consumed"`
	FuelUnit        string         `gorm:"size:20;default:liters" json:"fuel_unit"`
	AvgSpeed        float64        `gorm:"type:decimal(10,2)" json:"avg_speed"`
	MaxSpeed        float64        `gorm:"type:decimal(10,2)" json:"max_speed"`
	Weather         WeatherCondition `gorm:"type:varchar(20)" json:"weather"`
	WindSpeed       float64        `gorm:"type:decimal(10,2)" json:"wind_speed"`
	WindDirection   string         `gorm:"size:50" json:"wind_direction"`
	WaveHeight      float64        `gorm:"type:decimal(10,2)" json:"wave_height"`
	WaterTemperature float64       `gorm:"type:decimal(10,2)" json:"water_temperature"`
	AirTemperature  float64        `gorm:"type:decimal(10,2)" json:"air_temperature"`
	PassengerCount  int            `gorm:"default:0" json:"passenger_count"`
	CrewCount       int            `gorm:"default:0" json:"crew_count"`
	Notes           string         `gorm:"type:text" json:"notes"`
	Incidents       string         `gorm:"type:text" json:"incidents"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *VoyageLog) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

type CreateVoyageLogRequest struct {
	RentalID         string           `json:"rental_id" binding:"required,uuid"`
	LogDate          time.Time        `json:"log_date" binding:"required"`
	DeparturePort    string           `json:"departure_port"`
	ArrivalPort      string           `json:"arrival_port"`
	DepartureTime    *time.Time       `json:"departure_time"`
	ArrivalTime      *time.Time       `json:"arrival_time"`
	Distance         float64          `json:"distance"`
	EngineHours      float64          `json:"engine_hours"`
	FuelConsumed     float64          `json:"fuel_consumed"`
	FuelUnit         string           `json:"fuel_unit"`
	AvgSpeed         float64          `json:"avg_speed"`
	MaxSpeed         float64          `json:"max_speed"`
	Weather          WeatherCondition `json:"weather" binding:"required,oneof=sunny cloudy rainy stormy foggy"`
	WindSpeed        float64          `json:"wind_speed"`
	WindDirection    string           `json:"wind_direction"`
	WaveHeight       float64          `json:"wave_height"`
	WaterTemperature float64          `json:"water_temperature"`
	AirTemperature   float64          `json:"air_temperature"`
	PassengerCount   int              `json:"passenger_count"`
	CrewCount        int              `json:"crew_count"`
	Notes            string           `json:"notes"`
	Incidents        string           `json:"incidents"`
}

type ExportVoyageLogRequest struct {
	RentalID  string `form:"rental_id" binding:"required,uuid"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Format    string `form:"format" binding:"required,oneof=pdf csv"`
}
