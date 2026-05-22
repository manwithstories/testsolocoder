package service

import (
	"errors"
	"ticket-system/internal/dto"
	"ticket-system/internal/logger"
	"ticket-system/internal/models"
	"time"
)

type CheckInService struct{}

func NewCheckInService() *CheckInService {
	return &CheckInService{}
}

func (s *CheckInService) CheckIn(qrCode string) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	if err := models.DB.Where("qr_code = ?", qrCode).First(&checkIn).Error; err != nil {
		return nil, errors.New("无效的签到码")
	}

	if checkIn.CheckedIn {
		return nil, errors.New("已签到，请勿重复签到")
	}

	var order models.Order
	if err := models.DB.First(&order, checkIn.OrderID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.Status != models.OrderStatusPaid {
		return nil, errors.New("订单未支付，无法签到")
	}

	now := time.Now()
	checkIn.CheckedIn = true
	checkIn.CheckedAt = &now

	if err := models.DB.Save(&checkIn).Error; err != nil {
		return nil, errors.New("签到失败")
	}

	logger.Log.Infof("Check-in success: order=%s, user=%d", order.OrderNo, checkIn.UserID)
	return &checkIn, nil
}

func (s *CheckInService) GetList(req *dto.CheckInListRequest) ([]models.CheckIn, int64, error) {
	var checkIns []models.CheckIn
	var total int64

	query := models.DB.Model(&models.CheckIn{}).Preload("Activity").Preload("Order")

	if req.ActivityID > 0 {
		query = query.Where("activity_id = ?", req.ActivityID)
	}

	if req.CheckedIn != nil {
		query = query.Where("checked_in = ?", *req.CheckedIn)
	}

	if req.Keyword != "" {
		query = query.Joins("JOIN orders ON orders.id = check_ins.order_id").
			Where("orders.order_no LIKE ?", "%"+req.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&checkIns).Error; err != nil {
		return nil, 0, err
	}

	return checkIns, total, nil
}

func (s *CheckInService) GetByOrderID(orderID uint) ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	if err := models.DB.Where("order_id = ?", orderID).Find(&checkIns).Error; err != nil {
		return nil, err
	}
	return checkIns, nil
}

func (s *CheckInService) GetStatistics(activityID uint) (int64, int64, error) {
	var total int64
	var checked int64

	query := models.DB.Model(&models.CheckIn{})
	if activityID > 0 {
		query = query.Where("activity_id = ?", activityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, 0, err
	}

	if err := query.Where("checked_in = ?", true).Count(&checked).Error; err != nil {
		return 0, 0, err
	}

	return total, checked, nil
}
