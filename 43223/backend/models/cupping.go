package models

import (
	"time"

	"gorm.io/gorm"
)

type CuppingScore struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	ProductID     uint             `json:"product_id" gorm:"index;not null"`
	Product       *Product         `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	UserID        uint             `json:"user_id" gorm:"index;not null"`
	User          *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DryFragrance  float64          `json:"dry_fragrance" gorm:"type:decimal(4,2)"`
	WetAroma      float64          `json:"wet_aroma" gorm:"type:decimal(4,2)"`
	Body          float64          `json:"body" gorm:"type:decimal(4,2)"`
	Acidity       float64          `json:"acidity" gorm:"type:decimal(4,2)"`
	Sweetness     float64          `json:"sweetness" gorm:"type:decimal(4,2)"`
	Aftertaste    float64          `json:"aftertaste" gorm:"type:decimal(4,2)"`
	Balance       float64          `json:"balance" gorm:"type:decimal(4,2)"`
	OverallScore  float64          `json:"overall_score" gorm:"type:decimal(4,2)"`
	Notes         string           `json:"notes" gorm:"type:text"`
	CriteriaItems []CuppingCriteria `json:"criteria_items,omitempty" gorm:"foreignKey:CuppingScoreID"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `json:"-" gorm:"index"`
}

type CuppingCriteria struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	CuppingScoreID  uint           `json:"cupping_score_id" gorm:"index;not null"`
	Name            string         `json:"name" gorm:"size:50"`
	Score           float64        `json:"score" gorm:"type:decimal(4,2)"`
	Description     string         `json:"description" gorm:"size:500"`
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateCuppingScoreRequest struct {
	ProductID    uint    `json:"product_id" binding:"required"`
	DryFragrance float64 `json:"dry_fragrance" binding:"required,gte=0,lte=10"`
	WetAroma     float64 `json:"wet_aroma" binding:"required,gte=0,lte=10"`
	Body         float64 `json:"body" binding:"required,gte=0,lte=10"`
	Acidity      float64 `json:"acidity" binding:"required,gte=0,lte=10"`
	Sweetness    float64 `json:"sweetness" binding:"required,gte=0,lte=10"`
	Aftertaste   float64 `json:"aftertaste" binding:"required,gte=0,lte=10"`
	Balance      float64 `json:"balance" binding:"required,gte=0,lte=10"`
	Notes        string  `json:"notes"`
}

type CuppingStats struct {
	ProductID      uint    `json:"product_id"`
	ProductName    string  `json:"product_name"`
	TotalCount     int64   `json:"total_count"`
	AvgScore       float64 `json:"avg_score"`
	AvgDryFragrance float64 `json:"avg_dry_fragrance"`
	AvgWetAroma    float64 `json:"avg_wet_aroma"`
	AvgBody        float64 `json:"avg_body"`
	AvgAcidity     float64 `json:"avg_acidity"`
	AvgSweetness   float64 `json:"avg_sweetness"`
	AvgAftertaste  float64 `json:"avg_aftertaste"`
	AvgBalance     float64 `json:"avg_balance"`
}

type ScoreTrendItem struct {
	Date      string  `json:"date"`
	AvgScore  float64 `json:"avg_score"`
	Count     int64   `json:"count"`
}
