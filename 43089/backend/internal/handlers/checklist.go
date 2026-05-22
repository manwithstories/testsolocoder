package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"travel-planner/internal/database"
	"travel-planner/internal/logger"
	"travel-planner/internal/middleware"
	"travel-planner/internal/models"
	"travel-planner/internal/utils"
)

type CreateChecklistRequest struct {
	Title string `json:"title" validate:"required,max=200"`
	Type  string `json:"type" validate:"oneof=packing preparation other"`
}

type CreateChecklistItemRequest struct {
	Title       string `json:"title" validate:"required,max=200"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Quantity    int    `json:"quantity"`
	OrderIndex  int    `json:"order_index"`
}

func GetChecklists(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var checklists []models.Checklist
	if err := database.DB.Where("plan_id = ?", planUUID).
		Preload("Items").
		Order("created_at DESC").
		Find(&checklists).Error; err != nil {
		logger.Errorf("Failed to get checklists: %v", err)
		utils.InternalServerError(c, "Failed to get checklists")
		return
	}

	utils.Success(c, checklists)
}

func GetChecklist(c *gin.Context) {
	checklistID := c.Param("id")
	checklistUUID, err := uuid.Parse(checklistID)
	if err != nil {
		utils.BadRequest(c, "Invalid checklist ID")
		return
	}

	var checklist models.Checklist
	if err := database.DB.Preload("Items").First(&checklist, "id = ?", checklistUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Checklist not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	utils.Success(c, checklist)
}

func CreateChecklist(c *gin.Context) {
	planID := c.Param("plan_id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	userID := middleware.GetCurrentUserID(c)

	var req CreateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	checklist := models.Checklist{
		PlanID:    planUUID,
		Title:     req.Title,
		Type:      req.Type,
		CreatedBy: userID,
	}

	if err := database.DB.Create(&checklist).Error; err != nil {
		logger.Errorf("Failed to create checklist: %v", err)
		utils.InternalServerError(c, "Failed to create checklist")
		return
	}

	logger.Infof("Checklist created: %s in plan %s", checklist.ID, planID)
	utils.Created(c, checklist)
}

func UpdateChecklist(c *gin.Context) {
	checklistID := c.Param("id")
	checklistUUID, err := uuid.Parse(checklistID)
	if err != nil {
		utils.BadRequest(c, "Invalid checklist ID")
		return
	}

	var checklist models.Checklist
	if err := database.DB.First(&checklist, "id = ?", checklistUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Checklist not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var req CreateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}

	if err := database.DB.Model(&checklist).Updates(updates).Error; err != nil {
		logger.Errorf("Failed to update checklist: %v", err)
		utils.InternalServerError(c, "Failed to update checklist")
		return
	}

	var updatedChecklist models.Checklist
	database.DB.Preload("Items").First(&updatedChecklist, checklistUUID)

	logger.Infof("Checklist updated: %s", checklistID)
	utils.Success(c, updatedChecklist)
}

func DeleteChecklist(c *gin.Context) {
	checklistID := c.Param("id")
	checklistUUID, err := uuid.Parse(checklistID)
	if err != nil {
		utils.BadRequest(c, "Invalid checklist ID")
		return
	}

	tx := database.BeginTransaction()
	if tx.Error != nil {
		utils.InternalServerError(c, "Database error")
		return
	}

	if err := tx.Where("checklist_id = ?", checklistUUID).Delete(&models.ChecklistItem{}).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete checklist items: %v", err)
		utils.InternalServerError(c, "Failed to delete checklist")
		return
	}

	if err := tx.Delete(&models.Checklist{}, checklistUUID).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to delete checklist: %v", err)
		utils.InternalServerError(c, "Failed to delete checklist")
		return
	}

	if err := tx.Commit().Error; err != nil {
		logger.Errorf("Failed to commit transaction: %v", err)
		utils.InternalServerError(c, "Failed to delete checklist")
		return
	}

	logger.Infof("Checklist deleted: %s", checklistID)
	utils.Success(c, gin.H{"message": "Checklist deleted successfully"})
}

func AddChecklistItem(c *gin.Context) {
	checklistID := c.Param("id")
	checklistUUID, err := uuid.Parse(checklistID)
	if err != nil {
		utils.BadRequest(c, "Invalid checklist ID")
		return
	}

	var req CreateChecklistItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	item := models.ChecklistItem{
		ChecklistID: checklistUUID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Quantity:    req.Quantity,
		OrderIndex:  req.OrderIndex,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		logger.Errorf("Failed to create checklist item: %v", err)
		utils.InternalServerError(c, "Failed to create checklist item")
		return
	}

	logger.Infof("Checklist item created: %s", item.ID)
	utils.Created(c, item)
}

func UpdateChecklistItem(c *gin.Context) {
	itemID := c.Param("item_id")
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		utils.BadRequest(c, "Invalid item ID")
		return
	}

	var item models.ChecklistItem
	if err := database.DB.First(&item, "id = ?", itemUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Checklist item not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	allowedFields := map[string]bool{
		"title":        true,
		"description":  true,
		"category":     true,
		"quantity":     true,
		"is_completed": true,
		"order_index":  true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if isCompleted, ok := filteredUpdates["is_completed"].(bool); ok && isCompleted {
		userID := middleware.GetCurrentUserID(c)
		now := time.Now()
		filteredUpdates["completed_by"] = userID
		filteredUpdates["completed_at"] = &now
	}

	if err := database.DB.Model(&item).Updates(filteredUpdates).Error; err != nil {
		logger.Errorf("Failed to update checklist item: %v", err)
		utils.InternalServerError(c, "Failed to update checklist item")
		return
	}

	var updatedItem models.ChecklistItem
	database.DB.First(&updatedItem, itemUUID)

	logger.Infof("Checklist item updated: %s", itemID)
	utils.Success(c, updatedItem)
}

func DeleteChecklistItem(c *gin.Context) {
	itemID := c.Param("item_id")
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		utils.BadRequest(c, "Invalid item ID")
		return
	}

	if err := database.DB.Delete(&models.ChecklistItem{}, itemUUID).Error; err != nil {
		logger.Errorf("Failed to delete checklist item: %v", err)
		utils.InternalServerError(c, "Failed to delete checklist item")
		return
	}

	logger.Infof("Checklist item deleted: %s", itemID)
	utils.Success(c, gin.H{"message": "Checklist item deleted successfully"})
}
