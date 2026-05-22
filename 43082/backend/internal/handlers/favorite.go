package handlers

import (
	"multishop/internal/database"
	"multishop/internal/dto"
	"multishop/internal/middleware"
	"multishop/internal/models"
	"multishop/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct{}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{}
}

func (h *FavoriteHandler) ToggleShop(c *gin.Context) {
	userID := middleware.GetUserID(c)
	shopID64, _ := strconv.ParseUint(c.Param("shop_id"), 10, 32)
	shopID := uint(shopID64)

	var favorite models.Favorite
	err := database.DB.Where("user_id = ? AND shop_id = ?", userID, shopID).First(&favorite).Error

	if err == nil {
		database.DB.Delete(&favorite)
		utils.Success(c, gin.H{"favorited": false})
	} else {
		favorite = models.Favorite{
			UserID: userID,
			ShopID: &shopID,
		}
		database.DB.Create(&favorite)
		utils.Success(c, gin.H{"favorited": true})
	}
}

func (h *FavoriteHandler) ToggleProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)
	productID, _ := strconv.ParseUint(c.Param("product_id"), 10, 32)

	var favorite models.Favorite
	err := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&favorite).Error

	if err == nil {
		database.DB.Delete(&favorite)
		utils.Success(c, gin.H{"favorited": false})
	} else {
		pid := uint(productID)
		favorite = models.Favorite{
			UserID:    userID,
			ProductID: &pid,
		}
		database.DB.Create(&favorite)
		utils.Success(c, gin.H{"favorited": true})
	}
}

func (h *FavoriteHandler) GetMyShops(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var favorites []models.Favorite
	query := database.DB.Where("user_id = ? AND shop_id IS NOT NULL", userID)
	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Shop").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&favorites)

	result := make([]dto.ShopInfo, 0, len(favorites))
	for _, f := range favorites {
		if f.ShopID != nil {
			var shop models.Shop
			database.DB.First(&shop, *f.ShopID)
			result = append(result, dto.ShopInfo{
				ID:          shop.ID,
				Name:        shop.Name,
				Description: shop.Description,
				Logo:        shop.Logo,
				Rating:      shop.Rating,
				CreatedAt:   shop.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	utils.Paginated(c, result, total, page, pageSize)
}

func (h *FavoriteHandler) GetMyProducts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var favorites []models.Favorite
	query := database.DB.Where("user_id = ? AND product_id IS NOT NULL", userID)
	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&favorites)

	result := make([]dto.ProductInfo, 0, len(favorites))
	for _, f := range favorites {
		if f.ProductID != nil {
			var product models.Product
			database.DB.Preload("Shop").Preload("Category").First(&product, *f.ProductID)
			result = append(result, dto.ProductInfo{
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
			})
		}
	}

	utils.Paginated(c, result, total, page, pageSize)
}
