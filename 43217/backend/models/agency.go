package models

import (
	"time"

	"gorm.io/gorm"
)

type AgencyStatus int

const (
	AgencyStatusPending  AgencyStatus = 0
	AgencyStatusActive   AgencyStatus = 1
	AgencyStatusDisabled AgencyStatus = 2
)

type Agency struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	UnifiedCode     string         `gorm:"uniqueIndex;size:50" json:"unified_code"`
	LegalPerson     string         `gorm:"size:50" json:"legal_person"`
	ContactPhone    string         `gorm:"size:20" json:"contact_phone"`
	ContactEmail    string         `gorm:"size:100" json:"contact_email"`
	Address         string         `gorm:"size:255" json:"address"`
	BusinessLicense string         `gorm:"size:255" json:"business_license"`
	Description     string         `gorm:"type:text" json:"description"`
	Rating          float64        `gorm:"default:5.0" json:"rating"`
	Status          AgencyStatus   `gorm:"default:0" json:"status"`
	AdminUserID     *uint          `gorm:"index" json:"admin_user_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Packages        []Package      `gorm:"foreignKey:AgencyID" json:"packages,omitempty"`
}

type PackageStatus int

const (
	PackageStatusOffline PackageStatus = 0
	PackageStatusOnline  PackageStatus = 1
)

type Package struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AgencyID      uint           `gorm:"index;not null" json:"agency_id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	OriginalPrice float64        `gorm:"default:0" json:"original_price"`
	Price         float64        `gorm:"default:0" json:"price"`
	SuitableFor   string         `gorm:"size:50" json:"suitable_for"`
	GenderLimit   int            `gorm:"default:0" json:"gender_limit"`
	MinAge        int            `gorm:"default:0" json:"min_age"`
	MaxAge        int            `gorm:"default:150" json:"max_age"`
	Notes         string         `gorm:"type:text" json:"notes"`
	Status        PackageStatus  `gorm:"default:1" json:"status"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	SaleCount     int            `gorm:"default:0" json:"sale_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Agency        Agency         `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	Items         []PackageItem  `gorm:"foreignKey:PackageID" json:"items,omitempty"`
	TimeSlots     []TimeSlot     `gorm:"foreignKey:PackageID" json:"time_slots,omitempty"`
}

type PackageItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PackageID   uint           `gorm:"index;not null" json:"package_id"`
	ItemName    string         `gorm:"size:100;not null" json:"item_name"`
	ItemCode    string         `gorm:"size:50" json:"item_code"`
	Description string         `gorm:"type:text" json:"description"`
	Department  string         `gorm:"size:50" json:"department"`
	NormalRange string         `gorm:"size:255" json:"normal_range"`
	Unit        string         `gorm:"size:20" json:"unit"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsRequired  bool           `gorm:"default:true" json:"is_required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Package     Package        `gorm:"foreignKey:PackageID" json:"package,omitempty"`
}

type TimeSlot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PackageID uint      `gorm:"index;not null" json:"package_id"`
	Date      time.Time `gorm:"index;not null" json:"date"`
	StartTime string    `gorm:"size:10;not null" json:"start_time"`
	EndTime   string    `gorm:"size:10;not null" json:"end_time"`
	Total     int       `gorm:"default:10" json:"total"`
	Booked    int       `gorm:"default:0" json:"booked"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Package   Package   `gorm:"foreignKey:PackageID" json:"package,omitempty"`
}
