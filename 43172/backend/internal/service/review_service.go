package service

import (
	"context"
	"errors"
	"fmt"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/repository"
	"luxury-trading-platform/internal/utils"

	"gorm.io/gorm"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	orderRepo  *repository.OrderRepository
	userRepo   *repository.UserRepository
	redisClient *cache.RedisClient
	db         *gorm.DB
}

func NewReviewService(reviewRepo *repository.ReviewRepository, orderRepo *repository.OrderRepository, userRepo *repository.UserRepository, redisClient *cache.RedisClient, db *gorm.DB) *ReviewService {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		redisClient: redisClient,
		db:          db,
	}
}

type CreateReviewRequest struct {
	OrderID     uint              `json:"order_id" binding:"required"`
	RevieweeID  uint              `json:"reviewee_id" binding:"required"`
	Rating      model.ReviewRating `json:"rating" binding:"required,min=1,max=5"`
	Content     string            `json:"content"`
	Images      string            `json:"images"`
	IsAnonymous bool              `json:"is_anonymous"`
}

func (s *ReviewService) CreateReview(ctx context.Context, reviewerID uint, req *CreateReviewRequest) (*model.Review, error) {
	order, err := s.orderRepo.FindByID(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.BuyerID != reviewerID && order.SellerID != reviewerID {
		return nil, errors.New("permission denied: you are not part of this order")
	}

	if order.Status != model.OrderStatusCompleted {
		return nil, errors.New("can only review completed orders")
	}

	if req.RevieweeID == reviewerID {
		return nil, errors.New("cannot review yourself")
	}

	if req.RevieweeID != order.BuyerID && req.RevieweeID != order.SellerID {
		return nil, errors.New("reviewee must be part of this order")
	}

	existing, err := s.reviewRepo.FindByOrderID(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if existing != nil {
		return nil, errors.New("review already exists for this order")
	}

	review := &model.Review{
		OrderID:     req.OrderID,
		ReviewerID:  reviewerID,
		RevieweeID:  req.RevieweeID,
		Rating:      req.Rating,
		Content:     req.Content,
		Images:      req.Images,
		IsAnonymous: req.IsAnonymous,
	}

	err = s.reviewRepo.Create(review)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	s.updateUserCreditScore(req.RevieweeID)

	return review, nil
}

func (s *ReviewService) GetReview(ctx context.Context, id uint) (*model.Review, error) {
	review, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find review: %w", err)
	}
	if review == nil {
		return nil, errors.New("review not found")
	}
	return review, nil
}

func (s *ReviewService) ListReviews(page, pageSize int, revieweeID *uint, reviewerID *uint, minRating *model.ReviewRating) ([]model.Review, int64, error) {
	page, pageSize = utils.ValidatePage(page, pageSize)
	return s.reviewRepo.List(page, pageSize, revieweeID, reviewerID, minRating)
}

func (s *ReviewService) GetUserAverageRating(userID uint) (float64, error) {
	return s.reviewRepo.GetAverageRating(userID)
}

func (s *ReviewService) updateUserCreditScore(userID uint) {
	avgRating, err := s.reviewRepo.GetAverageRating(userID)
	if err != nil {
		return
	}

	newScore := 50 + int(avgRating*10)
	if newScore > 100 {
		newScore = 100
	}
	if newScore < 0 {
		newScore = 0
	}

	_ = s.userRepo.UpdateCreditScore(userID, newScore)
}
