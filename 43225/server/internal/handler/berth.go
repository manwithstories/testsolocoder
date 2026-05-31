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

type BerthHandler struct{}

func NewBerthHandler() *BerthHandler {
	return &BerthHandler{}
}

func (h *BerthHandler) CreateDock(c *gin.Context) {
	var req model.CreateDockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dock := model.Dock{
		Name:        req.Name,
		Address:     req.Address,
		City:        req.City,
		Country:     req.Country,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Amenities:   req.Amenities,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		IsActive:    true,
	}

	if err := database.DB.Create(&dock).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create dock")
		return
	}

	utils.Created(c, dock)
}

func (h *BerthHandler) GetDocks(c *gin.Context) {
	var docks []model.Dock
	query := database.DB.Preload("Berths").Where("is_active = ?", true)

	if city := c.Query("city"); city != "" {
		query = query.Where("city ILIKE ?", "%"+city+"%")
	}
	if country := c.Query("country"); country != "" {
		query = query.Where("country ILIKE ?", "%"+country+"%")
	}

	var total int64
	query.Model(&model.Dock{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "10"), 10)
	offset := (page - 1) * pageSize

	if err := query.Order("average_rating DESC").Offset(offset).Limit(pageSize).Find(&docks).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch docks")
		return
	}

	utils.Paginated(c, docks, total, page, pageSize)
}

func (h *BerthHandler) GetDock(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid dock ID")
		return
	}

	var dock model.Dock
	if err := database.DB.Preload("Berths").First(&dock, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Dock not found")
		return
	}

	utils.Success(c, dock)
}

func (h *BerthHandler) UpdateDock(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid dock ID")
		return
	}

	var dock model.Dock
	if err := database.DB.First(&dock, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Dock not found")
		return
	}

	var req model.CreateDockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name != "" {
		dock.Name = req.Name
	}
	if req.Address != "" {
		dock.Address = req.Address
	}
	if req.City != "" {
		dock.City = req.City
	}
	if req.Country != "" {
		dock.Country = req.Country
	}
	if req.Description != "" {
		dock.Description = req.Description
	}
	if req.ImageURL != "" {
		dock.ImageURL = req.ImageURL
	}
	if req.Amenities != "" {
		dock.Amenities = req.Amenities
	}
	if req.OpenTime != "" {
		dock.OpenTime = req.OpenTime
	}
	if req.CloseTime != "" {
		dock.CloseTime = req.CloseTime
	}

	if err := database.DB.Save(&dock).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update dock")
		return
	}

	utils.Success(c, dock)
}

func (h *BerthHandler) DeleteDock(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid dock ID")
		return
	}

	if err := database.DB.Delete(&model.Dock{}, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete dock")
		return
	}

	utils.Success(c, gin.H{"message": "Dock deleted successfully"})
}

func (h *BerthHandler) CreateBerth(c *gin.Context) {
	var req model.CreateBerthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dockID, _ := uuid.Parse(req.DockID)

	berth := model.Berth{
		DockID:      dockID,
		Number:      req.Number,
		BerthType:   req.BerthType,
		MaxLength:   req.MaxLength,
		MaxWidth:    req.MaxWidth,
		HourlyRate:  req.HourlyRate,
		DailyRate:   req.DailyRate,
		HasWater:    req.HasWater,
		HasElectric: req.HasElectric,
		HasInternet: req.HasInternet,
		Description: req.Description,
		Status:      model.BerthStatusAvailable,
	}

	if err := database.DB.Create(&berth).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create berth")
		return
	}

	utils.Created(c, berth)
}

func (h *BerthHandler) GetBerths(c *gin.Context) {
	var berths []model.Berth
	query := database.DB.Preload("Dock")

	if dockID := c.Query("dock_id"); dockID != "" {
		query = query.Where("dock_id = ?", dockID)
	}
	if berthType := c.Query("berth_type"); berthType != "" {
		query = query.Where("berth_type = ?", berthType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&model.Berth{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "10"), 10)
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("number ASC").Find(&berths).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch berths")
		return
	}

	utils.Paginated(c, berths, total, page, pageSize)
}

func (h *BerthHandler) GetBerth(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid berth ID")
		return
	}

	var berth model.Berth
	if err := database.DB.Preload("Reservations").First(&berth, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Berth not found")
		return
	}

	utils.Success(c, berth)
}

func (h *BerthHandler) UpdateBerth(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid berth ID")
		return
	}

	var berth model.Berth
	if err := database.DB.First(&berth, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Berth not found")
		return
	}

	var req model.CreateBerthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Number != "" {
		berth.Number = req.Number
	}
	if req.BerthType != "" {
		berth.BerthType = req.BerthType
	}
	if req.MaxLength != 0 {
		berth.MaxLength = req.MaxLength
	}
	if req.MaxWidth != 0 {
		berth.MaxWidth = req.MaxWidth
	}
	if req.HourlyRate != 0 {
		berth.HourlyRate = req.HourlyRate
	}
	if req.DailyRate != 0 {
		berth.DailyRate = req.DailyRate
	}
	berth.HasWater = req.HasWater
	berth.HasElectric = req.HasElectric
	berth.HasInternet = req.HasInternet
	if req.Description != "" {
		berth.Description = req.Description
	}

	if err := database.DB.Save(&berth).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update berth")
		return
	}

	utils.Success(c, berth)
}

func (h *BerthHandler) DeleteBerth(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid berth ID")
		return
	}

	if err := database.DB.Delete(&model.Berth{}, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete berth")
		return
	}

	utils.Success(c, gin.H{"message": "Berth deleted successfully"})
}

func (h *BerthHandler) CheckAvailability(c *gin.Context) {
	var req model.CheckAvailabilityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var conflictingReservations []model.BerthReservation
	err := database.DB.Where(
		"berth_id = ? AND status = ? AND (start_time < ? AND end_time > ?)",
		req.BerthID, "confirmed", req.EndTime, req.StartTime,
	).Find(&conflictingReservations).Error

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to check availability")
		return
	}

	available := len(conflictingReservations) == 0

	utils.Success(c, gin.H{
		"available": available,
		"berth_id":  req.BerthID,
		"conflicts": conflictingReservations,
	})
}

func (h *BerthHandler) CreateReservation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := uuid.Parse(userID.(string))

	var req model.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	berthID, _ := uuid.Parse(req.BerthID)

	var conflictingReservations []model.BerthReservation
	database.DB.Where(
		"berth_id = ? AND status = ? AND (start_time < ? AND end_time > ?)",
		berthID, "confirmed", req.EndTime, req.StartTime,
	).Find(&conflictingReservations)

	if len(conflictingReservations) > 0 {
		utils.Error(c, http.StatusConflict, "Berth is already reserved for this time period")
		return
	}

	var berth model.Berth
	if err := database.DB.First(&berth, "id = ?", berthID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Berth not found")
		return
	}

	duration := req.EndTime.Sub(req.StartTime)
	hours := duration.Hours()
	var totalAmount float64
	if hours <= 24 {
		totalAmount = berth.HourlyRate * hours
	} else {
		days := int(hours / 24)
		remainingHours := hours - float64(days*24)
		totalAmount = (berth.DailyRate * float64(days)) + (berth.HourlyRate * remainingHours)
	}

	var shipID *uuid.UUID
	if req.ShipID != "" {
		sid, err := uuid.Parse(req.ShipID)
		if err == nil {
			shipID = &sid
		}
	}

	reservation := model.BerthReservation{
		BerthID:     berthID,
		ShipID:      shipID,
		UserID:      uid,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		TotalAmount: utils.RoundTo(totalAmount, 2),
		Status:      "confirmed",
		Notes:       req.Notes,
	}

	if err := database.DB.Create(&reservation).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create reservation")
		return
	}

	utils.Created(c, reservation)
}

func (h *BerthHandler) GetReservations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := uuid.Parse(userID.(string))

	var reservations []model.BerthReservation
	query := database.DB.Preload("Berth").Preload("Ship").Where("user_id = ?", uid)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "10"), 10)
	offset := (page - 1) * pageSize

	var total int64
	query.Model(&model.BerthReservation{}).Count(&total)

	if err := query.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&reservations).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch reservations")
		return
	}

	utils.Paginated(c, reservations, total, page, pageSize)
}

func (h *BerthHandler) CancelReservation(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid reservation ID")
		return
	}

	userID, _ := c.Get("user_id")
	uid, _ := uuid.Parse(userID.(string))

	var reservation model.BerthReservation
	if err := database.DB.First(&reservation, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Reservation not found")
		return
	}

	if reservation.UserID != uid {
		role, _ := c.Get("role")
		if role.(string) != string(model.RoleAdmin) {
			utils.Error(c, http.StatusForbidden, "You can only cancel your own reservations")
			return
		}
	}

	if reservation.StartTime.Before(time.Now()) {
		utils.Error(c, http.StatusBadRequest, "Cannot cancel past reservations")
		return
	}

	reservation.Status = "cancelled"
	if err := database.DB.Save(&reservation).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to cancel reservation")
		return
	}

	utils.Success(c, gin.H{"message": "Reservation cancelled successfully"})
}

func (h *BerthHandler) RecordWaterLevel(c *gin.Context) {
	var req model.RecordWaterLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	dockID, _ := uuid.Parse(req.DockID)

	unit := req.Unit
	if unit == "" {
		unit = "meters"
	}

	waterLevel := model.WaterLevel{
		DockID:     dockID,
		Height:     req.Height,
		Unit:       unit,
		RecordedAt: req.RecordedAt,
	}

	if err := database.DB.Create(&waterLevel).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to record water level")
		return
	}

	utils.Created(c, waterLevel)
}

func (h *BerthHandler) GetWaterLevels(c *gin.Context) {
	dockID := c.Param("id")
	if !utils.IsValidUUID(dockID) {
		utils.Error(c, http.StatusBadRequest, "Invalid dock ID")
		return
	}

	var waterLevels []model.WaterLevel
	query := database.DB.Where("dock_id = ?", dockID)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" && endDate != "" {
		query = query.Where("recorded_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Order("recorded_at DESC").Limit(100).Find(&waterLevels).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch water levels")
		return
	}

	utils.Success(c, waterLevels)
}
