package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
)

// DeliveryService 配送安装业务逻辑层
type DeliveryService struct {
	deliveryRepo *repository.DeliveryRepository
	orderRepo    *repository.OrderRepository
	db           *gorm.DB
}

// NewDeliveryService 创建配送安装服务
func NewDeliveryService(
	deliveryRepo *repository.DeliveryRepository,
	orderRepo *repository.OrderRepository,
	db *gorm.DB,
) *DeliveryService {
	return &DeliveryService{
		deliveryRepo: deliveryRepo,
		orderRepo:    orderRepo,
		db:           db,
	}
}

// CreateDelivery 业主创建配送安装预约
func (s *DeliveryService) CreateDelivery(ownerID uint, req *dto.CreateDeliveryRequest) (*model.Delivery, error) {
	if !model.ValidDeliveryType(req.Type) {
		return nil, errors.New("配送安装类型不合法")
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
		return nil, errors.New("无权为他人的订单创建预约")
	}

	delivery := &model.Delivery{
		OrderID:      order.ID,
		OwnerID:      ownerID,
		Type:         req.Type,
		TimeSlot:     req.TimeSlot,
		Address:      req.Address,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       model.DeliveryStatusPending,
		InstallerID:  req.InstallerID,
		Remark:       req.Remark,
	}

	if err := s.deliveryRepo.Create(delivery); err != nil {
		return nil, err
	}
	return delivery, nil
}

// UpdateDelivery 更新配送安装（预约信息/确认/完成/取消）
func (s *DeliveryService) UpdateDelivery(id uint, userID uint, role string, req *dto.UpdateDeliveryRequest) (*model.Delivery, error) {
	delivery, err := s.deliveryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if delivery == nil {
		return nil, errors.New("配送安装记录不存在")
	}

	// 权限校验
	switch role {
	case model.RoleOwner:
		if delivery.OwnerID != userID {
			return nil, errors.New("无权操作他人的配送安装记录")
		}
	case model.RoleManufacturer, model.RoleDesigner:
		// 厂商/安装人员可对所有记录进行确认/完成
	default:
		return nil, errors.New("无权执行该操作")
	}

	// 更新字段（仅更新非零字段）
	if req.Type != "" {
		delivery.Type = req.Type
	}
	if req.TimeSlot != "" {
		delivery.TimeSlot = req.TimeSlot
	}
	if req.Address != "" {
		delivery.Address = req.Address
	}
	if req.ContactName != "" {
		delivery.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		delivery.ContactPhone = req.ContactPhone
	}
	if req.InstallerID > 0 {
		delivery.InstallerID = req.InstallerID
	}
	if req.Status != "" {
		if !model.ValidDeliveryStatus(req.Status) {
			return nil, errors.New("状态不合法")
		}
		// 业主只能取消自己的预约；厂商/安装人员才能确认、完成
		switch req.Status {
		case model.DeliveryStatusCancelled:
			if role != model.RoleOwner && role != model.RoleManufacturer {
				return nil, errors.New("无权取消")
			}
		case model.DeliveryStatusConfirmed, model.DeliveryStatusCompleted:
			if role != model.RoleManufacturer && role != model.RoleDesigner {
				return nil, errors.New("无权确认或完成")
			}
		}
		delivery.Status = req.Status
	}
	if req.Remark != "" {
		delivery.Remark = req.Remark
	}
	delivery.UpdatedAt = time.Now()

	if err := s.deliveryRepo.Update(delivery); err != nil {
		return nil, err
	}
	return delivery, nil
}

// GetByID 根据 ID 获取配送安装详情
func (s *DeliveryService) GetByID(id uint, userID uint, role string) (*model.Delivery, error) {
	delivery, err := s.deliveryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if delivery == nil {
		return nil, errors.New("配送安装记录不存在")
	}

	// 权限校验：业主只能查看自己的
	if role == model.RoleOwner && delivery.OwnerID != userID {
		return nil, errors.New("无权查看他人的配送安装记录")
	}
	return delivery, nil
}

// List 分页查询配送安装列表
func (s *DeliveryService) List(params *dto.DeliveryListRequest) ([]*model.Delivery, int64, error) {
	p := &repository.DeliveryListParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		OrderID:  params.OrderID,
		Type:     params.Type,
		Status:   params.Status,
		Role:     params.Role,
		UserID:   params.UserID,
	}
	return s.deliveryRepo.List(p)
}
