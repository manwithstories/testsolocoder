package repository

import (
	"drone-rental/internal/config"
	"drone-rental/internal/model"
)

type ServiceRepo struct{}

func NewServiceRepo() *ServiceRepo {
	return &ServiceRepo{}
}

func (r *ServiceRepo) Create(service *model.AerialService) error {
	return config.DB.Create(service).Error
}

func (r *ServiceRepo) GetByID(id uint) (*model.AerialService, error) {
	var service model.AerialService
	err := config.DB.Preload("User").Preload("Pilot").First(&service, id).Error
	return &service, err
}

func (r *ServiceRepo) Update(service *model.AerialService) error {
	return config.DB.Save(service).Error
}

func (r *ServiceRepo) UpdateStatus(id uint, status model.ServiceStatus) error {
	return config.DB.Model(&model.AerialService{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *ServiceRepo) List(page, pageSize int, userID, pilotID uint, status model.ServiceStatus, region string) ([]model.AerialService, int64, error) {
	var services []model.AerialService
	var total int64
	db := config.DB.Model(&model.AerialService{}).Preload("User").Preload("Pilot")
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	if pilotID != 0 {
		db = db.Where("pilot_id = ?", pilotID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if region != "" {
		db = db.Where("region = ?", region)
	}
	db.Count(&total)
	err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&services).Error
	return services, total, err
}

type BidRepo struct{}

func NewBidRepo() *BidRepo {
	return &BidRepo{}
}

func (r *BidRepo) Create(bid *model.ServiceBid) error {
	return config.DB.Create(bid).Error
}

func (r *BidRepo) GetByID(id uint) (*model.ServiceBid, error) {
	var bid model.ServiceBid
	err := config.DB.Preload("Pilot").First(&bid, id).Error
	return &bid, err
}

func (r *BidRepo) Update(bid *model.ServiceBid) error {
	return config.DB.Save(bid).Error
}

func (r *BidRepo) ListByServiceID(serviceID uint) ([]model.ServiceBid, error) {
	var bids []model.ServiceBid
	err := config.DB.Preload("Pilot").Where("service_id = ?", serviceID).
		Order("price ASC, id ASC").Find(&bids).Error
	return bids, err
}

func (r *BidRepo) RejectOtherBids(serviceID, acceptedBidID uint) error {
	return config.DB.Model(&model.ServiceBid{}).
		Where("service_id = ? AND id != ?", serviceID, acceptedBidID).
		Update("status", model.BidStatusRejected).Error
}
