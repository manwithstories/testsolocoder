package services

import (
	"errors"
	"time"

	"auction-system/internal/dto"
	"auction-system/internal/models"
	"auction-system/pkg/logger"
)

type ReviewService struct{}

func NewReviewService() *ReviewService {
	return &ReviewService{}
}

func (s *ReviewService) CreateReview(reviewerID uint, req *dto.CreateReviewRequest) (*models.Review, error) {
	var order models.Order
	if err := models.DB.First(&order, req.OrderID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.Status != OrderStatusCompleted {
		return nil, errors.New("订单未完成，无法评价")
	}

	var existingReview models.Review
	if models.DB.Where("order_id = ? AND reviewer_id = ?", req.OrderID, reviewerID).First(&existingReview).Error == nil {
		return nil, errors.New("您已评价过此订单")
	}

	var revieweeID uint
	if order.BuyerID == reviewerID {
		revieweeID = order.SellerID
	} else if order.SellerID == reviewerID {
		revieweeID = order.BuyerID
	} else {
		return nil, errors.New("无权评价此订单")
	}

	review := &models.Review{
		OrderID:    req.OrderID,
		ReviewerID: reviewerID,
		RevieweeID: revieweeID,
		Rating:     req.Rating,
		Content:    req.Content,
		CreatedAt:  time.Now(),
	}

	if err := models.DB.Create(review).Error; err != nil {
		logger.Error("Failed to create review: %v", err)
		return nil, errors.New("评价失败")
	}

	return review, nil
}

func (s *ReviewService) GetOrderReviews(orderID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := models.DB.Where("order_id = ?", orderID).Preload("Reviewer").Preload("Reviewee").Find(&reviews).Error
	return reviews, err
}

func (s *ReviewService) GetUserReviews(userID uint, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := models.DB.Model(&models.Review{}).Where("reviewee_id = ?", userID).Preload("Reviewer").Preload("Order.AuctionItem")
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews).Error
	return reviews, total, err
}

func (s *ReviewService) GetUserAverageRating(userID uint) (float64, int64, error) {
	var avgRating float64
	var total int64

	models.DB.Model(&models.Review{}).Where("reviewee_id = ?", userID).Count(&total)
	if total == 0 {
		return 5.0, 0, nil
	}

	row := models.DB.Model(&models.Review{}).Where("reviewee_id = ?", userID).Select("AVG(rating)").Row()
	row.Scan(&avgRating)

	return avgRating, total, nil
}

func (s *ReviewService) GetReviewByID(reviewID uint) (*models.Review, error) {
	var review models.Review
	if err := models.DB.Preload("Reviewer").Preload("Reviewee").Preload("Order.AuctionItem").First(&review, reviewID).Error; err != nil {
		return nil, err
	}
	return &review, nil
}
