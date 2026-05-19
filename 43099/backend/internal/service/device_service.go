package service

import (
	"errors"
	"fmt"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/internal/repository"
)

type DeviceService struct {
	deviceRepo *repository.DeviceRepository
	orderRepo  *repository.OrderRepository
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		deviceRepo: repository.NewDeviceRepository(),
		orderRepo:  repository.NewOrderRepository(),
	}
}

func (s *DeviceService) CreateCategory(req *dto.DeviceCategoryCreateRequest) (*model.DeviceCategory, error) {
	category := &model.DeviceCategory{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	err := s.deviceRepo.CreateCategory(category)
	return category, err
}

func (s *DeviceService) ListCategories() ([]model.DeviceCategory, error) {
	return s.deviceRepo.ListCategories()
}

func (s *DeviceService) Create(req *dto.DeviceCreateRequest) (*model.Device, error) {
	_, err := s.deviceRepo.GetCategoryByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	device := &model.Device{
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       req.Description,
		Specification:     req.Specification,
		StockQuantity:     req.StockQuantity,
		AvailableQuantity: req.StockQuantity,
		RentalPrice:       req.RentalPrice,
		DepositAmount:     req.DepositAmount,
		Status:            model.DeviceStatusOnline,
	}

	err = s.deviceRepo.Create(device)
	return device, err
}

func (s *DeviceService) GetByID(id uint) (*model.Device, error) {
	return s.deviceRepo.GetByID(id)
}

func (s *DeviceService) List(req *dto.DeviceListRequest) ([]model.Device, int64, error) {
	return s.deviceRepo.List(req)
}

func (s *DeviceService) Update(id uint, req *dto.DeviceUpdateRequest) (*model.Device, error) {
	device, err := s.deviceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.CategoryID > 0 {
		_, err := s.deviceRepo.GetCategoryByID(req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		device.CategoryID = req.CategoryID
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Description != "" {
		device.Description = req.Description
	}
	if req.Specification != "" {
		device.Specification = req.Specification
	}
	if req.StockQuantity >= 0 {
		diff := req.StockQuantity - device.StockQuantity
		device.StockQuantity = req.StockQuantity
		device.AvailableQuantity += diff
	}
	if req.RentalPrice >= 0 {
		device.RentalPrice = req.RentalPrice
	}
	if req.DepositAmount >= 0 {
		device.DepositAmount = req.DepositAmount
	}
	if req.Status != "" {
		device.Status = model.DeviceStatus(req.Status)
	}

	err = s.deviceRepo.Update(device)
	return device, err
}

func (s *DeviceService) UpdateStatus(id uint, status string) error {
	return s.deviceRepo.UpdateStatus(id, model.DeviceStatus(status))
}

func (s *DeviceService) GetAvailability(deviceID uint, date string) (*dto.DeviceAvailabilityResponse, error) {
	device, err := s.deviceRepo.GetByID(deviceID)
	if err != nil {
		return nil, errors.New("device not found")
	}

	bookedOrders, err := s.orderRepo.GetDeviceBookingsForDate(deviceID, date)
	if err != nil {
		return nil, err
	}

	bookedQuantity := 0
	for _, order := range bookedOrders {
		bookedQuantity += order.Quantity
	}

	availableQuantity := device.AvailableQuantity - bookedQuantity
	if availableQuantity < 0 {
		availableQuantity = 0
	}

	return &dto.DeviceAvailabilityResponse{
		Date:              date,
		DeviceID:          deviceID,
		TotalStock:        device.StockQuantity,
		AvailableQuantity: availableQuantity,
		BookedQuantity:    bookedQuantity,
	}, nil
}

func (s *DeviceService) BatchImport(items []dto.DeviceBatchImportItem) (*dto.DeviceBatchImportResponse, error) {
	response := &dto.DeviceBatchImportResponse{
		SuccessCount: 0,
		FailCount:    0,
		Errors:       make([]string, 0),
	}

	for i, item := range items {
		category, err := s.deviceRepo.GetCategoryByName(item.CategoryName)
		if err != nil {
			category, err = s.CreateCategory(&dto.DeviceCategoryCreateRequest{
				Name:        item.CategoryName,
				Description: "",
				SortOrder:   0,
			})
			if err != nil {
				response.FailCount++
				response.Errors = append(response.Errors, fmt.Sprintf("Row %d: Failed to create category - %v", i+1, err))
				continue
			}
		}

		_, err = s.Create(&dto.DeviceCreateRequest{
			CategoryID:    category.ID,
			Name:          item.Name,
			Description:   item.Description,
			Specification: item.Specification,
			StockQuantity: item.StockQuantity,
			RentalPrice:   item.RentalPrice,
			DepositAmount: item.DepositAmount,
		})

		if err != nil {
			response.FailCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Row %d: Failed to create device - %v", i+1, err))
			continue
		}

		response.SuccessCount++
	}

	return response, nil
}
