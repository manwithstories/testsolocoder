package repository

import (
	"car-rental/internal/model"
	cachedb "car-rental/internal/config"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{db: cachedb.DB}
}

func (r *ReviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) FindByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Preload("User").Preload("Car").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) FindByCarID(carID uint, page, pageSize int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	query := r.db.Model(&model.Review{}).Where("car_id = ? AND is_hidden = ?", carID, false)
	query.Count(&total)
	err := query.Preload("User").Preload("Car").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) FindByUserID(userID uint, page, pageSize int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	query := r.db.Model(&model.Review{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Preload("User").Preload("Car").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) FindAll(page, pageSize int, carID uint, minRating int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	query := r.db.Model(&model.Review{}).Where("is_hidden = ?", false)
	if carID > 0 {
		query = query.Where("car_id = ?", carID)
	}
	if minRating > 0 {
		query = query.Where("rating >= ?", minRating)
	}

	query.Count(&total)
	err := query.Preload("User").Preload("Car").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) Update(review *model.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.Review{}, id).Error
}

func (r *ReviewRepository) ExistsByBookingID(bookingID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Review{}).Where("booking_id = ?", bookingID).Count(&count).Error
	return count > 0, err
}

func (r *ReviewRepository) IncrementLikes(id uint) error {
	return r.db.Model(&model.Review{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + 1")).Error
}
