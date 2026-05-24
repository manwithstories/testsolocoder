package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"repair-platform/internal/models"
	"repair-platform/internal/utils"
	"repair-platform/pkg/logger"
)

type PartHandler struct{}

func NewPartHandler() *PartHandler {
	return &PartHandler{}
}

type CreatePartRequest struct {
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock"`
	MinStock    int     `json:"min_stock"`
	Image       string  `json:"image"`
}

func (h *PartHandler) CreatePart(c *gin.Context) {
	var req CreatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	part := models.Part{
		Name:        req.Name,
		Code:        req.Code,
		Category:    req.Category,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		Image:       req.Image,
		Status:      true,
	}

	if err := models.DB.Create(&part).Error; err != nil {
		logger.Errorf("Create part error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "创建失败")
		return
	}

	utils.Success(c, part)
}

func (h *PartHandler) GetPartList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	lowStock := c.Query("low_stock")

	var parts []models.Part
	query := models.DB.Where("status = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if lowStock == "true" {
		query = query.Where("stock <= min_stock")
	}

	var total int64
	query.Model(&models.Part{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&parts)

	var lowStockCount int64
	models.DB.Model(&models.Part{}).Where("status = ? AND stock <= min_stock", true).Count(&lowStockCount)

	utils.Success(c, gin.H{
		"list":          parts,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
		"low_stock_count": lowStockCount,
	})
}

func (h *PartHandler) GetPartDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var part models.Part
	if err := models.DB.First(&part, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "配件不存在")
		return
	}

	utils.Success(c, part)
}

type UpdatePartRequest struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	MinStock    int     `json:"min_stock"`
	Image       string  `json:"image"`
	Status      *bool   `json:"status"`
}

func (h *PartHandler) UpdatePart(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdatePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price != 0 {
		updates["price"] = req.Price
	}
	if req.Stock != 0 {
		updates["stock"] = req.Stock
	}
	if req.MinStock != 0 {
		updates["min_stock"] = req.MinStock
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := models.DB.Model(&models.Part{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		logger.Errorf("Update part error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

func (h *PartHandler) DeletePart(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := models.DB.Delete(&models.Part{}, id).Error; err != nil {
		logger.Errorf("Delete part error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

type CreatePartRequestRequest struct {
	Items  []PartRequestItemData `json:"items" binding:"required"`
	Remark string               `json:"remark"`
}

type PartRequestItemData struct {
	PartID   uint `json:"part_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

func (h *PartHandler) CreatePartRequest(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreatePartRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	tx := models.DB.Begin()

	partRequest := models.PartRequest{
		RequestNo:    utils.GenerateOrderNo("PR"),
		TechnicianID: userID,
		Status:       "pending",
		Remark:       req.Remark,
	}

	if err := tx.Create(&partRequest).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Create part request error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "创建申请失败")
		return
	}

	var totalAmount float64
	for _, item := range req.Items {
		var part models.Part
		if err := tx.First(&part, item.PartID).Error; err != nil {
			tx.Rollback()
			utils.Error(c, http.StatusNotFound, 404, "配件不存在")
			return
		}

		partRequestItem := models.PartRequestItem{
			PartRequestID: partRequest.ID,
			PartID:        item.PartID,
			Quantity:      item.Quantity,
			Price:         part.Price,
		}

		if err := tx.Create(&partRequestItem).Error; err != nil {
			tx.Rollback()
			logger.Errorf("Create part request item error: %v", err)
			utils.Error(c, http.StatusInternalServerError, 500, "创建申请失败")
			return
		}

		totalAmount += part.Price * float64(item.Quantity)
	}

	partRequest.TotalAmount = totalAmount
	if err := tx.Save(&partRequest).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "创建申请失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "创建申请失败")
		return
	}

	utils.Success(c, gin.H{
		"request_id": partRequest.ID,
		"request_no": partRequest.RequestNo,
		"message":    "配件申请已提交",
	})
}

func (h *PartHandler) GetPartRequestList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var requests []models.PartRequest
	query := models.DB.Preload("Technician")

	if role == "technician" {
		query = query.Where("technician_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.PartRequest{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests)

	utils.Success(c, gin.H{
		"list":     requests,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *PartHandler) GetPartRequestDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.PartRequest
	if err := models.DB.Preload("Technician").First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	var items []models.PartRequestItem
	models.DB.Preload("Part").Where("part_request_id = ?", id).Find(&items)

	utils.Success(c, gin.H{
		"request": request,
		"items":   items,
	})
}

func (h *PartHandler) ApprovePartRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.PartRequest
	if err := models.DB.First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	if request.Status != "pending" {
		utils.Error(c, http.StatusBadRequest, 400, "申请状态不允许审批")
		return
	}

	tx := models.DB.Begin()

	now := time.Now()
	if err := tx.Model(&request).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_at": now,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "审批失败")
		return
	}

	var items []models.PartRequestItem
	models.DB.Where("part_request_id = ?", id).Find(&items)
	for _, item := range items {
		var part models.Part
		tx.First(&part, item.PartID)
		if part.Stock >= item.Quantity {
			tx.Model(&part).Update("stock", gorm.Expr("stock - ?", item.Quantity))
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "审批失败")
		return
	}

	utils.Success(c, gin.H{"message": "审批通过"})
}

func (h *PartHandler) RejectPartRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	type RejectRequest struct {
		Remark string `json:"remark" binding:"required"`
	}

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var request models.PartRequest
	if err := models.DB.First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	if request.Status != "pending" {
		utils.Error(c, http.StatusBadRequest, 400, "申请状态不允许审批")
		return
	}

	models.DB.Model(&request).Updates(map[string]interface{}{
		"status": "rejected",
		"remark": req.Remark,
	})

	utils.Success(c, gin.H{"message": "已拒绝"})
}

func (h *PartHandler) ShipPartRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.PartRequest
	if err := models.DB.First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	if request.Status != "approved" {
		utils.Error(c, http.StatusBadRequest, 400, "申请状态不允许发货")
		return
	}

	now := time.Now()
	models.DB.Model(&request).Updates(map[string]interface{}{
		"status":     "shipped",
		"shipped_at": now,
	})

	utils.Success(c, gin.H{"message": "已发货"})
}

func (h *PartHandler) ReceivePartRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.PartRequest
	if err := models.DB.First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	if request.TechnicianID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权操作")
		return
	}

	if request.Status != "shipped" {
		utils.Error(c, http.StatusBadRequest, 400, "申请状态不允许收货")
		return
	}

	now := time.Now()
	models.DB.Model(&request).Updates(map[string]interface{}{
		"status":      "received",
		"received_at": now,
	})

	utils.Success(c, gin.H{"message": "已确认收货"})
}

type UsePartRequest struct {
	OrderID    uint `json:"order_id" binding:"required"`
	PartID     uint `json:"part_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required"`
}

func (h *PartHandler) UsePart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UsePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	if err := CheckOrderExists(req.OrderID); err != nil {
		utils.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	var part models.Part
	if err := models.DB.First(&part, req.PartID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "配件不存在")
		return
	}

	if part.Stock < req.Quantity {
		utils.Error(c, http.StatusBadRequest, 400, "库存不足")
		return
	}

	tx := models.DB.Begin()

	usage := models.PartUsage{
		OrderID:      req.OrderID,
		PartID:       req.PartID,
		TechnicianID: userID,
		Quantity:     req.Quantity,
		Price:        part.Price,
	}

	if err := tx.Create(&usage).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Create part usage error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "记录失败")
		return
	}

	tx.Model(&part).Update("stock", gorm.Expr("stock - ?", req.Quantity))

	if part.Stock-req.Quantity <= part.MinStock {
		logger.Warnf("Part %s stock low: %d", part.Name, part.Stock-req.Quantity)
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "记录失败")
		return
	}

	utils.Success(c, gin.H{"message": "使用记录已保存"})
}

func (h *PartHandler) CheckLowStock() {
	var lowStockParts []models.Part
	models.DB.Where("stock <= min_stock AND status = ?", true).Find(&lowStockParts)

	for _, part := range lowStockParts {
		logger.Warnf("Low stock alert: %s (Code: %s), Current: %d, Min: %d",
			part.Name, part.Code, part.Stock, part.MinStock)
	}
}
