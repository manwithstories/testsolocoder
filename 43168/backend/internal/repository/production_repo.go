package repository

import (
	"context"
	"fmt"
	"furniture-platform/internal/model"
	"gorm.io/gorm"
)

type ProductionRepo struct {
	db *gorm.DB
}

func NewProductionRepo(db *gorm.DB) *ProductionRepo {
	return &ProductionRepo{db: db}
}

func (r *ProductionRepo) Create(ctx context.Context, p *model.Production) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductionRepo) GetByID(ctx context.Context, id int64) (*model.Production, error) {
	var p model.Production
	err := r.db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductionRepo) GetByOrderID(ctx context.Context, orderID int64) (*model.Production, error) {
	var p model.Production
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductionRepo) Update(ctx context.Context, p *model.Production) error {
	return r.db.WithContext(ctx).Save(p).Error
}

type ProductionListFilter struct {
	Page     int
	PageSize int
	OrderID  int64
	Status   string
}

func (r *ProductionRepo) List(ctx context.Context, filter ProductionListFilter) ([]model.Production, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Production{})
	if filter.OrderID > 0 {
		db = db.Where("order_id = ?", filter.OrderID)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Production
	offset := (filter.Page - 1) * filter.PageSize
	if err := db.Order("id DESC").Limit(filter.PageSize).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	fmt.Println("list found:", len(list))
	return list, total, nil
}
