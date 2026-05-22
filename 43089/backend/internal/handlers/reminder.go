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

type CreateReminderRequest struct {
	Title         string `json:"title" validate:"required,max=200"`
	Description   string `json:"description"`
	ReminderTime  string `json:"reminder_time" validate:"required"`
	Channel       string `json:"channel" validate:"oneof=email app"`
	ActivityID    string `json:"activity_id"`
}

func GetReminders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var reminders []models.Reminder
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Plan").
		Preload("Activity").
		Order("reminder_time ASC").
		Find(&reminders).Error; err != nil {
		logger.Errorf("Failed to get reminders: %v", err)
		utils.InternalServerError(c, "Failed to get reminders")
		return
	}

	utils.Success(c, reminders)
}

func CreateReminder(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	userID := middleware.GetCurrentUserID(c)

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	reminderTime, err := time.Parse(time.RFC3339, req.ReminderTime)
	if err != nil {
		reminderTime, err = time.Parse("2006-01-02 15:04:05", req.ReminderTime)
		if err != nil {
			utils.BadRequest(c, "Invalid reminder time format, use RFC3339 or YYYY-MM-DD HH:MM:SS")
			return
		}
	}

	reminder := models.Reminder{
		PlanID:        planUUID,
		UserID:        userID,
		Title:         req.Title,
		Description:   req.Description,
		ReminderTime:  reminderTime,
		Channel:       req.Channel,
	}

	if req.ActivityID != "" {
		activityUUID, err := uuid.Parse(req.ActivityID)
		if err != nil {
			utils.BadRequest(c, "Invalid activity ID")
			return
		}
		reminder.ActivityID = activityUUID
	}

	if err := database.DB.Create(&reminder).Error; err != nil {
		logger.Errorf("Failed to create reminder: %v", err)
		utils.InternalServerError(c, "Failed to create reminder")
		return
	}

	logger.Infof("Reminder created: %s", reminder.ID)
	utils.Created(c, reminder)
}

func UpdateReminder(c *gin.Context) {
	reminderID := c.Param("id")
	reminderUUID, err := uuid.Parse(reminderID)
	if err != nil {
		utils.BadRequest(c, "Invalid reminder ID")
		return
	}

	var reminder models.Reminder
	if err := database.DB.First(&reminder, "id = ?", reminderUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Reminder not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	allowedFields := map[string]bool{
		"title":         true,
		"description":   true,
		"reminder_time":  true,
		"is_sent":       true,
		"channel":       true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if reminderTimeStr, ok := filteredUpdates["reminder_time"].(string); ok {
		reminderTime, err := time.Parse(time.RFC3339, reminderTimeStr)
		if err != nil {
			reminderTime, err = time.Parse("2006-01-02 15:04:05", reminderTimeStr)
			if err != nil {
				utils.BadRequest(c, "Invalid reminder time format")
				return
			}
		}
		filteredUpdates["reminder_time"] = reminderTime
	}

	if err := database.DB.Model(&reminder).Updates(filteredUpdates).Error; err != nil {
		logger.Errorf("Failed to update reminder: %v", err)
		utils.InternalServerError(c, "Failed to update reminder")
		return
	}

	var updatedReminder models.Reminder
	database.DB.Preload("Plan").Preload("Activity").First(&updatedReminder, reminderUUID)

	logger.Infof("Reminder updated: %s", reminderID)
	utils.Success(c, updatedReminder)
}

func DeleteReminder(c *gin.Context) {
	reminderID := c.Param("id")
	reminderUUID, err := uuid.Parse(reminderID)
	if err != nil {
		utils.BadRequest(c, "Invalid reminder ID")
		return
	}

	if err := database.DB.Delete(&models.Reminder{}, reminderUUID).Error; err != nil {
		logger.Errorf("Failed to delete reminder: %v", err)
		utils.InternalServerError(c, "Failed to delete reminder")
		return
	}

	logger.Infof("Reminder deleted: %s", reminderID)
	utils.Success(c, gin.H{"message": "Reminder deleted successfully"})
}

func GetMapData(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var activities []models.Activity
	if err := database.DB.Where("plan_id = ? AND latitude IS NOT NULL AND longitude IS NOT NULL", planUUID).
		Order("date ASC, order_index ASC").
		Find(&activities).Error; err != nil {
		logger.Errorf("Failed to get map data: %v", err)
		utils.InternalServerError(c, "Failed to get map data")
		return
	}

	type MapLocation struct {
		ID        string  `json:"id"`
		Title     string  `json:"title"`
		Type      string  `json:"type"`
		Location  string  `json:"location"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Date      string  `json:"date"`
		StartTime string  `json:"start_time"`
	}

	locations := make([]MapLocation, 0, len(activities))
	for _, activity := range activities {
		locations = append(locations, MapLocation{
			ID:        activity.ID.String(),
			Title:     activity.Title,
			Type:      activity.Type,
			Location:  activity.Location,
			Latitude:  activity.Latitude,
			Longitude: activity.Longitude,
			Date:      activity.Date.Format("2006-01-02"),
			StartTime: activity.StartTime,
		})
	}

	utils.Success(c, gin.H{
		"locations": locations,
		"total":     len(locations),
	})
}
