package services

import (
	"errors"
	"fmt"
	"time"

	"auction-system/config"
	"auction-system/internal/dto"
	"auction-system/internal/models"
	"auction-system/pkg/logger"
	"auction-system/pkg/redis"
)

type AutoBidService struct {
	bidService    *BidService
	notifyService *NotificationService
}

func NewAutoBidService() *AutoBidService {
	return &AutoBidService{
		bidService:    NewBidService(),
		notifyService: NewNotificationService(),
	}
}

func (s *AutoBidService) SetAutoBid(userID uint, req *dto.SetAutoBidRequest) (*models.AutoBid, error) {
	item, err := NewAuctionItemService().GetItemByID(req.AuctionItemID)
	if err != nil {
		return nil, errors.New("拍卖品不存在")
	}

	if item.Status != ItemStatusOnline {
		return nil, errors.New("只能对上架中的拍卖品设置自动出价")
	}

	if req.MaxPrice <= item.CurrentPrice {
		return nil, errors.New("自动出价上限必须高于当前价格")
	}

	user, err := NewUserService().GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.MaxPrice > user.Balance {
		return nil, errors.New("自动出价上限不能超过账户余额")
	}

	var existingAutoBid models.AutoBid
	err = models.DB.Where("auction_item_id = ? AND user_id = ?", req.AuctionItemID, userID).First(&existingAutoBid).Error

	if err == nil {
		existingAutoBid.MaxPrice = req.MaxPrice
		existingAutoBid.UpdatedAt = time.Now()
		if err := models.DB.Save(&existingAutoBid).Error; err != nil {
			return nil, err
		}
		return &existingAutoBid, nil
	}

	autoBid := &models.AutoBid{
		AuctionItemID: req.AuctionItemID,
		UserID:        userID,
		MaxPrice:      req.MaxPrice,
		CurrentBid:    item.CurrentPrice,
		Status:        1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := models.DB.Create(autoBid).Error; err != nil {
		logger.Error("Failed to create auto bid: %v", err)
		return nil, err
	}

	s.cacheAutoBid(req.AuctionItemID, userID, req.MaxPrice)

	return autoBid, nil
}

func (s *AutoBidService) cacheAutoBid(itemID, userID uint, maxPrice float64) {
	key := fmt.Sprintf(BidAutoBidKey, itemID, userID)
	redis.Set(key, maxPrice, 24*time.Hour)
}

func (s *AutoBidService) CancelAutoBid(userID uint, autoBidID uint) error {
	var autoBid models.AutoBid
	if err := models.DB.First(&autoBid, autoBidID).Error; err != nil {
		return errors.New("自动出价记录不存在")
	}

	if autoBid.UserID != userID {
		return errors.New("无权取消此自动出价")
	}

	key := fmt.Sprintf(BidAutoBidKey, autoBid.AuctionItemID, userID)
	redis.Del(key)

	return models.DB.Delete(&autoBid).Error
}

func (s *AutoBidService) GetUserAutoBids(userID uint) ([]models.AutoBid, error) {
	var autoBids []models.AutoBid
	err := models.DB.Where("user_id = ? AND status = 1", userID).Preload("AuctionItem.Images").Find(&autoBids).Error
	return autoBids, err
}

func (s *AutoBidService) ProcessAutoBids(itemID uint, newPrice float64, triggerUserID uint) {
	logger.Info("Processing auto bids for item %d, new price: %.2f", itemID, newPrice)

	var autoBids []models.AutoBid
	models.DB.Where("auction_item_id = ? AND user_id != ? AND status = 1 AND max_price > ?", itemID, triggerUserID, newPrice).Find(&autoBids)

	if len(autoBids) == 0 {
		return
	}

	var highestAutoBid *models.AutoBid
	for _, ab := range autoBids {
		if highestAutoBid == nil || ab.MaxPrice > highestAutoBid.MaxPrice {
			highestAutoBid = &ab
		}
	}

	if highestAutoBid == nil {
		return
	}

	autoIncrement := config.AppConfig.Auction.AutoBidIncrement
	nextBid := newPrice + autoIncrement

	if nextBid > highestAutoBid.MaxPrice {
		nextBid = highestAutoBid.MaxPrice
	}

	if nextBid <= newPrice {
		return
	}

	lockKey := fmt.Sprintf(BidLockKey, itemID)
	locked, err := redis.Lock(lockKey, 5*time.Second)
	if err != nil || !locked {
		logger.Warn("Could not acquire lock for auto bid on item %d", itemID)
		return
	}
	defer redis.Unlock(lockKey)

	user, err := NewUserService().GetUserByID(highestAutoBid.UserID)
	if err != nil {
		return
	}

	if nextBid > user.Balance {
		return
	}

	tx := models.BeginTransaction()

	bid := &models.Bid{
		AuctionItemID: itemID,
		UserID:        highestAutoBid.UserID,
		Amount:        nextBid,
		MaxAutoBid:    highestAutoBid.MaxPrice,
		IsAutoBid:     1,
		IsWinning:     1,
		CreatedAt:     time.Now(),
	}

	if err := tx.Create(bid).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to create auto bid: %v", err)
		return
	}

	if err := tx.Model(&models.AuctionItem{}).Where("id = ?", itemID).Update("current_price", nextBid).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Model(&models.Bid{}).Where("auction_item_id = ? AND id != ?", itemID, bid.ID).Update("is_winning", 0).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Model(highestAutoBid).Update("current_bid", nextBid).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	s.bidService.cacheCurrentBid(itemID, nextBid, highestAutoBid.UserID)
	s.bidService.addBidToHistory(itemID, nextBid, highestAutoBid.UserID)
	s.bidService.publishBidMessage(itemID, highestAutoBid.UserID, user.Username, nextBid, 1)

	s.notifyService.NotifyBidOutbid(triggerUserID, itemID, "", nextBid)

	logger.Info("Auto bid placed successfully: item=%d, user=%d, amount=%.2f", itemID, highestAutoBid.UserID, nextBid)

	go s.ProcessAutoBids(itemID, nextBid, highestAutoBid.UserID)
}

func (s *AutoBidService) GetAutoBid(itemID, userID uint) (*models.AutoBid, error) {
	var autoBid models.AutoBid
	err := models.DB.Where("auction_item_id = ? AND user_id = ? AND status = 1", itemID, userID).First(&autoBid).Error
	if err != nil {
		return nil, err
	}
	return &autoBid, nil
}
