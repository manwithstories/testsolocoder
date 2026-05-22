package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"time"
)

type StatsService struct {
	orderRepo   *repository.OrderRepository
	bookingRepo *repository.BookingRepository
	carRepo     *repository.CarRepository
}

func NewStatsService() *StatsService {
	return &StatsService{
		orderRepo:   repository.NewOrderRepository(),
		bookingRepo: repository.NewBookingRepository(),
		carRepo:     repository.NewCarRepository(),
	}
}

type DashboardStats struct {
	MonthlyRevenue  float64                `json:"monthly_revenue"`
	MonthlyOrders   int64                  `json:"monthly_orders"`
	CarUtilization  float64                `json:"car_utilization"`
	TotalCars       int64                  `json:"total_cars"`
	ActiveBookings  int64                  `json:"active_bookings"`
	PopularCars     []repository.CarStat   `json:"popular_cars"`
	RevenueTrend    []RevenuePoint         `json:"revenue_trend"`
	StatusBreakdown []repository.StatusCount `json:"status_breakdown"`
}

type RevenuePoint struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
}

func (s *StatsService) GetDashboardStats() (*DashboardStats, error) {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	monthlyRevenue, err := s.orderRepo.GetTotalRevenue(&firstDayOfMonth, &now)
	if err != nil {
		return nil, err
	}

	_, monthlyOrders, err := s.orderRepo.FindAll(1, 1, 0, "", &firstDayOfMonth, &now)
	if err != nil {
		return nil, err
	}

	totalCars, err := s.carRepo.CountCars()
	if err != nil {
		return nil, err
	}

	rentedCars, err := s.carRepo.CountCarsByStatus(string(model.CarStatusRented))
	if err != nil {
		return nil, err
	}

	var carUtilization float64
	if totalCars > 0 {
		carUtilization = float64(rentedCars) / float64(totalCars) * 100
	}

	_, activeBookings, err := s.bookingRepo.FindAll(1, 1, 0, string(model.BookingStatusConfirmed), 0, nil, nil)
	if err != nil {
		return nil, err
	}

	popularCars, err := s.carRepo.GetPopularCars(10)
	if err != nil {
		return nil, err
	}

	revenueTrend, err := s.getRevenueTrend(30)
	if err != nil {
		return nil, err
	}

	statusBreakdown, err := s.carRepo.GetStatusBreakdown()
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		MonthlyRevenue:  monthlyRevenue,
		MonthlyOrders:   monthlyOrders,
		CarUtilization:  carUtilization,
		TotalCars:       totalCars,
		ActiveBookings:  activeBookings,
		PopularCars:     popularCars,
		RevenueTrend:    revenueTrend,
		StatusBreakdown: statusBreakdown,
	}, nil
}

func (s *StatsService) getRevenueTrend(days int) ([]RevenuePoint, error) {
	var trend []RevenuePoint
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		revenue, err := s.orderRepo.GetTotalRevenue(&startOfDay, &endOfDay)
		if err != nil {
			return nil, err
		}

		trend = append(trend, RevenuePoint{
			Date:    date.Format("01-02"),
			Revenue: revenue,
		})
	}

	return trend, nil
}

func (s *StatsService) GetRevenueStats(startDate, endDate *time.Time, groupBy ...string) (map[string]interface{}, error) {
	group := "day"
	if len(groupBy) > 0 && groupBy[0] != "" {
		group = groupBy[0]
	}
	totalRevenue, err := s.orderRepo.GetTotalRevenue(startDate, endDate)
	if err != nil {
		return nil, err
	}

	_, orderCount, err := s.orderRepo.FindAll(1, 1, 0, "", startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_revenue": totalRevenue,
		"order_count":   orderCount,
		"avg_revenue":   func() float64 { if orderCount > 0 { return totalRevenue / float64(orderCount) }; return 0 }(),
		"group_by":      group,
	}, nil
}

func (s *StatsService) GetBookingStats(startDate, endDate *time.Time) (map[string]interface{}, error) {
	statuses := []string{
		string(model.BookingStatusPending),
		string(model.BookingStatusConfirmed),
		string(model.BookingStatusCancelled),
		string(model.BookingStatusCompleted),
	}

	result := make(map[string]interface{})
	statusCounts := make(map[string]int64)

	for _, status := range statuses {
		_, count, err := s.bookingRepo.FindAll(1, 1, 0, status, 0, startDate, endDate)
		if err != nil {
			return nil, err
		}
		statusCounts[status] = count
	}

	_, totalCount, err := s.bookingRepo.FindAll(1, 1, 0, "", 0, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result["total"] = totalCount
	result["by_status"] = statusCounts

	return result, nil
}

func (s *StatsService) GetCarStats() (map[string]interface{}, error) {
	totalCars, err := s.carRepo.CountCars()
	if err != nil {
		return nil, err
	}

	statusBreakdown, err := s.carRepo.GetStatusBreakdown()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["total"] = totalCars
	result["by_status"] = statusBreakdown

	return result, nil
}

func (s *StatsService) GetUserStats(startDate, endDate *time.Time) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["total"] = 0
	result["new_in_period"] = 0
	return result, nil
}

func (s *StatsService) GetMaintenanceStats() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["total"] = 0
	result["pending"] = 0
	result["in_progress"] = 0
	result["completed"] = 0
	return result, nil
}

func (s *StatsService) GetReviewStats(limit int) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["total"] = 0
	result["average_rating"] = 0
	result["top_rated"] = []interface{}{}
	return result, nil
}