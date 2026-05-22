package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"e-learning-platform/internal/config"
)

func ValidateUploadFile(file *multipart.FileHeader, fileType string, cfg *config.UploadConfig) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	maxSize := int64(cfg.MaxSizeMB) * 1024 * 1024

	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds limit: %dMB", cfg.MaxSizeMB)
	}

	switch fileType {
	case "image":
		if !contains(strings.Split(cfg.AllowedImage, ","), strings.TrimPrefix(ext, ".")) {
			return fmt.Errorf("unsupported image format: %s", ext)
		}
	case "video":
		if !contains(strings.Split(cfg.AllowedVideo, ","), strings.TrimPrefix(ext, ".")) {
			return fmt.Errorf("unsupported video format: %s", ext)
		}
	case "document":
		if !contains(strings.Split(cfg.AllowedDocument, ","), strings.TrimPrefix(ext, ".")) {
			return fmt.Errorf("unsupported document format: %s", ext)
		}
	default:
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	return nil
}

func GetFileCategory(ext string) string {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	switch ext {
	case "jpg", "jpeg", "png", "gif", "webp":
		return "image"
	case "mp4", "mov", "avi", "mkv":
		return "video"
	case "pdf", "doc", "docx", "xls", "xlsx":
		return "document"
	default:
		return "other"
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == strings.TrimSpace(item) {
			return true
		}
	}
	return false
}
