package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleStore  Role = "store"
	RoleKeeper Role = "keeper"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone        string         `gorm:"size:20" json:"phone"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Role         Role           `gorm:"type:varchar(20);not null;default:owner" json:"role"`
	AvatarURL    string         `gorm:"size:500" json:"avatar_url"`
	RealName     string         `gorm:"size:50" json:"real_name"`
	StoreInfo    *StoreInfo     `gorm:"foreignKey:UserID" json:"store_info,omitempty"`
	KeeperInfo   *KeeperInfo    `gorm:"foreignKey:UserID" json:"keeper_info,omitempty"`
	Status       string         `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type StoreInfo struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	StoreName     string    `gorm:"size:100;not null" json:"store_name"`
	Address       string    `gorm:"size:500" json:"address"`
	LicenseNo     string    `gorm:"size:50" json:"license_no"`
	BusinessHours string    `gorm:"size:200" json:"business_hours"`
	Description   string    `gorm:"size:1000" json:"description"`
	Rating        float64   `gorm:"default:0" json:"rating"`
	ReviewCount   int       `gorm:"default:0" json:"review_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type KeeperInfo struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	StoreID       uuid.UUID `gorm:"type:uuid;index" json:"store_id"`
	NickName      string    `gorm:"size:50" json:"nick_name"`
	Experience    int       `gorm:"default:0" json:"experience"`
	Specialty     string    `gorm:"size:500" json:"specialty"`
	Rating        float64   `gorm:"default:0" json:"rating"`
	ReviewCount   int       `gorm:"default:0" json:"review_count"`
	Certification string    `gorm:"size:500" json:"certification"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Pet struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	Name          string         `gorm:"size:50;not null" json:"name"`
	Species       string         `gorm:"size:20;not null" json:"species"`
	Breed         string         `gorm:"size:50" json:"breed"`
	Gender        string         `gorm:"size:10" json:"gender"`
	BirthDate     *time.Time     `json:"birth_date"`
	Weight        float64        `gorm:"default:0" json:"weight"`
	Color         string         `gorm:"size:30" json:"color"`
	AvatarURL     string         `gorm:"size:500" json:"avatar_url"`
	Allergies     string         `gorm:"size:500" json:"allergies"`
	DietHabit     string         `gorm:"size:500" json:"diet_habit"`
	Temperament   string         `gorm:"size:200" json:"temperament"`
	VaccineRecords []VaccineRecord `gorm:"foreignKey:PetID" json:"vaccine_records,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type VaccineRecord struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PetID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"pet_id"`
	VaccineName  string         `gorm:"size:100;not null" json:"vaccine_name"`
	VaccinatedAt time.Time      `gorm:"not null" json:"vaccinated_at"`
	ExpireAt     time.Time      `gorm:"not null" json:"expire_at"`
	Hospital     string         `gorm:"size:100" json:"hospital"`
	ProofURL     string         `gorm:"size:500" json:"proof_url"`
	IsValid      bool           `gorm:"default:true" json:"is_valid"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type DewormRecord struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PetID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"pet_id"`
	DewormType  string         `gorm:"size:50;not null" json:"deworm_type"`
	DewormedAt  time.Time      `gorm:"not null" json:"dewormed_at"`
	ExpireAt    time.Time      `gorm:"not null" json:"expire_at"`
	Medicine    string         `gorm:"size:100" json:"medicine"`
	IsValid     bool           `gorm:"default:true" json:"is_valid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type BoardingPackage struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Type         string         `gorm:"type:varchar(20);not null" json:"type"`
	Description  string         `gorm:"size:1000" json:"description"`
	PricePerDay  float64        `gorm:"not null" json:"price_per_day"`
	Capacity     int            `gorm:"default:1" json:"capacity"`
	Features     string         `gorm:"size:1000" json:"features"`
	IsAvailable  bool           `gorm:"default:true" json:"is_available"`
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	PackageTypeDaycare  = "daycare"
	PackageTypeBoarding = "boarding"
)

type Reservation struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNo        string         `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	PetID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"pet_id"`
	StoreID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	PackageID      uuid.UUID      `gorm:"type:uuid;not null" json:"package_id"`
	PackageType    string         `gorm:"type:varchar(20);not null" json:"package_type"`
	CheckInDate    time.Time      `gorm:"not null" json:"check_in_date"`
	CheckOutDate   time.Time      `gorm:"not null" json:"check_out_date"`
	TotalDays      int            `gorm:"not null" json:"total_days"`
	TotalAmount    float64        `gorm:"not null" json:"total_amount"`
	DepositAmount  float64        `gorm:"default:0" json:"deposit_amount"`
	Status         string         `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	KeeperID       *uuid.UUID     `gorm:"type:uuid;index" json:"keeper_id,omitempty"`
	Remark         string         `gorm:"size:500" json:"remark"`
	CancelReason   string         `gorm:"size:500" json:"cancel_reason"`
	CancelledAt    *time.Time     `json:"cancelled_at,omitempty"`
	CompletedAt    *time.Time     `json:"completed_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	ReservationStatusPending   = "pending"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusCheckedIn = "checked_in"
	ReservationStatusCompleted = "completed"
	ReservationStatusCancelled = "cancelled"
)

type DailyRecord struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ReservationID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"reservation_id"`
	PetID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"pet_id"`
	KeeperID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"keeper_id"`
	RecordDate     time.Time      `gorm:"not null" json:"record_date"`
	FeedStatus     string         `gorm:"size:200" json:"feed_status"`
	Activity       string         `gorm:"size:500" json:"activity"`
	HealthStatus   string         `gorm:"size:500" json:"health_status"`
	Mood           string         `gorm:"size:200" json:"mood"`
	Photos         string         `gorm:"size:2000" json:"photos"`
	Remark         string         `gorm:"size:500" json:"remark"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Review struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ReservationID  uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"reservation_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	StoreID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	KeeperID       uuid.UUID      `gorm:"type:uuid;index" json:"keeper_id,omitempty"`
	StoreRating    int            `gorm:"not null;check:store_rating >= 1 AND store_rating <= 5" json:"store_rating"`
	KeeperRating   int            `gorm:"check:keeper_rating >= 1 AND keeper_rating <= 5" json:"keeper_rating,omitempty"`
	Content        string         `gorm:"size:1000" json:"content"`
	Images         string         `gorm:"size:2000" json:"images"`
	Reply          string         `gorm:"size:1000" json:"reply"`
	ReplyAt        *time.Time     `json:"reply_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Order struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNo        string         `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	ReservationID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"reservation_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	StoreID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	Type           string         `gorm:"type:varchar(20);not null" json:"type"`
	Amount         float64        `gorm:"not null" json:"amount"`
	PayStatus      string         `gorm:"type:varchar(20);not null;default:unpaid" json:"pay_status"`
	PayMethod      string         `gorm:"size:30" json:"pay_method"`
	TransactionID  string         `gorm:"size:100" json:"transaction_id"`
	PaidAt         *time.Time     `json:"paid_at,omitempty"`
	RefundAmount   float64        `gorm:"default:0" json:"refund_amount"`
	RefundAt       *time.Time     `json:"refund_at,omitempty"`
	Remark         string         `gorm:"size:500" json:"remark"`
	AmountHash     string         `gorm:"size:100" json:"amount_hash"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	OrderTypePrepay    = "prepay"
	OrderTypeSettlement = "settlement"
	OrderTypeRefund    = "refund"

	PayStatusUnpaid = "unpaid"
	PayStatusPaid   = "paid"
	PayStatusRefund = "refunded"
)

type HealthAlert struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	PetID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"pet_id"`
	AlertType    string         `gorm:"type:varchar(30);not null" json:"alert_type"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Content      string         `gorm:"size:500" json:"content"`
	RecordID     uuid.UUID      `gorm:"type:uuid;index" json:"record_id"`
	ExpireAt     time.Time      `json:"expire_at"`
	IsRead       bool           `gorm:"default:false" json:"is_read"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	AlertTypeVaccineExpire = "vaccine_expire"
	AlertTypeDewormExpire  = "deworm_expire"
)

type OperationLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Username    string    `gorm:"size:50" json:"username"`
	Action      string    `gorm:"size:100;not null" json:"action"`
	Method      string    `gorm:"size:10" json:"method"`
	URL         string    `gorm:"size:500" json:"url"`
	IP          string    `gorm:"size:50" json:"ip"`
	Params      string    `gorm:"type:text" json:"params"`
	Result      string    `gorm:"type:text" json:"result"`
	Status      int       `json:"status"`
	ExecTime    int64     `json:"exec_time"`
	ErrorMsg    string    `gorm:"size:1000" json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}
