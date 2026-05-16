package controllers

import (
	"fmt"

	"finance-api/middleware"
	"finance-api/models"
	"finance-api/utils"

	"github.com/gin-gonic/gin"
)

type CategoryController struct{}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Type     string `json:"type" binding:"required,oneof=income expense"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"omitempty,max=50"`
}

func (ctrl *CategoryController) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	if req.ParentID != nil {
		var parent models.Category
		if result := utils.DB.Where("id = ? AND user_id = ? AND type = ?", *req.ParentID, userID, req.Type).First(&parent); result.Error != nil {
			utils.BadRequest(c, "Parent category not found or type mismatch")
			return
		}
		if parent.ParentID != nil {
			utils.BadRequest(c, "Cannot create subcategory under a subcategory (max 2 levels)")
			return
		}
	}

	category := models.Category{
		UserID:   userID,
		Name:     req.Name,
		Type:     req.Type,
		ParentID: req.ParentID,
	}

	if result := utils.DB.Create(&category); result.Error != nil {
		utils.InternalError(c, "Failed to create category: "+result.Error.Error())
		return
	}

	utils.Success(c, category)
}

func (ctrl *CategoryController) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	categoryType := c.Query("type")

	query := utils.DB.Where("user_id = ?", userID)
	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	var categories []models.Category
	if result := query.Find(&categories); result.Error != nil {
		utils.InternalError(c, "Failed to fetch categories: "+result.Error.Error())
		return
	}

	tree := buildCategoryTree(categories)
	utils.Success(c, tree)
}

func (ctrl *CategoryController) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var category models.Category
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&category); result.Error != nil {
		utils.NotFound(c, "Category not found")
		return
	}

	utils.Success(c, category)
}

func (ctrl *CategoryController) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var category models.Category
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&category); result.Error != nil {
		utils.NotFound(c, "Category not found")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if result := utils.DB.Save(&category); result.Error != nil {
		utils.InternalError(c, "Failed to update category: "+result.Error.Error())
		return
	}

	utils.Success(c, category)
}

func (ctrl *CategoryController) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var category models.Category
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&category); result.Error != nil {
		utils.NotFound(c, "Category not found")
		return
	}

	var childCount int64
	utils.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		utils.BadRequest(c, fmt.Sprintf("Cannot delete category: there are %d subcategories under it", childCount))
		return
	}

	var transactionCount int64
	utils.DB.Model(&models.Transaction{}).Where("category_id = ?", id).Count(&transactionCount)
	if transactionCount > 0 {
		utils.BadRequest(c, fmt.Sprintf("Cannot delete category: there are %d associated transaction records", transactionCount))
		return
	}

	if result := utils.DB.Delete(&category); result.Error != nil {
		utils.InternalError(c, "Failed to delete category: "+result.Error.Error())
		return
	}

	utils.Success(c, nil)
}

func buildCategoryTree(categories []models.Category) []models.Category {
	categoryMap := make(map[uint]*models.Category)
	var roots []models.Category

	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	for i := range categories {
		if categories[i].ParentID == nil {
			roots = append(roots, categories[i])
		} else {
			parent := categoryMap[*categories[i].ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, categories[i])
			}
		}
	}

	return roots
}
