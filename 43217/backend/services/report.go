package services

import (
	"errors"
	"fmt"
	"time"

	"health-platform/models"
	"health-platform/repository"
	"health-platform/utils"
)

type ReportService struct {
	reportRepo     *repository.ReportRepository
	reportItemRepo *repository.ReportItemRepository
	appointmentRepo *repository.AppointmentRepository
	healthRepo     *repository.HealthRecordRepository
}

func NewReportService() *ReportService {
	return &ReportService{
		reportRepo:     repository.NewReportRepository(),
		reportItemRepo: repository.NewReportItemRepository(),
		appointmentRepo: repository.NewAppointmentRepository(),
		healthRepo:     repository.NewHealthRecordRepository(),
	}
}

type CreateReportRequest struct {
	AppointmentID uint               `json:"appointment_id" binding:"required"`
	ReportDate    string             `json:"report_date" binding:"required"`
	DoctorName    string             `json:"doctor_name" binding:"max=50"`
	Summary       string             `json:"summary"`
	Suggestion    string             `json:"suggestion"`
	PdfFile       string             `json:"pdf_file"`
	Items         []ReportItemRequest `json:"items"`
}

type ReportItemRequest struct {
	PackageItemID *uint  `json:"package_item_id"`
	ItemName      string `json:"item_name" binding:"required,max=100"`
	ItemCode      string `json:"item_code" binding:"max=50"`
	Result        string `json:"result" binding:"max=100"`
	Unit          string `json:"unit" binding:"max=20"`
	NormalRange   string `json:"normal_range" binding:"max=255"`
	IsAbnormal    bool   `json:"is_abnormal"`
	AbnormalType  string `json:"abnormal_type" binding:"max=20"`
	Description   string `json:"description"`
	Suggestion    string `json:"suggestion"`
	Department    string `json:"department" binding:"max=50"`
}

func (s *ReportService) CreateReport(req *CreateReportRequest) (*models.Report, error) {
	var appointment models.Appointment
	if err := s.appointmentRepo.FindByID(&appointment, req.AppointmentID); err != nil {
		return nil, errors.New("预约不存在")
	}

	if appointment.Status != models.AppointmentStatusCompleted {
		return nil, errors.New("该预约尚未完成")
	}

	reportDate, err := time.Parse("2006-01-02", req.ReportDate)
	if err != nil {
		return nil, errors.New("日期格式错误")
	}

	reportNo := utils.GenerateOrderNo("RPT")

	hasAbnormal := false
	for _, item := range req.Items {
		if item.IsAbnormal {
			hasAbnormal = true
			break
		}
	}

	report := &models.Report{
		AppointmentID: req.AppointmentID,
		EmployeeID:    appointment.EmployeeID,
		AgencyID:      appointment.AgencyID,
		PackageID:     appointment.PackageID,
		ReportNo:      reportNo,
		ReportDate:    reportDate,
		DoctorName:    req.DoctorName,
		Summary:       req.Summary,
		Suggestion:    req.Suggestion,
		HasAbnormal:   hasAbnormal,
		Status:        models.ReportStatusUploaded,
		PdfFile:       req.PdfFile,
	}

	if err := s.reportRepo.Create(report); err != nil {
		return nil, fmt.Errorf("创建报告失败: %w", err)
	}

	for i, item := range req.Items {
		reportItem := &models.ReportItem{
			ReportID:      report.ID,
			PackageItemID: item.PackageItemID,
			ItemName:      item.ItemName,
			ItemCode:      item.ItemCode,
			Result:        item.Result,
			Unit:          item.Unit,
			NormalRange:   item.NormalRange,
			IsAbnormal:    item.IsAbnormal,
			AbnormalType:  item.AbnormalType,
			Description:   item.Description,
			Suggestion:    item.Suggestion,
			Department:    item.Department,
		}
		_ = i
		if err := s.reportRepo.Create(reportItem); err != nil {
			return nil, fmt.Errorf("添加报告项失败: %w", err)
		}
	}

	return report, nil
}

func (s *ReportService) GetReport(reportID uint) (*models.Report, error) {
	return s.reportRepo.GetWithItems(reportID)
}

func (s *ReportService) GetReportByAppointment(appointmentID uint) (*models.Report, error) {
	return s.reportRepo.FindByAppointmentID(appointmentID)
}

func (s *ReportService) GetEmployeeReports(employeeID uint, page, pageSize int) ([]models.Report, int64, error) {
	return s.reportRepo.FindByEmployeeID(employeeID, page, pageSize)
}

func (s *ReportService) GetCompanyReports(companyID uint, page, pageSize int) ([]models.Report, int64, error) {
	return s.reportRepo.FindByCompanyID(companyID, page, pageSize)
}

func (s *ReportService) MarkReportViewed(reportID uint) error {
	return s.reportRepo.UpdateViewed(reportID)
}

func (s *ReportService) GetReportByNo(reportNo string) (*models.Report, error) {
	return s.reportRepo.GetReportByNo(reportNo)
}

func (s *ReportService) GetAbnormalReports(employeeID uint) ([]models.Report, error) {
	return s.reportRepo.GetAbnormalReportsByEmployee(employeeID)
}

func (s *ReportService) GetAbnormalItems(employeeID uint) ([]models.ReportItem, error) {
	return s.reportItemRepo.GetAbnormalItemsByEmployee(employeeID)
}

func (s *ReportService) UpdateHealthRecordFromReport(reportID uint) error {
	report, err := s.reportRepo.GetWithItems(reportID)
	if err != nil {
		return err
	}

	recordYear := report.ReportDate.Year()
	existingRecord, _ := s.healthRepo.GetByEmployeeAndYear(report.EmployeeID, recordYear)

	var height, weight float64
	var bmi float64
	var bloodPressure string
	var heartRate int
	abnormalCount := 0
	var tags string

	for _, item := range report.Items {
		if item.ItemName == "身高" {
			fmt.Sscanf(item.Result, "%f", &height)
		}
		if item.ItemName == "体重" {
			fmt.Sscanf(item.Result, "%f", &weight)
		}
		if item.ItemName == "血压" {
			bloodPressure = item.Result
		}
		if item.ItemName == "心率" {
			fmt.Sscanf(item.Result, "%d", &heartRate)
		}
		if item.IsAbnormal {
			abnormalCount++
			if tags != "" {
				tags += ","
			}
			tags += item.ItemName
		}
	}

	if height > 0 && weight > 0 {
		bmi = utils.CalcBMI(height, weight)
	}

	if existingRecord != nil {
		existingRecord.Height = height
		existingRecord.Weight = weight
		existingRecord.BMI = bmi
		existingRecord.BloodPressure = bloodPressure
		existingRecord.HeartRate = heartRate
		existingRecord.HasAbnormal = abnormalCount > 0
		existingRecord.AbnormalCount = abnormalCount
		existingRecord.Tags = tags
		existingRecord.Summary = report.Summary
		return s.healthRepo.Save(existingRecord)
	}

	record := &models.HealthRecord{
		EmployeeID:    report.EmployeeID,
		CompanyID:     0,
		RecordYear:    recordYear,
		ReportID:      &report.ID,
		Height:        height,
		Weight:        weight,
		BMI:           bmi,
		BloodPressure: bloodPressure,
		HeartRate:     heartRate,
		HasAbnormal:   abnormalCount > 0,
		AbnormalCount: abnormalCount,
		Tags:          tags,
		Summary:       report.Summary,
	}

	return s.healthRepo.Create(record)
}
