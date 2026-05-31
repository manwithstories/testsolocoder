package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"ship-rental-platform/internal/config"
	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShipHandler struct{}

func NewShipHandler() *ShipHandler {
	return &ShipHandler{}
}

func (h *ShipHandler) CreateShip(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req model.CreateShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	ownerID, _ := uuid.Parse(userID.(string))

	ship := model.Ship{
		OwnerID:           ownerID,
		Name:              req.Name,
		Description:       req.Description,
		ShipType:          req.ShipType,
		Capacity:          req.Capacity,
		CabinCount:        req.CabinCount,
		BathroomCount:     req.BathroomCount,
		Length:            req.Length,
		Width:             req.Width,
		YearBuilt:         req.YearBuilt,
		Equipment:         req.Equipment,
		Features:          req.Features,
		SailingArea:       req.SailingArea,
		HomePort:          req.HomePort,
		LicenseNumber:     req.LicenseNumber,
		HourlyRate:        req.HourlyRate,
		DailyRate:         req.DailyRate,
		DepositAmount:     req.DepositAmount,
		InsuranceRequired: req.InsuranceRequired,
		CancellationPolicy: req.CancellationPolicy,
		Status:            model.ShipStatusAvailable,
	}

	if err := database.DB.Create(&ship).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create ship")
		return
	}

	utils.Created(c, ship)
}

func (h *ShipHandler) GetShips(c *gin.Context) {
	var req model.SearchShipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query := database.DB.Preload("Images").Where("status = ?", model.ShipStatusAvailable)

	if req.ShipType != nil {
		query = query.Where("ship_type = ?", *req.ShipType)
	}
	if req.MinCapacity != nil {
		query = query.Where("capacity >= ?", *req.MinCapacity)
	}
	if req.MaxCapacity != nil {
		query = query.Where("capacity <= ?", *req.MaxCapacity)
	}
	if req.MinPrice != nil {
		query = query.Where("daily_rate >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("daily_rate <= ?", *req.MaxPrice)
	}
	if req.Location != "" {
		query = query.Where("home_port ILIKE ? OR sailing_area ILIKE ?", "%"+req.Location+"%", "%"+req.Location+"%")
	}

	var total int64
	query.Model(&model.Ship{}).Count(&total)

	sortBy := req.SortBy
	sortOrder := req.SortOrder
	validSorts := map[string]bool{
		"rating":     true,
		"price_asc":  true,
		"price_desc": true,
		"created_at": true,
		"name":       true,
	}
	if !validSorts[sortBy] {
		sortBy = "rating"
	}

	orderClause := "average_rating DESC"
	switch sortBy {
	case "price_asc":
		orderClause = "daily_rate ASC"
	case "price_desc":
		orderClause = "daily_rate DESC"
	case "created_at":
		orderClause = "created_at DESC"
	case "name":
		orderClause = "name ASC"
	}

	offset := (req.Page - 1) * req.PageSize
	var ships []model.Ship
	if err := query.Order(orderClause).Offset(offset).Limit(req.PageSize).Find(&ships).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch ships")
		return
	}

	utils.Paginated(c, ships, total, req.Page, req.PageSize)
}

func (h *ShipHandler) GetShip(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid ship ID")
		return
	}

	var ship model.Ship
	if err := database.DB.Preload("Images").Preload("Owner").First(&ship, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Ship not found")
		return
	}

	ship.ViewCount++
	database.DB.Model(&ship).Update("view_count", ship.ViewCount)

	utils.Success(c, ship)
}

func (h *ShipHandler) UpdateShip(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid ship ID")
		return
	}

	userID, _ := c.Get("user_id")
	ownerID, _ := uuid.Parse(userID.(string))

	var ship model.Ship
	if err := database.DB.First(&ship, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Ship not found")
		return
	}

	if ship.OwnerID != ownerID {
		role, _ := c.Get("role")
		if role.(string) != string(model.RoleAdmin) {
			utils.Error(c, http.StatusForbidden, "You can only update your own ships")
			return
		}
	}

	var req model.UpdateShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name != "" {
		ship.Name = req.Name
	}
	if req.Description != "" {
		ship.Description = req.Description
	}
	if req.ShipType != nil {
		ship.ShipType = *req.ShipType
	}
	if req.Capacity != nil {
		ship.Capacity = *req.Capacity
	}
	if req.CabinCount != nil {
		ship.CabinCount = *req.CabinCount
	}
	if req.BathroomCount != nil {
		ship.BathroomCount = *req.BathroomCount
	}
	if req.Length != nil {
		ship.Length = *req.Length
	}
	if req.Width != nil {
		ship.Width = *req.Width
	}
	if req.YearBuilt != nil {
		ship.YearBuilt = *req.YearBuilt
	}
	if req.Equipment != "" {
		ship.Equipment = req.Equipment
	}
	if req.Features != "" {
		ship.Features = req.Features
	}
	if req.SailingArea != "" {
		ship.SailingArea = req.SailingArea
	}
	if req.HomePort != "" {
		ship.HomePort = req.HomePort
	}
	if req.LicenseNumber != "" {
		ship.LicenseNumber = req.LicenseNumber
	}
	if req.HourlyRate != nil {
		ship.HourlyRate = *req.HourlyRate
	}
	if req.DailyRate != nil {
		ship.DailyRate = *req.DailyRate
	}
	if req.DepositAmount != nil {
		ship.DepositAmount = *req.DepositAmount
	}
	if req.InsuranceRequired != nil {
		ship.InsuranceRequired = *req.InsuranceRequired
	}
	if req.CancellationPolicy != "" {
		ship.CancellationPolicy = req.CancellationPolicy
	}
	if req.Status != nil {
		ship.Status = *req.Status
	}

	if err := database.DB.Save(&ship).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update ship")
		return
	}

	utils.Success(c, ship)
}

func (h *ShipHandler) DeleteShip(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid ship ID")
		return
	}

	userID, _ := c.Get("user_id")
	ownerID, _ := uuid.Parse(userID.(string))

	var ship model.Ship
	if err := database.DB.First(&ship, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Ship not found")
		return
	}

	if ship.OwnerID != ownerID {
		role, _ := c.Get("role")
		if role.(string) != string(model.RoleAdmin) {
			utils.Error(c, http.StatusForbidden, "You can only delete your own ships")
			return
		}
	}

	if err := database.DB.Delete(&ship).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete ship")
		return
	}

	utils.Success(c, gin.H{"message": "Ship deleted successfully"})
}

func (h *ShipHandler) GetMyShips(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ownerID, _ := uuid.Parse(userID.(string))

	var ships []model.Ship
	if err := database.DB.Preload("Images").Where("owner_id = ?", ownerID).Order("created_at DESC").Find(&ships).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch ships")
		return
	}

	utils.Success(c, ships)
}

func (h *ShipHandler) UploadImage(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid ship ID")
		return
	}

	userID, _ := c.Get("user_id")
	ownerID, _ := uuid.Parse(userID.(string))

	var ship model.Ship
	if err := database.DB.First(&ship, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Ship not found")
		return
	}

	if ship.OwnerID != ownerID {
		utils.Error(c, http.StatusForbidden, "You can only upload images to your own ships")
		return
	}

	cfg := config.GetConfig()
	file, err := c.FormFile("image")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Failed to get image file")
		return
	}

	if !utils.IsAllowedFileType(file.Filename, cfg.Upload.AllowedTypes) {
		utils.Error(c, http.StatusBadRequest, "Invalid file type")
		return
	}

	maxSize := int64(cfg.Upload.MaxSize * 1024 * 1024)
	if file.Size > maxSize {
		utils.Error(c, http.StatusBadRequest, "File too large")
		return
	}

	shipID, _ := uuid.Parse(id)
	ext := filepath.Ext(file.Filename)
	imgID := uuid.New().String()
	newFilename := imgID + ext
	uploadPath := filepath.Join(cfg.Upload.Path, "ships", shipID.String())

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	filePath := filepath.Join(uploadPath, newFilename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	isPrimary, _ := strconv.ParseBool(c.DefaultPostForm("is_primary", "false"))
	sortOrder, _ := strconv.Atoi(c.DefaultPostForm("sort_order", "0"))

	image := model.ShipImage{
		ShipID:    shipID,
		URL:       "/uploads/ships/" + shipID.String() + "/" + newFilename,
		IsPrimary: isPrimary,
		SortOrder: sortOrder,
	}

	if err := database.DB.Create(&image).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save image record")
		return
	}

	utils.Created(c, image)
}

func (h *ShipHandler) DeleteImage(c *gin.Context) {
	shipID := c.Param("id")
	imageID := c.Param("imageId")

	if !utils.IsValidUUID(shipID) || !utils.IsValidUUID(imageID) {
		utils.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	userID, _ := c.Get("user_id")
	ownerID, _ := uuid.Parse(userID.(string))

	var ship model.Ship
	if err := database.DB.First(&ship, "id = ?", shipID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Ship not found")
		return
	}

	if ship.OwnerID != ownerID {
		utils.Error(c, http.StatusForbidden, "You can only delete images from your own ships")
		return
	}

	if err := database.DB.Delete(&model.ShipImage{}, "id = ? AND ship_id = ?", imageID, shipID).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	utils.Success(c, gin.H{"message": "Image deleted successfully"})
}
