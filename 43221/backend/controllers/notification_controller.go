package controllers

import (
	"net/http"
	"strconv"

	"consultation-platform/config"
	"consultation-platform/services"
	"consultation-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController(cfg *config.Config) *NotificationController {
	return &NotificationController{
		notificationService: services.NewNotificationService(cfg),
	}
}

func (ctrl *NotificationController) GetNotifications(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	onlyUnread := c.DefaultQuery("only_unread", "false") == "true"

	notifications, total, err := ctrl.notificationService.GetNotifications(userID, page, pageSize, onlyUnread)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    notifications,
	})
}

func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid notification ID")
		return
	}

	err = ctrl.notificationService.MarkAsRead(userID, notificationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	err = ctrl.notificationService.MarkAllAsRead(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *NotificationController) GetUnreadCount(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	count, err := ctrl.notificationService.GetUnreadCount(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"unread_count": count})
}

func (ctrl *NotificationController) GetTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	templates, total, err := ctrl.notificationService.GetTemplates(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    templates,
	})
}

func (ctrl *NotificationController) UpdateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid template ID")
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	err = ctrl.notificationService.UpdateTemplate(id, req.Title, req.Content)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

type NotificationControllerInterface interface {
	GetNotifications(c *gin.Context)
	MarkAsRead(c *gin.Context)
	MarkAllAsRead(c *gin.Context)
	GetUnreadCount(c *gin.Context)
	GetTemplates(c *gin.Context)
	UpdateTemplate(c *gin.Context)
}
