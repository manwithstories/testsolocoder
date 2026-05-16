package controllers

import (
	"net/http"
	"secondhand-trading/config"
	"secondhand-trading/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddFavorite(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := models.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var existingFavorite models.Favorite
	if models.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingFavorite).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already in favorites"})
		return
	}

	var count int64
	models.DB.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&count)
	if count >= int64(config.AppConfig.MaxFavorites) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Maximum number of favorites reached",
			"max":   config.AppConfig.MaxFavorites,
		})
		return
	}

	favorite := models.Favorite{
		UserID:    userID,
		ProductID: req.ProductID,
	}

	if err := models.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Favorite added successfully",
		"favorite": favorite,
	})
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("userID")
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var favorite models.Favorite
	if err := models.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&favorite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}

	if err := models.DB.Delete(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}

func GetFavorites(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := models.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)

	var total int64
	query.Count(&total)

	var favorites []models.Favorite
	offset := (page - 1) * pageSize
	query.Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Images", "is_primary = ?", true).
			Preload("Seller", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, username, avatar")
			})
	}).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&favorites)

	c.JSON(http.StatusOK, gin.H{
		"favorites":   favorites,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_page":  (total + int64(pageSize) - 1) / int64(pageSize),
		"max_allowed": config.AppConfig.MaxFavorites,
	})
}

func CheckFavorite(c *gin.Context) {
	userID := c.GetUint("userID")
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var count int64
	models.DB.Model(&models.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"is_favorited": count > 0,
	})
}
