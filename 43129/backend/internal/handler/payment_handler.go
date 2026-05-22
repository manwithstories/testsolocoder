package handler

import (
	"strconv"

	"beauty-salon-system/internal/service"
	"beauty-salon-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService      *service.PaymentService
	notificationService *service.NotificationService
	auditService        *service.AuditService
	customerService     *service.CustomerService
}

func NewPaymentHandler(
	paymentService *service.PaymentService,
	notificationService *service.NotificationService,
	auditService *service.AuditService,
	customerService *service.CustomerService,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService:      paymentService,
		notificationService: notificationService,
		auditService:        auditService,
		customerService:     customerService,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req service.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.paymentService.CreatePayment(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	customer, err := h.customerService.GetByID(req.CustomerID)
	if err == nil && customer != nil {
		h.notificationService.SendPaymentNotification(customer.UserID, result)
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "create", "payment", "Create payment", c.ClientIP())

	response.Success(c, result)
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.paymentService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "payment not found")
		return
	}

	response.Success(c, result)
}

func (h *PaymentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	filters := make(map[string]interface{})
	if v := c.Query("pay_method"); v != "" {
		filters["pay_method"] = v
	}
	if v := c.Query("start_date"); v != "" {
		filters["start_date"] = v
	}
	if v := c.Query("end_date"); v != "" {
		filters["end_date"] = v
	}

	result, total, err := h.paymentService.List(page, pageSize, filters)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  result,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *PaymentHandler) CreateMemberCard(c *gin.Context) {
	var req service.CreateMemberCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.paymentService.CreateMemberCard(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *PaymentHandler) GetMemberCards(c *gin.Context) {
	customerID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.paymentService.GetMemberCards(uint(customerID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *PaymentHandler) RechargeCard(c *gin.Context) {
	cardID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.paymentService.RechargeCard(uint(cardID), req.Amount); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *PaymentHandler) PurchasePackage(c *gin.Context) {
	var req service.PurchasePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.paymentService.PurchasePackage(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *PaymentHandler) GetCustomerPackages(c *gin.Context) {
	customerID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.paymentService.GetCustomerPackages(uint(customerID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}
