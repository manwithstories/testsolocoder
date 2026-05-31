package handlers

import (
	"translation-platform/internal/database"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListReviewTasks(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var tasks []models.ReviewTask
	query := database.DB.Preload("Project").Preload("Reviewer").Preload("Segment")

	if projectID := c.Query("project_id"); projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("reviewer_id = ?", userID)
	}

	var total int64
	query.Model(&models.ReviewTask{}).Count(&total)

	page, pageSize := parsePagination(c)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks)

	utils.Success(c, utils.PageResult{
		List:     tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func CreateReviewTask(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		ProjectID  uint   `json:"project_id" binding:"required"`
		SegmentID  uint   `json:"segment_id"`
		ReviewerID uint   `json:"reviewer_id" binding:"required"`
		Comment    string `json:"comment"`
		Suggestion string `json:"suggestion"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	var project models.Project
	if err := database.DB.First(&project, req.ProjectID).Error; err != nil {
		utils.NotFound(c, "项目不存在")
		return
	}

	var maxRound int64
	database.DB.Model(&models.ReviewTask{}).
		Where("project_id = ?", req.ProjectID).
		Select("COALESCE(MAX(round), 0)").
		Scan(&maxRound)

	task := models.ReviewTask{
		ProjectID:  req.ProjectID,
		ReviewerID: req.ReviewerID,
		SegmentID:  req.SegmentID,
		Round:      int(maxRound) + 1,
		Comment:    req.Comment,
		Suggestion: req.Suggestion,
		Status:     models.ReviewStatusPending,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		utils.InternalError(c, "创建审核任务失败")
		return
	}

	_ = userID
	utils.Success(c, task)
}

func ProcessReview(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var task models.ReviewTask
	if err := database.DB.First(&task, id).Error; err != nil {
		utils.NotFound(c, "审核任务不存在")
		return
	}

	if task.ReviewerID != userID.(uint) {
		utils.Forbidden(c, "无权处理此审核任务")
		return
	}

	var req struct {
		Status     string `json:"status" binding:"required,oneof=approved rejected"`
		Comment    string `json:"comment"`
		Suggestion string `json:"suggestion"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	tx := database.DB.Begin()

	updates := map[string]interface{}{
		"status":     req.Status,
		"comment":    req.Comment,
		"suggestion": req.Suggestion,
	}
	tx.Model(&task).Updates(updates)

	if task.SegmentID > 0 {
		tx.Model(&models.TranslationSegment{}).
			Where("id = ?", task.SegmentID).
			Update("review_status", req.Status)

		if req.Suggestion != "" {
			tx.Model(&models.TranslationSegment{}).
				Where("id = ?", task.SegmentID).
				Update("translated_text", req.Suggestion)
		}
	}

	tx.Commit()

	checkProjectReviewStatus(task.ProjectID)

	utils.Success(c, nil)
}

func checkProjectReviewStatus(projectID uint) {
	var pendingCount int64
	database.DB.Model(&models.ReviewTask{}).
		Where("project_id = ? AND status = ?", projectID, models.ReviewStatusPending).
		Count(&pendingCount)

	if pendingCount == 0 {
		var rejectedCount int64
		database.DB.Model(&models.ReviewTask{}).
			Where("project_id = ? AND status = ?", projectID, models.ReviewStatusRejected).
			Count(&rejectedCount)

		if rejectedCount == 0 {
			database.DB.Model(&models.Project{}).
				Where("id = ?", projectID).
				Update("review_status", models.ReviewStatusApproved)
		}
	}
}

func BatchReview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		TaskIDs    []uint `json:"task_ids" binding:"required"`
		Status     string `json:"status" binding:"required,oneof=approved rejected"`
		Comment    string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	tx := database.DB.Begin()

	for _, taskID := range req.TaskIDs {
		var task models.ReviewTask
		if err := tx.First(&task, taskID).Error; err != nil {
			continue
		}
		if task.ReviewerID != userID.(uint) {
			continue
		}

		updates := map[string]interface{}{
			"status":  req.Status,
			"comment": req.Comment,
		}
		tx.Model(&task).Updates(updates)

		if task.SegmentID > 0 {
			tx.Model(&models.TranslationSegment{}).
				Where("id = ?", task.SegmentID).
				Update("review_status", req.Status)
		}
	}

	tx.Commit()

	utils.Success(c, nil)
}

func GetProjectReviewSummary(c *gin.Context) {
	id := c.Param("id")

	var totalCount, pendingCount, approvedCount, rejectedCount int64

	database.DB.Model(&models.ReviewTask{}).
		Where("project_id = ?", id).Count(&totalCount)
	database.DB.Model(&models.ReviewTask{}).
		Where("project_id = ? AND status = ?", id, models.ReviewStatusPending).Count(&pendingCount)
	database.DB.Model(&models.ReviewTask{}).
		Where("project_id = ? AND status = ?", id, models.ReviewStatusApproved).Count(&approvedCount)
	database.DB.Model(&models.ReviewTask{}).
		Where("project_id = ? AND status = ?", id, models.ReviewStatusRejected).Count(&rejectedCount)

	utils.Success(c, gin.H{
		"total":     totalCount,
		"pending":   pendingCount,
		"approved":  approvedCount,
		"rejected":  rejectedCount,
	})
}
