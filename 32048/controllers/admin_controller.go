package controllers

import (
	"net/http"
	"secondhand-trading/models"
	"secondhand-trading/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewProductRequest struct {
	Status  string `json:"status" binding:"required,oneof=approved rejected"`
	Reason  string `json:"reason"`
}

func GetPendingReviewProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := models.DB.Model(&models.Product{}).Where("needs_review = ? AND is_reviewed = ?", true, false)

	var total int64
	query.Count(&total)

	var products []models.Product
	offset := (page - 1) * pageSize
	query.Preload("Images", "is_primary = ?", true).
		Preload("Seller", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, avatar, credit_score")
		}).
		Preload("Category").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"products":   products,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func ReviewProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req ReviewProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if !product.NeedsReview || product.IsReviewed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not pending review"})
		return
	}

	err = models.DB.Transaction(func(tx *gorm.DB) error {
		if req.Status == "approved" {
			if err := tx.Model(&product).Updates(map[string]interface{}{
				"is_reviewed":  true,
				"needs_review": false,
			}).Error; err != nil {
				return err
			}

			utils.CreateNotification(product.SellerID, "product_approved", "商品审核通过",
				"您的商品已通过审核，现已上架", &product.ID)
		} else {
			if err := tx.Model(&product).Updates(map[string]interface{}{
				"is_reviewed":  true,
				"needs_review": false,
				"status":       "off_shelf",
			}).Error; err != nil {
				return err
			}

			reason := ""
			if req.Reason != "" {
				reason = "，原因：" + req.Reason
			}
			utils.CreateNotification(product.SellerID, "product_rejected", "商品审核未通过",
				"您的商品未通过审核"+reason, &product.ID)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product reviewed successfully",
		"product": product,
	})
}

func BatchReviewProducts(c *gin.Context) {
	var req struct {
		ProductIDs []uint `json:"product_ids" binding:"required,min=1"`
		Status     string `json:"status" binding:"required,oneof=approved rejected"`
		Reason     string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := models.DB.Model(&models.Product{}).
		Where("id IN ? AND needs_review = ? AND is_reviewed = ?",
			req.ProductIDs, true, false)

	var products []models.Product
	result.Find(&products)

	if len(products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid products to review"})
		return
	}

	updates := map[string]interface{}{
		"is_reviewed":  true,
		"needs_review": false,
	}

	if req.Status == "rejected" {
		updates["status"] = "off_shelf"
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).
			Where("id IN ?", req.ProductIDs).
			Updates(updates).Error; err != nil {
			return err
		}

		reason := ""
		if req.Reason != "" {
			reason = "，原因：" + req.Reason
		}

		action := "已通过审核，现已上架"
		notifyType := "product_approved"
		if req.Status == "rejected" {
			action = "未通过审核" + reason
			notifyType = "product_rejected"
		}

		for _, p := range products {
			utils.CreateNotification(p.SellerID, notifyType, "商品审核结果",
				"您的商品"+action, &p.ID)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to batch review products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Batch review completed",
		"reviewed_count": len(products),
	})
}
