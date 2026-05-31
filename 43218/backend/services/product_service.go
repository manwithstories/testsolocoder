package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"secondhand-platform/cache"
	"secondhand-platform/database"
	"secondhand-platform/models"
	"secondhand-platform/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) CreateProduct(sellerID uint, title, category, brand, model, condition string, price, originalPrice float64, description string, warrantyDays int, images string) (*models.Product, error) {
	var seller models.User
	if err := database.DB.First(&seller, sellerID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if seller.Role != models.RoleSeller {
		return nil, errors.New("只有卖家可以发布商品")
	}

	if !seller.IsAuthenticated {
		return nil, errors.New("卖家需要先通过实名认证")
	}

	product := &models.Product{
		SellerID:      sellerID,
		Title:         title,
		Category:      category,
		Brand:         brand,
		Model:         model,
		Condition:     condition,
		Price:         price,
		OriginalPrice: originalPrice,
		Description:   description,
		WarrantyDays:  warrantyDays,
		Images:        images,
		Status:        models.ProductStatusPending,
	}

	result := database.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)
	cached, err := cache.Get(nil, cacheKey)
	if err == nil && cached != "" {
		var product models.Product
		if jsonErr := json.Unmarshal([]byte(cached), &product); jsonErr == nil {
			return &product, nil
		}
	}

	var product models.Product
	if err := database.DB.Preload("Seller").First(&product, id).Error; err != nil {
		return nil, err
	}

	product.ViewCount++
	database.DB.Model(&product).Update("view_count", product.ViewCount)

	productJSON, _ := json.Marshal(product)
	cache.Set(nil, cacheKey, productJSON, 5*time.Minute)

	return &product, nil
}

func (s *ProductService) ListProducts(page, pageSize int, category, condition string, minPrice, maxPrice float64, keyword, sortBy string, sellerID uint) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	db := database.DB.Model(&models.Product{}).Where("status = ?", models.ProductStatusOnSale)

	if category != "" {
		db = db.Where("category = ?", category)
	}
	if condition != "" {
		db = db.Where("condition = ?", condition)
	}
	if minPrice > 0 {
		db = db.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		db = db.Where("price <= ?", maxPrice)
	}
	if keyword != "" {
		db = db.Where("title LIKE ? OR description LIKE ? OR brand LIKE ? OR model LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sellerID > 0 {
		db = db.Where("seller_id = ?", sellerID)
	}

	switch sortBy {
	case "price_asc":
		db = db.Order("price ASC")
	case "price_desc":
		db = db.Order("price DESC")
	case "newest":
		db = db.Order("created_at DESC")
	case "sold":
		db = db.Order("sold_count DESC")
	default:
		db = db.Order("created_at DESC")
	}

	db.Count(&total)
	if err := db.Preload("Seller").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) UpdateProduct(userID, productID uint, updates map[string]interface{}) error {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return errors.New("商品不存在")
	}

	if product.SellerID != userID {
		return errors.New("无权修改此商品")
	}

	if product.Status == models.ProductStatusSoldOut {
		return errors.New("已售出的商品无法修改")
	}

	return database.DB.Model(&product).Updates(updates).Error
}

func (s *ProductService) DeleteProduct(userID, productID uint) error {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return errors.New("商品不存在")
	}

	if product.SellerID != userID {
		return errors.New("无权删除此商品")
	}

	cache.Delete(nil, fmt.Sprintf("product:%d", productID))
	return database.DB.Delete(&product).Error
}

func (s *ProductService) ReviewProduct(productID uint, approved bool, rejectReason string, reviewerID uint) error {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return errors.New("商品不存在")
	}

	now := time.Now()
	product.ReviewedAt = &now
	product.ReviewedBy = &reviewerID

	if approved {
		product.Status = models.ProductStatusOnSale
	} else {
		product.Status = models.ProductStatusRejected
		product.RejectReason = rejectReason
	}

	cache.Delete(nil, fmt.Sprintf("product:%d", productID))
	return database.DB.Save(&product).Error
}

func (s *ProductService) OffShelfProduct(userID, productID uint) error {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return errors.New("商品不存在")
	}

	if product.SellerID != userID {
		return errors.New("无权操作此商品")
	}

	product.Status = models.ProductStatusOffShelf
	cache.Delete(nil, fmt.Sprintf("product:%d", productID))
	return database.DB.Save(&product).Error
}

func (s *ProductService) ListPendingProducts(page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	db := database.DB.Model(&models.Product{}).Where("status = ?", models.ProductStatusPending)
	db.Count(&total)
	if err := db.Preload("Seller").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) ListMyProducts(userID uint, page, pageSize int, status int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	db := database.DB.Model(&models.Product{}).Where("seller_id = ?", userID)
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) ToggleFavorite(userID, productID uint) (bool, error) {
	var favorite models.Favorite
	err := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&favorite).Error

	if err == nil {
		database.DB.Delete(&favorite)
		database.DB.Model(&models.Product{}).Where("id = ?", productID).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1"))
		return false, nil
	}

	favorite = models.Favorite{
		UserID:    userID,
		ProductID: productID,
	}
	database.DB.Create(&favorite)
	database.DB.Model(&models.Product{}).Where("id = ?", productID).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1"))
	return true, nil
}

func (s *ProductService) ListFavorites(userID uint, page, pageSize int) ([]models.Favorite, int64, error) {
	var favorites []models.Favorite
	var total int64

	db := database.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)
	db.Count(&total)
	if err := db.Preload("Product").Preload("Product.Seller").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

func (s *ProductService) GetCategories() []string {
	return models.ProductCategories
}

func (s *ProductService) GetConditions() []string {
	return models.ProductConditions
}

func (s *ProductService) SaveProductImages(productID uint, images []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("product_id = ?", productID).Delete(&models.ProductImage{})

		for i, img := range images {
			productImage := models.ProductImage{
				ProductID: productID,
				ImageURL:  img,
				SortOrder: i,
			}
			if err := tx.Create(&productImage).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ProductService) GetHotProducts(limit int) ([]models.Product, error) {
	var products []models.Product
	if err := database.DB.Where("status = ?", models.ProductStatusOnSale).
		Order("view_count DESC").
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func init() {
	logrus.Info("Product service initialized")
}
