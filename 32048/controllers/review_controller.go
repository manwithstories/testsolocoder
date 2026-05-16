package controllers

import (
	"net/http"
	"secondhand-trading/models"
	"secondhand-trading/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateReviewRequest struct {
	TransactionID uint   `json:"transaction_id" binding:"required"`
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	Comment       string `json:"comment" binding:"max=1000"`
}

func CreateReview(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	if err := models.DB.First(&transaction, req.TransactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is not completed yet"})
		return
	}

	var existingReview models.Review
	if models.DB.Where("transaction_id = ? AND reviewer_id = ?", req.TransactionID, userID).First(&existingReview).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this transaction"})
		return
	}

	isBuyer := transaction.BuyerID == userID
	isSeller := transaction.SellerID == userID

	if !isBuyer && !isSeller {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not part of this transaction"})
		return
	}

	if isBuyer && transaction.BuyerReviewed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this transaction"})
		return
	}
	if isSeller && transaction.SellerReviewed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this transaction"})
		return
	}

	targetID := transaction.SellerID
	if isBuyer {
		targetID = transaction.SellerID
	} else {
		targetID = transaction.BuyerID
	}

	review := models.Review{
		TransactionID: req.TransactionID,
		ReviewerID:    userID,
		TargetID:      targetID,
		Rating:        req.Rating,
		Comment:       req.Comment,
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		if isBuyer {
			if err := tx.Model(&transaction).Update("buyer_reviewed", true).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&transaction).Update("seller_reviewed", true).Error; err != nil {
				return err
			}
		}

		scoreChange := utils.CalculateCreditScoreChange(req.Rating)
		if err := tx.Model(&models.User{}).Where("id = ?", targetID).
			UpdateColumn("credit_score", gorm.Expr("credit_score + ?", scoreChange)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

func GetTransactionReviews(c *gin.Context) {
	transactionID := c.Param("transaction_id")

	var reviews []models.Review
	if err := models.DB.Preload("Reviewer", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar")
	}).Where("transaction_id = ?", transactionID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func GetUserReviews(c *gin.Context) {
	userID := c.Param("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := models.DB.Model(&models.Review{}).Where("target_id = ?", userID)

	var total int64
	query.Count(&total)

	var reviews []models.Review
	offset := (page - 1) * pageSize
	query.Preload("Reviewer", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar")
	}).Preload("Transaction", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, product_id").Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title")
		})
	}).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews)

	var avgRating float64
	models.DB.Model(&models.Review{}).Where("target_id = ?", userID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)

	c.JSON(http.StatusOK, gin.H{
		"reviews":    reviews,
		"total":      total,
		"avg_rating": avgRating,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func GetMyReceivedReviews(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := models.DB.Model(&models.Review{}).Where("target_id = ?", userID)

	var total int64
	query.Count(&total)

	var reviews []models.Review
	offset := (page - 1) * pageSize
	query.Preload("Reviewer", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar")
	}).Preload("Transaction", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, product_id").Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title")
		})
	}).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews)

	var avgRating float64
	models.DB.Model(&models.Review{}).Where("target_id = ?", userID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)

	c.JSON(http.StatusOK, gin.H{
		"reviews":    reviews,
		"total":      total,
		"avg_rating": avgRating,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
