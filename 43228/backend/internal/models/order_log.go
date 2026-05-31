package models

import (
	"time"
)

type OrderLog struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	OrderID     uint        `gorm:"index;not null" json:"order_id"`
	Status      string      `gorm:"size:32;not null" json:"status"`
	Description string      `gorm:"size:512;not null" json:"description"`
	OperatorID  uint        `gorm:"index" json:"operator_id"`
	CreatedAt   time.Time   `json:"created_at"`

	Order    *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Operator *User  `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}

func (OrderLog) TableName() string {
	return "order_logs"
}
