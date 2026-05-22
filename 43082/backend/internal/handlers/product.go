package handlers

import (
	"encoding/json"
	"multishop/internal/database"
	"multishop/internal/dto"
	"multishop/internal/middleware"
	"multishop/internal/models"
	"multishop/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var shop models.Shop
	if err := database.DB.Where("user_id = ? AND status = ?", userID, models.ShopStatusApproved).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "请先申请并通过店铺审核")
		return
	}

	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	tx := database.DB.Begin()

	product := models.Product{
		ShopID:      shop.ID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		MainImage:   req.MainImage,
		Price:       req.Price,
		Stock:       req.Stock,
		Weight:      req.Weight,
		Status:      models.ProductStatusDraft,
		IsHot:       req.IsHot,
		IsRecommend: req.IsRecommend,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, "创建商品失败")
		return
	}

	for i, imgURL := range req.Images {
		img := models.ProductImage{
			ProductID: product.ID,
			URL:       imgURL,
			Sort:      i,
		}
		if err := tx.Create(&img).Error; err != nil {
			tx.Rollback()
			utils.Error(c, http.StatusInternalServerError, "保存图片失败")
			return
		}
	}

	for _, spec := range req.Specs {
		valuesJSON, _ := json.Marshal(spec.Values)
		productSpec := models.ProductSpec{
			ProductID: product.ID,
			Name:      spec.Name,
			Values:    string(valuesJSON),
		}
		if err := tx.Create(&productSpec).Error; err != nil {
			tx.Rollback()
			utils.Error(c, http.StatusInternalServerError, "保存规格失败")
			return
		}
	}

	for _, sku := range req.Skus {
		specsJSON, _ := json.Marshal(sku.Specs)
		skuModel := models.SKU{
			ProductID: product.ID,
			Specs:     string(specsJSON),
			Price:     sku.Price,
			Stock:     sku.Stock,
			SKUCode:   sku.SKUCode,
		}
		if err := tx.Create(&skuModel).Error; err != nil {
			tx.Rollback()
			utils.Error(c, http.StatusInternalServerError, "保存SKU失败")
			return
		}
	}

	tx.Commit()
	utils.Success(c, gin.H{"product_id": product.ID})
}

func (h *ProductHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var shop models.Shop
	if err := database.DB.Where("user_id = ? AND status = ?", userID, models.ShopStatusApproved).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "店铺未审核通过")
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND shop_id = ?", id, shop.ID).First(&product).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.MainImage != "" {
		updates["main_image"] = req.MainImage
	}
	if req.Price >= 0 {
		updates["price"] = req.Price
	}
	if req.Stock >= 0 {
		updates["stock"] = req.Stock
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.IsHot != nil {
		updates["is_hot"] = *req.IsHot
	}
	if req.IsRecommend != nil {
		updates["is_recommend"] = *req.IsRecommend
	}

	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.Success(c, nil)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var shop models.Shop
	database.DB.Where("user_id = ?", userID).First(&shop)

	var product models.Product
	if err := database.DB.Where("id = ? AND shop_id = ?", id, shop.ID).First(&product).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	database.DB.Delete(&product)
	utils.Success(c, nil)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var product models.Product
	if err := database.DB.Preload("Shop").Preload("Category").Preload("Images").Preload("Specs").Preload("Skus").First(&product, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	images := make([]string, 0, len(product.Images))
	for _, img := range product.Images {
		images = append(images, img.URL)
	}

	specs := make([]dto.ProductSpecInfo, 0, len(product.Specs))
	for _, spec := range product.Specs {
		var values []string
		json.Unmarshal([]byte(spec.Values), &values)
		specs = append(specs, dto.ProductSpecInfo{
			ID:     spec.ID,
			Name:   spec.Name,
			Values: values,
		})
	}

	skus := make([]dto.SKUInfo, 0, len(product.Skus))
	for _, sku := range product.Skus {
		var specMap map[string]string
		json.Unmarshal([]byte(sku.Specs), &specMap)
		skus = append(skus, dto.SKUInfo{
			ID:      sku.ID,
			Specs:   specMap,
			Price:   sku.Price,
			Stock:   sku.Stock,
			SKUCode: sku.SKUCode,
		})
	}

	utils.Success(c, dto.ProductDetail{
		ProductInfo: dto.ProductInfo{
			ID:           product.ID,
			ShopID:       product.ShopID,
			CategoryID:   product.CategoryID,
			Name:         product.Name,
			MainImage:    product.MainImage,
			Price:        product.Price,
			Stock:        product.Stock,
			Sales:        product.Sales,
			Status:       product.Status,
			IsHot:        product.IsHot,
			IsRecommend:  product.IsRecommend,
			CreatedAt:    product.CreatedAt.Format(time.RFC3339),
			ShopName:     product.Shop.Name,
			CategoryName: product.Category.Name,
		},
		Description: product.Description,
		Images:      images,
		Specs:       specs,
		Skus:        skus,
		Shop: dto.ShopInfo{
			ID:          product.Shop.ID,
			Name:        product.Shop.Name,
			Logo:        product.Shop.Logo,
			Rating:      product.Shop.Rating,
		},
	})
}

func (h *ProductHandler) List(c *gin.Context) {
	var query dto.ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := database.DB.Model(&models.Product{}).Where("status = ?", models.ProductStatusOnSale)

	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}
	if query.CategoryID != nil {
		db = db.Where("category_id = ?", *query.CategoryID)
	}
	if query.ShopID != nil {
		db = db.Where("shop_id = ?", *query.ShopID)
	}
	if query.MinPrice != nil {
		db = db.Where("price >= ?", *query.MinPrice)
	}
	if query.MaxPrice != nil {
		db = db.Where("price <= ?", *query.MaxPrice)
	}
	if query.IsHot != nil {
		db = db.Where("is_hot = ?", *query.IsHot)
	}

	var total int64
	db.Count(&total)

	sortField := "created_at"
	sortOrder := "DESC"
	if query.SortBy != "" {
		sortField = query.SortBy
	}
	if query.SortOrder != "" {
		sortOrder = query.SortOrder
	}

	var products []models.Product
	offset := (query.Page - 1) * query.PageSize
	db.Preload("Shop").Preload("Category").
		Offset(offset).Limit(query.PageSize).
		Order(sortField + " " + sortOrder).
		Find(&products)

	productInfos := make([]dto.ProductInfo, 0, len(products))
	for _, p := range products {
		productInfos = append(productInfos, dto.ProductInfo{
			ID:           p.ID,
			ShopID:       p.ShopID,
			CategoryID:   p.CategoryID,
			Name:         p.Name,
			MainImage:    p.MainImage,
			Price:        p.Price,
			Stock:        p.Stock,
			Sales:        p.Sales,
			Status:       p.Status,
			IsHot:        p.IsHot,
			IsRecommend:  p.IsRecommend,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
			ShopName:     p.Shop.Name,
			CategoryName: p.Category.Name,
		})
	}

	utils.Paginated(c, productInfos, total, query.Page, query.PageSize)
}

func (h *ProductHandler) MyProducts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	query := database.DB.Model(&models.Product{}).Where("shop_id = ?", shop.ID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	offset := (page - 1) * pageSize
	query.Preload("Category").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products)

	productInfos := make([]dto.ProductInfo, 0, len(products))
	for _, p := range products {
		productInfos = append(productInfos, dto.ProductInfo{
			ID:           p.ID,
			ShopID:       p.ShopID,
			CategoryID:   p.CategoryID,
			Name:         p.Name,
			MainImage:    p.MainImage,
			Price:        p.Price,
			Stock:        p.Stock,
			Sales:        p.Sales,
			Status:       p.Status,
			IsHot:        p.IsHot,
			IsRecommend:  p.IsRecommend,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
			ShopName:     shop.Name,
			CategoryName: p.Category.Name,
		})
	}

	utils.Paginated(c, productInfos, total, page, pageSize)
}
