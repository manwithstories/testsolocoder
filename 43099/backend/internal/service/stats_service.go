package service

import (
	"time"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"
)

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

func (s *StatsService) GetOverview() (*dto.StatsOverviewResponse, error) {
	var overview dto.StatsOverviewResponse

	database.DB.Model(&model.Order{}).Where("status IN ?", []string{
		string(model.OrderStatusConfirmed),
		string(model.OrderStatusPaid),
		string(model.OrderStatusCompleted),
	}).Count(&overview.TotalBookings)

	var totalRevenue float64
	database.DB.Model(&model.Order{}).Where("status IN ?", []string{
		string(model.OrderStatusPaid),
		string(model.OrderStatusCompleted),
	}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)
	overview.TotalRevenue = totalRevenue

	database.DB.Model(&model.User{}).Count(&overview.TotalUsers)
	database.DB.Model(&model.Venue{}).Where("status = ?", model.VenueStatusOnline).Count(&overview.TotalVenues)
	database.DB.Model(&model.Device{}).Where("status = ?", model.DeviceStatusOnline).Count(&overview.TotalDevices)

	database.DB.Model(&model.Order{}).Where("status = ?", model.OrderStatusPending).Count(&overview.PendingOrders)

	today := time.Now().Format("2006-01-02")
	database.DB.Model(&model.Order{}).Where("DATE(created_at) = ? AND status IN ?", today, []string{
		string(model.OrderStatusConfirmed),
		string(model.OrderStatusPaid),
		string(model.OrderStatusCompleted),
	}).Count(&overview.TodayBookings)

	var todayRevenue float64
	database.DB.Model(&model.Order{}).Where("DATE(created_at) = ? AND status IN ?", today, []string{
		string(model.OrderStatusPaid),
		string(model.OrderStatusCompleted),
	}).Select("COALESCE(SUM(total_amount), 0)").Scan(&todayRevenue)
	overview.TodayRevenue = todayRevenue

	return &overview, nil
}

func (s *StatsService) GetBookingStats(startDate, endDate string) ([]dto.BookingStatsResponse, error) {
	var results []struct {
		Date  string
		Count int64
	}

	err := database.DB.Model(&model.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).
		Where("status IN ?", []string{
			string(model.OrderStatusConfirmed),
			string(model.OrderStatusPaid),
			string(model.OrderStatusCompleted),
		}).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make([]dto.BookingStatsResponse, len(results))
	for i, r := range results {
		stats[i] = dto.BookingStatsResponse{
			Date:  r.Date,
			Count: r.Count,
		}
	}

	return stats, nil
}

func (s *StatsService) GetRevenueStats(startDate, endDate string) ([]dto.RevenueStatsResponse, error) {
	var results []struct {
		Date   string
		Amount float64
	}

	err := database.DB.Model(&model.Order{}).
		Select("DATE(created_at) as date, COALESCE(SUM(total_amount), 0) as amount").
		Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).
		Where("status IN ?", []string{
			string(model.OrderStatusPaid),
			string(model.OrderStatusCompleted),
		}).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make([]dto.RevenueStatsResponse, len(results))
	for i, r := range results {
		stats[i] = dto.RevenueStatsResponse{
			Date:   r.Date,
			Amount: r.Amount,
		}
	}

	return stats, nil
}

func (s *StatsService) GetPopularVenues(startDate, endDate string) ([]dto.PopularVenue, error) {
	var results []struct {
		ID       uint
		Name     string
		Bookings int64
		Revenue  float64
	}

	err := database.DB.Model(&model.Order{}).
		Select("item_id as id, item_name as name, COUNT(*) as bookings, COALESCE(SUM(total_amount), 0) as revenue").
		Where("type = ? AND DATE(created_at) BETWEEN ? AND ?", model.OrderTypeVenue, startDate, endDate).
		Where("status IN ?", []string{
			string(model.OrderStatusPaid),
			string(model.OrderStatusCompleted),
		}).
		Group("item_id, item_name").
		Order("bookings DESC").
		Limit(10).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	venues := make([]dto.PopularVenue, len(results))
	for i, r := range results {
		venues[i] = dto.PopularVenue{
			ID:       r.ID,
			Name:     r.Name,
			Bookings: r.Bookings,
			Revenue:  r.Revenue,
		}
	}

	return venues, nil
}
