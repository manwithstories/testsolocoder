package dto

import "time"

type CreateServiceReq struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Region      string     `json:"region" binding:"required"`
	Address     string     `json:"address"`
	ServiceDate *time.Time `json:"service_date"`
	ServiceTime string     `json:"service_time"`
	Duration    int        `json:"duration"`
	BudgetMin   float64    `json:"budget_min"`
	BudgetMax   float64    `json:"budget_max"`
	Images      string     `json:"images"`
	Remark      string     `json:"remark"`
}

type ServiceQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status"`
	Region   string `form:"region"`
	UserID   uint   `form:"user_id"`
	PilotID  uint   `form:"pilot_id"`
}

type CreateBidReq struct {
	ServiceID uint    `json:"service_id" binding:"required"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	Message   string  `json:"message"`
}

type AcceptBidReq struct {
	BidID uint `json:"bid_id" binding:"required"`
}

type UpdateServiceStatusReq struct {
	ServiceID uint   `json:"service_id" binding:"required"`
	Status    string `json:"status" binding:"required,oneof=progress completed cancelled"`
}
