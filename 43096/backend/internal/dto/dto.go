package dto

import "time"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int `json:"pages"`
}

type PaginatedResponse struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type UserUpdateRequest struct {
	RealName string `json:"real_name"`
	IDCard   string `json:"id_card"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type ShowCreateRequest struct {
	Name        string    `json:"name" binding:"required,max=200"`
	Description string    `json:"description"`
	Poster      string    `json:"poster"`
	Artist      string    `json:"artist"`
	Duration    int       `json:"duration"`
	Venue       string    `json:"venue" binding:"required"`
	Address     string    `json:"address"`
	Organizer   string    `json:"organizer"`
	Status      int       `json:"status"`
}

type ShowUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty,max=200"`
	Description string `json:"description"`
	Poster      string `json:"poster"`
	Artist      string `json:"artist"`
	Duration    int    `json:"duration"`
	Venue       string `json:"venue"`
	Address     string `json:"address"`
	Organizer   string `json:"organizer"`
	Status      int    `json:"status"`
}

type ShowListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   int    `form:"status"`
	Keyword  string `form:"keyword"`
}

type SessionCreateRequest struct {
	ShowID    uint64    `json:"show_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Status    int       `json:"status,default=1"`
}

type SeatAreaCreateRequest struct {
	SessionID uint64  `json:"session_id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Color     string  `json:"color"`
	Price     float64 `json:"price" binding:"required,min=0"`
	SortOrder int     `json:"sort_order"`
}

type SeatCreateRequest struct {
	SessionID uint64 `json:"session_id" binding:"required"`
	AreaID    uint64 `json:"area_id" binding:"required"`
	Row       string `json:"row" binding:"required"`
	Col       int    `json:"col" binding:"required"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type SeatBatchCreateRequest struct {
	SessionID uint64 `json:"session_id" binding:"required"`
	AreaID    uint64 `json:"area_id" binding:"required"`
	StartRow  string `json:"start_row" binding:"required"`
	EndRow    string `json:"end_row" binding:"required"`
	StartCol  int    `json:"start_col" binding:"required"`
	EndCol    int    `json:"end_col" binding:"required"`
	XOffset   int    `json:"x_offset"`
	YOffset   int    `json:"y_offset"`
	SeatWidth int    `json:"seat_width,default=30"`
	SeatHeight int   `json:"seat_height,default=30"`
	RowGap    int    `json:"row_gap,default=10"`
	ColGap    int    `json:"col_gap,default=5"`
}

type SeatLockRequest struct {
	SessionID uint64   `json:"session_id" binding:"required"`
	SeatIDs   []uint64 `json:"seat_ids" binding:"required,min=1"`
}

type SeatUnlockRequest struct {
	SessionID uint64   `json:"session_id" binding:"required"`
	SeatIDs   []uint64 `json:"seat_ids" binding:"required,min=1"`
}

type OrderCreateRequest struct {
	SessionID    uint64   `json:"session_id" binding:"required"`
	SeatIDs      []uint64 `json:"seat_ids" binding:"required,min=1"`
	RealName     string   `json:"real_name" binding:"required"`
	IDCard       string   `json:"id_card" binding:"required"`
	Phone        string   `json:"phone" binding:"required"`
	Email        string   `json:"email" binding:"omitempty,email"`
	Remark       string   `json:"remark"`
	CouponCode   string   `json:"coupon_code"`
}

type OrderPayRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
	PayType int    `json:"pay_type" binding:"required"`
}

type OrderListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   int    `form:"status"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Keyword  string `form:"keyword"`
}

type RefundCreateRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

type RefundAuditRequest struct {
	RefundNo    string `json:"refund_no" binding:"required"`
	Status      int    `json:"status" binding:"required,oneof=1 2"`
	AuditRemark string `json:"audit_remark"`
}

type CheckinRequest struct {
	TicketNo string `json:"ticket_no" binding:"required"`
}

type StatisticsRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	ShowID    uint64 `form:"show_id"`
	SessionID uint64 `form:"session_id"`
	AreaID    uint64 `form:"area_id"`
}

type SeatChartUpdateRequest struct {
	SessionID  uint64 `json:"session_id" binding:"required"`
	Config     string `json:"config" binding:"required"`
	Background string `json:"background"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

type CouponExchangeRequest struct {
	Points int `json:"points" binding:"required,min=100"`
}

type PointsAddRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
	Points int    `json:"points" binding:"required,min=1"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=20"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

type PDFExportRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	ShowID    uint64 `form:"show_id"`
}
