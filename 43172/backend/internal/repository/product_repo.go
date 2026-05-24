package repository

import (
	"errors"

	"luxury-trading-platform/internal/model"

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

func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Seller").
		Preload("Brand").
		Preload("Images").
		First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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

func (r *ProductRepository) List(page, pageSize int, category model.ProductCategory, brandID *uint, status model.ProductStatus, sortBy string, keyword string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if brandID != nil {
		query = query.Where("brand_id = ?", *brandID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR brand_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	orderClause := "created_at DESC"
	switch sortBy {
	case "price_asc":
		orderClause = "price ASC"
	case "price_desc":
		orderClause = "price DESC"
	case "views":
		orderClause = "views DESC"
	case "newest":
		orderClause = "created_at DESC"
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Seller").
		Preload("Brand").
		Preload("Images").
		Offset(offset).
		Limit(pageSize).
		Order(orderClause).
		Find(&products).Error

	return products, total, err
}

func (r *ProductRepository) ListBySeller(sellerID uint, page, pageSize int, status model.ProductStatus) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Where("seller_id = ?", sellerID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Images").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *ProductRepository) UpdateStatus(id uint, status model.ProductStatus) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ProductRepository) UpdateStock(id uint, stock int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("stock", stock).Error
}

func (r *ProductRepository) DecrementStock(id uint, quantity int) error {
	return r.db.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", id, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *ProductRepository) IncrementStock(id uint, quantity int) error {
	return r.db.Model(&model.Product{}).
		Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *ProductRepository) IncrementViews(id uint) error {
	return r.db.Model(&model.Product{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (r *ProductRepository) SetAuthenticated(id uint, authenticated bool) error {
	return r.db.Model(&model.Product{}).
		Where("id = ?", id).
		Update("is_authenticated", authenticated).Error
}

func (r *ProductRepository) CreateImage(image *model.ProductImage) error {
	return r.db.Create(image).Error
}

func (r *ProductRepository) DeleteImages(productID uint) error {
	return r.db.Where("product_id = ?", productID).Delete(&model.ProductImage{}).Error
}

func (r *ProductRepository) GetImages(productID uint) ([]model.ProductImage, error) {
	var images []model.ProductImage
	err := r.db.Where("product_id = ?", productID).Order("sort_order ASC").Find(&images).Error
	return images, err
}

func (r *ProductRepository) CreateBrand(brand *model.Brand) error {
	return r.db.Create(brand).Error
}

func (r *ProductRepository) FindBrandByID(id uint) (*model.Brand, error) {
	var brand model.Brand
	err := r.db.First(&brand, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &brand, nil
}

func (r *ProductRepository) FindBrandByName(name string) (*model.Brand, error) {
	var brand model.Brand
	err := r.db.Where("name = ?", name).First(&brand).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &brand, nil
}

func (r *ProductRepository) ListBrands(category model.ProductCategory) ([]model.Brand, error) {
	var brands []model.Brand
	query := r.db.Model(&model.Brand{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Order("popularity DESC").Find(&brands).Error
	return brands, err
}

func (r *ProductRepository) IncrementBrandPopularity(brandID uint) error {
	return r.db.Model(&model.Brand{}).
		Where("id = ?", brandID).
		UpdateColumn("popularity", gorm.Expr("popularity + 1")).Error
}
