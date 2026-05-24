package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

func (h *StatsHandler) GetActivityStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	activityType := c.Query("activity_type")

	query := database.DB.Model(&models.JobPosting{})

	if startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}
	if activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}

	var totalJobs int64
	query.Count(&totalJobs)

	var recruitingCount int64
	database.DB.Model(&models.JobPosting{}).Where("status = ?", "recruiting").Count(&recruitingCount)

	var completedCount int64
	database.DB.Model(&models.JobPosting{}).Where("status = ?", "completed").Count(&completedCount)

	var totalHeadcount int64
	database.DB.Model(&models.JobPosting{}).Select("COALESCE(SUM(headcount), 0)").Scan(&totalHeadcount)

	var totalHired int64
	database.DB.Model(&models.JobPosting{}).Select("COALESCE(SUM(hired_count), 0)").Scan(&totalHired)

	type ActivityTypeStat struct {
		ActivityType string `json:"activity_type"`
		Count        int64  `json:"count"`
	}
	var byActivityType []ActivityTypeStat
	database.DB.Model(&models.JobPosting{}).
		Select("activity_type, COUNT(*) as count").
		Group("activity_type").
		Scan(&byActivityType)

	type MonthlyStat struct {
		Month string `json:"month"`
		Count int64  `json:"count"`
	}
	var byMonth []MonthlyStat
	database.DB.Raw(`
		SELECT TO_CHAR(start_date, 'YYYY-MM') as month, COUNT(*) as count
		FROM job_postings
		WHERE start_date >= DATE_TRUNC('year', CURRENT_DATE)
		GROUP BY month
		ORDER BY month
	`).Scan(&byMonth)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"total_jobs":        totalJobs,
			"recruiting_jobs":   recruitingCount,
			"completed_jobs":    completedCount,
			"total_headcount":   totalHeadcount,
			"total_hired":       totalHired,
			"by_activity_type":  byActivityType,
			"by_month":          byMonth,
		},
	})
}

func (h *StatsHandler) GetPersonnelStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	var employerCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleEmployer).Count(&employerCount)

	var agentCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleAgent).Count(&agentCount)

	var tempCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleTemporary).Count(&tempCount)

	var activeUsers int64
	database.DB.Model(&models.User{}).Where("status = ?", models.UserStatusActive).Count(&activeUsers)

	var newUsers int64
	newUsersQuery := database.DB.Model(&models.User{})
	if startDate != "" {
		newUsersQuery = newUsersQuery.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		newUsersQuery = newUsersQuery.Where("created_at <= ?", endDate)
	}
	newUsersQuery.Count(&newUsers)

	type CreditScoreDist struct {
		Range string `json:"range"`
		Count int64  `json:"count"`
	}
	var creditDist []CreditScoreDist
	database.DB.Raw(`
		SELECT 
			CASE 
				WHEN credit_score >= 80 THEN '80-100'
				WHEN credit_score >= 60 THEN '60-79'
				WHEN credit_score >= 40 THEN '40-59'
				ELSE '0-39'
			END as range,
			COUNT(*) as count
		FROM users
		WHERE role = 'temporary'
		GROUP BY range
		ORDER BY range
	`).Scan(&creditDist)

	type TopRatedTemp struct {
		ID          string `json:"id"`
		RealName    string `json:"real_name"`
		CreditScore int    `json:"credit_score"`
		RatingCount int64  `json:"rating_count"`
	}
	var topRated []TopRatedTemp
	database.DB.Raw(`
		SELECT u.id, u.real_name, u.credit_score, COUNT(e.id) as rating_count
		FROM users u
		LEFT JOIN evaluations e ON e.to_user_id = u.id
		WHERE u.role = 'temporary'
		GROUP BY u.id, u.real_name, u.credit_score
		ORDER BY u.credit_score DESC
		LIMIT 10
	`).Scan(&topRated)

	type DailyNewUser struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var dailyNew []DailyNewUser
	database.DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM users
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY date
		ORDER BY date
	`).Scan(&dailyNew)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"total_users":        totalUsers,
			"employer_count":     employerCount,
			"agent_count":        agentCount,
			"temporary_count":    tempCount,
			"active_users":       activeUsers,
			"new_users":          newUsers,
			"credit_distribution": creditDist,
			"top_rated_temps":    topRated,
			"daily_new_users":    dailyNew,
		},
	})
}

func (h *StatsHandler) GetSalaryStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var totalSalary float64
	salaryQuery := database.DB.Model(&models.SalaryRecord{}).Where("status = ?", "paid")
	if startDate != "" {
		salaryQuery = salaryQuery.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		salaryQuery = salaryQuery.Where("created_at <= ?", endDate)
	}
	salaryQuery.Select("COALESCE(SUM(total_salary), 0)").Scan(&totalSalary)

	var pendingSalary float64
	database.DB.Model(&models.SalaryRecord{}).Where("status = ?", "pending").Select("COALESCE(SUM(total_salary), 0)").Scan(&pendingSalary)

	var totalPaid int64
	database.DB.Model(&models.SalaryRecord{}).Where("status = ?", "paid").Count(&totalPaid)

	var totalHours float64
	database.DB.Model(&models.SalaryRecord{}).Where("status = ?", "paid").Select("COALESCE(SUM(total_hours), 0)").Scan(&totalHours)

	var totalOvertimePay float64
	database.DB.Model(&models.SalaryRecord{}).Where("status = ?", "paid").Select("COALESCE(SUM(overtime_pay), 0)").Scan(&totalOvertimePay)

	var totalDeductions float64
	database.DB.Model(&models.SalaryRecord{}).Where("status = ?", "paid").Select("COALESCE(SUM(deductions), 0)").Scan(&totalDeductions)

	type MonthlySalary struct {
		Month   string  `json:"month"`
		Total   float64 `json:"total"`
		Count   int64   `json:"count"`
	}
	var monthlySalary []MonthlySalary
	database.DB.Raw(`
		SELECT TO_CHAR(created_at, 'YYYY-MM') as month, 
			COALESCE(SUM(total_salary), 0) as total, 
			COUNT(*) as count
		FROM salary_records
		WHERE status = 'paid' AND created_at >= DATE_TRUNC('year', CURRENT_DATE)
		GROUP BY month
		ORDER BY month
	`).Scan(&monthlySalary)

	type TopEarners struct {
		ID           string  `json:"id"`
		RealName     string  `json:"real_name"`
		TotalEarnings float64 `json:"total_earnings"`
		JobCount     int64   `json:"job_count"`
	}
	var topEarners []TopEarners
	database.DB.Raw(`
		SELECT u.id, u.real_name, 
			COALESCE(SUM(s.total_salary), 0) as total_earnings,
			COUNT(s.id) as job_count
		FROM users u
		JOIN salary_records s ON s.temporary_id = u.id
		WHERE s.status = 'paid'
		GROUP BY u.id, u.real_name
		ORDER BY total_earnings DESC
		LIMIT 10
	`).Scan(&topEarners)

	type SalaryByPosition struct {
		Position      string  `json:"position"`
		TotalSalary   float64 `json:"total_salary"`
		AvgSalary     float64 `json:"avg_salary"`
		JobCount      int64   `json:"job_count"`
	}
	var salaryByPosition []SalaryByPosition
	database.DB.Raw(`
		SELECT j.position,
			COALESCE(SUM(s.total_salary), 0) as total_salary,
			COALESCE(AVG(s.total_salary), 0) as avg_salary,
			COUNT(s.id) as job_count
		FROM salary_records s
		JOIN job_postings j ON s.job_id = j.id
		WHERE s.status = 'paid'
		GROUP BY j.position
		ORDER BY total_salary DESC
		LIMIT 10
	`).Scan(&salaryByPosition)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"total_paid_salary":   totalSalary,
			"pending_salary":      pendingSalary,
			"total_paid_count":    totalPaid,
			"total_hours":         totalHours,
			"total_overtime_pay":  totalOvertimePay,
			"total_deductions":    totalDeductions,
			"monthly_salary":      monthlySalary,
			"top_earners":         topEarners,
			"salary_by_position":  salaryByPosition,
		},
	})
}

func (h *StatsHandler) GetOverviewStats(c *gin.Context) {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	var totalJobs int64
	database.DB.Model(&models.JobPosting{}).Count(&totalJobs)

	var activeJobs int64
	database.DB.Model(&models.JobPosting{}).Where("status = ?", "recruiting").Count(&activeJobs)

	var totalSchedules int64
	database.DB.Model(&models.Schedule{}).Count(&totalSchedules)

	var todayCheckIns int64
	database.DB.Model(&models.CheckIn{}).Where("DATE(check_in_time) = ?", now.Format("2006-01-02")).Count(&todayCheckIns)

	var monthlySalary float64
	database.DB.Model(&models.SalaryRecord{}).
		Where("status = ? AND created_at >= ?", "paid", firstDayOfMonth).
		Select("COALESCE(SUM(total_salary), 0)").Scan(&monthlySalary)

	var newUsersThisMonth int64
	database.DB.Model(&models.User{}).
		Where("created_at >= ?", firstDayOfMonth).
		Count(&newUsersThisMonth)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"total_users":        totalUsers,
			"total_jobs":         totalJobs,
			"active_jobs":        activeJobs,
			"total_schedules":    totalSchedules,
			"today_check_ins":    todayCheckIns,
			"monthly_salary":     monthlySalary,
			"new_users_month":    newUsersThisMonth,
		},
	})
}
