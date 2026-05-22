package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"travel-planner/internal/database"
	"travel-planner/internal/logger"
	"travel-planner/internal/middleware"
	"travel-planner/internal/models"
	"travel-planner/internal/utils"
)

type CreatePlanRequest struct {
	Title       string    `json:"title" validate:"required,max=200"`
	Description string    `json:"description"`
	Destination string    `json:"destination" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     string    `json:"end_date" validate:"required"`
	Budget      float64   `json:"budget"`
	Currency    string    `json:"currency"`
	CoverImage  string    `json:"cover_image"`
	IsPublic    bool      `json:"is_public"`
}

type UpdatePlanRequest struct {
	Title       string    `json:"title" validate:"max=200"`
	Description string    `json:"description"`
	Destination string    `json:"destination"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	Budget      float64   `json:"budget"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CoverImage  string    `json:"cover_image"`
	IsPublic    *bool     `json:"is_public"`
	Version     int       `json:"version"`
}

type PlanDetailResponse struct {
	models.TravelPlan
	TotalSpent     float64           `json:"total_spent"`
	ActivityCount  int64             `json:"activity_count"`
	ParticipantCount int64           `json:"participant_count"`
}

func GetPlans(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	params := utils.GetPaginationParams(c)

	query := database.DB.Model(&models.TravelPlan{}).
		Joins("LEFT JOIN plan_participants pp ON travel_plans.id = pp.plan_id").
		Where("travel_plans.owner_id = ? OR pp.user_id = ?", userID, userID).
		Group("travel_plans.id")

	if params.Search != "" {
		query = query.Where("title ILIKE ? OR destination ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	status := c.DefaultQuery("status", "")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var plans []models.TravelPlan
	query.Preload("Owner").
		Order(params.Order()).
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&plans)

	utils.Paginated(c, plans, total, params.Page, params.PageSize)
}

func GetPlan(c *gin.Context) {
	planID := c.Param("id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var plan models.TravelPlan
	if err := database.DB.Preload("Owner").
		Preload("Participants.User").
		Preload("Activities").
		Preload("Files").
		First(&plan, "id = ?", planUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Plan not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var activities []models.Activity
	database.DB.Where("plan_id = ?", planUUID).Find(&activities)
	
	var totalSpent float64
	targetCurrency := plan.Currency
	if targetCurrency == "" {
		targetCurrency = "CNY"
	}
	
	for _, activity := range activities {
		if activity.Cost <= 0 {
			continue
		}
		activityCurrency := activity.Currency
		if activityCurrency == "" {
			activityCurrency = targetCurrency
		}
		convertedCost, err := utils.ConvertCurrency(activity.Cost, activityCurrency, targetCurrency)
		if err != nil {
			logger.Warnf("Failed to convert currency for activity %s: %v", activity.ID, err)
			convertedCost = activity.Cost
		}
		totalSpent += convertedCost
	}

	var activityCount int64
	database.DB.Model(&models.Activity{}).Where("plan_id = ?", planUUID).Count(&activityCount)

	var participantCount int64
	database.DB.Model(&models.PlanParticipant{}).Where("plan_id = ?", planUUID).Count(&participantCount)

	utils.Success(c, PlanDetailResponse{
		TravelPlan:      plan,
		TotalSpent:      totalSpent,
		ActivityCount:   activityCount,
		ParticipantCount: participantCount + 1,
	})
}

func CreatePlan(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.BadRequest(c, "Invalid start date format, use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.BadRequest(c, "Invalid end date format, use YYYY-MM-DD")
		return
	}

	if endDate.Before(startDate) {
		utils.BadRequest(c, "End date must be after start date")
		return
	}

	plan := models.TravelPlan{
		Title:       req.Title,
		Description: req.Description,
		Destination: req.Destination,
		StartDate:   startDate,
		EndDate:     endDate,
		Budget:      req.Budget,
		Currency:    req.Currency,
		OwnerID:     userID,
		CoverImage:  req.CoverImage,
		IsPublic:    req.IsPublic,
		Status:      "draft",
	}

	if err := database.DB.Create(&plan).Error; err != nil {
		logger.Errorf("Failed to create plan: %v", err)
		utils.InternalServerError(c, "Failed to create plan")
		return
	}

	logger.Infof("Plan created: %s by user %s", plan.ID, userID)
	utils.Created(c, plan)
}

func UpdatePlan(c *gin.Context) {
	planID := c.Param("id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var plan models.TravelPlan
	if err := database.DB.First(&plan, "id = ?", planUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Plan not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var req UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	if req.Version > 0 && req.Version != plan.Version {
		utils.ErrorWithDetails(c, http.StatusConflict, "Conflict detected", gin.H{
			"message":      "This plan has been modified by another user. Please refresh and try again.",
			"your_version": req.Version,
			"current_version": plan.Version,
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Destination != "" {
		updates["destination"] = req.Destination
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			utils.BadRequest(c, "Invalid start date format, use YYYY-MM-DD")
			return
		}
		updates["start_date"] = startDate
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			utils.BadRequest(c, "Invalid end date format, use YYYY-MM-DD")
			return
		}
		updates["end_date"] = endDate
	}
	if req.Budget >= 0 {
		updates["budget"] = req.Budget
	}
	if req.Currency != "" {
		updates["currency"] = req.Currency
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	updates["version"] = plan.Version + 1

	result := database.DB.Model(&plan).Where("version = ?", plan.Version).Updates(updates)
	if result.Error != nil {
		logger.Errorf("Failed to update plan: %v", result.Error)
		utils.InternalServerError(c, "Failed to update plan")
		return
	}

	if result.RowsAffected == 0 {
		var currentPlan models.TravelPlan
		database.DB.First(&currentPlan, planUUID)
		utils.ErrorWithDetails(c, http.StatusConflict, "Conflict detected", gin.H{
			"message":         "This plan has been modified by another user. Please refresh and try again.",
			"current_version": currentPlan.Version,
		})
		return
	}

	var updatedPlan models.TravelPlan
	database.DB.First(&updatedPlan, planUUID)

	logger.Infof("Plan updated: %s, version: %d -> %d", planID, plan.Version, updatedPlan.Version)
	utils.Success(c, updatedPlan)
}

func DeletePlan(c *gin.Context) {
	planID := c.Param("id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	tx := database.BeginTransaction()
	if tx.Error != nil {
		utils.InternalServerError(c, "Database error")
		return
	}

	if err := tx.Where("plan_id = ?", planUUID).Delete(&models.Activity{}).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete activities: %v", err)
		utils.InternalServerError(c, "Failed to delete plan")
		return
	}

	if err := tx.Where("plan_id = ?", planUUID).Delete(&models.PlanParticipant{}).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete participants: %v", err)
		utils.InternalServerError(c, "Failed to delete plan")
		return
	}

	if err := tx.Where("plan_id = ?", planUUID).Delete(&models.File{}).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete files: %v", err)
		utils.InternalServerError(c, "Failed to delete plan")
		return
	}

	if err := tx.Where("plan_id = ?", planUUID).Delete(&models.Checklist{}).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete checklists: %v", err)
		utils.InternalServerError(c, "Failed to delete plan")
		return
	}

	if err := tx.Where("plan_id = ?", planUUID).Delete(&models.Reminder{}).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete reminders: %v", err)
		utils.InternalServerError(c, "Failed to delete plan")
		return
	}

	if err := tx.Delete(&models.TravelPlan{}, planUUID).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete plan: %v", err)
		utils.InternalServerError(c, "Failed to delete plan")
		return
	}

	if err := tx.Commit().Error; err != nil {
		logger.Errorf("Failed to commit transaction: %v", err)
		utils.InternalServerError(c, "Failed to delete plan")
		return
	}

	logger.Infof("Plan deleted: %s", planID)
	utils.Success(c, gin.H{"message": "Plan deleted successfully"})
}

func AddParticipant(c *gin.Context) {
	planID := c.Param("id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var req struct {
		UserID    string `json:"user_id" validate:"required"`
		Role      string `json:"role" validate:"oneof=owner editor viewer"`
		CanEdit   bool   `json:"can_edit"`
		CanDelete bool   `json:"can_delete"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userUUID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	var existing models.PlanParticipant
	if err := database.DB.Where("plan_id = ? AND user_id = ?", planUUID, userUUID).First(&existing).Error; err == nil {
		utils.BadRequest(c, "User is already a participant")
		return
	}

	participant := models.PlanParticipant{
		PlanID:    planUUID,
		UserID:    userUUID,
		Role:      req.Role,
		CanEdit:   req.CanEdit,
		CanDelete: req.CanDelete,
	}

	if err := database.DB.Create(&participant).Error; err != nil {
		logger.Errorf("Failed to add participant: %v", err)
		utils.InternalServerError(c, "Failed to add participant")
		return
	}

	logger.Infof("Participant added to plan %s: %s", planID, req.UserID)
	utils.Created(c, participant)
}

func RemoveParticipant(c *gin.Context) {
	planID := c.Param("id")
	participantID := c.Param("participant_id")

	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	participantUUID, err := uuid.Parse(participantID)
	if err != nil {
		utils.BadRequest(c, "Invalid participant ID")
		return
	}

	if err := database.DB.Where("plan_id = ? AND id = ?", planUUID, participantUUID).Delete(&models.PlanParticipant{}).Error; err != nil {
		logger.Errorf("Failed to remove participant: %v", err)
		utils.InternalServerError(c, "Failed to remove participant")
		return
	}

	logger.Infof("Participant removed from plan %s: %s", planID, participantID)
	utils.Success(c, gin.H{"message": "Participant removed successfully"})
}
