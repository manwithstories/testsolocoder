package model

import (
	"fmt"
	"time"

	"luxury-trading-platform/internal/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRole string

const (
	RoleBuyer       UserRole = "buyer"
	RoleSeller      UserRole = "seller"
	RoleAuthenticator UserRole = "authenticator"
	RoleAdmin       UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
	UserStatusBanned   UserStatus = "banned"
)

type ProductCategory string

const (
	CategoryBag     ProductCategory = "bag"
	CategoryJewelry ProductCategory = "jewelry"
	CategoryWatch   ProductCategory = "watch"
	CategoryClothing ProductCategory = "clothing"
	CategoryShoes   ProductCategory = "shoes"
	CategoryOther   ProductCategory = "other"
)

type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusOnSale    ProductStatus = "on_sale"
	ProductStatusSold      ProductStatus = "sold"
	ProductStatusRemoved   ProductStatus = "removed"
)

type AuthenticatorStatus string

const (
	AuthenticatorStatusPending  AuthenticatorStatus = "pending"
	AuthenticatorStatusApproved AuthenticatorStatus = "approved"
	AuthenticatorStatusRejected AuthenticatorStatus = "rejected"
)

type AuthenticationStatus string

const (
	AuthenticationStatusPending   AuthenticationStatus = "pending"
	AuthenticationStatusAccepted  AuthenticationStatus = "accepted"
	AuthenticationStatusCompleted AuthenticationStatus = "completed"
	AuthenticationStatusRejected  AuthenticationStatus = "rejected"
	AuthenticationStatusCancelled AuthenticationStatus = "cancelled"
)

type AuthenticationResult string

const (
	AuthenticationResultGenuine    AuthenticationResult = "genuine"
	AuthenticationResultCounterfeit AuthenticationResult = "counterfeit"
	AuthenticationResultInconclusive AuthenticationResult = "inconclusive"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type ReviewRating int

const (
	RatingOne   ReviewRating = 1
	RatingTwo   ReviewRating = 2
	RatingThree ReviewRating = 3
	RatingFour  ReviewRating = 4
	RatingFive  ReviewRating = 5
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	BaseModel
	Username     string          `gorm:"uniqueIndex;size:100;not null" json:"username" validate:"required,min=3,max=100"`
	Email        string          `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
	Phone        string          `gorm:"size:20" json:"phone"`
	Password     string          `gorm:"size:255;not null" json:"-"`
	RealName     string          `gorm:"size:100" json:"real_name"`
	Avatar       string          `gorm:"size:255" json:"avatar"`
	Role         UserRole        `gorm:"size:20;not null;index" json:"role"`
	Status       UserStatus      `gorm:"size:20;not null;default:active" json:"status"`
	CreditScore  int             `gorm:"default:100" json:"credit_score"`
	Address      string          `gorm:"size:500" json:"address"`

	AuthenticatorProfile *AuthenticatorProfile `json:"authenticator_profile,omitempty"`
	SellerProfile        *SellerProfile        `json:"seller_profile,omitempty"`
	BuyerProfile         *BuyerProfile         `json:"buyer_profile,omitempty"`

	Products        []Product         `gorm:"foreignKey:SellerID" json:"products,omitempty"`
	Orders          []Order           `gorm:"foreignKey:BuyerID" json:"orders,omitempty"`
	ReviewsGiven    []Review          `gorm:"foreignKey:ReviewerID" json:"reviews_given,omitempty"`
	ReviewsReceived []Review          `gorm:"foreignKey:RevieweeID" json:"reviews_received,omitempty"`
}

type AuthenticatorProfile struct {
	BaseModel
	UserID           uint                `gorm:"uniqueIndex;not null" json:"user_id"`
	LicenseNumber    string              `gorm:"size:50;uniqueIndex" json:"license_number"`
	LicenseFile      string              `gorm:"size:255" json:"license_file"`
	Certifications   string              `gorm:"type:text" json:"certifications"`
	Specialties      string              `gorm:"size:500" json:"specialties"`
	Status           AuthenticatorStatus `gorm:"size:20;not null;default:pending" json:"status"`
	Rating           float64             `gorm:"default:0" json:"rating"`
	CompletedCount   int                 `gorm:"default:0" json:"completed_count"`
	RejectionReason  string              `gorm:"size:500" json:"rejection_reason"`
	VerifiedAt       *time.Time          `json:"verified_at"`

	User             *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Authentications  []Authentication   `gorm:"foreignKey:AuthenticatorID" json:"authentications,omitempty"`
}

type SellerProfile struct {
	BaseModel
	UserID          uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	StoreName       string `gorm:"size:200" json:"store_name"`
	StoreLogo       string `gorm:"size:255" json:"store_logo"`
	Description     string `gorm:"type:text" json:"description"`
	BusinessLicense string `gorm:"size:255" json:"business_license"`
	TotalSales      int    `gorm:"default:0" json:"total_sales"`

	User   *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type BuyerProfile struct {
	BaseModel
	UserID           uint    `gorm:"uniqueIndex;not null" json:"user_id"`
	PreferredBrands  string  `gorm:"size:500" json:"preferred_brands"`
	TotalPurchases   int     `gorm:"default:0" json:"total_purchases"`
	TotalSpent       float64 `gorm:"default:0" json:"total_spent"`

	User   *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Brand struct {
	BaseModel
	Name        string          `gorm:"uniqueIndex;size:100;not null" json:"name"`
	NameCN      string          `gorm:"size:100" json:"name_cn"`
	Logo        string          `gorm:"size:255" json:"logo"`
	Country     string          `gorm:"size:50" json:"country"`
	Description string          `gorm:"type:text" json:"description"`
	Category    ProductCategory `gorm:"size:20" json:"category"`
	Popularity  int             `gorm:"default:0" json:"popularity"`

	Products []Product `gorm:"foreignKey:BrandID" json:"products,omitempty"`
}

type Product struct {
	BaseModel
	SellerID      uint            `gorm:"index;not null" json:"seller_id"`
	Title         string          `gorm:"size:200;not null" json:"title" validate:"required,min=5,max=200"`
	Description   string          `gorm:"type:text;not null" json:"description" validate:"required,min=10"`
	Category      ProductCategory `gorm:"size:20;not null;index" json:"category"`
	BrandID       *uint           `gorm:"index" json:"brand_id"`
	BrandName     string          `gorm:"size:100" json:"brand_name"`
	OriginalPrice float64         `gorm:"type:decimal(12,2)" json:"original_price"`
	Price         float64         `gorm:"type:decimal(12,2);not null" json:"price" validate:"required,gt=0"`
	Condition     string          `gorm:"size:50" json:"condition"`
	Color         string          `gorm:"size:50" json:"color"`
	Size          string          `gorm:"size:50" json:"size"`
	Material      string          `gorm:"size:100" json:"material"`
	Stock         int             `gorm:"default:1" json:"stock" validate:"required,gt=0"`
	Status        ProductStatus   `gorm:"size:20;not null;default:draft" json:"status"`
	Views         int             `gorm:"default:0" json:"views"`
	Favorites     int             `gorm:"default:0" json:"favorites"`
	IsAuthenticated bool          `gorm:"default:false" json:"is_authenticated"`

	Seller      *User           `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Brand       *Brand          `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Images      []ProductImage  `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Order       *Order          `json:"order,omitempty"`
}

type ProductImage struct {
	BaseModel
	ProductID   uint   `gorm:"index;not null" json:"product_id"`
	ImageURL    string `gorm:"size:500;not null" json:"image_url"`
	ImageType   string `gorm:"size:20" json:"image_type"`
	IsPrimary   bool   `gorm:"default:false" json:"is_primary"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`

	Product     *Product `gorm:"foreignKey:ProductID" json:"-"`
}

type Authentication struct {
	BaseModel
	OrderID          uint                 `gorm:"index;not null" json:"order_id"`
	ProductID        uint                 `gorm:"index;not null" json:"product_id"`
	BuyerID          uint                 `gorm:"index;not null" json:"buyer_id"`
	AuthenticatorID  *uint                `gorm:"index" json:"authenticator_id"`
	Status           AuthenticationStatus `gorm:"size:20;not null;default:pending" json:"status"`
	Result           *AuthenticationResult `gorm:"size:20" json:"result"`
	ReportFile       string               `gorm:"size:255" json:"report_file"`
	ReportContent    string               `gorm:"type:text" json:"report_content"`
	AuthenticatorNotes string             `gorm:"type:text" json:"authenticator_notes"`
	RejectionReason  string               `gorm:"size:500" json:"rejection_reason"`
	AcceptedAt       *time.Time           `json:"accepted_at"`
	CompletedAt      *time.Time           `json:"completed_at"`

	Order           *Order           `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product         *Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Buyer           *User            `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Authenticator   *User            `gorm:"foreignKey:AuthenticatorID" json:"authenticator,omitempty"`
}

type Order struct {
	BaseModel
	OrderNumber     string        `gorm:"uniqueIndex;size:50;not null" json:"order_number"`
	BuyerID         uint          `gorm:"index;not null" json:"buyer_id"`
	SellerID        uint          `gorm:"index;not null" json:"seller_id"`
	ProductID       uint          `gorm:"index;not null" json:"product_id"`
	Price           float64       `gorm:"type:decimal(12,2);not null" json:"price"`
	Status          OrderStatus   `gorm:"size:20;not null;default:pending;index" json:"status"`
	PaymentStatus   PaymentStatus `gorm:"size:20;not null;default:pending" json:"payment_status"`
	PaymentMethod   string        `gorm:"size:50" json:"payment_method"`
	PaymentTime     *time.Time    `json:"payment_time"`
	ShippingAddress string        `gorm:"size:500" json:"shipping_address"`
	TrackingNumber  string        `gorm:"size:50" json:"tracking_number"`
	ShippedAt       *time.Time    `json:"shipped_at"`
	DeliveredAt     *time.Time    `json:"delivered_at"`
	CompletedAt     *time.Time    `json:"completed_at"`
	CancelledAt     *time.Time    `json:"cancelled_at"`
	CancelReason    string        `gorm:"size:500" json:"cancel_reason"`
	NeedAuth        bool          `gorm:"default:false" json:"need_auth"`
	Remark          string        `gorm:"size:500" json:"remark"`

	Buyer           *User           `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller          *User           `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Product         *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Authentication  *Authentication `gorm:"foreignKey:OrderID" json:"authentication,omitempty"`
	Review          *Review         `json:"review,omitempty"`
}

type Review struct {
	BaseModel
	OrderID     uint         `gorm:"uniqueIndex;not null" json:"order_id"`
	ReviewerID  uint         `gorm:"index;not null" json:"reviewer_id"`
	RevieweeID  uint         `gorm:"index;not null" json:"reviewee_id"`
	Rating      ReviewRating `gorm:"not null" json:"rating" validate:"required,min=1,max=5"`
	Content     string       `gorm:"type:text" json:"content"`
	Images      string       `gorm:"type:text" json:"images"`
	IsAnonymous bool         `gorm:"default:false" json:"is_anonymous"`

	Order       *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Reviewer    *User   `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Reviewee    *User   `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
}

type AuditLog struct {
	BaseModel
	UserID    *uint      `gorm:"index" json:"user_id"`
	Username  string     `gorm:"size:100" json:"username"`
	Action    string     `gorm:"size:100;not null;index" json:"action"`
	Module    string     `gorm:"size:50;index" json:"module"`
	IP        string     `gorm:"size:50" json:"ip"`
	UserAgent string     `gorm:"size:500" json:"user_agent"`
	Method    string     `gorm:"size:10" json:"method"`
	Path      string     `gorm:"size:500" json:"path"`
	Params    string     `gorm:"type:text" json:"params"`
	Result    string     `gorm:"size:20" json:"result"`
	Detail    string     `gorm:"type:text" json:"detail"`

	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Statistic struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Date        time.Time `gorm:"uniqueIndex" json:"date"`
	TotalOrders int       `json:"total_orders"`
	TotalAmount float64   `json:"total_amount"`
	NewUsers    int       `json:"new_users"`
	AuthCount   int       `json:"auth_count"`
	AuthPassRate float64 `json:"auth_pass_rate"`
}

var DB *gorm.DB

func InitDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * 60)

	if err := db.AutoMigrate(
		&User{},
		&AuthenticatorProfile{},
		&SellerProfile{},
		&BuyerProfile{},
		&Brand{},
		&Product{},
		&ProductImage{},
		&Authentication{},
		&Order{},
		&Review{},
		&AuditLog{},
		&Statistic{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	DB = db
	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}

func SetLogger(log *logrus.Logger) {
	DB.Logger = logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             200,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
