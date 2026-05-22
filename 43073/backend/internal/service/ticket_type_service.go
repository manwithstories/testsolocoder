package service

import (
	"errors"
	"ticket-system/internal/dto"
	"ticket-system/internal/logger"
	"ticket-system/internal/models"
	"ticket-system/internal/redis"
)

type TicketTypeService struct{}

func NewTicketTypeService() *TicketTypeService {
	return &TicketTypeService{}
}

func (s *TicketTypeService) Create(req *dto.TicketTypeCreateRequest) (*models.TicketType, error) {
	var activity models.Activity
	if err := models.DB.First(&activity, req.ActivityID).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	ticketType := &models.TicketType{
		ActivityID: req.ActivityID,
		Name:       req.Name,
		Type:       req.Type,
		Price:      req.Price,
		Stock:      req.Stock,
		SoldCount:  0,
		Status:     models.TicketStatusOnSale,
	}

	if err := models.DB.Create(ticketType).Error; err != nil {
		logger.Log.Errorf("Create ticket type failed: %v", err)
		return nil, errors.New("创建票型失败")
	}

	if err := redis.SetStock(ticketType.ID, ticketType.Stock); err != nil {
		logger.Log.Warnf("Set redis stock failed: %v", err)
	}

	logger.Log.Infof("Ticket type created: %d - %s", ticketType.ID, ticketType.Name)
	return ticketType, nil
}

func (s *TicketTypeService) GetList(req *dto.TicketTypeListRequest) ([]models.TicketType, int64, error) {
	var ticketTypes []models.TicketType
	var total int64

	query := models.DB.Model(&models.TicketType{}).Preload("Activity")

	if req.ActivityID > 0 {
		query = query.Where("activity_id = ?", req.ActivityID)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&ticketTypes).Error; err != nil {
		return nil, 0, err
	}

	return ticketTypes, total, nil
}

func (s *TicketTypeService) GetByID(id uint) (*models.TicketType, error) {
	var ticketType models.TicketType
	if err := models.DB.Preload("Activity").First(&ticketType, id).Error; err != nil {
		return nil, err
	}

	redisStock, err := redis.GetStock(id)
	if err == nil {
		ticketType.Stock = int(redisStock)
	}

	return &ticketType, nil
}

func (s *TicketTypeService) Update(id uint, req *dto.TicketTypeUpdateRequest) (*models.TicketType, error) {
	ticketType, err := s.GetByID(id)
	if err != nil {
		return nil, errors.New("票型不存在")
	}

	if req.Name != "" {
		ticketType.Name = req.Name
	}
	if req.Type != "" {
		ticketType.Type = req.Type
	}
	if req.Price >= 0 {
		ticketType.Price = req.Price
	}
	if req.Status != "" {
		ticketType.Status = req.Status
	}

	if req.Stock > 0 {
		oldStock := ticketType.Stock
		ticketType.Stock = req.Stock
		diff := req.Stock - oldStock
		if diff != 0 {
			if err := redis.IncrementStock(id, diff); err != nil {
				logger.Log.Warnf("Update redis stock failed: %v", err)
			}
		}
	}

	if err := models.DB.Save(ticketType).Error; err != nil {
		return nil, err
	}

	logger.Log.Infof("Ticket type updated: %d", id)
	return ticketType, nil
}

func (s *TicketTypeService) Delete(id uint) error {
	if err := models.DB.Delete(&models.TicketType{}, id).Error; err != nil {
		return err
	}
	logger.Log.Infof("Ticket type deleted: %d", id)
	return nil
}

func (s *TicketTypeService) UpdateSoldStatus(id uint) error {
	ticketType, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if ticketType.Stock <= 0 && ticketType.Status == models.TicketStatusOnSale {
		ticketType.Status = models.TicketStatusSoldOut
		if err := models.DB.Save(ticketType).Error; err != nil {
			return err
		}
		logger.Log.Infof("Ticket type sold out: %d", id)
	}

	return nil
}
