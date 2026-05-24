package dto

type CreateClaimReq struct {
	OrderID       uint    `json:"order_id" binding:"required"`
	DamageDesc    string  `json:"damage_desc" binding:"required"`
	DamageImages  string  `json:"damage_images"`
	EstimatedCost float64 `json:"estimated_cost" binding:"required,gt=0"`
}

type ClaimQuery struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	Status    string `form:"status"`
	OrderID   uint   `form:"order_id"`
}

type ReviewClaimReq struct {
	ClaimID        uint    `json:"claim_id" binding:"required"`
	Status         string  `json:"status" binding:"required,oneof=approved rejected"`
	ActualCost     float64 `json:"actual_cost"`
	ReviewRemark   string  `json:"review_remark"`
	DeductedAmount float64 `json:"deducted_amount"`
}

type CreateReviewReq struct {
	Type      string `json:"type" binding:"required,oneof=rental service"`
	OrderID   *uint  `json:"order_id"`
	ServiceID *uint  `json:"service_id"`
	RevieweeID uint  `json:"reviewee_id" binding:"required"`
	DroneID   *uint  `json:"drone_id"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Content   string `json:"content" binding:"required"`
	Images    string `json:"images"`
}

type ReviewQuery struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	Type      string `form:"type"`
	OrderID   uint   `form:"order_id"`
	ServiceID uint   `form:"service_id"`
	RevieweeID uint  `form:"reviewee_id"`
	DroneID   uint   `form:"drone_id"`
}

type ReplyReviewReq struct {
	ReviewID uint   `json:"review_id" binding:"required"`
	Reply    string `json:"reply" binding:"required"`
}

type StatsQuery struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Region    string `form:"region"`
}

type DroneStats struct {
	DroneID     uint    `json:"drone_id"`
	DroneName   string  `json:"drone_name"`
	TotalDays   int     `json:"total_days"`
	Utilization float64 `json:"utilization"`
	Income      float64 `json:"income"`
}

type RevenueStats struct {
	Date    string  `json:"date"`
	Amount  float64 `json:"amount"`
	Count   int     `json:"count"`
}

type RegionStats struct {
	Region string  `json:"region"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
}
