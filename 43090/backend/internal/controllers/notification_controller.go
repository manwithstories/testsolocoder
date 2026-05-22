package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"auction-system/internal/middleware"
	"auction-system/internal/services"
	"auction-system/pkg/response"
)

type NotificationController struct {
	notifyService *services.NotificationService
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		notifyService: services.NewNotificationService(),
	}
}

func (ctrl *NotificationController) GetMyNotifications(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	var isRead *int
	if s, err := strconv.Atoi(c.Query("is_read")); err == nil {
		isRead = &s
	}

	notifications, total, err := ctrl.notifyService.GetUserNotifications(userID, page, pageSize, isRead)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      notifications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req struct {
		NotificationIDs []uint `json:"notification_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.notifyService.MarkAsRead(userID, req.NotificationIDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	if err := ctrl.notifyService.MarkAllAsRead(userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *NotificationController) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	count, err := ctrl.notifyService.GetUnreadCount(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"count": count})
}
