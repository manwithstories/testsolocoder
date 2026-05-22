package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
	"auction-system/internal/dto"
	"auction-system/internal/middleware"
	"auction-system/internal/services"
	"auction-system/pkg/response"
)

type StatisticsController struct {
	statsService *services.StatisticsService
}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{
		statsService: services.NewStatisticsService(),
	}
}

func (ctrl *StatisticsController) GetOverall(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	stats, err := ctrl.statsService.GetOverallStatistics(&query)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

func (ctrl *StatisticsController) GetMyStatistics(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	stats, err := ctrl.statsService.GetUserBidStatistics(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

func (ctrl *StatisticsController) ExportOrders(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	filename, err := ctrl.statsService.ExportOrdersToCSV(&query)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	defer os.Remove(filename)

	c.FileAttachment(filename, "orders.csv")
}

func (ctrl *StatisticsController) ExportBids(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	filename, err := ctrl.statsService.ExportBidsToCSV(&query)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	defer os.Remove(filename)

	c.FileAttachment(filename, "bids.csv")
}
