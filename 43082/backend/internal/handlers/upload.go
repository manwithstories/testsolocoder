package handlers

import (
	"fmt"
	"multishop/internal/config"
	"multishop/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	cfg *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}

	if file.Size > h.cfg.MaxUploadSize {
		utils.Error(c, http.StatusBadRequest, "文件大小超过限制")
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		utils.Error(c, http.StatusBadRequest, "不支持的文件类型")
		return
	}

	if err := os.MkdirAll(h.cfg.UploadPath, 0755); err != nil {
		utils.Error(c, http.StatusInternalServerError, "创建目录失败")
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(h.cfg.UploadPath, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		utils.Error(c, http.StatusInternalServerError, "保存文件失败")
		return
	}

	url := "/uploads/" + filename
	utils.Success(c, gin.H{"url": url})
}

func (h *UploadHandler) UploadMultiple(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}

	if err := os.MkdirAll(h.cfg.UploadPath, 0755); err != nil {
		utils.Error(c, http.StatusInternalServerError, "创建目录失败")
		return
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	urls := make([]string, 0, len(files))
	for _, file := range files {
		if file.Size > h.cfg.MaxUploadSize {
			continue
		}
		ext := filepath.Ext(file.Filename)
		if !allowedExts[ext] {
			continue
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filepath := filepath.Join(h.cfg.UploadPath, filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			continue
		}

		urls = append(urls, "/uploads/"+filename)
	}

	utils.Success(c, gin.H{"urls": urls})
}
