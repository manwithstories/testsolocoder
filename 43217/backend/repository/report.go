package repository

import (
	"health-platform/models"
	"time"
)

type ReportRepository struct {
	*BaseRepository
}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *ReportRepository) FindByEmployeeID(employeeID uint, page, pageSize int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := r.DB.Model(&models.Report{}).Where("employee_id = ?", employeeID)
	query.Count(&total)

	err := query.Preload("Items").Preload("Agency").Preload("Appointment").
		Order("report_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reports).Error
	return reports, total, err
}

func (r *ReportRepository) FindByAppointmentID(appointmentID uint) (*models.Report, error) {
	var report models.Report
	err := r.DB.Where("appointment_id = ?", appointmentID).
		Preload("Items").Preload("Appointment").First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) GetWithItems(reportID uint) (*models.Report, error) {
	var report models.Report
	err := r.DB.Preload("Items").Preload("Appointment").Preload("Employee").
		First(&report, reportID).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) FindByCompanyID(companyID uint, page, pageSize int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := r.DB.Model(&models.Report{}).
		Joins("JOIN appointments ON appointments.id = reports.appointment_id").
		Where("appointments.company_id = ?", companyID)
	query.Count(&total)

	err := query.Preload("Items").Preload("Employee").Preload("Agency").
		Order("report_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reports).Error
	return reports, total, err
}

func (r *ReportRepository) FindByCompanyIDAndDateRange(companyID uint, startDate, endDate time.Time, page, pageSize int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := r.DB.Model(&models.Report{}).
		Joins("JOIN appointments ON appointments.id = reports.appointment_id").
		Where("appointments.company_id = ? AND reports.report_date >= ? AND reports.report_date <= ?", 
			companyID, startDate, endDate)
	query.Count(&total)

	err := query.Preload("Items").Preload("Employee").Preload("Agency").
		Order("report_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reports).Error
	return reports, total, err
}

func (r *ReportRepository) UpdateViewed(reportID uint) error {
	return r.DB.Model(&models.Report{}).Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"status":    models.ReportStatusViewed,
			"viewed_at": time.Now(),
		}).Error
}

func (r *ReportRepository) GetReportByNo(reportNo string) (*models.Report, error) {
	var report models.Report
	err := r.DB.Where("report_no = ?", reportNo).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) GetEmployeeReportsByYear(employeeID uint, year int) ([]models.Report, error) {
	var reports []models.Report
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)
	
	err := r.DB.Where("employee_id = ? AND report_date >= ? AND report_date < ?",
		employeeID, startDate, endDate).
		Preload("Items").
		Order("report_date DESC").
		Find(&reports).Error
	return reports, err
}

func (r *ReportRepository) GetAbnormalReportsByEmployee(employeeID uint) ([]models.Report, error) {
	var reports []models.Report
	err := r.DB.Where("employee_id = ? AND has_abnormal = ?", employeeID, true).
		Preload("Items").
		Order("report_date DESC").
		Find(&reports).Error
	return reports, err
}

type ReportItemRepository struct {
	*BaseRepository
}

func NewReportItemRepository() *ReportItemRepository {
	return &ReportItemRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *ReportItemRepository) FindByReportID(reportID uint) ([]models.ReportItem, error) {
	var items []models.ReportItem
	err := r.DB.Where("report_id = ?", reportID).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *ReportItemRepository) GetAbnormalItemsByEmployee(employeeID uint) ([]models.ReportItem, error) {
	var items []models.ReportItem
	err := r.DB.Joins("JOIN reports ON reports.id = report_items.report_id").
		Where("reports.employee_id = ? AND report_items.is_abnormal = ?", employeeID, true).
		Order("reports.report_date DESC").
		Find(&items).Error
	return items, err
}
