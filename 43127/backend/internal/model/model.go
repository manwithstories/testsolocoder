package model

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50"`
	Password     string    `json:"-"`
	RealName     string    `json:"realName"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Avatar       string    `json:"avatar"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Property struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Title         string    `json:"title"`
	Community     string    `json:"community"`
	Address       string    `json:"address"`
	Area          float64   `json:"area"`
	Layout        string    `json:"layout"`
	Floor         string    `json:"floor"`
	Rent          float64   `json:"rent"`
	Deposit       float64   `json:"deposit"`
	PaymentType   string    `json:"paymentType"`
	Description   string    `json:"description"`
	Status        int       `json:"status"`
	OwnerID       uint      `json:"ownerId"`
	Owner         *User    `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Images        []PropertyImage `json:"images,omitempty" gorm:"foreignKey:PropertyID"`
	Facilities    []Facility `json:"facilities,omitempty" gorm:"many2many:property_facilities;"`
	Region        string    `json:"region"`
	Building      string    `json:"building"`
	RoomNo        string    `json:"roomNo"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type PropertyImage struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	PropertyID uint   `json:"propertyId"`
	URL        string `json:"url"`
	Sort       int    `json:"sort"`
}

type Facility struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Tenant struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	IDCard      string    `json:"idCard"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Appointment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenantId"`
	Tenant      *Tenant   `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	PropertyID  uint      `json:"propertyId"`
	Property    *Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
	VisitTime   time.Time `json:"visitTime"`
	Status      int       `json:"status"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Contract struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenantId"`
	Tenant      *Tenant   `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	PropertyID  uint      `json:"propertyId"`
	Property    *Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Rent        float64   `json:"rent"`
	Deposit     float64   `json:"deposit"`
	PaymentType string    `json:"paymentType"`
	Status      int       `json:"status"`
	FileURL     string    `json:"fileUrl"`
	IsReminded  bool      `json:"isReminded"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RentRecord struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ContractID  uint      `json:"contractId"`
	Contract    *Contract `json:"contract,omitempty" gorm:"foreignKey:ContractID"`
	TenantID    uint      `json:"tenantId"`
	PropertyID  uint      `json:"propertyId"`
	Month       string    `json:"month"`
	Amount      float64   `json:"amount"`
	Status      int       `json:"status"`
	DueDate     time.Time `json:"dueDate"`
	PaidAt      *time.Time `json:"paidAt"`
	LateFee     float64   `json:"lateFee"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RepairOrder struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenantId"`
	Tenant      *Tenant   `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	PropertyID  uint      `json:"propertyId"`
	Property    *Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Images      string    `json:"images"`
	Priority    int       `json:"priority"`
	Status      int       `json:"status"`
	HandlerID    *uint     `json:"handlerId"`
	Handler     *User     `json:"handler,omitempty" gorm:"foreignKey:HandlerID"`
	ProcessNote   string    `json:"processNote"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UtilityFee struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PropertyID  uint      `json:"propertyId"`
	Property    *Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
	Type        string    `json:"type"`
	Month       string    `json:"month"`
	TotalAmount float64   `json:"totalAmount"`
	Units       float64   `json:"units"`
	UnitPrice   float64   `json:"unitPrice"`
	Status      int       `json:"status"`
	DueDate     time.Time `json:"dueDate"`
	PaidAt      *time.Time `json:"paidAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Notice struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Type        int       `json:"type"`
	Building    string    `json:"building"`
	PublisherID uint      `json:"publisherId"`
	Publisher   *User     `json:"publisher,omitempty" gorm:"foreignKey:PublisherID"`
	IsTop       int       `json:"isTop"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Statistics struct {
	OccupancyRate float64 `json:"occupancyRate"`
	TotalIncome  float64 `json:"totalIncome"`
	TotalOrders  int     `json:"totalOrders"`
}
