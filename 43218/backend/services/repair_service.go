package services

import (
	"errors"
	"time"

	"secondhand-platform/database"
	"secondhand-platform/models"
	"secondhand-platform/utils"

	"github.com/sirupsen/logrus"
)

type RepairService struct{}

func NewRepairService() *RepairService {
	return &RepairService{}
}

func (s *RepairService) CreateService(technicianID uint, serviceType, title, description string, price, minPrice, maxPrice float64, estimatedDays int, images string) (*models.RepairService, error) {
	var technician models.User
	if err := database.DB.First(&technician, technicianID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if technician.Role != models.RoleTechnician {
		return nil, errors.New("只有维修技师可以发布维修服务")
	}

	if !technician.IsAuthenticated {
		return nil, errors.New("维修技师需要先通过资质审核")
	}

	service := &models.RepairService{
		TechnicianID:  technicianID,
		ServiceType:   serviceType,
		Title:         title,
		Description:   description,
		Price:         price,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		EstimatedDays: estimatedDays,
		Images:        images,
		Status:        models.ServiceStatusActive,
	}

	result := database.DB.Create(service)
	if result.Error != nil {
		return nil, result.Error
	}

	return service, nil
}

func (s *RepairService) GetServiceByID(id uint) (*models.RepairService, error) {
	var service models.RepairService
	if err := database.DB.Preload("Technician").First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *RepairService) ListServices(page, pageSize int, serviceType, keyword string, sortBy string, technicianID uint) ([]models.RepairService, int64, error) {
	var services []models.RepairService
	var total int64

	db := database.DB.Model(&models.RepairService{}).Where("status = ?", models.ServiceStatusActive)

	if serviceType != "" {
		db = db.Where("service_type = ?", serviceType)
	}
	if keyword != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if technicianID > 0 {
		db = db.Where("technician_id = ?", technicianID)
	}

	switch sortBy {
	case "price_asc":
		db = db.Order("price ASC")
	case "price_desc":
		db = db.Order("price DESC")
	case "rating":
		db = db.Order("rating DESC")
	case "orders":
		db = db.Order("order_count DESC")
	default:
		db = db.Order("created_at DESC")
	}

	db.Count(&total)
	if err := db.Preload("Technician").Offset((page - 1) * pageSize).Limit(pageSize).Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (s *RepairService) UpdateService(userID, serviceID uint, updates map[string]interface{}) error {
	var service models.RepairService
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		return errors.New("服务不存在")
	}

	if service.TechnicianID != userID {
		return errors.New("无权修改此服务")
	}

	return database.DB.Model(&service).Updates(updates).Error
}

func (s *RepairService) DeleteService(userID, serviceID uint) error {
	var service models.RepairService
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		return errors.New("服务不存在")
	}

	if service.TechnicianID != userID {
		return errors.New("无权删除此服务")
	}

	return database.DB.Delete(&service).Error
}

func (s *RepairService) GetServiceTypes() []string {
	return models.ServiceTypes
}

func (s *RepairService) CreateRepairOrder(buyerID, technicianID, serviceID uint, deviceType, deviceBrand, deviceModel, faultDescription, contactName, contactPhone, address string, servicePrice float64) (*models.RepairOrder, error) {
	var service models.RepairService
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		return nil, errors.New("服务不存在")
	}

	warrantyDays := 90
	warrantyUntil := time.Now().AddDate(0, 0, warrantyDays)

	order := &models.RepairOrder{
		OrderNo:          utils.GenerateOrderNo(),
		BuyerID:          buyerID,
		TechnicianID:     technicianID,
		ServiceID:        serviceID,
		DeviceType:       deviceType,
		DeviceBrand:      deviceBrand,
		DeviceModel:      deviceModel,
		FaultDescription: faultDescription,
		ContactName:      contactName,
		ContactPhone:     contactPhone,
		Address:          address,
		ServicePrice:     servicePrice,
		FinalPrice:       servicePrice,
		Status:           models.RepairStatusPending,
		WarrantyDays:     warrantyDays,
		WarrantyUntil:    &warrantyUntil,
	}

	result := database.DB.Create(order)
	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (s *RepairService) GetRepairOrderByID(id uint) (*models.RepairOrder, error) {
	var order models.RepairOrder
	if err := database.DB.Preload("Buyer").Preload("Technician").Preload("Service").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *RepairService) GetRepairOrderByOrderNo(orderNo string) (*models.RepairOrder, error) {
	var order models.RepairOrder
	if err := database.DB.Preload("Buyer").Preload("Technician").Preload("Service").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *RepairService) ListRepairOrders(page, pageSize int, userID uint, userRole string, status int) ([]models.RepairOrder, int64, error) {
	var orders []models.RepairOrder
	var total int64

	db := database.DB.Model(&models.RepairOrder{})

	if userRole == models.RoleBuyer {
		db = db.Where("buyer_id = ?", userID)
	} else if userRole == models.RoleTechnician {
		db = db.Where("technician_id = ?", userID)
	}

	if status > 0 {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)
	if err := db.Preload("Buyer").Preload("Technician").Preload("Service").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *RepairService) AcceptRepairOrder(technicianID, orderID uint) error {
	var order models.RepairOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.TechnicianID != technicianID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.RepairStatusPending {
		return errors.New("订单状态不允许接单")
	}

	now := time.Now()
	order.Status = models.RepairStatusAccepted
	order.AcceptedAt = &now

	return database.DB.Save(&order).Error
}

func (s *RepairService) StartRepair(technicianID, orderID uint) error {
	var order models.RepairOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.TechnicianID != technicianID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.RepairStatusAccepted {
		return errors.New("订单状态不允许开始维修")
	}

	order.Status = models.RepairStatusRepairing
	return database.DB.Save(&order).Error
}

func (s *RepairService) CompleteRepair(technicianID, orderID uint, finalPrice float64) error {
	var order models.RepairOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.TechnicianID != technicianID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.RepairStatusRepairing {
		return errors.New("订单状态不允许完成维修")
	}

	now := time.Now()
	order.Status = models.RepairStatusCompleted
	order.CompletedAt = &now
	if finalPrice > 0 {
		order.FinalPrice = finalPrice
	}

	database.DB.Model(&models.RepairService{}).Where("id = ?", order.ServiceID).
		UpdateColumn("order_count", gorm.Expr("order_count + 1"))

	return database.DB.Save(&order).Error
}

func (s *RepairService) PickUpDevice(userID, orderID uint) error {
	var order models.RepairOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.RepairStatusCompleted {
		return errors.New("订单状态不允许取件")
	}

	now := time.Now()
	order.Status = models.RepairStatusPickedUp
	order.PickedUpAt = &now

	return database.DB.Save(&order).Error
}

func (s *RepairService) CancelRepairOrder(userID, orderID uint, userRole string) error {
	var order models.RepairOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if userRole == models.RoleBuyer && order.BuyerID != userID {
		return errors.New("无权操作此订单")
	}
	if userRole == models.RoleTechnician && order.TechnicianID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status > models.RepairStatusAccepted {
		return errors.New("订单状态不允许取消")
	}

	order.Status = models.RepairStatusCancelled
	return database.DB.Save(&order).Error
}

func (s *RepairService) UpdateRepairOrderStatus(orderID uint, status int) error {
	var order models.RepairOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	order.Status = status
	return database.DB.Save(&order).Error
}

func init() {
	logrus.Info("Repair service initialized")
}
