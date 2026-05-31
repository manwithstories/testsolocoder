package models

import (
	"time"

	"gorm.io/gorm"
)

type HealthRecord struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	EmployeeID    uint           `gorm:"index;not null" json:"employee_id"`
	CompanyID     uint           `gorm:"index;not null" json:"company_id"`
	RecordYear    int            `gorm:"index;not null" json:"record_year"`
	ReportID      *uint          `gorm:"index" json:"report_id"`
	Height        float64        `json:"height"`
	Weight        float64        `json:"weight"`
	BMI           float64        `json:"bmi"`
	BloodPressure string         `gorm:"size:20" json:"blood_pressure"`
	HeartRate     int            `json:"heart_rate"`
	HasAbnormal   bool           `gorm:"default:false" json:"has_abnormal"`
	AbnormalCount int            `gorm:"default:0" json:"abnormal_count"`
	Tags          string         `gorm:"size:500" json:"tags"`
	Summary       string         `gorm:"type:text" json:"summary"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Employee      Employee       `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Company       Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Report        *Report        `gorm:"foreignKey:ReportID" json:"report,omitempty"`
}

type AbnormalItem struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	EmployeeID     uint           `gorm:"index;not null" json:"employee_id"`
	HealthRecordID uint           `gorm:"index;not null" json:"health_record_id"`
	ItemName       string         `gorm:"size:100;not null" json:"item_name"`
	ItemCode       string         `gorm:"size:50" json:"item_code"`
	Result         string         `gorm:"size:100" json:"result"`
	NormalRange    string         `gorm:"size:255" json:"normal_range"`
	AbnormalType   string         `gorm:"size:20" json:"abnormal_type"`
	Level          int            `gorm:"default:1" json:"level"`
	Suggestion     string         `gorm:"type:text" json:"suggestion"`
	RecheckDate    *time.Time     `json:"recheck_date"`
	RecheckStatus  int            `gorm:"default:0" json:"recheck_status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Employee       Employee       `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	HealthRecord   HealthRecord   `gorm:"foreignKey:HealthRecordID" json:"health_record,omitempty"`
}

type RecheckReminder struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	EmployeeID    uint           `gorm:"index;not null" json:"employee_id"`
	AbnormalID    uint           `gorm:"index;not null" json:"abnormal_id"`
	RemindDate    time.Time      `gorm:"not null" json:"remind_date"`
	RemindType    string         `gorm:"size:50" json:"remind_type"`
	Content       string         `gorm:"type:text" json:"content"`
	IsRead        bool           `gorm:"default:false" json:"is_read"`
	IsSent        bool           `gorm:"default:false" json:"is_sent"`
	SentAt        *time.Time     `json:"sent_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
