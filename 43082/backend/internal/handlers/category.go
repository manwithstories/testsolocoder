package handlers

import (
	"multishop/internal/database"
	"multishop/internal/dto"
	"multishop/internal/models"
	"multishop/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	level := 1
	if req.ParentID != nil {
		var parent models.Category
		if err := database.DB.First(&parent, *req.ParentID).Error; err != nil {
			utils.Error(c, http.StatusNotFound, "父分类不存在")
			return
		}
		level = parent.Level + 1
	}

	category := models.Category{
		Name:     req.Name,
		Icon:     req.Icon,
		ParentID: req.ParentID,
		Level:    level,
		Sort:     req.Sort,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	utils.Success(c, gin.H{"category_id": category.ID})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Sort > 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := database.DB.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.Success(c, nil)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var count int64
	database.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		utils.Error(c, http.StatusBadRequest, "该分类下还有商品，无法删除")
		return
	}

	database.DB.Delete(&models.Category{}, id)
	utils.Success(c, nil)
}

func (h *CategoryHandler) List(c *gin.Context) {
	var categories []models.Category
	database.DB.Where("status = ?", "active").Order("sort ASC, id ASC").Find(&categories)

	categoryMap := make(map[uint]*dto.CategoryInfo)
	rootCategories := make([]dto.CategoryInfo, 0)

	for _, cat := range categories {
		info := dto.CategoryInfo{
			ID:       cat.ID,
			Name:     cat.Name,
			Icon:     cat.Icon,
			ParentID: cat.ParentID,
			Level:    cat.Level,
			Sort:     cat.Sort,
			Status:   cat.Status,
			Children: make([]dto.CategoryInfo, 0),
		}
		categoryMap[cat.ID] = &info

		if cat.ParentID == nil {
			rootCategories = append(rootCategories, info)
		} else {
			if parent, ok := categoryMap[*cat.ParentID]; ok {
				parent.Children = append(parent.Children, info)
			}
		}
	}

	utils.Success(c, rootCategories)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	var categories []models.Category
	database.DB.Order("sort ASC, id ASC").Find(&categories)

	result := make([]dto.CategoryInfo, 0, len(categories))
	for _, cat := range categories {
		result = append(result, dto.CategoryInfo{
			ID:       cat.ID,
			Name:     cat.Name,
			Icon:     cat.Icon,
			ParentID: cat.ParentID,
			Level:    cat.Level,
			Sort:     cat.Sort,
			Status:   cat.Status,
		})
	}

	utils.Success(c, result)
}
