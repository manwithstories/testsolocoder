package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	FirstName    string         `gorm:"size:100" json:"first_name"`
	LastName     string         `gorm:"size:100" json:"last_name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Clients     []Client     `gorm:"foreignKey:UserID" json:"-"`
	Projects    []Project    `gorm:"foreignKey:UserID" json:"-"`
	TimeEntries []TimeEntry  `gorm:"foreignKey:UserID" json:"-"`
	Invoices    []Invoice    `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type Client struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Email         string         `gorm:"size:255;not null" json:"email"`
	Phone         string         `gorm:"size:50" json:"phone"`
	Address       string         `gorm:"size:500" json:"address"`
	Company       string         `gorm:"size:255" json:"company"`
	ContractURL   string         `gorm:"size:500" json:"contract_url"`
	DefaultRate   float64        `gorm:"default:0" json:"default_rate"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Projects []Project `gorm:"foreignKey:ClientID" json:"projects,omitempty"`
}

type ProjectStatus string

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusArchived  ProjectStatus = "archived"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	ClientID    uint           `gorm:"not null;index" json:"client_id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      ProjectStatus  `gorm:"size:50;not null;default:draft" json:"status"`
	HourlyRate  float64        `gorm:"default:0" json:"hourly_rate"`
	Deadline    *time.Time     `json:"deadline"`
	Budget      float64        `gorm:"default:0" json:"budget"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Client      *Client      `json:"client,omitempty"`
	Milestones  []Milestone  `json:"milestones,omitempty"`
	TimeEntries []TimeEntry  `json:"time_entries,omitempty"`
	Invoices    []Invoice    `json:"invoices,omitempty"`
}

type Milestone struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProjectID   uint           `gorm:"not null;index" json:"project_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	DueDate     *time.Time     `json:"due_date"`
	Completed   bool           `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type TimeEntry struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	ProjectID   uint           `gorm:"not null;index" json:"project_id"`
	Date        time.Time      `gorm:"not null;index" json:"date"`
	Hours       float64        `gorm:"not null" json:"hours"`
	Description string         `gorm:"type:text" json:"description"`
	StartTime   *time.Time     `json:"start_time"`
	EndTime     *time.Time     `json:"end_time"`
	IsTimer     bool           `gorm:"default:false" json:"is_timer"`
	Billable    bool           `gorm:"default:true" json:"billable"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Project *Project `json:"project,omitempty"`
}

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusSent      InvoiceStatus = "sent"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	ClientID      uint           `gorm:"not null;index" json:"client_id"`
	ProjectID     *uint          `gorm:"index" json:"project_id"`
	InvoiceNumber string         `gorm:"uniqueIndex;size:50;not null" json:"invoice_number"`
	Status        InvoiceStatus  `gorm:"size:50;not null;default:draft" json:"status"`
	IssueDate     time.Time      `gorm:"not null" json:"issue_date"`
	DueDate       time.Time      `gorm:"not null" json:"due_date"`
	Subtotal      float64        `gorm:"default:0" json:"subtotal"`
	TaxRate       float64        `gorm:"default:0" json:"tax_rate"`
	TaxAmount     float64        `gorm:"default:0" json:"tax_amount"`
	Total         float64        `gorm:"default:0" json:"total"`
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Client *Client       `json:"client,omitempty"`
	Project *Project     `json:"project,omitempty"`
	Items  []InvoiceItem `json:"items,omitempty"`
}

type InvoiceItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	InvoiceID   uint      `gorm:"not null;index" json:"invoice_id"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	UnitPrice   float64   `gorm:"not null" json:"unit_price"`
	Amount      float64   `gorm:"not null" json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InvoiceCounter struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Year      int       `gorm:"uniqueIndex;not null" json:"year"`
	Count     int       `gorm:"default:0" json:"count"`
	UpdatedAt time.Time `json:"updated_at"`
}
