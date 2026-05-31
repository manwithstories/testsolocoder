package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"tea-platform/config"
)

func UploadFile(c *gin.Context, file *multipart.FileHeader, uploadType string) (string, error) {
	cfg := config.Get()

	var uploadPath string
	switch uploadType {
	case "tea":
		uploadPath = cfg.Upload.TeaImagesPath
	case "packaging":
		uploadPath = cfg.Upload.PackagingImagesPath
	case "evaluation":
		uploadPath = cfg.Upload.EvaluationImagesPath
	default:
		uploadPath = "./uploads"
	}

	if err := validateFile(file, cfg); err != nil {
		return "", err
	}

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	fileName := generateFileName(ext)
	fullPath := filepath.Join(uploadPath, fileName)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	return fullPath, nil
}

func validateFile(file *multipart.FileHeader, cfg *config.Config) error {
	maxSize := int64(cfg.Upload.MaxSize) * 1024 * 1024
	if file.Size > maxSize {
		return fmt.Errorf("文件大小超过限制 (最大 %dMB)", cfg.Upload.MaxSize)
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
	allowed := false
	for _, t := range cfg.Upload.AllowedTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

func generateFileName(ext string) string {
	now := time.Now()
	return fmt.Sprintf("%d_%d%s", now.UnixNano(), now.Nanosecond(), ext)
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
