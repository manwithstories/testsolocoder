package handlers

import (
	"net/http"
	"strconv"
	"time"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct{}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{}
}

func updateInventoryStatus(inventory *models.Inventory) {
	now := time.Now()
	daysUntilExpiry := inventory.ExpiryDate.Sub(now).Hours() / 24

	if inventory.Quantity <= 0 {
		inventory.Status = "sold_out"
	} else if inventory.Quantity < inventory.Threshold {
		inventory.Status = "low_stock"
	} else if daysUntilExpiry < 30 {
		inventory.Status = "expiring_soon"
	} else if daysUntilExpiry < 0 {
		inventory.Status = "expired"
	} else {
		inventory.Status = "in_stock"
	}
}

func (h *InventoryHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	query := database.DB.Model(&models.Inventory{}).Where("user_id = ?", userID)

	if honeyType := c.Query("honey_type"); honeyType != "" {
		query = query.Where("honey_type = ?", honeyType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if grade := c.Query("grade"); grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if batchCode := c.Query("batch_code"); batchCode != "" {
		query = query.Where("batch_code ILIKE ?", "%"+batchCode+"%")
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var inventories []models.Inventory
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("Harvest.Beehive").Find(&inventories)

	utils.SuccessWithTotal(c, inventories, total)
}

func (h *InventoryHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var inventory models.Inventory
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Harvest.Beehive").First(&inventory).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inventory not found")
		return
	}

	utils.Success(c, inventory)
}

func (h *InventoryHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var inventory models.Inventory
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&inventory).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inventory not found")
		return
	}

	if req.Quantity != nil {
		inventory.Quantity = *req.Quantity
	}
	if req.ExpiryDate != nil {
		expiryDate, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err != nil {
			utils.Fail(c, http.StatusBadRequest, "invalid date format")
			return
		}
		inventory.ExpiryDate = expiryDate
	}
	if req.InspectionReport != nil {
		inventory.InspectionReport = *req.InspectionReport
	}
	if req.Grade != nil {
		inventory.Grade = *req.Grade
	}
	if req.Threshold != nil {
		inventory.Threshold = *req.Threshold
	}
	if req.Price != nil {
		inventory.Price = *req.Price
	}

	updateInventoryStatus(&inventory)

	database.DB.Save(&inventory)

	utils.Success(c, inventory)
}

func (h *InventoryHandler) GetAlerts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var lowStockItems []models.Inventory
	database.DB.Where("user_id = ? AND status = ?", userID, "low_stock").
		Preload("Harvest.Beehive").Find(&lowStockItems)

	var expiringSoonItems []models.Inventory
	database.DB.Where("user_id = ? AND status = ?", userID, "expiring_soon").
		Preload("Harvest.Beehive").Find(&expiringSoonItems)

	utils.Success(c, gin.H{
		"low_stock":      lowStockItems,
		"expiring_soon":  expiringSoonItems,
	})
}

func (h *InventoryHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Inventory{})
	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to delete inventory")
		return
	}
	if result.RowsAffected == 0 {
		utils.Fail(c, http.StatusNotFound, "inventory not found")
		return
	}

	utils.Success(c, nil)
}
