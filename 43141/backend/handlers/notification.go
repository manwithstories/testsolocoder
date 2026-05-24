package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/models"
	"sports-league/pkg/auth"
)

type NotificationHandler struct {
	db *gorm.DB
}

func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{db: db}
}

func (h *NotificationHandler) List(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*auth.Claims)
	var notifications []models.Notification
	h.db.Where("user_id = ?", userClaims.UserID).Order("created_at DESC").Find(&notifications)
	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var n models.Notification
	if err := h.db.First(&n, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	n.IsRead = true
	h.db.Save(&n)
	c.JSON(http.StatusOK, n)
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*auth.Claims)
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userClaims.UserID, false).Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"message": "all marked read"})
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*auth.Claims)
	var count int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userClaims.UserID, false).Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
}
