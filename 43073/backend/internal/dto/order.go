package dto

type OrderTicket struct {
	TicketTypeID uint `json:"ticketTypeId" validate:"required"`
	Quantity    int  `json:"quantity" validate:"required,min=1"`
}

type OrderCreateRequest struct {
	ActivityID uint         `json:"activityId" validate:"required"`
	Tickets  []OrderTicket `json:"tickets" validate:"required,min=1"`
	CouponCode string    `json:"couponCode"`
	Remark     string    `json:"remark"`
}

type OrderListRequest struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"pageSize,default=10"`
	ActivityID uint   `form:"activityId"`
	Status     string `form:"status"`
	OrderNo    string `form:"orderNo"`
}

type OrderPayRequest struct {
	PayMethod string `json:"payMethod" validate:"required"`
}

type OrderCancelRequest struct {
	Reason string `json:"reason"`
}
