package handlers

import (
	"net/http"
	"strconv"
	"time"

	"freelancer-management/internal/config"
	"freelancer-management/internal/database"
	"freelancer-management/internal/logger"
	"freelancer-management/internal/middleware"
	"freelancer-management/internal/models"
	"freelancer-management/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TimeEntryHandler struct {
	db *gorm.DB
}

func NewTimeEntryHandler() *TimeEntryHandler {
	return &TimeEntryHandler{db: database.GetDB()}
}

type CreateTimeEntryRequest struct {
	ProjectID   uint       `json:"project_id" binding:"required"`
	Date        string     `json:"date" binding:"required"`
	Hours       float64    `json:"hours"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	IsTimer     bool       `json:"is_timer"`
	Billable    bool       `json:"billable"`
}

type UpdateTimeEntryRequest struct {
	Hours       float64    `json:"hours"`
	Description string     `json:"description"`
	Date        string     `json:"date"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Billable    *bool      `json:"billable"`
}

type StartTimerRequest struct {
	ProjectID   uint   `json:"project_id" binding:"required"`
	Description string `json:"description"`
}

type StopTimerRequest struct {
	Description *string `json:"description"`
}

func (h *TimeEntryHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	cfg := config.GetConfig()

	var req CreateTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", req.ProjectID, userID).First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	entryDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		return
	}

	if req.Hours <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Hours must be greater than 0")
		return
	}

	var totalHours float64
	h.db.Model(&models.TimeEntry{}).
		Where("user_id = ? AND project_id = ? AND date = ?", userID, req.ProjectID, entryDate).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&totalHours)

	if totalHours+req.Hours > cfg.App.MaxHoursPerDay {
		utils.ErrorResponse(c, http.StatusBadRequest, "Daily hours limit exceeded. Maximum 24 hours per day per project.")
		return
	}

	timeEntry := models.TimeEntry{
		UserID:      userID,
		ProjectID:   req.ProjectID,
		Date:        entryDate,
		Hours:       req.Hours,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsTimer:     req.IsTimer,
		Billable:    req.Billable,
	}

	if err := h.db.Create(&timeEntry).Error; err != nil {
		logger.LogError("Failed to create time entry: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create time entry")
		return
	}

	logger.LogOperation(userID, "create_time_entry", "Time entry created for project: "+project.Name)
	utils.SuccessResponse(c, timeEntry)
}

func (h *TimeEntryHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	projectID := c.Query("project_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	offset := (page - 1) * perPage

	var entries []models.TimeEntry
	var total int64

	query := h.db.Model(&models.TimeEntry{}).Where("user_id = ?", userID)
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	query.Count(&total)
	query.Preload("Project").Offset(offset).Limit(perPage).Order("date DESC, created_at DESC").Find(&entries)

	utils.PaginatedSuccessResponse(c, entries, utils.PaginationMeta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	})
}

func (h *TimeEntryHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var entry models.TimeEntry
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Project").First(&entry).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Time entry not found")
		return
	}

	utils.SuccessResponse(c, entry)
}

func (h *TimeEntryHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	cfg := config.GetConfig()
	id, _ := strconv.Atoi(c.Param("id"))

	var entry models.TimeEntry
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&entry).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Time entry not found")
		return
	}

	var req UpdateTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Hours != 0 && req.Hours <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Hours must be greater than 0")
		return
	}

	if req.Date != "" {
		entryDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
			return
		}
		entry.Date = entryDate
	}

	if req.Hours != 0 {
		var totalHours float64
		h.db.Model(&models.TimeEntry{}).
			Where("user_id = ? AND project_id = ? AND date = ? AND id != ?", userID, entry.ProjectID, entry.Date, id).
			Select("COALESCE(SUM(hours), 0)").
			Scan(&totalHours)

		if totalHours+req.Hours > cfg.App.MaxHoursPerDay {
			utils.ErrorResponse(c, http.StatusBadRequest, "Daily hours limit exceeded. Maximum 24 hours per day per project.")
			return
		}
		entry.Hours = req.Hours
	}

	if req.Description != "" {
		entry.Description = req.Description
	}
	if req.StartTime != nil {
		entry.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		entry.EndTime = req.EndTime
	}
	if req.Billable != nil {
		entry.Billable = *req.Billable
	}

	if err := h.db.Save(&entry).Error; err != nil {
		logger.LogError("Failed to update time entry: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update time entry")
		return
	}

	logger.LogOperation(userID, "update_time_entry", "Time entry updated")
	utils.SuccessResponse(c, entry)
}

func (h *TimeEntryHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var entry models.TimeEntry
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&entry).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Time entry not found")
		return
	}

	if err := h.db.Delete(&entry).Error; err != nil {
		logger.LogError("Failed to delete time entry: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete time entry")
		return
	}

	logger.LogOperation(userID, "delete_time_entry", "Time entry deleted")
	utils.SuccessResponseWithMessage(c, "Time entry deleted successfully", nil)
}

func (h *TimeEntryHandler) StartTimer(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req StartTimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var existingTimer models.TimeEntry
	if err := h.db.Where("user_id = ? AND is_timer = ? AND end_time IS NULL", userID, true).First(&existingTimer).Error; err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "There is already an active timer. Please stop it first.")
		return
	}

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", req.ProjectID, userID).First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	now := time.Now()
	timeEntry := models.TimeEntry{
		UserID:      userID,
		ProjectID:   req.ProjectID,
		Date:        time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
		Hours:       0,
		Description: req.Description,
		StartTime:   &now,
		IsTimer:     true,
		Billable:    true,
	}

	if err := h.db.Create(&timeEntry).Error; err != nil {
		logger.LogError("Failed to start timer: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to start timer")
		return
	}

	logger.LogOperation(userID, "start_timer", "Timer started for project: "+project.Name)
	utils.SuccessResponse(c, timeEntry)
}

func (h *TimeEntryHandler) StopTimer(c *gin.Context) {
	userID := middleware.GetUserID(c)
	cfg := config.GetConfig()
	id, _ := strconv.Atoi(c.Param("id"))

	var entry models.TimeEntry
	if err := h.db.Where("id = ? AND user_id = ? AND is_timer = ? AND end_time IS NULL", id, userID, true).First(&entry).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Active timer not found")
		return
	}

	var req StopTimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = StopTimerRequest{}
	}

	now := time.Now()
	entry.EndTime = &now
	hours := now.Sub(*entry.StartTime).Hours()

	if hours <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Timer duration is too short")
		return
	}

	var totalHours float64
	h.db.Model(&models.TimeEntry{}).
		Where("user_id = ? AND project_id = ? AND date = ? AND id != ?", userID, entry.ProjectID, entry.Date, id).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&totalHours)

	if totalHours+hours > cfg.App.MaxHoursPerDay {
		utils.ErrorResponse(c, http.StatusBadRequest, "Daily hours limit exceeded. Maximum 24 hours per day per project.")
		return
	}

	entry.Hours = hours
	if req.Description != nil && *req.Description != "" {
		entry.Description = *req.Description
	}

	if err := h.db.Save(&entry).Error; err != nil {
		logger.LogError("Failed to stop timer: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to stop timer")
		return
	}

	logger.LogOperation(userID, "stop_timer", "Timer stopped, duration: "+strconv.FormatFloat(hours, 'f', 2, 64)+" hours")
	utils.SuccessResponse(c, entry)
}

func (h *TimeEntryHandler) GetActiveTimer(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var entry models.TimeEntry
	if err := h.db.Where("user_id = ? AND is_timer = ? AND end_time IS NULL", userID, true).Preload("Project").First(&entry).Error; err != nil {
		utils.SuccessResponse(c, nil)
		return
	}

	utils.SuccessResponse(c, entry)
}
