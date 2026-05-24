package model

import (
	"time"

	"gorm.io/gorm"
)

type RevenueType string

const (
	RevenueTypePlay      RevenueType = "play"
	RevenueTypeSubscription RevenueType = "subscription"
	RevenueTypeTicket    RevenueType = "ticket"
)

type RevenueRecord struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	ArtistID    uint           `json:"artist_id" gorm:"index;not null"`
	WorkID      *uint          `json:"work_id" gorm:"index"`
	OrderID     *uint          `json:"order_id" gorm:"index"`
	Type        RevenueType    `json:"type" gorm:"size:50;not null"`
	Amount      float64        `json:"amount"`
	PlayCount   int64          `json:"play_count"`
	Rate        float64        `json:"rate"`
	Status      int            `json:"status" gorm:"default:0"`
	Period      string         `json:"period" gorm:"size:20"`
	SettledAt   *time.Time     `json:"settled_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (RevenueRecord) TableName() string {
	return "revenue_records"
}

type WithdrawStatus int

const (
	WithdrawStatusPending  WithdrawStatus = 0
	WithdrawStatusApproved WithdrawStatus = 1
	WithdrawStatusRejected WithdrawStatus = 2
	WithdrawStatusPaid     WithdrawStatus = 3
	WithdrawStatusFailed   WithdrawStatus = 4
)

type WithdrawRequest struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"index;not null"`
	ArtistID      uint           `json:"artist_id" gorm:"index;not null"`
	Amount        float64        `json:"amount"`
	Fee           float64        `json:"fee"`
	ActualAmount  float64        `json:"actual_amount"`
	Method        string         `json:"method" gorm:"size:50"`
	Account       string         `json:"account" gorm:"size:100"`
	AccountName   string         `json:"account_name" gorm:"size:100"`
	BankName      string         `json:"bank_name" gorm:"size:100"`
	Status        WithdrawStatus `json:"status" gorm:"default:0"`
	Remark        string         `json:"remark" gorm:"size:255"`
	ApprovedAt    *time.Time     `json:"approved_at"`
	ApprovedBy    *uint          `json:"approved_by"`
	RejectReason  string         `json:"reject_reason" gorm:"size:255"`
	PaidAt        *time.Time     `json:"paid_at"`
	TransactionNo string         `json:"transaction_no" gorm:"size:100"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (WithdrawRequest) TableName() string {
	return "withdraw_requests"
}

type Subscription struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	ArtistID    uint           `json:"artist_id" gorm:"index;not null"`
	Plan        string         `json:"plan" gorm:"size:50"`
	Amount      float64        `json:"amount"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

type DailyStats struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Date           time.Time      `json:"date" gorm:"index"`
	UserID         uint           `json:"user_id" gorm:"index;not null"`
	ArtistID       uint           `json:"artist_id" gorm:"index;not null"`
	NewFollowers   int64          `json:"new_followers" gorm:"default:0"`
	NewPlays       int64          `json:"new_plays" gorm:"default:0"`
	NewLikes       int64          `json:"new_likes" gorm:"default:0"`
	NewShares      int64          `json:"new_shares" gorm:"default:0"`
	NewComments    int64          `json:"new_comments" gorm:"default:0"`
	Revenue        float64        `json:"revenue" gorm:"default:0"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (DailyStats) TableName() string {
	return "daily_stats"
}

type OperationLog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index"`
	Username    string         `json:"username" gorm:"size:50"`
	Module      string         `json:"module" gorm:"size:50"`
	Operation   string         `json:"operation" gorm:"size:100"`
	Method      string         `json:"method" gorm:"size:20"`
	Path        string         `json:"path" gorm:"size:255"`
	IP          string         `json:"ip" gorm:"size:45"`
	UserAgent   string         `json:"user_agent" gorm:"size:500"`
	Params      string         `json:"params" gorm:"type:text"`
	Result      string         `json:"result" gorm:"type:text"`
	Status      int            `json:"status"`
	Duration    int64          `json:"duration"`
	ErrorMsg    string         `json:"error_msg" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
