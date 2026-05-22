package dto

type OccupancyRateRequest struct {
	StartDate string `form:"startDate" json:"startDate" binding:"required"`
	EndDate   string `form:"endDate" json:"endDate" binding:"required"`
}

type RevenueRequest struct {
	StartDate string `form:"startDate" json:"startDate" binding:"required"`
	EndDate   string `form:"endDate" json:"endDate" binding:"required"`
}

type OccupancyRateResponse struct {
	Date          string  `json:"date"`
	TotalRooms    int64   `json:"totalRooms"`
	OccupiedRooms int64   `json:"occupiedRooms"`
	OccupancyRate float64 `json:"occupancyRate"`
}

type RevenueResponse struct {
	Date           string  `json:"date"`
	RoomRevenue    float64 `json:"roomRevenue"`
	OtherRevenue   float64 `json:"otherRevenue"`
	TotalRevenue   float64 `json:"totalRevenue"`
	CheckOutCount  int64   `json:"checkOutCount"`
}

type ReportExportRequest struct {
	StartDate string `form:"startDate" json:"startDate" binding:"required"`
	EndDate   string `form:"endDate" json:"endDate" binding:"required"`
}
