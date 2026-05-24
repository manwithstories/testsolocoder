package handlers

import (
	"smart-energy-platform/models"
	"smart-energy-platform/services"
	"smart-energy-platform/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateGroupRequest struct {
	FamilyID    uint   `json:"familyId" binding:"required"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required,oneof=room floor function custom"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name" binding:"min=2,max=100"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"omitempty,oneof=room floor function custom"`
}

type AddDeviceToGroupRequest struct {
	DeviceID uint `json:"deviceId" binding:"required"`
}

type BatchControlRequest struct {
	Action string `json:"action" binding:"required,oneof=on off"`
	Value  string `json:"value"`
}

func ListGroups(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var groups []models.DeviceGroup
	query := models.DB.Preload("Devices")

	if familyID > 0 {
		query = query.Where("family_id = ?", familyID)
	} else {
		var familyIDs []uint
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
		if len(familyIDs) > 0 {
			query = query.Where("family_id IN ?", familyIDs)
		} else {
			utils.Success(c, []models.DeviceGroup{})
			return
		}
	}

	query.Find(&groups)
	utils.Success(c, groups)
}

func CreateGroup(c *gin.Context) {
	userID := c.GetUint("userId")

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !hasFamilyAccess(userID, req.FamilyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	group := models.DeviceGroup{
		FamilyID:    req.FamilyID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	if err := models.DB.Create(&group).Error; err != nil {
		utils.InternalError(c, "Failed to create group")
		return
	}

	utils.Success(c, group)
}

func GetGroup(c *gin.Context) {
	userID := c.GetUint("userId")
	groupID := parseUintParam(c, "id")

	var group models.DeviceGroup
	if err := models.DB.Preload("Devices").First(&group, groupID).Error; err != nil {
		utils.NotFound(c, "Group not found")
		return
	}

	if !hasFamilyAccess(userID, group.FamilyID) {
		utils.Forbidden(c, "No access to this group")
		return
	}

	utils.Success(c, group)
}

func UpdateGroup(c *gin.Context) {
	userID := c.GetUint("userId")
	groupID := parseUintParam(c, "id")

	var group models.DeviceGroup
	if err := models.DB.First(&group, groupID).Error; err != nil {
		utils.NotFound(c, "Group not found")
		return
	}

	if !hasFamilyAdminAccess(userID, group.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	var req UpdateGroupRequest
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
	if req.Type != "" {
		updates["type"] = req.Type
	}

	if err := models.DB.Model(&group).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update group")
		return
	}

	utils.Success(c, group)
}

func DeleteGroup(c *gin.Context) {
	userID := c.GetUint("userId")
	groupID := parseUintParam(c, "id")

	var group models.DeviceGroup
	if err := models.DB.First(&group, groupID).Error; err != nil {
		utils.NotFound(c, "Group not found")
		return
	}

	if !hasFamilyAdminAccess(userID, group.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	models.DB.Model(&group).Association("Devices").Clear()
	models.DB.Delete(&group)

	utils.Success(c, nil)
}

func AddDeviceToGroup(c *gin.Context) {
	userID := c.GetUint("userId")
	groupID := parseUintParam(c, "id")

	var group models.DeviceGroup
	if err := models.DB.Preload("Devices").First(&group, groupID).Error; err != nil {
		utils.NotFound(c, "Group not found")
		return
	}

	if !hasFamilyAccess(userID, group.FamilyID) {
		utils.Forbidden(c, "No access to this group")
		return
	}

	var req AddDeviceToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var device models.Device
	if err := models.DB.First(&device, req.DeviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	if device.FamilyID != group.FamilyID {
		utils.Error(c, 400, 400, "Device must belong to the same family")
		return
	}

	models.DB.Model(&group).Association("Devices").Append(&device)

	utils.Success(c, nil)
}

func RemoveDeviceFromGroup(c *gin.Context) {
	userID := c.GetUint("userId")
	groupID := parseUintParam(c, "id")
	deviceID := parseUintParam(c, "deviceId")

	var group models.DeviceGroup
	if err := models.DB.First(&group, groupID).Error; err != nil {
		utils.NotFound(c, "Group not found")
		return
	}

	if !hasFamilyAccess(userID, group.FamilyID) {
		utils.Forbidden(c, "No access to this group")
		return
	}

	var device models.Device
	if err := models.DB.First(&device, deviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	models.DB.Model(&group).Association("Devices").Delete(&device)

	utils.Success(c, nil)
}

func BatchControlGroup(c *gin.Context) {
	userID := c.GetUint("userId")
	groupID := parseUintParam(c, "id")

	var group models.DeviceGroup
	if err := models.DB.Preload("Devices").First(&group, groupID).Error; err != nil {
		utils.NotFound(c, "Group not found")
		return
	}

	if !hasFamilyAccess(userID, group.FamilyID) {
		utils.Forbidden(c, "No access to this group")
		return
	}

	var req BatchControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	now := time.Now()
	for i := range group.Devices {
		group.Devices[i].Status = req.Action
		if req.Action == "on" || req.Action == "online" {
			group.Devices[i].LastOnlineTime = &now
		}
		models.DB.Save(&group.Devices[i])
		services.UpdateDeviceStatus(group.Devices[i].ID, req.Action)
	}

	utils.Success(c, gin.H{
		"group":   group.Name,
		"action":  req.Action,
		"devices": len(group.Devices),
	})
}

func GetGroupEnergy(c *gin.Context) {
	userID := c.GetUint("userId")
	groupID := parseUintParam(c, "id")
	period := c.DefaultQuery("period", "day")

	var group models.DeviceGroup
	if err := models.DB.Preload("Devices").First(&group, groupID).Error; err != nil {
		utils.NotFound(c, "Group not found")
		return
	}

	if !hasFamilyAccess(userID, group.FamilyID) {
		utils.Forbidden(c, "No access to this group")
		return
	}

	var startTime time.Time
	now := time.Now()
	switch period {
	case "hour":
		startTime = now.Add(-1 * time.Hour)
	case "week":
		startTime = now.Add(-7 * 24 * time.Hour)
	case "month":
		startTime = now.Add(-30 * 24 * time.Hour)
	default:
		startTime = now.Add(-24 * time.Hour)
	}

	var deviceIDs []uint
	for _, d := range group.Devices {
		deviceIDs = append(deviceIDs, d.ID)
	}

	var totalEnergy float64
	if len(deviceIDs) > 0 {
		models.DB.Model(&models.EnergyData{}).
			Where("device_id IN ? AND timestamp >= ?", deviceIDs, startTime).
			Select("COALESCE(SUM(energy_used), 0)").
			Scan(&totalEnergy)
	}

	utils.Success(c, gin.H{
		"group":       group.Name,
		"totalEnergy": totalEnergy,
		"period":      period,
		"deviceCount": len(group.Devices),
	})
}
