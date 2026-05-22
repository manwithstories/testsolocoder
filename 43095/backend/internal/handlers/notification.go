package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		service: services.NewNotificationService(),
	}
}

func RegisterNotificationRoutes(api *gin.RouterGroup) {
	handler := NewNotificationHandler()

	notifications := api.Group("/notifications")
	notifications.Use(middleware.Auth())
	{
		notifications.GET("", handler.List)
		notifications.GET("/unread-count", handler.GetUnreadCount)
		notifications.POST("", middleware.AdminRequired(), handler.Send)
		notifications.PUT("/:id/read", handler.MarkAsRead)
		notifications.PUT("/read-all", handler.MarkAllAsRead)
		notifications.DELETE("/:id", handler.Delete)
		notifications.POST("/send-reminders", middleware.AdminRequired(), handler.SendReminders)
	}
}

func (h *NotificationHandler) List(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var isRead *bool
	isReadStr := c.Query("is_read")
	if isReadStr != "" {
		read, err := strconv.ParseBool(isReadStr)
		if err == nil {
			isRead = &read
		}
	}

	notifications, total, err := h.service.GetNotificationList(user.UserID, page, pageSize, isRead)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.SuccessWithPagination(c, notifications, total, page, pageSize)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	count, err := h.service.GetUnreadCount(user.UserID)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"unread_count": count,
	})
}

func (h *NotificationHandler) Send(c *gin.Context) {
	var req services.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	notification, err := h.service.SendInAppNotification(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, notification)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的通知ID")
		return
	}

	user := utils.GetCurrentUser(c)
	if err := h.service.MarkAsRead(user.UserID, uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	user := utils.GetCurrentUser(c)
	if err := h.service.MarkAllAsRead(user.UserID); err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的通知ID")
		return
	}

	user := utils.GetCurrentUser(c)
	if err := h.service.DeleteNotification(user.UserID, uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *NotificationHandler) SendReminders(c *gin.Context) {
	if err := h.service.SendAppointmentReminders(); err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "就诊提醒已发送",
	})
}
