package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *services.UploadService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		uploadService: services.NewUploadService(),
	}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件: "+err.Error())
		return
	}

	result, err := h.uploadService.UploadFile(file)
	if err != nil {
		if err.Error() == "不支持的文件类型，仅支持图片和PDF文件" ||
			err.Error() == "不支持的文件扩展名" ||
			strings.Contains(err.Error(), "文件大小超过限制") {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, result)
}

func (h *UploadHandler) GetFile(c *gin.Context) {
	dateDir := c.Param("date")
	filename := c.Param("filename")

	relativePath := filepath.Join(dateDir, filename)

	filePath, contentType, err := h.uploadService.GetFilePath(relativePath)
	if err != nil {
		if err.Error() == "文件不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		if err.Error() == "非法的文件路径" {
			utils.Forbidden(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=31536000")
	c.File(filePath)
}

func (h *UploadHandler) DeleteFile(c *gin.Context) {
	relativePath := c.Query("path")
	if relativePath == "" {
		utils.BadRequest(c, "请提供文件路径参数 path")
		return
	}

	decodedPath, err := url.QueryUnescape(relativePath)
	if err != nil {
		utils.BadRequest(c, "路径解码失败")
		return
	}

	if err := h.uploadService.DeleteFile(decodedPath); err != nil {
		if err.Error() == "文件不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		if err.Error() == "非法的文件路径" {
			utils.Forbidden(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, nil)
}

func RegisterUploadRoutes(api *gin.RouterGroup) {
	handler := NewUploadHandler()

	uploads := api.Group("/uploads")
	{
		uploads.GET("/:date/*filename", handler.GetFile)

		auth := uploads.Group("")
		auth.Use(middleware.Auth())
		{
			auth.POST("", handler.UploadFile)
			auth.DELETE("", handler.DeleteFile)
		}
	}
}
