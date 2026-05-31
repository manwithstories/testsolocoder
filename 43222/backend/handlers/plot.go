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

type CreatePlotRequest struct {
	Name              string  `json:"name" binding:"required,max=100"`
	Description       string  `json:"description"`
	SoilType          string  `json:"soil_type"`
	Sunlight          string  `json:"sunlight"`
	Area              float64 `json:"area"`
	Location          string  `json:"location"`
	GridConfig        string  `json:"grid_config"`
	IrrigationDevice  string  `json:"irrigation_device"`
	SensorData        string  `json:"sensor_data"`
}

type UpdatePlotRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	SoilType          string  `json:"soil_type"`
	Sunlight          string  `json:"sunlight"`
	Area              float64 `json:"area"`
	Location          string  `json:"location"`
	GridConfig        string  `json:"grid_config"`
	IrrigationDevice  string  `json:"irrigation_device"`
	SensorData        string  `json:"sensor_data"`
}

func CreatePlot(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreatePlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plot := models.Plot{
		ID:               uuid.New(),
		UserID:           userID,
		Name:             req.Name,
		Description:      req.Description,
		SoilType:         req.SoilType,
		Sunlight:         req.Sunlight,
		Area:             req.Area,
		Location:         req.Location,
		GridConfig:       req.GridConfig,
		IrrigationDevice: req.IrrigationDevice,
		SensorData:       req.SensorData,
	}

	if err := database.DB.Create(&plot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plot"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plot created successfully",
		"plot":    plot,
	})
}

func GetPlots(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var plots []models.Plot
	var total int64

	query := database.DB.Where("user_id = ?", userID)

	query.Model(&models.Plot{}).Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&plots)

	c.JSON(http.StatusOK, gin.H{
		"plots":     plots,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetPlot(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var plot models.Plot
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("PlantingRecords.Plant").
		First(&plot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plot not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plot": plot})
}

func UpdatePlot(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var plot models.Plot
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&plot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plot not found"})
		return
	}

	var req UpdatePlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		plot.Name = req.Name
	}
	if req.Description != "" {
		plot.Description = req.Description
	}
	if req.SoilType != "" {
		plot.SoilType = req.SoilType
	}
	if req.Sunlight != "" {
		plot.Sunlight = req.Sunlight
	}
	if req.Area > 0 {
		plot.Area = req.Area
	}
	if req.Location != "" {
		plot.Location = req.Location
	}
	if req.GridConfig != "" {
		plot.GridConfig = req.GridConfig
	}
	if req.IrrigationDevice != "" {
		plot.IrrigationDevice = req.IrrigationDevice
	}
	if req.SensorData != "" {
		plot.SensorData = req.SensorData
	}

	if err := database.DB.Save(&plot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plot updated successfully",
		"plot":    plot,
	})
}

func DeletePlot(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Plot{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plot"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plot not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plot deleted successfully"})
}
