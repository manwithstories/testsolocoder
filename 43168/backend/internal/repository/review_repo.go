package repository

import (
	"errors"

	"gorm.io/gorm"

	"furniture-platform/internal/model"
)

// ReviewRepository 售后评价数据访问层
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository 创建售后评价数据访问层
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create 创建评价
func (r *ReviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

// GetByID 根据 ID 查询评价
func (r *ReviewRepository) GetByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.First(&review, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

// Update 更新评价
func (r *ReviewRepository) Update(review *model.Review) error {
	return r.db.Save(review).Error
}

// Delete 根据 ID 删除评价
func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.Review{}, id).Error
}

// GetByOrderAndProduct 根据订单+产品查询已有评价
func (r *ReviewRepository) GetByOrderAndProduct(orderID, productID uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("order_id = ? AND product_id = ?", orderID, productID).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

// ReviewListParams 评价列表查询参数
type ReviewListParams struct {
	Page      int
	PageSize  int
	OrderID   uint
	ProductID uint
	OwnerID   uint
}

// List 分页查询评价列表
func (r *ReviewRepository) List(params *ReviewListParams) ([]*model.Review, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	query := r.db.Model(&model.Review{})

	if params.OrderID > 0 {
		query = query.Where("order_id = ?", params.OrderID)
	}
	if params.ProductID > 0 {
		query = query.Where("product_id = ?", params.ProductID)
	}
	if params.OwnerID > 0 {
		query = query.Where("owner_id = ?", params.OwnerID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*model.Review
	err := query.Order("id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
