package handlers

import (
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"

	"github.com/gin-gonic/gin"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) GetLogs(c *gin.Context) {
	db := database.GetDB()

	var logs []models.OperationLog
	var total int64

	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")

	module := c.Query("module")
	action := c.Query("action")

	query := db.Model(&models.OperationLog{})

	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs)

	response.Paginated(c, logs, total, page, pageSize)
}

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()

	var notifications []models.Notification
	db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&notifications)

	var unreadCount int64
	db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	response.Success(c, gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.GetUint("id")

	db := database.GetDB()

	db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true)

	response.Success(c, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()

	db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true)

	response.Success(c, gin.H{"message": "All notifications marked as read"})
}
