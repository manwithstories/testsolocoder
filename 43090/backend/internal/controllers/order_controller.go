package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"auction-system/internal/dto"
	"auction-system/internal/middleware"
	"auction-system/internal/services"
	"auction-system/pkg/response"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{
		orderService: services.NewOrderService(),
	}
}

func (ctrl *OrderController) Create(c *gin.Context) {
	buyerID := middleware.GetCurrentUserID(c)

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	order, err := ctrl.orderService.CreateOrder(buyerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, order)
}

func (ctrl *OrderController) GetByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	orderID, _ := strconv.Atoi(c.Param("id"))

	order, err := ctrl.orderService.GetOrderByID(uint(orderID), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, order)
}

func (ctrl *OrderController) GetByNo(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	orderNo := c.Param("order_no")

	order, err := ctrl.orderService.GetOrderByNo(orderNo, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, order)
}

func (ctrl *OrderController) GetBuyerOrders(c *gin.Context) {
	buyerID := middleware.GetCurrentUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var status *int
	if s, err := strconv.Atoi(c.Query("status")); err == nil {
		status = &s
	}

	orders, total, err := ctrl.orderService.GetBuyerOrders(buyerID, page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *OrderController) GetSellerOrders(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var status *int
	if s, err := strconv.Atoi(c.Query("status")); err == nil {
		status = &s
	}

	orders, total, err := ctrl.orderService.GetSellerOrders(sellerID, page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *OrderController) Pay(c *gin.Context) {
	buyerID := middleware.GetCurrentUserID(c)
	orderID, _ := strconv.Atoi(c.Param("id"))

	var req dto.PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	payment, err := ctrl.orderService.PayOrder(uint(orderID), buyerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, payment)
}

func (ctrl *OrderController) Ship(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	orderID, _ := strconv.Atoi(c.Param("id"))
	trackingNo := c.DefaultQuery("tracking_no", "")

	if err := ctrl.orderService.ShipOrder(uint(orderID), sellerID, trackingNo); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *OrderController) ConfirmDelivery(c *gin.Context) {
	buyerID := middleware.GetCurrentUserID(c)
	orderID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.orderService.ConfirmDelivery(uint(orderID), buyerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *OrderController) Complete(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.orderService.CompleteOrder(uint(orderID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *OrderController) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.DefaultQuery("keyword", "")
	var status *int
	if s, err := strconv.Atoi(c.Query("status")); err == nil {
		status = &s
	}

	orders, total, err := ctrl.orderService.GetAllOrders(page, pageSize, status, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
