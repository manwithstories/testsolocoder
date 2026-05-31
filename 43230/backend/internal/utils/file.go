package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("ORD%s%s",
		now.Format("20060102150405"),
		fmt.Sprintf("%04d", now.Nanosecond()/1000000))
}

func GenerateTransactionNo() string {
	now := time.Now()
	return fmt.Sprintf("TXN%s%s",
		now.Format("20060102150405"),
		fmt.Sprintf("%06d", now.Nanosecond()/1000))
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func CheckPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

func HashFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	h := sha256.New()
	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func HashFilePath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func MD5Hash(data string) string {
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}

func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func IsAllowedExtension(ext string, allowedExtensions []string) bool {
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
