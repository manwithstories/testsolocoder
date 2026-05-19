package controllers

import (
	"activity-management/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time"`
	Capacity    int       `json:"capacity" binding:"required,min=1"`
	Deadline    time.Time `json:"deadline" binding:"required"`
}

func CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Deadline.After(req.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration deadline must be before event start time"})
		return
	}

	userID, _ := c.Get("userID")

	event := models.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Capacity:    req.Capacity,
		Deadline:    req.Deadline,
		OrganizerID: userID.(uint),
		CreatedAt:   time.Now(),
	}

	if err := models.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func GetEvents(c *gin.Context) {
	var events []models.Event
	models.DB.Preload("Organizer").Order("created_at desc").Find(&events)

	userID, exists := c.Get("userID")

	type EventWithStatus struct {
		models.Event
		RegisteredCount int64  `json:"registered_count"`
		IsRegistered    bool   `json:"is_registered"`
		IsFull          bool   `json:"_is_full"`
		IsDeadlinePassed bool  `json:"is_deadline_passed"`
		CanRegister     bool   `json:"can_register"`
		IsOrganizer     bool   `json:"is_organizer"`
		CancelCount     int    `json:"cancel_count"`
	}

	result := make([]EventWithStatus, len(events))
	for i, event := range events {
		var count int64
		models.DB.Model(&models.Registration{}).Where("event_id = ? AND status = ?", event.ID, "registered").Count(&count)

		isDeadlinePassed := time.Now().After(event.Deadline)
		isFull := count >= int64(event.Capacity)

		isRegistered := false
		cancelCount := 0
		isOrganizer := false

		if exists {
			var reg models.Registration
			if models.DB.Where("user_id = ? AND event_id = ?", userID, event.ID).First(&reg).Error == nil {
				isRegistered = reg.Status == "registered"
				cancelCount = reg.CancelCount
			}
			isOrganizer = event.OrganizerID == userID.(uint)
		}

		canRegister := !isDeadlinePassed && !isFull && !isRegistered && cancelCount < 3

		result[i] = EventWithStatus{
			Event:            event,
			RegisteredCount:  count,
			IsRegistered:     isRegistered,
			IsFull:           isFull,
			IsDeadlinePassed: isDeadlinePassed,
			CanRegister:      canRegister,
			IsOrganizer:      isOrganizer,
			CancelCount:      cancelCount,
		}
	}

	c.JSON(http.StatusOK, result)
}

func GetEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	if err := models.DB.Preload("Organizer").Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var count int64
	models.DB.Model(&models.Registration{}).Where("event_id = ? AND status = ?", id, "registered").Count(&count)

	userID, exists := c.Get("userID")
	isRegistered := false
	cancelCount := 0
	isOrganizer := false

	if exists {
		var reg models.Registration
		if models.DB.Where("user_id = ? AND event_id = ?", userID, id).First(&reg).Error == nil {
			isRegistered = reg.Status == "registered"
			cancelCount = reg.CancelCount
		}
		isOrganizer = event.OrganizerID == userID.(uint)
	}

	isDeadlinePassed := time.Now().After(event.Deadline)
	isFull := count >= int64(event.Capacity)
	canRegister := !isDeadlinePassed && !isFull && !isRegistered && cancelCount < 3

	c.JSON(http.StatusOK, gin.H{
		"event":              event,
		"registered_count":   count,
		"is_registered":      isRegistered,
		"is_full":            isFull,
		"is_deadline_passed": isDeadlinePassed,
		"can_register":       canRegister,
		"is_organizer":       isOrganizer,
		"cancel_count":       cancelCount,
	})
}

func UpdateEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, _ := c.Get("userID")

	var event models.Event
	if err := models.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OrganizerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.Title = req.Title
	event.Description = req.Description
	event.Location = req.Location
	event.StartTime = req.StartTime
	event.EndTime = req.EndTime
	event.Capacity = req.Capacity
	event.Deadline = req.Deadline

	models.DB.Save(&event)
	c.JSON(http.StatusOK, event)
}

func DeleteEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, _ := c.Get("userID")

	var event models.Event
	if err := models.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OrganizerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	models.DB.Where("event_id = ?", id).Delete(&models.Registration{})
	models.DB.Delete(&event)

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
