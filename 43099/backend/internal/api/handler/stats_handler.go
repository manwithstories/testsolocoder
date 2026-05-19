package handler

import (
	"net/http"
	"venue-booking/internal/dto"
	"venue-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: service.NewStatsService(),
	}
}

func (h *StatsHandler) GetOverview(c *gin.Context) {
	overview, err := h.statsService.GetOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get overview stats"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(overview))
}

func (h *StatsHandler) GetBookingStats(c *gin.Context) {
	var req dto.DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	stats, err := h.statsService.GetBookingStats(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get booking stats"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(stats))
}

func (h *StatsHandler) GetRevenueStats(c *gin.Context) {
	var req dto.DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	stats, err := h.statsService.GetRevenueStats(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get revenue stats"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(stats))
}

func (h *StatsHandler) GetPopularVenues(c *gin.Context) {
	var req dto.DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	venues, err := h.statsService.GetPopularVenues(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get popular venues"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(venues))
}
