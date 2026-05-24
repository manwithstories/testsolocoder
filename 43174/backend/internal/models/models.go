package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleStudent    UserRole = "student"
	RoleMerchant   UserRole = "merchant"
	RoleAdmin      UserRole = "admin"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusActive   UserStatus = "active"
	UserStatusRejected UserStatus = "rejected"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username        string      `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email           string      `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password        string      `gorm:"size:255;not null" json:"-"`
	Phone           string      `gorm:"size:20" json:"phone"`
	Role            UserRole    `gorm:"type:varchar(20);not null;default:student" json:"role"`
	Status          UserStatus  `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	RealName        string      `gorm:"size:50" json:"real_name"`
	SchoolName      string      `gorm:"size:100" json:"school_name"`
	StudentID       string      `gorm:"size:50" json:"student_id"`
	StudentCardURL  string      `gorm:"size:500" json:"student_card_url"`
	BusinessLicense string      `gorm:"size:500" json:"business_license"`
	Avatar          string      `gorm:"size:500" json:"avatar"`
	Rating          float64     `gorm:"default:5.0" json:"rating"`
	RatingCount     int         `gorm:"default:0" json:"rating_count"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Textbooks     []Textbook     `gorm:"foreignKey:SellerID" json:"textbooks,omitempty"`
	Notes         []Note         `gorm:"foreignKey:UploaderID" json:"notes,omitempty"`
	BuyerOrders   []Order        `gorm:"foreignKey:BuyerID" json:"buyer_orders,omitempty"`
	SellerOrders  []Order        `gorm:"foreignKey:SellerID" json:"seller_orders,omitempty"`
	SentMessages  []Message      `gorm:"foreignKey:SenderID" json:"sent_messages,omitempty"`
	ReceiveMessages []Message    `gorm:"foreignKey:ReceiverID" json:"receive_messages,omitempty"`
}

type TextbookCondition string

const (
	ConditionNew        TextbookCondition = "new"
	ConditionLikeNew    TextbookCondition = "like_new"
	ConditionGood       TextbookCondition = "good"
	ConditionFair       TextbookCondition = "fair"
)

type TextbookStatus string

const (
	TextbookStatusAvailable TextbookStatus = "available"
	TextbookStatusReserved  TextbookStatus = "reserved"
	TextbookStatusSold      TextbookStatus = "sold"
)

type Textbook struct {
	ID           uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ISBN         string            `gorm:"size:20;index" json:"isbn"`
	Title        string            `gorm:"size:200;not null" json:"title"`
	Author       string            `gorm:"size:100" json:"author"`
	CourseName   string            `gorm:"size:200;index" json:"course_name"`
	Edition      string            `gorm:"size:50" json:"edition"`
	Publisher    string            `gorm:"size:100" json:"publisher"`
	OriginalPrice float64          `json:"original_price"`
	Price        float64           `gorm:"not null" json:"price"`
	Condition    TextbookCondition `gorm:"type:varchar(20)" json:"condition"`
	Description  string            `gorm:"type:text" json:"description"`
	CoverImage   string            `gorm:"size:500" json:"cover_image"`
	Status       TextbookStatus    `gorm:"type:varchar(20);default:available" json:"status"`
	SellerID     uuid.UUID         `gorm:"type:uuid;index" json:"seller_id"`
	CategoryID   *uuid.UUID        `gorm:"type:uuid;index" json:"category_id"`
	ViewCount    int               `gorm:"default:0" json:"view_count"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `gorm:"index" json:"-"`

	Seller       *User             `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Category     *Category         `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	OrderItems   []OrderItem       `gorm:"foreignKey:TextbookID" json:"order_items,omitempty"`
	Reviews      []Review          `gorm:"foreignKey:TextbookID" json:"reviews,omitempty"`
}

type Category struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name       string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	ParentID   *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id"`
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

	Parent     *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children   []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Textbooks  []Textbook     `gorm:"foreignKey:CategoryID" json:"textbooks,omitempty"`
	Notes      []Note         `gorm:"foreignKey:CategoryID" json:"notes,omitempty"`
}

type Note struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Subject     string         `gorm:"size:100;index" json:"subject"`
	CourseName  string         `gorm:"size:200;index" json:"course_name"`
	Description string         `gorm:"type:text" json:"description"`
	FileURL     string         `gorm:"size:500;not null" json:"file_url"`
	FileType    string         `gorm:"size:20" json:"file_type"`
	FileSize    int64          `json:"file_size"`
	CoverImage  string         `gorm:"size:500" json:"cover_image"`
	UploaderID  uuid.UUID      `gorm:"type:uuid;index" json:"uploader_id"`
	CategoryID  *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	DownloadCount int          `gorm:"default:0" json:"download_count"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	Rating      float64        `gorm:"default:0" json:"rating"`
	RatingCount int            `gorm:"default:0" json:"rating_count"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Uploader    *User          `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Reviews     []Review       `gorm:"foreignKey:NoteID" json:"reviews,omitempty"`
}

type TransactionType string

const (
	TransactionTypeSell   TransactionType = "sell"
	TransactionTypeExchange TransactionType = "exchange"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusNegotiating TransactionStatus = "negotiating"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

type Transaction struct {
	ID             uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TextbookID     uuid.UUID         `gorm:"type:uuid;index" json:"textbook_id"`
	SellerID       uuid.UUID         `gorm:"type:uuid;index" json:"seller_id"`
	BuyerID        uuid.UUID         `gorm:"type:uuid;index" json:"buyer_id"`
	Type           TransactionType   `gorm:"type:varchar(20)" json:"type"`
	AgreedPrice    float64           `json:"agreed_price"`
	Status         TransactionStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	ExchangeItem   string            `gorm:"size:500" json:"exchange_item"`
	NegotiationHistory string        `gorm:"type:text" json:"negotiation_history"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`

	Textbook       *Textbook         `gorm:"foreignKey:TextbookID" json:"textbook,omitempty"`
	Seller         *User             `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Buyer          *User             `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Order          *Order            `gorm:"foreignKey:TransactionID" json:"order,omitempty"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
)

type Order struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderNo       string       `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	BuyerID       uuid.UUID    `gorm:"type:uuid;index" json:"buyer_id"`
	SellerID      uuid.UUID    `gorm:"type:uuid;index" json:"seller_id"`
	TotalAmount   float64      `gorm:"not null" json:"total_amount"`
	Status        OrderStatus  `gorm:"type:varchar(20);default:pending" json:"status"`
	PaymentMethod string       `gorm:"size:50" json:"payment_method"`
	PaymentStatus string       `gorm:"size:20" json:"payment_status"`
	TransactionID *uuid.UUID   `gorm:"type:uuid;index" json:"transaction_id"`
	ShippingAddress string     `gorm:"size:500" json:"shipping_address"`
	TrackingNumber  string     `gorm:"size:50" json:"tracking_number"`
	Remark        string       `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Buyer         *User        `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller        *User        `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Transaction   *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	Items         []OrderItem  `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	StatusHistory []OrderStatusHistory `gorm:"foreignKey:OrderID" json:"status_history,omitempty"`
}

type OrderItem struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID    uuid.UUID   `gorm:"type:uuid;index" json:"order_id"`
	TextbookID uuid.UUID   `gorm:"type:uuid;index" json:"textbook_id"`
	Quantity   int         `gorm:"default:1" json:"quantity"`
	Price      float64     `json:"price"`
	Subtotal   float64     `json:"subtotal"`
	CreatedAt  time.Time   `json:"created_at"`

	Order      *Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Textbook   *Textbook   `gorm:"foreignKey:TextbookID" json:"textbook,omitempty"`
}

type OrderStatusHistory struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID     uuid.UUID   `gorm:"type:uuid;index" json:"order_id"`
	Status      OrderStatus `gorm:"type:varchar(20)" json:"status"`
	Remark      string      `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time   `json:"created_at"`

	Order       *Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

type Message struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SenderID   uuid.UUID      `gorm:"type:uuid;index" json:"sender_id"`
	ReceiverID uuid.UUID      `gorm:"type:uuid;index" json:"receiver_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	RelatedOrderID *uuid.UUID `gorm:"type:uuid;index" json:"related_order_id"`
	IsDispute  bool           `gorm:"default:false" json:"is_dispute"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

	Sender     *User          `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver   *User          `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
}

type ReviewTargetType string

const (
	ReviewTargetTextbook ReviewTargetType = "textbook"
	ReviewTargetNote     ReviewTargetType = "note"
)

type Review struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID     uuid.UUID        `gorm:"type:uuid;index" json:"user_id"`
	TargetType ReviewTargetType `gorm:"type:varchar(20)" json:"target_type"`
	TextbookID *uuid.UUID       `gorm:"type:uuid;index" json:"textbook_id"`
	NoteID     *uuid.UUID       `gorm:"type:uuid;index" json:"note_id"`
	Rating     int              `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Content    string           `gorm:"type:text" json:"content"`
	IsHidden   bool             `gorm:"default:false" json:"is_hidden"`
	IsMalicious bool            `gorm:"default:false" json:"is_malicious"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `gorm:"index" json:"-"`

	User       *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Textbook   *Textbook        `gorm:"foreignKey:TextbookID" json:"textbook,omitempty"`
	Note       *Note            `gorm:"foreignKey:NoteID" json:"note,omitempty"`
}

type Notification struct {
	ID      uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID  uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	Type    string         `gorm:"size:50" json:"type"`
	Title   string         `gorm:"size:200" json:"title"`
	Content string         `gorm:"type:text" json:"content"`
	IsRead  bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time    `json:"created_at"`

	User    *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
