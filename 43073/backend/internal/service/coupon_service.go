package service

import (
	"errors"
	"ticket-system/internal/dto"
	"ticket-system/internal/logger"
	"ticket-system/internal/models"
	"ticket-system/internal/util"
	"time"
)

type CouponService struct{}

func NewCouponService() *CouponService {
	return &CouponService{}
}

func (s *CouponService) Create(req *dto.CouponCreateRequest) (*models.Coupon, error) {
	var startTime, endTime time.Time
	var err error

	if req.StartTime != "" {
		startTime, err = time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
		if err != nil {
			return nil, errors.New("开始时间格式错误")
		}
	}

	if req.EndTime != "" {
		endTime, err = time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local)
		if err != nil {
			return nil, errors.New("结束时间格式错误")
		}
	}

	if !endTime.IsZero() && !startTime.IsZero() && endTime.Before(startTime) {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	coupon := &models.Coupon{
		Code:       util.GenerateCouponCode(),
		Type:       req.Type,
		Value:      req.Value,
		MinAmount:  req.MinAmount,
		TotalCount: req.TotalCount,
		UsedCount:  0,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     models.CouponStatusActive,
	}

	if err := models.DB.Create(coupon).Error; err != nil {
		logger.Log.Errorf("Create coupon failed: %v", err)
		return nil, errors.New("创建优惠券失败")
	}

	logger.Log.Infof("Coupon created: %s", coupon.Code)
	return coupon, nil
}

func (s *CouponService) GetList(req *dto.CouponListRequest) ([]models.Coupon, int64, error) {
	var coupons []models.Coupon
	var total int64

	query := models.DB.Model(&models.Coupon{})

	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&coupons).Error; err != nil {
		return nil, 0, err
	}

	return coupons, total, nil
}

func (s *CouponService) GetByID(id uint) (*models.Coupon, error) {
	var coupon models.Coupon
	if err := models.DB.First(&coupon, id).Error; err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (s *CouponService) GetByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	if err := models.DB.Where("code = ?", code).First(&coupon).Error; err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (s *CouponService) Validate(coupon *models.Coupon, amount float64) (float64, error) {
	if coupon.Status != models.CouponStatusActive {
		return 0, errors.New("优惠券不可用")
	}

	now := time.Now()
	if !coupon.StartTime.IsZero() && now.Before(coupon.StartTime) {
		return 0, errors.New("优惠券尚未生效")
	}

	if !coupon.EndTime.IsZero() && now.After(coupon.EndTime) {
		return 0, errors.New("优惠券已过期")
	}

	if coupon.UsedCount >= coupon.TotalCount {
		return 0, errors.New("优惠券已用完")
	}

	if amount < coupon.MinAmount {
		return 0, errors.New("订单金额不满足优惠券使用条件")
	}

	var discount float64
	if coupon.Type == models.CouponTypeFixed {
		discount = coupon.Value
	} else if coupon.Type == models.CouponTypeDiscount {
		discount = amount * (1 - coupon.Value/10)
	}

	if discount > amount {
		discount = amount
	}

	return discount, nil
}

func (s *CouponService) UseCoupon(id uint) error {
	coupon, err := s.GetByID(id)
	if err != nil {
		return err
	}

	coupon.UsedCount++
	if coupon.UsedCount >= coupon.TotalCount {
		coupon.Status = models.CouponStatusUsed
	}

	return models.DB.Save(coupon).Error
}

func (s *CouponService) Delete(id uint) error {
	if err := models.DB.Delete(&models.Coupon{}, id).Error; err != nil {
		return err
	}
	logger.Log.Infof("Coupon deleted: %d", id)
	return nil
}
