package handler

import (
	"strconv"

	"luxury-trading-platform/internal/service"
	resp "luxury-trading-platform/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type StatisticHandler struct {
	statService *service.StatisticService
}

func NewStatisticHandler(statService *service.StatisticService) *StatisticHandler {
	return &StatisticHandler{statService: statService}
}

func (h *StatisticHandler) GetDashboardStats(c *gin.Context) {
	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	stats, err := h.statService.GetDashboardStats(c.Request.Context(), days)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, stats)
}

func (h *StatisticHandler) InvalidateCache(c *gin.Context) {
	h.statService.InvalidateCache(c.Request.Context())
	resp.Success(c, gin.H{"message": "cache invalidated successfully"})
}
