package handlers

import (
	"sort"
	"translation-platform/internal/database"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

func RecommendTranslators(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := database.DB.Preload("ExpertiseTags").First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	var translators []models.User
	database.DB.Preload("LanguagePairs").Preload("ExpertiseTags").
		Where("role = ? AND status = ?", models.RoleTranslator, "active").
		Find(&translators)

	type scoredTranslator struct {
		User  models.User `json:"user"`
		Score float64     `json:"score"`
	}

	var results []scoredTranslator
	for _, t := range translators {
		score := utils.CalculateTranslatorScore(&t, &project)
		if score > 0 {
			results = append(results, scoredTranslator{
				User:  t,
				Score: score,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	limit := 10
	if len(results) > limit {
		results = results[:limit]
	}

	utils.Success(c, results)
}

func AutoAssignTranslator(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := database.DB.Preload("ExpertiseTags").First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	if project.Status != models.ProjectStatusApproved {
		utils.BadRequest(c, "项目状态不允许自动派单")
		return
	}

	var translators []models.User
	database.DB.Preload("LanguagePairs").Preload("ExpertiseTags").
		Where("role = ? AND status = ?", models.RoleTranslator, "active").
		Find(&translators)

	var bestTranslator *models.User
	bestScore := 0.0

	for i := range translators {
		score := utils.CalculateTranslatorScore(&translators[i], &project)
		if score > bestScore {
			bestScore = score
			bestTranslator = &translators[i]
		}
	}

	if bestTranslator == nil {
		utils.BadRequest(c, "没有合适的译者")
		return
	}

	updates := map[string]interface{}{
		"translator_id": bestTranslator.ID,
		"status":        models.ProjectStatusAssigned,
	}

	database.DB.Model(&project).Updates(updates)

	utils.Success(c, gin.H{
		"translator": bestTranslator,
		"score":      bestScore,
	})
}

func ListLanguagePairs(c *gin.Context) {
	var pairs []models.LanguagePair
	database.DB.Find(&pairs)
	utils.Success(c, pairs)
}

func ListExpertiseTags(c *gin.Context) {
	var tags []models.ExpertiseTag
	database.DB.Find(&tags)
	utils.Success(c, tags)
}

func CreateLanguagePair(c *gin.Context) {
	var req struct {
		SourceLang  string `json:"source_lang" binding:"required"`
		TargetLang  string `json:"target_lang" binding:"required"`
		DisplayName string `json:"display_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var existing models.LanguagePair
	if database.DB.Where("source_lang = ? AND target_lang = ?", req.SourceLang, req.TargetLang).
		First(&existing).Error == nil {
		utils.BadRequest(c, "语言对已存在")
		return
	}

	pair := models.LanguagePair{
		SourceLang:  req.SourceLang,
		TargetLang:  req.TargetLang,
		DisplayName: req.DisplayName,
	}

	if err := database.DB.Create(&pair).Error; err != nil {
		utils.InternalError(c, "创建语言对失败")
		return
	}

	utils.Success(c, pair)
}

func CreateExpertiseTag(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var existing models.ExpertiseTag
	if database.DB.Where("name = ?", req.Name).First(&existing).Error == nil {
		utils.BadRequest(c, "标签已存在")
		return
	}

	tag := models.ExpertiseTag{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		utils.InternalError(c, "创建标签失败")
		return
	}

	utils.Success(c, tag)
}
