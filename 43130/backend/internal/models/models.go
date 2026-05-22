package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password     string    `json:"-" gorm:"size:255;not null"`
	FullName     string    `json:"full_name" gorm:"size:100"`
	Phone        string    `json:"phone" gorm:"size:20"`
	Avatar       string    `json:"avatar" gorm:"size:255"`
	Role         string    `json:"role" gorm:"size:20;default:couple"`
	Status       string    `json:"status" gorm:"size:20;default:active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Wedding struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"index;not null"`
	Title         string    `json:"title" gorm:"size:200;not null"`
	GroomName     string    `json:"groom_name" gorm:"size:100"`
	BrideName     string    `json:"bride_name" gorm:"size:100"`
	WeddingDate   time.Time `json:"wedding_date"`
	Budget        float64   `json:"budget" gorm:"type:decimal(12,2);default:0"`
	Style         string    `json:"style" gorm:"size:100"`
	ThemeColor    string    `json:"theme_color" gorm:"size:50"`
	GuestCount    int       `json:"guest_count" gorm:"default:0"`
	Venue         string    `json:"venue" gorm:"size:200"`
	VenueAddress  string    `json:"venue_address" gorm:"size:500"`
	VenueCapacity int       `json:"venue_capacity" gorm:"default:0"`
	Description   string    `json:"description" gorm:"type:text"`
	Status        string    `json:"status" gorm:"size:20;default:planning"`
	IsTemplate    bool      `json:"is_template" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Vendor struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"index;not null"`
	WeddingID     *uint     `json:"wedding_id" gorm:"index"`
	Name          string    `json:"name" gorm:"size:200;not null"`
	Category      string    `json:"category" gorm:"size:100;not null"`
	ContactPerson string    `json:"contact_person" gorm:"size:100"`
	Phone         string    `json:"phone" gorm:"size:20"`
	Email         string    `json:"email" gorm:"size:100"`
	Address       string    `json:"address" gorm:"size:500"`
	Website       string    `json:"website" gorm:"size:200"`
	ServiceArea   string    `json:"service_area" gorm:"size:500"`
	PriceRange    string    `json:"price_range" gorm:"size:100"`
	Rating        float64   `json:"rating" gorm:"type:decimal(2,1);default:0"`
	ReviewCount   int       `json:"review_count" gorm:"default:0"`
	Notes         string    `json:"notes" gorm:"type:text"`
	Status        string    `json:"status" gorm:"size:20;default:active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type VendorReview struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	VendorID  uint      `json:"vendor_id" gorm:"index;not null"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Rating    int       `json:"rating" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

type Guest struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WeddingID   uint      `json:"wedding_id" gorm:"index;not null"`
	FirstName   string    `json:"first_name" gorm:"size:50;not null"`
	LastName    string    `json:"last_name" gorm:"size:50;not null"`
	FullName    string    `json:"full_name" gorm:"size:100"`
	Email       string    `json:"email" gorm:"size:100"`
	Phone       string    `json:"phone" gorm:"size:20"`
	Group       string    `json:"group" gorm:"size:50"`
	Relation    string    `json:"relation" gorm:"size:50"`
	RSVPStatus  string    `json:"rsvp_status" gorm:"size:20;default:pending"`
	PlusOne     bool      `json:"plus_one" gorm:"default:false"`
	PlusOneName string    `json:"plus_one_name" gorm:"size:100"`
	TableID     *uint     `json:"table_id" gorm:"index"`
	SeatNumber  int       `json:"seat_number" gorm:"default:0"`
	Notes       string    `json:"notes" gorm:"type:text"`
	IsVIP       bool      `json:"is_vip" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GuestTable struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WeddingID   uint      `json:"wedding_id" gorm:"index;not null"`
	TableName   string    `json:"table_name" gorm:"size:100;not null"`
	TableNumber int       `json:"table_number" gorm:"default:0"`
	Capacity    int       `json:"capacity" gorm:"default:10"`
	SeatsJSON   string    `json:"seats_json" gorm:"type:text"`
	Notes       string    `json:"notes" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BudgetItem struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	WeddingID     uint      `json:"wedding_id" gorm:"index;not null"`
	VendorID      *uint     `json:"vendor_id" gorm:"index"`
	Category      string    `json:"category" gorm:"size:100;not null"`
	Description   string    `json:"description" gorm:"size:500"`
	EstimatedCost float64   `json:"estimated_cost" gorm:"type:decimal(12,2);default:0"`
	ActualCost    float64   `json:"actual_cost" gorm:"type:decimal(12,2);default:0"`
	PaidAmount    float64   `json:"paid_amount" gorm:"type:decimal(12,2);default:0"`
	Status        string    `json:"status" gorm:"size:20;default:pending"`
	DueDate       *time.Time `json:"due_date"`
	Notes         string    `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Payment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BudgetItemID uint    `json:"budget_item_id" gorm:"index;not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(12,2);not null"`
	Method      string    `json:"method" gorm:"size:50"`
	Status      string    `json:"status" gorm:"size:20;default:pending"`
	PaidAt      *time.Time `json:"paid_at"`
	Reference   string    `json:"reference" gorm:"size:100"`
	Notes       string    `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WeddingID   uint      `json:"wedding_id" gorm:"index;not null"`
	Title       string    `json:"title" gorm:"size:200;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category" gorm:"size:50"`
	Assignee    string    `json:"assignee" gorm:"size:100"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string    `json:"priority" gorm:"size:20;default:medium"`
	Status      string    `json:"status" gorm:"size:20;default:pending"`
	CompletedAt *time.Time `json:"completed_at"`
	ParentID    *uint     `json:"parent_id" gorm:"index"`
	TemplateID  *uint     `json:"template_id" gorm:"index"`
	IsTemplate  bool      `json:"is_template" gorm:"default:false"`
	Order       int       `json:"order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:200;not null"`
	Category  string    `json:"category" gorm:"size:50"`
	TasksJSON string    `json:"tasks_json" gorm:"type:text"`
	IsDefault bool      `json:"is_default" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

type Document struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WeddingID   uint      `json:"wedding_id" gorm:"index;not null"`
	VendorID    *uint     `json:"vendor_id" gorm:"index"`
	BudgetItemID *uint    `json:"budget_item_id" gorm:"index"`
	FileName    string    `json:"file_name" gorm:"size:255;not null"`
	FilePath    string    `json:"file_path" gorm:"size:500;not null"`
	FileSize    int64     `json:"file_size"`
	FileType    string    `json:"file_type" gorm:"size:50"`
	Category    string    `json:"category" gorm:"size:50"`
	Version     int       `json:"version" gorm:"default:1"`
	ParentID    *uint     `json:"parent_id" gorm:"index"`
	UploadedBy  uint      `json:"uploaded_by" gorm:"index"`
	Notes       string    `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OperationLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index"`
	Module      string    `json:"module" gorm:"size:50"`
	Action      string    `json:"action" gorm:"size:50"`
	TargetID    uint      `json:"target_id"`
	Detail      string    `json:"detail" gorm:"type:text"`
	IPAddress   string    `json:"ip_address" gorm:"size:50"`
	CreatedAt   time.Time `json:"created_at"`
}

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Type      string    `json:"type" gorm:"size:50"`
	Title     string    `json:"title" gorm:"size:200"`
	Content   string    `json:"content" gorm:"type:text"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	RelatedID uint      `json:"related_id"`
	CreatedAt time.Time `json:"created_at"`
}
