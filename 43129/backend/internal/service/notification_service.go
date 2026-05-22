package service

import (
	"fmt"
	"time"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	userRepo         *repository.UserRepository
}

func NewNotificationService(notificationRepo *repository.NotificationRepository, userRepo *repository.UserRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

type SendNotificationRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Type    string `json:"type"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (s *NotificationService) SendNotification(req *SendNotificationRequest) error {
	notification := &model.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Content: req.Content,
		IsRead:  false,
	}

	return s.notificationRepo.Create(notification)
}

func (s *NotificationService) SendAppointmentNotification(appointment *model.Appointment) error {
	customerTitle := "预约成功提醒"
	customerContent := fmt.Sprintf("您的预约已成功！\n服务：%s\n技师：%s\n时间：%s %s-%s",
		appointment.Service.Name,
		appointment.Technician.Name,
		appointment.AppointmentDate.Format("2006-01-02"),
		appointment.StartTime,
		appointment.EndTime,
	)

	if appointment.Customer != nil && appointment.Customer.User != nil {
		if err := s.SendNotification(&SendNotificationRequest{
			UserID:  appointment.Customer.User.ID,
			Type:    "appointment",
			Title:   customerTitle,
			Content: customerContent,
		}); err != nil {
			return err
		}
	}

	techTitle := "新预约提醒"
	techContent := fmt.Sprintf("您有新的预约！\n服务：%s\n顾客：%s\n时间：%s %s-%s",
		appointment.Service.Name,
		appointment.Customer.Name,
		appointment.AppointmentDate.Format("2006-01-02"),
		appointment.StartTime,
		appointment.EndTime,
	)

	if appointment.Technician != nil && appointment.Technician.User != nil {
		return s.SendNotification(&SendNotificationRequest{
			UserID:  appointment.Technician.User.ID,
			Type:    "appointment",
			Title:   techTitle,
			Content: techContent,
		})
	}

	return nil
}

func (s *NotificationService) SendCancelNotification(appointment *model.Appointment) error {
	customerTitle := "预约取消提醒"
	customerContent := fmt.Sprintf("您的预约已取消！\n服务：%s\n时间：%s %s-%s\n原因：%s",
		appointment.Service.Name,
		appointment.AppointmentDate.Format("2006-01-02"),
		appointment.StartTime,
		appointment.EndTime,
		appointment.CancelReason,
	)

	if appointment.Customer != nil && appointment.Customer.User != nil {
		if err := s.SendNotification(&SendNotificationRequest{
			UserID:  appointment.Customer.User.ID,
			Type:    "appointment_cancel",
			Title:   customerTitle,
			Content: customerContent,
		}); err != nil {
			return err
		}
	}

	techTitle := "预约取消通知"
	techContent := fmt.Sprintf("顾客取消了预约！\n服务：%s\n顾客：%s\n时间：%s %s-%s",
		appointment.Service.Name,
		appointment.Customer.Name,
		appointment.AppointmentDate.Format("2006-01-02"),
		appointment.StartTime,
		appointment.EndTime,
	)

	if appointment.Technician != nil && appointment.Technician.User != nil {
		return s.SendNotification(&SendNotificationRequest{
			UserID:  appointment.Technician.User.ID,
			Type:    "appointment_cancel",
			Title:   techTitle,
			Content: techContent,
		})
	}

	return nil
}

func (s *NotificationService) SendDailySchedule(technicians []model.Technician, appointments map[uint][]model.Appointment) error {
	for _, tech := range technicians {
		todayAppointments := appointments[tech.ID]
		if len(todayAppointments) == 0 {
			continue
		}

		content := fmt.Sprintf("今日工作安排（%d个预约）：\n", len(todayAppointments))
		for i, appt := range todayAppointments {
			content += fmt.Sprintf("%d. %s-%s %s - %s\n",
				i+1, appt.StartTime, appt.EndTime, appt.Customer.Name, appt.Service.Name)
		}

		user, err := s.userRepo.GetByID(tech.UserID)
		if err == nil {
			s.SendNotification(&SendNotificationRequest{
				UserID:  user.ID,
				Type:    "daily_schedule",
				Title:   "今日工作安排",
				Content: content,
			})
		}
	}
	return nil
}

func (s *NotificationService) SendLowStockAlert(products []model.Product) error {
	if len(products) == 0 {
		return nil
	}

	content := "以下产品库存不足：\n"
	for _, product := range products {
		content += fmt.Sprintf("- %s：当前库存 %d，阈值 %d\n",
			product.Name, product.Stock, product.Threshold)
	}

	return s.notificationRepo.Create(&model.Notification{
		UserID:  1,
		Type:    "low_stock",
		Title:   "库存预警",
		Content: content,
		IsRead:  false,
	})
}

func (s *NotificationService) GetUserNotifications(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	return s.notificationRepo.GetByUserID(userID, page, pageSize)
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

func (s *NotificationService) MarkAsRead(userID, notificationID uint) error {
	return s.notificationRepo.MarkAsRead(userID, notificationID)
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *NotificationService) SendPaymentNotification(userID uint, payment *model.Payment) error {
	title := "支付成功提醒"
	content := fmt.Sprintf("您的支付已成功！\n金额：%.2f元\n支付方式：%s\n交易号：%s",
		payment.Amount,
		s.getPayMethodText(payment.PayMethod),
		payment.TransactionNo,
	)

	return s.SendNotification(&SendNotificationRequest{
		UserID:  userID,
		Type:    "payment",
		Title:   title,
		Content: content,
	})
}

func (s *NotificationService) getPayMethodText(method string) string {
	switch method {
	case "cash":
		return "现金"
	case "card":
		return "会员卡"
	case "points":
		return "积分抵扣"
	case "wechat":
		return "微信支付"
	case "alipay":
		return "支付宝"
	default:
		return method
	}
}

func (s *NotificationService) SendTechnicianLeaveNotification(technician *model.Technician, leaveDate time.Time) error {
	title := "请假审批通过"
	content := fmt.Sprintf("您的请假申请已通过！\n请假日期：%s", leaveDate.Format("2006-01-02"))

	user, err := s.userRepo.GetByID(technician.UserID)
	if err != nil {
		return err
	}

	return s.SendNotification(&SendNotificationRequest{
		UserID:  user.ID,
		Type:    "leave",
		Title:   title,
		Content: content,
	})
}

func (s *NotificationService) SendAppointmentRescheduleNotification(appointment *model.Appointment) error {
	if appointment.Customer == nil || appointment.Customer.User == nil {
		return nil
	}

	title := "预约调整通知"
	content := fmt.Sprintf("由于技师请假，您的预约需要调整\n服务：%s\n原时间：%s %s-%s\n请联系客服重新预约",
		appointment.Service.Name,
		appointment.AppointmentDate.Format("2006-01-02"),
		appointment.StartTime,
		appointment.EndTime,
	)

	return s.SendNotification(&SendNotificationRequest{
		UserID:  appointment.Customer.User.ID,
		Type:    "appointment_reschedule",
		Title:   title,
		Content: content,
	})
}
