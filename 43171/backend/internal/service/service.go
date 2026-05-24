package service

import (
	"drone-rental/internal/config"
	"drone-rental/internal/dto"
	"drone-rental/internal/model"
	"drone-rental/internal/pkg/utils"
	"drone-rental/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type ServiceService struct {
	serviceRepo *repository.ServiceRepo
	bidRepo     *repository.BidRepo
	userRepo    *repository.UserRepo
}

func NewServiceService() *ServiceService {
	return &ServiceService{
		serviceRepo: repository.NewServiceRepo(),
		bidRepo:     repository.NewBidRepo(),
		userRepo:    repository.NewUserRepo(),
	}
}

func (s *ServiceService) Create(userID uint, req *dto.CreateServiceReq) (*model.AerialService, error) {
	service := &model.AerialService{
		ServiceNo:   utils.GenerateOrderNo("SV"),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Region:      req.Region,
		Address:     req.Address,
		ServiceDate: req.ServiceDate,
		ServiceTime: req.ServiceTime,
		Duration:    req.Duration,
		BudgetMin:   req.BudgetMin,
		BudgetMax:   req.BudgetMax,
		Status:      model.ServiceStatusOpen,
		Images:      req.Images,
		Remark:      req.Remark,
	}
	if err := s.serviceRepo.Create(service); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceService) GetByID(id uint) (*model.AerialService, error) {
	return s.serviceRepo.GetByID(id)
}

func (s *ServiceService) List(page, pageSize int, userID, pilotID uint, status model.ServiceStatus, region string) ([]model.AerialService, int64, error) {
	return s.serviceRepo.List(page, pageSize, userID, pilotID, status, region)
}

func (s *ServiceService) CreateBid(pilotID uint, req *dto.CreateBidReq) error {
	service, err := s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		return errors.New("服务需求不存在")
	}
	if service.Status != model.ServiceStatusOpen {
		return errors.New("该需求已关闭或已分配")
	}
	pilot, err := s.userRepo.GetByID(pilotID)
	if err != nil {
		return errors.New("飞手不存在")
	}
	if pilot.VerifyStatus != model.VerifyApproved {
		return errors.New("飞手资质尚未通过认证")
	}
	bid := &model.ServiceBid{
		ServiceID: req.ServiceID,
		PilotID:   pilotID,
		Price:     req.Price,
		Message:   req.Message,
		Status:    model.BidStatusPending,
	}
	return s.bidRepo.Create(bid)
}

func (s *ServiceService) ListBids(serviceID uint) ([]model.ServiceBid, error) {
	return s.bidRepo.ListByServiceID(serviceID)
}

func (s *ServiceService) AcceptBid(userID uint, req *dto.AcceptBidReq) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		bid, err := s.bidRepo.GetByID(req.BidID)
		if err != nil {
			return errors.New("竞标不存在")
		}
		service, err := s.serviceRepo.GetByID(bid.ServiceID)
		if err != nil {
			return errors.New("服务需求不存在")
		}
		if service.UserID != userID {
			return errors.New("无权操作该需求")
		}
		if service.Status != model.ServiceStatusOpen {
			return errors.New("该需求已关闭或已分配")
		}
		bid.Status = model.BidStatusAccepted
		if err := tx.Save(bid).Error; err != nil {
			return err
		}
		if err := s.bidRepo.RejectOtherBids(service.ID, bid.ID); err != nil {
			return err
		}
		service.Status = model.ServiceStatusAssigned
		service.PilotID = &bid.PilotID
		service.FinalPrice = bid.Price
		return tx.Save(service).Error
	})
}

func (s *ServiceService) UpdateStatus(userID uint, req *dto.UpdateServiceStatusReq) error {
	service, err := s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		return errors.New("服务需求不存在")
	}
	if service.UserID != userID && (service.PilotID == nil || *service.PilotID != userID) {
		return errors.New("无权操作该需求")
	}
	return s.serviceRepo.UpdateStatus(req.ServiceID, model.ServiceStatus(req.Status))
}
