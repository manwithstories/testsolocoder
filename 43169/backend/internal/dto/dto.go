package dto

import "time"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
	Code     string `json:"code" binding:"required,len=6"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  UserInfo   `json:"user"`
}

type UserInfo struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Phone        string     `json:"phone"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	VerifyStatus string     `json:"verify_status"`
	RealName     string     `json:"real_name"`
	Avatar       string     `json:"avatar"`
	MemberLevel  string     `json:"member_level"`
	MemberExpire *time.Time `json:"member_expire"`
	Profile      *ProfileInfo `json:"profile,omitempty"`
}

type ProfileInfo struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Nickname        string    `json:"nickname"`
	Gender          string    `json:"gender"`
	Birthday        *time.Time `json:"birthday"`
	Age             int       `json:"age"`
	Height          int       `json:"height"`
	Weight          int       `json:"weight"`
	Education       string    `json:"education"`
	Occupation      string    `json:"occupation"`
	Income          string    `json:"income"`
	City            string    `json:"city"`
	District        string    `json:"district"`
	Intro           string    `json:"intro"`
	Hobbies         string    `json:"hobbies"`
	Tags            string    `json:"tags"`
	Photos          []string  `json:"photos"`
}

type ProfileUpdateRequest struct {
	Nickname        string    `json:"nickname"`
	Gender          string    `json:"gender"`
	Birthday        *time.Time `json:"birthday"`
	Height          int       `json:"height"`
	Weight          int       `json:"weight"`
	Education       string    `json:"education"`
	Occupation      string    `json:"occupation"`
	Income          string    `json:"income"`
	City            string    `json:"city"`
	District        string    `json:"district"`
	Address         string    `json:"address"`
	Intro           string    `json:"intro"`
	Hobbies         string    `json:"hobbies"`
	Tags            string    `json:"tags"`
	MinAge          int       `json:"min_age"`
	MaxAge          int       `json:"max_age"`
	MinHeight       int       `json:"min_height"`
	MaxHeight       int       `json:"max_height"`
	PreferEducation string    `json:"prefer_education"`
	PreferIncome    string    `json:"prefer_income"`
	PreferCity      string    `json:"prefer_city"`
}

type VerifyRequest struct {
	RealName    string `json:"real_name" binding:"required"`
	IDCard      string `json:"id_card" binding:"required"`
	IDCardFront string `json:"id_card_front" binding:"required"`
	IDCardBack  string `json:"id_card_back" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	SmsCode     string `json:"sms_code" binding:"required,len=6"`
}

type MatchFilterRequest struct {
	Gender        string   `form:"gender"`
	MinAge        int      `form:"min_age"`
	MaxAge        int      `form:"max_age"`
	MinHeight     int      `form:"min_height"`
	MaxHeight     int      `form:"max_height"`
	Education     string   `form:"education"`
	Income        string   `form:"income"`
	City          string   `form:"city"`
	Tags          []string `form:"tags"`
	Page          int      `form:"page"`
	PageSize      int      `form:"page_size"`
}

type MatchResultItem struct {
	UserID      uint        `json:"user_id"`
	ProfileInfo ProfileInfo `json:"profile"`
	MatchScore  float64     `json:"match_score"`
	MatchReason string      `json:"match_reason"`
	IsFavorited bool        `json:"is_favorited"`
	IsBlocked   bool        `json:"is_blocked"`
}

type DateInviteRequest struct {
	ReceiverID uint      `json:"receiver_id" binding:"required"`
	Title      string    `json:"title" binding:"required"`
	Location   string    `json:"location"`
	DateAt     time.Time `json:"date_at" binding:"required"`
	Duration   int       `json:"duration"`
	Note       string    `json:"note"`
}

type DateReviewRequest struct {
	DateID  uint   `json:"date_id" binding:"required"`
	TargetID uint  `json:"target_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
}

type ChatMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Type       string `json:"type"`
	Content    string `json:"content" binding:"required"`
}

type MemberUpgradeRequest struct {
	Level  string `json:"level" binding:"required"`
	Months int    `json:"months" binding:"required,min=1,max=24"`
}

type MatchmakerServiceRequest struct {
	MemberAID   uint   `json:"member_a_id" binding:"required"`
	MemberBID   uint   `json:"member_b_id" binding:"required"`
	ServiceType string `json:"service_type" binding:"required"`
	Note        string `json:"note"`
	DateID      *uint  `json:"date_id"`
}

type MemberManageRequest struct {
	MemberID uint `json:"member_id" binding:"required"`
}

type StatsFilterRequest struct {
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	MatchmakerID uint  `form:"matchmaker_id"`
}
