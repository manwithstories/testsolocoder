package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"translation-platform/internal/config"
	"translation-platform/internal/database"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var AllowedExtensions = map[string]models.DocumentType{
	".doc":  models.DocTypeWord,
	".docx": models.DocTypeWord,
	".xls":  models.DocTypeExcel,
	".xlsx": models.DocTypeExcel,
	".pdf":  models.DocTypePDF,
	".ppt":  models.DocTypePPT,
	".pptx": models.DocTypePPT,
	".txt":  models.DocTypeTXT,
}

func UploadDocument(c *gin.Context) {
	projectID := c.Param("project_id")
	userID, _ := c.Get("user_id")

	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "获取上传文件失败")
		return
	}
	defer file.Close()

	maxSize := config.Cfg.Upload.MaxSize
	if header.Size > maxSize {
		utils.BadRequest(c, fmt.Sprintf("文件大小超过限制（%d MB）", maxSize/(1024*1024)))
		return
	}

	ext := filepath.Ext(header.Filename)
	docType, ok := AllowedExtensions[ext]
	if !ok {
		docType = models.DocTypeOther
	}

	safeName := utils.SanitizeFilename(header.Filename)
	newFileName := fmt.Sprintf("%d_%s_%s", project.ID, uuid.New().String(), safeName)
	uploadPath := config.Cfg.Upload.Path

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		utils.InternalError(c, "创建上传目录失败")
		return
	}

	fullPath := filepath.Join(uploadPath, newFileName)
	dst, err := os.Create(fullPath)
	if err != nil {
		utils.InternalError(c, "保存文件失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.InternalError(c, "保存文件失败")
		return
	}

	isSource := c.DefaultPostForm("is_source", "true") == "true"

	var existingDoc models.Document
	query := database.DB.Where("project_id = ? AND is_source = ?", project.ID, isSource)
	if isSource {
		query = query.Order("version DESC")
	}
	query.First(&existingDoc)

	version := 1
	if existingDoc.ID > 0 {
		version = existingDoc.Version + 1
	}

	document := models.Document{
		ProjectID:  project.ID,
		FileName:   safeName,
		FilePath:   fullPath,
		FileType:   docType,
		FileSize:   header.Size,
		Version:    version,
		IsSource:   isSource,
		UploadedBy: userID.(uint),
	}

	if err := database.DB.Create(&document).Error; err != nil {
		utils.InternalError(c, "保存文档记录失败")
		return
	}

	utils.Success(c, document)
}

func ListDocuments(c *gin.Context) {
	projectID := c.Param("project_id")

	var documents []models.Document
	query := database.DB.Where("project_id = ?", projectID)

	if isSource := c.Query("is_source"); isSource != "" {
		query = query.Where("is_source = ?", isSource == "true")
	}

	query.Order("version DESC").Find(&documents)

	utils.Success(c, documents)
}

func DownloadDocument(c *gin.Context) {
	id := c.Param("id")

	var document models.Document
	if err := database.DB.First(&document, id).Error; err != nil {
		utils.NotFound(c, "文档不存在")
		return
	}

	if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
		utils.NotFound(c, "文件不存在")
		return
	}

	c.FileAttachment(document.FilePath, document.FileName)
}

func DeleteDocument(c *gin.Context) {
	id := c.Param("id")

	var document models.Document
	if err := database.DB.First(&document, id).Error; err != nil {
		utils.NotFound(c, "文档不存在")
		return
	}

	if err := os.Remove(document.FilePath); err != nil && !os.IsNotExist(err) {
		utils.InternalError(c, "删除文件失败")
		return
	}

	database.DB.Delete(&document)
	utils.Success(c, nil)
}

func GetDocumentVersions(c *gin.Context) {
	projectID := c.Param("project_id")
	isSource := c.DefaultQuery("is_source", "true") == "true"

	var documents []models.Document
	database.DB.Where("project_id = ? AND is_source = ?", projectID, isSource).
		Order("version DESC").Find(&documents)

	utils.Success(c, documents)
}

type ExtractSegmentsRequest struct {
	MaxLen int `json:"max_len"`
}

func ExtractDocumentSegments(c *gin.Context) {
	id := c.Param("id")

	var document models.Document
	if err := database.DB.First(&document, id).Error; err != nil {
		utils.NotFound(c, "文档不存在")
		return
	}

	content, err := readFileContent(document.FilePath)
	if err != nil {
		utils.InternalError(c, "读取文件内容失败")
		return
	}

	maxLen := 500
	var req ExtractSegmentsRequest
	if c.ShouldBindJSON(&req) == nil && req.MaxLen > 0 {
		maxLen = req.MaxLen
	}

	segments := utils.ExtractSegments(content, maxLen)

	var translationSegments []models.TranslationSegment
	for _, text := range segments {
		segment := models.TranslationSegment{
			ProjectID:  document.ProjectID,
			DocumentID: document.ID,
			SourceText: text,
			Status:     "pending",
		}
		translationSegments = append(translationSegments, segment)
	}

	if len(translationSegments) > 0 {
		database.DB.Create(&translationSegments)
	}

	wordCount := utils.CountWords(content)
	database.DB.Model(&document).Update("word_count", wordCount)

	utils.Success(c, gin.H{
		"segments":  translationSegments,
		"word_count": wordCount,
		"count":     len(translationSegments),
	})
}

func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetProjectSegments(c *gin.Context) {
	projectID := c.Param("project_id")

	var segments []models.TranslationSegment
	query := database.DB.Where("project_id = ?", projectID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("id ASC").Find(&segments)

	utils.Success(c, segments)
}

func UpdateSegmentTranslation(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var segment models.TranslationSegment
	if err := database.DB.First(&segment, id).Error; err != nil {
		utils.NotFound(c, "翻译片段不存在")
		return
	}

	var req struct {
		TranslatedText string `json:"translated_text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	updates := map[string]interface{}{
		"translated_text": req.TranslatedText,
		"status":          "translated",
	}

	if err := database.DB.Model(&segment).Updates(updates).Error; err != nil {
		utils.InternalError(c, "更新翻译失败")
		return
	}

	var project models.Project
	database.DB.First(&project, segment.ProjectID)
	if project.TranslatorID != nil && *project.TranslatorID == userID.(uint) {
		database.DB.Model(&models.User{}).Where("id = ?", userID).
			Update("current_workload", gorm.Expr("current_workload + 1"))
	}

	go saveToTranslationMemory(segment)

	utils.Success(c, nil)
}

func saveToTranslationMemory(segment models.TranslationSegment) {
	if segment.TranslatedText == "" {
		return
	}

	var project models.Project
	database.DB.First(&project, segment.ProjectID)

	memory := models.TranslationMemory{
		SourceText:     segment.SourceText,
		TranslatedText: segment.TranslatedText,
		SourceLang:     project.SourceLang,
		TargetLang:     project.TargetLang,
		ProjectID:      &project.ID,
		UsageCount:     1,
	}

	var existingMemory models.TranslationMemory
	database.DB.Where("source_text = ? AND source_lang = ? AND target_lang = ?",
		segment.SourceText, project.SourceLang, project.TargetLang).
		First(&existingMemory)

	if existingMemory.ID > 0 {
		database.DB.Model(&existingMemory).Update("usage_count", gorm.Expr("usage_count + 1"))
	} else {
		database.DB.Create(&memory)
	}
}

func GetMemorySuggestions(c *gin.Context) {
	sourceText := c.Query("text")
	sourceLang := c.Query("source_lang")
	targetLang := c.Query("target_lang")

	if sourceText == "" || sourceLang == "" || targetLang == "" {
		utils.BadRequest(c, "缺少必要参数")
		return
	}

	var memories []models.TranslationMemory
	database.DB.Where("source_lang = ? AND target_lang = ?", sourceLang, targetLang).
		Order("usage_count DESC").Limit(100).Find(&memories)

	var suggestions []gin.H
	for _, mem := range memories {
		similarity := utils.CalculateSimilarity(sourceText, mem.SourceText)
		if similarity >= 0.6 {
			suggestions = append(suggestions, gin.H{
				"source_text":     mem.SourceText,
				"translated_text": mem.TranslatedText,
				"similarity":      similarity,
				"usage_count":     mem.UsageCount,
			})
		}
	}

	utils.Success(c, suggestions)
}
