package handler

import (
	"property-management/internal/database"
	"property-management/internal/model"
	"property-management/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeeHandler struct{}

func NewFeeHandler() *FeeHandler {
	return &FeeHandler{}
}

type FeeRequest struct {
	PropertyID  uint    `json:"propertyId" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Month       string  `json:"month" binding:"required"`
	TotalAmount float64 `json:"totalAmount" binding:"required"`
	Units       float64 `json:"units"`
	UnitPrice   float64 `json:"unitPrice"`
}

func (h *FeeHandler) Create(c *gin.Context) {
	var req FeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	dueDate, _ := time.Parse("2006-01-02", req.Month+"-05")

	fee := model.UtilityFee{
		PropertyID:  req.PropertyID,
		Type:        req.Type,
		Month:       req.Month,
		TotalAmount: req.TotalAmount,
		Units:       req.Units,
		UnitPrice:   req.UnitPrice,
		Status:      0,
		DueDate:     dueDate,
	}

	if err := database.DB.Create(&fee).Error; err != nil {
		utils.Error(c, 500, "Failed to create utility fee")
		return
	}

	utils.Success(c, fee)
}

func (h *FeeHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var fee model.UtilityFee
	if err := database.DB.First(&fee, id).Error; err != nil {
		utils.Error(c, 404, "Utility fee not found")
		return
	}

	var req FeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&fee).Updates(map[string]interface{}{
		"type":         req.Type,
		"month":        req.Month,
		"total_amount": req.TotalAmount,
		"units":        req.Units,
		"unit_price":   req.UnitPrice,
	})

	utils.Success(c, nil)
}

func (h *FeeHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&model.UtilityFee{}, id)
	utils.Success(c, nil)
}

func (h *FeeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.UtilityFee{})

	if month := c.Query("month"); month != "" {
		query = query.Where("month = ?", month)
	}
	if feeType := c.Query("type"); feeType != "" {
		query = query.Where("type = ?", feeType)
	}
	if status := c.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	var total int64
	query.Count(&total)

	var fees []model.UtilityFee
	offset := (page - 1) * pageSize
	query.Preload("Property").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&fees)

	var totalAmount float64
	var paidAmount float64
	for _, fee := range fees {
		totalAmount += fee.TotalAmount
		if fee.Status == 1 {
			paidAmount += fee.TotalAmount
		}
	}

	utils.Success(c, gin.H{
		"list":        fees,
		"total":       total,
		"page":        page,
		"pageSize":    pageSize,
		"totalAmount": totalAmount,
		"paidAmount":  paidAmount,
	})
}

func (h *FeeHandler) Pay(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var fee model.UtilityFee
	if err := database.DB.First(&fee, id).Error; err != nil {
		utils.Error(c, 404, "Utility fee not found")
		return
	}

	now := time.Now()
	database.DB.Model(&fee).Updates(map[string]interface{}{
		"status":  1,
		"paid_at": now,
	})

	utils.Success(c, nil)
}

func (h *FeeHandler) BatchGenerate(c *gin.Context) {
	var req struct {
		Month       string   `json:"month" binding:"required"`
		Type        string   `json:"type" binding:"required"`
		PropertyIDs []uint   `json:"propertyIds" binding:"required"`
		Amount      float64  `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	dueDate, _ := time.Parse("2006-01-02", req.Month+"-05")
	generatedCount := 0

	for _, pid := range req.PropertyIDs {
		var count int64
		database.DB.Model(&model.UtilityFee{}).
			Where("property_id = ? AND month = ? AND type = ?", pid, req.Month, req.Type).
			Count(&count)

		if count == 0 {
			fee := model.UtilityFee{
				PropertyID:  pid,
				Type:        req.Type,
				Month:       req.Month,
				TotalAmount: req.Amount,
				Status:      0,
				DueDate:     dueDate,
			}
			database.DB.Create(&fee)
			generatedCount++
		}
	}

	utils.Success(c, gin.H{
		"generated": generatedCount,
	})
}
