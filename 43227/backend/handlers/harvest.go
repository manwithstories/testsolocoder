package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type HarvestHandler struct{}

func NewHarvestHandler() *HarvestHandler {
	return &HarvestHandler{}
}

func generateBatchCode() string {
	now := time.Now()
	return fmt.Sprintf("BATCH-%s-%d", now.Format("20060102150405"), now.UnixNano()%1000)
}

func (h *HarvestHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateHarvestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var beehive models.Beehive
	if err := database.DB.Where("id = ? AND user_id = ?", req.BeehiveID, userID).First(&beehive).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "beehive not found")
		return
	}

	harvestDate, err := time.Parse("2006-01-02", req.HarvestDate)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid date format")
		return
	}

	batchCode := req.BatchCode
	if batchCode == "" {
		batchCode = generateBatchCode()
	}

	var existingHarvest models.Harvest
	if err := database.DB.Where("batch_code = ?", batchCode).First(&existingHarvest).Error; err == nil {
		utils.Fail(c, http.StatusConflict, "batch code already exists")
		return
	}

	quality := req.Quality
	if quality == "" {
		quality = "normal"
	}

	harvest := models.Harvest{
		UserID:      userID.(uint),
		BeehiveID:   req.BeehiveID,
		HarvestDate: harvestDate,
		HoneyType:   req.HoneyType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Quality:     quality,
		BatchCode:   batchCode,
		Notes:       req.Notes,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to begin transaction")
		return
	}

	if err := tx.Create(&harvest).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create harvest", err)
		return
	}

	expiryDate := harvestDate.AddDate(2, 0, 0)
	inventory := models.Inventory{
		UserID:     userID.(uint),
		HarvestID:  harvest.ID,
		HoneyType:  req.HoneyType,
		BatchCode:  batchCode,
		Quantity:   req.Quantity,
		Unit:       req.Unit,
		ExpiryDate: expiryDate,
		Grade:      "ungraded",
		Status:     "in_stock",
		Threshold:  10,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create inventory", err)
		return
	}

	now := time.Now()
	if err := tx.Model(&beehive).Updates(map[string]interface{}{
		"status":          "harvesting",
		"last_inspection": &now,
	}).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to update beehive", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	utils.Success(c, gin.H{
		"harvest":   harvest,
		"inventory": inventory,
	})
}

func (h *HarvestHandler) List(c *gin.Context) {
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

	query := database.DB.Model(&models.Harvest{}).Where("user_id = ?", userID)

	if beehiveID := c.Query("beehive_id"); beehiveID != "" {
		query = query.Where("beehive_id = ?", beehiveID)
	}
	if honeyType := c.Query("honey_type"); honeyType != "" {
		query = query.Where("honey_type = ?", honeyType)
	}
	if quality := c.Query("quality"); quality != "" {
		query = query.Where("quality = ?", quality)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("harvest_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("harvest_date <= ?", endDate)
	}

	sortBy := c.DefaultQuery("sort_by", "harvest_date")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var harvests []models.Harvest
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("Beehive").Find(&harvests)

	utils.SuccessWithTotal(c, harvests, total)
}

func (h *HarvestHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var harvest models.Harvest
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Beehive").First(&harvest).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "harvest not found")
		return
	}

	utils.Success(c, harvest)
}

func (h *HarvestHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var harvest models.Harvest
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&harvest).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "harvest not found")
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to begin transaction")
		return
	}

	if err := tx.Where("harvest_id = ?", id).Delete(&models.Inventory{}).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to delete inventory", err)
		return
	}

	if err := tx.Delete(&harvest).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to delete harvest", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	utils.Success(c, nil)
}
