package model

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusNoShow    BookingStatus = "no_show"
)

type Booking struct {
	ID            uint          `gorm:"primarykey" json:"id"`
	BookingNo     string        `gorm:"uniqueIndex;size:50;not null" json:"booking_no"`
	UserID        uint          `gorm:"index;not null" json:"user_id"`
	User          *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CarID         uint          `gorm:"index;not null" json:"car_id"`
	Car           *Car          `gorm:"foreignKey:CarID" json:"car,omitempty"`
	PickupStoreID uint          `gorm:"index;not null" json:"pickup_store_id"`
	PickupStore   *Store        `gorm:"foreignKey:PickupStoreID" json:"pickup_store,omitempty"`
	ReturnStoreID uint          `gorm:"index;not null" json:"return_store_id"`
	ReturnStore   *Store        `gorm:"foreignKey:ReturnStoreID" json:"return_store,omitempty"`
	PickupTime    time.Time     `gorm:"not null" json:"pickup_time"`
	ReturnTime    time.Time     `gorm:"not null" json:"return_time"`
	ActualPickupTime *time.Time `json:"actual_pickup_time"`
	ActualReturnTime *time.Time `json:"actual_return_time"`
	TotalDays     int           `gorm:"not null" json:"total_days"`
	BasePrice     float64       `gorm:"type:decimal(10,2);not null" json:"base_price"`
	TotalPrice    float64       `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Discount      float64       `gorm:"type:decimal(10,2);default:0" json:"discount"`
	FinalPrice    float64       `gorm:"type:decimal(10,2);not null" json:"final_price"`
	Deposit       float64       `gorm:"type:decimal(10,2);default:0" json:"deposit"`
	PromoCodeID   *uint         `gorm:"index" json:"promo_code_id"`
	PromoCode     *PromoCode    `gorm:"foreignKey:PromoCodeID" json:"promo_code,omitempty"`
	Status        BookingStatus `gorm:"size:20;default:pending" json:"status"`
	Remark        string        `gorm:"type:text" json:"remark"`
	CancelReason  string        `gorm:"size:500" json:"cancel_reason"`
	CancelledAt   *time.Time    `json:"cancelled_at"`
	OrderID       *uint         `gorm:"index" json:"order_id"`
	Order         *Order        `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	PickupReminderSent bool    `gorm:"default:false" json:"pickup_reminder_sent"`
	ReturnReminderSent bool    `gorm:"default:false" json:"return_reminder_sent"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     *time.Time    `gorm:"index" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}

type PromoCode struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Code        string     `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Type        string     `gorm:"size:20;not null" json:"type"`
	Value       float64    `gorm:"type:decimal(10,2);not null" json:"value"`
	MinAmount   float64    `gorm:"type:decimal(10,2);default:0" json:"min_amount"`
	MaxDiscount float64    `gorm:"type:decimal(10,2);default:0" json:"max_discount"`
	UsageLimit  int        `gorm:"default:1" json:"usage_limit"`
	UsedCount   int        `gorm:"default:0" json:"used_count"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (PromoCode) TableName() string {
	return "promo_codes"
}

type PricingRule struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	RuleType    string     `gorm:"size:20;not null" json:"rule_type"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Weekdays    string     `gorm:"size:50" json:"weekdays"`
	Multiplier  float64    `gorm:"type:decimal(5,2);not null" json:"multiplier"`
	CarModel    string     `gorm:"size:100" json:"car_model"`
	MinDays     int        `gorm:"default:1" json:"min_days"`
	MaxDays     int        `gorm:"default:0" json:"max_days"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	Priority    int        `gorm:"default:0" json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (PricingRule) TableName() string {
	return "pricing_rules"
}