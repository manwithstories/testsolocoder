package dto

type OrderListRequest struct {
	PaginationRequest
	Status string `form:"status"`
	Type   string `form:"type"`
}

type CreateOrderRequest struct {
	Type         string `json:"type" binding:"required,oneof=venue device"`
	ItemID       uint   `json:"item_id" binding:"required,min=1"`
	StartTime    string `json:"start_time" binding:"required"`
	EndTime      string `json:"end_time" binding:"required"`
	Quantity     int    `json:"quantity" binding:"min=1"`
	Purpose      string `json:"purpose" binding:"max=500"`
	ContactName  string `json:"contact_name" binding:"required,max=50"`
	ContactPhone string `json:"contact_phone" binding:"required,max=20"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"max=500"`
}

type ReviewOrderRequest struct {
	Note string `json:"note" binding:"max=500"`
}

type CalendarRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	Type      string `form:"type"`
	ItemID    uint   `form:"item_id"`
}

type CalendarEvent struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	ItemID   uint   `json:"item_id"`
	ItemName string `json:"item_name"`
	Color    string `json:"color"`
}
