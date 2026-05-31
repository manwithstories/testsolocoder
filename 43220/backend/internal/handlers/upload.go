package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pet-board/internal/config"
	"pet-board/internal/utils"
)

type UploadHandler struct {
	cfg config.UploadConfig
}

func NewUploadHandler(cfg config.UploadConfig) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "no file uploaded")
		return
	}

	if file.Size > h.cfg.MaxSize {
		utils.BadRequest(c, fmt.Sprintf("file size exceeds maximum allowed size of %d MB", h.cfg.MaxSize/(1024*1024)))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, a := range h.cfg.AllowedExt {
		if ext == a {
			allowed = true
			break
		}
	}
	if !allowed {
		utils.BadRequest(c, fmt.Sprintf("file type %s is not allowed", ext))
		return
	}

	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	filePath := fmt.Sprintf("%s/%s", h.cfg.UploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	url := fmt.Sprintf("/uploads/%s", fileName)
	utils.Success(c, gin.H{"url": url, "filename": fileName})
}

func (h *UploadHandler) UploadMultiple(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "no files uploaded")
		return
	}

	files := form.File["files"]
	var urls []string

	for _, file := range files {
		if file.Size > h.cfg.MaxSize {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowed := false
		for _, a := range h.cfg.AllowedExt {
			if ext == a {
				allowed = true
				break
			}
		}
		if !allowed {
			continue
		}

		fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
		filePath := fmt.Sprintf("%s/%s", h.cfg.UploadDir, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			continue
		}

		url := fmt.Sprintf("/uploads/%s", fileName)
		urls = append(urls, url)
	}

	utils.Success(c, gin.H{"urls": urls})
}
