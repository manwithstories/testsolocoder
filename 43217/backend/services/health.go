package services

import (
	"errors"
	"fmt"
	"time"

	"health-platform/models"
	"health-platform/repository"
)

type HealthService struct {
	healthRepo    *repository.HealthRecordRepository
	abnormalRepo  *repository.AbnormalItemRepository
	reminderRepo  *repository.RecheckReminderRepository
	reportRepo    *repository.ReportRepository
	employeeRepo  *repository.EmployeeRepository
}

func NewHealthService() *HealthService {
	return &HealthService{
		healthRepo:   repository.NewHealthRecordRepository(),
		abnormalRepo: repository.NewAbnormalItemRepository(),
		reminderRepo: repository.NewRecheckReminderRepository(),
		reportRepo:   repository.NewReportRepository(),
		employeeRepo: repository.NewEmployeeRepository(),
	}
}

type CreateAbnormalItemRequest struct {
	EmployeeID     uint       `json:"employee_id" binding:"required"`
	HealthRecordID uint      `json:"health_record_id" binding:"required"`
	ItemName       string     `json:"item_name" binding:"required,max=100"`
	ItemCode       string     `json:"item_code" binding:"max=50"`
	Result         string     `json:"result" binding:"max=100"`
	NormalRange    string     `json:"normal_range" binding:"max=255"`
	AbnormalType   string     `json:"abnormal_type" binding:"max=20"`
	Level          int        `json:"level"`
	Suggestion     string     `json:"suggestion"`
	RecheckDate    *time.Time `json:"recheck_date"`
}

type SetRecheckRequest struct {
	AbnormalID  uint       `json:"abnormal_id" binding:"required"`
	RecheckDate *time.Time `json:"recheck_date" binding:"required"`
	Suggestion  string     `json:"suggestion"`
}

type UpdateRecheckStatusRequest struct {
	AbnormalID uint `json:"abnormal_id" binding:"required"`
	Status     int  `json:"status" binding:"required"`
}

func (s *HealthService) GetHealthRecords(employeeID uint, page, pageSize int) ([]models.HealthRecord, int64, error) {
	return s.healthRepo.FindByEmployeeID(employeeID, page, pageSize)
}

func (s *HealthService) GetHealthRecordByYear(employeeID uint, year int) (*models.HealthRecord, error) {
	return s.healthRepo.GetByEmployeeAndYear(employeeID, year)
}

func (s *HealthService) GetAllHealthRecords(employeeID uint) ([]models.HealthRecord, error) {
	return s.healthRepo.GetAllByEmployee(employeeID)
}

func (s *HealthService) GetTrendData(employeeID uint) (map[string]interface{}, error) {
	records, err := s.healthRepo.GetAllByEmployee(employeeID)
	if err != nil {
		return nil, err
	}

	years := make([]int, 0)
	bmiData := make([]float64, 0)
	weightData := make([]float64, 0)
	heightData := make([]float64, 0)
	abnormalData := make([]int, 0)

	for _, record := range records {
		years = append(years, record.RecordYear)
		bmiData = append(bmiData, record.BMI)
		weightData = append(weightData, record.Weight)
		heightData = append(heightData, record.Height)
		abnormalData = append(abnormalData, record.AbnormalCount)
	}

	return map[string]interface{}{
		"years":     years,
		"bmi":       bmiData,
		"weight":    weightData,
		"height":    heightData,
		"abnormal":  abnormalData,
		"records":   records,
	}, nil
}

func (s *HealthService) CreateAbnormalItem(req *CreateAbnormalItemRequest) (*models.AbnormalItem, error) {
	item := &models.AbnormalItem{
		EmployeeID:     req.EmployeeID,
		HealthRecordID: req.HealthRecordID,
		ItemName:       req.ItemName,
		ItemCode:       req.ItemCode,
		Result:         req.Result,
		NormalRange:    req.NormalRange,
		AbnormalType:   req.AbnormalType,
		Level:          req.Level,
		Suggestion:     req.Suggestion,
		RecheckDate:    req.RecheckDate,
	}

	if err := s.abnormalRepo.Create(item); err != nil {
		return nil, fmt.Errorf("创建异常项失败: %w", err)
	}

	if req.RecheckDate != nil {
		reminder := &models.RecheckReminder{
			EmployeeID: req.EmployeeID,
			AbnormalID: item.ID,
			RemindDate: *req.RecheckDate,
			RemindType: "复查提醒",
			Content:    fmt.Sprintf("您的%s指标异常，建议在%s前复查", req.ItemName, req.RecheckDate.Format("2006-01-02")),
		}
		s.reminderRepo.Create(reminder)
	}

	return item, nil
}

func (s *HealthService) GetAbnormalItems(employeeID uint, page, pageSize int) ([]models.AbnormalItem, int64, error) {
	return s.abnormalRepo.FindByEmployeeID(employeeID, page, pageSize)
}

func (s *HealthService) GetAllAbnormalItems(employeeID uint) ([]models.AbnormalItem, error) {
	return s.abnormalRepo.GetAllByEmployee(employeeID)
}

func (s *HealthService) SetRecheckDate(req *SetRecheckRequest) error {
	var item models.AbnormalItem
	if err := s.abnormalRepo.FindByID(&item, req.AbnormalID); err != nil {
		return errors.New("异常项不存在")
	}

	item.RecheckDate = req.RecheckDate
	if req.Suggestion != "" {
		item.Suggestion = req.Suggestion
	}

	if err := s.abnormalRepo.Save(&item); err != nil {
		return err
	}

	reminder := &models.RecheckReminder{
		EmployeeID: item.EmployeeID,
		AbnormalID: item.ID,
		RemindDate: *req.RecheckDate,
		RemindType: "复查提醒",
		Content:    fmt.Sprintf("您的%s指标异常，建议在%s前复查", item.ItemName, req.RecheckDate.Format("2006-01-02")),
	}
	return s.reminderRepo.Create(reminder)
}

func (s *HealthService) UpdateRecheckStatus(req *UpdateRecheckStatusRequest) error {
	return s.abnormalRepo.UpdateRecheckStatus(req.AbnormalID, req.Status)
}

func (s *HealthService) GetNeedRecheckItems(employeeID uint) ([]models.AbnormalItem, error) {
	return s.abnormalRepo.GetNeedRecheckItems(employeeID)
}

func (s *HealthService) GetReminders(employeeID uint, page, pageSize int) ([]models.RecheckReminder, int64, error) {
	return s.reminderRepo.FindByEmployeeID(employeeID, page, pageSize)
}

func (s *HealthService) GetUnreadReminders(employeeID uint) ([]models.RecheckReminder, error) {
	return s.reminderRepo.GetUnreadByEmployee(employeeID)
}

func (s *HealthService) MarkReminderAsRead(reminderID uint) error {
	return s.reminderRepo.MarkAsRead(reminderID)
}

func (s *HealthService) GenerateHealthSummary(employeeID uint) (map[string]interface{}, error) {
	records, err := s.healthRepo.GetAllByEmployee(employeeID)
	if err != nil {
		return nil, err
	}

	var totalAbnormal int
	var totalRecords int
	var latestRecord *models.HealthRecord
	var tags map[string]int = make(map[string]int)

	for i := range records {
		totalRecords++
		totalAbnormal += records[i].AbnormalCount
		if latestRecord == nil || records[i].RecordYear > latestRecord.RecordYear {
			latestRecord = &records[i]
		}
	}

	abnormalItems, _ := s.abnormalRepo.GetAllByEmployee(employeeID)
	needRecheck, _ := s.abnormalRepo.GetNeedRecheckItems(employeeID)

	return map[string]interface{}{
		"total_records":    totalRecords,
		"total_abnormal":   totalAbnormal,
		"latest_record":    latestRecord,
		"need_recheck":     len(needRecheck),
		"abnormal_items":   len(abnormalItems),
		"tags":             tags,
	}, nil
}
