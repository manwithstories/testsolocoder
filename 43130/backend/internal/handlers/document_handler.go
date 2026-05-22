package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"wedding-planner/config"
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"
	"wedding-planner/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct{}

func NewDocumentHandler() *DocumentHandler {
	return &DocumentHandler{}
}

type DocumentRequest struct {
	Category     string `json:"category"`
	Notes        string `json:"notes"`
	VendorID     *uint  `json:"vendor_id"`
	BudgetItemID *uint  `json:"budget_item_id"`
}

func (h *DocumentHandler) Upload(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	userID := c.GetUint("user_id")

	category := c.PostForm("category")
	notes := c.PostForm("notes")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Please upload a file")
		return
	}

	maxSize := int64(config.AppConfig.Server.MaxUploadMB) * 1024 * 1024
	if file.Size > maxSize {
		response.BadRequest(c, fmt.Sprintf("File size exceeds maximum limit of %dMB", config.AppConfig.Server.MaxUploadMB))
		return
	}

	uploadDir := filepath.Join(config.AppConfig.Server.UploadPath, fmt.Sprintf("wedding_%d", weddingID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "Failed to create upload directory")
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), utils.RandomString(8), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.InternalError(c, "Failed to save file")
		return
	}

	db := database.GetDB()

	document := models.Document{
		WeddingID:  weddingID,
		FileName:   file.Filename,
		FilePath:   filePath,
		FileSize:   file.Size,
		FileType:   ext,
		Category:   category,
		Notes:      notes,
		UploadedBy: userID,
		Version:    1,
	}

	if err := db.Create(&document).Error; err != nil {
		os.Remove(filePath)
		response.InternalError(c, "Failed to save document info")
		return
	}

	response.Created(c, document)
}

func (h *DocumentHandler) GetList(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var documents []models.Document

	category := c.Query("category")
	query := db.Where("wedding_id = ?", weddingID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Order("created_at DESC").Find(&documents)

	response.Success(c, documents)
}

func (h *DocumentHandler) GetByID(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var document models.Document
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&document).Error; err != nil {
		response.NotFound(c, "Document not found")
		return
	}

	var versions []models.Document
	db.Where("parent_id = ? OR id = ?", id, id).Order("version DESC").Find(&versions)

	response.Success(c, gin.H{
		"document": document,
		"versions": versions,
	})
}

func (h *DocumentHandler) Download(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var document models.Document
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&document).Error; err != nil {
		response.NotFound(c, "Document not found")
		return
	}

	if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
		response.NotFound(c, "File not found")
		return
	}

	c.FileAttachment(document.FilePath, document.FileName)
}

func (h *DocumentHandler) Update(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var document models.Document
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&document).Error; err != nil {
		response.NotFound(c, "Document not found")
		return
	}

	db.Model(&document).Updates(map[string]interface{}{
		"category":      req.Category,
		"notes":         req.Notes,
		"vendor_id":     req.VendorID,
		"budget_item_id": req.BudgetItemID,
	})

	response.Success(c, document)
}

func (h *DocumentHandler) UploadNewVersion(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")
	userID := c.GetUint("user_id")

	db := database.GetDB()

	var originalDoc models.Document
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&originalDoc).Error; err != nil {
		response.NotFound(c, "Document not found")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Please upload a file")
		return
	}

	maxSize := int64(config.AppConfig.Server.MaxUploadMB) * 1024 * 1024
	if file.Size > maxSize {
		response.BadRequest(c, fmt.Sprintf("File size exceeds maximum limit of %dMB", config.AppConfig.Server.MaxUploadMB))
		return
	}

	uploadDir := filepath.Join(config.AppConfig.Server.UploadPath, fmt.Sprintf("wedding_%d", weddingID))
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s_v%d%s", time.Now().Unix(), utils.RandomString(8), originalDoc.Version+1, ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.InternalError(c, "Failed to save file")
		return
	}

	newVersion := originalDoc.Version + 1

	newDoc := models.Document{
		WeddingID:     weddingID,
		FileName:      file.Filename,
		FilePath:      filePath,
		FileSize:      file.Size,
		FileType:      ext,
		Category:      originalDoc.Category,
		Version:       newVersion,
		ParentID:      &id,
		UploadedBy:    userID,
	}

	if err := db.Create(&newDoc).Error; err != nil {
		os.Remove(filePath)
		response.InternalError(c, "Failed to save document info")
		return
	}

	response.Created(c, newDoc)
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var document models.Document
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&document).Error; err != nil {
		response.NotFound(c, "Document not found")
		return
	}

	db.Where("parent_id = ?", id).Delete(&models.Document{})

	if err := db.Delete(&document).Error; err != nil {
		response.InternalError(c, "Failed to delete document")
		return
	}

	if _, err := os.Stat(document.FilePath); err == nil {
		os.Remove(document.FilePath)
	}

	response.Success(c, gin.H{"message": "Document deleted successfully"})
}

func (h *DocumentHandler) GetCategories(c *gin.Context) {
	categories := []string{
		"合同", "发票", "收据", "报价单", "设计稿", "其他",
	}

	response.Success(c, categories)
}

func init() {
	uploadPath := config.AppConfig.Server.UploadPath
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, 0755)
	}
}

var _ = io.EOF
