package handlers

import (
	"garden-planner/database"
	"garden-planner/middleware"
	"garden-planner/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateCalendarEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	EventType   string    `json:"event_type"`
	EventDate   string    `json:"event_date" binding:"required"`
	Description string    `json:"description"`
	PlantID     *uuid.UUID `json:"plant_id"`
	PlotID      *uuid.UUID `json:"plot_id"`
}

type UpdateCalendarEventRequest struct {
	Title       string    `json:"title"`
	EventType   string    `json:"event_type"`
	EventDate   string    `json:"event_date"`
	Description string    `json:"description"`
	IsCompleted *bool     `json:"is_completed"`
}

func CreateCalendarEvent(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateCalendarEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := models.CalendarEvent{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		EventType:   req.EventType,
		Description: req.Description,
		PlantID:     req.PlantID,
		PlotID:      req.PlotID,
		IsCompleted: false,
	}

	if t, err := parseDate(req.EventDate); err == nil {
		event.EventDate = t
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create calendar event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Calendar event created successfully",
		"event":   event,
	})
}

func GetCalendarEvents(c *gin.Context) {
	userID := middleware.GetUserID(c)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	eventType := c.Query("event_type")

	var events []models.CalendarEvent

	query := database.DB.Where("user_id = ?", userID)

	if startDate != "" {
		if t, err := parseDate(startDate); err == nil {
			query = query.Where("event_date >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := parseDate(endDate); err == nil {
			query = query.Where("event_date <= ?", t)
		}
	}
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	query.Order("event_date ASC").Find(&events)

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func GetCalendarEvent(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var event models.CalendarEvent
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func UpdateCalendarEvent(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var event models.CalendarEvent
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar event not found"})
		return
	}

	var req UpdateCalendarEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		event.Title = req.Title
	}
	if req.EventType != "" {
		event.EventType = req.EventType
	}
	if req.EventDate != "" {
		if t, err := parseDate(req.EventDate); err == nil {
			event.EventDate = t
		}
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.IsCompleted != nil {
		event.IsCompleted = *req.IsCompleted
		if *req.IsCompleted {
			now := time.Now()
			event.CompletedAt = &now
		}
	}

	if err := database.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update calendar event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Calendar event updated successfully",
		"event":   event,
	})
}

func DeleteCalendarEvent(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CalendarEvent{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete calendar event"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Calendar event deleted successfully"})
}

func GetPlantingRecommendations(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	currentMonth := time.Now().Month()
	monthStr := strconv.Itoa(int(currentMonth))

	var plants []models.Plant
	query := database.DB.Where("planting_season ILIKE ? OR planting_season ILIKE ? OR planting_season ILIKE ?",
		"%"+monthStr+"%", "%全年%", "%四季%")

	if user.ClimateZone != "" {
		query = query.Where("climate_zone ILIKE ?", "%"+user.ClimateZone+"%")
	}

	query.Limit(10).Find(&plants)

	c.JSON(http.StatusOK, gin.H{
		"recommendations": plants,
		"month":           currentMonth,
		"climate_zone":    user.ClimateZone,
	})
}
