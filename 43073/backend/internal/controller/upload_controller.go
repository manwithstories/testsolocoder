package controller

import (
	"os"
	"path/filepath"
	"ticket-system/config"
	"ticket-system/internal/common/response"
	"ticket-system/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (c *UploadController) UploadImage(ctx *gin.Context) {
	if err := os.MkdirAll(config.App.UploadPath, 0755); err != nil {
		response.ServerError(ctx, "创建上传目录失败")
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.BadRequest(ctx, "获取上传文件失败")
		return
	}

	if file.Size > config.App.MaxUploadSize {
		response.BadRequest(ctx, "文件大小超过限制")
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		response.BadRequest(ctx, "不支持的文件类型")
		return
	}

	filename := uuid.New().String() + ext
	dateDir := time.Now().Format("20060102")
	saveDir := filepath.Join(config.App.UploadPath, dateDir)

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		response.ServerError(ctx, "创建上传目录失败")
		return
	}

	savePath := filepath.Join(saveDir, filename)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		logger.Log.Errorf("Save file failed: %v", err)
		response.ServerError(ctx, "保存文件失败")
		return
	}

	url := "/uploads/" + dateDir + "/" + filename
	response.Success(ctx, gin.H{
		"url":  url,
		"name": file.Filename,
		"size": file.Size,
	})
}
