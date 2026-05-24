package handlers

import (
	"smart-energy-platform/models"
	"smart-energy-platform/utils"

	"github.com/gin-gonic/gin"
)

func ListNotifications(c *gin.Context) {
	userID := c.GetUint("userId")
	isRead := c.Query("isRead")
	notifType := c.Query("type")

	var notifications []models.Notification
	query := models.DB.Where("user_id = ?", userID)

	if isRead != "" {
		query = query.Where("is_read = ?", isRead == "true")
	}
	if notifType != "" {
		query = query.Where("type = ?", notifType)
	}

	query.Order("created_at DESC").Limit(100).Find(&notifications)

	var unreadCount int64
	models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&unreadCount)

	utils.Success(c, gin.H{
		"notifications": notifications,
		"unreadCount":   unreadCount,
	})
}

func MarkNotificationRead(c *gin.Context) {
	userID := c.GetUint("userId")
	notificationID := parseUintParam(c, "id")

	var notification models.Notification
	if err := models.DB.Where("id = ? AND user_id = ?", notificationID, userID).
		First(&notification).Error; err != nil {
		utils.NotFound(c, "Notification not found")
		return
	}

	notification.IsRead = true
	models.DB.Save(&notification)

	utils.Success(c, nil)
}

func MarkAllRead(c *gin.Context) {
	userID := c.GetUint("userId")

	models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	utils.Success(c, nil)
}

func DeleteNotification(c *gin.Context) {
	userID := c.GetUint("userId")
	notificationID := parseUintParam(c, "id")

	models.DB.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{})

	utils.Success(c, nil)
}
