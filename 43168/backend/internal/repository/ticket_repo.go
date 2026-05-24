package repository

import (
	"errors"

	"gorm.io/gorm"

	"furniture-platform/internal/model"
)

// TicketRepository 售后工单数据访问层
type TicketRepository struct {
	db *gorm.DB
}

// NewTicketRepository 创建售后工单数据访问层
func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Create 创建工单
func (r *TicketRepository) Create(ticket *model.Ticket) error {
	return r.db.Create(ticket).Error
}

// GetByID 根据 ID 查询工单
func (r *TicketRepository) GetByID(id uint) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.First(&ticket, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ticket, nil
}

// Update 更新工单
func (r *TicketRepository) Update(ticket *model.Ticket) error {
	return r.db.Save(ticket).Error
}

// Delete 根据 ID 删除工单
func (r *TicketRepository) Delete(id uint) error {
	return r.db.Delete(&model.Ticket{}, id).Error
}

// TicketListParams 工单列表查询参数
type TicketListParams struct {
	Page     int
	PageSize int
	OrderID  uint
	Type     string
	Status   string
	Role     string
	UserID   uint
}

// List 分页查询工单列表
func (r *TicketRepository) List(params *TicketListParams) ([]*model.Ticket, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	query := r.db.Model(&model.Ticket{})

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

	var list []*model.Ticket
	err := query.Order("id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
