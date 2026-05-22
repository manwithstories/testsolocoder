package dto

type StatisticsRequest struct {
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
	ActivityID uint  `form:"activityId"`
}

type ActivityStatistics struct {
	ActivityID    uint    `json:"activityId"`
	ActivityTitle string  `json:"activityTitle"`
	TotalOrders   int64   `json:"totalOrders"`
	TotalAmount    float64 `json:"totalAmount"`
	TotalTickets  int64   `json:"totalTickets"`
}

type TicketTypeStatistics struct {
	TicketTypeID   uint    `json:"ticketTypeId"`
	TicketTypeName string  `json:"ticketTypeName"`
	TotalSold      int64   `json:"totalSold"`
	TotalAmount    float64 `json:"totalAmount"`
}

type DailyStatistics struct {
	Date         string  `json:"date"`
	TotalOrders  int64   `json:"totalOrders"`
	TotalAmount  float64 `json:"totalAmount"`
	TotalTickets int64   `json:"totalTickets"`
}
