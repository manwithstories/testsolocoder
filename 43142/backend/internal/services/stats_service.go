package services

import (
	"fmt"
	"time"

	"recruitment-platform/internal/repository"
	"recruitment-platform/internal/utils"
)

type StatisticsService struct {
	statsRepo       *repository.StatisticsRepository
	applicationRepo *repository.ApplicationRepository
	jobRepo         *repository.JobRepository
}

func NewStatisticsService(
	statsRepo *repository.StatisticsRepository,
	applicationRepo *repository.ApplicationRepository,
	jobRepo *repository.JobRepository,
) *StatisticsService {
	return &StatisticsService{
		statsRepo:       statsRepo,
		applicationRepo: applicationRepo,
		jobRepo:         jobRepo,
	}
}

func (s *StatisticsService) GetDailyStatistics(date time.Time) (*map[string]interface{}, error) {
	stats, err := s.statsRepo.GetDailyStats(date)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return &map[string]interface{}{
			"date":             date.Format("2006-01-02"),
			"new_jobs":         0,
			"new_applications": 0,
			"new_users":        0,
			"total_views":      0,
			"interviews":       0,
			"hires":            0,
		}, nil
	}

	return &map[string]interface{}{
		"date":             stats.Date.Format("2006-01-02"),
		"new_jobs":         stats.NewJobs,
		"new_applications": stats.NewApplications,
		"new_users":        stats.NewUsers,
		"total_views":      stats.TotalViews,
		"interviews":       stats.Interviews,
		"hires":            stats.Hires,
	}, nil
}

func (s *StatisticsService) GetDateRangeStatistics(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	statsList, err := s.statsRepo.GetDateRangeStats(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, stats := range statsList {
		result = append(result, map[string]interface{}{
			"date":             stats.Date.Format("2006-01-02"),
			"new_jobs":         stats.NewJobs,
			"new_applications": stats.NewApplications,
			"new_users":        stats.NewUsers,
			"total_views":      stats.TotalViews,
			"interviews":       stats.Interviews,
			"hires":            stats.Hires,
		})
	}

	return result, nil
}

func (s *StatisticsService) GetApplicationStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.applicationRepo.GetApplicationStats(startDate, endDate)
}

func (s *StatisticsService) GetJobStats(companyID uint, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	return s.statsRepo.GetJobStats(companyID, startDate, endDate)
}

func (s *StatisticsService) GetRecruitmentCycleStats(companyID uint) ([]map[string]interface{}, error) {
	return s.statsRepo.GetRecruitmentCycleStats(companyID)
}

func (s *StatisticsService) ExportJobStatistics(companyID uint, startDate, endDate time.Time) (string, error) {
	jobStats, err := s.statsRepo.GetJobStats(companyID, startDate, endDate)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("job_statistics_%s.xlsx", time.Now().Format("20060102_150405"))

	columns := []utils.ExcelColumn{
		{Header: "职位ID", Key: "job_id", Width: 10},
		{Header: "职位名称", Key: "title", Width: 30},
		{Header: "浏览量", Key: "total_views", Width: 15},
		{Header: "投递数", Key: "apply_count", Width: 15},
		{Header: "面试数", Key: "interview_count", Width: 15},
		{Header: "录用数", Key: "hire_count", Width: 15},
		{Header: "转化率(%)", Key: "conversion_rate", Width: 15},
	}

	return utils.ExportToExcel(utils.ExcelExportConfig{
		SheetName: "职位统计",
		Columns:   columns,
		Data:      jobStats,
		FileName:  fileName,
	})
}

func (s *StatisticsService) ExportApplicationStatistics(startDate, endDate time.Time) (string, error) {
	stats, err := s.applicationRepo.GetApplicationStats(startDate, endDate)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("application_statistics_%s.xlsx", time.Now().Format("20060102_150405"))

	columns := []utils.ExcelColumn{
		{Header: "指标", Key: "metric", Width: 20},
		{Header: "数值", Key: "value", Width: 20},
	}

	var data []map[string]interface{}
	if total, ok := stats["total"].(int64); ok {
		data = append(data, map[string]interface{}{
			"metric": "总投递数",
			"value":  total,
		})
	}

	if byStatus, ok := stats["by_status"].(map[string]int64); ok {
		statusNames := map[string]string{
			"pending":    "待处理",
			"viewed":     "已查看",
			"interested": "感兴趣",
			"interview":  "面试中",
			"accepted":   "已录用",
			"rejected":   "已拒绝",
			"withdrawn":  "已撤回",
		}
		for status, count := range byStatus {
			name := statusNames[status]
			if name == "" {
				name = status
			}
			data = append(data, map[string]interface{}{
				"metric": name + "数",
				"value":  count,
			})
		}
	}

	return utils.ExportToExcel(utils.ExcelExportConfig{
		SheetName: "投递统计",
		Columns:   columns,
		Data:      data,
		FileName:  fileName,
	})
}
