package handlers

import (
	"net/http"
	"strconv"

	"print3d-platform/internal/middleware"
	"print3d-platform/internal/models"
	"print3d-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) EstimatePrice(c *gin.Context) {
	var req service.PriceEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	estimate, err := h.orderService.EstimatePrice(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, estimate)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrderByNo(c *gin.Context) {
	orderNo := c.Param("order_no")

	order, err := h.orderService.GetOrderByNo(c.Request.Context(), orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) ListCustomerOrders(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	orders, total, err := h.orderService.ListCustomerOrders(c.Request.Context(), authUser.UserID, page, pageSize, statusPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  orders,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *OrderHandler) ListPrinterOrders(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only printers can view printer orders"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	orders, total, err := h.orderService.ListPrinterOrders(c.Request.Context(), authUser.UserID, page, pageSize, statusPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  orders,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *OrderHandler) AssignPrinter(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only printers can accept orders"})
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderService.AssignPrinter(c.Request.Context(), orderID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order assigned successfully"})
}

func (h *OrderHandler) StartPrinting(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderService.StartPrinting(c.Request.Context(), orderID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Printing started"})
}

func (h *OrderHandler) CompletePrinting(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderService.CompletePrinting(c.Request.Context(), orderID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Printing completed"})
}

func (h *OrderHandler) ApproveQuality(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderService.ApproveQuality(c.Request.Context(), orderID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quality check passed"})
}

func (h *OrderHandler) ShipOrder(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only printers can ship orders"})
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		TrackingNumber string `json:"tracking_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.orderService.ShipOrder(c.Request.Context(), orderID, req.TrackingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order shipped"})
}

func (h *OrderHandler) DeliverOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderService.DeliverOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order delivered"})
}

func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderService.CompleteOrder(c.Request.Context(), orderID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order completed"})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.orderService.CancelOrder(c.Request.Context(), orderID, authUser.UserID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled"})
}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	history, err := h.orderService.GetOrderHistory(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *OrderHandler) GetPendingOrders(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter && authUser.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	orders, err := h.orderService.GetPendingOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
