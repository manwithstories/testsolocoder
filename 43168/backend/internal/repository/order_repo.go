package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"furniture-platform/internal/model"
)

// OrderRepository 订单数据访问层
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单数据访问层
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateInquiry 创建询价单（订单+订单项，事务）
func (r *OrderRepository) CreateInquiry(order *model.Order, items []*model.OrderItem, history *model.OrderHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		if history != nil {
			history.OrderID = order.ID
			if err := tx.Create(history).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByID 根据 ID 查询订单
func (r *OrderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号查询
func (r *OrderRepository) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetItems 获取订单项
func (r *OrderRepository) GetItems(orderID uint) ([]*model.OrderItem, error) {
	var items []*model.OrderItem
	err := r.db.Where("order_id = ?", orderID).Order("id ASC").Find(&items).Error
	return items, err
}

// GetItemByID 根据 ID 获取订单项
func (r *OrderRepository) GetItemByID(id uint) (*model.OrderItem, error) {
	var item model.OrderItem
	err := r.db.First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetHistories 获取订单状态历史
func (r *OrderRepository) GetHistories(orderID uint) ([]*model.OrderHistory, error) {
	var list []*model.OrderHistory
	err := r.db.Where("order_id = ?", orderID).Order("id ASC").Find(&list).Error
	return list, err
}

// Update 更新订单
func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

// UpdateItemsQuote 更新订单项报价（事务内调用）
func (r *OrderRepository) UpdateItemsQuote(tx *gorm.DB, items []*model.OrderItem) error {
	for _, it := range items {
		if err := tx.Model(&model.OrderItem{}).
			Where("id = ?", it.ID).
			Updates(map[string]interface{}{
				"base_price": it.BasePrice,
				"subtotal":   it.Subtotal,
			}).Error; err != nil {
			return err
		}
	}
	return nil
}

// Quote 更新订单报价并记录历史（事务）
func (r *OrderRepository) Quote(order *model.Order, items []*model.OrderItem, history *model.OrderHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(order).Error; err != nil {
			return err
		}
		if err := r.UpdateItemsQuote(tx, items); err != nil {
			return err
		}
		if history != nil {
			if err := tx.Create(history).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateStatus 更新订单状态并记录历史（事务）
func (r *OrderRepository) UpdateStatus(order *model.Order, history *model.OrderHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Order{}).
			Where("id = ?", order.ID).
			Updates(map[string]interface{}{
				"status":     order.Status,
				"updated_at": order.UpdatedAt,
			}).Error; err != nil {
			return err
		}
		if history != nil {
			if err := tx.Create(history).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// LockForUpdate 锁定订单行（事务内）
func (r *OrderRepository) LockForUpdate(tx *gorm.DB, id uint) (*model.Order, error) {
	var order model.Order
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// OrderListParams 订单列表查询参数
type OrderListParams struct {
	Page     int
	PageSize int
	Keyword  string
	Status   string
	Role     string
	UserID   uint
}

// List 分页查询订单列表
func (r *OrderRepository) List(params *OrderListParams) ([]*model.Order, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	query := r.db.Model(&model.Order{})

	switch params.Role {
	case model.RoleOwner:
		query = query.Where("owner_id = ?", params.UserID)
	case model.RoleManufacturer:
		query = query.Where("manufacturer_id = ?", params.UserID)
	case model.RoleDesigner:
		query = query.Where("designer_id = ?", params.UserID)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where("order_no LIKE ? OR contact_name LIKE ? OR remark LIKE ?",
			like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*model.Order
	err := query.Order("id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
