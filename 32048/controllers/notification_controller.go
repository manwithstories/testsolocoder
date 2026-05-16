package controllers

import (
	"net/http"
	"secondhand-trading/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	userID := c.GetUint("userID")
	unreadOnly := c.Query("unread") == "true"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := models.DB.Model(&models.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	var total int64
	query.Count(&total)

	var notifications []models.Notification
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications)

	var unreadCount int64
	models.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"unread_count":  unreadCount,
		"page":          page,
		"page_size":     pageSize,
		"total_page":    (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func MarkNotificationRead(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var notification models.Notification
	if err := models.DB.First(&notification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to mark this notification"})
		return
	}

	if err := models.DB.Model(&notification).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func MarkAllNotificationsRead(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func DeleteNotification(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var notification models.Notification
	if err := models.DB.First(&notification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this notification"})
		return
	}

	if err := models.DB.Delete(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
