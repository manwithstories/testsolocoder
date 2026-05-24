package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

func GetNotifications(c *gin.Context) {
	userID, _ := c.Get("userId")
	isRead := c.Query("isRead")

	var notifications []models.Notification
	query := database.DB.Where("user_id = ?", userID)

	if isRead != "" {
		query = query.Where("is_read = ?", isRead == "true")
	}

	if err := query.Order("created_at DESC").Limit(50).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func MarkNotificationAsRead(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	if err := database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"is_read": true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	userID, _ := c.Get("userId")

	if err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{"is_read": true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func GetUnreadNotificationCount(c *gin.Context) {
	userID, _ := c.Get("userId")

	var count int64
	if err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unreadCount": count})
}

func DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}

func GetSubjects(c *gin.Context) {
	var subjects []models.Subject
	if err := database.DB.Where("is_active = ?", true).Order("sort_order").Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subjects"})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

func CreateSubject(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Category    string `json:"category"`
		Description string `json:"description"`
		IconURL     string `json:"iconUrl"`
		SortOrder   int    `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.Subject
	if database.DB.Where("name = ?", req.Name).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Subject already exists"})
		return
	}

	subject := models.Subject{
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		IconURL:     req.IconURL,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	if err := database.DB.Create(&subject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subject"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subject created", "subject": subject})
}

func UpdateSubject(c *gin.Context) {
	id := c.Param("id")

	var subject models.Subject
	if err := database.DB.Where("id = ?", id).First(&subject).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
		IconURL     string `json:"iconUrl"`
		SortOrder   int    `json:"sortOrder"`
		IsActive    bool   `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IconURL != "" {
		updates["icon_url"] = req.IconURL
	}
	updates["sort_order"] = req.SortOrder
	updates["is_active"] = req.IsActive

	if err := database.DB.Model(&subject).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subject"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subject updated successfully"})
}

func GetAdminStats(c *gin.Context) {
	var stats struct {
		TotalUsers      int64   `json:"totalUsers"`
		TotalTeachers   int64   `json:"totalTeachers"`
		TotalStudents   int64   `json:"totalStudents"`
		TotalBookings   int64   `json:"totalBookings"`
		TotalRevenue    float64 `json:"totalRevenue"`
		PendingApprovals int64  `json:"pendingApprovals"`
	}

	database.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleTeacher).Count(&stats.TotalTeachers)
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&stats.TotalStudents)
	database.DB.Model(&models.Booking{}).Count(&stats.TotalBookings)
	database.DB.Model(&models.TeacherProfile{}).Where("approval_status = ?", "pending").Count(&stats.PendingApprovals)

	var transactions []models.Transaction
	database.DB.Where("type = ?", models.TransactionTypeCommission).Find(&transactions)
	for _, t := range transactions {
		stats.TotalRevenue += t.Amount
	}

	c.JSON(http.StatusOK, stats)
}

func GetSystemLogs(c *gin.Context) {
	level := c.Query("level")
	module := c.Query("module")

	var logs []models.SystemLog
	query := database.DB

	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func GetAdminActions(c *gin.Context) {
	var actions []models.AdminAction
	if err := database.DB.Preload("Admin").Order("created_at DESC").Limit(100).Find(&actions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admin actions"})
		return
	}

	c.JSON(http.StatusOK, actions)
}
