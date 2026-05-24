package service

import (
	"drone-rental/internal/config"
	"drone-rental/internal/dto"
	"drone-rental/internal/model"
	"drone-rental/internal/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

type FlightService struct {
	flightRepo *repository.FlightRepo
	droneRepo  *repository.DroneRepo
}

func NewFlightService() *FlightService {
	return &FlightService{
		flightRepo: repository.NewFlightRepo(),
		droneRepo:  repository.NewDroneRepo(),
	}
}

func (s *FlightService) Create(pilotID uint, req *dto.CreateFlightRecordReq) (*model.FlightRecord, error) {
	drone, err := s.droneRepo.GetByID(req.DroneID)
	if err != nil {
		return nil, errors.New("设备不存在")
	}
	flightDate, err := time.ParseInLocation("2006-01-02", req.FlightDate, time.Local)
	if err != nil {
		flightDate = time.Now()
	}
	flight := &model.FlightRecord{
		OrderID:     req.OrderID,
		ServiceID:   req.ServiceID,
		DroneID:     req.DroneID,
		PilotID:     pilotID,
		StartPoint:  req.StartPoint,
		EndPoint:    req.EndPoint,
		Route:       req.Route,
		AltitudeMax: req.AltitudeMax,
		AltitudeAvg: req.AltitudeAvg,
		Duration:    req.Duration,
		Distance:    req.Distance,
		FlightDate:  flightDate,
		FlightLog:   req.FlightLog,
		Images:      req.Images,
		Remark:      req.Remark,
	}
	if err := s.flightRepo.Create(flight); err != nil {
		return nil, err
	}
	_ = drone
	return flight, nil
}

func (s *FlightService) GetByID(id uint) (*model.FlightRecord, error) {
	return s.flightRepo.GetByID(id)
}

func (s *FlightService) Update(id, pilotID uint, req *dto.CreateFlightRecordReq) error {
	flight, err := s.flightRepo.GetByID(id)
	if err != nil {
		return errors.New("飞行记录不存在")
	}
	if flight.PilotID != pilotID {
		return errors.New("无权修改该记录")
	}
	if req.StartPoint != "" {
		flight.StartPoint = req.StartPoint
	}
	if req.EndPoint != "" {
		flight.EndPoint = req.EndPoint
	}
	if req.Route != "" {
		flight.Route = req.Route
	}
	if req.AltitudeMax > 0 {
		flight.AltitudeMax = req.AltitudeMax
	}
	if req.AltitudeAvg > 0 {
		flight.AltitudeAvg = req.AltitudeAvg
	}
	if req.Duration > 0 {
		flight.Duration = req.Duration
	}
	if req.Distance > 0 {
		flight.Distance = req.Distance
	}
	if req.FlightLog != "" {
		flight.FlightLog = req.FlightLog
	}
	if req.Images != "" {
		flight.Images = req.Images
	}
	if req.Remark != "" {
		flight.Remark = req.Remark
	}
	return s.flightRepo.Update(flight)
}

func (s *FlightService) Delete(id, pilotID uint) error {
	flight, err := s.flightRepo.GetByID(id)
	if err != nil {
		return errors.New("飞行记录不存在")
	}
	if flight.PilotID != pilotID {
		return errors.New("无权删除该记录")
	}
	return s.flightRepo.Delete(id)
}

func (s *FlightService) List(page, pageSize int, droneID, pilotID, orderID, serviceID uint, startDate, endDate string) ([]model.FlightRecord, int64, error) {
	return s.flightRepo.List(page, pageSize, droneID, pilotID, orderID, serviceID, startDate, endDate)
}

type InsuranceService struct {
	claimRepo *repository.InsuranceRepo
	orderRepo *repository.OrderRepo
	userRepo  *repository.UserRepo
}

func NewInsuranceService() *InsuranceService {
	return &InsuranceService{
		claimRepo: repository.NewInsuranceRepo(),
		orderRepo: repository.NewOrderRepo(),
		userRepo:  repository.NewUserRepo(),
	}
}

func (s *InsuranceService) CreateClaim(userID uint, req *dto.CreateClaimReq) (*model.InsuranceClaim, error) {
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if order.UserID != userID {
		return nil, errors.New("无权操作该订单")
	}
	claim := &model.InsuranceClaim{
		ClaimNo:       "CL" + time.Now().Format("20060102150405"),
		OrderID:       req.OrderID,
		UserID:        userID,
		DamageDesc:    req.DamageDesc,
		DamageImages:  req.DamageImages,
		EstimatedCost: req.EstimatedCost,
		Status:        model.ClaimStatusPending,
	}
	if err := s.claimRepo.Create(claim); err != nil {
		return nil, err
	}
	return claim, nil
}

func (s *InsuranceService) GetClaimByID(id uint) (*model.InsuranceClaim, error) {
	return s.claimRepo.GetByID(id)
}

func (s *InsuranceService) ListClaims(page, pageSize int, orderID uint, status model.ClaimStatus) ([]model.InsuranceClaim, int64, error) {
	return s.claimRepo.List(page, pageSize, orderID, status)
}

func (s *InsuranceService) ReviewClaim(reviewerID uint, req *dto.ReviewClaimReq) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		claim, err := s.claimRepo.GetByID(req.ClaimID)
		if err != nil {
			return errors.New("理赔申请不存在")
		}
		if claim.Status != model.ClaimStatusPending && claim.Status != model.ClaimStatusReviewing {
			return errors.New("该申请已处理")
		}
		now := time.Now()
		claim.ReviewerID = &reviewerID
		claim.ReviewRemark = req.ReviewRemark
		claim.ReviewedAt = &now
		if req.ActualCost > 0 {
			claim.ActualCost = req.ActualCost
		}
		if req.Status == "approved" {
			claim.Status = model.ClaimStatusApproved
			claim.DeductedAmount = req.DeductedAmount
			user, _ := s.userRepo.GetByID(claim.UserID)
			user.Deposit -= req.DeductedAmount
			if err := tx.Save(user).Error; err != nil {
				return err
			}
		} else {
			claim.Status = model.ClaimStatusRejected
		}
		return tx.Save(claim).Error
	})
}

type ReviewService struct {
	reviewRepo *repository.ReviewRepo
	userRepo   *repository.UserRepo
	droneRepo  *repository.DroneRepo
}

func NewReviewService() *ReviewService {
	return &ReviewService{
		reviewRepo: repository.NewReviewRepo(),
		userRepo:   repository.NewUserRepo(),
		droneRepo:  repository.NewDroneRepo(),
	}
}

func (s *ReviewService) Create(reviewerID uint, req *dto.CreateReviewReq) (*model.Review, error) {
	review := &model.Review{
		Type:       model.ReviewType(req.Type),
		OrderID:    req.OrderID,
		ServiceID:  req.ServiceID,
		ReviewerID: reviewerID,
		RevieweeID: req.RevieweeID,
		DroneID:    req.DroneID,
		Rating:     req.Rating,
		Content:    req.Content,
		Images:     req.Images,
	}
	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}
	go s.updateUserRating(req.RevieweeID)
	if req.DroneID != nil && *req.DroneID > 0 {
		go s.updateDroneRating(*req.DroneID)
	}
	return review, nil
}

func (s *ReviewService) GetByID(id uint) (*model.Review, error) {
	return s.reviewRepo.GetByID(id)
}

func (s *ReviewService) List(page, pageSize int, orderID, serviceID, revieweeID, droneID uint, reviewType model.ReviewType) ([]model.Review, int64, error) {
	return s.reviewRepo.List(page, pageSize, orderID, serviceID, revieweeID, droneID, reviewType)
}

func (s *ReviewService) Reply(reviewerID uint, req *dto.ReplyReviewReq) error {
	review, err := s.reviewRepo.GetByID(req.ReviewID)
	if err != nil {
		return errors.New("评价不存在")
	}
	if review.RevieweeID != reviewerID {
		return errors.New("无权回复该评价")
	}
	review.Reply = req.Reply
	return s.reviewRepo.Update(review)
}

func (s *ReviewService) updateUserRating(userID uint) {
	avg, count, err := s.reviewRepo.GetAvgRatingByUser(userID)
	if err != nil {
		return
	}
	config.DB.Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{"rating": avg, "rating_count": count})
}

func (s *ReviewService) updateDroneRating(droneID uint) {
	avg, count, err := s.reviewRepo.GetAvgRatingByDrone(droneID)
	if err != nil {
		return
	}
	config.DB.Model(&model.Drone{}).Where("id = ?", droneID).
		Updates(map[string]interface{}{"rating": avg, "rating_count": count})
}

type StatsService struct {
	statsRepo *repository.StatsRepo
	droneRepo *repository.DroneRepo
}

func NewStatsService() *StatsService {
	return &StatsService{
		statsRepo: repository.NewStatsRepo(),
		droneRepo: repository.NewDroneRepo(),
	}
}

func (s *StatsService) GetRevenueStats(startDate, endDate string) ([]dto.RevenueStats, error) {
	start, _ := time.ParseInLocation("2006-01-02", startDate, time.Local)
	end, _ := time.ParseInLocation("2006-01-02", endDate, time.Local)
	results, err := s.statsRepo.GetRevenueStats(start, end)
	if err != nil {
		return nil, err
	}
	var stats []dto.RevenueStats
	for _, r := range results {
		var item dto.RevenueStats
		if date, ok := r["date"].([]uint8); ok {
			item.Date = string(date)
		}
		if amount, ok := r["amount"].(float64); ok {
			item.Amount = amount
		}
		if count, ok := r["count"].(int64); ok {
			item.Count = int(count)
		}
		stats = append(stats, item)
	}
	return stats, nil
}

func (s *StatsService) GetRegionStats(startDate, endDate string) ([]dto.RegionStats, error) {
	start, _ := time.ParseInLocation("2006-01-02", startDate, time.Local)
	end, _ := time.ParseInLocation("2006-01-02", endDate, time.Local)
	results, err := s.statsRepo.GetRegionStats(start, end)
	if err != nil {
		return nil, err
	}
	var stats []dto.RegionStats
	for _, r := range results {
		var item dto.RegionStats
		if region, ok := r["region"].(string); ok {
			item.Region = region
		}
		if count, ok := r["count"].(int64); ok {
			item.Count = int(count)
		}
		if amount, ok := r["amount"].(float64); ok {
			item.Amount = amount
		}
		stats = append(stats, item)
	}
	return stats, nil
}

func (s *StatsService) GetDroneStats(startDate, endDate string) ([]dto.DroneStats, error) {
	start, _ := time.ParseInLocation("2006-01-02", startDate, time.Local)
	end, _ := time.ParseInLocation("2006-01-02", endDate, time.Local)
	results, err := s.statsRepo.GetDroneStats(start, end)
	if err != nil {
		return nil, err
	}
	var stats []dto.DroneStats
	for _, r := range results {
		var item dto.DroneStats
		if droneID, ok := r["drone_id"].(uint); ok {
			item.DroneID = droneID
		}
		if totalDays, ok := r["total_days"].(int64); ok {
			item.TotalDays = int(totalDays)
		}
		if income, ok := r["income"].(float64); ok {
			item.Income = income
		}
		drone, err := s.droneRepo.GetByID(item.DroneID)
		if err == nil {
			item.DroneName = drone.Name
		}
		stats = append(stats, item)
	}
	return stats, nil
}
