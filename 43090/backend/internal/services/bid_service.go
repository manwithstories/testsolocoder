package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"auction-system/config"
	"auction-system/internal/dto"
	"auction-system/internal/models"
	"auction-system/pkg/logger"
	"auction-system/pkg/redis"
)

type BidService struct {
	notifyService *NotificationService
	itemService   *AuctionItemService
	autoBidService *AutoBidService
}

func NewBidService() *BidService {
	return &BidService{
		notifyService: NewNotificationService(),
		itemService:   NewAuctionItemService(),
		autoBidService: NewAutoBidService(),
	}
}

const (
	BidLockKey      = "bid:lock:%d"
	BidCurrentKey   = "bid:current:%d"
	BidHistoryKey   = "bid:history:%d"
	BidAutoBidKey   = "bid:auto:%d:%d"
	BidChannel      = "bid:channel:%d"
)

type BidMessage struct {
	UserID        uint      `json:"user_id"`
	Username      string    `json:"username"`
	AuctionItemID uint      `json:"auction_item_id"`
	Amount        float64   `json:"amount"`
	IsAutoBid     int       `json:"is_auto_bid"`
	CreatedAt     time.Time `json:"created_at"`
}

func (s *BidService) PlaceBid(userID uint, itemID uint, req *dto.BidRequest) (*models.Bid, error) {
	lockKey := fmt.Sprintf(BidLockKey, itemID)
	locked, err := redis.Lock(lockKey, 5*time.Second)
	if err != nil {
		logger.Error("Failed to acquire bid lock: %v", err)
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	if !locked {
		return nil, errors.New("其他用户正在出价，请稍后再试")
	}
	defer redis.Unlock(lockKey)

	item, err := s.itemService.GetItemByID(itemID)
	if err != nil {
		return nil, errors.New("拍卖品不存在")
	}

	if item.Status != ItemStatusOnline {
		return nil, errors.New("该拍卖品未上架")
	}

	user, err := NewUserService().GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != 1 {
		return nil, errors.New("用户账户异常")
	}

	currentPrice := s.getCurrentPrice(itemID, item.CurrentPrice)
	minIncrement := config.AppConfig.Auction.MinIncrement

	minRequiredPrice := currentPrice + minIncrement
	if req.Amount < minRequiredPrice {
		return nil, fmt.Errorf("出价必须至少为 ¥%.2f", minRequiredPrice)
	}

	if req.Amount > user.Balance {
		return nil, errors.New("余额不足")
	}

	if req.MaxAutoBid > 0 && req.MaxAutoBid < req.Amount {
		return nil, errors.New("自动出价上限不能低于当前出价")
	}

	prevHighestUserID := s.getHighestBidUserID(itemID)

	tx := models.BeginTransaction()

	bid := &models.Bid{
		AuctionItemID: itemID,
		UserID:        userID,
		Amount:        req.Amount,
		MaxAutoBid:    req.MaxAutoBid,
		IsAutoBid:     0,
		IsWinning:     1,
		CreatedAt:     time.Now(),
	}

	if err := tx.Create(bid).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to create bid: %v", err)
		return nil, errors.New("出价失败")
	}

	if err := tx.Model(&models.AuctionItem{}).Where("id = ?", itemID).Updates(map[string]interface{}{
		"current_price": req.Amount,
		"bid_count":     item.BidCount + 1,
	}).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to update auction item: %v", err)
		return nil, errors.New("出价失败")
	}

	if err := tx.Model(&models.Bid{}).Where("auction_item_id = ? AND id != ?", itemID, bid.ID).Update("is_winning", 0).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to update previous bids: %v", err)
		return nil, errors.New("出价失败")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to commit transaction: %v", err)
		return nil, errors.New("出价失败")
	}

	s.cacheCurrentBid(itemID, req.Amount, userID)
	s.addBidToHistory(itemID, req.Amount, userID)

	if prevHighestUserID > 0 && prevHighestUserID != userID {
		s.notifyService.NotifyBidOutbid(prevHighestUserID, itemID, item.Title, req.Amount)
	}

	go s.autoBidService.ProcessAutoBids(itemID, req.Amount, userID)

	s.publishBidMessage(itemID, userID, user.Username, req.Amount, 0)

	return bid, nil
}

func (s *BidService) getCurrentPrice(itemID uint, fallback float64) float64 {
	currentKey := fmt.Sprintf(BidCurrentKey, itemID)
	val, err := redis.Get(currentKey)
	if err != nil || val == "" {
		return fallback
	}
	var price float64
	fmt.Sscanf(val, "%f", &price)
	return price
}

func (s *BidService) cacheCurrentBid(itemID uint, amount float64, userID uint) {
	currentKey := fmt.Sprintf(BidCurrentKey, itemID)
	data := fmt.Sprintf("%.2f:%d", amount, userID)
	redis.Set(currentKey, data, 24*time.Hour)
}

func (s *BidService) addBidToHistory(itemID uint, amount float64, userID uint) {
	historyKey := fmt.Sprintf(BidHistoryKey, itemID)
	member := fmt.Sprintf("%d:%.2f:%d", time.Now().Unix(), amount, userID)
	redis.ZAdd(historyKey, float64(time.Now().Unix()), member)
}

func (s *BidService) getHighestBidUserID(itemID uint) uint {
	currentKey := fmt.Sprintf(BidCurrentKey, itemID)
	val, err := redis.Get(currentKey)
	if err != nil || val == "" {
		var bid models.Bid
		if err := models.DB.Where("auction_item_id = ? AND is_winning = 1", itemID).Order("created_at DESC").First(&bid).Error; err == nil {
			return bid.UserID
		}
		return 0
	}
	var userID uint
	fmt.Sscanf(val, "%*f:%d", &userID)
	return userID
}

func (s *BidService) publishBidMessage(itemID, userID uint, username string, amount float64, isAutoBid int) {
	channel := fmt.Sprintf(BidChannel, itemID)
	msg := BidMessage{
		UserID:        userID,
		Username:      username,
		AuctionItemID: itemID,
		Amount:        amount,
		IsAutoBid:     isAutoBid,
		CreatedAt:     time.Now(),
	}
	data, _ := json.Marshal(msg)
	redis.Publish(channel, string(data))
}

func (s *BidService) GetBidHistory(itemID uint, page, pageSize int) ([]models.Bid, int64, error) {
	var bids []models.Bid
	var total int64

	query := models.DB.Model(&models.Bid{}).Where("auction_item_id = ?", itemID).Preload("User")
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&bids).Error
	return bids, total, err
}

func (s *BidService) GetUserBids(userID uint, page, pageSize int) ([]models.Bid, int64, error) {
	var bids []models.Bid
	var total int64

	query := models.DB.Model(&models.Bid{}).Where("user_id = ?", userID).Preload("AuctionItem.Images")
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&bids).Error
	return bids, total, err
}

func (s *BidService) GetCurrentBid(itemID uint) (*models.Bid, error) {
	var bid models.Bid
	err := models.DB.Where("auction_item_id = ? AND is_winning = 1", itemID).Preload("User").First(&bid).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}
