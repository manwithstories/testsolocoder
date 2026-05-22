package dto

import (
	"hotel-system/internal/model"
	"time"
)

type CheckInCreateRequest struct {
	BookingID       *uint     `json:"bookingId"`
	RoomID          uint      `json:"roomId" binding:"required,min=1"`
	GuestName       string    `json:"guestName" binding:"required,max=50"`
	GuestPhone      string    `json:"guestPhone" binding:"required,max=20"`
	GuestIDCard     string    `json:"guestIdCard" binding:"max=20"`
	CheckInTime     time.Time `json:"checkInTime" binding:"required"`
	ExpectedCheckOut time.Time `json:"expectedCheckOut" binding:"required,gtfield=CheckInTime"`
	Deposit         float64   `json:"deposit" binding:"min=0"`
}

type CheckOutRequest struct {
	PaymentMethod model.PaymentMethod `json:"paymentMethod" binding:"required"`
	TransactionID string              `json:"transactionId" binding:"max=100"`
	Remark        string              `json:"remark" binding:"max=500"`
}

type ExtendStayRequest struct {
	ExtendDays int `json:"extendDays" binding:"required,min=1"`
}

type ExtraChargeRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0"`
	Description string  `json:"description" binding:"required,max=200"`
}

type CheckInResponse struct {
	ID               uint                 `json:"id"`
	CheckInNo        string               `json:"checkInNo"`
	BookingID        *uint                `json:"bookingId,omitempty"`
	RoomID           uint                 `json:"roomId"`
	Room             *RoomResponse        `json:"room,omitempty"`
	GuestName        string               `json:"guestName"`
	GuestPhone       string               `json:"guestPhone"`
	GuestIDCard      string               `json:"guestIdCard,omitempty"`
	CheckInTime      string               `json:"checkInTime"`
	ExpectedCheckOut string               `json:"expectedCheckOut"`
	ActualCheckOut   *string              `json:"actualCheckOut,omitempty"`
	Status           model.CheckInStatus  `json:"status"`
	Deposit          float64              `json:"deposit"`
	ExtraCharges     float64              `json:"extraCharges"`
	TotalAmount      float64              `json:"totalAmount"`
	CreatedAt        string               `json:"createdAt"`
	UpdatedAt        string               `json:"updatedAt"`
}

type CheckInListRequest struct {
	PaginationRequest
	CheckInNo  string              `form:"checkInNo"`
	BookingID  uint                `form:"bookingId"`
	RoomID     uint                `form:"roomId"`
	GuestName  string              `form:"guestName"`
	GuestPhone string              `form:"guestPhone"`
	Status     model.CheckInStatus `form:"status"`
	CheckInTime  time.Time         `form:"checkInTime"`
	CheckOutTime time.Time         `form:"checkOutTime"`
}

type CheckOutResponse struct {
	CheckIn     *CheckInResponse `json:"checkIn"`
	TotalAmount float64          `json:"totalAmount"`
	Deposit     float64          `json:"deposit"`
	ExtraCharges float64         `json:"extraCharges"`
	PayAmount   float64          `json:"payAmount"`
	PaymentNo   string           `json:"paymentNo,omitempty"`
}
