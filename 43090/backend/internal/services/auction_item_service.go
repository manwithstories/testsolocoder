package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"auction-system/internal/dto"
	"auction-system/internal/models"
	"auction-system/pkg/logger"
)

type AuctionItemService struct{}

func NewAuctionItemService() *AuctionItemService {
	return &AuctionItemService{}
}

const (
	ItemStatusDraft   = 0
	ItemStatusOnline  = 1
	ItemStatusOffline = 2
	ItemStatusSold    = 3
)

func (s *AuctionItemService) CreateItem(sellerID uint, req *dto.CreateAuctionItemRequest) (*models.AuctionItem, error) {
	item := &models.AuctionItem{
		Title:        req.Title,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		SellerID:     sellerID,
		StartPrice:   req.StartPrice,
		ReservePrice: req.ReservePrice,
		CurrentPrice: req.StartPrice,
		Location:     req.Location,
		Condition:    req.Condition,
		Status:       ItemStatusDraft,
	}

	if err := models.DB.Create(item).Error; err != nil {
		logger.Error("Failed to create auction item: %v", err)
		return nil, err
	}

	return item, nil
}

func (s *AuctionItemService) UpdateItem(itemID uint, sellerID uint, req *dto.UpdateAuctionItemRequest) error {
	var item models.AuctionItem
	if err := models.DB.First(&item, itemID).Error; err != nil {
		return errors.New("拍卖品不存在")
	}

	if item.SellerID != sellerID {
		return errors.New("无权修改此拍卖品")
	}

	if item.Status != ItemStatusDraft {
		return errors.New("只能修改草稿状态的拍卖品")
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.StartPrice > 0 {
		updates["start_price"] = req.StartPrice
	}
	if req.ReservePrice > 0 {
		updates["reserve_price"] = req.ReservePrice
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Condition != "" {
		updates["condition"] = req.Condition
	}
	updates["updated_at"] = time.Now()

	return models.DB.Model(&item).Updates(updates).Error
}

func (s *AuctionItemService) GetItemByID(itemID uint) (*models.AuctionItem, error) {
	var item models.AuctionItem
	if err := models.DB.Preload("Category").Preload("Seller").Preload("Images").First(&item, itemID).Error; err != nil {
		return nil, err
	}

	models.DB.Model(&item).Update("view_count", item.ViewCount+1)
	return &item, nil
}

func (s *AuctionItemService) GetItemList(query *dto.AuctionItemQuery) ([]models.AuctionItem, int64, error) {
	var items []models.AuctionItem
	var total int64

	dbQuery := models.DB.Model(&models.AuctionItem{}).Preload("Category").Preload("Images")

	if query.Keyword != "" {
		dbQuery = dbQuery.Where("title LIKE ? OR description LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.CategoryID > 0 {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
	}
	if query.Status != nil {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}

	dbQuery.Count(&total)

	offset := (query.Page - 1) * query.PageSize
	sortStr := fmt.Sprintf("%s %s", query.SortBy, query.SortOrder)
	err := dbQuery.Order(sortStr).Offset(offset).Limit(query.PageSize).Find(&items).Error
	return items, total, err
}

func (s *AuctionItemService) OnlineItem(itemID uint, sellerID uint) error {
	var item models.AuctionItem
	if err := models.DB.First(&item, itemID).Error; err != nil {
		return errors.New("拍卖品不存在")
	}

	if item.SellerID != sellerID {
		return errors.New("无权操作此拍卖品")
	}

	if item.Status != ItemStatusDraft {
		return errors.New("只能上架草稿状态的拍卖品")
	}

	return models.DB.Model(&item).Update("status", ItemStatusOnline).Error
}

func (s *AuctionItemService) OfflineItem(itemID uint, sellerID uint) error {
	var item models.AuctionItem
	if err := models.DB.First(&item, itemID).Error; err != nil {
		return errors.New("拍卖品不存在")
	}

	if item.SellerID != sellerID {
		return errors.New("无权操作此拍卖品")
	}

	if item.Status != ItemStatusOnline {
		return errors.New("只能下架上架状态的拍卖品")
	}

	return models.DB.Model(&item).Update("status", ItemStatusOffline).Error
}

func (s *AuctionItemService) DeleteItem(itemID uint, sellerID uint) error {
	var item models.AuctionItem
	if err := models.DB.First(&item, itemID).Error; err != nil {
		return errors.New("拍卖品不存在")
	}

	if item.SellerID != sellerID {
		return errors.New("无权删除此拍卖品")
	}

	if item.Status != ItemStatusDraft {
		return errors.New("只能删除草稿状态的拍卖品")
	}

	return models.DB.Delete(&item).Error
}

func (s *AuctionItemService) UploadImages(itemID uint, sellerID uint, c *gin.Context) ([]models.AuctionImage, error) {
	var item models.AuctionItem
	if err := models.DB.First(&item, itemID).Error; err != nil {
		return nil, errors.New("拍卖品不存在")
	}

	if item.SellerID != sellerID {
		return nil, errors.New("无权操作此拍卖品")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File["images"]
	var images []models.AuctionImage

	uploadDir := "./uploads/auction_items"
	os.MkdirAll(uploadDir, 0755)

	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		savePath := filepath.Join(uploadDir, newFilename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			logger.Error("Failed to save uploaded file: %v", err)
			continue
		}

		image := models.AuctionImage{
			AuctionItemID: itemID,
			URL:           fmt.Sprintf("/uploads/auction_items/%s", newFilename),
			Sort:          i,
			IsMain:        0,
			CreatedAt:     time.Now(),
		}

		if i == 0 {
			image.IsMain = 1
		}

		if err := models.DB.Create(&image).Error; err != nil {
			logger.Error("Failed to save image record: %v", err)
			continue
		}

		images = append(images, image)
	}

	return images, nil
}

func (s *AuctionItemService) GetSellerItems(sellerID uint, page, pageSize int, status *int) ([]models.AuctionItem, int64, error) {
	var items []models.AuctionItem
	var total int64

	query := models.DB.Model(&models.AuctionItem{}).Where("seller_id = ?", sellerID).Preload("Images")
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (s *AuctionItemService) UpdateCurrentPrice(itemID uint, price float64) error {
	return models.DB.Model(&models.AuctionItem{}).Where("id = ?", itemID).Update("current_price", price).Error
}

func (s *AuctionItemService) IncrementBidCount(itemID uint) error {
	return models.DB.Model(&models.AuctionItem{}).Where("id = ?", itemID).UpdateColumn("bid_count", models.DB.Raw("bid_count + ?", 1)).Error
}
