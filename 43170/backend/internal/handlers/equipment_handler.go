package handlers

import (
	"net/http"
	"photo-rental/internal/models"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EquipmentHandler struct{}

type CreateEquipmentRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Category    string  `json:"category" binding:"required,max=50"`
	Brand       string  `json:"brand" binding:"required,max=50"`
	Model       string  `json:"model" binding:"required,max=100"`
	PurchaseDate string `json:"purchaseDate"`
	Deposit     float64 `json:"deposit" binding:"required,min=0"`
	DailyRent   float64 `json:"dailyRent" binding:"required,min=0"`
	Description string  `json:"description"`
	Images      []string `json:"images"`
}

type UpdateEquipmentRequest struct {
	Name        string  `json:"name" binding:"omitempty,max=100"`
	Category    string  `json:"category" binding:"omitempty,max=50"`
	Brand       string  `json:"brand" binding:"omitempty,max=50"`
	Model       string  `json:"model" binding:"omitempty,max=100"`
	PurchaseDate string `json:"purchaseDate"`
	Status      string  `json:"status" binding:"omitempty,oneof=available rented maintenance"`
	Deposit     float64 `json:"deposit" binding:"omitempty,min=0"`
	DailyRent   float64 `json:"dailyRent" binding:"omitempty,min=0"`
	Description string  `json:"description"`
}

func NewEquipmentHandler() *EquipmentHandler {
	return &EquipmentHandler{}
}

func (h *EquipmentHandler) CreateEquipment(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	if userRole != "owner" && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only owners can create equipment")
		return
	}

	var req CreateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.Images) > 9 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Maximum 9 images allowed")
		return
	}

	equipment := models.Equipment{
		OwnerID:      userID,
		Name:         req.Name,
		Category:     req.Category,
		Brand:        req.Brand,
		Model:        req.Model,
		PurchaseDate: req.PurchaseDate,
		Deposit:      req.Deposit,
		DailyRent:    req.DailyRent,
		Description:  req.Description,
		Status:       "available",
		Rating:       0,
		ReviewCount:  0,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	if err := tx.Create(&equipment).Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to create equipment: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create equipment")
		return
	}

	for i, imgURL := range req.Images {
		img := models.EquipmentImage{
			EquipmentID: equipment.ID,
			ImageURL:    imgURL,
			SortOrder:   i,
		}
		if err := tx.Create(&img).Error; err != nil {
			tx.Rollback()
			utils.Logger.Error("Failed to create equipment image: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save images")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to commit transaction: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save equipment")
		return
	}

	utils.Logger.Info("Equipment created: %d by user %d", equipment.ID, userID)
	utils.SuccessResponse(c, equipment)
}

func (h *EquipmentHandler) GetEquipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	var equipment models.Equipment
	if err := database.DB.Preload("Images").Preload("Owner").First(&equipment, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Equipment not found")
		return
	}

	if equipment.Owner.Password != "" {
		equipment.Owner.Password = ""
	}

	utils.SuccessResponse(c, equipment)
}

func (h *EquipmentHandler) GetMyEquipments(c *gin.Context) {
	userID := c.GetUint("userId")

	var equipments []models.Equipment
	if err := database.DB.Preload("Images").Where("owner_id = ?", userID).Find(&equipments).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch equipments")
		return
	}

	utils.SuccessResponse(c, equipments)
}

func (h *EquipmentHandler) UpdateEquipment(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Equipment not found")
		return
	}

	if equipment.OwnerID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own equipment")
		return
	}

	var req UpdateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Brand != "" {
		updates["brand"] = req.Brand
	}
	if req.Model != "" {
		updates["model"] = req.Model
	}
	if req.PurchaseDate != "" {
		updates["purchase_date"] = req.PurchaseDate
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Deposit > 0 {
		updates["deposit"] = req.Deposit
	}
	if req.DailyRent > 0 {
		updates["daily_rent"] = req.DailyRent
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "No fields to update")
		return
	}

	result := database.DB.Model(&equipment).Updates(updates)
	if result.Error != nil {
		utils.Logger.Error("Failed to update equipment: %v", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update equipment")
		return
	}

	utils.Logger.Info("Equipment updated: %d by user %d", id, userID)
	utils.SuccessResponse(c, equipment)
}

func (h *EquipmentHandler) DeleteEquipment(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Equipment not found")
		return
	}

	if equipment.OwnerID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only delete your own equipment")
		return
	}

	var activeOrders int64
	database.DB.Model(&models.Order{}).Where("equipment_id = ? AND status IN ?", id, []string{"pending", "confirmed", "rented"}).Count(&activeOrders)
	if activeOrders > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot delete equipment with active orders")
		return
	}

	tx := database.DB.Begin()

	if err := tx.Where("equipment_id = ?", id).Delete(&models.EquipmentImage{}).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete images")
		return
	}

	if err := tx.Delete(&equipment).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete equipment")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit deletion")
		return
	}

	utils.Logger.Info("Equipment deleted: %d by user %d", id, userID)
	utils.SuccessResponse(c, nil)
}

func (h *EquipmentHandler) UploadImage(c *gin.Context) {
	userID := c.GetUint("userId")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Equipment not found")
		return
	}

	if equipment.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only upload images for your own equipment")
		return
	}

	var imageCount int64
	database.DB.Model(&models.EquipmentImage{}).Where("equipment_id = ?", id).Count(&imageCount)
	if imageCount >= 9 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Maximum 9 images allowed")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to get uploaded file")
		return
	}

	filename, err := utils.UploadImage(file)
	if err != nil {
		utils.Logger.Error("Failed to upload image: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	img := models.EquipmentImage{
		EquipmentID: uint(id),
		ImageURL:    filename,
		SortOrder:   int(imageCount),
	}

	if err := database.DB.Create(&img).Error; err != nil {
		utils.DeleteImage(filename)
		utils.Logger.Error("Failed to save image record: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save image")
		return
	}

	utils.Logger.Info("Image uploaded for equipment %d: %s", id, filename)
	utils.SuccessResponse(c, gin.H{
		"id":       img.ID,
		"imageUrl": filename,
	})
}

func (h *EquipmentHandler) DeleteImage(c *gin.Context) {
	userID := c.GetUint("userId")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid image ID")
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Equipment not found")
		return
	}

	if equipment.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only delete images for your own equipment")
		return
	}

	var img models.EquipmentImage
	if err := database.DB.First(&img, imageID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image not found")
		return
	}

	if img.EquipmentID != equipment.ID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Image does not belong to this equipment")
		return
	}

	utils.DeleteImage(img.ImageURL)

	if err := database.DB.Delete(&img).Error; err != nil {
		utils.Logger.Error("Failed to delete image record: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	utils.Logger.Info("Image deleted: %d from equipment %d", imageID, id)
	utils.SuccessResponse(c, nil)
}
