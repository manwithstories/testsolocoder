package handler

import (
	"property-management/internal/database"
	"property-management/internal/model"
	"property-management/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RepairHandler struct{}

func NewRepairHandler() *RepairHandler {
	return &RepairHandler{}
}

type RepairRequest struct {
	TenantID    uint   `json:"tenantId" binding:"required"`
	PropertyID  uint   `json:"propertyId" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category"`
	Images      string `json:"images"`
	Priority    int    `json:"priority"`
}

func (h *RepairHandler) Create(c *gin.Context) {
	var req RepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	order := model.RepairOrder{
		TenantID:    req.TenantID,
		PropertyID:  req.PropertyID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Images:      req.Images,
		Priority:    req.Priority,
		Status:      1,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		utils.Error(c, 500, "Failed to create repair order")
		return
	}

	utils.Success(c, order)
}

func (h *RepairHandler) Assign(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var order model.RepairOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.Error(c, 404, "Repair order not found")
		return
	}

	var req struct {
		HandlerID uint `json:"handlerId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&order).Updates(map[string]interface{}{
		"handler_id": req.HandlerID,
		"status":     2,
	})

	utils.Success(c, nil)
}

func (h *RepairHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var order model.RepairOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.Error(c, 404, "Repair order not found")
		return
	}

	var req struct {
		Status      int    `json:"status" binding:"required"`
		ProcessNote string `json:"processNote"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}

	if req.ProcessNote != "" {
		updates["process_note"] = req.ProcessNote
	}

	if req.Status == 3 {
		now := time.Now()
		updates["completed_at"] = now
	}

	database.DB.Model(&order).Updates(updates)
	utils.Success(c, nil)
}

func (h *RepairHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.RepairOrder{})

	if status := c.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if priority := c.Query("priority"); priority != "" {
		p, _ := strconv.Atoi(priority)
		query = query.Where("priority = ?", p)
	}

	var total int64
	query.Count(&total)

	var orders []model.RepairOrder
	offset := (page - 1) * pageSize
	query.Preload("Tenant").Preload("Property").Preload("Handler").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders)

	utils.Success(c, gin.H{
		"list":     orders,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *RepairHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var order model.RepairOrder
	if err := database.DB.Preload("Tenant").Preload("Property").Preload("Handler").
		First(&order, id).Error; err != nil {
		utils.Error(c, 404, "Repair order not found")
		return
	}
	utils.Success(c, order)
}

func (h *RepairHandler) MyOrders(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.RepairOrder{}).Where("handler_id = ?", userID)

	var total int64
	query.Count(&total)

	var orders []model.RepairOrder
	offset := (page - 1) * pageSize
	query.Preload("Tenant").Preload("Property").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders)

	utils.Success(c, gin.H{
		"list":     orders,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
