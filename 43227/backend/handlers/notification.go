package handlers

import (
	"net/http"
	"strconv"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	query := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID)

	if isRead := c.Query("is_read"); isRead != "" {
		query = query.Where("is_read = ?", isRead == "true")
	}

	if nType := c.Query("type"); nType != "" {
		query = query.Where("type = ?", nType)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var notifications []models.Notification
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).Find(&notifications)

	utils.SuccessWithTotal(c, notifications, total)
}

func (h *NotificationHandler) Read(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)

	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to mark notification as read")
		return
	}

	utils.Success(c, nil)
}

func (h *NotificationHandler) ReadAll(c *gin.Context) {
	userID, _ := c.Get("user_id")

	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	utils.Success(c, nil)
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var count int64
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count)

	utils.Success(c, gin.H{"count": count})
}
