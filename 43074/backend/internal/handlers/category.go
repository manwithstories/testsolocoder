package handlers

import (
	"net/http"

	"booklibrary/internal/database"
	"booklibrary/internal/errors"
	"booklibrary/internal/logger"
	"booklibrary/internal/models"
	"booklibrary/internal/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

type CreateCategoryRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *uint64 `json:"parent_id"`
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	flat := utils.GetStringQuery(c, "flat", "false") == "true"

	var categories []models.Category
	if flat {
		database.DB.Order("name asc").Find(&categories)
		c.JSON(http.StatusOK, categories)
		return
	}

	var rootCategories []models.Category
	database.DB.Where("parent_id IS NULL").
		Preload("Children.Children").
		Order("name asc").
		Find(&rootCategories)

	c.JSON(http.StatusOK, rootCategories)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var category models.Category
	result := database.DB.Preload("Children").First(&category, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	if req.ParentID != nil {
		var parent models.Category
		result := database.DB.First(&parent, *req.ParentID)
		if result.Error != nil {
			errors.ErrorResponse(c, http.StatusBadRequest, "父分类不存在")
			return
		}
	}

	category := models.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	result := database.DB.Create(&category)
	if result.Error != nil {
		logger.Errorf("Create category failed: %v", result.Error)
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	logger.Infof("Category created: %d - %s", category.ID, category.Name)
	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var category models.Category
	result := database.DB.First(&category, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			errors.ErrorResponse(c, http.StatusBadRequest, "不能将自己设为父分类")
			return
		}
		var parent models.Category
		result := database.DB.First(&parent, *req.ParentID)
		if result.Error != nil {
			errors.ErrorResponse(c, http.StatusBadRequest, "父分类不存在")
			return
		}
		category.ParentID = req.ParentID
	} else {
		category.ParentID = nil
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	database.DB.Save(&category)
	logger.Infof("Category updated: %d - %s", category.ID, category.Name)
	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var category models.Category
	result := database.DB.First(&category, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	tx := database.DB.Begin()

	var childCount int64
	tx.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		tx.Rollback()
		errors.ErrorResponse(c, http.StatusBadRequest, "该分类下有子分类，请先删除子分类")
		return
	}

	tx.Model(&category).Association("Books").Clear()
	tx.Delete(&category)
	tx.Commit()

	logger.Infof("Category deleted: %d", id)
	c.JSON(http.StatusNoContent, nil)
}
