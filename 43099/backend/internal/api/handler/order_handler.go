package handler

import (
	"net/http"
	"strconv"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/internal/service"
	"venue-booking/pkg/email"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
	logService   *service.OperationLogService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(),
		logService:   service.NewOperationLogService(),
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	userID, _ := c.Get("userID")
	order, err := h.orderService.Create(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, userID.(uint), "create_order", "order", map[string]interface{}{
		"order_id":   order.ID,
		"order_no":   order.OrderNo,
		"type":       order.Type,
		"item_id":    order.ItemID,
		"item_name":  order.ItemName,
		"total_amount": order.TotalAmount,
	})

	orderUser, _ := order.User, nil
	if orderUser != nil {
		emailService := email.NewEmailService()
		emailService.SendOrderConfirmation(orderUser.Email, order.OrderNo, order.ItemName)
	}

	c.JSON(http.StatusOK, dto.Success(order))
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid order ID"))
		return
	}

	order, err := h.orderService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "Order not found"))
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")
	if order.UserID != userID.(uint) && userRole != model.RoleAdmin && userRole != model.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, dto.Error(403, "You can only view your own orders"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(order))
}

func (h *OrderHandler) List(c *gin.Context) {
	var req dto.OrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	var filterUserID *uint
	if userRole != model.RoleAdmin && userRole != model.RoleSuperAdmin {
		uid := userID.(uint)
		filterUserID = &uid
	}

	orders, total, err := h.orderService.List(&req, filterUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get orders"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     orders,
	}))
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid order ID"))
		return
	}

	var req dto.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	userID, _ := c.Get("userID")
	err = h.orderService.Cancel(uint(id), userID.(uint), req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, userID.(uint), "cancel_order", "order", map[string]interface{}{
		"order_id": id,
		"reason":   req.Reason,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *OrderHandler) Confirm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid order ID"))
		return
	}

	var req dto.ReviewOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	userID, _ := c.Get("userID")
	err = h.orderService.Confirm(uint(id), userID.(uint), req.Note)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, userID.(uint), "confirm_order", "order", map[string]interface{}{
		"order_id": id,
		"note":     req.Note,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *OrderHandler) Complete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid order ID"))
		return
	}

	userID, _ := c.Get("userID")
	err = h.orderService.Complete(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, userID.(uint), "complete_order", "order", map[string]interface{}{
		"order_id": id,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *OrderHandler) GetCalendar(c *gin.Context) {
	var req dto.CalendarRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	events, err := h.orderService.GetCalendar(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get calendar data"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(events))
}
