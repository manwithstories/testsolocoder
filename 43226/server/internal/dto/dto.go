package dto

import "time"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
}

type CollectionCategoryRequest struct {
	Name      string `json:"name" binding:"required,max=64"`
	ParentID  *uint  `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
}

type CollectionRequest struct {
	Name        string  `json:"name" binding:"required,max=128"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	Code        string  `json:"code" binding:"required,max=64"`
	Era         string  `json:"era" binding:"max=64"`
	Material    string  `json:"material" binding:"max=128"`
	Size        string  `json:"size" binding:"max=64"`
	Source      string  `json:"source" binding:"max=255"`
	Condition   string  `json:"condition" binding:"max=64"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image_url" binding:"max=255"`
	Status      string  `json:"status"`
	Tags        string  `json:"tags" binding:"max=255"`
}

type CollectionListQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	Keyword    string `form:"keyword"`
	CategoryID uint   `form:"category_id"`
	Status     string `form:"status"`
	Era        string `form:"era"`
	Material   string `form:"material"`
	SortBy     string `form:"sort_by,default=created_at"`
	SortOrder  string `form:"sort_order,default=desc"`
}

type ExhibitionRequest struct {
	Title       string    `json:"title" binding:"required,max=128"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Location    string    `json:"location" binding:"max=128"`
	HallNumber  string    `json:"hall_number" binding:"max=32"`
	TicketPrice float64   `json:"ticket_price"`
	MaxVisitors int       `json:"max_visitors"`
	ImageUrl    string    `json:"image_url" binding:"max=255"`
	Status      string    `json:"status"`
	IsVirtual   bool      `json:"is_virtual"`
	VirtualUrl  string    `json:"virtual_url" binding:"max=255"`
	CollectionIDs []uint  `json:"collection_ids"`
}

type ExhibitionListQuery struct {
	Page      int       `form:"page,default=1"`
	PageSize  int       `form:"page_size,default=10"`
	Keyword   string    `form:"keyword"`
	Status    string    `form:"status"`
	StartDate time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02"`
}

type ReservationRequest struct {
	ExhibitionID uint   `json:"exhibition_id" binding:"required"`
	TimeSlotID   uint   `json:"time_slot_id" binding:"required"`
	VisitorCount int    `json:"visitor_count" binding:"required,min=1,max=10"`
	GuideType    string `json:"guide_type"`
}

type ReservationCancelRequest struct {
	Reason string `json:"reason"`
}

type TimeSlotRequest struct {
	Date        time.Time `json:"date" binding:"required"`
	StartTime   string    `json:"start_time" binding:"required"`
	EndTime     string    `json:"end_time" binding:"required"`
	MaxCapacity int       `json:"max_capacity" binding:"required"`
}

type BatchTimeSlotRequest struct {
	ExhibitionID uint      `json:"exhibition_id" binding:"required"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	StartTime    string    `json:"start_time" binding:"required"`
	EndTime      string    `json:"end_time" binding:"required"`
	Interval     int       `json:"interval" binding:"required"`
	MaxCapacity  int       `json:"max_capacity" binding:"required"`
}

type GuideScheduleRequest struct {
	Date        time.Time `json:"date" binding:"required"`
	StartTime   string    `json:"start_time" binding:"required"`
	EndTime     string    `json:"end_time" binding:"required"`
	IsAvailable bool      `json:"is_available"`
}

type GuideContentRequest struct {
	CollectionID uint   `json:"collection_id" binding:"required"`
	ExhibitionID *uint  `json:"exhibition_id"`
	Language     string `json:"language" binding:"required,max=10"`
	Content      string `json:"content" binding:"required"`
	AudioUrl     string `json:"audio_url" binding:"max=255"`
	SortOrder    int    `json:"sort_order"`
}

type ResearchApplicationRequest struct {
	CollectionID uint   `json:"collection_id" binding:"required"`
	Purpose      string `json:"purpose" binding:"required"`
	Institution  string `json:"institution" binding:"required,max=128"`
}

type ApplicationReviewRequest struct {
	Status        string `json:"status" binding:"required,oneof=approved rejected"`
	ReviewComment string `json:"review_comment"`
}

type StatisticsQuery struct {
	StartDate    time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate      time.Time `form:"end_date" time_format:"2006-01-02"`
	ExhibitionID uint      `form:"exhibition_id"`
	Type         string    `form:"type"`
}

type MuseumRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
	Address     string `json:"address" binding:"max=255"`
	Contact     string `json:"contact" binding:"max=64"`
	Phone       string `json:"phone" binding:"max=20"`
	Email       string `json:"email" binding:"max=128"`
	Logo        string `json:"logo" binding:"max=255"`
	OpenTime    string `json:"open_time" binding:"max=10"`
	CloseTime   string `json:"close_time" binding:"max=10"`
}
