package handler

import (
	"net/http"
	"survey-platform/internal/dto"
	"survey-platform/internal/service"
	"survey-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statisticsService *service.StatisticsService
}

func NewStatisticsHandler(statisticsService *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statisticsService: statisticsService}
}

func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.statisticsService.GetStatistics(&query)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *StatisticsHandler) CrossAnalysis(c *gin.Context) {
	var query dto.CrossAnalysisQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.statisticsService.CrossAnalysis(&query)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, result)
}
