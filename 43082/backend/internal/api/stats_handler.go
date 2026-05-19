package api

import (
	"gym-management/internal/middleware"
	"gym-management/internal/pkg/utils"
	"gym-management/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService service.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: service.NewStatsService(),
	}
}

func (h *StatsHandler) RegisterRoutes(r *gin.RouterGroup) {
	stats := r.Group("/stats")
	stats.Use(middleware.JWTAuth())
	{
		stats.GET("/dashboard", h.GetDashboardStats)
		stats.GET("/members", h.GetMemberStats)
		stats.GET("/courses", h.GetCourseStats)
		stats.GET("/coaches", h.GetCoachStats)
	}
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetMemberStats(c *gin.Context) {
	startDate, endDate := h.getDateRange(c)

	stats, err := h.statsService.GetMemberStats(startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetCourseStats(c *gin.Context) {
	startDate, endDate := h.getDateRange(c)

	stats, err := h.statsService.GetCourseStats(startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetCoachStats(c *gin.Context) {
	startDate, endDate := h.getDateRange(c)

	stats, err := h.statsService.GetCoachStats(startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) getDateRange(c *gin.Context) (time.Time, time.Time) {
	now := time.Now()
	endDate := now
	startDate := now.AddDate(0, 0, -30)

	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", startStr, time.Local); err == nil {
			startDate = t
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", endStr, time.Local); err == nil {
			endDate = t.Add(24 * time.Hour).Add(-time.Second)
		}
	}

	return startDate, endDate
}
