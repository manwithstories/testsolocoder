package controllers

import (
	"net/http"
	"secondhand-trading/config"
	"secondhand-trading/models"
	"secondhand-trading/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTransactionRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type ConfirmTransactionRequest struct {
	Role string `json:"role" binding:"required,oneof=buyer seller"`
}

func CreateTransaction(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := models.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.Status != "on_sale" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not available for sale"})
		return
	}

	if product.SellerID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot buy your own product"})
		return
	}

	var existingTransaction models.Transaction
	if models.DB.Where("product_id = ? AND status IN ?", req.ProductID,
		[]string{"pending", "waiting_confirmation", "confirmed"}).First(&existingTransaction).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is already in a transaction"})
		return
	}

	transaction := models.Transaction{
		ProductID: req.ProductID,
		BuyerID:   userID,
		SellerID:  product.SellerID,
		Price:     product.Price,
		Status:    "pending",
		ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.AutoCloseHours) * time.Hour),
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		if err := tx.Model(&product).Update("status", "reserved").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	utils.CreateNotification(product.SellerID, "new_order", "新订单通知",
		"您有一个新的订单，请及时处理", &req.ProductID)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

func GetTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")

	var transaction models.Transaction
	if err := models.DB.Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Images", "is_primary = ?", true)
	}).Preload("Buyer", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar, credit_score")
	}).Preload("Seller", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar, credit_score")
	}).First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.BuyerID != userID && transaction.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view this transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func ConfirmTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")

	var transaction models.Transaction
	if err := models.DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.BuyerID != userID && transaction.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to confirm this transaction"})
		return
	}

	if transaction.Status == "cancelled" || transaction.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is already closed"})
		return
	}

	isBuyer := transaction.BuyerID == userID

	if isBuyer && transaction.BuyerConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already confirmed this transaction"})
		return
	}
	if !isBuyer && transaction.SellerConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already confirmed this transaction"})
		return
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if isBuyer {
			if err := tx.Model(&transaction).Update("buyer_confirmed", true).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&transaction).Update("seller_confirmed", true).Error; err != nil {
				return err
			}
		}

		var updatedTx models.Transaction
		tx.First(&updatedTx, id)

		if updatedTx.BuyerConfirmed && updatedTx.SellerConfirmed {
			if err := tx.Model(&updatedTx).Update("status", "completed").Error; err != nil {
				return err
			}

			if err := tx.Model(&models.Product{}).Where("id = ?", transaction.ProductID).
				Update("status", "sold").Error; err != nil {
				return err
			}

			otherUserID := transaction.BuyerID
			if isBuyer {
				otherUserID = transaction.SellerID
			}

			utils.CreateNotification(otherUserID, "transaction_completed", "交易完成",
				"交易已完成，请对对方进行评价", &transaction.ProductID)
		} else {
			if transaction.Status == "pending" {
				if err := tx.Model(&transaction).Update("status", "waiting_confirmation").Error; err != nil {
					return err
				}

				otherUserID := transaction.SellerID
				if isBuyer {
					otherUserID = transaction.SellerID
				} else {
					otherUserID = transaction.BuyerID
				}
				role := "买家"
				if !isBuyer {
					role = "卖家"
				}
				utils.CreateNotification(otherUserID, "transaction_waiting", "等待确认",
					role+"已确认交易，请您确认", &transaction.ProductID)
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction confirmed successfully"})
}

func CancelTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")

	var transaction models.Transaction
	if err := models.DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.BuyerID != userID && transaction.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to cancel this transaction"})
		return
	}

	if transaction.Status == "cancelled" || transaction.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is already closed"})
		return
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&transaction).Update("status", "cancelled").Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Product{}).Where("id = ?", transaction.ProductID).
			Update("status", "on_sale").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel transaction"})
		return
	}

	otherUserID := transaction.SellerID
	if transaction.BuyerID == userID {
		otherUserID = transaction.SellerID
	} else {
		otherUserID = transaction.BuyerID
	}
	utils.CreateNotification(otherUserID, "transaction_cancelled", "交易已取消",
		"交易已被对方取消", &transaction.ProductID)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction cancelled successfully"})
}

func GetMyTransactions(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.Query("role")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := models.DB.Model(&models.Transaction{})

	switch role {
	case "buyer":
		query = query.Where("buyer_id = ?", userID)
	case "seller":
		query = query.Where("seller_id = ?", userID)
	default:
		query = query.Where("buyer_id = ? OR seller_id = ?", userID, userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var transactions []models.Transaction
	offset := (page - 1) * pageSize
	query.Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Images", "is_primary = ?", true)
	}).Preload("Buyer", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar")
	}).Preload("Seller", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar")
	}).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"total_page":   (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
