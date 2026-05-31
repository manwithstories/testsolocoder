package handlers

import (
	"time"
	"translation-platform/internal/database"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateProjectRequest struct {
	Title         string              `json:"title" binding:"required,max=200"`
	Description   string              `json:"description"`
	SourceLang    string              `json:"source_lang" binding:"required"`
	TargetLang    string              `json:"target_lang" binding:"required"`
	ExpertiseTags []uint              `json:"expertise_tag_ids"`
	WordCount     int                 `json:"word_count"`
	Urgency       models.UrgencyLevel `json:"urgency" binding:"required,oneof=normal fast urgent"`
	Deadline      string              `json:"deadline" binding:"required"`
}

func CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		utils.BadRequest(c, "截止日期格式错误")
		return
	}

	domainNames := getDomainNames(req.ExpertiseTags)
	feeResult := utils.CalculateFee(req.WordCount, req.SourceLang, req.TargetLang, req.Urgency, domainNames)

	project := models.Project{
		Title:       req.Title,
		Description: req.Description,
		ClientID:    userID.(uint),
		SourceLang:  req.SourceLang,
		TargetLang:  req.TargetLang,
		WordCount:   req.WordCount,
		Urgency:     req.Urgency,
		Difficulty:  feeResult["domain_diff"],
		UnitPrice:   feeResult["base_amount"] / float64(req.WordCount),
		TotalAmount: feeResult["total_amount"],
		Deadline:    deadline,
		Status:      models.ProjectStatusPending,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "创建项目失败")
		return
	}

	if len(req.ExpertiseTags) > 0 {
		var tags []models.ExpertiseTag
		tx.Find(&tags, req.ExpertiseTags)
		tx.Model(&project).Association("ExpertiseTags").Append(tags)
	}

	tx.Commit()

	utils.Success(c, gin.H{
		"project":    project,
		"fee_detail": feeResult,
	})
}

func ListProjects(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var projects []models.Project
	query := database.DB.Preload("Client").Preload("PM").Preload("Translator").
		Preload("Reviewer").Preload("ExpertiseTags").Preload("Documents")

	switch role.(string) {
	case "client":
		query = query.Where("client_id = ?", userID)
	case "translator":
		query = query.Where("translator_id = ?", userID)
	case "pm":
		query = query.Where("pm_id = ? OR status = ?", userID, models.ProjectStatusPending)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if sourceLang := c.Query("source_lang"); sourceLang != "" {
		query = query.Where("source_lang = ?", sourceLang)
	}
	if targetLang := c.Query("target_lang"); targetLang != "" {
		query = query.Where("target_lang = ?", targetLang)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Model(&models.Project{}).Count(&total)

	page, pageSize := parsePagination(c)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&projects)

	utils.Success(c, utils.PageResult{
		List:     projects,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func GetProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := database.DB.Preload("Client").Preload("PM").Preload("Translator").
		Preload("Reviewer").Preload("ExpertiseTags").Preload("Documents").
		Preload("Comments.User").First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	utils.Success(c, project)
}

func ApproveProject(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	if project.Status != models.ProjectStatusPending {
		utils.BadRequest(c, "项目状态不允许审核")
		return
	}

	pmID := userID.(uint)
	updates := map[string]interface{}{
		"status": models.ProjectStatusApproved,
		"pm_id":  pmID,
	}

	if err := database.DB.Model(&project).Updates(updates).Error; err != nil {
		utils.InternalError(c, "审核项目失败")
		return
	}

	utils.Success(c, nil)
}

func AssignTranslator(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		TranslatorID uint `json:"translator_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	if project.Status != models.ProjectStatusApproved && project.Status != models.ProjectStatusAssigned {
		utils.BadRequest(c, "项目状态不允许分配译者")
		return
	}

	var translator models.User
	if err := database.DB.Where("id = ? AND role = ?", req.TranslatorID, models.RoleTranslator).
		First(&translator).Error; err != nil {
		utils.NotFound(c, "译者不存在")
		return
	}

	updates := map[string]interface{}{
		"translator_id": req.TranslatorID,
		"status":        models.ProjectStatusAssigned,
	}

	if err := database.DB.Model(&project).Updates(updates).Error; err != nil {
		utils.InternalError(c, "分配译者失败")
		return
	}

	utils.Success(c, nil)
}

func StartProject(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	if project.TranslatorID == nil || *project.TranslatorID != userID.(uint) {
		utils.Forbidden(c, "只有项目译者可以开始翻译")
		return
	}

	if project.Status != models.ProjectStatusAssigned {
		utils.BadRequest(c, "项目状态不允许开始")
		return
	}

	database.DB.Model(&project).Update("status", models.ProjectStatusInProgress)
	utils.Success(c, nil)
}

func SubmitForReview(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	if project.TranslatorID == nil || *project.TranslatorID != userID.(uint) {
		utils.Forbidden(c, "只有项目译者可以提交审核")
		return
	}

	if project.Status != models.ProjectStatusInProgress {
		utils.BadRequest(c, "项目状态不允许提交审核")
		return
	}

	updates := map[string]interface{}{
		"status":        models.ProjectStatusReview,
		"review_status": models.ReviewStatusPending,
	}

	database.DB.Model(&project).Updates(updates)
	utils.Success(c, nil)
}

func CompleteProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	if project.ReviewStatus != models.ReviewStatusApproved {
		utils.BadRequest(c, "审核未通过，无法完成")
		return
	}

	tx := database.DB.Begin()

	tx.Model(&project).Update("status", models.ProjectStatusCompleted)

	if project.TranslatorID != nil {
		tx.Model(&models.User{}).Where("id = ?", *project.TranslatorID).
			Updates(map[string]interface{}{
				"completed_count":   gorm.Expr("completed_count + 1"),
				"current_workload":  gorm.Expr("GREATEST(current_workload - ?, 0)", project.WordCount),
			})
	}

	payment := models.Payment{
		ProjectID:     project.ID,
		ClientID:      project.ClientID,
		TranslatorID:  project.TranslatorID,
		Amount:        project.TotalAmount,
		BaseAmount:    project.TotalAmount,
		Status:        "pending",
	}
	tx.Create(&payment)

	tx.Commit()

	utils.Success(c, nil)
}

func CancelProject(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	if project.ClientID != userID.(uint) && project.Status == models.ProjectStatusPending {
		utils.BadRequest(c, "只能取消待审核的项目")
		return
	}

	database.DB.Model(&project).Update("status", models.ProjectStatusCancelled)
	utils.Success(c, nil)
}

func AddProjectComment(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	comment := models.ProjectComment{
		ProjectID: parseUintParam(id),
		UserID:    userID.(uint),
		Content:   req.Content,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		utils.InternalError(c, "添加评论失败")
		return
	}

	utils.Success(c, comment)
}

func getDomainNames(tagIDs []uint) []string {
	var tags []models.ExpertiseTag
	database.DB.Find(&tags, tagIDs)
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	return names
}

func parseUintParam(s string) uint {
	var n uint
	return n
}