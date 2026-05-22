package dto

import (
	"hotel-system/internal/model"
	"time"
)

type BookingCreateRequest struct {
	RoomID      uint      `json:"roomId" binding:"required,min=1"`
	MemberID    *uint     `json:"memberId"`
	GuestName   string    `json:"guestName" binding:"required,max=50"`
	GuestPhone  string    `json:"guestPhone" binding:"required,max=20"`
	GuestIDCard string    `json:"guestIdCard" binding:"max=20"`
	CheckInDate time.Time `json:"checkInDate" binding:"required"`
	CheckOutDate time.Time `json:"checkOutDate" binding:"required,gtfield=CheckInDate"`
	Remarks     string    `json:"remarks" binding:"max=500"`
}

type BookingUpdateRequest struct {
	RoomID      *uint     `json:"roomId" binding:"omitempty,min=1"`
	GuestName   string    `json:"guestName" binding:"omitempty,max=50"`
	GuestPhone  string    `json:"guestPhone" binding:"omitempty,max=20"`
	GuestIDCard string    `json:"guestIdCard" binding:"omitempty,max=20"`
	CheckInDate time.Time `json:"checkInDate" binding:"omitempty"`
	CheckOutDate time.Time `json:"checkOutDate" binding:"omitempty,gtfield=CheckInDate"`
	Remarks     string    `json:"remarks" binding:"omitempty,max=500"`
}

type BookingResponse struct {
	ID             uint                `json:"id"`
	BookingNo      string              `json:"bookingNo"`
	RoomID         uint                `json:"roomId"`
	Room           *RoomResponse       `json:"room,omitempty"`
	MemberID       *uint               `json:"memberId,omitempty"`
	Member         *MemberResponse     `json:"member,omitempty"`
	GuestName      string              `json:"guestName"`
	GuestPhone     string              `json:"guestPhone"`
	GuestIDCard    string              `json:"guestIdCard"`
	CheckInDate    string              `json:"checkInDate"`
	CheckOutDate   string              `json:"checkOutDate"`
	Days           int                 `json:"days"`
	TotalPrice     float64             `json:"totalPrice"`
	Status         model.BookingStatus `json:"status"`
	PaidAmount     float64             `json:"paidAmount"`
	Remarks        string              `json:"remarks,omitempty"`
	CancelDeadline *string             `json:"cancelDeadline,omitempty"`
	CreatedAt      string              `json:"createdAt"`
	UpdatedAt      string              `json:"updatedAt"`
}

type MemberResponse struct {
	ID        uint                `json:"id"`
	MemberNo  string              `json:"memberNo"`
	Name      string              `json:"name"`
	Phone     string              `json:"phone"`
	Email     string              `json:"email"`
	IDCard    string              `json:"idCard"`
	LevelID   uint                `json:"levelId"`
	Level     *MemberLevelResponse `json:"level,omitempty"`
	LevelName string              `json:"levelName,omitempty"`
	Points    int                 `json:"points"`
	Balance   float64             `json:"balance"`
	Status    model.MemberStatus `json:"status"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
}

type MemberLevelResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	DiscountRate float64 `json:"discountRate"`
	PointsRate   float64 `json:"pointsRate"`
	MinPoints    int     `json:"minPoints"`
	MaxPoints    int     `json:"maxPoints"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type BookingListRequest struct {
	PaginationRequest
	BookingNo  string              `form:"bookingNo"`
	RoomID     uint                `form:"roomId"`
	MemberID   uint                `form:"memberId"`
	GuestName  string              `form:"guestName"`
	GuestPhone string              `form:"guestPhone"`
	Status     model.BookingStatus `form:"status"`
	CheckInDate  time.Time         `form:"checkInDate"`
	CheckOutDate time.Time         `form:"checkOutDate"`
}

type BookingCancelRequest struct {
	Reason string `json:"reason" binding:"max=500"`
}

type BookingConfirmRequest struct {
	PaidAmount float64 `json:"paidAmount" binding:"min=0"`
}

type BookingPriceCalculationRequest struct {
	RoomID       uint      `json:"roomId" binding:"required,min=1"`
	MemberID     *uint     `json:"memberId"`
	CheckInDate  time.Time `json:"checkInDate" binding:"required"`
	CheckOutDate time.Time `json:"checkOutDate" binding:"required,gtfield=CheckInDate"`
}

type BookingPriceCalculationResponse struct {
	RoomPrice    float64 `json:"roomPrice"`
	Days         int     `json:"days"`
	OriginalPrice float64 `json:"originalPrice"`
	DiscountRate float64 `json:"discountRate"`
	Discount     float64 `json:"discount"`
	TotalPrice   float64 `json:"totalPrice"`
}
