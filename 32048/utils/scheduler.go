package utils

import (
	"log"
	"secondhand-trading/models"
	"time"

	"gorm.io/gorm"
)

func StartAutoCloseTransactionJob() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		autoCloseExpiredTransactions()
	}
}

func autoCloseExpiredTransactions() {
	var transactions []models.Transaction
	now := time.Now()

	result := models.DB.Where("status IN ? AND expires_at < ?",
		[]string{"pending", "waiting_confirmation"}, now).Find(&transactions)

	if result.Error != nil {
		log.Printf("Error finding expired transactions: %v", result.Error)
		return
	}

	for _, tx := range transactions {
		if err := closeTransaction(&tx); err != nil {
			log.Printf("Error closing transaction %d: %v", tx.ID, err)
		}
	}

	if len(transactions) > 0 {
		log.Printf("Closed %d expired transactions", len(transactions))
	}
}

func closeTransaction(tx *models.Transaction) error {
	return models.DB.Transaction(func(txDB *gorm.DB) error {
		if err := txDB.Model(tx).Update("status", "cancelled").Error; err != nil {
			return err
		}

		if err := txDB.Model(&models.Product{}).Where("id = ?", tx.ProductID).Update("status", "on_sale").Error; err != nil {
			return err
		}

		CreateNotification(tx.BuyerID, "transaction_cancelled", "交易已关闭",
			"您的交易已超时自动关闭，商品已重新上架", &tx.ProductID)
		CreateNotification(tx.SellerID, "transaction_cancelled", "交易已关闭",
			"您的交易已超时自动关闭，商品已重新上架", &tx.ProductID)

		return nil
	})
}

func CreateNotification(userID uint, notifyType, title, content string, productID *uint) {
	notification := models.Notification{
		UserID:    userID,
		Type:      notifyType,
		Title:     title,
		Content:   content,
		ProductID: productID,
	}
	models.DB.Create(&notification)
}

func CalculateCreditScoreChange(rating int) int {
	switch rating {
	case 5:
		return 5
	case 4:
		return 2
	case 3:
		return 0
	case 2:
		return -3
	case 1:
		return -5
	default:
		return 0
	}
}
