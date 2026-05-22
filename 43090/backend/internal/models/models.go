package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Balance   float64        `gorm:"type:decimal(10,2);default:0" json:"balance"`
	Role      string         `gorm:"size:20;default:'user'" json:"role"`
	Status    int            `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	ParentID  uint      `gorm:"default:0" json:"parent_id"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuctionItem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	CategoryID    uint           `gorm:"not null;index" json:"category_id"`
	SellerID      uint           `gorm:"not null;index" json:"seller_id"`
	StartPrice    float64        `gorm:"type:decimal(10,2);not null" json:"start_price"`
	ReservePrice  float64        `gorm:"type:decimal(10,2);not null" json:"reserve_price"`
	CurrentPrice  float64        `gorm:"type:decimal(10,2);default:0" json:"current_price"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	BidCount      int            `gorm:"default:0" json:"bid_count"`
	Status        int            `gorm:"default:0" json:"status"`
	Location      string         `gorm:"size:200" json:"location"`
	Condition     string         `gorm:"size:50" json:"condition"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Seller        User           `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Images        []AuctionImage `gorm:"foreignKey:AuctionItemID" json:"images,omitempty"`
}

type AuctionImage struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AuctionItemID uint      `gorm:"not null;index" json:"auction_item_id"`
	URL           string    `gorm:"size:255;not null" json:"url"`
	Sort          int       `gorm:"default:0" json:"sort"`
	IsMain        int       `gorm:"default:0" json:"is_main"`
	CreatedAt     time.Time `json:"created_at"`
}

type AuctionSession struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"size:200;not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	StartTime       time.Time `gorm:"not null" json:"start_time"`
	EndTime         time.Time `gorm:"not null" json:"end_time"`
	MinIncrement    float64   `gorm:"type:decimal(10,2);not null" json:"min_increment"`
	ExtendTime      int       `gorm:"default:300" json:"extend_time"`
	Status          int       `gorm:"default:0" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	AuctionItems    []AuctionItemSession `gorm:"foreignKey:SessionID" json:"auction_items,omitempty"`
}

type AuctionItemSession struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	SessionID     uint         `gorm:"not null;index" json:"session_id"`
	AuctionItemID uint         `gorm:"not null;index" json:"auction_item_id"`
	Sort          int          `gorm:"default:0" json:"sort"`
	StartTime     time.Time    `json:"start_time"`
	EndTime       time.Time    `json:"end_time"`
	Status        int          `gorm:"default:0" json:"status"`
	AuctionItem   AuctionItem  `gorm:"foreignKey:AuctionItemID" json:"auction_item,omitempty"`
}

type Bid struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AuctionItemID   uint      `gorm:"not null;index" json:"auction_item_id"`
	SessionID       uint      `gorm:"index" json:"session_id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	MaxAutoBid      float64   `gorm:"type:decimal(10,2);default:0" json:"max_auto_bid"`
	IsAutoBid       int       `gorm:"default:0" json:"is_auto_bid"`
	IsWinning       int       `gorm:"default:0" json:"is_winning"`
	CreatedAt       time.Time `json:"created_at"`
	AuctionItem     AuctionItem `gorm:"foreignKey:AuctionItemID" json:"auction_item,omitempty"`
	User            User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type AutoBid struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AuctionItemID   uint      `gorm:"not null;index" json:"auction_item_id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	MaxPrice        float64   `gorm:"type:decimal(10,2);not null" json:"max_price"`
	CurrentBid      float64   `gorm:"type:decimal(10,2);default:0" json:"current_bid"`
	Status          int       `gorm:"default:1" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	AuctionItem     AuctionItem `gorm:"foreignKey:AuctionItemID" json:"auction_item,omitempty"`
	User            User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Order struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderNo       string    `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	AuctionItemID uint      `gorm:"not null;index" json:"auction_item_id"`
	BuyerID       uint      `gorm:"not null;index" json:"buyer_id"`
	SellerID      uint      `gorm:"not null;index" json:"seller_id"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Status        int       `gorm:"default:0" json:"status"`
	PaymentTime   *time.Time `json:"payment_time"`
	ShippingInfo  string    `gorm:"type:text" json:"shipping_info"`
	TrackingNo    string    `gorm:"size:100" json:"tracking_no"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AuctionItem   AuctionItem `gorm:"foreignKey:AuctionItemID" json:"auction_item,omitempty"`
	Buyer         User      `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller        User      `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Payments      []Payment `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"not null;index" json:"order_id"`
	PaymentNo     string    `gorm:"size:50;uniqueIndex;not null" json:"payment_no"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Method        string    `gorm:"size:50" json:"method"`
	Status        int       `gorm:"default:0" json:"status"`
	TransactionID string    `gorm:"size:255" json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	PaidAt        *time.Time `json:"paid_at"`
}

type Review struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"not null;index" json:"order_id"`
	ReviewerID    uint      `gorm:"not null;index" json:"reviewer_id"`
	RevieweeID    uint      `gorm:"not null;index" json:"reviewee_id"`
	Rating        int       `gorm:"not null" json:"rating"`
	Content       string    `gorm:"type:text" json:"content"`
	CreatedAt     time.Time `json:"created_at"`
	Reviewer      User      `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Reviewee      User      `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
}

type Notification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Type       string    `gorm:"size:50;not null" json:"type"`
	Title      string    `gorm:"size:200;not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	RelatedID  uint      `json:"related_id"`
	IsRead     int       `gorm:"default:0" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type SystemLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Action    string    `gorm:"size:100;not null" json:"action"`
	Module    string    `gorm:"size:50" json:"module"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Params    string    `gorm:"type:text" json:"params"`
	Result    string    `gorm:"type:text" json:"result"`
	CreatedAt time.Time `json:"created_at"`
}
