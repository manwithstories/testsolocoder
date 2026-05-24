package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleTech     UserRole = "technician"
	RoleAdmin    UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusPending  UserStatus = "pending"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Phone        string         `gorm:"size:20" json:"phone"`
	Email        string         `gorm:"size:100" json:"email"`
	RealName     string         `gorm:"size:50" json:"real_name"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Role         UserRole       `gorm:"size:20;not null" json:"role"`
	Status       UserStatus     `gorm:"size:20;default:active" json:"status"`
	Balance      float64        `gorm:"default:0" json:"balance"`
	Address      string         `gorm:"size:255" json:"address"`
	Longitude    float64        `json:"longitude"`
	Latitude     float64        `json:"latitude"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type TechnicianProfile struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	User               User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CertificateImage   string         `gorm:"size:255" json:"certificate_image"`
	CertificateNo      string         `gorm:"size:50" json:"certificate_no"`
	Specialty          string         `gorm:"size:255" json:"specialty"`
	ExperienceYears    int            `gorm:"default:0" json:"experience_years"`
	Rating             float64        `gorm:"default:5.0" json:"rating"`
	TotalOrders        int            `gorm:"default:0" json:"total_orders"`
	CompletedOrders    int            `gorm:"default:0" json:"completed_orders"`
	ActiveOrders       int            `gorm:"default:0" json:"active_orders"`
	MaxActiveOrders    int            `gorm:"default:5" json:"max_active_orders"`
	IsVerified         bool           `gorm:"default:false" json:"is_verified"`
	VerifyStatus       string         `gorm:"size:20;default:pending" json:"verify_status"`
	VerifyRemark       string         `gorm:"size:255" json:"verify_remark"`
	ServiceRadius      float64        `gorm:"default:50.0" json:"service_radius"`
	BankAccount        string         `gorm:"size:100" json:"bank_account"`
	BankName           string         `gorm:"size:50" json:"bank_name"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Code      string         `gorm:"uniqueIndex;size:20;not null" json:"code"`
	Icon      string         `gorm:"size:255" json:"icon"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    bool           `gorm:"default:true" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ServiceItem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CategoryID    uint           `gorm:"not null" json:"category_id"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Description   string         `gorm:"size:500" json:"description"`
	MinPrice      float64        `gorm:"not null" json:"min_price"`
	MaxPrice      float64        `gorm:"not null" json:"max_price"`
	EstimatedTime int            `gorm:"not null;comment:预计时长（分钟）" json:"estimated_time"`
	Image         string         `gorm:"size:255" json:"image"`
	Sort          int            `gorm:"default:0" json:"sort"`
	Status        bool           `gorm:"default:true" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAssigned   OrderStatus = "assigned"
	OrderStatusAccepted   OrderStatus = "accepted"
	OrderStatusOnSite     OrderStatus = "on_site"
	OrderStatusRepairing  OrderStatus = "repairing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunding  OrderStatus = "refunding"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type Order struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	OrderNo          string         `gorm:"uniqueIndex;size:30;not null" json:"order_no"`
	CustomerID       uint           `gorm:"not null" json:"customer_id"`
	Customer         User           `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	TechnicianID     *uint          `json:"technician_id"`
	Technician       *User          `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	ServiceItemID    uint           `gorm:"not null" json:"service_item_id"`
	ServiceItem      ServiceItem    `gorm:"foreignKey:ServiceItemID" json:"service_item,omitempty"`
	Title            string         `gorm:"size:200;not null" json:"title"`
	Description      string         `gorm:"size:1000" json:"description"`
	Images           string         `gorm:"size:1000" json:"images"`
	Address          string         `gorm:"size:255;not null" json:"address"`
	Longitude        float64        `json:"longitude"`
	Latitude         float64        `json:"latitude"`
	ContactName      string         `gorm:"size:50;not null" json:"contact_name"`
	ContactPhone     string         `gorm:"size:20;not null" json:"contact_phone"`
	AppointmentTime  *time.Time     `json:"appointment_time"`
	QuotedPrice      float64        `gorm:"default:0" json:"quoted_price"`
	FinalPrice       float64        `gorm:"default:0" json:"final_price"`
	Status           OrderStatus    `gorm:"size:20;default:pending" json:"status"`
	UrgentLevel      int            `gorm:"default:0" json:"urgent_level"`
	CancelReason     string         `gorm:"size:500" json:"cancel_reason"`
	RefundReason     string         `gorm:"size:500" json:"refund_reason"`
	RefundAmount     float64        `gorm:"default:0" json:"refund_amount"`
	AssignedAt       *time.Time     `json:"assigned_at"`
	AcceptedAt       *time.Time     `json:"accepted_at"`
	ArrivedAt        *time.Time     `json:"arrived_at"`
	StartedAt        *time.Time     `json:"started_at"`
	CompletedAt      *time.Time     `json:"completed_at"`
	CancelledAt      *time.Time     `json:"cancelled_at"`
	RefundedAt       *time.Time     `json:"refunded_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"not null" json:"order_id"`
	Order     Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	UserID    uint      `json:"user_id"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	Content   string    `gorm:"size:500" json:"content"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderID       uint           `gorm:"uniqueIndex;not null" json:"order_id"`
	Order         Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	CustomerID    uint           `gorm:"not null" json:"customer_id"`
	Customer      User           `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	TechnicianID  uint           `gorm:"not null" json:"technician_id"`
	Technician    User           `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	Rating        int            `gorm:"not null" json:"rating"`
	Content       string         `gorm:"size:1000" json:"content"`
	Images        string         `gorm:"size:1000" json:"images"`
	Reply         string         `gorm:"size:500" json:"reply"`
	RepliedAt     *time.Time     `json:"replied_at"`
	IsIntervened  bool           `gorm:"default:false" json:"is_intervened"`
	InterveneNote string         `gorm:"size:500" json:"intervene_note"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Part struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Category    string         `gorm:"size:50" json:"category"`
	Description string         `gorm:"size:500" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	MinStock    int            `gorm:"default:10" json:"min_stock"`
	Image       string         `gorm:"size:255" json:"image"`
	Status      bool           `gorm:"default:true" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type PartRequest struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RequestNo    string         `gorm:"uniqueIndex;size:30;not null" json:"request_no"`
	TechnicianID uint           `gorm:"not null" json:"technician_id"`
	Technician   User           `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	Status       string         `gorm:"size:20;default:pending" json:"status"`
	TotalAmount  float64        `gorm:"default:0" json:"total_amount"`
	Remark       string         `gorm:"size:500" json:"remark"`
	ApprovedAt   *time.Time     `json:"approved_at"`
	ShippedAt    *time.Time     `json:"shipped_at"`
	ReceivedAt   *time.Time     `json:"received_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type PartRequestItem struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	PartRequestID uint        `gorm:"not null" json:"part_request_id"`
	PartRequest   PartRequest `gorm:"foreignKey:PartRequestID" json:"part_request,omitempty"`
	PartID        uint        `gorm:"not null" json:"part_id"`
	Part          Part        `gorm:"foreignKey:PartID" json:"part,omitempty"`
	Quantity      int         `gorm:"not null" json:"quantity"`
	Price         float64     `gorm:"not null" json:"price"`
}

type PartUsage struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"not null" json:"order_id"`
	Order         Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	PartID        uint      `gorm:"not null" json:"part_id"`
	Part          Part      `gorm:"foreignKey:PartID" json:"part,omitempty"`
	TechnicianID  uint      `gorm:"not null" json:"technician_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"not null" json:"price"`
	Note          string    `gorm:"size:500" json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

type WithdrawRequest struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	RequestNo     string         `gorm:"uniqueIndex;size:30;not null" json:"request_no"`
	TechnicianID  uint           `gorm:"not null" json:"technician_id"`
	Technician    User           `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	Amount        float64        `gorm:"not null" json:"amount"`
	BankAccount   string         `gorm:"size:100" json:"bank_account"`
	BankName      string         `gorm:"size:50" json:"bank_name"`
	Status        string         `gorm:"size:20;default:pending" json:"status"`
	Remark        string         `gorm:"size:500" json:"remark"`
	ApprovedAt    *time.Time     `json:"approved_at"`
	TransferredAt *time.Time     `json:"transferred_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TransactionNo string         `gorm:"uniqueIndex;size:30;not null" json:"transaction_no"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type          string         `gorm:"size:20;not null" json:"type"`
	Amount        float64        `gorm:"not null" json:"amount"`
	BalanceAfter  float64        `json:"balance_after"`
	OrderID       *uint          `json:"order_id"`
	Order         *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Description   string         `gorm:"size:500" json:"description"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type MonthlyReport struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Month           string    `gorm:"size:7;uniqueIndex;not null" json:"month"`
	TotalOrders     int       `gorm:"default:0" json:"total_orders"`
	TotalRevenue    float64   `gorm:"default:0" json:"total_revenue"`
	PlatformIncome  float64   `gorm:"default:0" json:"platform_income"`
	TechnicianPay   float64   `gorm:"default:0" json:"technician_pay"`
	TotalWithdraw   float64   `gorm:"default:0" json:"total_withdraw"`
	NewTechnicians  int       `gorm:"default:0" json:"new_technicians"`
	NewCustomers    int       `gorm:"default:0" json:"new_customers"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
