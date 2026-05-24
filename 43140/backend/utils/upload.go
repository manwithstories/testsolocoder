package utils

import (
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveUploadedFile(c *gin.Context, file *multipart.FileHeader, uploadDir string) (string, string, error) {
	ext := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFilename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", "", err
	}

	return filePath, newFilename, nil
}
