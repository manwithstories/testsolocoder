package services

import (
	"time"

	"business-registration-platform/database"
	"business-registration-platform/models"
)

type StatisticsService struct{}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

type OverviewStats struct {
	TotalApplications    int64   `json:"totalApplications"`
	PendingApplications  int64   `json:"pendingApplications"`
	ProcessingApps       int64   `json:"processingApps"`
	CompletedApps        int64   `json:"completedApps"`
	RejectedApps         int64   `json:"rejectedApps"`
	TotalRevenue         float64 `json:"totalRevenue"`
	TodayApplications    int64   `json:"todayApplications"`
	TodayRevenue         float64 `json:"todayRevenue"`
	TotalAgents          int64   `json:"totalAgents"`
	ActiveAgents         int64   `json:"activeAgents"`
	TotalEntrepreneurs   int64   `json:"totalEntrepreneurs"`
}

type StatusDistribution struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type CompanyTypeDistribution struct {
	CompanyType string `json:"companyType"`
	Count       int64  `json:"count"`
	TotalAmount float64 `json:"totalAmount"`
}

type AgentPerformance struct {
	AgentID        uint    `json:"agentId"`
	AgentName      string  `json:"agentName"`
	EmployeeNo     string  `json:"employeeNo"`
	TotalHandled   int64   `json:"totalHandled"`
	CompletedCount int64   `json:"completedCount"`
	InProgressCount int64  `json:"inProgressCount"`
	AvgDuration    float64 `json:"avgDuration"`
	TotalRevenue   float64 `json:"totalRevenue"`
}

type TimeSeriesData struct {
	Date  string  `json:"date"`
	Count int64   `json:"count"`
	Amount float64 `json:"amount"`
}

func (s *StatisticsService) GetOverviewStats(startDate, endDate *time.Time) (*OverviewStats, error) {
	var stats OverviewStats

	db := database.DB.Model(&models.Application{})

	if startDate != nil {
		db = db.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		db = db.Where("created_at <= ?", *endDate)
	}

	db.Count(&stats.TotalApplications)

	database.DB.Model(&models.Application{}).Where("status = ?", models.AppStatusPendingReview).Count(&stats.PendingApplications)
	database.DB.Model(&models.Application{}).Where("status = ?", models.AppStatusProcessing).Count(&stats.ProcessingApps)
	database.DB.Model(&models.Application{}).Where("status = ?", models.AppStatusCompleted).Count(&stats.CompletedApps)
	database.DB.Model(&models.Application{}).Where("status = ?", models.AppStatusRejected).Count(&stats.RejectedApps)

	var feeResult struct {
		Total float64
	}
	feeDB := database.DB.Model(&models.ApplicationFee{}).Where("status = ?", models.FeeStatusPaid)
	if startDate != nil {
		feeDB = feeDB.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		feeDB = feeDB.Where("created_at <= ?", *endDate)
	}
	feeDB.Select("COALESCE(SUM(paid_amount), 0) as total").Scan(&feeResult)
	stats.TotalRevenue = feeResult.Total

	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	database.DB.Model(&models.Application{}).Where("created_at >= ? AND created_at < ?", today, tomorrow).Count(&stats.TodayApplications)

	var todayFeeResult struct {
		Total float64
	}
	database.DB.Model(&models.ApplicationFee{}).
		Where("status = ? AND created_at >= ? AND created_at < ?", models.FeeStatusPaid, today, tomorrow).
		Select("COALESCE(SUM(paid_amount), 0) as total").Scan(&todayFeeResult)
	stats.TodayRevenue = todayFeeResult.Total

	database.DB.Model(&models.User{}).Where("role = ?", models.RoleAgent).Count(&stats.TotalAgents)
	database.DB.Model(&models.AgentProfile{}).Where("status = ?", "available").Count(&stats.ActiveAgents)
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleEntrepreneur).Count(&stats.TotalEntrepreneurs)

	return &stats, nil
}

func (s *StatisticsService) GetStatusDistribution() ([]StatusDistribution, error) {
	var distributions []StatusDistribution

	rows, err := database.DB.Model(&models.Application{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d StatusDistribution
		rows.Scan(&d.Status, &d.Count)
		distributions = append(distributions, d)
	}

	return distributions, nil
}

func (s *StatisticsService) GetCompanyTypeDistribution() ([]CompanyTypeDistribution, error) {
	var distributions []CompanyTypeDistribution

	rows, err := database.DB.Table("applications a").
		Select("a.company_type, COUNT(*) as count, COALESCE(SUM(f.paid_amount), 0) as total_amount").
		Joins("LEFT JOIN application_fees f ON f.application_id = a.id AND f.status = ?", models.FeeStatusPaid).
		Group("a.company_type").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d CompanyTypeDistribution
		rows.Scan(&d.CompanyType, &d.Count, &d.TotalAmount)
		distributions = append(distributions, d)
	}

	return distributions, nil
}

func (s *StatisticsService) GetAgentPerformance(startDate, endDate *time.Time) ([]AgentPerformance, error) {
	var performances []AgentPerformance

	db := database.DB.Table("users u").
		Select(`u.id as agent_id, u.real_name as agent_name, 
			ap.employee_no as employee_no,
			COUNT(a.id) as total_handled,
			SUM(CASE WHEN a.status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN a.status = 'processing' THEN 1 ELSE 0 END) as in_progress_count,
			COALESCE(SUM(f.paid_amount), 0) as total_revenue`).
		Joins("LEFT JOIN agent_profiles ap ON ap.user_id = u.id").
		Joins("LEFT JOIN applications a ON a.agent_id = u.id").
		Joins("LEFT JOIN application_fees f ON f.application_id = a.id AND f.status = ?", models.FeeStatusPaid).
		Where("u.role = ?", models.RoleAgent).
		Group("u.id")

	if startDate != nil {
		db = db.Where("a.created_at >= ?", *startDate)
	}
	if endDate != nil {
		db = db.Where("a.created_at <= ?", *endDate)
	}

	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p AgentPerformance
		rows.Scan(&p.AgentID, &p.AgentName, &p.EmployeeNo, &p.TotalHandled, &p.CompletedCount, &p.InProgressCount, &p.TotalRevenue)
		performances = append(performances, p)
	}

	return performances, nil
}

func (s *StatisticsService) GetApplicationTimeSeries(startDate, endDate time.Time, interval string) ([]TimeSeriesData, error) {
	var data []TimeSeriesData

	var dateFormat string
	switch interval {
	case "month":
		dateFormat = "%Y-%m"
	case "week":
		dateFormat = "%Y-%u"
	default:
		dateFormat = "%Y-%m-%d"
	}

	rows, err := database.DB.Table("applications a").
		Select(`DATE_FORMAT(a.created_at, ?) as date, 
			COUNT(*) as count, 
			COALESCE(SUM(f.paid_amount), 0) as amount`, dateFormat).
		Joins("LEFT JOIN application_fees f ON f.application_id = a.id AND f.status = ?", models.FeeStatusPaid).
		Where("a.created_at >= ? AND a.created_at <= ?", startDate, endDate).
		Group("date").
		Order("date ASC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d TimeSeriesData
		rows.Scan(&d.Date, &d.Count, &d.Amount)
		data = append(data, d)
	}

	return data, nil
}
