package service

import (
	"errors"

	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
)

// ReviewService 售后评价业务逻辑层
type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	orderRepo  *repository.OrderRepository
	db         *gorm.DB
}

// NewReviewService 创建售后评价服务
func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	orderRepo *repository.OrderRepository,
	db *gorm.DB,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
		db:         db,
	}
}

// CreateReview 业主创建评价
func (s *ReviewService) CreateReview(ownerID uint, req *dto.CreateReviewRequest) (*model.Review, error) {
	if !model.ValidRating(req.ProductRating) || !model.ValidRating(req.ServiceRating) {
		return nil, errors.New("评分必须在 1-5 之间")
	}

	// 校验订单存在且属于当前业主
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}
	if order.OwnerID != ownerID {
		return nil, errors.New("无权评价他人的订单")
	}

	// 校验订单是否包含该产品
	items, err := s.orderRepo.GetItems(order.ID)
	if err != nil {
		return nil, err
	}
	found := false
	for _, it := range items {
		if it.ProductID == req.ProductID {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("订单中不存在该产品")
	}

	// 校验是否已评价
	exist, err := s.reviewRepo.GetByOrderAndProduct(req.OrderID, req.ProductID)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, errors.New("已对该产品进行过评价")
	}

	review := &model.Review{
		OrderID:       req.OrderID,
		ProductID:     req.ProductID,
		OwnerID:       ownerID,
		ProductRating: req.ProductRating,
		ServiceRating: req.ServiceRating,
		Content:       req.Content,
	}
	if err := review.SetImages(req.Images); err != nil {
		return nil, err
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}
	return review, nil
}

// GetByID 根据 ID 获取评价
func (s *ReviewService) GetByID(id uint, userID uint, role string) (*model.Review, error) {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if review == nil {
		return nil, errors.New("评价不存在")
	}

	if role == model.RoleOwner && review.OwnerID != userID {
		return nil, errors.New("无权查看他人的评价")
	}
	return review, nil
}

// List 分页查询评价列表
func (s *ReviewService) List(params *dto.ReviewListRequest) ([]*model.Review, int64, error) {
	ownerID := params.OwnerID
	if params.Role == model.RoleOwner {
		ownerID = params.UserID
	}
	p := &repository.ReviewListParams{
		Page:      params.Page,
		PageSize:  params.PageSize,
		OrderID:   params.OrderID,
		ProductID: params.ProductID,
		OwnerID:   ownerID,
	}
	return s.reviewRepo.List(p)
}
