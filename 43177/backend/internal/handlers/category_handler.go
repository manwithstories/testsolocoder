package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"repair-platform/internal/models"
	"repair-platform/internal/utils"
	"repair-platform/pkg/logger"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
	Icon string `json:"icon"`
	Sort int    `json:"sort"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	category := models.Category{
		Name:   req.Name,
		Code:   req.Code,
		Icon:   req.Icon,
		Sort:   req.Sort,
		Status: true,
	}

	if err := models.DB.Create(&category).Error; err != nil {
		logger.Errorf("Create category error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "创建失败")
		return
	}

	utils.Success(c, category)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := models.DB.Where("status = ?", true).Order("sort ASC").Find(&categories).Error; err != nil {
		logger.Errorf("Get categories error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "获取失败")
		return
	}

	utils.Success(c, categories)
}

func (h *CategoryHandler) GetCategoryDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "分类不存在")
		return
	}

	utils.Success(c, category)
}

type UpdateCategoryRequest struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Sort   int    `json:"sort"`
	Status *bool  `json:"status"`
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := models.DB.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		logger.Errorf("Update category error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := models.DB.Delete(&models.Category{}, id).Error; err != nil {
		logger.Errorf("Delete category error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

type ServiceItemHandler struct{}

func NewServiceItemHandler() *ServiceItemHandler {
	return &ServiceItemHandler{}
}

type CreateServiceItemRequest struct {
	CategoryID    uint    `json:"category_id" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	MinPrice      float64 `json:"min_price" binding:"required"`
	MaxPrice      float64 `json:"max_price" binding:"required"`
	EstimatedTime int     `json:"estimated_time" binding:"required"`
	Image         string  `json:"image"`
	Sort          int     `json:"sort"`
}

func (h *ServiceItemHandler) CreateServiceItem(c *gin.Context) {
	var req CreateServiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	if req.MinPrice > req.MaxPrice {
		utils.Error(c, http.StatusBadRequest, 400, "最低价格不能高于最高价格")
		return
	}

	item := models.ServiceItem{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Description:   req.Description,
		MinPrice:      req.MinPrice,
		MaxPrice:      req.MaxPrice,
		EstimatedTime: req.EstimatedTime,
		Image:         req.Image,
		Sort:          req.Sort,
		Status:        true,
	}

	if err := models.DB.Create(&item).Error; err != nil {
		logger.Errorf("Create service item error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "创建失败")
		return
	}

	utils.Success(c, item)
}

func (h *ServiceItemHandler) GetServiceItems(c *gin.Context) {
	categoryID := c.Query("category_id")

	var items []models.ServiceItem
	query := models.DB.Preload("Category").Where("status = ?", true)

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Order("sort ASC").Find(&items).Error; err != nil {
		logger.Errorf("Get service items error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "获取失败")
		return
	}

	utils.Success(c, items)
}

func (h *ServiceItemHandler) GetServiceItemDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var item models.ServiceItem
	if err := models.DB.Preload("Category").First(&item, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "服务项目不存在")
		return
	}

	utils.Success(c, item)
}

type UpdateServiceItemRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	EstimatedTime int     `json:"estimated_time"`
	Image         string  `json:"image"`
	Sort          int     `json:"sort"`
	Status        *bool   `json:"status"`
}

func (h *ServiceItemHandler) UpdateServiceItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateServiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.MinPrice != 0 {
		updates["min_price"] = req.MinPrice
	}
	if req.MaxPrice != 0 {
		updates["max_price"] = req.MaxPrice
	}
	if req.EstimatedTime != 0 {
		updates["estimated_time"] = req.EstimatedTime
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := models.DB.Model(&models.ServiceItem{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		logger.Errorf("Update service item error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

func (h *ServiceItemHandler) DeleteServiceItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := models.DB.Delete(&models.ServiceItem{}, id).Error; err != nil {
		logger.Errorf("Delete service item error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}
