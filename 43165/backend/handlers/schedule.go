package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

type ScheduleHandler struct{}

func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{}
}

type CreateScheduleRequest struct {
	JobID       uuid.UUID `json:"job_id" binding:"required"`
	TemporaryID uuid.UUID `json:"temporary_id" binding:"required"`
	ShiftDate   string    `json:"shift_date" binding:"required"`
	StartTime   string    `json:"start_time" binding:"required"`
	EndTime     string    `json:"end_time" binding:"required"`
	Location    string    `json:"location"`
	Notes       string    `json:"notes"`
}

type UpdateScheduleRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
	Notes     string `json:"notes"`
	Status    string `json:"status"`
}

type BatchCreateScheduleRequest struct {
	JobID       uuid.UUID   `json:"job_id" binding:"required"`
	TemporaryIDs []uuid.UUID `json:"temporary_ids" binding:"required"`
	ShiftDate   string      `json:"shift_date" binding:"required"`
	StartTime   string      `json:"start_time" binding:"required"`
	EndTime     string      `json:"end_time" binding:"required"`
	Location    string      `json:"location"`
	Notes       string      `json:"notes"`
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	var existing models.Schedule
	if err := database.DB.Where("temporary_id = ? AND shift_date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
		req.TemporaryID, shiftDate, req.EndTime, req.StartTime, req.EndTime, req.StartTime).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, models.Response{
			Code:    409,
			Message: "Schedule conflict: This temporary staff already has a shift at this time",
		})
		return
	}

	schedule := models.Schedule{
		JobID:       req.JobID,
		TemporaryID: req.TemporaryID,
		ShiftDate:   shiftDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Location:    req.Location,
		Notes:       req.Notes,
		Status:      "scheduled",
		CreatedBy:   userID.(uuid.UUID),
	}

	if err := database.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to create schedule: " + err.Error(),
		})
		return
	}

	database.DB.Preload("Temporary").Preload("JobPosting").First(&schedule, schedule.ID)

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Schedule created successfully",
		Data:    schedule,
	})
}

func (h *ScheduleHandler) BatchCreateSchedules(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req BatchCreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	if len(req.TemporaryIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "At least one temporary staff ID is required",
		})
		return
	}

	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	var createdSchedules []models.Schedule
	var conflicts []string

	for _, tempID := range req.TemporaryIDs {
		var existing models.Schedule
		if err := database.DB.Where("temporary_id = ? AND shift_date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
			tempID, shiftDate, req.EndTime, req.StartTime, req.EndTime, req.StartTime).First(&existing).Error; err == nil {
			conflicts = append(conflicts, tempID.String())
			continue
		}

		schedule := models.Schedule{
			JobID:       req.JobID,
			TemporaryID: tempID,
			ShiftDate:   shiftDate,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
			Location:    req.Location,
			Notes:       req.Notes,
			Status:      "scheduled",
			CreatedBy:   userID.(uuid.UUID),
		}

		if err := database.DB.Create(&schedule).Error; err == nil {
			database.DB.Preload("Temporary").First(&schedule, schedule.ID)
			createdSchedules = append(createdSchedules, schedule)
		}
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: fmt.Sprintf("Created %d schedules. Conflicts: %d", len(createdSchedules), len(conflicts)),
		Data: gin.H{
			"created":   createdSchedules,
			"conflicts": conflicts,
		},
	})
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var schedules []models.Schedule
	var total int64

	query := database.DB.Model(&models.Schedule{})

	if jobID := c.Query("job_id"); jobID != "" {
		uid, err := uuid.Parse(jobID)
		if err == nil {
			query = query.Where("job_id = ?", uid)
		}
	}
	if tempID := c.Query("temporary_id"); tempID != "" {
		uid, err := uuid.Parse(tempID)
		if err == nil {
			query = query.Where("temporary_id = ?", uid)
		}
	}
	if date := c.Query("date"); date != "" {
		query = query.Where("shift_date = ?", date)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("shift_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("shift_date <= ?", endDate)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("Temporary").
		Preload("JobPosting").
		Order("shift_date ASC, start_time ASC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&schedules)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       schedules,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	scheduleID := id.(uuid.UUID)

	var schedule models.Schedule
	if err := database.DB.Preload("Temporary").Preload("JobPosting").First(&schedule, scheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Schedule not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data:    schedule,
	})
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	scheduleID := id.(uuid.UUID)

	var schedule models.Schedule
	if err := database.DB.First(&schedule, scheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Schedule not found",
		})
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.StartTime != "" {
		updates["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		updates["end_time"] = req.EndTime
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		database.DB.Model(&schedule).Updates(updates)
	}

	database.DB.Preload("Temporary").Preload("JobPosting").First(&schedule, scheduleID)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Schedule updated successfully",
		Data:    schedule,
	})
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	scheduleID := id.(uuid.UUID)

	var schedule models.Schedule
	if err := database.DB.First(&schedule, scheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Schedule not found",
		})
		return
	}

	database.DB.Delete(&schedule)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Schedule deleted successfully",
	})
}

func (h *ScheduleHandler) GetMySchedules(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var schedules []models.Schedule
	var total int64

	query := database.DB.Model(&models.Schedule{}).Where("temporary_id = ?", userID.(uuid.UUID))

	if date := c.Query("date"); date != "" {
		query = query.Where("shift_date = ?", date)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("shift_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("shift_date <= ?", endDate)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("JobPosting").
		Preload("JobPosting.Employer").
		Order("shift_date ASC, start_time ASC").
		Find(&schedules)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       schedules,
			Total:      total,
			Page:       1,
			PageSize:   int(total),
			TotalPages: 1,
		},
	})
}

func (h *ScheduleHandler) ExportSchedules(c *gin.Context) {
	query := database.DB.Model(&models.Schedule{})

	if jobID := c.Query("job_id"); jobID != "" {
		uid, err := uuid.Parse(jobID)
		if err == nil {
			query = query.Where("job_id = ?", uid)
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("shift_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("shift_date <= ?", endDate)
	}

	var schedules []models.Schedule
	query.Preload("Temporary").Preload("JobPosting").
		Order("shift_date ASC, start_time ASC").
		Find(&schedules)

	if len(schedules) == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "No schedules found for export",
		})
		return
	}

	filename := fmt.Sprintf("schedules_%s.xlsx", time.Now().Format("20060102_150405"))
	filepath, err := utils.ExportSchedulesToExcel(schedules, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to export schedules: " + err.Error(),
		})
		return
	}

	utils.ServeFile(c, filepath, filename)
}

func (h *ScheduleHandler) CheckConflict(c *gin.Context) {
	var req struct {
		TemporaryID uuid.UUID `json:"temporary_id" binding:"required"`
		ShiftDate   string    `json:"shift_date" binding:"required"`
		StartTime   string    `json:"start_time" binding:"required"`
		EndTime     string    `json:"end_time" binding:"required"`
		ExcludeID   *uuid.UUID `json:"exclude_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid date format",
		})
		return
	}

	query := database.DB.Where("temporary_id = ? AND shift_date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
		req.TemporaryID, shiftDate, req.EndTime, req.StartTime, req.EndTime, req.StartTime)

	if req.ExcludeID != nil {
		query = query.Where("id != ?", *req.ExcludeID)
	}

	var conflicts []models.Schedule
	query.Preload("JobPosting").Find(&conflicts)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"has_conflict": len(conflicts) > 0,
			"conflicts":    conflicts,
		},
	})
}
