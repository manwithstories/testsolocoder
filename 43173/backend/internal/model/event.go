package model

import (
	"time"

	"gorm.io/gorm"
)

type EventStatus int

const (
	EventStatusDraft     EventStatus = 0
	EventStatusPublished EventStatus = 1
	EventStatusSoldOut   EventStatus = 2
	EventStatusEnded     EventStatus = 3
	EventStatusCancelled EventStatus = 4
)

type Event struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"index;not null"`
	ArtistID     uint           `json:"artist_id" gorm:"index;not null"`
	Title        string         `json:"title" gorm:"size:200;not null"`
	Description  string         `json:"description" gorm:"type:text"`
	CoverURL     string         `json:"cover_url" gorm:"size:255"`
	Venue        string         `json:"venue" gorm:"size:200;not null"`
	Address      string         `json:"address" gorm:"size:255"`
	City         string         `json:"city" gorm:"size:100"`
	Longitude    float64        `json:"longitude"`
	Latitude     float64        `json:"latitude"`
	StartTime    time.Time      `json:"start_time"`
	EndTime      time.Time      `json:"end_time"`
	DoorTime     *time.Time     `json:"door_time"`
	TicketPrice  float64        `json:"ticket_price"`
	TotalTickets int            `json:"total_tickets"`
	SoldTickets  int            `json:"sold_tickets" gorm:"default:0"`
	MaxPerUser   int            `json:"max_per_user" gorm:"default:4"`
	Status       EventStatus    `json:"status" gorm:"default:0"`
	PublishedAt  *time.Time     `json:"published_at"`
	ViewCount    int64          `json:"view_count" gorm:"default:0"`
	LikeCount    int64          `json:"like_count" gorm:"default:0"`
	ShareCount   int64          `json:"share_count" gorm:"default:0"`
	SeatMap      string         `json:"seat_map" gorm:"type:text"`
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Event) TableName() string {
	return "events"
}

type TicketStatus int

const (
	TicketStatusPending   TicketStatus = 0
	TicketStatusPaid      TicketStatus = 1
	TicketStatusUsed      TicketStatus = 2
	TicketStatusCancelled TicketStatus = 3
	TicketStatusRefunded  TicketStatus = 4
)

type Ticket struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderID     uint           `json:"order_id" gorm:"index;not null"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	EventID     uint           `json:"event_id" gorm:"index;not null"`
	SeatRow     int            `json:"seat_row"`
	SeatCol     int            `json:"seat_col"`
	SeatLabel   string         `json:"seat_label" gorm:"size:20"`
	QRCode      string         `json:"qr_code" gorm:"size:255"`
	Status      TicketStatus   `json:"status" gorm:"default:0"`
	UsedAt      *time.Time     `json:"used_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Ticket) TableName() string {
	return "tickets"
}

type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusPaid      OrderStatus = 1
	OrderStatusCompleted OrderStatus = 2
	OrderStatusCancelled OrderStatus = 3
	OrderStatusRefunded  OrderStatus = 4
)

type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderNo     string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	EventID     uint           `json:"event_id" gorm:"index;not null"`
	ArtistID    uint           `json:"artist_id" gorm:"index;not null"`
	TotalAmount float64        `json:"total_amount"`
	Quantity    int            `json:"quantity"`
	Status      OrderStatus    `json:"status" gorm:"default:0"`
	PaymentMethod string       `json:"payment_method" gorm:"size:50"`
	PaidAt      *time.Time     `json:"paid_at"`
	Remark      string         `json:"remark" gorm:"size:255"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Event       *Event         `json:"event,omitempty" gorm:"foreignKey:EventID"`
	Tickets     []Ticket       `json:"tickets,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Order) TableName() string {
	return "orders"
}
