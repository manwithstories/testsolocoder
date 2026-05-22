package controller

import (
	"strconv"
	"ticket-system/internal/common/response"
	"ticket-system/internal/dto"
	"ticket-system/internal/service"

	"github.com/gin-gonic/gin"
)

type CheckInController struct {
	checkInService *service.CheckInService
}

func NewCheckInController() *CheckInController {
	return &CheckInController{
		checkInService: service.NewCheckInService(),
	}
}

func (c *CheckInController) CheckIn(ctx *gin.Context) {
	var req dto.CheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	checkIn, err := c.checkInService.CheckIn(req.QrCode)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, checkIn)
}

func (c *CheckInController) GetList(ctx *gin.Context) {
	var req dto.CheckInListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	checkIns, total, err := c.checkInService.GetList(&req)
	if err != nil {
		response.ServerError(ctx, "获取签到列表失败")
		return
	}

	response.Page(ctx, checkIns, total, req.Page, req.PageSize)
}

func (c *CheckInController) GetStatistics(ctx *gin.Context) {
	activityIDStr := ctx.Query("activityId")
	var activityID uint
	if activityIDStr != "" {
		id, err := strconv.ParseUint(activityIDStr, 10, 32)
		if err == nil {
			activityID = uint(id)
		}
	}

	total, checked, err := c.checkInService.GetStatistics(activityID)
	if err != nil {
		response.ServerError(ctx, "获取签到统计失败")
		return
	}

	response.Success(ctx, gin.H{
		"total":   total,
		"checked": checked,
	})
}
