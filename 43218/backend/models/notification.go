package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Type      string         `gorm:"size:30;not null" json:"type"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	OrderNo   string         `gorm:"size:50" json:"order_no"`
	IsRead    bool           `gorm:"default:false;index" json:"is_read"`
	ExtraData string         `gorm:"type:text" json:"extra_data"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}

const (
	NotificationTypeOrderStatus   = "order_status"
	NotificationTypePayment       = "payment"
	NotificationTypeRepairProgress = "repair_progress"
	NotificationTypeReview        = "review"
	NotificationTypeSystem        = "system"
	NotificationTypeReport        = "report"
	NotificationTypeWarranty      = "warranty"
	NotificationTypeMessage       = "message"
)

type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SenderID  uint           `gorm:"index;not null" json:"sender_id"`
	ReceiverID uint          `gorm:"index;not null" json:"receiver_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Type      string         `gorm:"size:20;default:text" json:"type"`
	IsRead    bool           `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Sender   User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}

type Transaction struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	Type            string         `gorm:"size:20;not null" json:"type"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	OrderNo         string         `gorm:"size:50" json:"order_no"`
	PaymentMethod   string         `gorm:"size:20" json:"payment_method"`
	TransactionNo   string         `gorm:"uniqueIndex;size:100" json:"transaction_no"`
	Status          int            `gorm:"default:1" json:"status"`
	Remark          string         `gorm:"size:255" json:"remark"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Transaction) TableName() string {
	return "transactions"
}

const (
	TransactionTypeRecharge = "recharge"
	TransactionTypeWithdraw = "withdraw"
	TransactionTypePayment  = "payment"
	TransactionTypeIncome   = "income"
	TransactionTypeRefund   = "refund"
	TransactionTypeCommission = "commission"
)

const (
	TransactionStatusPending  = 1
	TransactionStatusSuccess  = 2
	TransactionStatusFailed   = 3
	TransactionStatusCancelled = 4
)

type AdminLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	AdminID   uint           `gorm:"index;not null" json:"admin_id"`
	Action    string         `gorm:"size:50;not null" json:"action"`
	TargetType string        `gorm:"size:50" json:"target_type"`
	TargetID  *uint          `json:"target_id"`
	Content   string         `gorm:"type:text" json:"content"`
	IPAddress string         `gorm:"size:50" json:"ip_address"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AdminLog) TableName() string {
	return "admin_logs"
}
