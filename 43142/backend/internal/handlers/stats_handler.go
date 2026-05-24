package handlers

import (
	"time"

	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/services"
	"recruitment-platform/internal/utils"
	"recruitment-platform/pkg/logger"
)

type StatsHandler struct {
	statsService *services.StatisticsService
}

func NewStatsHandler(statsService *services.StatisticsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetDailyStats(c *gin.Context) {
	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			utils.BadRequest(c, "日期格式错误，应为YYYY-MM-DD")
			return
		}
		date = parsed
	} else {
		date = time.Now()
	}

	stats, err := h.statsService.GetDailyStatistics(date)
	if err != nil {
		logger.Error("获取日统计失败: %v", err)
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetDateRangeStats(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.BadRequest(c, "请提供开始日期和结束日期")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequest(c, "开始日期格式错误")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequest(c, "结束日期格式错误")
		return
	}

	stats, err := h.statsService.GetDateRangeStatistics(startDate, endDate)
	if err != nil {
		logger.Error("获取日期范围统计失败: %v", err)
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetApplicationStats(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	stats, err := h.statsService.GetApplicationStats(startDate, endDate)
	if err != nil {
		logger.Error("获取投递统计失败: %v", err)
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetJobStats(c *gin.Context) {
	companyID := c.GetUint("user_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	stats, err := h.statsService.GetJobStats(companyID, startDate, endDate)
	if err != nil {
		logger.Error("获取职位统计失败: %v", err)
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetRecruitmentCycleStats(c *gin.Context) {
	companyID := c.GetUint("user_id")

	stats, err := h.statsService.GetRecruitmentCycleStats(companyID)
	if err != nil {
		logger.Error("获取招聘周期统计失败: %v", err)
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) ExportJobStats(c *gin.Context) {
	companyID := c.GetUint("user_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	fileName, err := h.statsService.ExportJobStatistics(companyID, startDate, endDate)
	if err != nil {
		logger.Error("导出职位统计失败: %v", err)
		utils.InternalError(c, "导出失败")
		return
	}

	c.File(fileName)
}

func (h *StatsHandler) ExportApplicationStats(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	fileName, err := h.statsService.ExportApplicationStatistics(startDate, endDate)
	if err != nil {
		logger.Error("导出投递统计失败: %v", err)
		utils.InternalError(c, "导出失败")
		return
	}

	c.File(fileName)
}
