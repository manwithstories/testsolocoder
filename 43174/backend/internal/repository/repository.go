package repository

import (
	"gorm.io/gorm"

	"campus-trade-platform/internal/models"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *TransactionRepository) FindByID(id string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Textbook").Preload("Seller").Preload("Buyer").
		First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) FindAll(page, pageSize int, status string, textbookID string) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Preload("Textbook").Preload("Seller").Preload("Buyer")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if textbookID != "" {
		query = query.Where("textbook_id = ?", textbookID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepository) FindByBuyerID(buyerID string, page, pageSize int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("buyer_id = ?", buyerID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepository) FindBySellerID(sellerID string, page, pageSize int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("seller_id = ?", sellerID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepository) UpdateStatus(id string, status models.TransactionStatus) error {
	return r.db.Model(&models.Transaction{}).Where("id = ?", id).Update("status", status).Error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order, items []models.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		history := &models.OrderStatusHistory{
			OrderID: order.ID,
			Status:  order.Status,
			Remark:  "Order created",
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Buyer").Preload("Seller").Preload("Items.Textbook").
		Preload("StatusHistory").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Buyer").Preload("Seller").Preload("Items.Textbook").
		Preload("StatusHistory").First(&order, "order_no = ?", orderNo).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindAll(page, pageSize int, status string, buyerID string, sellerID string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).Preload("Buyer").Preload("Seller").Preload("Items.Textbook")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if buyerID != "" {
		query = query.Where("buyer_id = ?", buyerID)
	}
	if sellerID != "" {
		query = query.Where("seller_id = ?", sellerID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(id string, status models.OrderStatus, remark string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			return err
		}

		history := &models.OrderStatusHistory{
			OrderID: mustParseUUID(id),
			Status:  status,
			Remark:  remark,
		}
		return tx.Create(history).Error
	})
}

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) FindByID(id string) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("Sender").Preload("Receiver").First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) FindConversation(userID1, userID2 string, page, pageSize int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	query := r.db.Model(&models.Message{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1)

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *MessageRepository) GetUnreadCount(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).Where("receiver_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *MessageRepository) MarkAsRead(userID string) error {
	return r.db.Model(&models.Message{}).Where("receiver_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) FindByID(id string) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) FindByTextbookID(textbookID string, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.Model(&models.Review{}).Where("textbook_id = ? AND is_hidden = ?", textbookID, false)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) FindByNoteID(noteID string, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.Model(&models.Review{}).Where("note_id = ? AND is_hidden = ?", noteID, false)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) FindAll(page, pageSize int, isMalicious bool) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.Model(&models.Review{}).Preload("User")
	if isMalicious {
		query = query.Where("is_malicious = ?", true)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) HideReview(id string) error {
	return r.db.Model(&models.Review{}).Where("id = ?", id).Update("is_hidden", true).Error
}

func (r *ReviewRepository) MarkMalicious(id string, isMalicious bool) error {
	return r.db.Model(&models.Review{}).Where("id = ?", id).Update("is_malicious", isMalicious).Error
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id string) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

func (r *CategoryRepository) FindByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Children").First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Children").Where("parent_id IS NULL").Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func mustParseUUID(s string) [16]byte {
	var uuid [16]byte
	copy(uuid[:], []byte(s))
	return uuid
}
