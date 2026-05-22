package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"errors"
	"time"
)

type PromoService struct {
	promoRepo *repository.PromoCodeRepository
}

func NewPromoService() *PromoService {
	return &PromoService{
		promoRepo: repository.NewPromoCodeRepository(),
	}
}

type CreatePromoRequest struct {
	Name        string    `json:"name" binding:"required"`
	Code        string    `json:"code" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Value       float64   `json:"value" binding:"required"`
	MinAmount   float64   `json:"min_amount"`
	MaxDiscount float64   `json:"max_discount"`
	UsageLimit  int       `json:"usage_limit" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

func (s *PromoService) CreatePromo(req *CreatePromoRequest) (*model.PromoCode, error) {
	if s.promoRepo.ExistsByCode(req.Code) {
		return nil, errors.New("优惠码已存在")
	}

	if !req.EndDate.After(req.StartDate) {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	promo := &model.PromoCode{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Value:       req.Value,
		MinAmount:   req.MinAmount,
		MaxDiscount: req.MaxDiscount,
		UsageLimit:  req.UsageLimit,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsActive:    true,
	}

	err := s.promoRepo.Create(promo)
	if err != nil {
		return nil, err
	}

	return promo, nil
}

func (s *PromoService) GetPromoByID(id uint) (*model.PromoCode, error) {
	return s.promoRepo.FindByID(id)
}

func (s *PromoService) GetPromoByCode(code string) (*model.PromoCode, error) {
	return s.promoRepo.FindByCode(code)
}

func (s *PromoService) GetAllPromos(page, pageSize int, keyword string) ([]model.PromoCode, int64, error) {
	return s.promoRepo.FindAll(page, pageSize, keyword, nil)
}

func (s *PromoService) UpdatePromo(id uint, updates map[string]interface{}) error {
	promo, err := s.promoRepo.FindByID(id)
	if err != nil {
		return errors.New("优惠码不存在")
	}

	if name, ok := updates["name"]; ok {
		promo.Name = name.(string)
	}
	if value, ok := updates["value"]; ok {
		promo.Value = value.(float64)
	}
	if minAmount, ok := updates["min_amount"]; ok {
		promo.MinAmount = minAmount.(float64)
	}
	if maxDiscount, ok := updates["max_discount"]; ok {
		promo.MaxDiscount = maxDiscount.(float64)
	}
	if usageLimit, ok := updates["usage_limit"]; ok {
		promo.UsageLimit = int(usageLimit.(float64))
	}
	if isActive, ok := updates["is_active"]; ok {
		promo.IsActive = isActive.(bool)
	}

	return s.promoRepo.Update(promo)
}

func (s *PromoService) DeletePromo(id uint) error {
	return s.promoRepo.Delete(id)
}
