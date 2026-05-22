package services

import (
	"encoding/json"
	"fmt"
	"time"

	"auction-system/internal/models"
	"auction-system/pkg/logger"
	"auction-system/pkg/mail"
	"auction-system/pkg/redis"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

const (
	NotifyTypeBidOutbid     = "bid_outbid"
	NotifyTypeAuctionEnding = "auction_ending"
	NotifyTypeBidSuccess    = "bid_success"
	NotifyTypeBidFailed     = "bid_failed"
	NotifyTypePayment       = "payment"
	NotifyTypeSystem        = "system"
)

func (s *NotificationService) CreateNotification(userID uint, notifyType, title, content string, relatedID uint) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:    userID,
		Type:      notifyType,
		Title:     title,
		Content:   content,
		RelatedID: relatedID,
		IsRead:    0,
		CreatedAt: time.Now(),
	}

	if err := models.DB.Create(notification).Error; err != nil {
		logger.Error("Failed to create notification: %v", err)
		return nil, err
	}

	s.publishNotification(notification)

	go s.sendEmailNotification(userID, notifyType, title, content)

	return notification, nil
}

func (s *NotificationService) publishNotification(notification *models.Notification) {
	data, _ := json.Marshal(notification)
	channel := fmt.Sprintf("notification:%d", notification.UserID)
	redis.Publish(channel, string(data))
}

func (s *NotificationService) sendEmailNotification(userID uint, notifyType, title, content string) {
	user, err := NewUserService().GetUserByID(userID)
	if err != nil || user.Email == "" {
		return
	}

	mail.SendMail(user.Email, title, content)
}

func (s *NotificationService) GetUserNotifications(userID uint, page, pageSize int, isRead *int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := models.DB.Model(&models.Notification{}).Where("user_id = ?", userID)
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error
	return notifications, total, err
}

func (s *NotificationService) MarkAsRead(userID uint, notificationIDs []uint) error {
	return models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND id IN ?", userID, notificationIDs).
		Update("is_read", 1).Error
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = 0", userID).
		Update("is_read", 1).Error
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = 0", userID).
		Count(&count).Error
	return count, err
}

func (s *NotificationService) NotifyBidOutbid(userID uint, auctionItemID uint, auctionName string, currentPrice float64) {
	title := fmt.Sprintf("您在【%s】的出价被超越了", auctionName)
	content := fmt.Sprintf("您的出价已被超越，当前价格：¥%.2f", currentPrice)
	s.CreateNotification(userID, NotifyTypeBidOutbid, title, content, auctionItemID)

	user, _ := NewUserService().GetUserByID(userID)
	if user != nil {
		mail.SendBidOutbidNotification(user.Email, auctionName, fmt.Sprintf("%.2f", currentPrice))
	}
}

func (s *NotificationService) NotifyAuctionEnding(userID uint, auctionItemID uint, auctionName string, timeLeft time.Duration) {
	title := fmt.Sprintf("拍卖即将结束：%s", auctionName)
	content := fmt.Sprintf("拍卖将在 %s 后结束，请抓紧时间出价", formatDuration(timeLeft))
	s.CreateNotification(userID, NotifyTypeAuctionEnding, title, content, auctionItemID)

	user, _ := NewUserService().GetUserByID(userID)
	if user != nil {
		mail.SendAuctionEndingSoonNotification(user.Email, auctionName, formatDuration(timeLeft))
	}
}

func (s *NotificationService) NotifyBidSuccess(userID uint, auctionItemID uint, auctionName string, price float64) {
	title := fmt.Sprintf("恭喜您成功拍得【%s】", auctionName)
	content := fmt.Sprintf("成交价格：¥%.2f，请尽快完成支付", price)
	s.CreateNotification(userID, NotifyTypeBidSuccess, title, content, auctionItemID)

	user, _ := NewUserService().GetUserByID(userID)
	if user != nil {
		mail.SendBidSuccessNotification(user.Email, auctionName, fmt.Sprintf("%.2f", price))
	}
}

func (s *NotificationService) NotifyPaymentSuccess(userID uint, auctionItemID uint, auctionName string, orderNo string) {
	title := fmt.Sprintf("支付成功：%s", auctionName)
	content := fmt.Sprintf("订单 %s 支付成功，请等待卖家发货", orderNo)
	s.CreateNotification(userID, NotifyTypePayment, title, content, auctionItemID)

	user, _ := NewUserService().GetUserByID(userID)
	if user != nil {
		mail.SendPaymentSuccessNotification(user.Email, auctionName, orderNo)
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, seconds)
	}
	return fmt.Sprintf("%d秒", seconds)
}
