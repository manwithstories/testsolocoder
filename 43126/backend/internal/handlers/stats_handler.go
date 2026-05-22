package handlers

import (
	"time"

	"meeting-room/internal/services"
	"meeting-room/internal/utils"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: services.NewStatsService(),
	}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	department := c.Query("department")

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		utils.BadRequest(c, "开始日期格式错误")
		return
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		utils.BadRequest(c, "结束日期格式错误")
		return
	}
	endTime = endTime.AddDate(0, 0, 1)

	stats, err := h.statsService.GetStats(startTime, endTime, department)
	if err != nil {
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) ExportStats(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	department := c.Query("department")

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		utils.BadRequest(c, "开始日期格式错误")
		return
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		utils.BadRequest(c, "结束日期格式错误")
		return
	}
	endTime = endTime.AddDate(0, 0, 1)

	filename, err := h.statsService.ExportToExcel(startTime, endTime, department)
	if err != nil {
		utils.InternalError(c, "导出报表失败")
		return
	}

	c.FileAttachment("./uploads/"+filename, filename)
}
