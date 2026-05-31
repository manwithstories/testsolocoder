package services

import (
	"time"

	"health-platform/config"
	"health-platform/models"
	"health-platform/repository"
)

type StatisticsService struct {
	appointmentRepo *repository.AppointmentRepository
	employeeRepo    *repository.EmployeeRepository
	reportRepo      *repository.ReportRepository
	billingRepo     *repository.BillingRepository
	packageRepo     *repository.PackageRepository
	healthRepo      *repository.HealthRecordRepository
}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{
		appointmentRepo: repository.NewAppointmentRepository(),
		employeeRepo:    repository.NewEmployeeRepository(),
		reportRepo:      repository.NewReportRepository(),
		billingRepo:     repository.NewBillingRepository(),
		packageRepo:     repository.NewPackageRepository(),
		healthRepo:      repository.NewHealthRecordRepository(),
	}
}

type StatisticsQuery struct {
	CompanyID    uint   `json:"company_id"`
	DepartmentID *uint  `json:"department_id"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Year         int    `json:"year"`
	Gender       *int   `json:"gender"`
	AgeRange     string `json:"age_range"`
}

func (s *StatisticsService) GetCompanyStatistics(query StatisticsQuery) (map[string]interface{}, error) {
	cacheKey := "company_stats_"
	if query.CompanyID > 0 {
		cacheKey += string(rune(query.CompanyID))
	}
	if query.Year > 0 {
		cacheKey += "_" + string(rune(query.Year))
	}

	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		return map[string]interface{}{"cached": true, "data": cachedData}, nil
	}

	totalEmployees, _ := s.employeeRepo.CountByCompanyID(query.CompanyID)
	
	var startDate, endDate time.Time
	if query.StartDate != "" {
		startDate, _ = time.Parse("2006-01-02", query.StartDate)
	} else {
		startDate = time.Date(query.Year, 1, 1, 0, 0, 0, 0, time.Local)
	}
	if query.EndDate != "" {
		endDate, _ = time.Parse("2006-01-02", query.EndDate)
	} else {
		endDate = time.Now()
	}

	appointments, _, _ := s.appointmentRepo.FindByCompanyIDAndDateRange(query.CompanyID, startDate, endDate, 1, 100000)
	
	var completedCount, pendingCount, cancelledCount int64
	for _, appt := range appointments {
		switch appt.Status {
		case models.AppointmentStatusCompleted:
			completedCount++
		case models.AppointmentStatusPending, models.AppointmentStatusConfirmed:
			pendingCount++
		case models.AppointmentStatusCancelled:
			cancelledCount++
		}
	}

	completionRate := float64(0)
	if totalEmployees > 0 {
		completionRate = float64(completedCount) / float64(totalEmployees) * 100
	}

	reports, _ := s.reportRepo.FindByCompanyIDAndDateRange(query.CompanyID, startDate, endDate, 1, 100000)
	
	var abnormalCount int
	for _, report := range reports {
		if report.HasAbnormal {
			abnormalCount++
		}
	}

	billings, _ := s.billingRepo.FindByCompanyID(query.CompanyID, 1, 100000)
	var totalAmount, paidAmount float64
	for _, billing := range billings {
		totalAmount += billing.TotalAmount
		paidAmount += billing.PaidAmount
	}

	result := map[string]interface{}{
		"total_employees":  totalEmployees,
		"total_appointments": len(appointments),
		"completed_count":  completedCount,
		"pending_count":    pendingCount,
		"cancelled_count":  cancelledCount,
		"completion_rate":  completionRate,
		"total_reports":    len(reports),
		"abnormal_count":   abnormalCount,
		"abnormal_rate":    float64(abnormalCount) / float64(len(reports)) * 100,
		"total_billing":    totalAmount,
		"paid_amount":      paidAmount,
	}

	config.RedisClient.Set(config.Ctx, cacheKey, result, 30*time.Minute)

	return result, nil
}

func (s *StatisticsService) GetDepartmentStatistics(companyID uint, departmentID uint, year int) (map[string]interface{}, error) {
	employees, _ := s.employeeRepo.FindByDepartmentID(departmentID)
	
	totalEmployees := len(employees)
	
	appointments, _ := s.appointmentRepo.GetDepartmentAppointmentCount(departmentID, year)
	
	healthRecords := make([]models.HealthRecord, 0)
	for _, emp := range employees {
		records, _ := s.healthRepo.FindByEmployeeID(emp.ID, 1, 1000)
		healthRecords = append(healthRecords, records...)
	}

	var abnormalCount int
	for _, record := range healthRecords {
		if record.HasAbnormal {
			abnormalCount++
		}
	}

	return map[string]interface{}{
		"department_id":    departmentID,
		"total_employees":  totalEmployees,
		"appointment_count": appointments,
		"health_records":   len(healthRecords),
		"abnormal_count":   abnormalCount,
		"year":             year,
	}, nil
}

func (s *StatisticsService) GetAgeDistribution(companyID uint) (map[string]interface{}, error) {
	employees, _, _ := s.employeeRepo.FindByCompanyID(companyID, 1, 100000)
	
	distribution := map[string]int{
		"under_25":  0,
		"25_30":     0,
		"31_35":     0,
		"36_40":     0,
		"41_45":     0,
		"46_50":     0,
		"above_50":  0,
	}

	for _, emp := range employees {
		age := time.Now().Year() - emp.Birthday.Year()
		if emp.Birthday != nil {
			switch {
			case age < 25:
				distribution["under_25"]++
			case age <= 30:
				distribution["25_30"]++
			case age <= 35:
				distribution["31_35"]++
			case age <= 40:
				distribution["36_40"]++
			case age <= 45:
				distribution["41_45"]++
			case age <= 50:
				distribution["46_50"]++
			default:
				distribution["above_50"]++
			}
		}
	}

	return map[string]interface{}{
		"distribution": distribution,
		"total":        len(employees),
	}, nil
}

func (s *StatisticsService) GetGenderDistribution(companyID uint) (map[string]interface{}, error) {
	employees, _, _ := s.employeeRepo.FindByCompanyID(companyID, 1, 100000)
	
	maleCount := 0
	femaleCount := 0
	for _, emp := range employees {
		if emp.Gender == 1 {
			maleCount++
		} else {
			femaleCount++
		}
	}

	return map[string]interface{}{
		"male":   maleCount,
		"female": femaleCount,
		"total":  len(employees),
	}, nil
}

func (s *StatisticsService) GetAbnormalDistribution(companyID uint, year int) (map[string]interface{}, error) {
	healthRecords := make([]models.HealthRecord, 0)
	employees, _, _ := s.employeeRepo.FindByCompanyID(companyID, 1, 100000)
	for _, emp := range employees {
		records, _ := s.healthRepo.FindByEmployeeID(emp.ID, 1, 1000)
		healthRecords = append(healthRecords, records...)
	}

	abnormalTypeCount := make(map[string]int)
	for _, record := range healthRecords {
		if record.Tags != "" {
			tags := splitTags(record.Tags)
			for _, tag := range tags {
				abnormalTypeCount[tag]++
			}
		}
	}

	return map[string]interface{}{
		"abnormal_types": abnormalTypeCount,
		"total_records":  len(healthRecords),
	}, nil
}

func splitTags(tags string) []string {
	var result []string
	current := ""
	for _, char := range tags {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func (s *StatisticsService) GetAgencyRatingStatistics() (map[string]interface{}, error) {
	cacheKey := "agency_ranking"
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		return map[string]interface{}{"cached": true, "data": cachedData}, nil
	}

	agencies, _ := s.appointmentRepo.GetAgencyRanking()

	result := map[string]interface{}{
		"agency_rankings": agencies,
	}

	config.RedisClient.Set(config.Ctx, cacheKey, result, 30*time.Minute)

	return result, nil
}

func (s *StatisticsService) GetPackageRanking(limit int) ([]models.Package, error) {
	cacheKey := "package_ranking"
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		return nil, nil
	}
	_ = cachedData

	return s.packageRepo.GetHotPackages(limit)
}

func (s *StatisticsService) ClearCache() error {
	keys := []string{
		"company_stats_*",
		"agency_ranking",
		"package_ranking",
	}
	for _, key := range keys {
		config.RedisClient.Del(config.Ctx, key)
	}
	return nil
}
