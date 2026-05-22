package repository

import (
	"beauty-salon-system/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) List(page, pageSize int, category string, lowStock bool) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Where("status = 1")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if lowStock {
		query = query.Where("stock <= threshold")
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) ListAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("status = 1").Find(&products).Error
	return products, err
}

func (r *ProductRepository) DeductStock(id uint, quantity int, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&model.Product{}).Where("id = ? AND stock >= ?", id, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *ProductRepository) AddStock(id uint, quantity int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *ProductRepository) GetLowStockProducts() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("status = 1 AND stock <= threshold").Find(&products).Error
	return products, err
}

type ProductRecordRepository struct {
	db *gorm.DB
}

func NewProductRecordRepository(db *gorm.DB) *ProductRecordRepository {
	return &ProductRecordRepository{db: db}
}

func (r *ProductRecordRepository) Create(record *model.ProductRecord, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(record).Error
	}
	return r.db.Create(record).Error
}

func (r *ProductRecordRepository) List(page, pageSize int, productID uint, changeType string) ([]model.ProductRecord, int64, error) {
	var records []model.ProductRecord
	var total int64

	query := r.db.Model(&model.ProductRecord{}).Preload("Product")
	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}
	if changeType != "" {
		query = query.Where("change_type = ?", changeType)
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&records).Error
	return records, total, err
}

type ProductSaleRepository struct {
	db *gorm.DB
}

func NewProductSaleRepository(db *gorm.DB) *ProductSaleRepository {
	return &ProductSaleRepository{db: db}
}

func (r *ProductSaleRepository) Create(sale *model.ProductSale, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(sale).Error
	}
	return r.db.Create(sale).Error
}

func (r *ProductSaleRepository) List(page, pageSize int, customerID uint) ([]model.ProductSale, int64, error) {
	var sales []model.ProductSale
	var total int64

	query := r.db.Model(&model.ProductSale{}).Preload("Product").Preload("Customer.User")
	if customerID > 0 {
		query = query.Where("customer_id = ?", customerID)
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&sales).Error
	return sales, total, err
}
