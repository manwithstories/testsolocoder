package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

type EvaluationHandler struct{}

func NewEvaluationHandler() *EvaluationHandler {
	return &EvaluationHandler{}
}

type CreateEvaluationRequest struct {
	JobID       uuid.UUID `json:"job_id" binding:"required"`
	ToUserID    uuid.UUID `json:"to_user_id" binding:"required"`
	Rating      int       `json:"rating" binding:"required,min=1,max=5"`
	Content     string    `json:"content"`
	Tags        string    `json:"tags"`
	Type        string    `json:"type" binding:"required,oneof=employer_to_temp temp_to_employer"`
	IsAnonymous bool      `json:"is_anonymous"`
}

type UpdateEvaluationRequest struct {
	Rating  int    `json:"rating" binding:"min=1,max=5"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

func (h *EvaluationHandler) CreateEvaluation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	if req.ToUserID == userID.(uuid.UUID) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Cannot evaluate yourself",
		})
		return
	}

	var job models.JobPosting
	if err := database.DB.First(&job, req.JobID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Job posting not found",
		})
		return
	}

	var existing models.Evaluation
	if err := database.DB.Where("job_id = ? AND from_user_id = ? AND to_user_id = ? AND type = ?",
		req.JobID, userID.(uuid.UUID), req.ToUserID, req.Type).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, models.Response{
			Code:    409,
			Message: "You have already submitted an evaluation for this job",
		})
		return
	}

	evaluation := models.Evaluation{
		JobID:       req.JobID,
		FromUserID:  userID.(uuid.UUID),
		ToUserID:    req.ToUserID,
		Rating:      req.Rating,
		Content:     req.Content,
		Tags:        req.Tags,
		Type:        req.Type,
		IsAnonymous: req.IsAnonymous,
	}

	if err := database.DB.Create(&evaluation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to create evaluation: " + err.Error(),
		})
		return
	}

	h.updateCreditScore(req.ToUserID, req.Rating)

	database.DB.Preload("FromUser").Preload("ToUser").First(&evaluation, evaluation.ID)

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Evaluation submitted successfully",
		Data:    evaluation,
	})
}

func (h *EvaluationHandler) updateCreditScore(userID uuid.UUID, newRating int) {
	var evaluations []models.Evaluation
	database.DB.Where("to_user_id = ?", userID).Find(&evaluations)

	if len(evaluations) == 0 {
		return
	}

	var totalRating int
	for _, e := range evaluations {
		totalRating += e.Rating
	}

	avgRating := float64(totalRating) / float64(len(evaluations))
	creditScore := int(avgRating * 20)

	if creditScore < 0 {
		creditScore = 0
	}
	if creditScore > 100 {
		creditScore = 100
	}

	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("credit_score", creditScore)
}

func (h *EvaluationHandler) GetEvaluations(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var evaluations []models.Evaluation
	var total int64

	query := database.DB.Model(&models.Evaluation{})

	if toUserID := c.Query("to_user_id"); toUserID != "" {
		uid, err := uuid.Parse(toUserID)
		if err == nil {
			query = query.Where("to_user_id = ?", uid)
		}
	}
	if fromUserID := c.Query("from_user_id"); fromUserID != "" {
		uid, err := uuid.Parse(fromUserID)
		if err == nil {
			query = query.Where("from_user_id = ?", uid)
		}
	}
	if jobID := c.Query("job_id"); jobID != "" {
		uid, err := uuid.Parse(jobID)
		if err == nil {
			query = query.Where("job_id = ?", uid)
		}
	}
	if evalType := c.Query("type"); evalType != "" {
		query = query.Where("type = ?", evalType)
	}
	if minRating := c.Query("min_rating"); minRating != "" {
		query = query.Where("rating >= ?", minRating)
	}

	query.Count(&total)
	query.Preload("FromUser").
		Preload("ToUser").
		Preload("JobPosting").
		Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&evaluations)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       evaluations,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

func (h *EvaluationHandler) GetMyEvaluations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	pagination := utils.GetPagination(c)

	var evaluations []models.Evaluation
	var total int64

	query := database.DB.Model(&models.Evaluation{}).Where("to_user_id = ?", userID.(uuid.UUID))

	if evalType := c.Query("type"); evalType != "" {
		query = query.Where("type = ?", evalType)
	}

	query.Count(&total)
	query.Preload("FromUser").
		Preload("JobPosting").
		Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&evaluations)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       evaluations,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

func (h *EvaluationHandler) GetEvaluationStats(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	userID := id.(uuid.UUID)

	var evaluations []models.Evaluation
	database.DB.Where("to_user_id = ?", userID).Find(&evaluations)

	if len(evaluations) == 0 {
		c.JSON(http.StatusOK, models.Response{
			Code:    200,
			Message: "Success",
			Data: gin.H{
				"total_count":   0,
				"average_rating": 0,
				"rating_distribution": map[string]int{
					"5": 0, "4": 0, "3": 0, "2": 0, "1": 0,
				},
			},
		})
		return
	}

	var totalRating int
	ratingDist := map[string]int{"5": 0, "4": 0, "3": 0, "2": 0, "1": 0}

	for _, e := range evaluations {
		totalRating += e.Rating
		ratingDist[string(rune('0'+e.Rating))]++
	}

	avgRating := float64(totalRating) / float64(len(evaluations))

	var user models.User
	database.DB.First(&user, userID)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"total_count":         len(evaluations),
			"average_rating":      avgRating,
			"credit_score":        user.CreditScore,
			"rating_distribution": ratingDist,
		},
	})
}

func (h *EvaluationHandler) UpdateEvaluation(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	evalID := id.(uuid.UUID)
	userID, _ := c.Get("user_id")

	var evaluation models.Evaluation
	if err := database.DB.First(&evaluation, evalID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Evaluation not found",
		})
		return
	}

	if evaluation.FromUserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "You don't have permission to update this evaluation",
		})
		return
	}

	if time.Since(evaluation.CreatedAt).Hours() > 72 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Evaluation can only be updated within 72 hours",
		})
		return
	}

	var req UpdateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Rating > 0 {
		updates["rating"] = req.Rating
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}

	if len(updates) > 0 {
		database.DB.Model(&evaluation).Updates(updates)
	}

	database.DB.Preload("FromUser").Preload("ToUser").First(&evaluation, evalID)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Evaluation updated successfully",
		Data:    evaluation,
	})
}
