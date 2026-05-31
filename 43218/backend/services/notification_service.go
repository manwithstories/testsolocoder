package services

import (
	"encoding/json"
	"errors"
	"time"

	"secondhand-platform/database"
	"secondhand-platform/models"

	"github.com/sirupsen/logrus"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) CreateNotification(userID uint, notificationType, title, content, orderNo string, extraData map[string]interface{}) (*models.Notification, error) {
	var extraDataStr string
	if extraData != nil {
		data, _ := json.Marshal(extraData)
		extraDataStr = string(data)
	}

	notification := &models.Notification{
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Content:   content,
		OrderNo:   orderNo,
		ExtraData: extraDataStr,
		IsRead:    false,
	}

	result := database.DB.Create(notification)
	if result.Error != nil {
		return nil, result.Error
	}

	return notification, nil
}

func (s *NotificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := database.DB.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (s *NotificationService) ListNotifications(userID uint, page, pageSize int, isRead *bool) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	db := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID)
	if isRead != nil {
		db = db.Where("is_read = ?", *isRead)
	}

	db.Count(&total)
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (s *NotificationService) MarkAsRead(userID, notificationID uint) error {
	var notification models.Notification
	if err := database.DB.First(&notification, notificationID).Error; err != nil {
		return errors.New("通知不存在")
	}

	if notification.UserID != userID {
		return errors.New("无权操作此通知")
	}

	notification.IsRead = true
	return database.DB.Save(&notification).Error
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (s *NotificationService) DeleteNotification(userID, notificationID uint) error {
	var notification models.Notification
	if err := database.DB.First(&notification, notificationID).Error; err != nil {
		return errors.New("通知不存在")
	}

	if notification.UserID != userID {
		return errors.New("无权操作此通知")
	}

	return database.DB.Delete(&notification).Error
}

func (s *NotificationService) SendOrderStatusNotification(userID uint, orderNo string, status int, orderType string) {
	var statusText string
	if orderType == "product" {
		statusText = models.OrderStatusText[status]
	} else {
		statusText = models.RepairStatusText[status]
	}

	title := "订单状态更新"
	content := "您的订单 " + orderNo + " 状态已更新为: " + statusText

	s.CreateNotification(userID, models.NotificationTypeOrderStatus, title, content, orderNo, map[string]interface{}{
		"status":     status,
		"order_type": orderType,
	})
}

func (s *NotificationService) SendPaymentNotification(userID uint, orderNo string, amount float64, paymentType string) {
	title := "支付通知"
	content := "您有一笔 " + paymentType + " 支付，金额: ¥" + formatMoney(amount) + "，订单号: " + orderNo

	s.CreateNotification(userID, models.NotificationTypePayment, title, content, orderNo, map[string]interface{}{
		"amount":       amount,
		"payment_type": paymentType,
	})
}

func (s *NotificationService) SendReviewReminder(userID uint, orderNo string, orderType string) {
	title := "评价提醒"
	content := "您的订单 " + orderNo + " 已完成，请及时评价"

	s.CreateNotification(userID, models.NotificationTypeReview, title, content, orderNo, map[string]interface{}{
		"order_type": orderType,
	})
}

func (s *NotificationService) SendSystemNotification(userID uint, title, content string) {
	s.CreateNotification(userID, models.NotificationTypeSystem, title, content, "", nil)
}

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) SendMessage(senderID, receiverID uint, content, msgType string) (*models.Message, error) {
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType,
		IsRead:     false,
	}

	result := database.DB.Create(message)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

func (s *MessageService) ListMessages(userID, otherUserID uint, page, pageSize int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	db := database.DB.Model(&models.Message{}).Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, otherUserID, otherUserID, userID,
	)

	db.Count(&total)
	if err := db.Preload("Sender").Preload("Receiver").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (s *MessageService) MarkMessagesAsRead(userID, otherUserID uint) error {
	return database.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", otherUserID, userID, false).
		Update("is_read", true).Error
}

func (s *MessageService) GetUnreadMessageCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (s *MessageService) GetMessageContacts(userID uint) ([]map[string]interface{}, error) {
	type Contact struct {
		UserID    uint
		LastMsg   string
		LastTime  time.Time
		UnreadCnt int64
	}

	var contacts []Contact
	database.DB.Raw(`
		SELECT 
			CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as user_id,
			content as last_msg,
			created_at as last_time,
			(SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND sender_id = user_id AND is_read = false) as unread_cnt
		FROM messages 
		WHERE sender_id = ? OR receiver_id = ?
		ORDER BY created_at DESC
	`, userID, userID, userID, userID).Scan(&contacts)

	var result []map[string]interface{}
	seen := make(map[uint]bool)
	for _, c := range contacts {
		if !seen[c.UserID] {
			seen[c.UserID] = true
			result = append(result, map[string]interface{}{
				"user_id":     c.UserID,
				"last_msg":    c.LastMsg,
				"last_time":   c.LastTime,
				"unread_count": c.UnreadCnt,
			})
		}
	}

	return result, nil
}

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) GetDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	stats["total_users"] = userCount

	var sellerCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleSeller).Count(&sellerCount)
	stats["sellers"] = sellerCount

	var buyerCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleBuyer).Count(&buyerCount)
	stats["buyers"] = buyerCount

	var technicianCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleTechnician).Count(&technicianCount)
	stats["technicians"] = technicianCount

	var productCount int64
	database.DB.Model(&models.Product{}).Count(&productCount)
	stats["total_products"] = productCount

	var pendingProducts int64
	database.DB.Model(&models.Product{}).Where("status = ?", models.ProductStatusPending).Count(&pendingProducts)
	stats["pending_products"] = pendingProducts

	var orderCount int64
	database.DB.Model(&models.Order{}).Count(&orderCount)
	stats["total_orders"] = orderCount

	var repairOrderCount int64
	database.DB.Model(&models.RepairOrder{}).Count(&repairOrderCount)
	stats["total_repair_orders"] = repairOrderCount

	var totalAmount float64
	database.DB.Model(&models.Order{}).
		Where("status IN ?", []int{models.OrderStatusCompleted, models.OrderStatusDelivered}).
		Select("COALESCE(SUM(final_price), 0)").Scan(&totalAmount)
	stats["total_amount"] = totalAmount

	var totalCommission float64
	database.DB.Model(&models.Order{}).
		Where("status IN ?", []int{models.OrderStatusCompleted, models.OrderStatusDelivered}).
		Select("COALESCE(SUM(commission), 0)").Scan(&totalCommission)
	stats["total_commission"] = totalCommission

	return stats, nil
}

func (s *AdminService) GetTransactionStats(startDate, endDate string) ([]map[string]interface{}, error) {
	type DailyStat struct {
		Date        string  `gorm:"column:date"`
		OrderCount  int64   `gorm:"column:order_count"`
		TotalAmount float64 `gorm:"column:total_amount"`
		Commission  float64 `gorm:"column:commission"`
	}

	var stats []DailyStat
	database.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as order_count,
			COALESCE(SUM(final_price), 0) as total_amount,
			COALESCE(SUM(commission), 0) as commission
		FROM orders
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, startDate, endDate).Scan(&stats)

	var result []map[string]interface{}
	for _, s := range stats {
		result = append(result, map[string]interface{}{
			"date":         s.Date,
			"order_count":  s.OrderCount,
			"total_amount": s.TotalAmount,
			"commission":   s.Commission,
		})
	}

	return result, nil
}

func (s *AdminService) GetUserActivityStats(days int) ([]map[string]interface{}, error) {
	type DailyActivity struct {
		Date       string `gorm:"column:date"`
		LoginCount int64  `gorm:"column:login_count"`
		NewUsers   int64  `gorm:"column:new_users"`
	}

	var stats []DailyActivity
	database.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as new_users
		FROM users
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, days).Scan(&stats)

	var result []map[string]interface{}
	for _, s := range stats {
		result = append(result, map[string]interface{}{
			"date":       s.Date,
			"login_count": s.LoginCount,
			"new_users":  s.NewUsers,
		})
	}

	return result, nil
}

func (s *AdminService) ListAllTransactions(page, pageSize int, transactionType string, status int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	db := database.DB.Model(&models.Transaction{})
	if transactionType != "" {
		db = db.Where("type = ?", transactionType)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (s *AdminService) AddAdminLog(adminID uint, action, targetType string, targetID *uint, content, ipAddress string) {
	log := &models.AdminLog{
		AdminID:    adminID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Content:    content,
		IPAddress:  ipAddress,
	}
	database.DB.Create(log)
}

func (s *AdminService) ListAdminLogs(page, pageSize int) ([]models.AdminLog, int64, error) {
	var logs []models.AdminLog
	var total int64

	db := database.DB.Model(&models.AdminLog{})
	db.Count(&total)
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func formatMoney(amount float64) string {
	return floatToStr(amount)
}

func floatToStr(f float64) string {
	return floatToStrWithPrec(f, 2)
}

func floatToStrWithPrec(f float64, prec int) string {
	return formatFloat(f, prec)
}

func formatFloat(f float64, prec int) string {
	s := fmt.Sprintf("%.*f", prec, f)
	return s
}

func init() {
	logrus.Info("Notification service initialized")
}
