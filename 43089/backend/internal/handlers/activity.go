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

type CreateActivityRequest struct {
	Title        string  `json:"title" validate:"required,max=200"`
	Description  string  `json:"description"`
	Type         string  `json:"type" validate:"required,oneof=sightseeing transport accommodation food other"`
	Date         string  `json:"date" validate:"required"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	Location     string  `json:"location"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Cost         float64 `json:"cost"`
	Currency     string  `json:"currency"`
	Notes        string  `json:"notes"`
	Booked       bool    `json:"booked"`
	Confirmation string  `json:"confirmation"`
	ContactInfo  string  `json:"contact_info"`
	OrderIndex   int     `json:"order_index"`
	Version      int     `json:"version"`
}

func GetActivities(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	date := c.DefaultQuery("date", "")
	activityType := c.DefaultQuery("type", "")

	query := database.DB.Where("plan_id = ?", planUUID)

	if date != "" {
		query = query.Where("date = ?", date)
	}
	if activityType != "" {
		query = query.Where("type = ?", activityType)
	}

	var activities []models.Activity
	if err := query.Order("date ASC, order_index ASC, start_time ASC").Find(&activities).Error; err != nil {
		logger.Errorf("Failed to get activities: %v", err)
		utils.InternalServerError(c, "Failed to get activities")
		return
	}

	grouped := make(map[string][]models.Activity)
	for _, activity := range activities {
		dateStr := activity.Date.Format("2006-01-02")
		grouped[dateStr] = append(grouped[dateStr], activity)
	}

	utils.Success(c, gin.H{
		"list":    activities,
		"grouped": grouped,
	})
}

func GetActivity(c *gin.Context) {
	activityID := c.Param("id")
	activityUUID, err := uuid.Parse(activityID)
	if err != nil {
		utils.BadRequest(c, "Invalid activity ID")
		return
	}

	var activity models.Activity
	if err := database.DB.Preload("Files").First(&activity, "id = ?", activityUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Activity not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	utils.Success(c, activity)
}

func CreateActivity(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	userID := middleware.GetCurrentUserID(c)

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.BadRequest(c, "Invalid date format, use YYYY-MM-DD")
		return
	}

	var plan models.TravelPlan
	if err := database.DB.First(&plan, "id = ?", planUUID).Error; err != nil {
		utils.NotFound(c, "Plan not found")
		return
	}

	if date.Before(plan.StartDate) || date.After(plan.EndDate) {
		utils.BadRequest(c, "Activity date must be within plan date range")
		return
	}

	activity := models.Activity{
		PlanID:       planUUID,
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		Date:         date,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Location:     req.Location,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Cost:         req.Cost,
		Currency:     req.Currency,
		Notes:        req.Notes,
		Booked:       req.Booked,
		Confirmation: req.Confirmation,
		ContactInfo:  req.ContactInfo,
		OrderIndex:   req.OrderIndex,
		CreatedBy:    userID,
	}

	if err := database.DB.Create(&activity).Error; err != nil {
		logger.Errorf("Failed to create activity: %v", err)
		utils.InternalServerError(c, "Failed to create activity")
		return
	}

	logger.Infof("Activity created: %s in plan %s", activity.ID, planID)
	utils.Created(c, activity)
}

func UpdateActivity(c *gin.Context) {
	activityID := c.Param("id")
	activityUUID, err := uuid.Parse(activityID)
	if err != nil {
		utils.BadRequest(c, "Invalid activity ID")
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, "id = ?", activityUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Activity not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	if req.Version > 0 && req.Version != activity.Version {
		utils.ErrorWithDetails(c, http.StatusConflict, "Conflict detected", gin.H{
			"message":         "This activity has been modified by another user. Please refresh and try again.",
			"your_version":    req.Version,
			"current_version": activity.Version,
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
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			utils.BadRequest(c, "Invalid date format, use YYYY-MM-DD")
			return
		}
		updates["date"] = date
	}
	if req.StartTime != "" {
		updates["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		updates["end_time"] = req.EndTime
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Latitude != 0 {
		updates["latitude"] = req.Latitude
	}
	if req.Longitude != 0 {
		updates["longitude"] = req.Longitude
	}
	if req.Cost >= 0 {
		updates["cost"] = req.Cost
	}
	if req.Currency != "" {
		updates["currency"] = req.Currency
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	updates["booked"] = req.Booked
	if req.Confirmation != "" {
		updates["confirmation"] = req.Confirmation
	}
	if req.ContactInfo != "" {
		updates["contact_info"] = req.ContactInfo
	}
	if req.OrderIndex >= 0 {
		updates["order_index"] = req.OrderIndex
	}
	updates["version"] = activity.Version + 1

	result := database.DB.Model(&activity).Where("version = ?", activity.Version).Updates(updates)
	if result.Error != nil {
		logger.Errorf("Failed to update activity: %v", result.Error)
		utils.InternalServerError(c, "Failed to update activity")
		return
	}

	if result.RowsAffected == 0 {
		var currentActivity models.Activity
		database.DB.First(&currentActivity, activityUUID)
		utils.ErrorWithDetails(c, http.StatusConflict, "Conflict detected", gin.H{
			"message":         "This activity has been modified by another user. Please refresh and try again.",
			"current_version": currentActivity.Version,
		})
		return
	}

	var updatedActivity models.Activity
	database.DB.First(&updatedActivity, activityUUID)

	logger.Infof("Activity updated: %s, version: %d -> %d", activityID, activity.Version, updatedActivity.Version)
	utils.Success(c, updatedActivity)
}

func DeleteActivity(c *gin.Context) {
	activityID := c.Param("id")
	activityUUID, err := uuid.Parse(activityID)
	if err != nil {
		utils.BadRequest(c, "Invalid activity ID")
		return
	}

	if err := database.DB.Delete(&models.Activity{}, activityUUID).Error; err != nil {
		logger.Errorf("Failed to delete activity: %v", err)
		utils.InternalServerError(c, "Failed to delete activity")
		return
	}

	logger.Infof("Activity deleted: %s", activityID)
	utils.Success(c, gin.H{"message": "Activity deleted successfully"})
}

func GetBudgetSummary(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var plan models.TravelPlan
	if err := database.DB.First(&plan, "id = ?", planUUID).Error; err != nil {
		utils.NotFound(c, "Plan not found")
		return
	}

	var activities []models.Activity
	if err := database.DB.Where("plan_id = ?", planUUID).Find(&activities).Error; err != nil {
		logger.Errorf("Failed to get activities for budget: %v", err)
		utils.InternalServerError(c, "Failed to calculate budget")
		return
	}

	targetCurrency := plan.Currency
	if targetCurrency == "" {
		targetCurrency = "CNY"
	}

	type CategorySum struct {
		Type           string  `json:"type"`
		Total          float64 `json:"total"`
		Count          int64   `json:"count"`
		OriginalTotal  float64 `json:"original_total"`
		OriginalCurrency string `json:"original_currency"`
	}

	categoryMap := make(map[string]*CategorySum)
	var totalSpent float64
	currencyBreakdown := make(map[string]float64)

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
			logger.Warnf("Failed to convert currency %s to %s: %v, using original value", activityCurrency, targetCurrency, err)
			convertedCost = activity.Cost
		}

		if _, ok := categoryMap[activity.Type]; !ok {
			categoryMap[activity.Type] = &CategorySum{
				Type:  activity.Type,
				Total: 0,
				Count: 0,
			}
		}
		categoryMap[activity.Type].Total += convertedCost
		categoryMap[activity.Type].Count++
		currencyBreakdown[activityCurrency] += activity.Cost

		totalSpent += convertedCost
	}

	categorySums := make([]CategorySum, 0, len(categoryMap))
	for _, sum := range categoryMap {
		categorySums = append(categorySums, *sum)
	}

	budgetRemaining := plan.Budget - totalSpent
	budgetUsage := 0.0
	if plan.Budget > 0 {
		budgetUsage = (totalSpent / plan.Budget) * 100
	}

	utils.Success(c, gin.H{
		"plan_budget":        plan.Budget,
		"plan_currency":      plan.Currency,
		"total_spent":        totalSpent,
		"budget_remaining":   budgetRemaining,
		"budget_usage":       budgetUsage,
		"by_category":        categorySums,
		"currency_breakdown": currencyBreakdown,
		"exchange_rate_note": fmt.Sprintf("All amounts converted to %s using current exchange rates", targetCurrency),
	})
}
