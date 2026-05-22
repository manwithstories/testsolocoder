package model

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Phone        string    `json:"phone" gorm:"uniqueIndex;size:20;not null"`
	Password     string    `json:"-" gorm:"size:255;not null"`
	Nickname     string    `json:"nickname" gorm:"size:50"`
	Avatar       string    `json:"avatar" gorm:"size:255"`
	Role         string    `json:"role" gorm:"size:20;default:customer"`
	Status       int       `json:"status" gorm:"default:1"`
	LastLoginAt  time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type Customer struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	User            *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name            string    `json:"name" gorm:"size:50"`
	Gender          string    `json:"gender" gorm:"size:10"`
	Age             int       `json:"age"`
	SkinType        string    `json:"skin_type" gorm:"size:50"`
	HairPreference  string    `json:"hair_preference" gorm:"size:255"`
	AllergyHistory  string    `json:"allergy_history" gorm:"size:500"`
	Notes           string    `json:"notes" gorm:"size:500"`
	MemberLevel     int       `json:"member_level" gorm:"default:1"`
	Points          int       `json:"points" gorm:"default:0"`
	TotalSpent      float64   `json:"total_spent" gorm:"default:0"`
	VisitCount      int       `json:"visit_count" gorm:"default:0"`
	LastVisitAt     time.Time `json:"last_visit_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}

func (Customer) TableName() string {
	return "customers"
}

type Technician struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	User            *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name            string    `json:"name" gorm:"size:50"`
	Title           string    `json:"title" gorm:"size:50"`
	Avatar          string    `json:"avatar" gorm:"size:255"`
	Specialties     string    `json:"specialties" gorm:"type:text"`
	Description     string    `json:"description" gorm:"size:500"`
	Rating          float64   `json:"rating" gorm:"default:5.0"`
	ReviewCount     int       `json:"review_count" gorm:"default:0"`
	WorkStartTime   string    `json:"work_start_time" gorm:"size:10;default:09:00"`
	WorkEndTime     string    `json:"work_end_time" gorm:"size:10;default:21:00"`
	WorkDays        string    `json:"work_days" gorm:"size:50;default:1,2,3,4,5,6"`
	Status          int       `json:"status" gorm:"default:1"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}

func (Technician) TableName() string {
	return "technicians"
}

type TechnicianLeave struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TechnicianID uint      `json:"technician_id" gorm:"index;not null"`
	Technician   *Technician `json:"technician,omitempty" gorm:"foreignKey:TechnicianID"`
	LeaveDate    time.Time `json:"leave_date" gorm:"not null"`
	Reason       string    `json:"reason" gorm:"size:500"`
	Approved     bool      `json:"approved" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (TechnicianLeave) TableName() string {
	return "technician_leaves"
}

type Service struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"size:100;not null"`
	Category      string    `json:"category" gorm:"size:50;index"`
	Description   string    `json:"description" gorm:"size:500"`
	Price         float64   `json:"price" gorm:"not null"`
	Duration      int       `json:"duration" gorm:"not null"`
	RequiredSkill string    `json:"required_skill" gorm:"size:255"`
	Products      string    `json:"products" gorm:"type:text"`
	IsPackage     bool      `json:"is_package" gorm:"default:false"`
	PackageCount  int       `json:"package_count" gorm:"default:0"`
	DynamicPricing bool     `json:"dynamic_pricing" gorm:"default:false"`
	WeekendPrice  float64   `json:"weekend_price" gorm:"default:0"`
	HolidayPrice  float64   `json:"holiday_price" gorm:"default:0"`
	Status        int       `json:"status" gorm:"default:1"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
}

func (Service) TableName() string {
	return "services"
}

type PackageService struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ServiceID   uint      `json:"service_id" gorm:"index;not null"`
	Service     *Service  `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	ChildServiceID uint   `json:"child_service_id" gorm:"index;not null"`
	ChildService *Service `json:"child_service,omitempty" gorm:"foreignKey:ChildServiceID"`
	Count       int       `json:"count" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
}

func (PackageService) TableName() string {
	return "package_services"
}

type CustomerPackage struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CustomerID   uint      `json:"customer_id" gorm:"index;not null"`
	Customer     *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	ServiceID    uint      `json:"service_id" gorm:"index;not null"`
	Service      *Service  `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	TotalCount   int       `json:"total_count"`
	UsedCount    int       `json:"used_count" gorm:"default:0"`
	PurchaseDate time.Time `json:"purchase_date"`
	ExpireDate   time.Time `json:"expire_date"`
	Status       int       `json:"status" gorm:"default:1"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (CustomerPackage) TableName() string {
	return "customer_packages"
}

type Appointment struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CustomerID      uint      `json:"customer_id" gorm:"index;not null"`
	Customer        *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	TechnicianID    uint      `json:"technician_id" gorm:"index;not null"`
	Technician      *Technician `json:"technician,omitempty" gorm:"foreignKey:TechnicianID"`
	ServiceID       uint      `json:"service_id" gorm:"index;not null"`
	Service         *Service  `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	PackageID       *uint     `json:"package_id" gorm:"index"`
	AppointmentDate time.Time `json:"appointment_date" gorm:"index;not null"`
	StartTime       string    `json:"start_time" gorm:"size:10;not null"`
	EndTime         string    `json:"end_time" gorm:"size:10;not null"`
	Status          string    `json:"status" gorm:"size:20;default:pending"`
	Remark          string    `json:"remark" gorm:"size:255"`
	CancelReason    string    `json:"cancel_reason" gorm:"size:255"`
	PointsDeducted  int       `json:"points_deducted" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}

func (Appointment) TableName() string {
	return "appointments"
}

type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	AppointmentID uint      `json:"appointment_id" gorm:"index;not null"`
	Appointment   *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	CustomerID    uint      `json:"customer_id" gorm:"index;not null"`
	Amount        float64   `json:"amount" gorm:"not null"`
	PayMethod     string    `json:"pay_method" gorm:"size:20;not null"`
	PointsUsed    int       `json:"points_used" gorm:"default:0"`
	CardID        *uint     `json:"card_id" gorm:"index"`
	Status        string    `json:"status" gorm:"size:20;default:success"`
	TransactionNo string    `json:"transaction_no" gorm:"size:50;uniqueIndex"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}

type MemberCard struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CustomerID  uint      `json:"customer_id" gorm:"index;not null"`
	Customer    *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	CardNo      string    `json:"card_no" gorm:"size:50;uniqueIndex;not null"`
	CardType    string    `json:"card_type" gorm:"size:20"`
	Balance     float64   `json:"balance" gorm:"default:0"`
	Discount    float64   `json:"discount" gorm:"default:1.0"`
	Status      int       `json:"status" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

func (MemberCard) TableName() string {
	return "member_cards"
}

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Category    string    `json:"category" gorm:"size:50;index"`
	Description string    `json:"description" gorm:"size:500"`
	Unit        string    `json:"unit" gorm:"size:20"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Threshold   int       `json:"threshold" gorm:"default:10"`
	Price       float64   `json:"price" gorm:"default:0"`
	RetailPrice float64   `json:"retail_price" gorm:"default:0"`
	Supplier    string    `json:"supplier" gorm:"size:100"`
	Status      int       `json:"status" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

type ProductRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ProductID    uint      `json:"product_id" gorm:"index;not null"`
	Product      *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ChangeType   string    `json:"change_type" gorm:"size:20"`
	Quantity     int       `json:"quantity"`
	BeforeStock  int       `json:"before_stock"`
	AfterStock   int       `json:"after_stock"`
	AppointmentID *uint   `json:"appointment_id" gorm:"index"`
	OperatorID   uint      `json:"operator_id" gorm:"index"`
	Remark       string    `json:"remark" gorm:"size:255"`
	CreatedAt    time.Time `json:"created_at"`
}

func (ProductRecord) TableName() string {
	return "product_records"
}

type ProductSale struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CustomerID  uint      `json:"customer_id" gorm:"index;not null"`
	Customer    *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	ProductID   uint      `json:"product_id" gorm:"index;not null"`
	Product     *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	PayMethod   string    `json:"pay_method" gorm:"size:20"`
	OperatorID  uint      `json:"operator_id" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ProductSale) TableName() string {
	return "product_sales"
}

type Review struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	AppointmentID uint      `json:"appointment_id" gorm:"uniqueIndex;not null"`
	Appointment   *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	CustomerID    uint      `json:"customer_id" gorm:"index;not null"`
	Customer      *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	TechnicianID  uint      `json:"technician_id" gorm:"index;not null"`
	Technician    *Technician `json:"technician,omitempty" gorm:"foreignKey:TechnicianID"`
	ServiceID     uint      `json:"service_id" gorm:"index;not null"`
	Service       *Service  `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	Rating        int       `json:"rating"`
	Content       string    `json:"content" gorm:"size:500"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Review) TableName() string {
	return "reviews"
}

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type      string    `json:"type" gorm:"size:20"`
	Title     string    `json:"title" gorm:"size:100"`
	Content   string    `json:"content" gorm:"size:500"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action    string    `json:"action" gorm:"size:50"`
	Module    string    `json:"module" gorm:"size:50"`
	Detail    string    `json:"detail" gorm:"type:text"`
	IP        string    `json:"ip" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
