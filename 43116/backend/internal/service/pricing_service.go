package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"errors"
	"time"
)

type PricingService struct {
	pricingRepo *repository.PricingRuleRepository
}

func NewPricingService() *PricingService {
	return &PricingService{
		pricingRepo: repository.NewPricingRuleRepository(),
	}
}

type CreatePricingRuleRequest struct {
	Name       string     `json:"name" binding:"required"`
	RuleType   string     `json:"rule_type" binding:"required"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Weekdays   string     `json:"weekdays"`
	Multiplier float64    `json:"multiplier" binding:"required"`
	CarModel   string     `json:"car_model"`
	MinDays    int        `json:"min_days"`
	MaxDays    int        `json:"max_days"`
	Priority   int        `json:"priority"`
}

func (s *PricingService) CreateRule(req *CreatePricingRuleRequest) (*model.PricingRule, error) {
	if req.RuleType != "weekend" && req.RuleType != "holiday" && req.RuleType != "special" {
		return nil, errors.New("无效的规则类型")
	}

	if req.Multiplier <= 0 {
		return nil, errors.New("倍率必须大于0")
	}

	if req.RuleType == "holiday" {
		if req.StartDate == nil || req.EndDate == nil {
			return nil, errors.New("节假日规则必须设置开始和结束日期")
		}
		if !req.EndDate.After(*req.StartDate) {
			return nil, errors.New("结束日期必须晚于开始日期")
		}
	}

	rule := &model.PricingRule{
		Name:       req.Name,
		RuleType:   req.RuleType,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Weekdays:   req.Weekdays,
		Multiplier: req.Multiplier,
		CarModel:   req.CarModel,
		MinDays:    req.MinDays,
		MaxDays:    req.MaxDays,
		Priority:   req.Priority,
		IsActive:   true,
	}

	err := s.pricingRepo.Create(rule)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *PricingService) GetRuleByID(id uint) (*model.PricingRule, error) {
	return s.pricingRepo.FindByID(id)
}

func (s *PricingService) GetAllRules(page, pageSize int, ruleType string) ([]model.PricingRule, int64, error) {
	return s.pricingRepo.FindAll(page, pageSize, ruleType)
}

func (s *PricingService) UpdateRule(id uint, updates map[string]interface{}) error {
	rule, err := s.pricingRepo.FindByID(id)
	if err != nil {
		return errors.New("规则不存在")
	}

	if name, ok := updates["name"]; ok {
		rule.Name = name.(string)
	}
	if multiplier, ok := updates["multiplier"]; ok {
		rule.Multiplier = multiplier.(float64)
	}
	if priority, ok := updates["priority"]; ok {
		rule.Priority = int(priority.(float64))
	}
	if weekdays, ok := updates["weekdays"]; ok {
		rule.Weekdays = weekdays.(string)
	}
	if carModel, ok := updates["car_model"]; ok {
		rule.CarModel = carModel.(string)
	}

	return s.pricingRepo.Update(rule)
}

func (s *PricingService) DeleteRule(id uint) error {
	return s.pricingRepo.Delete(id)
}

func (s *PricingService) ToggleActive(id uint) error {
	return s.pricingRepo.ToggleActive(id)
}
