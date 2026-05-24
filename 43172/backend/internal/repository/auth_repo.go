package repository

import (
	"errors"
	"time"

	"luxury-trading-platform/internal/model"

	"gorm.io/gorm"
)

type AuthenticationRepository struct {
	db *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) *AuthenticationRepository {
	return &AuthenticationRepository{db: db}
}

func (r *AuthenticationRepository) Create(auth *model.Authentication) error {
	return r.db.Create(auth).Error
}

func (r *AuthenticationRepository) FindByID(id uint) (*model.Authentication, error) {
	var auth model.Authentication
	err := r.db.Preload("Order").
		Preload("Product").
		Preload("Buyer").
		Preload("Authenticator").
		First(&auth, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &auth, nil
}

func (r *AuthenticationRepository) FindByOrderID(orderID uint) (*model.Authentication, error) {
	var auth model.Authentication
	err := r.db.Where("order_id = ?", orderID).
		Preload("Order").
		Preload("Product").
		Preload("Buyer").
		Preload("Authenticator").
		First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &auth, nil
}

func (r *AuthenticationRepository) Update(auth *model.Authentication) error {
	return r.db.Save(auth).Error
}

func (r *AuthenticationRepository) List(page, pageSize int, status model.AuthenticationStatus, authenticatorID *uint, buyerID *uint) ([]model.Authentication, int64, error) {
	var auths []model.Authentication
	var total int64

	query := r.db.Model(&model.Authentication{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if authenticatorID != nil {
		query = query.Where("authenticator_id = ?", *authenticatorID)
	}
	if buyerID != nil {
		query = query.Where("buyer_id = ?", *buyerID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Order").
		Preload("Product").
		Preload("Buyer").
		Preload("Authenticator").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&auths).Error

	return auths, total, err
}

func (r *AuthenticationRepository) Accept(id uint, authenticatorID uint) error {
	now := time.Now()
	return r.db.Model(&model.Authentication{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           model.AuthenticationStatusAccepted,
			"authenticator_id": authenticatorID,
			"accepted_at":      now,
		}).Error
}

func (r *AuthenticationRepository) Complete(id uint, result model.AuthenticationResult, reportFile, reportContent, notes string) error {
	now := time.Now()
	return r.db.Model(&model.Authentication{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":             model.AuthenticationStatusCompleted,
			"result":             result,
			"report_file":        reportFile,
			"report_content":     reportContent,
			"authenticator_notes": notes,
			"completed_at":       now,
		}).Error
}

func (r *AuthenticationRepository) Reject(id uint, reason string) error {
	return r.db.Model(&model.Authentication{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           model.AuthenticationStatusRejected,
			"rejection_reason": reason,
		}).Error
}

func (r *AuthenticationRepository) Cancel(id uint) error {
	return r.db.Model(&model.Authentication{}).
		Where("id = ?", id).
		Update("status", model.AuthenticationStatusCancelled).Error
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) FindByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Preload("Order").
		Preload("Reviewer").
		Preload("Reviewee").
		First(&review, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) FindByOrderID(orderID uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("order_id = ?", orderID).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) List(page, pageSize int, revieweeID *uint, reviewerID *uint, minRating *model.ReviewRating) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	query := r.db.Model(&model.Review{})
	if revieweeID != nil {
		query = query.Where("reviewee_id = ?", *revieweeID)
	}
	if reviewerID != nil {
		query = query.Where("reviewer_id = ?", *reviewerID)
	}
	if minRating != nil {
		query = query.Where("rating >= ?", *minRating)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Order").
		Preload("Reviewer").
		Preload("Reviewee").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

func (r *ReviewRepository) GetAverageRating(revieweeID uint) (float64, error) {
	var result struct {
		Avg float64
	}
	err := r.db.Model(&model.Review{}).
		Where("reviewee_id = ?", revieweeID).
		Select("COALESCE(AVG(rating), 0) as avg").
		Scan(&result).Error
	return result.Avg, err
}

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) List(page, pageSize int, userID *uint, module string, action string) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.Model(&model.AuditLog{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

type StatisticRepository struct {
	db *gorm.DB
}

func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{db: db}
}

func (r *StatisticRepository) GetTransactionTrend(startDate, endDate time.Time) ([]model.Statistic, error) {
	var stats []model.Statistic
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date ASC").
		Find(&stats).Error
	return stats, err
}

func (r *StatisticRepository) GetTotalOrders() (int64, error) {
	var count int64
	err := r.db.Model(&model.Order{}).Count(&count).Error
	return count, err
}

func (r *StatisticRepository) GetTotalAmount() (float64, error) {
	var result struct {
		Total float64
	}
	err := r.db.Model(&model.Order{}).
		Where("status = ?", model.OrderStatusCompleted).
		Select("COALESCE(SUM(price), 0) as total").
		Scan(&result).Error
	return result.Total, err
}

func (r *StatisticRepository) GetTotalUsers() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *StatisticRepository) GetTotalProducts() (int64, error) {
	var count int64
	err := r.db.Model(&model.Product{}).
		Where("status = ?", model.ProductStatusOnSale).
		Count(&count).Error
	return count, err
}

func (r *StatisticRepository) GetBrandRankings() ([]struct {
	BrandName   string
	ProductCount int64
	OrderCount   int64
}, error) {
	var results []struct {
		BrandName   string
		ProductCount int64
		OrderCount   int64
	}
	err := r.db.Table("products p").
		Select("p.brand_name, COUNT(DISTINCT p.id) as product_count, COUNT(DISTINCT o.id) as order_count").
		Joins("LEFT JOIN orders o ON o.product_id = p.id").
		Where("p.status = ? AND p.brand_name != ?", model.ProductStatusOnSale, "").
		Group("p.brand_name").
		Order("order_count DESC").
		Limit(10).
		Scan(&results).Error
	return results, err
}

func (r *StatisticRepository) GetAuthStatistics() (struct {
	Total     int64
	Completed int64
	Passed    int64
	PassRate  float64
}, error) {
	var result struct {
		Total     int64
		Completed int64
		Passed    int64
		PassRate  float64
	}
	r.db.Model(&model.Authentication{}).Count(&result.Total)
	r.db.Model(&model.Authentication{}).
		Where("status = ?", model.AuthenticationStatusCompleted).
		Count(&result.Completed)
	r.db.Model(&model.Authentication{}).
		Where("result = ?", model.AuthenticationResultGenuine).
		Count(&result.Passed)
	if result.Completed > 0 {
		result.PassRate = float64(result.Passed) / float64(result.Completed) * 100
	}
	return result, nil
}
