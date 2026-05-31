package handlers

import (
	"io"
	"net/http"
	"strconv"

	"print3d-platform/internal/middleware"
	"print3d-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) InitiateUpload(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	var req service.InitiateUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.fileService.InitiateUpload(c.Request.Context(), authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *FileHandler) UploadChunk(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	uploadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload ID"})
		return
	}

	chunkNumber, err := strconv.Atoi(c.Query("chunk"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chunk number"})
		return
	}

	body := io.Reader(c.Request.Body)
	err = h.fileService.UploadChunk(c.Request.Context(), authUser.UserID, uploadID, chunkNumber, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chunk uploaded successfully"})
}

func (h *FileHandler) CompleteUpload(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	uploadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload ID"})
		return
	}

	upload, err := h.fileService.CompleteUpload(c.Request.Context(), authUser.UserID, uploadID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, upload)
}

func (h *FileHandler) GetUpload(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload ID"})
		return
	}

	upload, err := h.fileService.GetUpload(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	if upload.UserID != authUser.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, upload)
}

func (h *FileHandler) GetUserUploads(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	uploads, total, err := h.fileService.GetUserUploads(c.Request.Context(), authUser.UserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get uploads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  uploads,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *FileHandler) DeleteUpload(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload ID"})
		return
	}

	err = h.fileService.DeleteUpload(c.Request.Context(), id, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload deleted successfully"})
}

func (h *FileHandler) GetAccessLogs(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.fileService.GetAccessLogs(c.Request.Context(), fileID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
