package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelStatus string

const (
	ModelStatusDraft     ModelStatus = "draft"
	ModelStatusPublished ModelStatus = "published"
	ModelStatusRejected  ModelStatus = "rejected"
	ModelStatusBanned    ModelStatus = "banned"
)

type LicenseType string

const (
	LicensePerPurchase LicenseType = "per_purchase"
	LicenseSubscription LicenseType = "subscription"
	LicenseCommercial  LicenseType = "commercial"
)

type Model3D struct {
	ID             uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	DesignerID     uuid.UUID   `json:"designer_id" gorm:"type:uuid;not null;index"`
	Designer       *User       `json:"designer,omitempty" gorm:"foreignKey:DesignerID"`
	Title          string      `json:"title" gorm:"not null"`
	Description    string      `json:"description"`
	Category       string      `json:"category"`
	Tags           []string    `json:"tags" gorm:"type:text[]"`
	Price          float64     `json:"price" gorm:"not null"`
	LicenseType    LicenseType `json:"license_type" gorm:"not null"`
	SubscriptionPrice float64  `json:"subscription_price"`
	Status         ModelStatus `json:"status" gorm:"default:draft"`
	ThumbnailURL   string      `json:"thumbnail_url"`
	FileURL        string      `json:"-"`
	FileSize       int64       `json:"file_size"`
	FileType       string      `json:"file_type"`
	FileHash       string      `json:"-"`
	Version        string      `json:"version" gorm:"default:1.0.0"`
	Volume         float64     `json:"volume"`
	BoundingBox    string      `json:"bounding_box"`
	PrintTimeHours float64     `json:"print_time_hours"`
	RecommendedMaterials []string `json:"recommended_materials" gorm:"type:text[]"`
	PolygonCount   int         `json:"polygon_count"`
	ViewCount      int         `json:"view_count" gorm:"default:0"`
	DownloadCount  int         `json:"download_count" gorm:"default:0"`
	PurchaseCount  int         `json:"purchase_count" gorm:"default:0"`
	FavoriteCount  int         `json:"favorite_count" gorm:"default:0"`
	Rating         float64     `json:"rating" gorm:"default:5.0"`
	RatingCount    int         `json:"rating_count" gorm:"default:0"`
	IsFeatured     bool        `json:"is_featured" gorm:"default:false"`
	PreviousVersions []ModelVersion `json:"previous_versions,omitempty" gorm:"foreignKey:ModelID"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type ModelVersion struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	ModelID       uuid.UUID `json:"model_id" gorm:"type:uuid;not null;index"`
	VersionNumber string    `json:"version_number"`
	FileURL       string    `json:"-"`
	FileHash      string    `json:"-"`
	ChangeLog     string    `json:"change_log"`
	CreatedAt     time.Time `json:"created_at"`
}

type ModelPurchase struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	ModelID       uuid.UUID `json:"model_id" gorm:"type:uuid;not null;index"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Model         *Model3D  `json:"model,omitempty" gorm:"foreignKey:ModelID"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PurchaseType  LicenseType `json:"purchase_type"`
	Amount        float64   `json:"amount"`
	TransactionID string    `json:"transaction_id"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type ModelFavorite struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	ModelID   uuid.UUID `json:"model_id" gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `json:"created_at"`
}
