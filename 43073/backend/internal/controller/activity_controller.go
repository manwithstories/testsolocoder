package controller

import (
	"strconv"
	"ticket-system/internal/common/response"
	"ticket-system/internal/dto"
	"ticket-system/internal/service"

	"github.com/gin-gonic/gin"
)

type ActivityController struct {
	activityService *service.ActivityService
}

func NewActivityController() *ActivityController {
	return &ActivityController{
		activityService: service.NewActivityService(),
	}
}

func (c *ActivityController) Create(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var req dto.ActivityCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	activity, err := c.activityService.Create(userID, &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, activity)
}

func (c *ActivityController) GetList(ctx *gin.Context) {
	var req dto.ActivityListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	activities, total, err := c.activityService.GetList(&req)
	if err != nil {
		response.ServerError(ctx, "获取活动列表失败")
		return
	}

	response.Page(ctx, activities, total, req.Page, req.PageSize)
}

func (c *ActivityController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的活动ID")
		return
	}

	activity, err := c.activityService.GetByID(uint(id))
	if err != nil {
		response.NotFound(ctx, "活动不存在")
		return
	}

	response.Success(ctx, activity)
}

func (c *ActivityController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的活动ID")
		return
	}

	var req dto.ActivityUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	activity, err := c.activityService.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, activity)
}

func (c *ActivityController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的活动ID")
		return
	}

	var req dto.ActivityStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	activity, err := c.activityService.UpdateStatus(uint(id), req.Status)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, activity)
}

func (c *ActivityController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的活动ID")
		return
	}

	if err := c.activityService.Delete(uint(id)); err != nil {
		response.ServerError(ctx, "删除活动失败")
		return
	}

	response.Success(ctx, nil)
}
