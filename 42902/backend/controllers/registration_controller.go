package controllers

import (
	"activity-management/models"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterEvent(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, _ := c.Get("userID")

	var event models.Event
	if err := models.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if time.Now().After(event.Deadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration deadline has passed"})
		return
	}

	var count int64
	models.DB.Model(&models.Registration{}).Where("event_id = ? AND status = ?", eventID, "registered").Count(&count)
	if count >= int64(event.Capacity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event is full"})
		return
	}

	var existingReg models.Registration
	err = models.DB.Where("user_id = ? AND event_id = ?", userID, eventID).First(&existingReg).Error
	if err == nil {
		if existingReg.Status == "registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Already registered for this event"})
			return
		}
		if existingReg.CancelCount >= 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You have cancelled registration too many times for this event"})
			return
		}
		existingReg.Status = "registered"
		existingReg.UpdatedAt = time.Now()
		models.DB.Save(&existingReg)
		c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
		return
	}

	reg := models.Registration{
		UserID:    userID.(uint),
		EventID:   uint(eventID),
		Status:    "registered",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := models.DB.Create(&reg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func CancelRegistration(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, _ := c.Get("userID")

	var reg models.Registration
	if err := models.DB.Where("user_id = ? AND event_id = ?", userID, eventID).First(&reg).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registration not found"})
		return
	}

	if reg.Status != "registered" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are not registered for this event"})
		return
	}

	reg.Status = "cancelled"
	reg.CancelCount++
	reg.UpdatedAt = time.Now()
	models.DB.Save(&reg)

	c.JSON(http.StatusOK, gin.H{"message": "Registration cancelled successfully"})
}

func GetEventRegistrations(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, _ := c.Get("userID")

	var event models.Event
	if err := models.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OrganizerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view registrations"})
		return
	}

	var registrations []models.Registration
	models.DB.Preload("User").Where("event_id = ? AND status = ?", eventID, "registered").Find(&registrations)

	c.JSON(http.StatusOK, registrations)
}

func ExportRegistrations(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, _ := c.Get("userID")

	var event models.Event
	if err := models.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OrganizerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to export registrations"})
		return
	}

	var registrations []models.Registration
	models.DB.Preload("User").Where("event_id = ? AND status = ?", eventID, "registered").Find(&registrations)

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=registrations_"+strconv.Itoa(int(eventID))+".csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{"ID", "用户名", "邮箱", "报名时间"})
	for _, reg := range registrations {
		writer.Write([]string{
			strconv.Itoa(int(reg.ID)),
			reg.User.Username,
			reg.User.Email,
			reg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

func GetMyRegistrations(c *gin.Context) {
	userID, _ := c.Get("userID")

	var registrations []models.Registration
	models.DB.Preload("Event").Preload("Event.Organizer").Where("user_id = ?", userID).Order("created_at desc").Find(&registrations)

	c.JSON(http.StatusOK, registrations)
}
