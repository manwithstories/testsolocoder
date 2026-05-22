package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"travel-planner/internal/config"
	"travel-planner/internal/database"
	"travel-planner/internal/logger"
	"travel-planner/internal/middleware"
	"travel-planner/internal/models"
	"travel-planner/internal/utils"
)

func UploadFile(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	userID := middleware.GetCurrentUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "Failed to get file from request")
		return
	}
	defer file.Close()

	if header.Size > config.AppConfig.Upload.MaxSize {
		utils.BadRequest(c, fmt.Sprintf("File size exceeds maximum limit of %d bytes", config.AppConfig.Upload.MaxSize))
		return
	}

	contentType := header.Header.Get("Content-Type")
	allowed := false
	for _, t := range config.AppConfig.Upload.AllowedTypes {
		if strings.HasPrefix(contentType, t) || contentType == t {
			allowed = true
			break
		}
	}
	if !allowed {
		utils.BadRequest(c, "File type not allowed")
		return
	}

	uploadPath := config.AppConfig.Upload.Path
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Errorf("Failed to create upload directory: %v", err)
		utils.InternalServerError(c, "Failed to upload file")
		return
	}

	fileExt := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), fileExt)
	filePath := filepath.Join(uploadPath, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("Failed to create file: %v", err)
		utils.InternalServerError(c, "Failed to upload file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		logger.Errorf("Failed to save file: %v", err)
		utils.InternalServerError(c, "Failed to upload file")
		return
	}

	category := c.PostForm("category")
	description := c.PostForm("description")

	fileRecord := models.File{
		PlanID:       planUUID,
		UploadedBy:   userID,
		FileName:     newFileName,
		OriginalName: header.Filename,
		FileType:     contentType,
		FileSize:     header.Size,
		FileURL:      fmt.Sprintf("/api/files/%s", newFileName),
		Category:     category,
		Description:  description,
	}

	if err := database.DB.Create(&fileRecord).Error; err != nil {
		os.Remove(filePath)
		logger.Errorf("Failed to save file record: %v", err)
		utils.InternalServerError(c, "Failed to upload file")
		return
	}

	logger.Infof("File uploaded: %s by user %s", fileRecord.ID, userID)
	utils.Created(c, fileRecord)
}

func GetFiles(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	category := c.DefaultQuery("category", "")

	query := database.DB.Where("plan_id = ?", planUUID)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var files []models.File
	if err := query.Preload("Uploader").Order("created_at DESC").Find(&files).Error; err != nil {
		logger.Errorf("Failed to get files: %v", err)
		utils.InternalServerError(c, "Failed to get files")
		return
	}

	utils.Success(c, files)
}

func GetFile(c *gin.Context) {
	filename := c.Param("filename")

	filePath := filepath.Join(config.AppConfig.Upload.Path, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.NotFound(c, "File not found")
		return
	}

	var fileRecord models.File
	if err := database.DB.Where("file_name = ?", filename).First(&fileRecord).Error; err != nil {
		utils.NotFound(c, "File not found")
		return
	}

	c.FileAttachment(filePath, fileRecord.OriginalName)
}

func DownloadFile(c *gin.Context) {
	fileID := c.Param("id")
	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		utils.BadRequest(c, "Invalid file ID")
		return
	}

	var fileRecord models.File
	if err := database.DB.First(&fileRecord, "id = ?", fileUUID).Error; err != nil {
		utils.NotFound(c, "File not found")
		return
	}

	filePath := filepath.Join(config.AppConfig.Upload.Path, fileRecord.FileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.NotFound(c, "File not found on disk")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileRecord.OriginalName))
	c.Header("Content-Length", strconv.FormatInt(fileRecord.FileSize, 10))
	c.Header("Content-Type", fileRecord.FileType)
	c.File(filePath)
}

func DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		utils.BadRequest(c, "Invalid file ID")
		return
	}

	var fileRecord models.File
	if err := database.DB.First(&fileRecord, "id = ?", fileUUID).Error; err != nil {
		utils.NotFound(c, "File not found")
		return
	}

	filePath := filepath.Join(config.AppConfig.Upload.Path, fileRecord.FileName)
	if err := os.Remove(filePath); err != nil {
		logger.Warnf("Failed to delete file from disk: %v", err)
	}

	if err := database.DB.Delete(&fileRecord).Error; err != nil {
		logger.Errorf("Failed to delete file record: %v", err)
		utils.InternalServerError(c, "Failed to delete file")
		return
	}

	logger.Infof("File deleted: %s", fileID)
	utils.Success(c, gin.H{"message": "File deleted successfully"})
}
