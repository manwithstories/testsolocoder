package repository

import (
	"context"
	"errors"
	"time"

	"print3d-platform/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("DesignerProfile").
		Preload("PrinterProfile").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("DesignerProfile").
		Preload("PrinterProfile").
		First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("DesignerProfile").
		Preload("PrinterProfile").
		First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) UpdateLoginInfo(ctx context.Context, id uuid.UUID, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
		}).Error
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *UserRepository) CreateDesignerProfile(ctx context.Context, profile *models.DesignerProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *UserRepository) CreatePrinterProfile(ctx context.Context, profile *models.PrinterProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *UserRepository) UpdateDesignerProfile(ctx context.Context, profile *models.DesignerProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *UserRepository) UpdatePrinterProfile(ctx context.Context, profile *models.PrinterProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *UserRepository) UpdateDesignerRating(ctx context.Context, designerID uuid.UUID, rating float64) error {
	return r.db.WithContext(ctx).Model(&models.DesignerProfile{}).
		Where("user_id = ?", designerID).
		Updates(map[string]interface{}{
			"rating":       rating,
			"rating_count": gorm.Expr("rating_count + 1"),
		}).Error
}

func (r *UserRepository) UpdatePrinterRating(ctx context.Context, printerID uuid.UUID, rating float64) error {
	return r.db.WithContext(ctx).Model(&models.PrinterProfile{}).
		Where("user_id = ?", printerID).
		Updates(map[string]interface{}{
			"rating":       rating,
			"rating_count": gorm.Expr("rating_count + 1"),
		}).Error
}

func (r *UserRepository) ListDesigners(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).
		Preload("DesignerProfile").
		Where("role = ? AND status = ?", models.RoleDesigner, models.UserStatusActive)

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) ListPrinters(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).
		Preload("PrinterProfile").
		Where("role = ? AND status = ?", models.RolePrinter, models.UserStatusActive)

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) AddNotification(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *UserRepository) GetNotifications(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&notifications).Error

	return notifications, total, err
}

func (r *UserRepository) MarkNotificationRead(ctx context.Context, id, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *UserRepository) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(transaction).Error
	}
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *UserRepository) GetTransactions(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&transactions).Error

	return transactions, total, err
}

func (r *UserRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var modelCount int64
	r.db.WithContext(ctx).Model(&models.Model3D{}).Where("designer_id = ?", userID).Count(&modelCount)
	stats["model_count"] = modelCount

	var orderCount int64
	r.db.WithContext(ctx).Model(&models.PrintOrder{}).Where("customer_id = ? OR printer_id = ?", userID, userID).Count(&orderCount)
	stats["order_count"] = orderCount

	var totalSpent float64
	var totalEarned float64
	r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "expense").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalSpent)
	r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "income").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalEarned)
	stats["total_spent"] = totalSpent
	stats["total_earned"] = totalEarned

	return stats, nil
}

func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) ExecTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *UserRepository) LockForUpdate(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
