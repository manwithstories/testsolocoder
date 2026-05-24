package service

import (
	"drone-rental/internal/dto"
	"drone-rental/internal/model"
	"drone-rental/internal/repository"
	"errors"
)

type DroneService struct {
	droneRepo *repository.DroneRepo
}

func NewDroneService() *DroneService {
	return &DroneService{
		droneRepo: repository.NewDroneRepo(),
	}
}

func (s *DroneService) Create(ownerID uint, req *dto.CreateDroneReq) (*model.Drone, error) {
	existing, _ := s.droneRepo.GetBySerialNo(req.SerialNo)
	if existing != nil {
		return nil, errors.New("该序列号设备已存在")
	}
	drone := &model.Drone{
		OwnerID:      ownerID,
		Name:         req.Name,
		Brand:        req.Brand,
		Model:        req.Model,
		SerialNo:     req.SerialNo,
		Weight:       req.Weight,
		BatteryLife:  req.BatteryLife,
		GimbalSpec:   req.GimbalSpec,
		CameraSpec:   req.CameraSpec,
		MaxSpeed:     req.MaxSpeed,
		MaxAltitude:  req.MaxAltitude,
		Region:       req.Region,
		Description:  req.Description,
		Images:       req.Images,
		PricePerDay:  req.PricePerDay,
		Deposit:      req.Deposit,
		Status:       model.DroneStatusOffline,
		Rating:       5.0,
	}
	if req.AvailableFrom != "" {
		t, err := parseDate(req.AvailableFrom)
		if err == nil {
			drone.AvailableFrom = &t
		}
	}
	if req.AvailableTo != "" {
		t, err := parseDate(req.AvailableTo)
		if err == nil {
			drone.AvailableTo = &t
		}
	}
	if err := s.droneRepo.Create(drone); err != nil {
		return nil, err
	}
	return drone, nil
}

func (s *DroneService) GetByID(id uint) (*model.Drone, error) {
	return s.droneRepo.GetByID(id)
}

func (s *DroneService) Update(id, ownerID uint, req *dto.UpdateDroneReq) error {
	drone, err := s.droneRepo.GetByID(id)
	if err != nil {
		return err
	}
	if drone.OwnerID != ownerID {
		return errors.New("无权修改该设备")
	}
	if req.Name != "" {
		drone.Name = req.Name
	}
	if req.Brand != "" {
		drone.Brand = req.Brand
	}
	if req.Model != "" {
		drone.Model = req.Model
	}
	if req.Weight > 0 {
		drone.Weight = req.Weight
	}
	if req.BatteryLife > 0 {
		drone.BatteryLife = req.BatteryLife
	}
	if req.GimbalSpec != "" {
		drone.GimbalSpec = req.GimbalSpec
	}
	if req.CameraSpec != "" {
		drone.CameraSpec = req.CameraSpec
	}
	if req.MaxSpeed > 0 {
		drone.MaxSpeed = req.MaxSpeed
	}
	if req.MaxAltitude > 0 {
		drone.MaxAltitude = req.MaxAltitude
	}
	if req.Region != "" {
		drone.Region = req.Region
	}
	if req.Description != "" {
		drone.Description = req.Description
	}
	if req.Images != "" {
		drone.Images = req.Images
	}
	if req.PricePerDay > 0 {
		drone.PricePerDay = req.PricePerDay
	}
	if req.Deposit > 0 {
		drone.Deposit = req.Deposit
	}
	if req.AvailableFrom != "" {
		t, err := parseDate(req.AvailableFrom)
		if err == nil {
			drone.AvailableFrom = &t
		}
	}
	if req.AvailableTo != "" {
		t, err := parseDate(req.AvailableTo)
		if err == nil {
			drone.AvailableTo = &t
		}
	}
	return s.droneRepo.Update(drone)
}

func (s *DroneService) UpdateStatus(id, ownerID uint, status model.DroneStatus) error {
	drone, err := s.droneRepo.GetByID(id)
	if err != nil {
		return err
	}
	if drone.OwnerID != ownerID {
		return errors.New("无权修改该设备")
	}
	return s.droneRepo.UpdateStatus(id, status)
}

func (s *DroneService) List(page, pageSize int, ownerID uint, keyword, region, brand string, status model.DroneStatus) ([]model.Drone, int64, error) {
	return s.droneRepo.List(page, pageSize, ownerID, keyword, region, brand, status)
}

func (s *DroneService) SearchAvailable(startDate, endDate, region, keyword string, minPrice, maxPrice float64, page, pageSize int) ([]model.Drone, int64, error) {
	start, err := parseDate(startDate)
	if err != nil {
		return nil, 0, err
	}
	end, err := parseDate(endDate)
	if err != nil {
		return nil, 0, err
	}
	return s.droneRepo.SearchAvailable(start, end, region, keyword, minPrice, maxPrice, page, pageSize)
}

func (s *DroneService) BatchImport(ownerID uint, drones []dto.CreateDroneReq) (int, error) {
	var droneModels []model.Drone
	for _, req := range drones {
		existing, _ := s.droneRepo.GetBySerialNo(req.SerialNo)
		if existing != nil {
			continue
		}
		droneModels = append(droneModels, model.Drone{
			OwnerID:     ownerID,
			Name:        req.Name,
			Brand:       req.Brand,
			Model:       req.Model,
			SerialNo:    req.SerialNo,
			Weight:      req.Weight,
			BatteryLife: req.BatteryLife,
			GimbalSpec:  req.GimbalSpec,
			CameraSpec:  req.CameraSpec,
			MaxSpeed:    req.MaxSpeed,
			MaxAltitude: req.MaxAltitude,
			Region:      req.Region,
			Description: req.Description,
			Images:      req.Images,
			PricePerDay: req.PricePerDay,
			Deposit:     req.Deposit,
			Status:      model.DroneStatusOffline,
			Rating:      5.0,
		})
	}
	if len(droneModels) > 0 {
		if err := s.droneRepo.BatchCreate(droneModels); err != nil {
			return 0, err
		}
	}
	return len(droneModels), nil
}

func (s *DroneService) Delete(id, ownerID uint) error {
	drone, err := s.droneRepo.GetByID(id)
	if err != nil {
		return err
	}
	if drone.OwnerID != ownerID {
		return errors.New("无权删除该设备")
	}
	return s.droneRepo.Delete(id)
}
