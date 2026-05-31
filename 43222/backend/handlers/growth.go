package handlers

import (
	"garden-planner/config"
	"garden-planner/database"
	"garden-planner/middleware"
	"garden-planner/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateGrowthLogRequest struct {
	PlantingRecordID uuid.UUID `json:"planting_record_id" binding:"required"`
	Title             string    `json:"title" binding:"required"`
	Description       string    `json:"description"`
	LogType           string    `json:"log_type"`
	LogDate           string    `json:"log_date"`
}

func UploadFile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	if header.Size > config.AppConfig.MaxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit"})
		return
	}

	ext := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".pdf": true, ".doc": true, ".docx": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	uploadDir := filepath.Join(config.AppConfig.UploadDir, userID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	filename := uuid.New().String() + ext
	filepath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileURL := "/uploads/" + userID.String() + "/" + filename

	c.JSON(http.StatusOK, gin.H{
		"url":      fileURL,
		"filename": filename,
		"size":     header.Size,
	})
}

func CreateGrowthLog(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateGrowthLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var record models.PlantingRecord
	if err := database.DB.Where("id = ? AND user_id = ?", req.PlantingRecordID, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	log := models.GrowthLog{
		ID:               uuid.New(),
		PlantingRecordID: req.PlantingRecordID,
		Title:            req.Title,
		Description:      req.Description,
		LogType:          req.LogType,
	}

	if req.LogDate != "" {
		if t, err := parseDate(req.LogDate); err == nil {
			log.LogDate = t
		}
	} else {
		log.LogDate = time.Now()
	}

	if err := database.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create growth log"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Growth log created successfully",
		"log":     log,
	})
}

func GetGrowthLogs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	recordID := c.Query("planting_record_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var logs []models.GrowthLog
	var total int64

	query := database.DB.Joins("JOIN planting_records ON planting_records.id = growth_logs.planting_record_id").
		Where("planting_records.user_id = ?", userID)

	if recordID != "" {
		query = query.Where("growth_logs.planting_record_id = ?", recordID)
	}

	query.Model(&models.GrowthLog{}).Count(&total)
	query.Offset(offset).Limit(pageSize).Order("log_date DESC").Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetGrowthLog(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var log models.GrowthLog
	if err := database.DB.Joins("JOIN planting_records ON planting_records.id = growth_logs.planting_record_id").
		Where("growth_logs.id = ? AND planting_records.user_id = ?", id, userID).
		First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Growth log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": log})
}

func UpdateGrowthLog(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var log models.GrowthLog
	if err := database.DB.Joins("JOIN planting_records ON planting_records.id = growth_logs.planting_record_id").
		Where("growth_logs.id = ? AND planting_records.user_id = ?", id, userID).
		First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Growth log not found"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		LogType     string `json:"log_type"`
		LogDate     string `json:"log_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		log.Title = req.Title
	}
	if req.Description != "" {
		log.Description = req.Description
	}
	if req.ImageURL != "" {
		log.ImageURL = req.ImageURL
	}
	if req.LogType != "" {
		log.LogType = req.LogType
	}
	if req.LogDate != "" {
		if t, err := parseDate(req.LogDate); err == nil {
			log.LogDate = t
		}
	}

	if err := database.DB.Save(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update growth log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Growth log updated successfully",
		"log":     log,
	})
}

func DeleteGrowthLog(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Joins("JOIN planting_records ON planting_records.id = growth_logs.planting_record_id").
		Where("growth_logs.id = ? AND planting_records.user_id = ?", id, userID).
		Delete(&models.GrowthLog{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete growth log"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Growth log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Growth log deleted successfully"})
}

func ExportPlantingReport(c *gin.Context) {
	userID := middleware.GetUserID(c)
	recordID := c.Param("id")

	var record models.PlantingRecord
	if err := database.DB.Where("id = ? AND user_id = ?", recordID, userID).
		Preload("Plant").Preload("Plot").Preload("GrowthLogs").
		First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	report := gin.H{
		"plant":       record.Plant,
		"plot":        record.Plot,
		"record":      record,
		"growth_logs":  record.GrowthLogs,
		"generated_at": time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report generated successfully",
		"report":  report,
	})
}
