package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"skillshare/internal/models"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) GetMessages(senderID, receiverID uuid.UUID, page, pageSize int) ([]*models.Message, int64, error) {
	var messages []*models.Message
	var total int64

	query := r.db.Model(&models.Message{}).Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID,
	)

	query.Count(&total)
	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&messages)

	return messages, total, nil
}

func (r *MessageRepository) MarkAsRead(senderID, receiverID uuid.UUID) error {
	return r.db.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", senderID, receiverID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *MessageRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *MessageRepository) GetConversations(userID uuid.UUID) ([]*models.Message, error) {
	var messages []*models.Message

	err := r.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? OR receiver_id = ?) AND id IN (?)",
			userID, userID,
			r.db.Model(&models.Message{}).
				Select("MAX(id)").
				Where("sender_id = ? OR receiver_id = ?", userID, userID).
				Group(fmt.Sprintf("CASE WHEN sender_id = '%s' THEN receiver_id ELSE sender_id END", userID)),
		).
		Order("created_at DESC").
		Find(&messages).Error

	return messages, err
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Booking").Preload("Payer").Preload("Payee").
		First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) GetUserPayments(userID uuid.UUID, page, pageSize int) ([]*models.Payment, int64, error) {
	var payments []*models.Payment
	var total int64

	query := r.db.Model(&models.Payment{}).Where("payer_id = ? OR payee_id = ?", userID, userID)

	query.Count(&total)
	query.Preload("Booking").Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&payments)

	return payments, total, nil
}

func (r *PaymentRepository) ReleaseEscrowPayments() error {
	return r.db.Model(&models.Payment{}).
		Where("type = ? AND status = ? AND escrow_release_at <= ? AND released = ?",
			models.TransactionTypePayment,
			models.TransactionStatusCompleted,
			time.Now(),
			false,
		).Updates(map[string]interface{}{
			"released":   true,
			"released_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *PaymentRepository) GetUserWallet(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *PaymentRepository) CreateWallet(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *PaymentRepository) UpdateWallet(userID uuid.UUID, amount float64) error {
	return r.db.Model(&models.Wallet{}).Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *ScheduleRepository) Update(schedule *models.Schedule) error {
	return r.db.Save(schedule).Error
}

func (r *ScheduleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Schedule{}, id).Error
}

func (r *ScheduleRepository) GetUserSchedules(userID uuid.UUID) ([]*models.Schedule, error) {
	var schedules []*models.Schedule
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("day_of_week ASC, start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) FindByID(id uuid.UUID) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.First(&schedule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) GetUserAvailability(userID uuid.UUID, dayOfWeek models.DayOfWeek) ([]*models.Schedule, error) {
	var schedules []*models.Schedule
	err := r.db.Where("user_id = ? AND day_of_week = ? AND type = ? AND is_active = ?",
		userID, dayOfWeek, models.ScheduleTypeAvailability, true).
		Order("start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

type OperationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{db: db}
}

func (r *OperationLogRepository) Create(log *models.OperationLog) error {
	return r.db.Create(log).Error
}

func (r *OperationLogRepository) GetLogs(page, pageSize int, userID *uuid.UUID, operation *string) ([]*models.OperationLog, int64, error) {
	var logs []*models.OperationLog
	var total int64

	query := r.db.Model(&models.OperationLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if operation != nil {
		query = query.Where("operation = ?", *operation)
	}

	query.Count(&total)
	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	return logs, total, nil
}
