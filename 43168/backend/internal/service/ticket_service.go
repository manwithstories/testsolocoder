package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
)

// TicketService 售后工单业务逻辑层
type TicketService struct {
	ticketRepo *repository.TicketRepository
	orderRepo  *repository.OrderRepository
	db         *gorm.DB
}

// NewTicketService 创建售后工单服务
func NewTicketService(
	ticketRepo *repository.TicketRepository,
	orderRepo *repository.OrderRepository,
	db *gorm.DB,
) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		orderRepo:  orderRepo,
		db:         db,
	}
}

// CreateTicket 业主创建售后工单
func (s *TicketService) CreateTicket(ownerID uint, req *dto.CreateTicketRequest) (*model.Ticket, error) {
	if !model.ValidTicketType(req.Type) {
		return nil, errors.New("工单类型不合法")
	}

	// 校验订单存在且属于当前业主
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}
	if order.OwnerID != ownerID {
		return nil, errors.New("无权为他人的订单创建工单")
	}

	ticket := &model.Ticket{
		OrderID: req.OrderID,
		OwnerID: ownerID,
		Type:    req.Type,
		Title:   req.Title,
		Content: req.Content,
		Status:  model.TicketStatusOpen,
	}
	if err := ticket.SetImages(req.Images); err != nil {
		return nil, err
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

// UpdateTicket 更新工单（状态流转/内容修改）
func (s *TicketService) UpdateTicket(id uint, userID uint, role string, req *dto.UpdateTicketRequest) (*model.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, errors.New("工单不存在")
	}

	// 权限校验
	if role == model.RoleOwner && ticket.OwnerID != userID {
		return nil, errors.New("无权操作他人的工单")
	}

	if req.Type != "" {
		ticket.Type = req.Type
	}
	if req.Title != "" {
		ticket.Title = req.Title
	}
	if req.Content != "" {
		ticket.Content = req.Content
	}
	if req.Images != nil {
		if err := ticket.SetImages(req.Images); err != nil {
			return nil, err
		}
	}
	if req.Status != "" {
		if !model.ValidTicketStatus(req.Status) {
			return nil, errors.New("工单状态不合法")
		}
		// 业主只能关闭自己的工单；厂商/管理员可处理、解决、关闭
		switch req.Status {
		case model.TicketStatusProcessing, model.TicketStatusResolved:
			if role == model.RoleOwner {
				return nil, errors.New("业主无权执行该状态变更")
			}
		case model.TicketStatusClosed:
			if role == model.RoleOwner && ticket.OwnerID != userID {
				return nil, errors.New("无权关闭他人的工单")
			}
		}
		ticket.Status = req.Status
	}
	ticket.UpdatedAt = time.Now()

	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

// GetByID 根据 ID 获取工单详情
func (s *TicketService) GetByID(id uint, userID uint, role string) (*model.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, errors.New("工单不存在")
	}

	if role == model.RoleOwner && ticket.OwnerID != userID {
		return nil, errors.New("无权查看他人的工单")
	}
	return ticket, nil
}

// List 分页查询工单列表
func (s *TicketService) List(params *dto.TicketListRequest) ([]*model.Ticket, int64, error) {
	p := &repository.TicketListParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		OrderID:  params.OrderID,
		Type:     params.Type,
		Status:   params.Status,
		Role:     params.Role,
		UserID:   params.UserID,
	}
	return s.ticketRepo.List(p)
}
