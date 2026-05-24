package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileValidator struct {
	MaxFileSize  int64
	AllowedTypes []string
	UploadPath   string
}

func NewFileValidator(maxSizeMB int, allowedTypes []string, uploadPath string) *FileValidator {
	return &FileValidator{
		MaxFileSize:  int64(maxSizeMB) * 1024 * 1024,
		AllowedTypes: allowedTypes,
		UploadPath:   uploadPath,
	}
}

func (fv *FileValidator) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > fv.MaxFileSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %dMB", fv.MaxFileSize/(1024*1024))
	}

	fileType := file.Header.Get("Content-Type")
	if !fv.isAllowedType(fileType) {
		return errors.New("不支持的文件类型")
	}

	return nil
}

func (fv *FileValidator) isAllowedType(fileType string) bool {
	for _, t := range fv.AllowedTypes {
		if t == fileType {
			return true
		}
	}
	return false
}

func (fv *FileValidator) SaveFile(file *multipart.FileHeader, subDir string) (string, string, error) {
	if err := fv.ValidateFile(file); err != nil {
		return "", "", err
	}

	dir := filepath.Join(fv.UploadPath, subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", fmt.Errorf("创建上传目录失败: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	ext = strings.ToLower(ext)

	uniqueName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)
	fullPath := filepath.Join(dir, uniqueName)

	if err := saveUploadedFile(file, fullPath); err != nil {
		return "", "", fmt.Errorf("保存文件失败: %w", err)
	}

	relativePath := filepath.Join(subDir, uniqueName)
	relativePath = filepath.ToSlash(relativePath)

	return fullPath, relativePath, nil
}

func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	return err
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}
