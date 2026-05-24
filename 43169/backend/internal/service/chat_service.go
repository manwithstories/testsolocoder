package service

import (
	"errors"
	"strings"
	"time"

	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/repository"
	"matchmaking-platform/internal/utils"

	"gorm.io/gorm"
)

type ChatService struct {
	chatRepo       *repository.ChatRepo
	sensitiveRepo  *repository.SensitiveWordRepo
	logRepo        *repository.SystemLogRepo
}

func NewChatService() *ChatService {
	return &ChatService{
		chatRepo:      repository.NewChatRepo(),
		sensitiveRepo: repository.NewSensitiveWordRepo(),
		logRepo:       repository.NewSystemLogRepo(),
	}
}

func (s *ChatService) SendMessage(senderID, receiverID uint, msgType string, content string) (*model.ChatMessage, error) {
	if senderID == receiverID {
		return nil, errors.New("不能给自己发送消息")
	}

	sensitiveWords, _ := s.sensitiveRepo.ListAll()
	var wordList []string
	for _, w := range sensitiveWords {
		wordList = append(wordList, w.Word)
	}
	filteredContent := utils.FilterSensitiveWords(content, wordList)

	msg := &model.ChatMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Type:       model.MessageType(msgType),
		Content:    filteredContent,
	}

	if err := s.chatRepo.SendMessage(msg); err != nil {
		return nil, err
	}

	session, _ := s.chatRepo.GetOrCreateSession(senderID, receiverID)
	lastMsg := filteredContent
	if len(lastMsg) > 100 {
		lastMsg = lastMsg[:100]
	}

	if session.UserAID == senderID {
		session.UnreadB++
	} else {
		session.UnreadA++
	}
	session.LastMessage = lastMsg
	session.LastTime = time.Now()
	s.chatRepo.UpdateSession(session)

	s.logRepo.Create(&model.SystemLog{
		UserID: senderID,
		Module: "chat",
		Action: "send_message",
	})

	return msg, nil
}

func (s *ChatService) GetHistory(userID, otherID uint, page, pageSize int) ([]model.ChatMessage, int64, error) {
	s.chatRepo.MarkAsRead(otherID, userID)
	return s.chatRepo.GetHistory(userID, otherID, page, pageSize)
}

func (s *ChatService) GetUnreadCount(userID uint) (int64, error) {
	return s.chatRepo.GetUnreadCount(userID)
}

func (s *ChatService) GetSessions(userID uint, page, pageSize int) ([]model.ChatSession, int64, error) {
	return s.chatRepo.ListSessions(userID, page, pageSize)
}

func (s *ChatService) MarkAsRead(senderID, receiverID uint) error {
	return s.chatRepo.MarkAsRead(senderID, receiverID)
}

func (s *ChatService) ContainsSensitiveWords(content string) bool {
	sensitiveWords, _ := s.sensitiveRepo.ListAll()
	for _, w := range sensitiveWords {
		if strings.Contains(content, w.Word) {
			return true
		}
	}
	return false
}

type MemberService struct {
	memberRepo *repository.MemberRepo
	userRepo   *repository.UserRepo
	logRepo    *repository.SystemLogRepo
}

func NewMemberService() *MemberService {
	return &MemberService{
		memberRepo: repository.NewMemberRepo(),
		userRepo:   repository.NewUserRepo(),
		logRepo:    repository.NewSystemLogRepo(),
	}
}

func (s *MemberService) GetBenefits() ([]model.MemberBenefit, error) {
	return s.memberRepo.ListBenefits()
}

func (s *MemberService) GetBenefit(level string) (*model.MemberBenefit, error) {
	return s.memberRepo.GetBenefit(level)
}

func (s *MemberService) CreateOrder(userID uint, level string, months int) (*model.MemberOrder, error) {
	benefit, err := s.memberRepo.GetBenefit(level)
	if err != nil {
		return nil, errors.New("无效的会员等级")
	}

	amount := benefit.PricePerMonth * float64(months)

	order := &model.MemberOrder{
		UserID: userID,
		Level:  model.MemberLevel(level),
		Months: months,
		Amount: amount,
		Status: "pending",
	}

	if err := s.memberRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *MemberService) PayOrder(orderID uint) error {
	order, err := s.getOrder(orderID)
	if err != nil {
		return err
	}
	if order.Status != "pending" {
		return errors.New("订单状态异常")
	}

	now := time.Now()
	user, _ := s.userRepo.FindByID(order.UserID)

	var expireAt *time.Time
	if user.MemberExpire != nil && user.MemberExpire.After(now) {
		t := user.MemberExpire.AddDate(0, order.Months, 0)
		expireAt = &t
	} else {
		t := now.AddDate(0, order.Months, 0)
		expireAt = &t
	}

	return utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.MemberOrder{}).
			Where("id = ?", orderID).
			Updates(map[string]interface{}{
				"status":   "paid",
				"paid_at":  now,
				"expire_at": expireAt,
			}).Error; err != nil {
			return err
		}
		return tx.Model(&model.User{}).
			Where("id = ?", order.UserID).
			Updates(map[string]interface{}{
				"member_level":  order.Level,
				"member_expire": expireAt,
			}).Error
	})
}

func (s *MemberService) CheckAndDowngradeExpired() error {
	expiredUsers, err := s.memberRepo.GetExpiredMembers()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, user := range expiredUsers {
		if user.MemberExpire != nil && user.MemberExpire.Before(now) {
			s.userRepo.Update(user.ID, map[string]interface{}{
				"member_level":  model.MemberFree,
				"member_expire": nil,
			})
			s.logRepo.Create(&model.SystemLog{
				UserID: user.ID,
				Module: "member",
				Action: "auto_downgrade",
			})
		}
	}
	return nil
}

func (s *MemberService) GetUserOrders(userID uint, page, pageSize int) ([]model.MemberOrder, int64, error) {
	return s.memberRepo.ListOrders(userID, page, pageSize)
}

func (s *MemberService) getOrder(id uint) (*model.MemberOrder, error) {
	var order model.MemberOrder
	err := utils.DB.First(&order, id).Error
	return &order, err
}

func (s *MemberService) CountTodayInteracts(userID uint) (int64, error) {
	return s.memberRepo.CountTodayInteracts(userID)
}
