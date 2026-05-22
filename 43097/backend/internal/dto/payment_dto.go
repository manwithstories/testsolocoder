package dto

import (
	"hotel-system/internal/model"
)

type PaymentCreateRequest struct {
	OrderType     model.OrderType     `json:"orderType" binding:"required,oneof=booking checkin"`
	OrderID       uint                `json:"orderId" binding:"required"`
	Amount        float64             `json:"amount" binding:"required,gt=0"`
	PaymentMethod model.PaymentMethod `json:"paymentMethod" binding:"required,oneof=cash wechat alipay card transfer"`
	PaymentType   model.PaymentType   `json:"paymentType" binding:"required,oneof=prepaid extra refund deposit"`
	TransactionID string              `json:"transactionId"`
	Remark        string              `json:"remark"`
}

type PaymentResponse struct {
	ID             uint                `json:"id"`
	PaymentNo      string              `json:"paymentNo"`
	OrderType      model.OrderType     `json:"orderType"`
	OrderID        uint                `json:"orderId"`
	Amount         float64             `json:"amount"`
	PaymentMethod  model.PaymentMethod `json:"paymentMethod"`
	PaymentType    model.PaymentType   `json:"paymentType"`
	Status         model.PaymentStatus `json:"status"`
	TransactionID  string              `json:"transactionId"`
	Remark         string              `json:"remark"`
	CreatedAt      string              `json:"createdAt"`
}

type PaymentListRequest struct {
	PaginationRequest
	PaymentNo      string              `form:"paymentNo"`
	OrderType      model.OrderType     `form:"orderType"`
	OrderID        uint                `form:"orderId"`
	PaymentMethod  model.PaymentMethod `form:"paymentMethod"`
	PaymentType    model.PaymentType   `form:"paymentType"`
	Status         model.PaymentStatus `form:"status"`
	StartDate      string              `form:"startDate"`
	EndDate        string              `form:"endDate"`
}

type PaymentRefundRequest struct {
	PaymentID     uint                `json:"paymentId" binding:"required"`
	RefundAmount  float64             `json:"refundAmount" binding:"required,gt=0"`
	PaymentMethod model.PaymentMethod `json:"paymentMethod" binding:"required,oneof=cash wechat alipay card transfer"`
	Reason        string              `json:"reason" binding:"required"`
	TransactionID string              `json:"transactionId"`
}

type OrderPaymentsRequest struct {
	OrderType model.OrderType `form:"orderType" binding:"required,oneof=booking checkin"`
	OrderID   uint            `form:"orderId" binding:"required"`
}

type OrderPaymentsResponse struct {
	Payments     []PaymentResponse `json:"payments"`
	TotalPaid    float64           `json:"totalPaid"`
	TotalRefund  float64           `json:"totalRefund"`
}

type PaymentVoucherResponse struct {
	PaymentNo     string              `json:"paymentNo"`
	OrderType     model.OrderType     `json:"orderType"`
	OrderID       uint                `json:"orderId"`
	Amount        float64             `json:"amount"`
	PaymentMethod model.PaymentMethod `json:"paymentMethod"`
	PaymentType   model.PaymentType   `json:"paymentType"`
	Status        model.PaymentStatus `json:"status"`
	TransactionID string              `json:"transactionId"`
	Remark        string              `json:"remark"`
	CreatedAt     string              `json:"createdAt"`
	GuestName     string              `json:"guestName,omitempty"`
	RoomNo        string              `json:"roomNo,omitempty"`
}
