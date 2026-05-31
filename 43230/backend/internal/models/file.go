package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileUpload struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	FileName      string    `json:"file_name"`
	OriginalName  string    `json:"original_name"`
	FileType      string    `json:"file_type"`
	FileSize      int64     `json:"file_size"`
	FileHash      string    `json:"file_hash"`
	StoragePath   string    `json:"-"`
	UploadStatus  string    `json:"upload_status" gorm:"default:uploading"`
	TotalChunks   int       `json:"total_chunks"`
	UploadedChunks int      `json:"uploaded_chunks" gorm:"default:0"`
	ChunkSize     int64     `json:"chunk_size"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
}

type FileAccessLog struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	FileID     uuid.UUID `json:"file_id" gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	AccessType string    `json:"access_type"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

type DownloadRecord struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	ModelID    uuid.UUID `json:"model_id" gorm:"type:uuid;not null;index"`
	Model      *Model3D  `json:"model,omitempty" gorm:"foreignKey:ModelID"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PurchaseID uuid.UUID `json:"purchase_id" gorm:"type:uuid;index"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}

type Notification struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	RelatedID uuid.UUID `json:"related_id" gorm:"type:uuid"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
}

type Transaction struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrderID       uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	Order         *PrintOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	BalanceAfter  float64   `json:"balance_after"`
	Description   string    `json:"description"`
	PaymentMethod string    `json:"payment_method"`
	TransactionNo string    `json:"transaction_no" gorm:"uniqueIndex"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
