package controller

import (
	"strconv"
	"ticket-system/internal/common/response"
	"ticket-system/internal/dto"
	"ticket-system/internal/service"

	"github.com/gin-gonic/gin"
)

type TicketTypeController struct {
	ticketTypeService *service.TicketTypeService
}

func NewTicketTypeController() *TicketTypeController {
	return &TicketTypeController{
		ticketTypeService: service.NewTicketTypeService(),
	}
}

func (c *TicketTypeController) Create(ctx *gin.Context) {
	var req dto.TicketTypeCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	ticketType, err := c.ticketTypeService.Create(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, ticketType)
}

func (c *TicketTypeController) GetList(ctx *gin.Context) {
	var req dto.TicketTypeListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	ticketTypes, total, err := c.ticketTypeService.GetList(&req)
	if err != nil {
		response.ServerError(ctx, "获取票型列表失败")
		return
	}

	response.Page(ctx, ticketTypes, total, req.Page, req.PageSize)
}

func (c *TicketTypeController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的票型ID")
		return
	}

	ticketType, err := c.ticketTypeService.GetByID(uint(id))
	if err != nil {
		response.NotFound(ctx, "票型不存在")
		return
	}

	response.Success(ctx, ticketType)
}

func (c *TicketTypeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的票型ID")
		return
	}

	var req dto.TicketTypeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	ticketType, err := c.ticketTypeService.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, ticketType)
}

func (c *TicketTypeController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的票型ID")
		return
	}

	if err := c.ticketTypeService.Delete(uint(id)); err != nil {
		response.ServerError(ctx, "删除票型失败")
		return
	}

	response.Success(ctx, nil)
}
