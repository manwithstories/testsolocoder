package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaintenanceType string

const (
	MaintenanceTypeRoutine    MaintenanceType = "routine"
	MaintenanceTypeRepair     MaintenanceType = "repair"
	MaintenanceTypeInspection MaintenanceType = "inspection"
	MaintenanceTypeOverhaul   MaintenanceType = "overhaul"
)

type MaintenanceStatus string

const (
	MaintenanceStatusScheduled MaintenanceStatus = "scheduled"
	MaintenanceStatusInProgress MaintenanceStatus = "in_progress"
	MaintenanceStatusCompleted MaintenanceStatus = "completed"
	MaintenanceStatusCancelled MaintenanceStatus = "cancelled"
)

type MaintenanceRecord struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ShipID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"ship_id"`
	Ship            Ship           `gorm:"foreignKey:ShipID" json:"ship,omitempty"`
	MaintenanceType MaintenanceType `gorm:"type:varchar(20);not null" json:"maintenance_type"`
	Title           string         `gorm:"size:200;not null" json:"title" binding:"required"`
	Description     string         `gorm:"type:text" json:"description"`
	Status          MaintenanceStatus `gorm:"type:varchar(20);default:scheduled" json:"status"`
	PlannedDate     time.Time      `json:"planned_date"`
	StartDate       *time.Time     `json:"start_date"`
	CompletedDate   *time.Time     `json:"completed_date"`
	Cost            float64        `gorm:"type:decimal(12,2);default:0" json:"cost"`
	Currency        string         `gorm:"size:10;default:USD" json:"currency"`
	Provider        string         `gorm:"size:200" json:"provider"`
	Technician      string         `gorm:"size:200" json:"technician"`
	NextDueDate     *time.Time     `json:"next_due_date"`
	Priority        int            `gorm:"default:0" json:"priority"`
	Attachments     string         `gorm:"type:text" json:"attachments"`
	ReminderSent    bool           `gorm:"default:false" json:"reminder_sent"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type MaintenanceSchedule struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ShipID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"ship_id"`
	Ship          Ship           `gorm:"foreignKey:ShipID" json:"ship,omitempty"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	IntervalDays  int            `gorm:"not null;default:90" json:"interval_days"`
	LastCompleted *time.Time     `json:"last_completed"`
	NextDue       time.Time      `gorm:"not null" json:"next_due"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *MaintenanceRecord) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (ms *MaintenanceSchedule) BeforeCreate(tx *gorm.DB) error {
	if ms.ID == uuid.Nil {
		ms.ID = uuid.New()
	}
	return nil
}

type CreateMaintenanceRequest struct {
	ShipID          string          `json:"ship_id" binding:"required,uuid"`
	MaintenanceType MaintenanceType `json:"maintenance_type" binding:"required,oneof=routine repair inspection overhaul"`
	Title           string          `json:"title" binding:"required,min=2,max=200"`
	Description     string          `json:"description"`
	PlannedDate     time.Time       `json:"planned_date" binding:"required"`
	Cost            float64         `json:"cost"`
	Currency        string          `json:"currency"`
	Provider        string          `json:"provider"`
	Technician      string          `json:"technician"`
	NextDueDate     *time.Time      `json:"next_due_date"`
	Priority        int             `json:"priority"`
	Notes           string          `json:"notes"`
}

type UpdateMaintenanceRequest struct {
	Status        MaintenanceStatus `json:"status"`
	StartDate     *time.Time        `json:"start_date"`
	CompletedDate *time.Time        `json:"completed_date"`
	Cost          *float64          `json:"cost"`
	Provider      string            `json:"provider"`
	Technician    string            `json:"technician"`
	Notes         string            `json:"notes"`
}

type CreateMaintenanceScheduleRequest struct {
	ShipID        string    `json:"ship_id" binding:"required,uuid"`
	Title         string    `json:"title" binding:"required"`
	Description   string    `json:"description"`
	IntervalDays  int       `json:"interval_days" binding:"required,min=1"`
	NextDue       time.Time `json:"next_due" binding:"required"`
	IsActive      bool      `json:"is_active"`
}
