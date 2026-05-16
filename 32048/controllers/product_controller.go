package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"secondhand-trading/config"
	"secondhand-trading/models"
	"secondhand-trading/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateProductRequest struct {
	Title         string   `json:"title" binding:"required,max=200"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required,min=0"`
	OriginalPrice float64  `json:"original_price"`
	Condition     string   `json:"condition" binding:"required"`
	CategoryID    uint     `json:"category_id" binding:"required"`
	Tags          []string `json:"tags"`
	Location      string   `json:"location"`
}

type UpdateProductRequest struct {
	Title         string   `json:"title" binding:"omitempty,max=200"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"omitempty,min=0"`
	OriginalPrice float64  `json:"original_price"`
	Condition     string   `json:"condition"`
	CategoryID    uint     `json:"category_id"`
	Tags          []string `json:"tags"`
	Location      string   `json:"location"`
	Status        string   `json:"status"`
}

type ProductListQuery struct {
	Keyword    string  `form:"keyword"`
	CategoryID uint    `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	Condition  string  `form:"condition"`
	Page       int     `form:"page,default=1"`
	PageSize   int     `form:"page_size,default=20"`
	SortBy     string  `form:"sort_by,default=relevance"`
}

func CreateProduct(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.CreditScore < config.AppConfig.MinCreditScore {
		c.JSON(http.StatusForbidden, gin.H{"error": "Credit score too low to publish products"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Title:         req.Title,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Condition:     req.Condition,
		Status:        "on_sale",
		SellerID:      userID,
		CategoryID:    req.CategoryID,
		Tags:          strings.Join(req.Tags, ","),
		Location:      req.Location,
		IsReviewed:    true,
	}

	if err := models.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func UploadProductImages(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := c.GetUint("userID")

	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this product"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images uploaded"})
		return
	}

	if len(files) > 9 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 9 images allowed"})
		return
	}

	uploadDir := "./uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	var imageURLs []string
	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("%s_%d%s", uuid.New().String(), i, ext)
		filePath := filepath.Join(uploadDir, newFilename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imageURL := fmt.Sprintf("/uploads/products/%s", newFilename)
		isPrimary := i == 0

		productImage := models.ProductImage{
			ProductID: uint(productID),
			ImageURL:  imageURL,
			IsPrimary: isPrimary,
			SortOrder: i,
		}

		if err := models.DB.Create(&productImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image record"})
			return
		}

		imageURLs = append(imageURLs, imageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Images uploaded successfully",
		"image_urls": imageURLs,
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := models.DB.Preload("Images").Preload("Seller", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar, credit_score")
	}).Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	models.DB.Model(&product).UpdateColumn("view_count", product.ViewCount+1)

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := c.GetUint("userID")

	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this product"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldPrice := product.Price
	oldStatus := product.Status

	if req.Title != "" {
		product.Title = req.Title
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.OriginalPrice > 0 {
		product.OriginalPrice = req.OriginalPrice
	}
	if req.Condition != "" {
		product.Condition = req.Condition
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if len(req.Tags) > 0 {
		product.Tags = strings.Join(req.Tags, ",")
	}
	if req.Location != "" {
		product.Location = req.Location
	}
	if req.Status != "" {
		if oldStatus == "off_shelf" && req.Status == "on_sale" {
			product.NeedsReview = true
			product.IsReviewed = false
		}
		product.Status = req.Status
	}

	if err := models.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	if oldPrice > 0 && req.Price > 0 && req.Price < oldPrice {
		notifyPriceDrop(product.ID, oldPrice, req.Price)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

func notifyPriceDrop(productID uint, oldPrice, newPrice float64) {
	var favorites []models.Favorite
	models.DB.Where("product_id = ?", productID).Find(&favorites)

	for _, fav := range favorites {
		title := "商品降价通知"
		content := fmt.Sprintf("您收藏的商品价格从 %.2f 降至 %.2f", oldPrice, newPrice)
		utils.CreateNotification(fav.UserID, "price_drop", title, content, &productID)
	}
}

func DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := c.GetUint("userID")

	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this product"})
		return
	}

	if err := models.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func GetMyProducts(c *gin.Context) {
	userID := c.GetUint("userID")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := models.DB.Model(&models.Product{}).Where("seller_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	offset := (page - 1) * pageSize
	query.Preload("Images").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"products":   products,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func SearchProducts(c *gin.Context) {
	var query ProductListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbQuery := models.DB.Model(&models.Product{}).Where("status = ? AND is_reviewed = ?", "on_sale", true)

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		dbQuery = dbQuery.Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?", keyword, keyword, keyword)
	}

	if query.CategoryID > 0 {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
	}

	if query.MinPrice > 0 {
		dbQuery = dbQuery.Where("price >= ?", query.MinPrice)
	}

	if query.MaxPrice > 0 {
		dbQuery = dbQuery.Where("price <= ?", query.MaxPrice)
	}

	if query.Condition != "" {
		dbQuery = dbQuery.Where("condition = ?", query.Condition)
	}

	var total int64
	dbQuery.Count(&total)

	var products []models.Product
	offset := (query.Page - 1) * query.PageSize

	orderBy := "created_at DESC"
	switch query.SortBy {
	case "price_asc":
		orderBy = "price ASC"
	case "price_desc":
		orderBy = "price DESC"
	case "newest":
		orderBy = "created_at DESC"
	case "relevance":
		if query.Keyword != "" {
			escapedKeyword := escapeSQLLike(query.Keyword)
			titleMatch := "'" + escapedKeyword + "%'"
			anyMatch := "'%" + escapedKeyword + "%'"
			dbQuery = dbQuery.Order("CASE WHEN title ILIKE " + titleMatch + " THEN 1 WHEN description ILIKE " + anyMatch + " THEN 2 ELSE 3 END ASC, created_at DESC")
		}
	}

	dbQuery.Preload("Images", "is_primary = ?", true).
		Preload("Seller", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, avatar, credit_score")
		}).
		Preload("Category").
		Order(orderBy).
		Offset(offset).
		Limit(query.PageSize).
		Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"products":   products,
		"total":      total,
		"page":       query.Page,
		"page_size":  query.PageSize,
		"total_page": (total + int64(query.PageSize) - 1) / int64(query.PageSize),
	})
}

func RelistProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := c.GetUint("userID")

	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this product"})
		return
	}

	if product.Status != "off_shelf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only off-shelf products can be relisted"})
		return
	}

	product.Status = "on_sale"
	product.NeedsReview = true
	product.IsReviewed = false

	if err := models.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to relist product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product relisted successfully, waiting for review",
		"product": product,
	})
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	models.DB.Order("sort_order ASC, name ASC").Find(&categories)

	categoryTree := buildCategoryTree(categories, nil)

	c.JSON(http.StatusOK, gin.H{"categories": categoryTree})
}

func buildCategoryTree(categories []models.Category, parentID *uint) []map[string]interface{} {
	var result []map[string]interface{}

	for _, cat := range categories {
		if (parentID == nil && cat.ParentID == nil) || (parentID != nil && cat.ParentID != nil && *cat.ParentID == *parentID) {
			item := map[string]interface{}{
				"id":         cat.ID,
				"name":       cat.Name,
				"icon":       cat.Icon,
				"sort_order": cat.SortOrder,
			}

			children := buildCategoryTree(categories, &cat.ID)
			if len(children) > 0 {
				item["children"] = children
			}

			result = append(result, item)
		}
	}

	return result
}

func escapeSQLLike(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "'", "''")
	return s
}
