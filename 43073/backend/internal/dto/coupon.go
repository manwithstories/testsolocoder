package dto

type CouponCreateRequest struct {
	Type      string  `json:"type" validate:"required,oneof=fixed discount"`
	Value     float64 `json:"value" validate:"required,min=0"`
	MinAmount float64 `json:"minAmount" validate:"min=0"`
	TotalCount int     `json:"totalCount" validate:"required,min=1"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
}

type CouponListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Code     string `form:"code"`
	Status   string `form:"status"`
	Type     string `form:"type"`
}

type CouponUseRequest struct {
	Code string `json:"code" validate:"required"`
}
