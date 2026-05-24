package dto

import "furniture-platform/internal/model"

type UpdateProductionStatusRequest struct {
	Status model.ProductionStatus `json:"status" binding:"required"`
	Remark string                 `json:"remark"`
}

type ProductionListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	OrderID  int64  `form:"order_id"`
	Status   string `form:"status"`
}

type ProductionResponse struct {
	ID              int64                     `json:"id"`
	OrderID         int64                     `json:"order_id"`
	OrderNo         string                    `json:"order_no"`
	Status          model.ProductionStatus    `json:"status"`
	ProgressPercent int                       `json:"progress_percent"`
	OperatorID      int64                     `json:"operator_id"`
	Remark          string                    `json:"remark"`
	CreatedAt       string                    `json:"created_at"`
	UpdatedAt       string                    `json:"updated_at"`
}

type ProductionListResponse struct {
	List  []ProductionResponse `json:"list"`
	Total int64                `json:"total"`
}
