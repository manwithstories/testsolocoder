package handler

import (
	"strconv"

	"beauty-salon-system/internal/service"
	"beauty-salon-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type ServiceItemHandler struct {
	serviceItemService *service.ServiceItemService
}

func NewServiceItemHandler(serviceItemService *service.ServiceItemService) *ServiceItemHandler {
	return &ServiceItemHandler{
		serviceItemService: serviceItemService,
	}
}

func (h *ServiceItemHandler) Create(c *gin.Context) {
	var req service.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.serviceItemService.Create(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ServiceItemHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.serviceItemService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "service not found")
		return
	}

	response.Success(c, result)
}

func (h *ServiceItemHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.serviceItemService.Update(uint(id), &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ServiceItemHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.serviceItemService.Delete(uint(id)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ServiceItemHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	isPackage, _ := strconv.ParseBool(c.DefaultQuery("is_package", "false"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, total, err := h.serviceItemService.List(page, pageSize, category, isPackage)
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

func (h *ServiceItemHandler) ListAll(c *gin.Context) {
	result, err := h.serviceItemService.ListAll()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ServiceItemHandler) AddPackageService(c *gin.Context) {
	var req service.AddPackageServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.serviceItemService.AddPackageService(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ServiceItemHandler) GetPackageServices(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.serviceItemService.GetPackageServices(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ServiceItemHandler) DeletePackageServices(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.serviceItemService.DeletePackageServices(uint(id)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}
