package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	AllowedImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	AllowedFileTypes = map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
	}
)

func UploadFile(file *multipart.FileHeader, uploadPath string, allowedTypes map[string]bool, maxSize int64) (string, error) {
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	fileType := http.DetectContentType(buffer)
	if !allowedTypes[fileType] {
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}

	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to reset file position: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	ext = strings.ToLower(ext)

	uniqueName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)

	subDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadPath, subDir)

	err = os.MkdirAll(fullDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	filePath := filepath.Join(fullDir, uniqueName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}

	relativePath := filepath.Join("/uploads", subDir, uniqueName)
	relativePath = filepath.ToSlash(relativePath)

	return relativePath, nil
}

func DeleteFile(filePath, uploadPath string) error {
	if filePath == "" {
		return nil
	}

	relativePath := strings.TrimPrefix(filePath, "/uploads/")
	fullPath := filepath.Join(uploadPath, relativePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func GetFullPath(relativePath, uploadPath string) string {
	relativePath = strings.TrimPrefix(relativePath, "/uploads/")
	return filepath.Join(uploadPath, relativePath)
}
