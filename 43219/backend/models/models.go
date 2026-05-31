package models

import "time"

type Role string

const (
	RoleCompany  Role = "company"
	RoleStaff    Role = "staff"
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username" gorm:"uniqueIndex;size:64;not null"`
	Password    string    `json:"-" gorm:"size:255;not null"`
	RealName    string    `json:"real_name" gorm:"size:64"`
	Phone       string    `json:"phone" gorm:"size:32"`
	Role        Role      `json:"role" gorm:"size:32;not null;index"`
	CompanyID   *uint     `json:"company_id,omitempty" gorm:"index"`
	Company     *User     `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Avatar      string    `json:"avatar" gorm:"size:255"`
	Status      string    `json:"status" gorm:"size:32;default:active"`
	Rating      float64   `json:"rating" gorm:"default:5"`
	Level       int       `json:"level" gorm:"default:1"`
	Suspended   bool      `json:"suspended" gorm:"default:false"`
	Skills      string    `json:"skills" gorm:"size:512"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StaffProfile struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	User         *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	IDCard       string    `json:"id_card" gorm:"size:64"`
	CertFiles    string    `json:"cert_files" gorm:"size:1024"`
	HealthFiles  string    `json:"health_files" gorm:"size:1024"`
	CertVerified bool      `json:"cert_verified" gorm:"default:false"`
	Intro        string    `json:"intro" gorm:"size:1024"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Service struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CompanyID uint      `json:"company_id" gorm:"index;not null"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	Category  string    `json:"category" gorm:"size:64;index"`
	Desc      string    `json:"desc" gorm:"size:1024"`
	MinPrice  float64   `json:"min_price"`
	MaxPrice  float64   `json:"max_price"`
	Duration  int       `json:"duration"`
	Skills    string    `json:"skills" gorm:"size:512"`
	Status    string    `json:"status" gorm:"size:32;default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCanceled  BookingStatus = "canceled"
	BookingRejected  BookingStatus = "rejected"
	BookingCompleted BookingStatus = "completed"
)

type Booking struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	CustomerID      uint          `json:"customer_id" gorm:"index;not null"`
	Customer        *User         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	ServiceID       uint          `json:"service_id" gorm:"index;not null"`
	Service         *Service      `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	StaffID         *uint         `json:"staff_id,omitempty" gorm:"index"`
	Staff           *User         `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
	StartAt         time.Time     `json:"start_at"`
	EndAt           time.Time     `json:"end_at"`
	Address         string        `json:"address" gorm:"size:255"`
	Remark          string        `json:"remark" gorm:"size:512"`
	Price           float64       `json:"price"`
	Status          BookingStatus `json:"status" gorm:"size:32;index;default:pending"`
	RescheduleCount int           `json:"reschedule_count" gorm:"default:0"`
	NeedReview      bool          `json:"need_review" gorm:"default:false"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type OrderStatus string

const (
	OrderCreated   OrderStatus = "created"
	OrderReported  OrderStatus = "reported"
	OrderConfirmed OrderStatus = "confirmed"
	OrderPaid      OrderStatus = "paid"
	OrderRefunding OrderStatus = "refunding"
	OrderRefunded  OrderStatus = "refunded"
	OrderClosed    OrderStatus = "closed"
)

type Order struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	BookingID    uint        `json:"booking_id" gorm:"index;not null"`
	Booking      *Booking    `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	CustomerID   uint        `json:"customer_id" gorm:"index;not null"`
	StaffID      uint        `json:"staff_id" gorm:"index;not null"`
	CompanyID    uint        `json:"company_id" gorm:"index;not null"`
	TotalAmount  float64     `json:"total_amount"`
	Status       OrderStatus `json:"status" gorm:"size:32;index;default:created"`
	ReportText   string      `json:"report_text" gorm:"size:2000"`
	ReportImages string      `json:"report_images" gorm:"size:2000"`
	ReportedAt   *time.Time  `json:"reported_at,omitempty"`
	ConfirmedAt  *time.Time  `json:"confirmed_at,omitempty"`
	PaidAt       *time.Time  `json:"paid_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type Review struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	OrderID    uint      `json:"order_id" gorm:"index;not null"`
	StaffID    uint      `json:"staff_id" gorm:"index;not null"`
	Staff      *User     `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
	CustomerID uint      `json:"customer_id" gorm:"index;not null"`
	Rating     int       `json:"rating" gorm:"not null"`
	Content    string    `json:"content" gorm:"size:2000"`
	Images     string    `json:"images" gorm:"size:2000"`
	CreatedAt  time.Time `json:"created_at"`
}

type TicketStatus string

const (
	TicketOpen     TicketStatus = "open"
	TicketAssigned TicketStatus = "assigned"
	TicketPending  TicketStatus = "pending"
	TicketResolved TicketStatus = "resolved"
	TicketClosed   TicketStatus = "closed"
	TicketEscalate TicketStatus = "escalated"
)

type TicketType string

const (
	TicketComplaint TicketType = "complaint"
	TicketRefund    TicketType = "refund"
)

type Ticket struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	OrderID      *uint        `json:"order_id,omitempty" gorm:"index"`
	CustomerID   uint         `json:"customer_id" gorm:"index;not null"`
	StaffID      *uint        `json:"staff_id,omitempty" gorm:"index"`
	AgentID      *uint        `json:"agent_id,omitempty" gorm:"index"`
	Type         TicketType   `json:"type" gorm:"size:32;index"`
	Title        string       `json:"title" gorm:"size:255"`
	Content      string       `json:"content" gorm:"size:2000"`
	Status       TicketStatus `json:"status" gorm:"size:32;index;default:open"`
	Escalated    bool         `json:"escalated" gorm:"default:false"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	LastActionAt time.Time    `json:"last_action_at"`
}

type Wallet struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	Balance   float64   `json:"balance" gorm:"default:0"`
	Frozen    float64   `json:"frozen" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Settlement struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	OrderID      uint      `json:"order_id" gorm:"index;not null"`
	CompanyID    uint      `json:"company_id" gorm:"index;not null"`
	StaffID      uint      `json:"staff_id" gorm:"index;not null"`
	TotalAmount  float64   `json:"total_amount"`
	CompanyShare float64   `json:"company_share"`
	StaffShare   float64   `json:"staff_share"`
	Status       string    `json:"status" gorm:"size:32;default:paid"`
	CreatedAt    time.Time `json:"created_at"`
}

type Withdrawal struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StaffID   uint      `json:"staff_id" gorm:"index;not null"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status" gorm:"size:32;default:pending"`
	Account   string    `json:"account" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
