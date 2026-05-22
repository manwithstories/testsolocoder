package utils

import (
	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"time"
)

func SendMessage(userID uint, msgType models.MessageType, title, content string, relatedID *uint) error {
	msg := models.Message{
		UserID:    userID,
		Type:      msgType,
		Title:     title,
		Content:   content,
		RelatedID: relatedID,
		CreatedAt: time.Now(),
	}
	return config.DB.Create(&msg).Error
}

func SendOrderNotification(userID uint, orderNo string, orderStatus string) {
	title := "订单状态更新"
	content := "您的订单 " + orderNo + " 状态已更新为: " + orderStatus
	_ = SendMessage(userID, models.MessageTypeOrder, title, content, nil)
}

func SendInvitationNotification(providerID uint, orderNo string) {
	title := "新的预约邀请"
	content := "您收到一个新的预约邀请，订单号: " + orderNo + "，请及时处理"
	_ = SendMessage(providerID, models.MessageTypeInvitation, title, content, nil)
}

func SendReviewNotification(userID uint, orderNo string) {
	title := "评价提醒"
	content := "订单 " + orderNo + " 已完成，请及时评价服务"
	_ = SendMessage(userID, models.MessageTypeReview, title, content, nil)
}

func SendSystemNotification(userID uint, title, content string) {
	_ = SendMessage(userID, models.MessageTypeSystem, title, content, nil)
}
