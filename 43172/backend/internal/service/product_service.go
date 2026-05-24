package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/repository"
	"luxury-trading-platform/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepo *repository.ProductRepository
	redisClient *cache.RedisClient
	db          *gorm.DB
	log         *logrus.Logger
	uploadPath  string
}

func NewProductService(productRepo *repository.ProductRepository, redisClient *cache.RedisClient, db *gorm.DB, log *logrus.Logger, uploadPath string) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		redisClient: redisClient,
		db:          db,
		log:         log,
		uploadPath:  uploadPath,
	}
}

type CreateProductRequest struct {
	Title         string              `json:"title" binding:"required,min=5,max=200"`
	Description   string              `json:"description" binding:"required,min=10"`
	Category      model.ProductCategory `json:"category" binding:"required,oneof=bag jewelry watch clothing shoes other"`
	BrandID       *uint               `json:"brand_id"`
	BrandName     string              `json:"brand_name"`
	OriginalPrice float64             `json:"original_price"`
	Price         float64             `json:"price" binding:"required,gt=0"`
	Condition     string              `json:"condition"`
	Color         string              `json:"color"`
	Size          string              `json:"size"`
	Material      string              `json:"material"`
	Stock         int                 `json:"stock" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Title         *string             `json:"title"`
	Description   *string             `json:"description"`
	Category      *model.ProductCategory `json:"category"`
	BrandID       *uint               `json:"brand_id"`
	BrandName     *string             `json:"brand_name"`
	OriginalPrice *float64            `json:"original_price"`
	Price         *float64            `json:"price"`
	Condition     *string             `json:"condition"`
	Color         *string             `json:"color"`
	Size          *string             `json:"size"`
	Material      *string             `json:"material"`
	Stock         *int                `json:"stock"`
}

func (s *ProductService) CreateProduct(ctx context.Context, sellerID uint, req *CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		SellerID:      sellerID,
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		BrandID:       req.BrandID,
		BrandName:     req.BrandName,
		OriginalPrice: req.OriginalPrice,
		Price:         req.Price,
		Condition:     req.Condition,
		Color:         req.Color,
		Size:          req.Size,
		Material:      req.Material,
		Stock:         req.Stock,
		Status:        model.ProductStatusDraft,
	}

	err := s.productRepo.Create(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uint) (*model.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	if s.redisClient != nil {
		cached, err := s.redisClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			_ = s.productRepo.IncrementViews(id)
			return s.productRepo.FindByID(id)
		}
	}

	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	_ = s.productRepo.IncrementViews(id)

	if s.redisClient != nil {
		_ = s.redisClient.Set(ctx, cacheKey, product, 1*time.Hour)
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uint, sellerID uint, req *UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if product.SellerID != sellerID {
		return nil, errors.New("permission denied: you can only update your own products")
	}

	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.BrandID != nil {
		product.BrandID = req.BrandID
	}
	if req.BrandName != nil {
		product.BrandName = *req.BrandName
	}
	if req.OriginalPrice != nil {
		product.OriginalPrice = *req.OriginalPrice
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Condition != nil {
		product.Condition = *req.Condition
	}
	if req.Color != nil {
		product.Color = *req.Color
	}
	if req.Size != nil {
		product.Size = *req.Size
	}
	if req.Material != nil {
		product.Material = *req.Material
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	if s.redisClient != nil {
		_ = s.redisClient.Del(ctx, cacheKey)
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint, sellerID uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return errors.New("product not found")
	}

	if product.SellerID != sellerID {
		return errors.New("permission denied: you can only delete your own products")
	}

	if product.Status == model.ProductStatusSold {
		return errors.New("cannot delete a sold product")
	}

	err = s.productRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	if s.redisClient != nil {
		_ = s.redisClient.Del(ctx, cacheKey)
	}

	return nil
}

func (s *ProductService) ListProducts(page, pageSize int, category model.ProductCategory, brandID *uint, status model.ProductStatus, sortBy string, keyword string) ([]model.Product, int64, error) {
	page, pageSize = utils.ValidatePage(page, pageSize)
	return s.productRepo.List(page, pageSize, category, brandID, status, sortBy, keyword)
}

func (s *ProductService) ListSellerProducts(sellerID uint, page, pageSize int, status model.ProductStatus) ([]model.Product, int64, error) {
	page, pageSize = utils.ValidatePage(page, pageSize)
	return s.productRepo.ListBySeller(sellerID, page, pageSize, status)
}

func (s *ProductService) UpdateProductStatus(ctx context.Context, id uint, sellerID uint, status model.ProductStatus) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return errors.New("product not found")
	}

	if product.SellerID != sellerID {
		return errors.New("permission denied")
	}

	err = s.productRepo.UpdateStatus(id, status)
	if err != nil {
		return fmt.Errorf("failed to update product status: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	if s.redisClient != nil {
		_ = s.redisClient.Del(ctx, cacheKey)
	}

	return nil
}

func (s *ProductService) UploadImages(ctx context.Context, productID uint, sellerID uint, files []*multipart.FileHeader) ([]model.ProductImage, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if product.SellerID != sellerID {
		return nil, errors.New("permission denied")
	}

	uploadDir := filepath.Join(s.uploadPath, fmt.Sprintf("products/%d", productID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	var images []model.ProductImage
	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := saveUploadedFile(file, filePath); err != nil {
			s.log.Errorf("Failed to save file: %v", err)
			continue
		}

		image := model.ProductImage{
			ProductID: productID,
			ImageURL:  fmt.Sprintf("/uploads/products/%d/%s", productID, filename),
			ImageType: ext,
			IsPrimary: i == 0,
			SortOrder: i,
		}

		if err := s.productRepo.CreateImage(&image); err != nil {
			s.log.Errorf("Failed to create image record: %v", err)
			continue
		}

		images = append(images, image)
	}

	return images, nil
}

func (s *ProductService) GetProductImages(productID uint) ([]model.ProductImage, error) {
	return s.productRepo.GetImages(productID)
}

func (s *ProductService) DeleteProductImages(ctx context.Context, productID uint, sellerID uint) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return errors.New("product not found")
	}

	if product.SellerID != sellerID {
		return errors.New("permission denied")
	}

	return s.productRepo.DeleteImages(productID)
}

func (s *ProductService) CreateBrand(ctx context.Context, brand *model.Brand) (*model.Brand, error) {
	existing, err := s.productRepo.FindBrandByName(brand.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check brand: %w", err)
	}
	if existing != nil {
		return nil, errors.New("brand already exists")
	}

	err = s.productRepo.CreateBrand(brand)
	if err != nil {
		return nil, fmt.Errorf("failed to create brand: %w", err)
	}

	return brand, nil
}

func (s *ProductService) ListBrands(category model.ProductCategory) ([]model.Brand, error) {
	return s.productRepo.ListBrands(category)
}

func (s *ProductService) GetBrand(id uint) (*model.Brand, error) {
	brand, err := s.productRepo.FindBrandByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find brand: %w", err)
	}
	if brand == nil {
		return nil, errors.New("brand not found")
	}
	return brand, nil
}

func saveUploadedFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	return err
}
