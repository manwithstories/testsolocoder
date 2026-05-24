package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"business-registration-platform/config"
	"github.com/google/uuid"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
}

func UploadFile(file io.Reader, filename string, folder string) (string, error) {
	ext := filepath.Ext(filename)
	ext = strings.ToLower(ext)

	if !allowedExtensions[ext] {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	uploadPath := config.AppConfig.File.UploadPath
	fullPath := filepath.Join(uploadPath, folder)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", err
	}

	newFilename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)
	filePath := filepath.Join(fullPath, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	relativePath := filepath.Join(folder, newFilename)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	return relativePath, nil
}

func GetFileURL(relativePath string) string {
	if relativePath == "" {
		return ""
	}
	return fmt.Sprintf("/uploads/%s", relativePath)
}

func DeleteFile(relativePath string) error {
	if relativePath == "" {
		return nil
	}

	fullPath := filepath.Join(config.AppConfig.File.UploadPath, relativePath)
	return os.Remove(fullPath)
}

func ValidateFileSize(size int64) error {
	if size > config.AppConfig.File.MaxSize {
		return fmt.Errorf("file size exceeds maximum limit of %d bytes", config.AppConfig.File.MaxSize)
	}
	return nil
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
