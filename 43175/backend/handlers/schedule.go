package handlers

import (
	"net/http"
	"smart-energy-platform/models"
	"smart-energy-platform/services"
	"smart-energy-platform/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateScheduleRequest struct {
	FamilyID    uint   `json:"familyId" binding:"required"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	DeviceID    uint   `json:"deviceId" binding:"required"`
	Action      string `json:"action" binding:"required,oneof=on off toggle dim set_temp"`
	Value       string `json:"value"`
	CronExpr    string `json:"cronExpr" binding:"required"`
	IsEnabled   bool   `json:"isEnabled"`
}

type UpdateScheduleRequest struct {
	Name        string `json:"name" binding:"min=2,max=100"`
	Description string `json:"description"`
	Action      string `json:"action" binding:"omitempty,oneof=on off toggle dim set_temp"`
	Value       string `json:"value"`
	CronExpr    string `json:"cronExpr"`
	IsEnabled   *bool  `json:"isEnabled"`
}

func ListSchedules(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")
	deviceID := parseQueryUint(c, "deviceId")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var familyIDs []uint
	if familyID > 0 {
		familyIDs = []uint{familyID}
	} else {
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
	}

	var schedules []models.Schedule
	query := models.DB

	if len(familyIDs) > 0 {
		query = query.Where("family_id IN ?", familyIDs)
	}
	if deviceID > 0 {
		query = query.Where("device_id = ?", deviceID)
	}

	query.Order("created_at DESC").Find(&schedules)
	utils.Success(c, schedules)
}

func CreateSchedule(c *gin.Context) {
	userID := c.GetUint("userId")

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !hasFamilyAccess(userID, req.FamilyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	if err := services.ParseCronField(req.CronExpr); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "Invalid cron expression: "+err.Error())
		return
	}

	if services.CheckTimeConflict(req.CronExpr, nil) {
		utils.Error(c, http.StatusConflict, 409, "Time conflict detected with existing schedule")
		return
	}

	var device models.Device
	if err := models.DB.First(&device, req.DeviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	if device.FamilyID != req.FamilyID {
		utils.Error(c, http.StatusBadRequest, 400, "Device must belong to the specified family")
		return
	}

	schedule := models.Schedule{
		FamilyID:    req.FamilyID,
		Name:        req.Name,
		Description: req.Description,
		DeviceID:    req.DeviceID,
		Action:      req.Action,
		Value:       req.Value,
		CronExpr:    req.CronExpr,
		IsEnabled:   req.IsEnabled,
	}

	if err := models.DB.Create(&schedule).Error; err != nil {
		utils.InternalError(c, "Failed to create schedule")
		return
	}

	if req.IsEnabled {
		services.ReloadSchedule(schedule.ID)
	}

	utils.Success(c, schedule)
}

func GetSchedule(c *gin.Context) {
	userID := c.GetUint("userId")
	scheduleID := parseUintParam(c, "id")

	var schedule models.Schedule
	if err := models.DB.First(&schedule, scheduleID).Error; err != nil {
		utils.NotFound(c, "Schedule not found")
		return
	}

	if !hasFamilyAccess(userID, schedule.FamilyID) {
		utils.Forbidden(c, "No access to this schedule")
		return
	}

	utils.Success(c, schedule)
}

func UpdateSchedule(c *gin.Context) {
	userID := c.GetUint("userId")
	scheduleID := parseUintParam(c, "id")

	var schedule models.Schedule
	if err := models.DB.First(&schedule, scheduleID).Error; err != nil {
		utils.NotFound(c, "Schedule not found")
		return
	}

	if !hasFamilyAdminAccess(userID, schedule.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Action != "" {
		updates["action"] = req.Action
	}
	if req.Value != "" {
		updates["value"] = req.Value
	}
	if req.CronExpr != "" {
		if err := services.ParseCronField(req.CronExpr); err != nil {
			utils.Error(c, http.StatusBadRequest, 400, "Invalid cron expression")
			return
		}
		if services.CheckTimeConflict(req.CronExpr, &scheduleID) {
			utils.Error(c, http.StatusConflict, 409, "Time conflict detected")
			return
		}
		updates["cron_expr"] = req.CronExpr
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}

	if err := models.DB.Model(&schedule).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update schedule")
		return
	}

	services.ReloadSchedule(scheduleID)

	models.DB.First(&schedule, scheduleID)
	utils.Success(c, schedule)
}

func DeleteSchedule(c *gin.Context) {
	userID := c.GetUint("userId")
	scheduleID := parseUintParam(c, "id")

	var schedule models.Schedule
	if err := models.DB.First(&schedule, scheduleID).Error; err != nil {
		utils.NotFound(c, "Schedule not found")
		return
	}

	if !hasFamilyAdminAccess(userID, schedule.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	services.RemoveSchedule(scheduleID)

	models.DB.Where("schedule_id = ?", scheduleID).Delete(&models.ScheduleLog{})
	models.DB.Delete(&schedule)

	utils.Success(c, nil)
}

func ListScheduleLogs(c *gin.Context) {
	userID := c.GetUint("userId")
	scheduleID := parseUintParam(c, "id")

	var schedule models.Schedule
	if err := models.DB.First(&schedule, scheduleID).Error; err != nil {
		utils.NotFound(c, "Schedule not found")
		return
	}

	if !hasFamilyAccess(userID, schedule.FamilyID) {
		utils.Forbidden(c, "No access to this schedule")
		return
	}

	var logs []models.ScheduleLog
	models.DB.Where("schedule_id = ?", scheduleID).
		Order("executed_at DESC").
		Limit(50).
		Find(&logs)

	utils.Success(c, logs)
}

func validateTimeRange(startTime, endTime time.Time) error {
	if endTime.Before(startTime) {
		return http.ErrNotSupported
	}
	return nil
}
