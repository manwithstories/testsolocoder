package handler

import (
	"strconv"
	"strings"

	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/service"
	resp "luxury-trading-platform/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	buyerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), buyerID.(uint), &req)
	if err != nil {
		if strings.Contains(err.Error(), "out of stock") {
			resp.Conflict(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Created(c, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), uint(id))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, order)
}

func (h *OrderHandler) GetOrderByNumber(c *gin.Context) {
	orderNumber := c.Param("order_number")

	order, err := h.orderService.GetOrderByNumber(c.Request.Context(), orderNumber)
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	role, _ := c.Get("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var buyerID, sellerID *uint
	if exists {
		uid := userID.(uint)
		if role == model.RoleBuyer {
			buyerID = &uid
		} else if role == model.RoleSeller {
			sellerID = &uid
		}
	}

	orders, total, err := h.orderService.ListOrders(page, pageSize, buyerID, sellerID, model.OrderStatus(status))
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.SuccessWithPage(c, orders, total, page, pageSize)
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	buyerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req service.PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	order, err := h.orderService.PayOrder(c.Request.Context(), uint(id), buyerID.(uint), &req)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, order)
}

func (h *OrderHandler) ShipOrder(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req service.ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	order, err := h.orderService.ShipOrder(c.Request.Context(), uint(id), sellerID.(uint), &req)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, order)
}

func (h *OrderHandler) ConfirmDelivery(c *gin.Context) {
	buyerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	order, err := h.orderService.ConfirmDelivery(c.Request.Context(), uint(id), buyerID.(uint))
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	order, err := h.orderService.CancelOrder(c.Request.Context(), uint(id), userID.(uint), req.Reason)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, order)
}
