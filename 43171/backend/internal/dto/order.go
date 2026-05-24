package dto

import "time"

type CreateOrderReq struct {
	DroneID      uint   `json:"drone_id" binding:"required"`
	StartDate    string `json:"start_date" binding:"required"`
	EndDate      string `json:"end_date" binding:"required"`
	Region       string `json:"region" binding:"required"`
	Address      string `json:"address"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required"`
	Remark       string `json:"remark"`
}

type OrderQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status"`
	DroneID  uint   `form:"drone_id"`
}

type PayOrderReq struct {
	OrderID uint   `json:"order_id" binding:"required"`
	PayType string `json:"pay_type" binding:"required,oneof=balance alipay wechat"`
}

type CancelOrderReq struct {
	OrderID      uint   `json:"order_id" binding:"required"`
	CancelReason string `json:"cancel_reason" binding:"required"`
}

type ConfirmReturnReq struct {
	OrderID     uint       `json:"order_id" binding:"required"`
	ReturnDate  *time.Time `json:"return_date"`
}

type SearchDroneReq struct {
	StartDate string  `form:"start_date" binding:"required"`
	EndDate   string  `form:"end_date" binding:"required"`
	Region    string  `form:"region"`
	Keyword   string  `form:"keyword"`
	MinPrice  float64 `form:"min_price"`
	MaxPrice  float64 `form:"max_price"`
	Page      int     `form:"page,default=1"`
	PageSize  int     `form:"page_size,default=10"`
}
