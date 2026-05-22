package handlers

import (
	"strconv"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
)

func GetServiceCategories(c *gin.Context) {
	var categories []models.ServiceCategory
	config.DB.Where("is_active = ?", true).Order("sort_order ASC").Find(&categories)

	utils.Success(c, categories)
}

func CreateServiceCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var existing models.ServiceCategory
	if result := config.DB.Where("name = ?", req.Name).First(&existing); result.Error == nil {
		utils.BadRequest(c, "分类名称已存在")
		return
	}

	category := models.ServiceCategory{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	if result := config.DB.Create(&category); result.Error != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	operatorID := c.GetUint("user_id")
	go utils.LogOperation(operatorID, "admin", "service", "create_category", &category.ID, "service_category", "创建服务分类", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, category)
}

func UpdateServiceCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Icon        *string `json:"icon"`
		SortOrder   *int    `json:"sort_order"`
		IsActive    *bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	config.DB.Model(&models.ServiceCategory{}).Where("id = ?", id).Updates(updates)

	operatorID := c.GetUint("user_id")
	uid := uint(id)
	go utils.LogOperation(operatorID, "admin", "service", "update_category", &uid, "service_category", "更新服务分类", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func DeleteServiceCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var serviceItemCount int64
	config.DB.Model(&models.ServiceItem{}).Where("category_id = ?", id).Count(&serviceItemCount)
	if serviceItemCount > 0 {
		utils.BadRequest(c, "该分类下存在服务项目，无法删除")
		return
	}

	config.DB.Delete(&models.ServiceCategory{}, id)

	operatorID := c.GetUint("user_id")
	uid := uint(id)
	go utils.LogOperation(operatorID, "admin", "service", "delete_category", &uid, "service_category", "删除服务分类", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func GetServiceList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	minRating := c.Query("min_rating")
	keyword := c.Query("keyword")

	query := config.DB.Model(&models.ServiceItem{}).Where("is_active = ?", true)

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if minPrice != "" {
		query = query.Where("base_price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("base_price <= ?", maxPrice)
	}
	if minRating != "" {
		query = query.Where("rating >= ?", minRating)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var services []models.ServiceItem
	offset := (page - 1) * pageSize
	query.Preload("Category").Preload("Provider").Offset(offset).Limit(pageSize).Order("rating DESC, order_count DESC").Find(&services)

	utils.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      services,
	})
}

func GetServiceDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var service models.ServiceItem
	if result := config.DB.Preload("Category").Preload("Provider").Preload("ServiceAreas").First(&service, id); result.Error != nil {
		utils.NotFound(c, "服务不存在")
		return
	}

	utils.Success(c, service)
}

func CreateServiceItem(c *gin.Context) {
	providerID := c.GetUint("user_id")

	var req struct {
		CategoryID    uint     `json:"category_id" binding:"required"`
		Name          string   `json:"name" binding:"required"`
		Description   string   `json:"description"`
		Images        string   `json:"images"`
		BasePrice     float64  `json:"base_price" binding:"required"`
		PriceUnit     string   `json:"price_unit"`
		MinDuration   int      `json:"min_duration"`
		MaxDuration   int      `json:"max_duration"`
		ServiceAreaIDs []uint  `json:"service_area_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.ValidatePrice(req.BasePrice) {
		utils.BadRequest(c, "价格不能为负数")
		return
	}

	service := models.ServiceItem{
		CategoryID:  req.CategoryID,
		ProviderID:  providerID,
		Name:        req.Name,
		Description: req.Description,
		Images:      req.Images,
		BasePrice:   req.BasePrice,
		PriceUnit:   req.PriceUnit,
		MinDuration: req.MinDuration,
		MaxDuration: req.MaxDuration,
		IsActive:    true,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&service).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "创建失败")
		return
	}

	if len(req.ServiceAreaIDs) > 0 {
		var areas []models.ServiceArea
		tx.Find(&areas, req.ServiceAreaIDs)
		tx.Model(&service).Association("ServiceAreas").Append(areas)
	}

	tx.Commit()

	go utils.LogOperation(providerID, "service_provider", "service", "create_service", &service.ID, "service_item", "创建服务项目", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, service)
}

func UpdateServiceItem(c *gin.Context) {
	providerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var service models.ServiceItem
	if result := config.DB.Where("id = ? AND provider_id = ?", id, providerID).First(&service); result.Error != nil {
		utils.NotFound(c, "服务不存在")
		return
	}

	var req struct {
		Name          *string  `json:"name"`
		Description   *string  `json:"description"`
		Images        *string  `json:"images"`
		BasePrice     *float64 `json:"base_price"`
		PriceUnit     *string  `json:"price_unit"`
		MinDuration   *int     `json:"min_duration"`
		MaxDuration   *int     `json:"max_duration"`
		IsActive      *bool    `json:"is_active"`
		ServiceAreaIDs []uint  `json:"service_area_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Images != nil {
		updates["images"] = *req.Images
	}
	if req.BasePrice != nil {
		updates["base_price"] = *req.BasePrice
	}
	if req.PriceUnit != nil {
		updates["price_unit"] = *req.PriceUnit
	}
	if req.MinDuration != nil {
		updates["min_duration"] = *req.MinDuration
	}
	if req.MaxDuration != nil {
		updates["max_duration"] = *req.MaxDuration
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	tx := config.DB.Begin()
	if err := tx.Model(&service).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "更新失败")
		return
	}

	if req.ServiceAreaIDs != nil {
		var areas []models.ServiceArea
		tx.Find(&areas, req.ServiceAreaIDs)
		tx.Model(&service).Association("ServiceAreas").Replace(areas)
	}

	tx.Commit()

	uid := uint(id)
	go utils.LogOperation(providerID, "service_provider", "service", "update_service", &uid, "service_item", "更新服务项目", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func DeleteServiceItem(c *gin.Context) {
	providerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var service models.ServiceItem
	if result := config.DB.Where("id = ? AND provider_id = ?", id, providerID).First(&service); result.Error != nil {
		utils.NotFound(c, "服务不存在")
		return
	}

	config.DB.Delete(&service)

	uid := uint(id)
	go utils.LogOperation(providerID, "service_provider", "service", "delete_service", &uid, "service_item", "删除服务项目", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func GetMyServices(c *gin.Context) {
	providerID := c.GetUint("user_id")

	var services []models.ServiceItem
	config.DB.Where("provider_id = ?", providerID).Preload("Category").Order("id DESC").Find(&services)

	utils.Success(c, services)
}

func GetServiceAreas(c *gin.Context) {
	var areas []models.ServiceArea
	config.DB.Find(&areas)

	utils.Success(c, areas)
}

func CreateServiceArea(c *gin.Context) {
	var req struct {
		Province string `json:"province" binding:"required"`
		City     string `json:"city" binding:"required"`
		District string `json:"district" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	area := models.ServiceArea{
		Province: req.Province,
		City:     req.City,
		District: req.District,
	}

	if result := config.DB.Create(&area); result.Error != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.Success(c, area)
}
