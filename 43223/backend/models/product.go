package models

import (
	"time"

	"gorm.io/gorm"
)

type RoastLevel string

const (
	RoastLevelLight      RoastLevel = "light"
	RoastLevelMedium     RoastLevel = "medium"
	RoastLevelMediumDark RoastLevel = "medium_dark"
	RoastLevelDark       RoastLevel = "dark"
)

type ProcessMethod string

const (
	ProcessWashed    ProcessMethod = "washed"
	ProcessNatural   ProcessMethod = "natural"
	ProcessHoney     ProcessMethod = "honey"
	ProcessAnaerobic ProcessMethod = "anaerobic"
	ProcessWetHulled ProcessMethod = "wet_hulled"
)

type ProductStatus string

const (
	ProductStatusOnSale  ProductStatus = "on_sale"
	ProductStatusOffline ProductStatus = "offline"
	ProductStatusDraft   ProductStatus = "draft"
)

type Product struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"size:200;not null"`
	Origin        string         `json:"origin" gorm:"size:100;not null"`
	Farm          string         `json:"farm" gorm:"size:200"`
	Variety       string         `json:"variety" gorm:"size:100"`
	Altitude      string         `json:"altitude" gorm:"size:50"`
	ProcessMethod ProcessMethod  `json:"process_method" gorm:"size:50;not null"`
	RoastLevel    RoastLevel     `json:"roast_level" gorm:"size:50;not null"`
	FlavorNotes   string         `json:"flavor_notes" gorm:"type:text"`
	CuppingScore  float64        `json:"cupping_score" gorm:"default:0"`
	Description   string         `json:"description" gorm:"type:text"`
	Price         float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Weight        int            `json:"weight" gorm:"not null"`
	Stock         int            `json:"stock" gorm:"default:0"`
	Status        ProductStatus  `json:"status" gorm:"size:20;default:draft"`
	RoasterID     uint           `json:"roaster_id"`
	Roaster       *User          `json:"roaster,omitempty" gorm:"foreignKey:RoasterID"`
	Images        []ProductImage `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID uint           `json:"product_id" gorm:"index;not null"`
	URL       string         `json:"url" gorm:"size:500;not null"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	IsCover   bool           `json:"is_cover" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateProductRequest struct {
	Name          string        `json:"name" binding:"required,max=200"`
	Origin        string        `json:"origin" binding:"required"`
	Farm          string        `json:"farm"`
	Variety       string        `json:"variety"`
	Altitude      string        `json:"altitude"`
	ProcessMethod ProcessMethod `json:"process_method" binding:"required"`
	RoastLevel    RoastLevel    `json:"roast_level" binding:"required"`
	FlavorNotes   string        `json:"flavor_notes"`
	CuppingScore  float64       `json:"cupping_score"`
	Description   string        `json:"description"`
	Price         float64       `json:"price" binding:"required,gt=0"`
	Weight        int           `json:"weight" binding:"required,gt=0"`
	Stock         int           `json:"stock"`
	Status        ProductStatus `json:"status"`
}

type UpdateProductRequest struct {
	Name          *string        `json:"name"`
	Origin        *string        `json:"origin"`
	Farm          *string        `json:"farm"`
	Variety       *string        `json:"variety"`
	Altitude      *string        `json:"altitude"`
	ProcessMethod *ProcessMethod `json:"process_method"`
	RoastLevel    *RoastLevel    `json:"roast_level"`
	FlavorNotes   *string        `json:"flavor_notes"`
	CuppingScore  *float64       `json:"cupping_score"`
	Description   *string        `json:"description"`
	Price         *float64       `json:"price"`
	Weight        *int           `json:"weight"`
	Stock         *int           `json:"stock"`
}

type UpdateProductStatusRequest struct {
	Status ProductStatus `json:"status" binding:"required,oneof=on_sale offline draft"`
}

type CSVProductRow struct {
	Name          string `csv:"name"`
	Origin        string `csv:"origin"`
	Farm          string `csv:"farm"`
	Variety       string `csv:"variety"`
	Altitude      string `csv:"altitude"`
	ProcessMethod string `csv:"process_method"`
	RoastLevel    string `csv:"roast_level"`
	FlavorNotes   string `csv:"flavor_notes"`
	CuppingScore  string `csv:"cupping_score"`
	Description   string `csv:"description"`
	Price         string `csv:"price"`
	Weight        string `csv:"weight"`
	Stock         string `csv:"stock"`
}
