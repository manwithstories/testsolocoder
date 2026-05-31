package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RentalType string

const (
	RentalTypeDaily  RentalType = "daily"
	RentalTypeHourly RentalType = "hourly"
	RentalTypeVoyage RentalType = "voyage"
)

type RentalStatus string

const (
	RentalStatusPending     RentalStatus = "pending"
	RentalStatusConfirmed   RentalStatus = "confirmed"
	RentalStatusActive      RentalStatus = "active"
	RentalStatusCompleted   RentalStatus = "completed"
	RentalStatusCancelled   RentalStatus = "cancelled"
	RentalStatusRefunded    RentalStatus = "refunded"
)

type InsuranceType string

const (
	InsuranceTypeBasic    InsuranceType = "basic"
	InsuranceTypePremium  InsuranceType = "premium"
	InsuranceTypeNone     InsuranceType = "none"
)

type Rental struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Tenant            User           `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	ShipID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"ship_id"`
	Ship              Ship           `gorm:"foreignKey:ShipID" json:"ship,omitempty"`
	RentalType        RentalType     `gorm:"type:varchar(20);not null" json:"rental_type"`
	StartDate         time.Time      `gorm:"not null" json:"start_date"`
	EndDate           time.Time      `gorm:"not null" json:"end_date"`
	StartLocation     string         `gorm:"size:200" json:"start_location"`
	EndLocation       string         `gorm:"size:200" json:"end_location"`
	BaseAmount        float64        `gorm:"type:decimal(12,2);not null" json:"base_amount"`
	InsuranceType     InsuranceType  `gorm:"type:varchar(20);default:none" json:"insurance_type"`
	InsuranceAmount   float64        `gorm:"type:decimal(12,2);default:0" json:"insurance_amount"`
	PlatformFee       float64        `gorm:"type:decimal(12,2);default:0" json:"platform_fee"`
	TotalAmount       float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	DepositAmount     float64        `gorm:"type:decimal(12,2);default:0" json:"deposit_amount"`
	DepositStatus     string         `gorm:"type:varchar(20);default:unpaid" json:"deposit_status"`
	Currency          string         `gorm:"size:10;default:USD" json:"currency"`
	Status            RentalStatus   `gorm:"type:varchar(20);default:pending" json:"status"`
	EmergencyContact  string         `gorm:"size:200" json:"emergency_contact"`
	EmergencyPhone    string         `gorm:"size:20" json:"emergency_phone"`
	Notes             string         `gorm:"type:text" json:"notes"`
	CrewCount         int            `gorm:"default:0" json:"crew_count"`
	PassengerCount    int            `gorm:"default:0" json:"passenger_count"`
	CancellationReason string        `gorm:"type:text" json:"cancellation_reason"`
	CancelledAt       *time.Time     `json:"cancelled_at"`
	CompletedAt       *time.Time     `json:"completed_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	VoyageLogs        []VoyageLog    `gorm:"foreignKey:RentalID" json:"voyage_logs,omitempty"`
}

func (r *Rental) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type CreateRentalRequest struct {
	ShipID          string        `json:"ship_id" binding:"required,uuid"`
	RentalType      RentalType    `json:"rental_type" binding:"required,oneof=daily hourly voyage"`
	StartDate       time.Time     `json:"start_date" binding:"required"`
	EndDate         time.Time     `json:"end_date" binding:"required"`
	StartLocation   string        `json:"start_location"`
	EndLocation     string        `json:"end_location"`
	InsuranceType   InsuranceType `json:"insurance_type" binding:"required,oneof=basic premium none"`
	EmergencyContact string       `json:"emergency_contact" binding:"required"`
	EmergencyPhone  string        `json:"emergency_phone" binding:"required"`
	Notes           string        `json:"notes"`
	CrewCount       int           `json:"crew_count"`
	PassengerCount  int           `json:"passenger_count"`
	Currency        string        `json:"currency"`
}

type UpdateRentalStatusRequest struct {
	Status             RentalStatus `json:"status" binding:"required,oneof=confirmed active completed cancelled refunded"`
	CancellationReason string       `json:"cancellation_reason"`
}

type PaymentRecord struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RentalID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"rental_id"`
	Rental          Rental         `gorm:"foreignKey:RentalID" json:"rental,omitempty"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Amount          float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string         `gorm:"size:10;default:USD" json:"currency"`
	PaymentType     string         `gorm:"type:varchar(20);not null" json:"payment_type"`
	TransactionID   string         `gorm:"size:100" json:"transaction_id"`
	Status          string         `gorm:"type:varchar(20);default:pending" json:"status"`
	PaymentMethod   string         `gorm:"type:varchar(50)" json:"payment_method"`
	RefundedAt      *time.Time     `json:"refunded_at"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *PaymentRecord) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
