package models

import (
	"time"

	"gorm.io/gorm"
)

type AppointmentStatus int

const (
	AppointmentStatusPending   AppointmentStatus = 0
	AppointmentStatusConfirmed AppointmentStatus = 1
	AppointmentStatusCompleted AppointmentStatus = 2
	AppointmentStatusCancelled AppointmentStatus = 3
	AppointmentStatusExpired   AppointmentStatus = 4
)

type Appointment struct {
	ID             uint              `gorm:"primaryKey" json:"id"`
	EmployeeID     uint              `gorm:"index;not null" json:"employee_id"`
	CompanyID      uint              `gorm:"index;not null" json:"company_id"`
	AgencyID       uint              `gorm:"index;not null" json:"agency_id"`
	PackageID      uint              `gorm:"index;not null" json:"package_id"`
	TimeSlotID     uint              `gorm:"index;not null" json:"time_slot_id"`
	AppointmentNo  string            `gorm:"uniqueIndex;size:50;not null" json:"appointment_no"`
	AppointmentDate time.Time        `gorm:"not null" json:"appointment_date"`
	StartTime      string            `gorm:"size:10;not null" json:"start_time"`
	EndTime        string            `gorm:"size:10;not null" json:"end_time"`
	Status         AppointmentStatus `gorm:"default:0" json:"status"`
	Remark         string            `gorm:"type:text" json:"remark"`
	IsReminded     bool              `gorm:"default:false" json:"is_reminded"`
	RemindedAt     *time.Time        `json:"reminded_at"`
	CancelledAt    *time.Time        `json:"cancelled_at"`
	CancelReason   string            `gorm:"size:255" json:"cancel_reason"`
	CompletedAt    *time.Time        `json:"completed_at"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `gorm:"index" json:"-"`
	Employee       Employee          `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Company        Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Agency         Agency            `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	Package        Package           `gorm:"foreignKey:PackageID" json:"package,omitempty"`
	TimeSlot       TimeSlot          `gorm:"foreignKey:TimeSlotID" json:"time_slot,omitempty"`
	Report         *Report           `gorm:"foreignKey:AppointmentID" json:"report,omitempty"`
}

type ReportStatus int

const (
	ReportStatusPending  ReportStatus = 0
	ReportStatusUploaded ReportStatus = 1
	ReportStatusViewed   ReportStatus = 2
)

type Report struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppointmentID uint           `gorm:"uniqueIndex;not null" json:"appointment_id"`
	EmployeeID    uint           `gorm:"index;not null" json:"employee_id"`
	AgencyID      uint           `gorm:"index;not null" json:"agency_id"`
	PackageID     uint           `gorm:"index;not null" json:"package_id"`
	ReportNo      string         `gorm:"uniqueIndex;size:50;not null" json:"report_no"`
	ReportDate    time.Time      `gorm:"not null" json:"report_date"`
	DoctorName    string         `gorm:"size:50" json:"doctor_name"`
	Summary       string         `gorm:"type:text" json:"summary"`
	Suggestion    string         `gorm:"type:text" json:"suggestion"`
	HasAbnormal   bool           `gorm:"default:false" json:"has_abnormal"`
	Status        ReportStatus   `gorm:"default:0" json:"status"`
	PdfFile       string         `gorm:"size:255" json:"pdf_file"`
	ViewedAt      *time.Time     `json:"viewed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Appointment   Appointment    `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
	Employee      Employee       `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Items         []ReportItem   `gorm:"foreignKey:ReportID" json:"items,omitempty"`
}

type ReportItem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ReportID      uint           `gorm:"index;not null" json:"report_id"`
	PackageItemID *uint          `gorm:"index" json:"package_item_id"`
	ItemName      string         `gorm:"size:100;not null" json:"item_name"`
	ItemCode      string         `gorm:"size:50" json:"item_code"`
	Result        string         `gorm:"size:100" json:"result"`
	Unit          string         `gorm:"size:20" json:"unit"`
	NormalRange   string         `gorm:"size:255" json:"normal_range"`
	IsAbnormal    bool           `gorm:"default:false" json:"is_abnormal"`
	AbnormalType  string         `gorm:"size:20" json:"abnormal_type"`
	Description   string         `gorm:"type:text" json:"description"`
	Suggestion    string         `gorm:"type:text" json:"suggestion"`
	Department    string         `gorm:"size:50" json:"department"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Report        Report         `gorm:"foreignKey:ReportID" json:"report,omitempty"`
	PackageItem   *PackageItem   `gorm:"foreignKey:PackageItemID" json:"package_item,omitempty"`
}
