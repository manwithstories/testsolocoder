package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type VideoSessionRequest struct {
	BookingID uuid.UUID `json:"bookingId" binding:"required"`
}

type StartSessionRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

type EndSessionRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

func CreateVideoSession(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	var req VideoSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	if err := database.DB.Where("id = ?", req.BookingID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	uid := userID.(uuid.UUID)

	if userRole == models.RoleStudent && booking.StudentID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}
	if userRole == models.RoleTeacher && booking.TeacherID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var existingSession models.VideoSession
	if database.DB.Where("booking_id = ?", req.BookingID).First(&existingSession).Error == nil {
		c.JSON(http.StatusOK, gin.H{"session": existingSession})
		return
	}

	sessionID := fmt.Sprintf("session-%s", uuid.New().String()[:8])
	roomName := fmt.Sprintf("room-%s", uuid.New().String()[:8])

	session := models.VideoSession{
		BookingID: req.BookingID,
		SessionID: sessionID,
		RoomName:  roomName,
		Token:     fmt.Sprintf("token-%s", uuid.New().String()),
		JoinURL:   fmt.Sprintf("https://video.example.com/join/%s", roomName),
		Status:    "scheduled",
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video session"})
		return
	}

	log := models.VideoSessionLog{
		SessionID: session.ID,
		EventType: "session_created",
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Video session created for booking %s", req.BookingID),
		UserID:    &uid,
	}
	database.DB.Create(&log)

	c.JSON(http.StatusCreated, gin.H{"message": "Video session created", "session": session})
}

func GetVideoSession(c *gin.Context) {
	id := c.Param("id")

	var session models.VideoSession
	if err := database.DB.Preload("Logs").Where("id = ?", id).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func GetSessionByBooking(c *gin.Context) {
	bookingID := c.Param("bookingId")

	var session models.VideoSession
	if err := database.DB.Preload("Logs").Where("booking_id = ?", bookingID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func StartVideoSession(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var session models.VideoSession
	if err := database.DB.Where("session_id = ?", req.SessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	uid := userID.(uuid.UUID)
	now := time.Now()
	session.Status = "in_progress"
	session.ActualStartAt = &now

	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}

	log := models.VideoSessionLog{
		SessionID: session.ID,
		EventType: "session_started",
		Timestamp: now,
		Details:   "Video session started",
		UserID:    &uid,
	}
	database.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{"message": "Session started", "session": session})
}

func EndVideoSession(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req EndSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var session models.VideoSession
	if err := database.DB.Where("session_id = ?", req.SessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.Status != "in_progress" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is not in progress"})
		return
	}

	uid := userID.(uuid.UUID)
	now := time.Now()
	duration := 0
	if session.ActualStartAt != nil {
		duration = int(now.Sub(*session.ActualStartAt).Minutes())
	}

	session.Status = "completed"
	session.ActualEndAt = &now
	session.ActualDuration = duration
	session.QualityScore = calculateQualityScore(&session)

	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end session"})
		return
	}

	log := models.VideoSessionLog{
		SessionID: session.ID,
		EventType: "session_ended",
		Timestamp: now,
		Details:   fmt.Sprintf("Video session ended. Duration: %d minutes", duration),
		UserID:    &uid,
	}
	database.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{"message": "Session ended", "session": session})
}

func HandleSessionEvent(c *gin.Context) {
	sessionID := c.Param("sessionId")
	eventType := c.Query("event")

	var session models.VideoSession
	if err := database.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	log := models.VideoSessionLog{
		SessionID: session.ID,
		EventType: eventType,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Event: %s", eventType),
	}
	database.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{"message": "Event recorded"})
}

func GetSessionQuality(c *gin.Context) {
	id := c.Param("id")

	var session models.VideoSession
	if err := database.DB.Preload("Logs").Where("id = ?", id).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	qualityReport := map[string]interface{}{
		"sessionId":     session.SessionID,
		"status":        session.Status,
		"actualDuration": session.ActualDuration,
		"qualityScore":  session.QualityScore,
		"totalEvents":   len(session.Logs),
		"reconnects":    countReconnects(&session),
	}

	c.JSON(http.StatusOK, qualityReport)
}

func calculateQualityScore(session *models.VideoSession) float64 {
	score := 100.0

	reconnectCount := countReconnects(session)
	score -= float64(reconnectCount) * 5

	if session.ActualDuration > 0 {
		expectedDuration := 60
		completionRate := float64(session.ActualDuration) / float64(expectedDuration)
		if completionRate < 0.8 {
			score -= 10
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}

func countReconnects(session *models.VideoSession) int {
	count := 0
	for _, log := range session.Logs {
		if log.EventType == "reconnect" {
			count++
		}
	}
	return count
}
