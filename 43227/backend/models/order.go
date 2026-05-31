package models

import (
	"time"
)

type Order struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	OrderNo        string     `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	BuyerID        uint       `gorm:"not null;index" json:"buyer_id"`
	SellerID       uint       `gorm:"not null;index" json:"seller_id"`
	ProductID      uint       `gorm:"not null;index" json:"product_id"`
	Quantity       float64    `gorm:"type:decimal(10,2);not null" json:"quantity"`
	UnitPrice      float64    `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalAmount    float64    `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status         string     `gorm:"size:20;default:pending;index" json:"status"`
	PaymentStatus  string     `gorm:"size:20;default:unpaid" json:"payment_status"`
	PaymentTime    *time.Time `json:"payment_time"`
	ShippingAddress string    `gorm:"type:text;not null" json:"shipping_address"`
	TrackingNumber string     `gorm:"size:100" json:"tracking_number"`
	TrackingStatus string     `gorm:"size:50" json:"tracking_status"`
	BuyerRating    *float64   `gorm:"type:decimal(2,1)" json:"buyer_rating"`
	SellerRating   *float64   `gorm:"type:decimal(2,1)" json:"seller_rating"`
	BuyerComment   string     `gorm:"type:text" json:"buyer_comment"`
	SellerComment  string     `gorm:"type:text" json:"seller_comment"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Buyer          User       `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller         User       `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Product        Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type CreateOrderRequest struct {
	ProductID       uint    `json:"product_id" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`
	ShippingAddress string  `json:"shipping_address" binding:"required"`
}

type UpdateOrderRequest struct {
	Status         *string `json:"status"`
	TrackingNumber *string `json:"tracking_number"`
	TrackingStatus *string `json:"tracking_status"`
}

type RateOrderRequest struct {
	Rating  float64 `json:"rating" binding:"required,gte=1,lte=5"`
	Comment string  `json:"comment"`
}

type Post struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	Category      string    `gorm:"size:50;not null;index" json:"category"`
	Tags          []string  `gorm:"serializer:json" json:"tags"`
	Images        []string  `gorm:"serializer:json" json:"images"`
	ViewCount     int       `gorm:"default:0" json:"view_count"`
	LikeCount     int       `gorm:"default:0" json:"like_count"`
	CommentCount  int       `gorm:"default:0" json:"comment_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}

type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,max=200"`
	Content  string   `json:"content" binding:"required"`
	Category string   `json:"category" binding:"required,oneof=disease_control harvest_technique seasonal_management equipment market general"`
	Tags     []string `json:"tags"`
	Images   []string `json:"images"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	LikeCount int       `gorm:"default:0" json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

type CreateCommentRequest struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	RelatedID *uint     `json:"related_id"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

type OperationLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       *uint     `gorm:"index" json:"user_id"`
	Operation    string    `gorm:"size:50;not null" json:"operation"`
	Module       string    `gorm:"size:50" json:"module"`
	Description  string    `gorm:"type:text" json:"description"`
	IPAddress    string    `gorm:"size:45" json:"ip_address"`
	UserAgent    string    `gorm:"size:255" json:"user_agent"`
	Status       string    `gorm:"size:20;default:success" json:"status"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

type Upload struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	FileName  string    `gorm:"size:255;not null" json:"file_name"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileSize  int64     `json:"file_size"`
	FileType  string    `gorm:"size:50" json:"file_type"`
	Category  string    `gorm:"size:50" json:"category"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Upload) TableName() string {
	return "uploads"
}
