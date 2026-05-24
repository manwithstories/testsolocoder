package dto

// StatisticsRequest 通用统计请求
type StatisticsRequest struct {
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
	Granularity string `form:"granularity" binding:"omitempty,oneof=day week month"`
}

// SalesTrendItem 销售趋势单项
type SalesTrendItem struct {
	Period     string  `json:"period"`
	Amount     float64 `json:"amount"`
	OrderCount int64   `json:"order_count"`
}

// SalesTrendResponse 销售趋势响应
type SalesTrendResponse struct {
	Granularity string           `json:"granularity"`
	StartDate   string           `json:"start_date"`
	EndDate     string           `json:"end_date"`
	List        []SalesTrendItem `json:"list"`
	TotalAmount float64          `json:"total_amount"`
	TotalOrders int64            `json:"total_orders"`
}

// RoleDistribution 角色分布
type RoleDistribution struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

// RegionDistribution 地区分布
type RegionDistribution struct {
	Region string `json:"region"`
	Count  int64  `json:"count"`
}

// CustomerProfileResponse 客户画像响应
type CustomerProfileResponse struct {
	TotalUsers         int64                `json:"total_users"`
	RoleDistribution   []RoleDistribution   `json:"role_distribution"`
	RegionDistribution []RegionDistribution `json:"region_distribution"`
}

// ExportRequest 导出请求
type ExportRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Type      string `form:"type" binding:"required,oneof=sales orders users"`
}
