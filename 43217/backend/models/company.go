package models

import (
	"time"

	"gorm.io/gorm"
)

type CompanyStatus int

const (
	CompanyStatusPending  CompanyStatus = 0
	CompanyStatusActive   CompanyStatus = 1
	CompanyStatusDisabled CompanyStatus = 2
)

type Company struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"size:100;not null" json:"name"`
	UnifiedCode      string         `gorm:"uniqueIndex;size:50" json:"unified_code"`
	LegalPerson      string         `gorm:"size:50" json:"legal_person"`
	ContactPhone     string         `gorm:"size:20" json:"contact_phone"`
	ContactEmail     string         `gorm:"size:100" json:"contact_email"`
	Address          string         `gorm:"size:255" json:"address"`
	BusinessLicense  string         `gorm:"size:255" json:"business_license"`
	Status           CompanyStatus  `gorm:"default:0" json:"status"`
	AnnualBudget     float64        `gorm:"default:0" json:"annual_budget"`
	UsedBudget       float64        `gorm:"default:0" json:"used_budget"`
	PaymentType      int            `gorm:"default:1" json:"payment_type"`
	Balance          float64        `gorm:"default:0" json:"balance"`
	CreditLimit      float64        `gorm:"default:0" json:"credit_limit"`
	HRUserID         *uint          `gorm:"index" json:"hr_user_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	Departments      []Department   `gorm:"foreignKey:CompanyID" json:"departments,omitempty"`
	Employees        []Employee     `gorm:"foreignKey:CompanyID" json:"employees,omitempty"`
}

type Department struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CompanyID   uint           `gorm:"index;not null" json:"company_id"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	ManagerName string         `gorm:"size:50" json:"manager_name"`
	ManagerPhone string        `gorm:"size:20" json:"manager_phone"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Company     Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Employees   []Employee     `gorm:"foreignKey:DepartmentID" json:"employees,omitempty"`
}

type Employee struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CompanyID     uint           `gorm:"index;not null" json:"company_id"`
	DepartmentID  uint           `gorm:"index;not null" json:"department_id"`
	UserID        *uint          `gorm:"index" json:"user_id"`
	EmployeeNo    string         `gorm:"size:50" json:"employee_no"`
	RealName      string         `gorm:"size:50;not null" json:"real_name"`
	Gender        int            `gorm:"default:0" json:"gender"`
	Birthday      *time.Time     `json:"birthday"`
	IDCard        string         `gorm:"size:20" json:"id_card"`
	Phone         string         `gorm:"size:20" json:"phone"`
	Email         string         `gorm:"size:100" json:"email"`
	Position      string         `gorm:"size:50" json:"position"`
	EntryDate     *time.Time     `json:"entry_date"`
	Status        int            `gorm:"default:1" json:"status"`
	Quota         int            `gorm:"default:1" json:"quota"`
	UsedQuota     int            `gorm:"default:0" json:"used_quota"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Company       Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Department    Department     `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Appointments  []Appointment  `gorm:"foreignKey:EmployeeID" json:"appointments,omitempty"`
	HealthRecords []HealthRecord `gorm:"foreignKey:EmployeeID" json:"health_records,omitempty"`
}
