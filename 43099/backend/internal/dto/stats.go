package dto

type StatsOverviewResponse struct {
	TotalBookings    int64   `json:"total_bookings"`
	TotalRevenue     float64 `json:"total_revenue"`
	TotalUsers       int64   `json:"total_users"`
	TotalVenues      int64   `json:"total_venues"`
	TotalDevices     int64   `json:"total_devices"`
	PendingOrders    int64   `json:"pending_orders"`
	TodayBookings    int64   `json:"today_bookings"`
	TodayRevenue     float64 `json:"today_revenue"`
}

type BookingStatsResponse struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type RevenueStatsResponse struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type PopularVenue struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Bookings int64 `json:"bookings"`
	Revenue float64 `json:"revenue"`
}
