package handlers

import (
	"bytes"
	"fmt"
	"time"
	"translation-platform/internal/database"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/tealeg/xlsx/v3"
)

func GetProjectStatistics(c *gin.Context) {
	var totalProjects, pendingProjects, inProgressProjects, completedProjects int64
	var totalWords, completedWords int64

	database.DB.Model(&models.Project{}).Count(&totalProjects)
	database.DB.Model(&models.Project{}).Where("status = ?", models.ProjectStatusPending).Count(&pendingProjects)
	database.DB.Model(&models.Project{}).Where("status = ?", models.ProjectStatusInProgress).Count(&inProgressProjects)
	database.DB.Model(&models.Project{}).Where("status = ?", models.ProjectStatusCompleted).Count(&completedProjects)

	var projects []models.Project
	database.DB.Find(&projects)

	for _, p := range projects {
		totalWords += int64(p.WordCount)
		if p.Status == models.ProjectStatusCompleted {
			completedWords += int64(p.WordCount)
		}
	}

	completionRate := 0.0
	if totalProjects > 0 {
		completionRate = float64(completedProjects) / float64(totalProjects) * 100
	}

	utils.Success(c, gin.H{
		"total_projects":      totalProjects,
		"pending_projects":    pendingProjects,
		"in_progress_projects": inProgressProjects,
		"completed_projects":  completedProjects,
		"total_words":         totalWords,
		"completed_words":     completedWords,
		"completion_rate":     fmt.Sprintf("%.1f%%", completionRate),
	})
}

func GetTranslatorStatistics(c *gin.Context) {
	var translators []models.User
	database.DB.Preload("LanguagePairs").Preload("ExpertiseTags").
		Where("role = ?", models.RoleTranslator).Find(&translators)

	type translatorStat struct {
		User           models.User `json:"user"`
		CompletedCount int         `json:"completed_count"`
		TotalWords     int         `json:"total_words"`
		AvgRating      float64     `json:"avg_rating"`
		Efficiency     string      `json:"efficiency"`
	}

	var stats []translatorStat
	for _, t := range translators {
		var completedProjects []models.Project
		database.DB.Where("translator_id = ? AND status = ?", t.ID, models.ProjectStatusCompleted).
			Find(&completedProjects)

		totalWords := 0
		for _, p := range completedProjects {
			totalWords += p.WordCount
		}

		stats = append(stats, translatorStat{
			User:           t,
			CompletedCount: len(completedProjects),
			TotalWords:     totalWords,
			AvgRating:      t.Rating,
			Efficiency:     fmt.Sprintf("%d词/日", t.DailyCapacity),
		})
	}

	utils.Success(c, stats)
}

func GetRevenueTrend(c *gin.Context) {
	months := 12
	if m := c.Query("months"); m != "" {
		if n, err := parseInt(m); err == nil && n > 0 {
			months = n
		}
	}

	type monthlyRevenue struct {
		Month   string  `json:"month"`
		Revenue float64 `json:"revenue"`
		Count   int64   `json:"count"`
	}

	var trend []monthlyRevenue
	now := time.Now()

	for i := months - 1; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0)

		var payments []models.Payment
		database.DB.Where("created_at >= ? AND created_at < ? AND status = ?",
			startOfMonth, endOfMonth, "paid").Find(&payments)

		var revenue float64
		for _, p := range payments {
			revenue += p.Amount
		}

		trend = append(trend, monthlyRevenue{
			Month:   date.Format("2006-01"),
			Revenue: revenue,
			Count:   int64(len(payments)),
		})
	}

	utils.Success(c, trend)
}

func ExportStatisticsExcel(c *gin.Context) {
	wb := xlsx.NewFile()
	sheet, _ := wb.AddSheet("项目统计")

	headers := []string{"项目名称", "客户", "译者", "语言对", "字数", "金额", "状态", "创建时间"}
	for i, h := range headers {
		cell := sheet.Cell(0, i)
		cell.Value = h
		cell.GetStyle().Font.Bold = true
	}

	var projects []models.Project
	database.DB.Preload("Client").Preload("Translator").Find(&projects)

	for i, p := range projects {
		row := i + 1
		sheet.Cell(row, 0).Value = p.Title
		if p.Client.ID > 0 {
			sheet.Cell(row, 1).Value = p.Client.Username
		}
		if p.Translator != nil {
			sheet.Cell(row, 2).Value = p.Translator.Username
		}
		sheet.Cell(row, 3).Value = p.SourceLang + "-" + p.TargetLang
		sheet.Cell(row, 4).Value = fmt.Sprintf("%d", p.WordCount)
		sheet.Cell(row, 5).Value = fmt.Sprintf("%.2f", p.TotalAmount)
		sheet.Cell(row, 6).Value = string(p.Status)
		sheet.Cell(row, 7).Value = p.CreatedAt.Format("2006-01-02")
	}

	var buf bytes.Buffer
	wb.Write(&buf)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=statistics_%s.xlsx",
		time.Now().Format("20060102")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func ExportStatisticsPDF(c *gin.Context) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	err := pdf.AddTTFFont("default", "/System/Library/Fonts/PingFang.ttc")
	if err != nil {
		_ = err
	}

	pdf.SetFont("default", "", 14)

	pdf.SetX(40)
	pdf.SetY(40)
	pdf.Cell(nil, "翻译服务平台数据统计报表")
	pdf.Br(20)

	pdf.SetFont("default", "", 10)
	pdf.SetX(40)
	pdf.Cell(nil, fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Br(30)

	var totalProjects, completedProjects int64
	var totalRevenue float64

	database.DB.Model(&models.Project{}).Count(&totalProjects)
	database.DB.Model(&models.Project{}).Where("status = ?", models.ProjectStatusCompleted).Count(&completedProjects)

	var payments []models.Payment
	database.DB.Where("status = ?", "paid").Find(&payments)
	for _, p := range payments {
		totalRevenue += p.Amount
	}

	stats := []struct {
		Label string
		Value string
	}{
		{"项目总数", fmt.Sprintf("%d", totalProjects)},
		{"已完成项目", fmt.Sprintf("%d", completedProjects)},
		{"总收入", fmt.Sprintf("%.2f", totalRevenue)},
		{"项目完成率", fmt.Sprintf("%.1f%%", float64(completedProjects)/float64(totalProjects)*100)},
	}

	for _, s := range stats {
		pdf.SetX(40)
		pdf.Cell(nil, fmt.Sprintf("%s: %s", s.Label, s.Value))
		pdf.Br(15)
	}

	var buf bytes.Buffer
	pdf.Write(&buf)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=statistics_%s.pdf",
		time.Now().Format("20060102")))
	c.Data(200, "application/pdf", buf.Bytes())
}

func GetLanguagePairStatistics(c *gin.Context) {
	type pairStat struct {
		LanguagePair string `json:"language_pair"`
		ProjectCount int64  `json:"project_count"`
		TotalWords   int    `json:"total_words"`
		TotalRevenue float64 `json:"total_revenue"`
	}

	var projects []models.Project
	database.DB.Find(&projects)

	pairMap := make(map[string]*pairStat)

	for _, p := range projects {
		key := p.SourceLang + "-" + p.TargetLang
		if _, ok := pairMap[key]; !ok {
			pairMap[key] = &pairStat{LanguagePair: key}
		}
		pairMap[key].ProjectCount++
		pairMap[key].TotalWords += p.WordCount
		pairMap[key].TotalRevenue += p.TotalAmount
	}

	var result []pairStat
	for _, v := range pairMap {
		result = append(result, *v)
	}

	utils.Success(c, result)
}

func ListOperationLogs(c *gin.Context) {
	var logs []models.OperationLog
	query := database.DB.Preload("User")

	if module := c.Query("module"); module != "" {
		query = query.Where("module = ?", module)
	}
	if action := c.Query("action"); action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}

	var total int64
	query.Model(&models.OperationLog{}).Count(&total)

	page, pageSize := parsePagination(c)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	utils.Success(c, utils.PageResult{
		List:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
