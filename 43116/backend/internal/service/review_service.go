package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"errors"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	carRepo    *repository.CarRepository
	bookingRepo *repository.BookingRepository
}

func NewReviewService() *ReviewService {
	return &ReviewService{
		reviewRepo:  repository.NewReviewRepository(),
		carRepo:     repository.NewCarRepository(),
		bookingRepo: repository.NewBookingRepository(),
	}
}

type CreateReviewRequest struct {
	BookingID   uint   `json:"booking_id" binding:"required"`
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Content     string `json:"content"`
	Images      string `json:"images"`
	IsAnonymous bool   `json:"is_anonymous"`
}

func (s *ReviewService) CreateReview(userID uint, req *CreateReviewRequest) (*model.Review, error) {
	booking, err := s.bookingRepo.FindByID(req.BookingID)
	if err != nil {
		return nil, errors.New("预订不存在")
	}

	if booking.UserID != userID {
		return nil, errors.New("无权对该预订进行评价")
	}

	if booking.Status != model.BookingStatusCompleted {
		return nil, errors.New("只能对已完成的订单进行评价")
	}

	exists, _ := s.reviewRepo.ExistsByBookingID(req.BookingID)
	if exists {
		return nil, errors.New("您已评价过该订单")
	}

	review := &model.Review{
		UserID:      userID,
		CarID:       booking.CarID,
		BookingID:   req.BookingID,
		Rating:      req.Rating,
		Content:     req.Content,
		Images:      req.Images,
		IsAnonymous: req.IsAnonymous,
	}

	err = s.reviewRepo.Create(review)
	if err != nil {
		return nil, err
	}

	err = s.carRepo.UpdateRating(booking.CarID)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) GetReviewByID(id uint) (*model.Review, error) {
	return s.reviewRepo.FindByID(id)
}

func (s *ReviewService) GetCarReviews(carID uint, page, pageSize int) ([]model.Review, int64, error) {
	return s.reviewRepo.FindByCarID(carID, page, pageSize)
}

func (s *ReviewService) GetUserReviews(userID uint, page, pageSize int) ([]model.Review, int64, error) {
	return s.reviewRepo.FindByUserID(userID, page, pageSize)
}

func (s *ReviewService) GetAllReviews(page, pageSize int, carID uint, minRating int) ([]model.Review, int64, error) {
	return s.reviewRepo.FindAll(page, pageSize, carID, minRating)
}

func (s *ReviewService) UpdateReview(id uint, userID uint, updates map[string]interface{}) error {
	review, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return errors.New("评价不存在")
	}

	if review.UserID != userID {
		return errors.New("无权修改该评价")
	}

	if rating, ok := updates["rating"]; ok {
		review.Rating = int(rating.(float64))
	}
	if content, ok := updates["content"]; ok {
		review.Content = content.(string)
	}
	if images, ok := updates["images"]; ok {
		review.Images = images.(string)
	}

	err = s.reviewRepo.Update(review)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateRating(review.CarID)
}

func (s *ReviewService) DeleteReview(id uint, userID uint) error {
	review, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return errors.New("评价不存在")
	}

	if review.UserID != userID {
		return errors.New("无权删除该评价")
	}

	err = s.reviewRepo.Delete(id)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateRating(review.CarID)
}

func (s *ReviewService) ToggleReviewHidden(id uint, isHidden bool) error {
	review, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return errors.New("评价不存在")
	}

	review.IsHidden = isHidden
	err = s.reviewRepo.Update(review)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateRating(review.CarID)
}

func (s *ReviewService) LikeReview(id uint) error {
	return s.reviewRepo.IncrementLikes(id)
}