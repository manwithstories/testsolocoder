package service

import (
	"encoding/json"
	"fmt"
	"time"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
	"beauty-salon-system/internal/repository/redis"
	"beauty-salon-system/internal/utils"
)

type ServiceItemService struct {
	serviceRepo *repository.ServiceRepository
}

func NewServiceItemService(serviceRepo *repository.ServiceRepository) *ServiceItemService {
	return &ServiceItemService{
		serviceRepo: serviceRepo,
	}
}

type CreateServiceRequest struct {
	Name           string  `json:"name" binding:"required"`
	Category       string  `json:"category" binding:"required"`
	Description    string  `json:"description"`
	Price          float64 `json:"price" binding:"required"`
	Duration       int     `json:"duration" binding:"required"`
	RequiredSkill  string  `json:"required_skill"`
	Products       string  `json:"products"`
	IsPackage      bool    `json:"is_package"`
	PackageCount   int     `json:"package_count"`
	DynamicPricing bool    `json:"dynamic_pricing"`
	WeekendPrice   float64 `json:"weekend_price"`
	HolidayPrice   float64 `json:"holiday_price"`
}

type UpdateServiceRequest struct {
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Duration       int     `json:"duration"`
	RequiredSkill  string  `json:"required_skill"`
	Products       string  `json:"products"`
	IsPackage      bool    `json:"is_package"`
	PackageCount   int     `json:"package_count"`
	DynamicPricing bool    `json:"dynamic_pricing"`
	WeekendPrice   float64 `json:"weekend_price"`
	HolidayPrice   float64 `json:"holiday_price"`
	Status         int     `json:"status"`
}

type AddPackageServiceRequest struct {
	ServiceID      uint `json:"service_id" binding:"required"`
	ChildServiceID uint `json:"child_service_id" binding:"required"`
	Count          int  `json:"count"`
}

func (s *ServiceItemService) Create(req *CreateServiceRequest) (*model.Service, error) {
	cacheKey := "services:all"
	redis.Delete(cacheKey)

	service := &model.Service{
		Name:           req.Name,
		Category:       req.Category,
		Description:    req.Description,
		Price:          req.Price,
		Duration:       req.Duration,
		RequiredSkill:  req.RequiredSkill,
		Products:       req.Products,
		IsPackage:      req.IsPackage,
		PackageCount:   req.PackageCount,
		DynamicPricing: req.DynamicPricing,
		WeekendPrice:   req.WeekendPrice,
		HolidayPrice:   req.HolidayPrice,
		Status:         1,
	}

	if err := s.serviceRepo.Create(service); err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}

	return service, nil
}

func (s *ServiceItemService) GetByID(id uint) (*model.Service, error) {
	cacheKey := fmt.Sprintf("service:%d", id)
	cached, err := redis.Get(cacheKey)
	if err == nil && cached != "" {
		var service model.Service
		json.Unmarshal([]byte(cached), &service)
		return &service, nil
	}

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(service)
	redis.Set(cacheKey, string(data), 10*time.Minute)

	return service, nil
}

func (s *ServiceItemService) Update(id uint, req *UpdateServiceRequest) (*model.Service, error) {
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	if req.Name != "" {
		service.Name = req.Name
	}
	if req.Category != "" {
		service.Category = req.Category
	}
	if req.Description != "" {
		service.Description = req.Description
	}
	if req.Price > 0 {
		service.Price = req.Price
	}
	if req.Duration > 0 {
		service.Duration = req.Duration
	}
	if req.RequiredSkill != "" {
		service.RequiredSkill = req.RequiredSkill
	}
	if req.Products != "" {
		service.Products = req.Products
	}
	service.IsPackage = req.IsPackage
	if req.PackageCount > 0 {
		service.PackageCount = req.PackageCount
	}
	service.DynamicPricing = req.DynamicPricing
	if req.WeekendPrice > 0 {
		service.WeekendPrice = req.WeekendPrice
	}
	if req.HolidayPrice > 0 {
		service.HolidayPrice = req.HolidayPrice
	}
	if req.Status > 0 {
		service.Status = req.Status
	}

	if err := s.serviceRepo.Update(service); err != nil {
		return nil, fmt.Errorf("update service: %w", err)
	}

	redis.Delete(fmt.Sprintf("service:%d", id))
	redis.Delete("services:all")

	return service, nil
}

func (s *ServiceItemService) Delete(id uint) error {
	redis.Delete(fmt.Sprintf("service:%d", id))
	redis.Delete("services:all")
	return s.serviceRepo.Delete(id)
}

func (s *ServiceItemService) List(page, pageSize int, category string, isPackage bool) ([]model.Service, int64, error) {
	return s.serviceRepo.List(page, pageSize, category, isPackage)
}

func (s *ServiceItemService) ListAll() ([]model.Service, error) {
	cacheKey := "services:all"
	cached, err := redis.Get(cacheKey)
	if err == nil && cached != "" {
		var services []model.Service
		json.Unmarshal([]byte(cached), &services)
		return services, nil
	}

	services, err := s.serviceRepo.ListAll()
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(services)
	redis.Set(cacheKey, string(data), 10*time.Minute)

	return services, nil
}

func (s *ServiceItemService) AddPackageService(req *AddPackageServiceRequest) error {
	_, err := s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}

	_, err = s.serviceRepo.GetByID(req.ChildServiceID)
	if err != nil {
		return fmt.Errorf("child service not found: %w", err)
	}

	count := req.Count
	if count <= 0 {
		count = 1
	}

	return s.serviceRepo.AddPackageService(&model.PackageService{
		ServiceID:      req.ServiceID,
		ChildServiceID: req.ChildServiceID,
		Count:          count,
	})
}

func (s *ServiceItemService) GetPackageServices(serviceID uint) ([]model.PackageService, error) {
	return s.serviceRepo.GetPackageServices(serviceID)
}

func (s *ServiceItemService) DeletePackageServices(serviceID uint) error {
	return s.serviceRepo.DeletePackageServices(serviceID)
}

func (s *ServiceItemService) CalculatePrice(serviceID uint, date time.Time) (float64, error) {
	service, err := s.GetByID(serviceID)
	if err != nil {
		return 0, err
	}

	if !service.DynamicPricing {
		return service.Price, nil
	}

	if utils.IsWeekend(date) && service.WeekendPrice > 0 {
		return service.WeekendPrice, nil
	}

	return service.Price, nil
}
