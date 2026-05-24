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

type ServiceHandler struct {
	serviceService *service.ServiceService
}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{
		serviceService: service.NewServiceService(),
	}
}

func (h *ServiceHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.CreateServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	svc, err := h.serviceService.Create(userID, &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, svc)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	svc, err := h.serviceService.GetByID(uint(id))
	if err != nil {
		response.ErrNotFound(c, "服务需求不存在")
		return
	}
	response.Success(c, svc)
}

func (h *ServiceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := model.ServiceStatus(c.Query("status"))
	region := c.Query("region")
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	pilotID, _ := strconv.ParseUint(c.Query("pilot_id"), 10, 64)
	services, total, err := h.serviceService.List(page, pageSize, uint(userID), uint(pilotID), status, region)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, services, total, page, pageSize)
}

func (h *ServiceHandler) MyServices(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := model.ServiceStatus(c.Query("status"))
	services, total, err := h.serviceService.List(page, pageSize, userID, 0, status, "")
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, services, total, page, pageSize)
}

func (h *ServiceHandler) CreateBid(c *gin.Context) {
	pilotID := middleware.GetUserID(c)
	var req dto.CreateBidReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.serviceService.CreateBid(pilotID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *ServiceHandler) ListBids(c *gin.Context) {
	serviceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	bids, err := h.serviceService.ListBids(uint(serviceID))
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Success(c, bids)
}

func (h *ServiceHandler) AcceptBid(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.AcceptBidReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.serviceService.AcceptBid(userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *ServiceHandler) UpdateStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdateServiceStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.serviceService.UpdateStatus(userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}
