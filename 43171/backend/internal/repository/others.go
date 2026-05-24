package repository

import (
	"drone-rental/internal/config"
	"drone-rental/internal/model"
	"time"
)

type FlightRepo struct{}

func NewFlightRepo() *FlightRepo {
	return &FlightRepo{}
}

func (r *FlightRepo) Create(flight *model.FlightRecord) error {
	return config.DB.Create(flight).Error
}

func (r *FlightRepo) GetByID(id uint) (*model.FlightRecord, error) {
	var flight model.FlightRecord
	err := config.DB.Preload("Drone").Preload("Pilot").Preload("Order").Preload("Service").First(&flight, id).Error
	return &flight, err
}

func (r *FlightRepo) Update(flight *model.FlightRecord) error {
	return config.DB.Save(flight).Error
}

func (r *FlightRepo) Delete(id uint) error {
	return config.DB.Delete(&model.FlightRecord{}, id).Error
}

func (r *FlightRepo) List(page, pageSize int, droneID, pilotID, orderID, serviceID uint, startDate, endDate string) ([]model.FlightRecord, int64, error) {
	var flights []model.FlightRecord
	var total int64
	db := config.DB.Model(&model.FlightRecord{}).Preload("Drone").Preload("Pilot")
	if droneID != 0 {
		db = db.Where("drone_id = ?", droneID)
	}
	if pilotID != 0 {
		db = db.Where("pilot_id = ?", pilotID)
	}
	if orderID != 0 {
		db = db.Where("order_id = ?", orderID)
	}
	if serviceID != 0 {
		db = db.Where("service_id = ?", serviceID)
	}
	if startDate != "" {
		db = db.Where("flight_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("flight_date <= ?", endDate)
	}
	db.Count(&total)
	err := db.Order("flight_date DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&flights).Error
	return flights, total, err
}

type InsuranceRepo struct{}

func NewInsuranceRepo() *InsuranceRepo {
	return &InsuranceRepo{}
}

func (r *InsuranceRepo) Create(claim *model.InsuranceClaim) error {
	return config.DB.Create(claim).Error
}

func (r *InsuranceRepo) GetByID(id uint) (*model.InsuranceClaim, error) {
	var claim model.InsuranceClaim
	err := config.DB.Preload("Order").Preload("User").Preload("Reviewer").First(&claim, id).Error
	return &claim, err
}

func (r *InsuranceRepo) Update(claim *model.InsuranceClaim) error {
	return config.DB.Save(claim).Error
}

func (r *InsuranceRepo) List(page, pageSize int, orderID uint, status model.ClaimStatus) ([]model.InsuranceClaim, int64, error) {
	var claims []model.InsuranceClaim
	var total int64
	db := config.DB.Model(&model.InsuranceClaim{}).Preload("Order").Preload("User")
	if orderID != 0 {
		db = db.Where("order_id = ?", orderID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)
	err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&claims).Error
	return claims, total, err
}

type ReviewRepo struct{}

func NewReviewRepo() *ReviewRepo {
	return &ReviewRepo{}
}

func (r *ReviewRepo) Create(review *model.Review) error {
	return config.DB.Create(review).Error
}

func (r *ReviewRepo) GetByID(id uint) (*model.Review, error) {
	var review model.Review
	err := config.DB.Preload("Reviewer").Preload("Reviewee").First(&review, id).Error
	return &review, err
}

func (r *ReviewRepo) Update(review *model.Review) error {
	return config.DB.Save(review).Error
}

func (r *ReviewRepo) List(page, pageSize int, orderID, serviceID, revieweeID, droneID uint, reviewType model.ReviewType) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64
	db := config.DB.Model(&model.Review{}).Preload("Reviewer").Preload("Reviewee")
	if orderID != 0 {
		db = db.Where("order_id = ?", orderID)
	}
	if serviceID != 0 {
		db = db.Where("service_id = ?", serviceID)
	}
	if revieweeID != 0 {
		db = db.Where("reviewee_id = ?", revieweeID)
	}
	if droneID != 0 {
		db = db.Where("drone_id = ?", droneID)
	}
	if reviewType != "" {
		db = db.Where("type = ?", reviewType)
	}
	db.Count(&total)
	err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepo) GetAvgRatingByUser(userID uint) (float64, int64, error) {
	var avg float64
	var count int64
	db := config.DB.Model(&model.Review{}).Where("reviewee_id = ?", userID)
	db.Count(&count)
	if count == 0 {
		return 5.0, 0, nil
	}
	err := db.Select("COALESCE(AVG(rating), 5)").Scan(&avg).Error
	return avg, count, err
}

func (r *ReviewRepo) GetAvgRatingByDrone(droneID uint) (float64, int64, error) {
	var avg float64
	var count int64
	db := config.DB.Model(&model.Review{}).Where("drone_id = ?", droneID)
	db.Count(&count)
	if count == 0 {
		return 5.0, 0, nil
	}
	err := db.Select("COALESCE(AVG(rating), 5)").Scan(&avg).Error
	return avg, count, err
}

type StatsRepo struct{}

func NewStatsRepo() *StatsRepo {
	return &StatsRepo{}
}

func (r *StatsRepo) GetRevenueStats(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db := config.DB.Model(&model.RentalOrder{}).
		Select("DATE(created_at) as date, SUM(total_amount) as amount, COUNT(*) as count").
		Where("status IN ? AND created_at BETWEEN ? AND ?",
			[]model.OrderStatus{model.OrderStatusCompleted, model.OrderStatusReturned},
			startDate, endDate).
		Group("DATE(created_at)").
		Order("date ASC")
	err := db.Scan(&results).Error
	return results, err
}

func (r *StatsRepo) GetRegionStats(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db := config.DB.Model(&model.RentalOrder{}).
		Select("region, COUNT(*) as count, SUM(total_amount) as amount").
		Where("status IN ? AND created_at BETWEEN ? AND ?",
			[]model.OrderStatus{model.OrderStatusCompleted, model.OrderStatusReturned},
			startDate, endDate).
		Group("region").
		Order("count DESC")
	err := db.Scan(&results).Error
	return results, err
}

func (r *StatsRepo) GetDroneStats(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db := config.DB.Model(&model.RentalOrder{}).
		Select("drone_id, COUNT(*) as count, SUM(total_days) as total_days, SUM(rental_fee) as income").
		Where("status IN ? AND created_at BETWEEN ? AND ?",
			[]model.OrderStatus{model.OrderStatusCompleted, model.OrderStatusReturned},
			startDate, endDate).
		Group("drone_id").
		Order("income DESC")
	err := db.Scan(&results).Error
	return results, err
}
