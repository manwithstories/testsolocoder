package service

import (
	"errors"

	"ticket-system/internal/dto"
	"ticket-system/internal/logging"
	"ticket-system/internal/models"
	"ticket-system/internal/repository"
)

type CheckinService struct {
	orderRepo *repository.OrderRepository
}

func NewCheckinService() *CheckinService {
	return &CheckinService{
		orderRepo: repository.NewOrderRepository(),
	}
}

func (s *CheckinService) Checkin(ticketNo string, operatorID uint64) (*dto.Response, error) {
	ticket, err := s.orderRepo.GetTicketByTicketNo(ticketNo)
	if err != nil {
		_ = s.logCheckin(0, ticketNo, operatorID, 0, "票码无效")
		return nil, errors.New("票码无效")
	}

	if ticket.Status == models.TicketStatusRefunded {
		_ = s.logCheckin(ticket.ID, ticketNo, operatorID, 0, "该票已退款")
		return nil, errors.New("该票已退款，无法入场")
	}

	order, err := s.orderRepo.GetByID(ticket.OrderID)
	if err != nil || order.Status != models.OrderStatusPaid {
		_ = s.logCheckin(ticket.ID, ticketNo, operatorID, 0, "订单未支付")
		return nil, errors.New("订单未支付")
	}

	if ticket.CheckedIn == 1 {
		_ = s.logCheckin(ticket.ID, ticketNo, operatorID, 0, "该票已入场")
		return nil, errors.New("该票已入场，请勿重复验票")
	}

	err = s.orderRepo.MarkTicketCheckedIn(ticket.ID)
	if err != nil {
		_ = s.logCheckin(ticket.ID, ticketNo, operatorID, 0, "验票失败")
		return nil, errors.New("验票失败，请重试")
	}

	_ = s.logCheckin(ticket.ID, ticketNo, operatorID, 1, "验票成功")

	logging.Infof("Ticket checked in: %s, operator=%d", ticketNo, operatorID)
	return &dto.Response{
		Code:    200,
		Message: "验票成功",
		Data: map[string]interface{}{
			"ticket_no":  ticket.TicketNo,
			"seat_info":  ticket.SeatInfo,
			"real_name":  ticket.RealName,
			"checkin_time": ticket.CheckinTime,
		},
	}, nil
}

func (s *CheckinService) logCheckin(ticketID uint64, ticketNo string, operatorID uint64, status int, message string) error {
	log := &models.CheckinLog{
		TicketID:   ticketID,
		TicketNo:   ticketNo,
		OperatorID: operatorID,
		Status:     status,
		Message:    message,
	}
	return s.orderRepo.CreateCheckinLog(log)
}

func (s *CheckinService) GetCheckinLogs(ticketNo string, page, pageSize int) ([]models.CheckinLog, int64, error) {
	var logs []models.CheckinLog
	var total int64

	_ = repository.NewOrderRepository()
	if ticketNo != "" {
		// 简化处理，实际应在repo中实现
	}

	return logs, total, nil
}
