package handlers

import (
	"translation-platform/internal/database"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListTranslationMemories(c *gin.Context) {
	var memories []models.TranslationMemory
	query := database.DB

	if sourceLang := c.Query("source_lang"); sourceLang != "" {
		query = query.Where("source_lang = ?", sourceLang)
	}
	if targetLang := c.Query("target_lang"); targetLang != "" {
		query = query.Where("target_lang = ?", targetLang)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("source_text LIKE ? OR translated_text LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Model(&models.TranslationMemory{}).Count(&total)

	page, pageSize := parsePagination(c)
	query.Order("usage_count DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&memories)

	utils.Success(c, utils.PageResult{
		List:     memories,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func CreateTranslationMemory(c *gin.Context) {
	var req struct {
		SourceText     string `json:"source_text" binding:"required"`
		TranslatedText string `json:"translated_text" binding:"required"`
		SourceLang     string `json:"source_lang" binding:"required"`
		TargetLang     string `json:"target_lang" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	memory := models.TranslationMemory{
		SourceText:     req.SourceText,
		TranslatedText: req.TranslatedText,
		SourceLang:     req.SourceLang,
		TargetLang:     req.TargetLang,
		UsageCount:     1,
	}

	if err := database.DB.Create(&memory).Error; err != nil {
		utils.InternalError(c, "创建翻译记忆失败")
		return
	}

	utils.Success(c, memory)
}

func UpdateTranslationMemory(c *gin.Context) {
	id := c.Param("id")

	var memory models.TranslationMemory
	if err := database.DB.First(&memory, id).Error; err != nil {
		utils.NotFound(c, "翻译记忆不存在")
		return
	}

	var req struct {
		SourceText     string `json:"source_text"`
		TranslatedText string `json:"translated_text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.SourceText != "" {
		updates["source_text"] = req.SourceText
	}
	if req.TranslatedText != "" {
		updates["translated_text"] = req.TranslatedText
	}

	database.DB.Model(&memory).Updates(updates)
	utils.Success(c, nil)
}

func DeleteTranslationMemory(c *gin.Context) {
	id := c.Param("id")

	var memory models.TranslationMemory
	if err := database.DB.First(&memory, id).Error; err != nil {
		utils.NotFound(c, "翻译记忆不存在")
		return
	}

	database.DB.Delete(&memory)
	utils.Success(c, nil)
}

func ListGlossaryTerms(c *gin.Context) {
	var terms []models.GlossaryTerm
	query := database.DB

	if sourceLang := c.Query("source_lang"); sourceLang != "" {
		query = query.Where("source_lang = ?", sourceLang)
	}
	if targetLang := c.Query("target_lang"); targetLang != "" {
		query = query.Where("target_lang = ?", targetLang)
	}
	if domain := c.Query("domain"); domain != "" {
		query = query.Where("domain = ?", domain)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("source_term LIKE ? OR target_term LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Model(&models.GlossaryTerm{}).Count(&total)

	page, pageSize := parsePagination(c)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&terms)

	utils.Success(c, utils.PageResult{
		List:     terms,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func CreateGlossaryTerm(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		SourceTerm   string `json:"source_term" binding:"required"`
		TargetTerm   string `json:"target_term" binding:"required"`
		SourceLang   string `json:"source_lang" binding:"required"`
		TargetLang   string `json:"target_lang" binding:"required"`
		Domain       string `json:"domain"`
		Definition   string `json:"definition"`
		PartOfSpeech string `json:"part_of_speech"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	term := models.GlossaryTerm{
		SourceTerm:   req.SourceTerm,
		TargetTerm:   req.TargetTerm,
		SourceLang:   req.SourceLang,
		TargetLang:   req.TargetLang,
		Domain:       req.Domain,
		Definition:   req.Definition,
		PartOfSpeech: req.PartOfSpeech,
		CreatedBy:    userID.(uint),
	}

	if err := database.DB.Create(&term).Error; err != nil {
		utils.InternalError(c, "创建术语失败")
		return
	}

	utils.Success(c, term)
}

func UpdateGlossaryTerm(c *gin.Context) {
	id := c.Param("id")

	var term models.GlossaryTerm
	if err := database.DB.First(&term, id).Error; err != nil {
		utils.NotFound(c, "术语不存在")
		return
	}

	var req struct {
		SourceTerm   string `json:"source_term"`
		TargetTerm   string `json:"target_term"`
		Domain       string `json:"domain"`
		Definition   string `json:"definition"`
		PartOfSpeech string `json:"part_of_speech"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.SourceTerm != "" {
		updates["source_term"] = req.SourceTerm
	}
	if req.TargetTerm != "" {
		updates["target_term"] = req.TargetTerm
	}
	if req.Domain != "" {
		updates["domain"] = req.Domain
	}
	if req.Definition != "" {
		updates["definition"] = req.Definition
	}
	if req.PartOfSpeech != "" {
		updates["part_of_speech"] = req.PartOfSpeech
	}

	database.DB.Model(&term).Updates(updates)
	utils.Success(c, nil)
}

func DeleteGlossaryTerm(c *gin.Context) {
	id := c.Param("id")

	var term models.GlossaryTerm
	if err := database.DB.First(&term, id).Error; err != nil {
		utils.NotFound(c, "术语不存在")
		return
	}

	database.DB.Delete(&term)
	utils.Success(c, nil)
}
