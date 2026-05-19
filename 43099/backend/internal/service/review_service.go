package service

import (
	"errors"
	"time"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/internal/repository"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	orderRepo  *repository.OrderRepository
}

func NewReviewService() *ReviewService {
	return &ReviewService{
		reviewRepo: repository.NewReviewRepository(),
		orderRepo:  repository.NewOrderRepository(),
	}
}

func (s *ReviewService) Create(req *dto.ReviewCreateRequest, userID uint) (*model.Review, error) {
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("you can only review your own orders")
	}

	if order.Status != model.OrderStatusCompleted {
		return nil, errors.New("only completed orders can be reviewed")
	}

	existing, _ := s.reviewRepo.GetByOrderID(req.OrderID)
	if existing != nil {
		return nil, errors.New("this order has already been reviewed")
	}

	review := &model.Review{
		OrderID: req.OrderID,
		UserID:  userID,
		Rating:  req.Rating,
		Content: req.Content,
		Status:  model.ReviewStatusPending,
	}

	err = s.reviewRepo.Create(review)
	return review, err
}

func (s *ReviewService) GetByID(id uint) (*model.Review, error) {
	return s.reviewRepo.GetByID(id)
}

func (s *ReviewService) List(req *dto.ReviewListRequest) ([]model.Review, int64, error) {
	return s.reviewRepo.List(req)
}

func (s *ReviewService) Approve(id uint, adminID uint, note string) error {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return errors.New("review not found")
	}

	if review.Status != model.ReviewStatusPending {
		return errors.New("only pending reviews can be approved")
	}

	now := time.Now()
	review.Status = model.ReviewStatusApproved
	review.ReviewNote = note
	review.ReviewedBy = &adminID
	review.ReviewedAt = &now

	return s.reviewRepo.Update(review)
}

func (s *ReviewService) Reject(id uint, adminID uint, note string) error {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return errors.New("review not found")
	}

	if review.Status != model.ReviewStatusPending {
		return errors.New("only pending reviews can be rejected")
	}

	now := time.Now()
	review.Status = model.ReviewStatusRejected
	review.ReviewNote = note
	review.ReviewedBy = &adminID
	review.ReviewedAt = &now

	return s.reviewRepo.Update(review)
}
