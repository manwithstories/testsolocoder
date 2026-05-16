package models

import (
	"fmt"
	"log"
	"secondhand-trading/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;unique;not null" json:"username"`
	Email        string    `gorm:"size:100;unique;not null" json:"email"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	Phone        string    `gorm:"size:20" json:"phone"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	CreditScore  int       `gorm:"default:100" json:"credit_score"`
	IsBanned     bool      `gorm:"default:false" json:"is_banned"`
	Role         string    `gorm:"size:20;default:'user'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Products     []Product `gorm:"foreignKey:SellerID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:BuyerID" json:"-"`
	Reviews      []Review  `gorm:"foreignKey:ReviewerID" json:"-"`
	Favorites    []Favorite `gorm:"foreignKey:UserID" json:"-"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;unique;not null" json:"name"`
	ParentID  *uint     `gorm:"index" json:"parent_id,omitempty"`
	Icon      string    `gorm:"size:255" json:"icon"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	Products  []Product `gorm:"foreignKey:CategoryID" json:"-"`
}

type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64        `gorm:"type:decimal(10,2)" json:"original_price"`
	Condition     string         `gorm:"size:20;not null" json:"condition"`
	Status        string         `gorm:"size:20;default:'on_sale';index" json:"status"`
	SellerID      uint           `gorm:"not null;index" json:"seller_id"`
	CategoryID    uint           `gorm:"not null;index" json:"category_id"`
	Tags          string         `gorm:"size:500" json:"tags"`
	Location      string         `gorm:"size:100" json:"location"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	NeedsReview   bool           `gorm:"default:false" json:"needs_review"`
	IsReviewed    bool           `gorm:"default:true" json:"is_reviewed"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Images        []ProductImage `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Seller        *User          `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Category      *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Favorites     []Favorite     `gorm:"foreignKey:ProductID" json:"-"`
	Transactions  []Transaction  `gorm:"foreignKey:ProductID" json:"-"`
}

type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	ImageURL  string    `gorm:"size:500;not null" json:"image_url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ProductID       uint      `gorm:"not null;index" json:"product_id"`
	BuyerID         uint      `gorm:"not null;index" json:"buyer_id"`
	SellerID        uint      `gorm:"not null;index" json:"seller_id"`
	Price           float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Status          string    `gorm:"size:20;default:'pending';index" json:"status"`
	BuyerConfirmed  bool      `gorm:"default:false" json:"buyer_confirmed"`
	SellerConfirmed bool      `gorm:"default:false" json:"seller_confirmed"`
	BuyerReviewed   bool      `gorm:"default:false" json:"buyer_reviewed"`
	SellerReviewed  bool      `gorm:"default:false" json:"seller_reviewed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	Product         *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Buyer           *User     `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller          *User     `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type Review struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"not null;uniqueIndex:idx_transaction_reviewer,priority:1" json:"transaction_id"`
	ReviewerID    uint      `gorm:"not null;uniqueIndex:idx_transaction_reviewer,priority:2;index" json:"reviewer_id"`
	TargetID      uint      `gorm:"not null;index" json:"target_id"`
	Rating        int       `gorm:"not null;check:rating BETWEEN 1 AND 5" json:"rating"`
	Comment       string    `gorm:"type:text" json:"comment"`
	CreatedAt     time.Time `json:"created_at"`
	Transaction   *Transaction `gorm:"foreignKey:TransactionID" json:"-"`
	Reviewer      *User     `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Target        *User     `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	ProductID *uint     `gorm:"index" json:"product_id,omitempty"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&User{},
		&Category{},
		&Product{},
		&ProductImage{},
		&Transaction{},
		&Favorite{},
		&Review{},
		&Notification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrated successfully")
	seedCategories()
}

func seedCategories() {
	var count int64
	DB.Model(&Category{}).Count(&count)
	if count > 0 {
		return
	}

	categories := []Category{
		{Name: "电子产品", SortOrder: 1},
		{Name: "数码配件", ParentID: uintPtr(1), SortOrder: 1},
		{Name: "手机", ParentID: uintPtr(1), SortOrder: 2},
		{Name: "电脑", ParentID: uintPtr(1), SortOrder: 3},
		{Name: "家居生活", SortOrder: 2},
		{Name: "家具", ParentID: uintPtr(5), SortOrder: 1},
		{Name: "厨具", ParentID: uintPtr(5), SortOrder: 2},
		{Name: "服装鞋帽", SortOrder: 3},
		{Name: "男装", ParentID: uintPtr(8), SortOrder: 1},
		{Name: "女装", ParentID: uintPtr(8), SortOrder: 2},
		{Name: "图书文具", SortOrder: 4},
		{Name: "书籍", ParentID: uintPtr(11), SortOrder: 1},
		{Name: "运动户外", SortOrder: 5},
		{Name: "母婴用品", SortOrder: 6},
		{Name: "其他", SortOrder: 99},
	}

	for _, cat := range categories {
		DB.Create(&cat)
	}
	log.Println("Categories seeded successfully")
}

func uintPtr(n uint) *uint {
	return &n
}
