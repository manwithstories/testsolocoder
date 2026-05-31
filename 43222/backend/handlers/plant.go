package handlers

import (
	"garden-planner/database"
	"garden-planner/middleware"
	"garden-planner/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePlantRequest struct {
	Name            string  `json:"name" binding:"required,max=100"`
	LatinName       string  `json:"latin_name"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	GrowthCycle     int     `json:"growth_cycle"`
	WaterFrequency  string  `json:"water_frequency"`
	FertilizerNeed  string  `json:"fertilizer_need"`
	SunlightNeed    string  `json:"sunlight_need"`
	SoilPH          string  `json:"soil_ph"`
	PlantingSeason  string  `json:"planting_season"`
	HarvestSeason   string  `json:"harvest_season"`
	DiseaseInfo     string  `json:"disease_info"`
	PestInfo        string  `json:"pest_info"`
	ImageURL        string  `json:"image_url"`
	SowingDepth     string  `json:"sowing_depth"`
	Spacing         string  `json:"spacing"`
	Difficulty      string  `json:"difficulty"`
	ClimateZone     string  `json:"climate_zone"`
}

type CreatePlantingRecordRequest struct {
	PlotID        uuid.UUID `json:"plot_id" binding:"required"`
	PlantID       uuid.UUID `json:"plant_id" binding:"required"`
	Quantity      int       `json:"quantity"`
	PlantingDate  string    `json:"planting_date"`
	Notes         string    `json:"notes"`
}

type UpdatePlantingRecordRequest struct {
	Quantity           int    `json:"quantity"`
	Status             string `json:"status"`
	Notes              string `json:"notes"`
	ExpectedHarvestDate string `json:"expected_harvest_date"`
	ActualHarvestDate   string `json:"actual_harvest_date"`
}

func GetPlants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	category := c.Query("category")
	difficulty := c.Query("difficulty")
	climateZone := c.Query("climate_zone")
	search := c.Query("search")

	var plants []models.Plant
	var total int64

	query := database.DB.Model(&models.Plant{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if climateZone != "" {
		query = query.Where("climate_zone ILIKE ?", "%"+climateZone+"%")
	}
	if search != "" {
		query = query.Where("name ILIKE ? OR latin_name ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("name ASC").Find(&plants)

	c.JSON(http.StatusOK, gin.H{
		"plants":    plants,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetPlant(c *gin.Context) {
	id := c.Param("id")

	var plant models.Plant
	if err := database.DB.First(&plant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plant": plant})
}

func CreatePlant(c *gin.Context) {
	userType := middleware.GetUserType(c)
	if userType != "expert" && userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only experts and admins can create plants"})
		return
	}

	var req CreatePlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plant := models.Plant{
		ID:             uuid.New(),
		Name:           req.Name,
		LatinName:      req.LatinName,
		Category:       req.Category,
		Description:    req.Description,
		GrowthCycle:    req.GrowthCycle,
		WaterFrequency: req.WaterFrequency,
		FertilizerNeed: req.FertilizerNeed,
		SunlightNeed:   req.SunlightNeed,
		SoilPH:         req.SoilPH,
		PlantingSeason: req.PlantingSeason,
		HarvestSeason:  req.HarvestSeason,
		DiseaseInfo:    req.DiseaseInfo,
		PestInfo:       req.PestInfo,
		ImageURL:       req.ImageURL,
		SowingDepth:    req.SowingDepth,
		Spacing:        req.Spacing,
		Difficulty:     req.Difficulty,
		ClimateZone:    req.ClimateZone,
	}

	if err := database.DB.Create(&plant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plant created successfully",
		"plant":   plant,
	})
}

func CreatePlantingRecord(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreatePlantingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plot models.Plot
	if err := database.DB.Where("id = ? AND user_id = ?", req.PlotID, userID).First(&plot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plot not found"})
		return
	}

	var plant models.Plant
	if err := database.DB.First(&plant, req.PlantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
		return
	}

	record := models.PlantingRecord{
		ID:       uuid.New(),
		PlotID:   req.PlotID,
		PlantID:  req.PlantID,
		UserID:   userID,
		Quantity: req.Quantity,
		Notes:    req.Notes,
		Status:   "planted",
	}

	if req.PlantingDate != "" {
		if t, err := parseDate(req.PlantingDate); err == nil {
			record.PlantingDate = t
		}
	}

	if plant.GrowthCycle > 0 {
		expectedHarvest := record.PlantingDate.AddDate(0, 0, plant.GrowthCycle)
		record.ExpectedHarvestDate = &expectedHarvest
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create planting record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Planting record created successfully",
		"record":  record,
	})
}

func GetPlantingRecords(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	plotID := c.Query("plot_id")
	status := c.Query("status")

	var records []models.PlantingRecord
	var total int64

	query := database.DB.Where("user_id = ?", userID)

	if plotID != "" {
		query = query.Where("plot_id = ?", plotID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.PlantingRecord{}).Count(&total)
	query.Preload("Plant").Preload("Plot").
		Offset(offset).Limit(pageSize).
		Order("planting_date DESC").
		Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"records":   records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetPlantingRecord(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var record models.PlantingRecord
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Plant").Preload("Plot").Preload("GrowthLogs").
		First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"record": record})
}

func UpdatePlantingRecord(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var record models.PlantingRecord
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	var req UpdatePlantingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity > 0 {
		record.Quantity = req.Quantity
	}
	if req.Status != "" {
		record.Status = req.Status
	}
	if req.Notes != "" {
		record.Notes = req.Notes
	}
	if req.ExpectedHarvestDate != "" {
		if t, err := parseDate(req.ExpectedHarvestDate); err == nil {
			record.ExpectedHarvestDate = &t
		}
	}
	if req.ActualHarvestDate != "" {
		if t, err := parseDate(req.ActualHarvestDate); err == nil {
			record.ActualHarvestDate = &t
		}
	}

	if err := database.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update planting record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Planting record updated successfully",
		"record":  record,
	})
}

func DeletePlantingRecord(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PlantingRecord{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete planting record"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Planting record deleted successfully"})
}
