package models

import (
	"time"

	"gorm.io/gorm"
)

type TeaType string

const (
	TeaTypeGreen  TeaType = "green_tea"
	TeaTypeBlack  TeaType = "black_tea"
	TeaTypeOolong TeaType = "oolong"
	TeaTypePuer   TeaType = "puer"
	TeaTypeWhite  TeaType = "white_tea"
	TeaTypeYellow TeaType = "yellow_tea"
	TeaTypeDark   TeaType = "dark_tea"
	TeaTypeFlower TeaType = "flower_tea"
)

type TeaGrade string

const (
	TeaGradeSpecial TeaGrade = "special"
	TeaGrade1       TeaGrade = "grade1"
	TeaGrade2       TeaGrade = "grade2"
	TeaGrade3       TeaGrade = "grade3"
)

type ProcessType string

const (
	ProcessTypeManual    ProcessType = "manual"
	ProcessTypeSemi      ProcessType = "semi_manual"
	ProcessTypeMechanism ProcessType = "mechanism"
)

type Tea struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"size:128;not null;index" json:"name"`
	Type             TeaType        `gorm:"size:32;not null;index" json:"type"`
	Origin           string         `gorm:"size:128" json:"origin"`
	Year             int            `json:"year"`
	Grade            TeaGrade       `gorm:"size:32" json:"grade"`
	ProcessType      ProcessType    `gorm:"size:32" json:"process_type"`
	StorageCondition string         `gorm:"size:255" json:"storage_condition"`
	Description      string         `gorm:"type:text" json:"description"`
	Price            float64        `gorm:"type:decimal(12,2);not null;default:0" json:"price"`
	Stock            int            `gorm:"not null;default:0" json:"stock"`
	SellerID         uint           `gorm:"index;not null" json:"seller_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Images          []TeaImage        `gorm:"foreignKey:TeaID" json:"images,omitempty"`
	TastingRecords  []TastingRecord   `gorm:"foreignKey:TeaID" json:"tasting_records,omitempty"`
	Collections     []Collection      `gorm:"foreignKey:TeaID" json:"collections,omitempty"`
	Orders          []Order           `gorm:"foreignKey:TeaID" json:"orders,omitempty"`
	Traceability    []Traceability    `gorm:"foreignKey:TeaID" json:"traceability,omitempty"`
	AppraisalReport *AppraisalReport  `gorm:"foreignKey:TeaID" json:"appraisal_report,omitempty"`
}

func (Tea) TableName() string {
	return "teas"
}
