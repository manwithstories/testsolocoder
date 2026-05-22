package handler

import (
	"car-rental/internal/middleware"
	"car-rental/internal/model"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"
	"time"

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

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	order, err := h.orderService.GetOrderByID(uint(id))
	if err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) GetOrderByNo(c *gin.Context) {
	orderNo := c.Param("no")

	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	status := c.Query("status")

	var startDate, endDate *time.Time
	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			startDate = &t
		}
	}
	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			endDate = &t
		}
	}

	orders, total, err := h.orderService.GetAllOrders(page, pageSize, uint(userID), status, startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, orders, total, page, pageSize)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	page, pageSize, _, _ := utils.ParsePageParams(c)

	orders, total, err := h.orderService.GetUserOrders(user.UserID, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, orders, total, page, pageSize)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Status model.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.orderService.UpdateOrderStatus(uint(id), req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *OrderHandler) RefundOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	err := h.orderService.RefundOrder(uint(id), req.Reason)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *OrderHandler) ExportOrders(c *gin.Context) {
	status := c.Query("status")

	var startDate, endDate *time.Time
	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			startDate = &t
		}
	}
	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			endDate = &t
		}
	}

	filePath, err := h.orderService.ExportOrders(status, startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, "导出失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"file_path": filePath})
}
