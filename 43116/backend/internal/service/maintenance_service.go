package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"errors"
	"time"
)

type MaintenanceService struct {
	maintenanceRepo *repository.MaintenanceRepository
	carRepo         *repository.CarRepository
}

func NewMaintenanceService() *MaintenanceService {
	return &MaintenanceService{
		maintenanceRepo: repository.NewMaintenanceRepository(),
		carRepo:         repository.NewCarRepository(),
	}
}

type CreateMaintenanceRequest struct {
	CarID       uint      `json:"car_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Cost        float64   `json:"cost"`
	Notes       string    `json:"notes"`
}

func (s *MaintenanceService) CreateMaintenance(createdBy uint, req *CreateMaintenanceRequest) (*model.MaintenancePlan, error) {
	_, err := s.carRepo.FindByID(req.CarID)
	if err != nil {
		return nil, errors.New("车辆不存在")
	}

	if !req.EndDate.After(req.StartDate) {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	plan := &model.MaintenancePlan{
		CarID:       req.CarID,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Cost:        req.Cost,
		Notes:       req.Notes,
		Status:      model.MaintenanceStatusScheduled,
		CreatedBy:   createdBy,
	}

	err = s.maintenanceRepo.Create(plan)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (s *MaintenanceService) GetMaintenanceByID(id uint) (*model.MaintenancePlan, error) {
	return s.maintenanceRepo.FindByID(id)
}

func (s *MaintenanceService) GetAllMaintenance(page, pageSize int, carID uint, status string) ([]model.MaintenancePlan, int64, error) {
	return s.maintenanceRepo.FindAll(page, pageSize, carID, status)
}

func (s *MaintenanceService) GetCarMaintenance(carID uint, page, pageSize int) ([]model.MaintenancePlan, int64, error) {
	return s.maintenanceRepo.FindByCarID(carID, page, pageSize)
}

func (s *MaintenanceService) StartMaintenance(id uint) error {
	plan, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return errors.New("维护计划不存在")
	}

	if plan.Status != model.MaintenanceStatusScheduled {
		return errors.New("维护计划状态不允许开始")
	}

	now := time.Now()
	plan.ActualStart = &now
	plan.Status = model.MaintenanceStatusInProgress

	err = s.maintenanceRepo.Update(plan)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateStatus(plan.CarID, model.CarStatusMaintenance)
}

func (s *MaintenanceService) CompleteMaintenance(id uint) error {
	plan, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return errors.New("维护计划不存在")
	}

	if plan.Status != model.MaintenanceStatusInProgress {
		return errors.New("维护计划状态不允许完成")
	}

	now := time.Now()
	plan.ActualEnd = &now
	plan.Status = model.MaintenanceStatusCompleted

	err = s.maintenanceRepo.Update(plan)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateStatus(plan.CarID, model.CarStatusAvailable)
}

func (s *MaintenanceService) CancelMaintenance(id uint) error {
	plan, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return errors.New("维护计划不存在")
	}

	if plan.Status == model.MaintenanceStatusCompleted {
		return errors.New("已完成的维护计划不能取消")
	}

	plan.Status = model.MaintenanceStatusCancelled
	err = s.maintenanceRepo.Update(plan)
	if err != nil {
		return err
	}

	if plan.Status == model.MaintenanceStatusInProgress {
		return s.carRepo.UpdateStatus(plan.CarID, model.CarStatusAvailable)
	}

	return nil
}

func (s *MaintenanceService) UpdateMaintenance(id uint, updates map[string]interface{}) error {
	plan, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return errors.New("维护计划不存在")
	}

	if title, ok := updates["title"]; ok {
		plan.Title = title.(string)
	}
	if description, ok := updates["description"]; ok {
		plan.Description = description.(string)
	}
	if startDate, ok := updates["start_date"]; ok {
		plan.StartDate, _ = time.Parse(time.RFC3339, startDate.(string))
	}
	if endDate, ok := updates["end_date"]; ok {
		plan.EndDate, _ = time.Parse(time.RFC3339, endDate.(string))
	}
	if cost, ok := updates["cost"]; ok {
		plan.Cost = cost.(float64)
	}
	if notes, ok := updates["notes"]; ok {
		plan.Notes = notes.(string)
	}

	return s.maintenanceRepo.Update(plan)
}

func (s *MaintenanceService) DeleteMaintenance(id uint) error {
	plan, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return errors.New("维护计划不存在")
	}

	if plan.Status == model.MaintenanceStatusInProgress {
		return errors.New("进行中的维护计划不能删除")
	}

	return s.maintenanceRepo.Delete(id)
}

func (s *MaintenanceService) GetUpcomingMaintenance() ([]model.MaintenancePlan, error) {
	return s.maintenanceRepo.FindUpcoming()
}