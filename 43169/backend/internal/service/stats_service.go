package service

import (
	"bytes"
	"fmt"
	"time"

	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/repository"
	"matchmaking-platform/internal/utils"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type StatsService struct {
	userRepo   *repository.UserRepo
	dateRepo   *repository.DateRepo
	matchRepo  *repository.MatchRepo
	chatRepo   *repository.ChatRepo
	logRepo    *repository.SystemLogRepo
}

func NewStatsService() *StatsService {
	return &StatsService{
		userRepo:  repository.NewUserRepo(),
		dateRepo:  repository.NewDateRepo(),
		matchRepo: repository.NewMatchRepo(),
		chatRepo:  repository.NewChatRepo(),
		logRepo:   repository.NewSystemLogRepo(),
	}
}

type PlatformStats struct {
	TotalUsers       int64   `json:"total_users"`
	ActiveToday      int64   `json:"active_today"`
	VerifiedUsers    int64   `json:"verified_users"`
	TotalMatchmakers int64   `json:"total_matchmakers"`
	TotalDates       int64   `json:"total_dates"`
	CompletedDates   int64   `json:"completed_dates"`
	MatchSuccessRate float64 `json:"match_success_rate"`
	TotalMessages    int64   `json:"total_messages"`
	NewUsersToday    int64   `json:"new_users_today"`
}

func (s *StatsService) GetPlatformStats() (*PlatformStats, error) {
	totalUsers, _ := s.userRepo.CountTotal()
	activeToday, _ := s.userRepo.CountActiveToday()
	verifiedUsers, _ := s.userRepo.CountByVerifyStatus("verified")
	totalMatchmakers, _ := s.userRepo.CountByRole("matchmaker")
	totalDates, _ := s.dateRepo.CountTotal()
	completedDates, _ := s.dateRepo.CountByStatus("completed")

	var newUsersToday int64
	utils.DB.Model(&model.User{}).Where("DATE(created_at) = CURDATE()").Count(&newUsersToday)

	var totalMessages int64
	utils.DB.Model(&model.ChatMessage{}).Count(&totalMessages)

	matchRate := 0.0
	if totalDates > 0 {
		matchRate = float64(completedDates) / float64(totalDates) * 100
	}

	return &PlatformStats{
		TotalUsers:       totalUsers,
		ActiveToday:      activeToday,
		VerifiedUsers:    verifiedUsers,
		TotalMatchmakers: totalMatchmakers,
		TotalDates:       totalDates,
		CompletedDates:   completedDates,
		MatchSuccessRate: matchRate,
		TotalMessages:    totalMessages,
		NewUsersToday:    newUsersToday,
	}, nil
}

type DailyStats struct {
	Date          string `json:"date"`
	NewUsers      int64  `json:"new_users"`
	ActiveUsers   int64  `json:"active_users"`
	DatesCreated  int64  `json:"dates_created"`
	DatesCompleted int64 `json:"dates_completed"`
	MessagesCount int64  `json:"messages_count"`
}

func (s *StatsService) GetDailyStats(startDate, endDate time.Time) ([]DailyStats, error) {
	var stats []DailyStats

	type row struct {
		Date          string
		NewUsers      int64
		ActiveUsers   int64
		DatesCreated  int64
		DatesCompleted int64
		MessagesCount int64
	}

	var rows []row
	utils.DB.Raw(`
		SELECT 
			DATE(d.date) as date,
			COALESCE(new_users, 0) as new_users,
			COALESCE(active_users, 0) as active_users,
			COALESCE(dates_created, 0) as dates_created,
			COALESCE(dates_completed, 0) as dates_completed,
			COALESCE(messages_count, 0) as messages_count
		FROM (
			SELECT CURDATE() - INTERVAL n DAY as date
			FROM (SELECT 0 n UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 
				  UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9
				  UNION SELECT 10 UNION SELECT 11 UNION SELECT 12 UNION SELECT 13 UNION SELECT 14) days
		) d
		LEFT JOIN (
			SELECT DATE(created_at) as dt, COUNT(*) as new_users FROM users WHERE DATE(created_at) BETWEEN ? AND ? GROUP BY DATE(created_at)
		) u ON d.date = u.dt
		LEFT JOIN (
			SELECT DATE(last_login_at) as dt, COUNT(DISTINCT id) as active_users FROM users WHERE DATE(last_login_at) BETWEEN ? AND ? GROUP BY DATE(last_login_at)
		) a ON d.date = a.dt
		LEFT JOIN (
			SELECT DATE(created_at) as dt, COUNT(*) as dates_created FROM date_records WHERE DATE(created_at) BETWEEN ? AND ? GROUP BY DATE(created_at)
		) dc ON d.date = dc.dt
		LEFT JOIN (
			SELECT DATE(updated_at) as dt, COUNT(*) as dates_completed FROM date_records WHERE status = 'completed' AND DATE(updated_at) BETWEEN ? AND ? GROUP BY DATE(updated_at)
		) dco ON d.date = dco.dt
		LEFT JOIN (
			SELECT DATE(created_at) as dt, COUNT(*) as messages_count FROM chat_messages WHERE DATE(created_at) BETWEEN ? AND ? GROUP BY DATE(created_at)
		) m ON d.date = m.dt
		WHERE d.date BETWEEN ? AND ?
		ORDER BY d.date
	`, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate).Scan(&rows)

	for _, r := range rows {
		stats = append(stats, DailyStats{
			Date:           r.Date,
			NewUsers:       r.NewUsers,
			ActiveUsers:    r.ActiveUsers,
			DatesCreated:   r.DatesCreated,
			DatesCompleted: r.DatesCompleted,
			MessagesCount:  r.MessagesCount,
		})
	}

	return stats, nil
}

func (s *StatsService) GetMatchmakerStats(startDate, endDate time.Time, matchmakerID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db := utils.DB.Table("matchmaker_stats ms").
		Joins("LEFT JOIN users u ON u.id = ms.matchmaker_id")

	if matchmakerID > 0 {
		db = db.Where("ms.matchmaker_id = ?", matchmakerID)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		db = db.Where("ms.updated_at BETWEEN ? AND ?", startDate, endDate)
	}

	db.Select("u.username as matchmaker_name, u.avatar, ms.*").
		Order("ms.avg_rating DESC").
		Scan(&results)

	return results, nil
}

func (s *StatsService) ExportExcel(startDate, endDate time.Time) ([]byte, error) {
	f := excelize.NewFile()

	stats, _ := s.GetDailyStats(startDate, endDate)
	sheet1 := "每日数据"
	index1, _ := f.NewSheet(sheet1)

	headers := []string{"日期", "新增用户", "活跃用户", "新建约会", "完成约会", "消息数"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet1, cell, h)
	}

	for i, s := range stats {
		row := i + 2
		f.SetCellValue(sheet1, fmt.Sprintf("A%d", row), s.Date)
		f.SetCellValue(sheet1, fmt.Sprintf("B%d", row), s.NewUsers)
		f.SetCellValue(sheet1, fmt.Sprintf("C%d", row), s.ActiveUsers)
		f.SetCellValue(sheet1, fmt.Sprintf("D%d", row), s.DatesCreated)
		f.SetCellValue(sheet1, fmt.Sprintf("E%d", row), s.DatesCompleted)
		f.SetCellValue(sheet1, fmt.Sprintf("F%d", row), s.MessagesCount)
	}

	f.SetActiveSheet(index1)

	sheet2 := "平台概览"
	index2, _ := f.NewSheet(sheet2)
	platformStats, _ := s.GetPlatformStats()

	overviewHeaders := []string{"指标", "数值"}
	for i, h := range overviewHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet2, cell, h)
	}

	overviewData := [][]interface{}{
		{"总用户数", platformStats.TotalUsers},
		{"今日活跃", platformStats.ActiveToday},
		{"已认证用户", platformStats.VerifiedUsers},
		{"红娘总数", platformStats.TotalMatchmakers},
		{"约会总数", platformStats.TotalDates},
		{"完成约会", platformStats.CompletedDates},
		{"匹配成功率(%)", fmt.Sprintf("%.2f", platformStats.MatchSuccessRate)},
		{"消息总数", platformStats.TotalMessages},
		{"今日新增", platformStats.NewUsersToday},
	}

	for i, row := range overviewData {
		r := i + 2
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", r), row[0])
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", r), row[1])
	}

	f.SetActiveSheet(index2)

	sheet3 := "红娘业绩"
	index3, _ := f.NewSheet(sheet3)
	matchmakerHeaders := []string{"红娘", "会员数", "服务数", "约会数", "成功率", "平均分"}
	for i, h := range matchmakerHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet3, cell, h)
	}

	matchmakerStats, _ := s.GetMatchmakerStats(startDate, endDate, 0)
	for i, m := range matchmakerStats {
		row := i + 2
		f.SetCellValue(sheet3, fmt.Sprintf("A%d", row), m["matchmaker_name"])
		f.SetCellValue(sheet3, fmt.Sprintf("B%d", row), m["total_members"])
		f.SetCellValue(sheet3, fmt.Sprintf("C%d", row), m["total_services"])
		f.SetCellValue(sheet3, fmt.Sprintf("D%d", row), m["total_dates"])
		if totalDates, ok := m["total_dates"].(int64); ok && totalDates > 0 {
			successRate := float64(m["success_dates"].(int64)) / float64(totalDates) * 100
			f.SetCellValue(sheet3, fmt.Sprintf("E%d", row), fmt.Sprintf("%.2f%%", successRate))
		}
		f.SetCellValue(sheet3, fmt.Sprintf("F%d", row), m["avg_rating"])
	}

	f.SetActiveSheet(index3)
	f.DeleteSheet("Sheet1")

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *StatsService) ExportPDF(startDate, endDate time.Time) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Matchmaking Platform Report")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 10)
	dateRange := fmt.Sprintf("Date Range: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	pdf.Cell(40, 8, dateRange)
	pdf.Ln(12)

	platformStats, _ := s.GetPlatformStats()

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Platform Overview")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)

	overviewData := [][]string{
		{"Total Users", fmt.Sprintf("%d", platformStats.TotalUsers)},
		{"Active Today", fmt.Sprintf("%d", platformStats.ActiveToday)},
		{"Verified Users", fmt.Sprintf("%d", platformStats.VerifiedUsers)},
		{"Total Matchmakers", fmt.Sprintf("%d", platformStats.TotalMatchmakers)},
		{"Total Dates", fmt.Sprintf("%d", platformStats.TotalDates)},
		{"Completed Dates", fmt.Sprintf("%d", platformStats.CompletedDates)},
		{"Success Rate", fmt.Sprintf("%.2f%%", platformStats.MatchSuccessRate)},
		{"Total Messages", fmt.Sprintf("%d", platformStats.TotalMessages)},
		{"New Users Today", fmt.Sprintf("%d", platformStats.NewUsersToday)},
	}

	for _, row := range overviewData {
		pdf.Cell(50, 7, row[0])
		pdf.Cell(30, 7, row[1])
		pdf.Ln(7)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Daily Statistics")
	pdf.Ln(8)

	stats, _ := s.GetDailyStats(startDate, endDate)
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"Date", "New Users", "Active", "Dates", "Completed", "Messages"}
	colWidths := []float64{28, 22, 18, 18, 22, 22}
	for i, h := range headers {
		pdf.Cell(colWidths[i], 7, h)
	}
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 9)
	for _, s := range stats {
		pdf.Cell(colWidths[0], 6, s.Date)
		pdf.Cell(colWidths[1], 6, fmt.Sprintf("%d", s.NewUsers))
		pdf.Cell(colWidths[2], 6, fmt.Sprintf("%d", s.ActiveUsers))
		pdf.Cell(colWidths[3], 6, fmt.Sprintf("%d", s.DatesCreated))
		pdf.Cell(colWidths[4], 6, fmt.Sprintf("%d", s.DatesCompleted))
		pdf.Cell(colWidths[5], 6, fmt.Sprintf("%d", s.MessagesCount))
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *StatsService) GetSystemLogs(page, pageSize int, filter map[string]interface{}) ([]model.SystemLog, int64, error) {
	return s.logRepo.List(page, pageSize, filter)
}
