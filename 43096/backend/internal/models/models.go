package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	BaseModel
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password     string    `gorm:"type:varchar(255);not null" json:"-"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone        string    `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	RealName     string    `gorm:"type:varchar(50)" json:"real_name"`
	IDCard       string    `gorm:"type:varchar(18);uniqueIndex" json:"id_card"`
	MemberLevel  int       `gorm:"default:1;not null" json:"member_level"`
	Points       int       `gorm:"default:0;not null" json:"points"`
	Status       int       `gorm:"default:1;not null" json:"status"`
	Avatar       string    `gorm:"type:varchar(255)" json:"avatar"`
	Role         string    `gorm:"type:varchar(20);default:'user'" json:"role"`
}

type MemberLevel struct {
	BaseModel
	Name        string  `gorm:"type:varchar(50);not null" json:"name"`
	Level       int     `gorm:"uniqueIndex;not null" json:"level"`
	Discount    float64 `gorm:"type:decimal(3,2);not null" json:"discount"`
	MinPoints   int     `gorm:"default:0;not null" json:"min_points"`
	Priority    int     `gorm:"default:0;not null" json:"priority"`
	Description string  `gorm:"type:text" json:"description"`
}

type Coupon struct {
	BaseModel
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Type        int       `gorm:"not null" json:"type"`
	Value       float64 `gorm:"type:decimal(10,2);not null" json:"value"`
	MinAmount   float64 `gorm:"type:decimal(10,2);default:0" json:"min_amount"`
	UserID      uint64    `gorm:"index" json:"user_id"`
	Status      int       `gorm:"default:1;not null" json:"status"`
	ExpireAt    time.Time `json:"expire_at"`
}

type Show struct {
	BaseModel
	Name         string    `gorm:"type:varchar(200);not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	Poster       string    `gorm:"type:varchar(500)" json:"poster"`
	Artist       string    `gorm:"type:varchar(200)" json:"artist"`
	Duration     int       `gorm:"default:0" json:"duration"`
	Status         int       `gorm:"default:0;not null" json:"status"`
	OrganizerID  string    `gorm:"type:varchar(200)" json:"organizer"`
	Venue        string    `gorm:"type:varchar(200);not null" json:"venue"`
	Address      string    `gorm:"type:varchar(500)" json:"address"`
	Sessions     []Session `gorm:"foreignKey:ShowID" json:"sessions,omitempty"`
}

type Session struct {
	BaseModel
	ShowID           uint64    `gorm:"index;not null" json:"show_id"`
	StartTime        time.Time `gorm:"not null" json:"start_time"`
	EndTime          time.Time `gorm:"not null" json:"end_time"`
	PresaleStartTime time.Time `gorm:"" json:"presale_start_time"`
	Status           int       `gorm:"default:1;not null" json:"status"`
	TotalSeats       int       `gorm:"default:0;not null" json:"total_seats"`
	SoldSeats        int       `gorm:"default:0;not null" json:"sold_seats"`
	SeatAreas        []SeatArea `gorm:"foreignKey:SessionID" json:"seat_areas,omitempty"`
}

type SeatArea struct {
	BaseModel
	SessionID   uint64  `gorm:"index;not null" json:"session_id"`
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Color       string  `gorm:"type:varchar(20)" json:"color"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	TotalSeats  int     `gorm:"default:0;not null" json:"total_seats"`
	SoldSeats   int     `gorm:"default:0;not null" json:"sold_seats"`
	SortOrder   int     `gorm:"default:0;not null" json:"sort_order"`
}

type Seat struct {
	BaseModel
	SessionID   uint64 `gorm:"index;not null" json:"session_id"`
	AreaID      uint64 `gorm:"index;not null" json:"area_id"`
	Row         string `gorm:"type:varchar(20);not null" json:"row"`
	Col         int    `gorm:"not null" json:"col"`
	SeatNo      string `gorm:"type:varchar(50);not null" json:"seat_no"`
	Status      int    `gorm:"default:0;not null" json:"status"`
	X           int    `gorm:"default:0" json:"x"`
	Y           int    `gorm:"default:0" json:"y"`
	Width       int    `gorm:"default:30" json:"width"`
	Height      int    `gorm:"default:30" json:"height"`
}

type SeatChart struct {
	BaseModel
	SessionID   uint64 `gorm:"uniqueIndex;not null" json:"session_id"`
	Config      string `gorm:"type:text;not null" json:"config"`
	Background  string `gorm:"type:varchar(500)" json:"background"`
	Width       int    `gorm:"default:800" json:"width"`
	Height      int    `gorm:"default:600" json:"height"`
}

type Order struct {
	BaseModel
	OrderNo     string    `gorm:"type:varchar(32);uniqueIndex;not null" json:"order_no"`
	UserID      uint64    `gorm:"index;not null" json:"user_id"`
	ShowID      uint64    `gorm:"index;not null" json:"show_id"`
	SessionID   uint64    `gorm:"index;not null" json:"session_id"`
	TotalAmount float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Discount    float64   `gorm:"type:decimal(10,2);default:0" json:"discount"`
	PayAmount   float64   `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	Status      int       `gorm:"default:0;not null" json:"status"`
	PayType     int       `json:"pay_type"`
	PayTime     *time.Time `json:"pay_time"`
	RealName    string    `gorm:"type:varchar(50)" json:"real_name"`
	IDCard      string    `gorm:"type:varchar(18)" json:"id_card"`
	Phone       string    `gorm:"type:varchar(20)" json:"phone"`
	Email       string    `gorm:"type:varchar(100)" json:"email"`
	Remark      string    `gorm:"type:varchar(500)" json:"remark"`
	Tickets     []Ticket  `gorm:"foreignKey:OrderID" json:"tickets,omitempty"`
	Refund      *Refund   `gorm:"foreignKey:OrderID" json:"refund,omitempty"`
}

type Ticket struct {
	BaseModel
	OrderID     uint64    `gorm:"index;not null" json:"order_id"`
	UserID      uint64    `gorm:"index;not null" json:"user_id"`
	ShowID      uint64    `gorm:"index;not null" json:"show_id"`
	SessionID   uint64    `gorm:"index;not null" json:"session_id"`
	SeatID      uint64    `gorm:"index;not null" json:"seat_id"`
	AreaID      uint64    `gorm:"index;not null" json:"area_id"`
	TicketNo    string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"ticket_no"`
	QrCode      string    `gorm:"type:varchar(500);not null" json:"qr_code"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	SeatInfo    string    `gorm:"type:varchar(100);not null" json:"seat_info"`
	Status      int       `gorm:"default:0;not null" json:"status"`
	CheckedIn   int       `gorm:"default:0;not null" json:"checked_in"`
	CheckinTime *time.Time `json:"checkin_time"`
	RealName    string    `gorm:"type:varchar(50)" json:"real_name"`
	IDCard      string    `gorm:"type:varchar(18)" json:"id_card"`
}

type Refund struct {
	BaseModel
	OrderID       uint64    `gorm:"index;not null" json:"order_id"`
	UserID        uint64    `gorm:"index;not null" json:"user_id"`
	RefundNo      string    `gorm:"type:varchar(32);uniqueIndex;not null" json:"refund_no"`
	RefundAmount  float64   `gorm:"type:decimal(10,2);not null" json:"refund_amount"`
	Reason        string    `gorm:"type:varchar(500)" json:"reason"`
	Status        int       `gorm:"default:0;not null" json:"status"`
	AuditTime     *time.Time `json:"audit_time"`
	AuditRemark   string    `gorm:"type:varchar(500)" json:"audit_remark"`
}

type CheckinLog struct {
	BaseModel
	TicketID    uint64    `gorm:"index;not null" json:"ticket_id"`
	TicketNo    string    `gorm:"type:varchar(50);index;not null" json:"ticket_no"`
	OperatorID  uint64    `gorm:"index;not null" json:"operator_id"`
	Status      int       `gorm:"not null" json:"status"`
	Message     string    `gorm:"type:varchar(500)" json:"message"`
}

type PaymentLog struct {
	BaseModel
	OrderID     uint64    `gorm:"index;not null" json:"order_id"`
	PayType     int       `gorm:"not null" json:"pay_type"`
	TradeNo     string    `gorm:"type:varchar(100)" json:"trade_no"`
	Amount      float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status      int       `gorm:"not null" json:"status"`
	Request     string    `gorm:"type:text" json:"request"`
	Response    string    `gorm:"type:text" json:"response"`
	RetryCount  int       `gorm:"default:0" json:"retry_count"`
}

type OperationLog struct {
	BaseModel
	UserID      uint64    `gorm:"index" json:"user_id"`
	Module      string    `gorm:"type:varchar(50);not null" json:"module"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"`
	Content     string    `gorm:"type:text" json:"content"`
	IP          string    `gorm:"type:varchar(50)" json:"ip"`
}

const (
	SeatStatusAvailable = 0
	SeatStatusLocked    = 1
	SeatStatusSold      = 2
	SeatStatusDisabled  = 3

	OrderStatusPending   = 0
	OrderStatusPaid    = 1
	OrderStatusCanceled = 2
	OrderStatusRefunded = 3
	OrderStatusRefunding = 4

	TicketStatusValid   = 0
	TicketStatusUsed  = 1
	TicketStatusRefunded = 2

	RefundStatusPending = 0
	RefundStatusApproved = 1
	RefundStatusRejected = 2

	MemberLevelNormal = 1
	MemberLevelSilver = 2
	MemberLevelGold = 3
	MemberLevelPlatinum = 4
)
