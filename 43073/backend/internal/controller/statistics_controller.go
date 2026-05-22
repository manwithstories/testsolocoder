package controller

import (
	"ticket-system/internal/common/response"
	"ticket-system/internal/dto"
	"ticket-system/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsController struct {
	statisticsService *service.StatisticsService
}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{
		statisticsService: service.NewStatisticsService(),
	}
}

func (c *StatisticsController) GetActivityStatistics(ctx *gin.Context) {
	var req dto.StatisticsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	stats, err := c.statisticsService.GetActivityStatistics(&req)
	if err != nil {
		response.ServerError(ctx, "获取活动统计失败")
		return
	}

	response.Success(ctx, stats)
}

func (c *StatisticsController) GetTicketTypeStatistics(ctx *gin.Context) {
	var req dto.StatisticsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	stats, err := c.statisticsService.GetTicketTypeStatistics(&req)
	if err != nil {
		response.ServerError(ctx, "获取票型统计失败")
		return
	}

	response.Success(ctx, stats)
}

func (c *StatisticsController) GetDailyStatistics(ctx *gin.Context) {
	var req dto.StatisticsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	stats, err := c.statisticsService.GetDailyStatistics(&req)
	if err != nil {
		response.ServerError(ctx, "获取每日统计失败")
		return
	}

	response.Success(ctx, stats)
}

func (c *StatisticsController) ExportExcel(ctx *gin.Context) {
	var req dto.StatisticsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	data, err := c.statisticsService.ExportExcel(&req)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	filename := "statistics_" + time.Now().Format("20060102150405") + ".xlsx"
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
