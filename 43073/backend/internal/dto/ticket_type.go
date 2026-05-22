package dto

type TicketTypeCreateRequest struct {
	ActivityID uint    `json:"activityId" validate:"required"`
	Name       string  `json:"name" validate:"required,max=100"`
	Type       string  `json:"type" validate:"required,oneof=normal vip early_bird"`
	Price      float64 `json:"price" validate:"required,min=0"`
	Stock      int     `json:"stock" validate:"required,min=0"`
}

type TicketTypeUpdateRequest struct {
	Name   string  `json:"name" validate:"omitempty,max=100"`
	Type   string  `json:"type" validate:"omitempty,oneof=normal vip early_bird"`
	Price  float64 `json:"price" validate:"omitempty,min=0"`
	Stock  int     `json:"stock" validate:"omitempty,min=0"`
	Status string  `json:"status" validate:"omitempty,oneof=on_sale sold_out off_sale"`
}

type TicketTypeListRequest struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"pageSize,default=10"`
	ActivityID uint   `form:"activityId"`
	Status     string `form:"status"`
}
