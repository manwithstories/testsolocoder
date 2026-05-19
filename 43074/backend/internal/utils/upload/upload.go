package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	apperrors "booklibrary/internal/errors"
	"booklibrary/internal/config"
)

func ValidateImage(file *multipart.FileHeader) error {
	if file.Size > config.AppConfig.Upload.MaxSize {
		return apperrors.ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedExt := range config.AppConfig.Upload.AllowedExt {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return apperrors.ErrFileTypeNotAllowed
	}

	return nil
}

func SaveImage(file *multipart.FileHeader) (string, error) {
	if err := ValidateImage(file); err != nil {
		return "", err
	}

	if err := os.MkdirAll(config.AppConfig.Upload.SavePath, 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := generateFilename() + ext
	filepath := filepath.Join(config.AppConfig.Upload.SavePath, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", config.AppConfig.Upload.AccessURL, filename), nil
}

func DeleteImage(accessURL string) error {
	if accessURL == "" {
		return nil
	}

	filename := filepath.Base(accessURL)
	filepath := filepath.Join(config.AppConfig.Upload.SavePath, filename)

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return os.Remove(filepath)
}

func generateFilename() string {
	now := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%d", now)))
	return hex.EncodeToString(hash[:])
}
