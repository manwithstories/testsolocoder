package models

import (
	"time"
)

const (
	UserRoleNormal    = "normal"
	UserRoleAdmin     = "admin"
	UserRoleGuide     = "guide"
	UserRoleResearcher = "researcher"

	ReservationStatusPending   = "pending"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusCancelled = "cancelled"
	ReservationStatusCompleted = "completed"

	ApplicationStatusPending  = "pending"
	ApplicationStatusApproved = "approved"
	ApplicationStatusRejected = "rejected"

	CollectionStatusActive   = "active"
	CollectionStatusInactive = "inactive"
	CollectionStatusRepair   = "repair"

	ExhibitionStatusDraft     = "draft"
	ExhibitionStatusPublished = "published"
	ExhibitionStatusClosed    = "closed"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:64;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:128;not null"`
	Phone        string    `json:"phone" gorm:"size:20"`
	Password     string    `json:"-" gorm:"size:255;not null"`
	Nickname     string    `json:"nickname" gorm:"size:64"`
	Avatar       string    `json:"avatar" gorm:"size:255"`
	Role         string    `json:"role" gorm:"size:20;default:normal"`
	MuseumID     *uint     `json:"museum_id,omitempty"`
	Status       int       `json:"status" gorm:"default:1"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Museum       *Museum   `json:"museum,omitempty" gorm:"foreignKey:MuseumID"`
}

type Museum struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Address     string    `json:"address" gorm:"size:255"`
	Contact     string    `json:"contact" gorm:"size:64"`
	Phone       string    `json:"phone" gorm:"size:20"`
	Email       string    `json:"email" gorm:"size:128"`
	Logo        string    `json:"logo" gorm:"size:255"`
	OpenTime    string    `json:"open_time" gorm:"size:10"`
	CloseTime   string    `json:"close_time" gorm:"size:10"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CollectionCategory struct {
	ID        uint                `json:"id" gorm:"primaryKey"`
	Name      string              `json:"name" gorm:"size:64;not null"`
	ParentID  *uint               `json:"parent_id,omitempty"`
	SortOrder int                 `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Parent    *CollectionCategory `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children  []CollectionCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

type Collection struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:64;not null"`
	Era         string    `json:"era" gorm:"size:64"`
	Material    string    `json:"material" gorm:"size:128"`
	Size        string    `json:"size" gorm:"size:64"`
	Source      string    `json:"source" gorm:"size:255"`
	Condition   string    `json:"condition" gorm:"size:64"`
	Description string    `json:"description" gorm:"type:text"`
	ImageUrl    string    `json:"image_url" gorm:"size:255"`
	Status      string    `json:"status" gorm:"size:20;default:active"`
	Tags        string    `json:"tags" gorm:"size:255"`
	MuseumID    uint      `json:"museum_id" gorm:"not null"`
	ViewCount   int       `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Category    *CollectionCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Museum      *Museum   `json:"museum,omitempty" gorm:"foreignKey:MuseumID"`
}

type Exhibition struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"size:128;not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location" gorm:"size:128"`
	HallNumber  string    `json:"hall_number" gorm:"size:32"`
	TicketPrice float64   `json:"ticket_price" gorm:"default:0"`
	MaxVisitors int       `json:"max_visitors" gorm:"default:50"`
	ImageUrl    string    `json:"image_url" gorm:"size:255"`
	Status      string    `json:"status" gorm:"size:20;default:draft"`
	IsVirtual   bool      `json:"is_virtual" gorm:"default:false"`
	VirtualUrl  string    `json:"virtual_url,omitempty" gorm:"size:255"`
	MuseumID    uint      `json:"museum_id" gorm:"not null"`
	ViewCount   int       `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Collections []Collection `json:"collections,omitempty" gorm:"many2many:exhibition_collections;"`
	Museum      *Museum   `json:"museum,omitempty" gorm:"foreignKey:MuseumID"`
}

type ExhibitionCollection struct {
	ExhibitionID uint       `json:"exhibition_id" gorm:"primaryKey"`
	CollectionID uint       `json:"collection_id" gorm:"primaryKey"`
	SortOrder    int        `json:"sort_order" gorm:"default:0"`
	CreatedAt    time.Time  `json:"created_at"`
}

type TimeSlot struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ExhibitionID uint      `json:"exhibition_id" gorm:"not null"`
	Date         time.Time `json:"date"`
	StartTime    string    `json:"start_time" gorm:"size:10"`
	EndTime      string    `json:"end_time" gorm:"size:10"`
	MaxCapacity  int       `json:"max_capacity"`
	BookedCount  int       `json:"booked_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Exhibition   *Exhibition `json:"exhibition,omitempty" gorm:"foreignKey:ExhibitionID"`
}

type Reservation struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	ExhibitionID uint      `json:"exhibition_id" gorm:"not null"`
	TimeSlotID   uint      `json:"time_slot_id" gorm:"not null"`
	VisitorCount int       `json:"visitor_count" gorm:"default:1"`
	GuideType    string    `json:"guide_type" gorm:"size:20;default:standard"`
	TotalPrice   float64   `json:"total_price" gorm:"default:0"`
	Status       string    `json:"status" gorm:"size:20;default:pending"`
	QRCode       string    `json:"qr_code" gorm:"size:64"`
	Remark       string    `json:"remark" gorm:"size:255"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Exhibition   *Exhibition `json:"exhibition,omitempty" gorm:"foreignKey:ExhibitionID"`
	TimeSlot     *TimeSlot `json:"time_slot,omitempty" gorm:"foreignKey:TimeSlotID"`
}

type VisitRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	ReservationID uint     `json:"reservation_id" gorm:"not null"`
	ExhibitionID uint      `json:"exhibition_id" gorm:"not null"`
	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	Favorite     bool      `json:"favorite" gorm:"default:false"`
	Rating       int       `json:"rating" gorm:"default:0"`
	Comment      string    `json:"comment" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GuideSchedule struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	GuideID     uint      `json:"guide_id" gorm:"not null"`
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time" gorm:"size:10"`
	EndTime     string    `json:"end_time" gorm:"size:10"`
	IsAvailable bool      `json:"is_available" gorm:"default:true"`
	ReservationID *uint   `json:"reservation_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Guide       *User     `json:"guide,omitempty" gorm:"foreignKey:GuideID"`
	Reservation *Reservation `json:"reservation,omitempty" gorm:"foreignKey:ReservationID"`
}

type GuideContent struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CollectionID  uint      `json:"collection_id" gorm:"not null"`
	ExhibitionID  *uint     `json:"exhibition_id,omitempty"`
	Language      string    `json:"language" gorm:"size:10;default:zh"`
	Content       string    `json:"content" gorm:"type:text"`
	AudioUrl      string    `json:"audio_url,omitempty" gorm:"size:255"`
	SortOrder     int       `json:"sort_order" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Collection    *Collection `json:"collection,omitempty" gorm:"foreignKey:CollectionID"`
}

type ResearchApplication struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	CollectionID  uint      `json:"collection_id" gorm:"not null"`
	Purpose       string    `json:"purpose" gorm:"type:text"`
	Institution   string    `json:"institution" gorm:"size:128"`
	Status        string    `json:"status" gorm:"size:20;default:pending"`
	ReviewerID    *uint     `json:"reviewer_id,omitempty"`
	ReviewComment string    `json:"review_comment" gorm:"type:text"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Collection    *Collection `json:"collection,omitempty" gorm:"foreignKey:CollectionID"`
	Reviewer      *User     `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
}

type Statistic struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	StatDate      time.Time `json:"stat_date"`
	ExhibitionID  uint      `json:"exhibition_id"`
	VisitorCount  int       `json:"visitor_count" gorm:"default:0"`
	ReservationCount int   `json:"reservation_count" gorm:"default:0"`
	Revenue       float64   `json:"revenue" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CollectionTag struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"uniqueIndex;size:32;not null"`
	CreatedAt    time.Time `json:"created_at"`
}
