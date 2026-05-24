package model

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleBuyer    Role = "buyer"
	RoleSeller   Role = "seller"
	RoleAppraiser Role = "appraiser"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password      string         `gorm:"not null" json:"-"`
	Role          Role           `gorm:"size:16;not null" json:"role"`
	Email         string         `gorm:"size:128" json:"email"`
	Phone         string         `gorm:"size:32" json:"phone"`
	Avatar        string         `gorm:"size:255" json:"avatar"`
	CreditScore   int            `gorm:"default:100" json:"credit_score"`
	ReviewCount   int            `gorm:"default:0" json:"review_count"`
	RealName      string         `gorm:"size:64" json:"real_name"`
	Certified     bool           `gorm:"default:false" json:"certified"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Watch struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SellerID     uint           `gorm:"index;not null" json:"seller_id"`
	Brand        string         `gorm:"size:64;index;not null" json:"brand"`
	Model        string         `gorm:"size:128;not null" json:"model"`
	ReferenceNo  string         `gorm:"size:64" json:"reference_no"`
	Year         int            `json:"year"`
	Movement     string         `gorm:"size:64" json:"movement"`
	CaseSizeMM   float64        `json:"case_size_mm"`
	CaseMaterial string         `gorm:"size:64" json:"case_material"`
	DialColor    string         `gorm:"size:64" json:"dial_color"`
	Bracelet     string         `gorm:"size:64" json:"bracelet"`
	Condition    string         `gorm:"size:32" json:"condition"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        float64        `json:"price"`
	Status       string         `gorm:"size:32;index;default:on_sale" json:"status"`
	Authed       bool           `gorm:"default:false" json:"authed"`
	Photos       []WatchPhoto   `gorm:"foreignKey:WatchID" json:"photos,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type WatchPhoto struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WatchID   uint      `gorm:"index;not null" json:"watch_id"`
	URL       string    `gorm:"size:512;not null" json:"url"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthOrderStatus string

const (
	AuthPending  AuthOrderStatus = "pending"
	AuthAssigned AuthOrderStatus = "assigned"
	AuthReported AuthOrderStatus = "reported"
	AuthRejected AuthOrderStatus = "rejected"
)

type AuthOrder struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	UserID      uint            `gorm:"index;not null" json:"user_id"`
	WatchID     uint            `gorm:"index" json:"watch_id"`
	AppraiserID uint            `gorm:"index" json:"appraiser_id"`
	Status      AuthOrderStatus `gorm:"size:32;index;not null;default:pending" json:"status"`
	Note        string          `gorm:"type:text" json:"note"`
	Photos      []AuthPhoto     `gorm:"foreignKey:AuthOrderID" json:"photos,omitempty"`
	Report      *AuthReport     `gorm:"foreignKey:AuthOrderID" json:"report,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type AuthPhoto struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AuthOrderID uint      `gorm:"index;not null" json:"auth_order_id"`
	URL         string    `gorm:"size:512;not null" json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

type AuthReport struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AuthOrderID   uint      `gorm:"uniqueIndex;not null" json:"auth_order_id"`
	AppraiserID   uint      `gorm:"index;not null" json:"appraiser_id"`
	Conclusion    string    `gorm:"size:32;not null" json:"conclusion"`
	Authentic     bool      `json:"authentic"`
	Details       string    `gorm:"type:text" json:"details"`
	EstimatedValue float64 `json:"estimated_value"`
	PDFPath       string    `gorm:"size:512" json:"pdf_path"`
	CreatedAt     time.Time `json:"created_at"`
}

type TradeStatus string

const (
	TradeOpen        TradeStatus = "open"
	TradeBidding     TradeStatus = "bidding"
	TradePendingDeal TradeStatus = "pending_deal"
	TradeDealed      TradeStatus = "dealed"
	TradeShipped     TradeStatus = "shipped"
	TradeCompleted   TradeStatus = "completed"
	TradeCanceled    TradeStatus = "canceled"
)

type Trade struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	SellerID   uint        `gorm:"index;not null" json:"seller_id"`
	BuyerID    uint        `gorm:"index" json:"buyer_id"`
	WatchID    uint        `gorm:"index;not null" json:"watch_id"`
	StartPrice float64     `json:"start_price"`
	FinalPrice float64     `json:"final_price"`
	Status     TradeStatus `gorm:"size:32;index;not null;default:open" json:"status"`
	Remark     string      `gorm:"type:text" json:"remark"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type TradeBid struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TradeID   uint      `gorm:"index;not null" json:"trade_id"`
	BuyerID   uint      `gorm:"index;not null" json:"buyer_id"`
	Price     float64   `json:"price"`
	Message   string    `gorm:"size:512" json:"message"`
	Accepted  bool      `gorm:"default:false" json:"accepted"`
	CreatedAt time.Time `json:"created_at"`
}

type FavoriteGroup struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Brand     string         `gorm:"size:64;index;not null" json:"brand"`
	Name      string         `gorm:"size:64" json:"name"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
}

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	WatchID   uint      `gorm:"index;not null" json:"watch_id"`
	GroupID   uint      `gorm:"index" json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TradeID    uint      `gorm:"uniqueIndex:trade_role;not null" json:"trade_id"`
	FromUserID uint      `gorm:"index;not null" json:"from_user_id"`
	ToUserID   uint      `gorm:"index;not null" json:"to_user_id"`
	Role       string    `gorm:"size:16;uniqueIndex:trade_role" json:"role"`
	Rating     int       `gorm:"not null" json:"rating"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Type      string    `gorm:"size:32;index" json:"type"`
	Title     string    `gorm:"size:256" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	RefType   string    `gorm:"size:32" json:"ref_type"`
	RefID     uint      `json:"ref_id"`
	Read      bool      `gorm:"default:false" json:"read"`
	CreatedAt time.Time `json:"created_at"`
}
