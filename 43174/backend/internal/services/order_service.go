package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"campus-trade-platform/internal/models"
	"campus-trade-platform/internal/repository"
	"campus-trade-platform/internal/utils"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	textbookRepo    *repository.TextbookRepository
	userRepo          *repository.UserRepository
	db               *gorm.DB
}

func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	textbookRepo *repository.TextbookRepository,
	userRepo *repository.UserRepository,
	db *gorm.DB,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		textbookRepo:    textbookRepo,
		userRepo:          userRepo,
		db:               db,
	}
}

type CreateTransactionRequest struct {
	TextbookID   string  `json:"textbook_id" binding:"required"`
	SellerID     string  `json:"seller_id" binding:"required"`
	BuyerID      string  `json:"buyer_id" binding:"required"`
	Type         string  `json:"type" binding:"required,oneof=sell exchange"`
	AgreedPrice  float64 `json:"agreed_price"`
	ExchangeItem string  `json:"exchange_item"`
}

func (s *TransactionService) CreateTransaction(req *CreateTransactionRequest) (*models.Transaction, error) {
	textbookID, err := uuid.Parse(req.TextbookID)
	if err != nil {
		return nil, errors.New("invalid textbook ID")
	}

	sellerID, err := uuid.Parse(req.SellerID)
	if err != nil {
		return nil, errors.New("invalid seller ID")
	}

	buyerID, err := uuid.Parse(req.BuyerID)
	if err != nil {
		return nil, errors.New("invalid buyer ID")
	}

	textbook, err := s.textbookRepo.FindByID(req.TextbookID)
	if err != nil {
		return nil, errors.New("textbook not found")
	}

	if textbook.Status != models.TextbookStatusAvailable {
		return nil, errors.New("textbook is not available for transaction")
	}

	transaction := &models.Transaction{
		ID:             uuid.New(),
		TextbookID:     textbookID,
		SellerID:       sellerID,
		BuyerID:        buyerID,
		Type:           models.TransactionType(req.Type),
		AgreedPrice:    req.AgreedPrice,
		Status:         models.TransactionStatusPending,
		ExchangeItem:   req.ExchangeItem,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.transactionRepo.Create(transaction); err != nil {
			return err
		}

		if err := s.textbookRepo.UpdateStatus(req.TextbookID, models.TextbookStatusReserved); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactionByID(id string) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	return transaction, nil
}

func (s *TransactionService) GetAllTransactions(page, pageSize int, status, textbookID string) ([]models.Transaction, int64, error) {
	return s.transactionRepo.FindAll(page, pageSize, status, textbookID)
}

func (s *TransactionService) GetBuyerTransactions(buyerID string, page, pageSize int) ([]models.Transaction, int64, error) {
	return s.transactionRepo.FindByBuyerID(buyerID, page, pageSize)
}

func (s *TransactionService) GetSellerTransactions(sellerID string, page, pageSize int) ([]models.Transaction, int64, error) {
	return s.transactionRepo.FindBySellerID(sellerID, page, pageSize)
}

func (s *TransactionService) ConfirmTransaction(id string) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status != models.TransactionStatusPending &&
		transaction.Status != models.TransactionStatusNegotiating {
		return errors.New("transaction cannot be confirmed in current status")
	}

	return s.transactionRepo.UpdateStatus(id, models.TransactionStatusConfirmed)
}

func (s *TransactionService) CompleteTransaction(id string) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status != models.TransactionStatusConfirmed {
		return errors.New("transaction must be confirmed before completing")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.transactionRepo.UpdateStatus(id, models.TransactionStatusCompleted); err != nil {
			return err
		}

		if err := s.textbookRepo.UpdateStatus(transaction.TextbookID.String(), models.TextbookStatusSold); err != nil {
			return err
		}

		return nil
	})
}

func (s *TransactionService) CancelTransaction(id string, reason string) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status == models.TransactionStatusCompleted {
		return errors.New("completed transactions cannot be cancelled")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.transactionRepo.UpdateStatus(id, models.TransactionStatusCancelled); err != nil {
			return err
		}

		if err := s.textbookRepo.UpdateStatus(transaction.TextbookID.String(), models.TextbookStatusAvailable); err != nil {
			return err
		}

		return nil
	})
}

func (s *TransactionService) StartNegotiation(id string, price float64) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status != models.TransactionStatusPending {
		return errors.New("negotiation can only start from pending status")
	}

	transaction.AgreedPrice = price
	transaction.Status = models.TransactionStatusNegotiating
	return s.transactionRepo.Update(transaction)
}

type OrderService struct {
	orderRepo      *repository.OrderRepository
	textbookRepo   *repository.TextbookRepository
	userRepo        *repository.UserRepository
	db             *gorm.DB
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	textbookRepo *repository.TextbookRepository,
	userRepo *repository.UserRepository,
	db *gorm.DB,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		textbookRepo: textbookRepo,
		userRepo:    userRepo,
		db:         db,
	}
}

type OrderItemInput struct {
	TextbookID string  `json:"textbook_id" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required,min=1"`
	Price      float64 `json:"price" binding:"required"`
}

type CreateOrderRequest struct {
	BuyerID         string           `json:"buyer_id" binding:"required"`
	SellerID        string           `json:"seller_id" binding:"required"`
	Items           []OrderItemInput `json:"items" binding:"required"`
	ShippingAddress string         `json:"shipping_address"`
	PaymentMethod   string         `json:"payment_method"`
	TransactionID   string         `json:"transaction_id"`
}

func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*models.Order, error) {
	buyerID, err := uuid.Parse(req.BuyerID)
	if err != nil {
		return nil, errors.New("invalid buyer ID")
	}

	sellerID, err := uuid.Parse(req.SellerID)
	if err != nil {
		return nil, errors.New("invalid seller ID")
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range req.Items {
		textbook, err := s.textbookRepo.FindByID(item.TextbookID)
		if err != nil {
			return nil, fmt.Errorf("textbook %s not found", item.TextbookID)
		}

		if textbook.Status != models.TextbookStatusAvailable && textbook.Status != models.TextbookStatusReserved {
			return nil, fmt.Errorf("textbook %s is not available", textbook.Title)
		}

		subtotal := item.Price * float64(item.Quantity)
		totalAmount += subtotal

		tbID, _ := uuid.Parse(item.TextbookID)
		orderItems = append(orderItems, models.OrderItem{
			ID:         uuid.New(),
			TextbookID: tbID,
			Quantity:   item.Quantity,
			Price:      item.Price,
			Subtotal:   subtotal,
		})
	}

	var transactionID *uuid.UUID
	if req.TransactionID != "" {
		tid, err := uuid.Parse(req.TransactionID)
		if err == nil {
			transactionID = &tid
		}
	}

	order := &models.Order{
		ID:              uuid.New(),
		OrderNo:         utils.GenerateOrderNo(),
		BuyerID:         buyerID,
		SellerID:        sellerID,
		TotalAmount:     totalAmount,
		Status:          models.OrderStatusPending,
		PaymentMethod:   req.PaymentMethod,
		PaymentStatus:   "unpaid",
		TransactionID:   transactionID,
		ShippingAddress: req.ShippingAddress,
	}

	err = s.orderRepo.Create(order, orderItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) GetOrderByOrderNo(orderNo string) (*models.Order, error) {
	order, err := s.orderRepo.FindByOrderNo(orderNo)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) GetAllOrders(page, pageSize int, status, buyerID, sellerID string) ([]models.Order, int64, error) {
	return s.orderRepo.FindAll(page, pageSize, status, buyerID, sellerID)
}

func (s *OrderService) GetUserOrders(userID string, page, pageSize int) ([]models.Order, int64, error) {
	var allOrders []models.Order
	var total int64

	buyerOrders, buyerTotal, err := s.orderRepo.FindAll(page, pageSize, "", userID, "")
	if err != nil {
		return nil, 0, err
	}

	sellerOrders, sellerTotal, err := s.orderRepo.FindAll(page, pageSize, "", "", userID)
	if err != nil {
		return nil, 0, err
	}

	allOrders = append(allOrders, buyerOrders...)
	allOrders = append(allOrders, sellerOrders...)
	total = buyerTotal + sellerTotal

	return allOrders, total, nil
}

func (s *OrderService) UpdateOrderStatus(id string, status models.OrderStatus, remark string) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPending:   {models.OrderStatusPaid, models.OrderStatusCancelled},
		models.OrderStatusPaid:      {models.OrderStatusShipped, models.OrderStatusRefunded},
		models.OrderStatusShipped:   {models.OrderStatusDelivered, models.OrderStatusRefunded},
		models.OrderStatusDelivered: {models.OrderStatusCompleted},
	}

	allowedStatuses, ok := validTransitions[order.Status]
	if !ok {
		return errors.New("invalid order status transition")
	}

	isValid := false
	for _, s := range allowedStatuses {
		if s == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("cannot transition from %s to %s", order.Status, status)
	}

	return s.orderRepo.UpdateStatus(id, status, remark)
}

func (s *OrderService) PayOrder(id string) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("only pending orders can be paid")
	}

	order.PaymentStatus = "paid"
	order.Status = models.OrderStatusPaid

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepo.Update(order); err != nil {
			return err
		}

		history := &models.OrderStatusHistory{
			OrderID: order.ID,
			Status:  models.OrderStatusPaid,
			Remark:  "Payment completed successfully",
		}
		return tx.Create(history).Error
	})
}

func (s *OrderService) ShipOrder(id, trackingNumber string) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != models.OrderStatusPaid {
		return errors.New("only paid orders can be shipped")
	}

	order.TrackingNumber = trackingNumber
	order.Status = models.OrderStatusShipped

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepo.Update(order); err != nil {
			return err
		}

		history := &models.OrderStatusHistory{
			OrderID: order.ID,
			Status:  models.OrderStatusShipped,
			Remark:  fmt.Sprintf("Order shipped with tracking number: %s", trackingNumber),
		}
		return tx.Create(history).Error
	})
}

func (s *OrderService) DeliverOrder(id string) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != models.OrderStatusShipped {
		return errors.New("only shipped orders can be delivered")
	}

	return s.orderRepo.UpdateStatus(id, models.OrderStatusDelivered, "Order delivered successfully")
}

func (s *OrderService) CompleteOrder(id string) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != models.OrderStatusDelivered {
		return errors.New("only delivered orders can be completed")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepo.UpdateStatus(id, models.OrderStatusCompleted, "Order completed"); err != nil {
			return err
		}

		for _, item := range order.Items {
			if err := s.textbookRepo.UpdateStatus(item.TextbookID.String(), models.TextbookStatusSold); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *OrderService) CancelOrder(id, reason string) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status == models.OrderStatusCompleted || order.Status == models.OrderStatusDelivered {
		return errors.New("cannot cancel completed or delivered orders")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepo.UpdateStatus(id, models.OrderStatusCancelled, reason); err != nil {
			return err
		}

		for _, item := range order.Items {
			if err := s.textbookRepo.UpdateStatus(item.TextbookID.String(), models.TextbookStatusAvailable); err != nil {
				return err
			}
		}

		return nil
	})
}

type MessageService struct {
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
	db            *gorm.DB
}

func NewMessageService(messageRepo *repository.MessageRepository, userRepo *repository.UserRepository, db *gorm.DB) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		db:            db,
	}
}

type CreateMessageRequest struct {
	SenderID       string `json:"sender_id" binding:"required"`
	ReceiverID     string `json:"receiver_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
	RelatedOrderID string `json:"related_order_id"`
	IsDispute      bool   `json:"is_dispute"`
}

func (s *MessageService) CreateMessage(req *CreateMessageRequest) (*models.Message, error) {
	senderID, err := uuid.Parse(req.SenderID)
	if err != nil {
		return nil, errors.New("invalid sender ID")
	}

	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		return nil, errors.New("invalid receiver ID")
	}

	var relatedOrderID *uuid.UUID
	if req.RelatedOrderID != "" {
		oid, err := uuid.Parse(req.RelatedOrderID)
		if err == nil {
			relatedOrderID = &oid
		}
	}

	message := &models.Message{
		ID:             uuid.New(),
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Content:        req.Content,
		IsRead:         false,
		RelatedOrderID: relatedOrderID,
		IsDispute:      req.IsDispute,
	}

	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return message, nil
}

func (s *MessageService) GetConversation(userID1, userID2 string, page, pageSize int) ([]models.Message, int64, error) {
	return s.messageRepo.FindConversation(userID1, userID2, page, pageSize)
}

func (s *MessageService) GetUnreadCount(userID string) (int64, error) {
	return s.messageRepo.GetUnreadCount(userID)
}

func (s *MessageService) MarkAsRead(userID string) error {
	return s.messageRepo.MarkAsRead(userID)
}

type ReviewService struct {
	reviewRepo   *repository.ReviewRepository
	userRepo     *repository.UserRepository
	textbookRepo *repository.TextbookRepository
	noteRepo     *repository.NoteRepository
	db           *gorm.DB
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	userRepo *repository.UserRepository,
	textbookRepo *repository.TextbookRepository,
	noteRepo *repository.NoteRepository,
	db *gorm.DB,
) *ReviewService {
	return &ReviewService{
		reviewRepo:   reviewRepo,
		userRepo:     userRepo,
		textbookRepo: textbookRepo,
		noteRepo:     noteRepo,
		db:           db,
	}
}

type CreateReviewRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	TargetType string `json:"target_type" binding:"required,oneof=textbook note"`
	TextbookID string `json:"textbook_id"`
	NoteID     string `json:"note_id"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Content    string `json:"content"`
}

func (s *ReviewService) CreateReview(req *CreateReviewRequest) (*models.Review, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var textbookID *uuid.UUID
	var noteID *uuid.UUID

	if req.TargetType == "textbook" {
		if req.TextbookID == "" {
			return nil, errors.New("textbook_id is required for textbook reviews")
		}
		tid, err := uuid.Parse(req.TextbookID)
		if err != nil {
			return nil, errors.New("invalid textbook ID")
		}
		textbookID = &tid

		_, err = s.textbookRepo.FindByID(req.TextbookID)
		if err != nil {
			return nil, errors.New("textbook not found")
		}
	} else if req.TargetType == "note" {
		if req.NoteID == "" {
			return nil, errors.New("note_id is required for note reviews")
		}
		nid, err := uuid.Parse(req.NoteID)
		if err != nil {
			return nil, errors.New("invalid note ID")
		}
		noteID = &nid

		_, err = s.noteRepo.FindByID(req.NoteID)
		if err != nil {
			return nil, errors.New("note not found")
		}
	}

	review := &models.Review{
		ID:         uuid.New(),
		UserID:     userID,
		TargetType: models.ReviewTargetType(req.TargetType),
		TextbookID: textbookID,
		NoteID:     noteID,
		Rating:     req.Rating,
		Content:    req.Content,
		IsHidden:  false,
		IsMalicious: false,
	}

	err = s.reviewRepo.Create(review)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (s *ReviewService) GetTextbookReviews(textbookID string, page, pageSize int) ([]models.Review, int64, error) {
	return s.reviewRepo.FindByTextbookID(textbookID, page, pageSize)
}

func (s *ReviewService) GetNoteReviews(noteID string, page, pageSize int) ([]models.Review, int64, error) {
	return s.reviewRepo.FindByNoteID(noteID, page, pageSize)
}

func (s *ReviewService) GetAllReviews(page, pageSize int, isMalicious bool) ([]models.Review, int64, error) {
	return s.reviewRepo.FindAll(page, pageSize, isMalicious)
}

func (s *ReviewService) HideReview(id string) error {
	_, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return errors.New("review not found")
	}
	return s.reviewRepo.HideReview(id)
}

func (s *ReviewService) MarkMalicious(id string, isMalicious bool) error {
	_, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return errors.New("review not found")
	}
	return s.reviewRepo.MarkMalicious(id, isMalicious)
}

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	db            *gorm.DB
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, db *gorm.DB) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		db:           db,
	}
}

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID string `json:"parent_id"`
	SortOrder int   `json:"sort_order"`
}

func (s *CategoryService) CreateCategory(req *CreateCategoryRequest) (*models.Category, error) {
	var parentID *uuid.UUID
	if req.ParentID != "" {
		pid, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent ID")
		}
		parentID = &pid
	}

	category := &models.Category{
		ID:        uuid.New(),
		Name:      req.Name,
		ParentID:  parentID,
		SortOrder: req.SortOrder,
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *CategoryService) GetCategoryByID(id string) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *CategoryService) UpdateCategory(id string, updates map[string]interface{}) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	if name, ok := updates["name"]; ok {
		category.Name = name.(string)
	}
	if sortOrder, ok := updates["sort_order"]; ok {
		category.SortOrder = sortOrder.(int)
	}

	category.UpdatedAt = time.Now()
	return s.categoryRepo.Update(category)
}

func (s *CategoryService) DeleteCategory(id string) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.categoryRepo.Delete(id)
}

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) CreateNotification(userID uuid.UUID, notificationType, title, content string) error {
	notification := &models.Notification{
		ID:      uuid.New(),
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Content: content,
		IsRead:  false,
	}
	return s.db.Create(notification).Error
}

func (s *NotificationService) GetUserNotifications(userID string, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	offset := (page - 1) * pageSize

	s.db.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total)
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error

	return notifications, total, err
}

func (s *NotificationService) MarkAsRead(notificationID string) error {
	return s.db.Model(&models.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func (s *NotificationService) MarkAllAsRead(userID string) error {
	return s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (s *NotificationService) GetUnreadCount(userID string) (int64, error) {
	var count int64
	err := s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
