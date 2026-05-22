package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type FileHandler struct {
	cfg *config.Config
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{cfg: cfg}
}

func (h *FileHandler) Upload(c *gin.Context) {
	userID, _ := c.Get("user_id")

	fileType := c.DefaultPostForm("file_type", "image")
	if fileType != "image" && fileType != "video" && fileType != "document" {
		utils.BadRequest(c, "Invalid file type")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "No file provided")
		return
	}
	defer file.Close()

	if err := utils.ValidateUploadFile(header, fileType, &h.cfg.Upload); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)

	dir := fmt.Sprintf("%s/%s", h.cfg.Upload.UploadDir, fileType)
	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.InternalError(c, "Failed to create upload directory")
		return
	}

	fullPath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(header, fullPath); err != nil {
		utils.InternalError(c, "Failed to save file")
		return
	}

	category := models.FileType(fileType)
	fileRecord := models.File{
		UserID:   userID.(uuid.UUID),
		FileType: category,
		FileName: header.Filename,
		FilePath: fullPath,
		FileURL:  fmt.Sprintf("/uploads/%s/%s", fileType, filename),
		FileSize: header.Size,
		MimeType: header.Header.Get("Content-Type"),
	}

	if err := database.DB.Create(&fileRecord).Error; err != nil {
		utils.InternalError(c, "Failed to save file record")
		return
	}

	utils.Created(c, gin.H{
		"id":        fileRecord.ID,
		"url":       fileRecord.FileURL,
		"file_url":  fileRecord.FileURL,
		"name":      fileRecord.FileName,
		"file_name": fileRecord.FileName,
		"size":      fileRecord.FileSize,
		"file_size": fileRecord.FileSize,
		"type":      fileRecord.FileType,
		"file_type": fileRecord.FileType,
	})
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	fileType := c.Query("file_type")

	query := database.DB.Where("user_id = ?", userID)
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	var files []models.File
	var total int64

	query.Model(&models.File{}).Count(&total)

	page, pageSize := getPagination(c)
	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&files)

	utils.Paginated(c, files, total, page, pageSize)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	id := c.Param("id")
	fileID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid file ID")
		return
	}

	var file models.File
	if err := database.DB.First(&file, fileID).Error; err != nil {
		utils.NotFound(c, "File not found")
		return
	}

	if role != "admin" && file.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	if err := os.Remove(file.FilePath); err != nil {
		utils.Logger.Warnf("Failed to remove physical file: %v", err)
	}

	database.DB.Delete(&file)
	utils.Success(c, gin.H{"message": "File deleted"})
}
