package model

import "time"

type ProductionStatus string

const (
	ProductionStatusQueued     ProductionStatus = "queued"
	ProductionStatusCutting    ProductionStatus = "cutting"
	ProductionStatusAssembling ProductionStatus = "assembling"
	ProductionStatusFinishing  ProductionStatus = "finishing"
	ProductionStatusPackaging  ProductionStatus = "packaging"
	ProductionStatusCompleted  ProductionStatus = "completed"
)

func (s ProductionStatus) Valid() bool {
	switch s {
	case ProductionStatusQueued, ProductionStatusCutting, ProductionStatusAssembling,
		ProductionStatusFinishing, ProductionStatusPackaging, ProductionStatusCompleted:
		return true
	}
	return false
}

func (s ProductionStatus) Next() ProductionStatus {
	switch s {
	case ProductionStatusQueued:
		return ProductionStatusCutting
	case ProductionStatusCutting:
		return ProductionStatusAssembling
	case ProductionStatusAssembling:
		return ProductionStatusFinishing
	case ProductionStatusFinishing:
		return ProductionStatusPackaging
	case ProductionStatusPackaging:
		return ProductionStatusCompleted
	}
	return s
}

func (s ProductionStatus) ProgressPercent() int {
	switch s {
	case ProductionStatusQueued:
		return 0
	case ProductionStatusCutting:
		return 20
	case ProductionStatusAssembling:
		return 40
	case ProductionStatusFinishing:
		return 60
	case ProductionStatusPackaging:
		return 80
	case ProductionStatusCompleted:
		return 100
	}
	return 0
}

type Production struct {
	ID              int64            `json:"id" gorm:"primaryKey"`
	OrderID         int64            `json:"order_id" gorm:"index;not null"`
	Status          ProductionStatus `json:"status" gorm:"size:32;not null"`
	ProgressPercent int              `json:"progress_percent" gorm:"not null"`
	OperatorID      int64            `json:"operator_id" gorm:"index"`
	Remark          string           `json:"remark" gorm:"size:500"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func (Production) TableName() string {
	return "productions"
}
