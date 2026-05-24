package repository

import (
	"time"

	"gorm.io/gorm"

	"recruitment-platform/internal/models"
)

type StatisticsRepository struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) *StatisticsRepository {
	return &StatisticsRepository{db: db}
}

func (r *StatisticsRepository) GetDailyStats(date time.Time) (*models.DailyStatistics, error) {
	var stats models.DailyStatistics
	err := r.db.Where("date = ?", date.Format("2006-01-02")).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *StatisticsRepository) SaveDailyStats(stats *models.DailyStatistics) error {
	return r.db.Save(stats).Error
}

func (r *StatisticsRepository) GetDateRangeStats(startDate, endDate time.Time) ([]models.DailyStatistics, error) {
	var stats []models.DailyStatistics
	err := r.db.Where("date >= ? AND date <= ?",
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("date ASC").Find(&stats).Error
	return stats, err
}

func (r *StatisticsRepository) GetJobStats(companyID uint, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []struct {
		JobID       uint
		Title       string
		TotalViews  int64
		ApplyCount  int64
		InterviewCount int64
		HireCount   int64
	}

	query := r.db.Table("jobs").
		Select(`jobs.id as job_id, jobs.title, 
			COALESCE(view_stats.view_count, 0) as total_views,
			jobs.apply_count,
			COALESCE(interview_stats.interview_count, 0) as interview_count,
			COALESCE(hire_stats.hire_count, 0) as hire_count`).
		Joins(`LEFT JOIN (SELECT job_id, COUNT(*) as view_count FROM job_view_logs 
			WHERE viewed_at >= ? AND viewed_at <= ? GROUP BY job_id) as view_stats 
			ON view_stats.job_id = jobs.id`, startDate, endDate).
		Joins(`LEFT JOIN (SELECT applications.job_id, COUNT(*) as interview_count FROM applications 
			WHERE applications.status = 'interview' AND applications.last_update_at >= ? 
			AND applications.last_update_at <= ? GROUP BY applications.job_id) as interview_stats 
			ON interview_stats.job_id = jobs.id`, startDate, endDate).
		Joins(`LEFT JOIN (SELECT applications.job_id, COUNT(*) as hire_count FROM applications 
			WHERE applications.status = 'accepted' AND applications.last_update_at >= ? 
			AND applications.last_update_at <= ? GROUP BY applications.job_id) as hire_stats 
			ON hire_stats.job_id = jobs.id`, startDate, endDate)

	if companyID > 0 {
		query = query.Where("jobs.company_id = ?", companyID)
	}

	err := query.Scan(&results).Error

	var stats []map[string]interface{}
	for _, r := range results {
		stats = append(stats, map[string]interface{}{
			"job_id":          r.JobID,
			"title":           r.Title,
			"total_views":     r.TotalViews,
			"apply_count":     r.ApplyCount,
			"interview_count": r.InterviewCount,
			"hire_count":      r.HireCount,
			"conversion_rate": float64(0),
		})
		if len(stats) > 0 && r.ApplyCount > 0 {
			stats[len(stats)-1]["conversion_rate"] = float64(r.HireCount) / float64(r.ApplyCount) * 100
		}
	}

	return stats, err
}

func (r *StatisticsRepository) GetRecruitmentCycleStats(companyID uint) ([]map[string]interface{}, error) {
	var results []struct {
		JobID        uint
		Title        string
		AppliedCount int64
		AvgDays      float64
	}

	query := r.db.Table("applications").
		Select(`applications.job_id, jobs.title,
			COUNT(*) as applied_count,
			AVG(JULIANDAY(applications.last_update_at) - JULIANDAY(applications.applied_at)) * 24 as avg_days`).
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("applications.status IN ?", []models.ApplicationStatus{
			models.ApplicationStatusAccepted,
			models.ApplicationStatusRejected,
		})

	if companyID > 0 {
		query = query.Where("jobs.company_id = ?", companyID)
	}

	err := query.Group("applications.job_id").Scan(&results).Error

	var stats []map[string]interface{}
	for _, r := range results {
		stats = append(stats, map[string]interface{}{
			"job_id":        r.JobID,
			"title":         r.Title,
			"applied_count": r.AppliedCount,
			"avg_days":      r.AvgDays / 24,
		})
	}

	return stats, err
}
