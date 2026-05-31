package models

import (
	"time"

	"gorm.io/gorm"
)

type RoastingRecord struct {
	ID             uint                 `json:"id" gorm:"primaryKey"`
	ProductID      uint                 `json:"product_id" gorm:"index;not null"`
	Product        *Product             `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	RoasterID      uint                 `json:"roaster_id" gorm:"index;not null"`
	Roaster        *User                `json:"roaster,omitempty" gorm:"foreignKey:RoasterID"`
	BatchNumber    string               `json:"batch_number" gorm:"uniqueIndex;size:50;not null"`
	GreenBeanWeight float64             `json:"green_bean_weight" gorm:"type:decimal(10,2)"`
	RoastedWeight  float64              `json:"roasted_weight" gorm:"type:decimal(10,2)"`
	InputTemp      float64              `json:"input_temp"`
	TurningPoint   float64              `json:"turning_point"`
	TurningTime    int                  `json:"turning_time"`
	FirstCrackTemp float64              `json:"first_crack_temp"`
	FirstCrackTime int                  `json:"first_crack_time"`
	SecondCrackTemp float64             `json:"second_crack_temp"`
	SecondCrackTime int                 `json:"second_crack_time"`
	DropTemp       float64              `json:"drop_temp"`
	TotalRoastTime int                  `json:"total_roast_time"`
	Notes          string               `json:"notes" gorm:"type:text"`
	DataPoints     []RoastingDataPoint  `json:"data_points,omitempty" gorm:"foreignKey:RoastingRecordID"`
	RoastedAt      time.Time            `json:"roasted_at"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	DeletedAt      gorm.DeletedAt       `json:"-" gorm:"index"`
}

type RoastingDataPoint struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	RoastingRecordID  uint           `json:"roasting_record_id" gorm:"index;not null"`
	TimeElapsed       int            `json:"time_elapsed"`
	BeanTemp          float64        `json:"bean_temp"`
	EnvTemp           float64        `json:"env_temp"`
	RateOfRise        float64        `json:"rate_of_rise"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateRoastingRecordRequest struct {
	ProductID       uint                    `json:"product_id" binding:"required"`
	BatchNumber     string                  `json:"batch_number" binding:"required,max=50"`
	GreenBeanWeight float64                 `json:"green_bean_weight"`
	RoastedWeight   float64                 `json:"roasted_weight"`
	InputTemp       float64                 `json:"input_temp"`
	TurningPoint    float64                 `json:"turning_point"`
	TurningTime     int                     `json:"turning_time"`
	FirstCrackTemp  float64                 `json:"first_crack_temp"`
	FirstCrackTime  int                     `json:"first_crack_time"`
	SecondCrackTemp float64                 `json:"second_crack_temp"`
	SecondCrackTime int                     `json:"second_crack_time"`
	DropTemp        float64                 `json:"drop_temp"`
	TotalRoastTime  int                     `json:"total_roast_time"`
	Notes           string                  `json:"notes"`
	DataPoints      []RoastingDataPointItem `json:"data_points"`
	RoastedAt       time.Time               `json:"roasted_at"`
}

type RoastingDataPointItem struct {
	TimeElapsed int     `json:"time_elapsed"`
	BeanTemp    float64 `json:"bean_temp"`
	EnvTemp     float64 `json:"env_temp"`
	RateOfRise  float64 `json:"rate_of_rise"`
}

type RoastingComparisonResult struct {
	Records   []RoastingRecord `json:"records"`
	AvgInputTemp   float64 `json:"avg_input_temp"`
	AvgDropTemp    float64 `json:"avg_drop_temp"`
	AvgTotalTime   float64 `json:"avg_total_time"`
	WeightLossRate float64 `json:"weight_loss_rate"`
}
