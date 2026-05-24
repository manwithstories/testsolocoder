package handler

import (
	"drone-rental/internal/dto"
	"drone-rental/internal/middleware"
	"drone-rental/internal/model"
	"drone-rental/internal/pkg/response"
	"drone-rental/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(),
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	order, err := h.orderService.Create(userID, &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	order, err := h.orderService.GetByID(uint(id))
	if err != nil {
		response.ErrNotFound(c, "订单不存在")
		return
	}
	response.Success(c, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	droneID, _ := strconv.ParseUint(c.Query("drone_id"), 10, 64)
	status := model.OrderStatus(c.Query("status"))
	orders, total, err := h.orderService.List(page, pageSize, 0, uint(droneID), status)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, orders, total, page, pageSize)
}

func (h *OrderHandler) MyOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := model.OrderStatus(c.Query("status"))
	orders, total, err := h.orderService.List(page, pageSize, userID, 0, status)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, orders, total, page, pageSize)
}

func (h *OrderHandler) Pay(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.PayOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.orderService.PayOrder(userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.CancelOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.orderService.CancelOrder(userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) Pickup(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.orderService.PickupOrder(userID, uint(id)); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) ConfirmReturn(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.ConfirmReturnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.orderService.ConfirmReturn(userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) Complete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.orderService.CompleteOrder(uint(id)); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}
