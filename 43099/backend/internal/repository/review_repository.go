package repository

import (
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"
)

type ReviewRepository struct{}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{}
}

func (r *ReviewRepository) Create(review *model.Review) error {
	return database.DB.Create(review).Error
}

func (r *ReviewRepository) GetByID(id uint) (*model.Review, error) {
	var review model.Review
	err := database.DB.Preload("Order").Preload("User").First(&review, id).Error
	return &review, err
}

func (r *ReviewRepository) GetByOrderID(orderID uint) (*model.Review, error) {
	var review model.Review
	err := database.DB.Where("order_id = ?", orderID).First(&review).Error
	return &review, err
}

func (r *ReviewRepository) List(req *dto.ReviewListRequest) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	query := database.DB.Model(&model.Review{}).Preload("Order").Preload("User")

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) Update(review *model.Review) error {
	return database.DB.Save(review).Error
}
