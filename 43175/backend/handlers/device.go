package handlers

import (
	"smart-energy-platform/models"
	"smart-energy-platform/services"
	"smart-energy-platform/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateDeviceRequest struct {
	FamilyID   uint    `json:"familyId" binding:"required"`
	Name       string  `json:"name" binding:"required,min=2,max=100"`
	DeviceType string  `json:"deviceType" binding:"required,oneof=light ac heater curtain camera sensor switch plug other"`
	Vendor     string  `json:"vendor" binding:"max=100"`
	Location   string  `json:"location" binding:"max=100"`
	Power      float64 `json:"power" binding:"min=0,max=10000"`
	Protocol   string  `json:"protocol" binding:"required,oneof=wifi zigbee bluetooth zwave mqtt other"`
}

type UpdateDeviceRequest struct {
	Name       string  `json:"name" binding:"min=2,max=100"`
	DeviceType string  `json:"deviceType" binding:"omitempty,oneof=light ac heater curtain camera sensor switch plug other"`
	Vendor     string  `json:"vendor" binding:"max=100"`
	Location   string  `json:"location" binding:"max=100"`
	Power      float64 `json:"power" binding:"omitempty,min=0,max=10000"`
	Protocol   string  `json:"protocol" binding:"omitempty,oneof=wifi zigbee bluetooth zwave mqtt other"`
	Status     string  `json:"status" binding:"omitempty,oneof=online offline"`
}

type UpdateDeviceStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=on off online offline"`
}

func ListDevices(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")
	deviceType := c.Query("deviceType")
	status := c.Query("status")
	location := c.Query("location")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var devices []models.Device
	query := models.DB

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
			utils.Success(c, []models.Device{})
			return
		}
	}

	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}

	query.Find(&devices)

	utils.Success(c, devices)
}

func CreateDevice(c *gin.Context) {
	userID := c.GetUint("userId")

	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !hasFamilyAccess(userID, req.FamilyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	device := models.Device{
		FamilyID:   req.FamilyID,
		Name:       req.Name,
		DeviceType: req.DeviceType,
		Vendor:     req.Vendor,
		Location:   req.Location,
		Power:      req.Power,
		Protocol:   req.Protocol,
		Status:     "offline",
	}

	if err := models.DB.Create(&device).Error; err != nil {
		utils.InternalError(c, "Failed to create device")
		return
	}

	services.UpdateDeviceStatus(device.ID, device.Status)
	services.InvalidateFamilyDevicesCache(device.FamilyID)

	utils.Success(c, device)
}

func GetDevice(c *gin.Context) {
	userID := c.GetUint("userId")
	deviceID := parseUintParam(c, "id")

	var device models.Device
	if err := models.DB.First(&device, deviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	if !hasFamilyAccess(userID, device.FamilyID) {
		utils.Forbidden(c, "No access to this device")
		return
	}

	utils.Success(c, device)
}

func UpdateDevice(c *gin.Context) {
	userID := c.GetUint("userId")
	deviceID := parseUintParam(c, "id")

	var device models.Device
	if err := models.DB.First(&device, deviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	if !hasFamilyAdminAccess(userID, device.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.DeviceType != "" {
		updates["device_type"] = req.DeviceType
	}
	if req.Vendor != "" {
		updates["vendor"] = req.Vendor
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Power > 0 {
		updates["power"] = req.Power
	}
	if req.Protocol != "" {
		updates["protocol"] = req.Protocol
	}
	if req.Status != "" {
		updates["status"] = req.Status
		if req.Status == "online" {
			now := time.Now()
			updates["last_online_time"] = &now
		}
	}

	if err := models.DB.Model(&device).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update device")
		return
	}

	if req.Status != "" {
		services.UpdateDeviceStatus(device.ID, req.Status)
	}
	services.InvalidateFamilyDevicesCache(device.FamilyID)

	models.DB.First(&device, deviceID)
	utils.Success(c, device)
}

func DeleteDevice(c *gin.Context) {
	userID := c.GetUint("userId")
	deviceID := parseUintParam(c, "id")

	var device models.Device
	if err := models.DB.First(&device, deviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	if !hasFamilyAdminAccess(userID, device.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	models.DB.Where("device_id = ?", deviceID).Delete(&models.EnergyData{})
	models.DB.Where("device_id = ?", deviceID).Delete(&models.ScheduleLog{})
	models.DB.Where("device_id = ?", deviceID).Delete(&models.Schedule{})
	models.DB.Where("device_id = ?", deviceID).Delete(&models.SceneAction{})
	models.DB.Where("device_id = ?", deviceID).Delete(&models.EnergyAlert{})
	models.DB.Delete(&device)

	services.InvalidateFamilyDevicesCache(device.FamilyID)

	utils.Success(c, nil)
}

func UpdateDeviceStatus(c *gin.Context) {
	userID := c.GetUint("userId")
	deviceID := parseUintParam(c, "id")

	var device models.Device
	if err := models.DB.First(&device, deviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	if !hasFamilyAccess(userID, device.FamilyID) {
		utils.Forbidden(c, "No access to this device")
		return
	}

	var req UpdateDeviceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	device.Status = req.Status
	if req.Status == "online" || req.Status == "on" {
		now := time.Now()
		device.LastOnlineTime = &now
	}
	models.DB.Save(&device)

	services.UpdateDeviceStatus(device.ID, req.Status)

	utils.Success(c, gin.H{
		"id":     device.ID,
		"status": device.Status,
	})
}

func GetDeviceEnergy(c *gin.Context) {
	userID := c.GetUint("userId")
	deviceID := parseUintParam(c, "id")
	period := c.DefaultQuery("period", "day")

	var device models.Device
	if err := models.DB.First(&device, deviceID).Error; err != nil {
		utils.NotFound(c, "Device not found")
		return
	}

	if !hasFamilyAccess(userID, device.FamilyID) {
		utils.Forbidden(c, "No access to this device")
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

	var energyData []models.EnergyData
	models.DB.Where("device_id = ? AND timestamp >= ?", deviceID, startTime).
		Order("timestamp ASC").
		Find(&energyData)

	var totalEnergy float64
	for _, d := range energyData {
		totalEnergy += d.EnergyUsed
	}

	utils.Success(c, gin.H{
		"data":        energyData,
		"totalEnergy": totalEnergy,
		"period":      period,
		"deviceId":    deviceID,
	})
}

func parseQueryUint(c *gin.Context, name string) uint {
	val := c.Query(name)
	if val == "" {
		return 0
	}
	n, _ := strconv.ParseUint(val, 10, 64)
	return uint(n)
}
