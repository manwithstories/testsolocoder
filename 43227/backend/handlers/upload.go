package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"beehive-platform/config"
	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	config *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{config: cfg}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	userID, _ := c.Get("user_id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "no file provided")
		return
	}
	defer file.Close()

	if header.Size > h.config.Upload.MaxSize {
		utils.Fail(c, http.StatusBadRequest, "file size exceeds limit")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := false
	for _, allowedExt := range h.config.Upload.AllowExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		utils.Fail(c, http.StatusBadRequest, "file type not allowed")
		return
	}

	category := c.DefaultPostForm("category", "other")

	dateDir := time.Now().Format("2006/01/02")
	uploadDir := filepath.Join(h.config.Upload.Path, category, dateDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	newFilename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	filePath := filepath.Join(uploadDir, newFilename)

	if err := c.SaveUploadedFile(header, filePath); err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	relativePath := filepath.Join("uploads", category, dateDir, newFilename)

	upload := models.Upload{
		UserID:   userID.(uint),
		FileName: header.Filename,
		FilePath: relativePath,
		FileSize: header.Size,
		FileType: ext,
		Category: category,
	}

	database.DB.Create(&upload)

	utils.Success(c, gin.H{
		"id":        upload.ID,
		"file_name": upload.FileName,
		"file_path": upload.FilePath,
		"url":       "/static/" + relativePath,
		"file_size": upload.FileSize,
		"category":  upload.Category,
	})
}

func (h *UploadHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 20}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 20
	}

	query := database.DB.Model(&models.Upload{}).Where("user_id = ?", userID)

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	query.Count(&total)

	var uploads []models.Upload
	query.Order("created_at DESC").
		Offset(pageParams.GetOffset()).
		Limit(pageParams.PageSize).
		Find(&uploads)

	utils.SuccessWithTotal(c, uploads, total)
}

func (h *UploadHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var upload models.Upload
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&upload).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "file not found")
		return
	}

	fullPath := filepath.Join(h.config.Upload.Path, "..", upload.FilePath)
	os.Remove(fullPath)

	database.DB.Delete(&upload)

	utils.Success(c, nil)
}
