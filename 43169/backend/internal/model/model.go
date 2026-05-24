package model

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleMatchmaker UserRole = "matchmaker"
	RoleAdmin     UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusBanned   UserStatus = "banned"
)

type VerifyStatus string

const (
	VerifyStatusPending  VerifyStatus = "pending"
	VerifyStatusVerified VerifyStatus = "verified"
	VerifyStatusRejected VerifyStatus = "rejected"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type Education string

const (
	EduHighSchool Education = "high_school"
	EduCollege    Education = "college"
	EduBachelor   Education = "bachelor"
	EduMaster     Education = "master"
	EduPhd        Education = "phd"
)

type IncomeLevel string

const (
	IncomeLow    IncomeLevel = "low"
	IncomeMid    IncomeLevel = "mid"
	IncomeHigh   IncomeLevel = "high"
	IncomeLuxury IncomeLevel = "luxury"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Phone        string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Email        string         `gorm:"uniqueIndex;size:100" json:"email"`
	Role         UserRole       `gorm:"size:20;default:user" json:"role"`
	Status       UserStatus     `gorm:"size:20;default:active" json:"status"`
	VerifyStatus VerifyStatus   `gorm:"size:20;default:pending" json:"verify_status"`
	RealName     string         `gorm:"size:50" json:"real_name"`
	IDCard       string         `gorm:"size:20" json:"-"`
	IDCardFront  string         `gorm:"size:255" json:"id_card_front"`
	IDCardBack   string         `gorm:"size:255" json:"id_card_back"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	MemberLevel  string         `gorm:"size:20;default:free" json:"member_level"`
	MemberExpire *time.Time     `json:"member_expire"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `gorm:"size:50" json:"last_login_ip"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Profile *Profile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

type Profile struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	UserID          uint        `gorm:"uniqueIndex;not null" json:"user_id"`
	Nickname        string      `gorm:"size:50" json:"nickname"`
	Gender          Gender      `gorm:"size:10" json:"gender"`
	Birthday        *time.Time  `json:"birthday"`
	Age             int         `json:"age"`
	Height          int         `json:"height"`
	Weight          int         `json:"weight"`
	Education       Education   `gorm:"size:20" json:"education"`
	Occupation      string      `gorm:"size:100" json:"occupation"`
	Income          IncomeLevel `gorm:"size:20" json:"income"`
	City            string      `gorm:"size:50" json:"city"`
	District        string      `gorm:"size:50" json:"district"`
	Address         string      `gorm:"size:255" json:"address"`
	Latitude        float64     `json:"latitude"`
	Longitude       float64     `json:"longitude"`
	Intro           string      `gorm:"type:text" json:"intro"`
	Hobbies         string      `gorm:"type:text" json:"hobbies"`
	Tags            string      `gorm:"type:text" json:"tags"`
	Photos          string      `gorm:"type:text" json:"photos"`
	MinAge          int         `json:"min_age"`
	MaxAge          int         `json:"max_age"`
	MinHeight       int         `json:"min_height"`
	MaxHeight       int         `json:"max_height"`
	PreferEducation Education   `gorm:"size:20" json:"prefer_education"`
	PreferIncome    IncomeLevel `gorm:"size:20" json:"prefer_income"`
	PreferCity      string      `gorm:"size:50" json:"prefer_city"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type MatchRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	TargetID    uint      `gorm:"index;not null" json:"target_id"`
	MatchScore  float64   `json:"match_score"`
	MatchReason string    `gorm:"type:text" json:"match_reason"`
	IsFavorited bool      `gorm:"default:false" json:"is_favorited"`
	IsBlocked   bool      `gorm:"default:false" json:"is_blocked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DateStatus string

const (
	DateStatusPending  DateStatus = "pending"
	DateStatusAccepted DateStatus = "accepted"
	DateStatusRejected DateStatus = "rejected"
	DateStatusCanceled DateStatus = "canceled"
	DateStatusCompleted DateStatus = "completed"
)

type DateRecord struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	InitiatorID uint       `gorm:"index;not null" json:"initiator_id"`
	ReceiverID  uint       `gorm:"index;not null" json:"receiver_id"`
	MatchmakerID *uint     `gorm:"index" json:"matchmaker_id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Location    string     `gorm:"size:200" json:"location"`
	DateAt      time.Time  `json:"date_at"`
	Duration    int        `json:"duration"`
	Status      DateStatus `gorm:"size:20;default:pending" json:"status"`
	Note        string     `gorm:"type:text" json:"note"`
	Reminded    bool       `gorm:"default:false" json:"reminded"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type DateReview struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DateID     uint      `gorm:"index;not null" json:"date_id"`
	ReviewerID uint      `gorm:"index;not null" json:"reviewer_id"`
	TargetID   uint      `gorm:"index;not null" json:"target_id"`
	Rating     int       `gorm:"not null" json:"rating"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type MatchmakerMember struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MatchmakerID uint      `gorm:"index;not null" json:"matchmaker_id"`
	MemberID     uint      `gorm:"index;not null" json:"member_id"`
	Status       string    `gorm:"size:20;default:active" json:"status"`
	JoinedAt     time.Time `json:"joined_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MatchmakerService struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	MatchmakerID uint       `gorm:"index;not null" json:"matchmaker_id"`
	MemberAID    uint       `gorm:"index;not null" json:"member_a_id"`
	MemberBID    uint       `gorm:"index;not null" json:"member_b_id"`
	DateID       *uint      `gorm:"index" json:"date_id"`
	ServiceType  string     `gorm:"size:50;not null" json:"service_type"`
	Note         string     `gorm:"type:text" json:"note"`
	Status       string     `gorm:"size:20;default:progress" json:"status"`
	Progress     int        `gorm:"default:0" json:"progress"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type MatchmakerStats struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	MatchmakerID   uint      `gorm:"uniqueIndex;not null" json:"matchmaker_id"`
	TotalMembers   int       `gorm:"default:0" json:"total_members"`
	TotalServices  int       `gorm:"default:0" json:"total_services"`
	TotalDates     int       `gorm:"default:0" json:"total_dates"`
	SuccessDates   int       `gorm:"default:0" json:"success_dates"`
	AvgRating      float64   `gorm:"default:5.0" json:"avg_rating"`
	TotalRating    int       `gorm:"default:0" json:"total_rating"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MessageType string

const (
	MsgTypeText  MessageType = "text"
	MsgTypeImage MessageType = "image"
	MsgTypeVoice MessageType = "voice"
)

type ChatMessage struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	SenderID  uint        `gorm:"index;not null" json:"sender_id"`
	ReceiverID uint       `gorm:"index;not null" json:"receiver_id"`
	Type      MessageType `gorm:"size:20;default:text" json:"type"`
	Content   string      `gorm:"type:text;not null" json:"content"`
	IsRead    bool        `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time   `json:"created_at"`
}

type ChatSession struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserAID     uint      `gorm:"index;not null" json:"user_a_id"`
	UserBID     uint      `gorm:"index;not null" json:"user_b_id"`
	LastMessage string    `gorm:"size:500" json:"last_message"`
	LastTime    time.Time `json:"last_time"`
	UnreadA     int       `gorm:"default:0" json:"unread_a"`
	UnreadB     int       `gorm:"default:0" json:"unread_b"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MemberLevel string

const (
	MemberFree    MemberLevel = "free"
	MemberSilver  MemberLevel = "silver"
	MemberGold    MemberLevel = "gold"
	MemberDiamond MemberLevel = "diamond"
)

type MemberBenefit struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	Level           MemberLevel `gorm:"uniqueIndex;size:20;not null" json:"level"`
	DailyInteract   int         `gorm:"default:5" json:"daily_interact"`
	UnlimitedChat   bool        `gorm:"default:false" json:"unlimited_chat"`
	ViewWhoLiked    bool        `gorm:"default:false" json:"view_who_liked"`
	PriorityMatch   bool        `gorm:"default:false" json:"priority_match"`
	AdvancedFilter  bool        `gorm:"default:false" json:"advanced_filter"`
	VideoChat       bool        `gorm:"default:false" json:"video_chat"`
	HideOnline      bool        `gorm:"default:false" json:"hide_online"`
	NoAds           bool        `gorm:"default:false" json:"no_ads"`
	MatchmakerAssist bool       `gorm:"default:false" json:"matchmaker_assist"`
	PricePerMonth   float64     `json:"price_per_month"`
	Description     string      `gorm:"type:text" json:"description"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type MemberOrder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Level       MemberLevel `gorm:"size:20;not null" json:"level"`
	Months      int       `json:"months"`
	Amount      float64   `json:"amount"`
	Status      string    `gorm:"size:20;default:pending" json:"status"`
	PaidAt      *time.Time `json:"paid_at"`
	ExpireAt    *time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InteractLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	TargetID    uint      `gorm:"index" json:"target_id"`
	Action      string    `gorm:"size:50" json:"action"`
	CreatedAt   time.Time `json:"created_at"`
}

type SensitiveWord struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Word     string `gorm:"uniqueIndex;size:50;not null" json:"word"`
	Category string `gorm:"size:20" json:"category"`
}

type SystemLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Module    string    `gorm:"size:50" json:"module"`
	Action    string    `gorm:"size:50" json:"action"`
	IP        string    `gorm:"size:50" json:"ip"`
	Detail    string    `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}
