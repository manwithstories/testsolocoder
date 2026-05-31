package handler

import (
	"net/http"
	"time"

	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MaintenanceHandler struct{}

func NewMaintenanceHandler() *MaintenanceHandler {
	return &MaintenanceHandler{}
}

func (h *MaintenanceHandler) CreateMaintenance(c *gin.Context) {
	var req model.CreateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	shipID, _ := uuid.Parse(req.ShipID)

	var ship model.Ship
	if err := database.DB.First(&ship, "id = ?", shipID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Ship not found")
		return
	}

	maintenance := model.MaintenanceRecord{
		ShipID:          shipID,
		MaintenanceType: req.MaintenanceType,
		Title:           req.Title,
		Description:     req.Description,
		Status:          model.MaintenanceStatusScheduled,
		PlannedDate:     req.PlannedDate,
		Cost:            req.Cost,
		Currency:        req.Currency,
		Provider:        req.Provider,
		Technician:      req.Technician,
		NextDueDate:     req.NextDueDate,
		Priority:        req.Priority,
		Notes:           req.Notes,
	}

	if err := database.DB.Create(&maintenance).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create maintenance record")
		return
	}

	utils.Created(c, maintenance)
}

func (h *MaintenanceHandler) GetMaintenances(c *gin.Context) {
	var maintenances []model.MaintenanceRecord
	query := database.DB.Preload("Ship")

	if shipID := c.Query("ship_id"); shipID != "" {
		query = query.Where("ship_id = ?", shipID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if maintenanceType := c.Query("maintenance_type"); maintenanceType != "" {
		query = query.Where("maintenance_type = ?", maintenanceType)
	}

	showDueSoon := c.Query("due_soon")
	if showDueSoon == "true" {
		nextWeek := time.Now().AddDate(0, 0, 7)
		query = query.Where("next_due_date BETWEEN ? AND ?", time.Now(), nextWeek)
	}

	var total int64
	query.Model(&model.MaintenanceRecord{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "10"), 10)
	offset := (page - 1) * pageSize

	if err := query.Order("planned_date DESC").Offset(offset).Limit(pageSize).Find(&maintenances).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch maintenance records")
		return
	}

	utils.Paginated(c, maintenances, total, page, pageSize)
}

func (h *MaintenanceHandler) GetMaintenance(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid maintenance ID")
		return
	}

	var maintenance model.MaintenanceRecord
	if err := database.DB.Preload("Ship").First(&maintenance, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Maintenance record not found")
		return
	}

	utils.Success(c, maintenance)
}

func (h *MaintenanceHandler) UpdateMaintenance(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid maintenance ID")
		return
	}

	var maintenance model.MaintenanceRecord
	if err := database.DB.First(&maintenance, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Maintenance record not found")
		return
	}

	var req model.UpdateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Status != "" {
		maintenance.Status = req.Status
		if req.Status == model.MaintenanceStatusInProgress {
			now := time.Now()
			maintenance.StartDate = &now
		} else if req.Status == model.MaintenanceStatusCompleted {
			now := time.Now()
			maintenance.CompletedDate = &now
		}
	}
	if req.StartDate != nil {
		maintenance.StartDate = req.StartDate
	}
	if req.CompletedDate != nil {
		maintenance.CompletedDate = req.CompletedDate
	}
	if req.Cost != nil {
		maintenance.Cost = *req.Cost
	}
	if req.Provider != "" {
		maintenance.Provider = req.Provider
	}
	if req.Technician != "" {
		maintenance.Technician = req.Technician
	}
	if req.Notes != "" {
		maintenance.Notes = req.Notes
	}

	if err := database.DB.Save(&maintenance).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update maintenance record")
		return
	}

	utils.Success(c, maintenance)
}

func (h *MaintenanceHandler) DeleteMaintenance(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid maintenance ID")
		return
	}

	if err := database.DB.Delete(&model.MaintenanceRecord{}, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete maintenance record")
		return
	}

	utils.Success(c, gin.H{"message": "Maintenance record deleted successfully"})
}

func (h *MaintenanceHandler) CreateSchedule(c *gin.Context) {
	var req model.CreateMaintenanceScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	shipID, _ := uuid.Parse(req.ShipID)

	schedule := model.MaintenanceSchedule{
		ShipID:       shipID,
		Title:        req.Title,
		Description:  req.Description,
		IntervalDays: req.IntervalDays,
		NextDue:      req.NextDue,
		IsActive:     req.IsActive,
	}

	if err := database.DB.Create(&schedule).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create maintenance schedule")
		return
	}

	utils.Created(c, schedule)
}

func (h *MaintenanceHandler) GetSchedules(c *gin.Context) {
	var schedules []model.MaintenanceSchedule
	query := database.DB.Preload("Ship")

	if shipID := c.Query("ship_id"); shipID != "" {
		query = query.Where("ship_id = ?", shipID)
	}

	showDueSoon := c.Query("due_soon")
	if showDueSoon == "true" {
		nextWeek := time.Now().AddDate(0, 0, 7)
		query = query.Where("next_due BETWEEN ? AND ? AND is_active = ?", time.Now(), nextWeek, true)
	}

	if err := query.Order("next_due ASC").Find(&schedules).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch maintenance schedules")
		return
	}

	utils.Success(c, schedules)
}

func (h *MaintenanceHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var schedule model.MaintenanceSchedule
	if err := database.DB.First(&schedule, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Maintenance schedule not found")
		return
	}

	var req model.CreateMaintenanceScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Title != "" {
		schedule.Title = req.Title
	}
	if req.Description != "" {
		schedule.Description = req.Description
	}
	if req.IntervalDays != 0 {
		schedule.IntervalDays = req.IntervalDays
	}
	if !req.NextDue.IsZero() {
		schedule.NextDue = req.NextDue
	}
	schedule.IsActive = req.IsActive

	if err := database.DB.Save(&schedule).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update maintenance schedule")
		return
	}

	utils.Success(c, schedule)
}

func (h *MaintenanceHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	if err := database.DB.Delete(&model.MaintenanceSchedule{}, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete maintenance schedule")
		return
	}

	utils.Success(c, gin.H{"message": "Maintenance schedule deleted successfully"})
}
