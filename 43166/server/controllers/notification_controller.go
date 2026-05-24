package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"business-registration-platform/models"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		notificationService: services.NewNotificationService(),
	}
}

func (ctrl *NotificationController) GetUserNotifications(c *gin.Context) {
	userID, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var isRead *bool
	if v := c.Query("isRead"); v != "" {
		b, _ := strconv.ParseBool(v)
		isRead = &b
	}

	notifications, total, err := ctrl.notificationService.GetUserNotifications(userID.(uint), page, pageSize, isRead)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     notifications,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("userID")

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid notification ID")
		return
	}

	if err := ctrl.notificationService.MarkAsRead(userID.(uint), uint(notificationID)); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID, _ := c.Get("userID")

	if err := ctrl.notificationService.MarkAllAsRead(userID.(uint)); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *NotificationController) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("userID")

	count, err := ctrl.notificationService.GetUnreadCount(userID.(uint))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"count": count})
}

func (ctrl *NotificationController) GetNotificationTemplates(c *gin.Context) {
	templates, err := ctrl.notificationService.GetNotificationTemplates()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, templates)
}

func (ctrl *NotificationController) CreateNotificationTemplate(c *gin.Context) {
	var template struct {
		Code      string `json:"code" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Variables string `json:"variables"`
		IsActive  bool   `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&template); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	dbTemplate := &models.NotificationTemplate{
		Code:      template.Code,
		Name:      template.Name,
		Type:      template.Type,
		Title:     template.Title,
		Content:   template.Content,
		Variables: template.Variables,
		IsActive:  template.IsActive,
	}

	if err := ctrl.notificationService.CreateNotificationTemplate(dbTemplate); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, dbTemplate)
}

func (ctrl *NotificationController) UpdateNotificationTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid template ID")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.notificationService.UpdateNotificationTemplate(uint(id), data); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *NotificationController) DeleteNotificationTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid template ID")
		return
	}

	if err := ctrl.notificationService.DeleteNotificationTemplate(uint(id)); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *NotificationController) SendNotification(c *gin.Context) {
	var req services.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.notificationService.SendNotification(&req); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
