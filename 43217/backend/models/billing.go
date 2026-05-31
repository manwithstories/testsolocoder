package models

import (
	"time"

	"gorm.io/gorm"
)

type BillingStatus int

const (
	BillingStatusPending  BillingStatus = 0
	BillingStatusPaid     BillingStatus = 1
	BillingStatusOverdue  BillingStatus = 2
	BillingStatusCancelled BillingStatus = 3
)

type Billing struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	BillingNo   string         `gorm:"uniqueIndex;size:50;not null" json:"billing_no"`
	CompanyID   uint           `gorm:"index;not null" json:"company_id"`
	AgencyID    uint           `gorm:"index;not null" json:"agency_id"`
	Period      string         `gorm:"size:20;not null" json:"period"`
	TotalAmount float64        `gorm:"default:0" json:"total_amount"`
	PaidAmount  float64        `gorm:"default:0" json:"paid_amount"`
	Status      BillingStatus  `gorm:"default:0" json:"status"`
	DueDate     *time.Time     `json:"due_date"`
	PaidAt      *time.Time     `json:"paid_at"`
	InvoiceNo   string         `gorm:"size:50" json:"invoice_no"`
	InvoiceFile string         `gorm:"size:255" json:"invoice_file"`
	InvoiceStatus int          `gorm:"default:0" json:"invoice_status"`
	Remark      string         `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Company     Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Agency      Agency         `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	Items       []BillingItem  `gorm:"foreignKey:BillingID" json:"items,omitempty"`
}

type BillingItem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	BillingID     uint           `gorm:"index;not null" json:"billing_id"`
	AppointmentID uint           `gorm:"index;not null" json:"appointment_id"`
	EmployeeID    uint           `gorm:"index;not null" json:"employee_id"`
	PackageID     uint           `gorm:"index;not null" json:"package_id"`
	PackageName   string         `gorm:"size:100" json:"package_name"`
	UnitPrice     float64        `gorm:"default:0" json:"unit_price"`
	Quantity      int            `gorm:"default:1" json:"quantity"`
	Amount        float64        `gorm:"default:0" json:"amount"`
	AppointmentDate time.Time    `json:"appointment_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Billing       Billing        `gorm:"foreignKey:BillingID" json:"billing,omitempty"`
}

type Transaction struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TransactionNo string         `gorm:"uniqueIndex;size:50;not null" json:"transaction_no"`
	CompanyID     uint           `gorm:"index;not null" json:"company_id"`
	Type          string         `gorm:"size:20;not null" json:"type"`
	Amount        float64        `gorm:"default:0" json:"amount"`
	Balance       float64        `gorm:"default:0" json:"balance"`
	PaymentMethod string         `gorm:"size:20" json:"payment_method"`
	Status        int            `gorm:"default:1" json:"status"`
	Remark        string         `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Company       Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

type CompanyBudget struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CompanyID   uint           `gorm:"index;not null" json:"company_id"`
	Year        int            `gorm:"index;not null" json:"year"`
	TotalBudget float64        `gorm:"default:0" json:"total_budget"`
	UsedBudget  float64        `gorm:"default:0" json:"used_budget"`
	Frequency   int            `gorm:"default:1" json:"frequency"`
	PeriodStart *time.Time     `json:"period_start"`
	PeriodEnd   *time.Time     `json:"period_end"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Company     Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

type DepartmentAppointment struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CompanyID      uint           `gorm:"index;not null" json:"company_id"`
	DepartmentID   uint           `gorm:"index;not null" json:"department_id"`
	AgencyID       uint           `gorm:"index;not null" json:"agency_id"`
	PackageID      uint           `gorm:"index;not null" json:"package_id"`
	Year           int            `gorm:"index;not null" json:"year"`
	TotalQuota     int            `gorm:"default:0" json:"total_quota"`
	UsedQuota      int            `gorm:"default:0" json:"used_quota"`
	StartDate      *time.Time     `json:"start_date"`
	EndDate        *time.Time     `json:"end_date"`
	Status         int            `gorm:"default:1" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Company        Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Department     Department     `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Agency         Agency         `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	Package        Package        `gorm:"foreignKey:PackageID" json:"package,omitempty"`
}
