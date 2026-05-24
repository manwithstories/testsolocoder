package repository

import (
	"errors"

	"gorm.io/gorm"

	"furniture-platform/internal/model"
)

// DeliveryRepository 配送安装数据访问层
type DeliveryRepository struct {
	db *gorm.DB
}

// NewDeliveryRepository 创建配送安装数据访问层
func NewDeliveryRepository(db *gorm.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

// Create 创建配送安装预约
func (r *DeliveryRepository) Create(delivery *model.Delivery) error {
	return r.db.Create(delivery).Error
}

// GetByID 根据 ID 查询配送安装
func (r *DeliveryRepository) GetByID(id uint) (*model.Delivery, error) {
	var d model.Delivery
	err := r.db.First(&d, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

// Update 更新配送安装
func (r *DeliveryRepository) Update(delivery *model.Delivery) error {
	return r.db.Save(delivery).Error
}

// Delete 根据 ID 删除配送安装
func (r *DeliveryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Delivery{}, id).Error
}

// DeliveryListParams 配送安装列表查询参数
type DeliveryListParams struct {
	Page     int
	PageSize int
	OrderID  uint
	Type     string
	Status   string
	Role     string
	UserID   uint
}

// List 分页查询配送安装列表
func (r *DeliveryRepository) List(params *DeliveryListParams) ([]*model.Delivery, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	query := r.db.Model(&model.Delivery{})

	switch params.Role {
	case model.RoleOwner:
		query = query.Where("owner_id = ?", params.UserID)
	}

	if params.OrderID > 0 {
		query = query.Where("order_id = ?", params.OrderID)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*model.Delivery
	err := query.Order("id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
