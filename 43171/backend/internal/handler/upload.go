package handler

import (
	"drone-rental/internal/config"
	"drone-rental/internal/pkg/response"
	"drone-rental/internal/pkg/utils"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrParam(c, "请选择文件")
		return
	}
	if err := validateFile(file); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := time.Now().Format("20060102150405") + "_" + utils.GenerateUUID()[:8] + ext
	savePath := config.Cfg.Upload.SavePath + "/" + filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.ErrServer(c, "文件保存失败")
		return
	}
	response.Success(c, gin.H{
		"url":  "/uploads/" + filename,
		"name": filename,
	})
}

func (h *UploadHandler) UploadLicense(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrParam(c, "请选择文件")
		return
	}
	if err := validateFile(file); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := "license_" + time.Now().Format("20060102150405") + ext
	savePath := config.Cfg.Upload.SavePath + "/" + filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.ErrServer(c, "文件保存失败")
		return
	}
	response.Success(c, gin.H{
		"url":  "/uploads/" + filename,
		"name": filename,
	})
}

func (h *UploadHandler) UploadDamage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrParam(c, "请选择文件")
		return
	}
	if err := validateFile(file); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := "damage_" + time.Now().Format("20060102150405") + ext
	savePath := config.Cfg.Upload.SavePath + "/" + filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.ErrServer(c, "文件保存失败")
		return
	}
	response.Success(c, gin.H{
		"url":  "/uploads/" + filename,
		"name": filename,
	})
}

func validateFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, e := range config.Cfg.Upload.AllowedExt {
		if e == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("不支持的文件类型")
	}
	maxSize := int64(config.Cfg.Upload.MaxSize * 1024 * 1024)
	if file.Size > maxSize {
		return errors.New("文件大小超过限制")
	}
	return nil
}
