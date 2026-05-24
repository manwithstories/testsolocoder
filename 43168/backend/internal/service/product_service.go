package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
)

// 热门产品缓存 key 与 TTL
const (
	HotProductsCacheKey = "hot_products"
	HotProductsTTL      = 5 * time.Minute
)

// ProductService 产品业务逻辑层
type ProductService struct {
	repo  *repository.ProductRepository
	db    *gorm.DB
	rdb   *redis.Client
}

// NewProductService 创建产品服务
func NewProductService(repo *repository.ProductRepository, db *gorm.DB, rdb *redis.Client) *ProductService {
	return &ProductService{
		repo: repo,
		db:   db,
		rdb:  rdb,
	}
}

// toOptionModels 将 DTO 转换为模型
func toOptionModels(opts []dto.ProductOptionDTO) []*model.ProductOption {
	list := make([]*model.ProductOption, 0, len(opts))
	for _, o := range opts {
		list = append(list, &model.ProductOption{
			OptionType:  o.OptionType,
			OptionValue: o.OptionValue,
			PriceAdjust: o.PriceAdjust,
			Sort:        o.Sort,
		})
	}
	return list
}

// toImageModels 将 DTO 转换为模型
func toImageModels(imgs []dto.ProductImageDTO) []*model.ProductImage {
	list := make([]*model.ProductImage, 0, len(imgs))
	for _, img := range imgs {
		list = append(list, &model.ProductImage{
			ImageURL: img.ImageURL,
			Sort:     img.Sort,
		})
	}
	return list
}

// Create 发布产品
func (s *ProductService) Create(manufacturerID uint, req *dto.CreateProductRequest) (*model.Product, []*model.ProductOption, []*model.ProductImage, error) {
	if !model.ValidProductStatus(req.Status) {
		return nil, nil, nil, errors.New("产品状态不合法")
	}
	for _, o := range req.Options {
		if !model.ValidOptionType(o.OptionType) {
			return nil, nil, nil, errors.New("选项类型不合法")
		}
	}

	product := &model.Product{
		ManufacturerID: manufacturerID,
		Name:           req.Name,
		Category:       req.Category,
		Description:    req.Description,
		BasePrice:      req.BasePrice,
		Stock:          req.Stock,
		Status:         req.Status,
		IsHot:          req.IsHot,
	}

	options := toOptionModels(req.Options)
	images := toImageModels(req.Images)

	if err := s.repo.Create(product, options, images); err != nil {
		return nil, nil, nil, err
	}

	// 若标记为热门则清除缓存
	if product.IsHot {
		s.DeleteHotCache()
	}
	return product, options, images, nil
}

// Update 编辑产品
func (s *ProductService) Update(productID uint, manufacturerID uint, req *dto.UpdateProductRequest) (*model.Product, []*model.ProductOption, []*model.ProductImage, error) {
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, nil, nil, err
	}
	if product == nil {
		return nil, nil, nil, errors.New("产品不存在")
	}
	if product.ManufacturerID != manufacturerID {
		return nil, nil, nil, errors.New("无权修改他人的产品")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.BasePrice != nil {
		product.BasePrice = *req.BasePrice
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	oldStatus := product.Status
	if req.Status != nil {
		if !model.ValidProductStatus(*req.Status) {
			return nil, nil, nil, errors.New("产品状态不合法")
		}
		product.Status = *req.Status
	}
	wasHot := product.IsHot
	if req.IsHot != nil {
		product.IsHot = *req.IsHot
	}

	if err := s.repo.Update(product); err != nil {
		return nil, nil, nil, err
	}

	var options []*model.ProductOption
	var images []*model.ProductImage
	if req.Options != nil {
		for _, o := range req.Options {
			if !model.ValidOptionType(o.OptionType) {
				return nil, nil, nil, errors.New("选项类型不合法")
			}
		}
		options = toOptionModels(req.Options)
		if err := s.repo.ReplaceOptions(productID, options); err != nil {
			return nil, nil, nil, err
		}
	} else {
		options, _ = s.repo.GetOptions(productID)
	}
	if req.Images != nil {
		images = toImageModels(req.Images)
		if err := s.repo.ReplaceImages(productID, images); err != nil {
			return nil, nil, nil, err
		}
	} else {
		images, _ = s.repo.GetImages(productID)
	}

	// 热门标记或上架状态变化时，刷新缓存
	if wasHot != product.IsHot || oldStatus != product.Status {
		s.DeleteHotCache()
	}
	return product, options, images, nil
}

// Delete 删除产品
func (s *ProductService) Delete(productID uint, manufacturerID uint) error {
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("产品不存在")
	}
	if product.ManufacturerID != manufacturerID {
		return errors.New("无权删除他人的产品")
	}
	if err := s.repo.Delete(productID); err != nil {
		return err
	}
	if product.IsHot {
		s.DeleteHotCache()
	}
	return nil
}

// GetByID 根据 ID 获取产品详情（含选项、图片）
func (s *ProductService) GetByID(productID uint) (*model.Product, []*model.ProductOption, []*model.ProductImage, error) {
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, nil, nil, err
	}
	if product == nil {
		return nil, nil, nil, errors.New("产品不存在")
	}
	options, err := s.repo.GetOptions(productID)
	if err != nil {
		return nil, nil, nil, err
	}
	images, err := s.repo.GetImages(productID)
	if err != nil {
		return nil, nil, nil, err
	}
	return product, options, images, nil
}

// List 分页查询产品列表（返回产品与关联的选项、图片）
func (s *ProductService) List(params *dto.ProductListRequest) ([]*model.Product, map[uint][]*model.ProductOption, map[uint][]*model.ProductImage, int64, error) {
	p := &repository.ProductListParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.Keyword,
		Category: params.Category,
		Status:   params.Status,
		IsHot:    params.IsHot,
	}
	products, total, err := s.repo.List(p)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	ids := make([]uint, 0, len(products))
	for _, p := range products {
		ids = append(ids, p.ID)
	}
	opts, err := s.repo.BatchGetOptions(ids)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	imgs, err := s.repo.BatchGetImages(ids)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	return products, opts, imgs, total, nil
}

// HotProductItem 热门产品缓存项
type HotProductItem struct {
	ID             uint    `json:"id"`
	ManufacturerID uint    `json:"manufacturer_id"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	BasePrice      float64 `json:"base_price"`
	Cover          string  `json:"cover"`
}

// ListHot 获取热门产品列表（优先从 Redis 缓存获取）
func (s *ProductService) ListHot() ([]HotProductItem, error) {
	if s.rdb != nil {
		ctx := context.Background()
		data, err := s.rdb.Get(ctx, HotProductsCacheKey).Bytes()
		if err == nil {
			var items []HotProductItem
			if jErr := json.Unmarshal(data, &items); jErr == nil {
				return items, nil
			}
		}
	}

	products, err := s.repo.ListHot(20)
	if err != nil {
		return nil, err
	}
	items := make([]HotProductItem, 0, len(products))
	for _, p := range products {
		cover, _ := s.repo.GetCoverImage(p.ID)
		items = append(items, HotProductItem{
			ID:             p.ID,
			ManufacturerID: p.ManufacturerID,
			Name:           p.Name,
			Category:       p.Category,
			BasePrice:      p.BasePrice,
			Cover:          cover,
		})
	}

	// 写回缓存
	if s.rdb != nil {
		data, jErr := json.Marshal(items)
		if jErr == nil {
			ctx := context.Background()
			_ = s.rdb.Set(ctx, HotProductsCacheKey, data, HotProductsTTL).Err()
		}
	}
	return items, nil
}

// DeleteHotCache 删除热门产品缓存
func (s *ProductService) DeleteHotCache() {
	if s.rdb == nil {
		return
	}
	ctx := context.Background()
	_ = s.rdb.Del(ctx, HotProductsCacheKey).Err()
}

// GetOptions 获取产品选项
func (s *ProductService) GetOptions(productID uint) ([]*model.ProductOption, error) {
	return s.repo.GetOptions(productID)
}

// GetImages 获取产品图片
func (s *ProductService) GetImages(productID uint) ([]*model.ProductImage, error) {
	return s.repo.GetImages(productID)
}
