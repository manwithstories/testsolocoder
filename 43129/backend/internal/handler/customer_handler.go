package handler

import (
	"strconv"

	"beauty-salon-system/internal/service"
	"beauty-salon-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req service.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.customerService.Create(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.customerService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "customer not found")
		return
	}

	response.Success(c, result)
}

func (h *CustomerHandler) GetMyProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	result, err := h.customerService.GetByUserID(userID.(uint))
	if err != nil {
		response.Error(c, 404, "customer not found")
		return
	}

	response.Success(c, result)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.customerService.Update(uint(id), &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *CustomerHandler) UpdateMyProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	customer, err := h.customerService.GetByUserID(userID.(uint))
	if err != nil {
		response.Error(c, 404, "customer not found")
		return
	}

	var req service.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.customerService.Update(customer.ID, &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *CustomerHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, total, err := h.customerService.List(page, pageSize, keyword)
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
