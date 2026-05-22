package handler

import (
	"strconv"

	"beauty-salon-system/internal/service"
	"beauty-salon-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type TechnicianHandler struct {
	technicianService *service.TechnicianService
}

func NewTechnicianHandler(technicianService *service.TechnicianService) *TechnicianHandler {
	return &TechnicianHandler{
		technicianService: technicianService,
	}
}

func (h *TechnicianHandler) Create(c *gin.Context) {
	var req service.CreateTechnicianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.technicianService.Create(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *TechnicianHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.technicianService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "technician not found")
		return
	}

	response.Success(c, result)
}

func (h *TechnicianHandler) GetMyProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	result, err := h.technicianService.GetByUserID(userID.(uint))
	if err != nil {
		response.Error(c, 404, "technician not found")
		return
	}

	response.Success(c, result)
}

func (h *TechnicianHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.UpdateTechnicianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.technicianService.Update(uint(id), &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *TechnicianHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, total, err := h.technicianService.List(page, pageSize, status)
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

func (h *TechnicianHandler) ListAll(c *gin.Context) {
	result, err := h.technicianService.ListAll()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *TechnicianHandler) AddLeave(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Date   string `json:"date" binding:"required"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.technicianService.AddLeave(uint(id), req.Date, req.Reason); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TechnicianHandler) GetLeaves(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))

	result, err := h.technicianService.GetLeaves(uint(id), month)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}
