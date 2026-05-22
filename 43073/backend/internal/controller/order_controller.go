package controller

import (
	"strconv"
	"ticket-system/internal/common/response"
	"ticket-system/internal/dto"
	"ticket-system/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{
		orderService: service.NewOrderService(),
	}
}

func (c *OrderController) Create(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var req dto.OrderCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	order, err := c.orderService.Create(userID, &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, order)
}

func (c *OrderController) GetList(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	role := ctx.GetString("role")
	isAdmin := role == "admin"

	var req dto.OrderListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	orders, total, err := c.orderService.GetList(&req, userID, isAdmin)
	if err != nil {
		response.ServerError(ctx, "获取订单列表失败")
		return
	}

	response.Page(ctx, orders, total, req.Page, req.PageSize)
}

func (c *OrderController) Get(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	role := ctx.GetString("role")
	isAdmin := role == "admin"

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的订单ID")
		return
	}

	order, err := c.orderService.GetByID(uint(id), userID, isAdmin)
	if err != nil {
		response.NotFound(ctx, "订单不存在")
		return
	}

	response.Success(ctx, order)
}

func (c *OrderController) Pay(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	role := ctx.GetString("role")
	isAdmin := role == "admin"

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的订单ID")
		return
	}

	order, err := c.orderService.Pay(uint(id), userID, isAdmin)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, order)
}

func (c *OrderController) Cancel(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	role := ctx.GetString("role")
	isAdmin := role == "admin"

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的订单ID")
		return
	}

	order, err := c.orderService.Cancel(uint(id), userID, isAdmin)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, order)
}
