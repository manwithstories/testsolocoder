package services

import (
	"errors"
	"fmt"
	"medical-platform/internal/config"
	"medical-platform/pkg/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

var allowedContentTypes = map[string]string{
	"image/jpeg":      "jpg",
	"image/png":       "png",
	"image/gif":       "gif",
	"image/bmp":       "bmp",
	"image/webp":      "webp",
	"application/pdf": "pdf",
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
	".pdf":  true,
}

type UploadResult struct {
	Filename    string `json:"filename"`
	OriginalName string `json:"original_name"`
	FilePath    string `json:"file_path"`
	FileURL     string `json:"file_url"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
}

func (s *UploadService) UploadFile(file *multipart.FileHeader) (*UploadResult, error) {
	if file.Size > config.AppConfig.MaxUploadSize {
		return nil, fmt.Errorf("文件大小超过限制，最大允许 %d MB", config.AppConfig.MaxUploadSize/1024/1024)
	}

	contentType := file.Header.Get("Content-Type")
	ext, allowed := allowedContentTypes[contentType]
	if !allowed {
		return nil, errors.New("不支持的文件类型，仅支持图片和PDF文件")
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[fileExt] {
		return nil, errors.New("不支持的文件扩展名")
	}

	if err := os.MkdirAll(config.AppConfig.UploadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	randomStr := utils.GenerateRandomString(8)
	newFilename := fmt.Sprintf("%s_%s.%s", timestamp, randomStr, ext)

	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(config.AppConfig.UploadDir, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日期目录失败: %w", err)
	}

	filePath := filepath.Join(fullDir, newFilename)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	fileURL := fmt.Sprintf("/api/v1/uploads/%s/%s", dateDir, newFilename)

	return &UploadResult{
		Filename:     newFilename,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     file.Size,
		ContentType:  contentType,
	}, nil
}

func (s *UploadService) GetFilePath(relativePath string) (string, string, error) {
	filePath := filepath.Join(config.AppConfig.UploadDir, relativePath)
	cleanPath := filepath.Clean(filePath)

	uploadDir, err := filepath.Abs(config.AppConfig.UploadDir)
	if err != nil {
		return "", "", errors.New("获取上传目录绝对路径失败")
	}

	absFilePath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", "", errors.New("获取文件绝对路径失败")
	}

	if !strings.HasPrefix(absFilePath, uploadDir) {
		return "", "", errors.New("非法的文件路径")
	}

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return "", "", errors.New("文件不存在")
	}

	contentType := getContentType(absFilePath)

	return absFilePath, contentType, nil
}

func (s *UploadService) DeleteFile(relativePath string) error {
	filePath := filepath.Join(config.AppConfig.UploadDir, relativePath)
	cleanPath := filepath.Clean(filePath)

	uploadDir, err := filepath.Abs(config.AppConfig.UploadDir)
	if err != nil {
		return errors.New("获取上传目录绝对路径失败")
	}

	absFilePath, err := filepath.Abs(cleanPath)
	if err != nil {
		return errors.New("获取文件绝对路径失败")
	}

	if !strings.HasPrefix(absFilePath, uploadDir) {
		return errors.New("非法的文件路径")
	}

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}

	if err := os.Remove(absFilePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func init() {
	_ = http.DetectContentType
}
