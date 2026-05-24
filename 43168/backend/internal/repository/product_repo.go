package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"furniture-platform/internal/model"
)

// ProductRepository 产品数据访问层
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品数据访问层
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create 创建产品（包含选项和图片）
func (r *ProductRepository) Create(product *model.Product, options []*model.ProductOption, images []*model.ProductImage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		for i := range options {
			options[i].ProductID = product.ID
		}
		if len(options) > 0 {
			if err := tx.Create(&options).Error; err != nil {
				return err
			}
		}
		for i := range images {
			images[i].ProductID = product.ID
		}
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByID 根据 ID 查询产品
func (r *ProductRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// Update 更新产品基本信息
func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// ReplaceOptions 替换产品选项（删除旧的，批量插入新的）
func (r *ProductRepository) ReplaceOptions(productID uint, options []*model.ProductOption) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&model.ProductOption{}).Error; err != nil {
			return err
		}
		for i := range options {
			options[i].ID = 0
			options[i].ProductID = productID
		}
		if len(options) > 0 {
			return tx.Create(&options).Error
		}
		return nil
	})
}

// ReplaceImages 替换产品图片
func (r *ProductRepository) ReplaceImages(productID uint, images []*model.ProductImage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&model.ProductImage{}).Error; err != nil {
			return err
		}
		for i := range images {
			images[i].ID = 0
			images[i].ProductID = productID
		}
		if len(images) > 0 {
			return tx.Create(&images).Error
		}
		return nil
	})
}

// Delete 根据 ID 删除产品（级联删除选项和图片）
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&model.ProductOption{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&model.ProductImage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Product{}, id).Error
	})
}

// GetOptions 获取产品选项
func (r *ProductRepository) GetOptions(productID uint) ([]*model.ProductOption, error) {
	var options []*model.ProductOption
	err := r.db.Where("product_id = ?", productID).Order("sort ASC, id ASC").Find(&options).Error
	return options, err
}

// GetImages 获取产品图片
func (r *ProductRepository) GetImages(productID uint) ([]*model.ProductImage, error) {
	var images []*model.ProductImage
	err := r.db.Where("product_id = ?", productID).Order("sort ASC, id ASC").Find(&images).Error
	return images, err
}

// ProductListParams 产品列表查询参数
type ProductListParams struct {
	Page     int
	PageSize int
	Keyword  string
	Category string
	Status   *int
	IsHot    *bool
}

// List 分页查询产品列表
func (r *ProductRepository) List(params *ProductListParams) ([]*model.Product, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	query := r.db.Model(&model.Product{})

	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ? OR category LIKE ? OR description LIKE ?",
			like, like, like)
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.IsHot != nil {
		query = query.Where("is_hot = ?", *params.IsHot)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []*model.Product
	err := query.Order("id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// ListHot 获取热门产品列表（上架状态）
func (r *ProductRepository) ListHot(limit int) ([]*model.Product, error) {
	if limit <= 0 {
		limit = 20
	}
	var products []*model.Product
	err := r.db.Model(&model.Product{}).
		Where("is_hot = ? AND status = ?", true, model.ProductStatusOnSale).
		Order("id DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// GetCoverImage 获取产品封面图（第一张）
func (r *ProductRepository) GetCoverImage(productID uint) (string, error) {
	var img model.ProductImage
	err := r.db.Where("product_id = ?", productID).
		Order("sort ASC, id ASC").
		First(&img).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return img.ImageURL, nil
}

// LockForUpdate 锁定产品行（用于事务内扣库存等）
func (r *ProductRepository) LockForUpdate(tx *gorm.DB, id uint) (*model.Product, error) {
	var product model.Product
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}
