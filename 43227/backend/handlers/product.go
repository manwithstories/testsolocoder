package handlers

import (
	"net/http"
	"strconv"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var inventory models.Inventory
	if err := database.DB.Where("id = ? AND user_id = ?", req.InventoryID, userID).First(&inventory).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inventory not found")
		return
	}

	if req.Stock > inventory.Quantity {
		utils.Fail(c, http.StatusBadRequest, "stock cannot exceed inventory quantity")
		return
	}

	product := models.Product{
		UserID:      userID.(uint),
		InventoryID: req.InventoryID,
		Title:       req.Title,
		Description: req.Description,
		HoneyType:   req.HoneyType,
		BatchCode:   req.BatchCode,
		Price:       req.Price,
		Stock:       req.Stock,
		Unit:        req.Unit,
		Images:      req.Images,
		Grade:       req.Grade,
		Status:      "on_sale",
	}

	if err := database.DB.Create(&product).Error; err != nil {
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create product", err)
		return
	}

	utils.Success(c, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	query := database.DB.Model(&models.Product{}).Where("status = ?", "on_sale")

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if honeyType := c.Query("honey_type"); honeyType != "" {
		query = query.Where("honey_type = ?", honeyType)
	}
	if grade := c.Query("grade"); grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		price, _ := strconv.ParseFloat(minPrice, 64)
		query = query.Where("price >= ?", price)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		price, _ := strconv.ParseFloat(maxPrice, 64)
		query = query.Where("price <= ?", price)
	}
	if sellerID := c.Query("seller_id"); sellerID != "" {
		query = query.Where("user_id = ?", sellerID)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var products []models.Product
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("User").Preload("Inventory.Harvest.Beehive").Find(&products)

	utils.SuccessWithTotal(c, products, total)
}

func (h *ProductHandler) ListMyProducts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	query := database.DB.Model(&models.Product{}).Where("user_id = ?", userID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var products []models.Product
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("Inventory.Harvest.Beehive").Find(&products)

	utils.SuccessWithTotal(c, products, total)
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var product models.Product
	if err := database.DB.Preload("User").Preload("Inventory.Harvest.Beehive").First(&product, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "product not found")
		return
	}

	database.DB.Model(&product).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	utils.Success(c, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&product).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "product not found")
		return
	}

	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Images != nil {
		product.Images = *req.Images
	}
	if req.Status != nil {
		product.Status = *req.Status
	}

	database.DB.Save(&product)

	utils.Success(c, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Product{})
	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to delete product")
		return
	}
	if result.RowsAffected == 0 {
		utils.Fail(c, http.StatusNotFound, "product not found")
		return
	}

	utils.Success(c, nil)
}
