package service

import (
	"fmt"
	"time"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
	"beauty-salon-system/internal/utils"
	"gorm.io/gorm"
)

type PaymentService struct {
	paymentRepo       *repository.PaymentRepository
	appointmentRepo   *repository.AppointmentRepository
	customerRepo      *repository.CustomerRepository
	customerPkgRepo   *repository.CustomerPackageRepository
	memberCardRepo    *repository.MemberCardRepository
	serviceItemService *ServiceItemService
}

func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	appointmentRepo *repository.AppointmentRepository,
	customerRepo *repository.CustomerRepository,
	customerPkgRepo *repository.CustomerPackageRepository,
	memberCardRepo *repository.MemberCardRepository,
	serviceItemService *ServiceItemService,
) *PaymentService {
	return &PaymentService{
		paymentRepo:       paymentRepo,
		appointmentRepo:   appointmentRepo,
		customerRepo:      customerRepo,
		customerPkgRepo:   customerPkgRepo,
		memberCardRepo:    memberCardRepo,
		serviceItemService: serviceItemService,
	}
}

type CreatePaymentRequest struct {
	AppointmentID uint    `json:"appointment_id" binding:"required"`
	CustomerID    uint    `json:"customer_id" binding:"required"`
	PayMethod     string  `json:"pay_method" binding:"required"`
	Amount        float64 `json:"amount"`
	PointsUsed    int     `json:"points_used"`
	CardID        *uint   `json:"card_id"`
	PackageID     *uint   `json:"package_id"`
}

func (s *PaymentService) CreatePayment(req *CreatePaymentRequest) (*model.Payment, error) {
	appointment, err := s.appointmentRepo.GetByID(req.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("appointment not found: %w", err)
	}

	customer, err := s.customerRepo.GetByID(req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	price, err := s.serviceItemService.CalculatePrice(appointment.ServiceID, appointment.AppointmentDate)
	if err != nil {
		return nil, fmt.Errorf("calculate price: %w", err)
	}

	actualAmount := price

	if req.PointsUsed > 0 {
		if customer.Points < req.PointsUsed {
			return nil, fmt.Errorf("insufficient points")
		}
		pointDeduction := float64(req.PointsUsed) / 100
		actualAmount -= pointDeduction
		if actualAmount < 0 {
			actualAmount = 0
		}
	}

	if req.PackageID != nil {
		customerPkg, err := s.customerPkgRepo.GetByID(*req.PackageID)
		if err != nil {
			return nil, fmt.Errorf("package not found: %w", err)
		}
		if customerPkg.UsedCount >= customerPkg.TotalCount {
			return nil, fmt.Errorf("package has been fully used")
		}
		if customerPkg.ExpireDate.Before(time.Now()) {
			return nil, fmt.Errorf("package has expired")
		}
		actualAmount = 0
	}

	if req.CardID != nil {
		card, err := s.memberCardRepo.GetByID(*req.CardID)
		if err != nil {
			return nil, fmt.Errorf("member card not found: %w", err)
		}
		actualAmount = actualAmount * card.Discount
	}

	tx := s.paymentRepo.BeginTransaction()
	if tx.Error != nil {
		return nil, fmt.Errorf("begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	payment := &model.Payment{
		AppointmentID: req.AppointmentID,
		CustomerID:    req.CustomerID,
		Amount:        actualAmount,
		PayMethod:     req.PayMethod,
		PointsUsed:    req.PointsUsed,
		CardID:        req.CardID,
		Status:        "success",
		TransactionNo: utils.GenerateRandomToken(),
	}

	if err := s.paymentRepo.Create(payment, tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create payment: %w", err)
	}

	if req.PointsUsed > 0 {
		if err := tx.Model(&model.Customer{}).Where("id = ?", req.CustomerID).
			UpdateColumn("points", gorm.Expr("points - ?", req.PointsUsed)).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("deduct points: %w", err)
		}
	}

	if req.PackageID != nil {
		if err := s.customerPkgRepo.IncrementUsedCount(*req.PackageID, tx); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("use package: %w", err)
		}
	}

	if req.CardID != nil && actualAmount > 0 {
		if err := s.memberCardRepo.DeductBalance(*req.CardID, actualAmount, tx); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("deduct card balance: %w", err)
		}
	}

	if err := tx.Model(&model.Appointment{}).Where("id = ?", req.AppointmentID).
		Update("status", "paid").Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("update appointment: %w", err)
	}

	pointsEarned := int(actualAmount / 10)
	if pointsEarned > 0 {
		if err := tx.Model(&model.Customer{}).Where("id = ?", req.CustomerID).
			UpdateColumn("points", gorm.Expr("points + ?", pointsEarned)).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("earn points: %w", err)
		}
	}

	if err := tx.Model(&model.Customer{}).Where("id = ?", req.CustomerID).
		UpdateColumn("total_spent", gorm.Expr("total_spent + ?", actualAmount)).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("update total spent: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return payment, nil
}

func (s *PaymentService) GetByID(id uint) (*model.Payment, error) {
	return s.paymentRepo.GetByID(id)
}

func (s *PaymentService) GetByAppointmentID(appointmentID uint) (*model.Payment, error) {
	return s.paymentRepo.GetByAppointmentID(appointmentID)
}

func (s *PaymentService) List(page, pageSize int, filters map[string]interface{}) ([]model.Payment, int64, error) {
	return s.paymentRepo.List(page, pageSize, filters)
}

func (s *PaymentService) GetRevenueByDateRange(startDate, endDate string) (float64, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	return s.paymentRepo.GetRevenueByDateRange(start, end)
}

func (s *PaymentService) GetRevenueByTechnician(technicianID uint, startDate, endDate string) (float64, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	return s.paymentRepo.GetRevenueByTechnician(technicianID, start, end)
}

func (s *PaymentService) GetRevenueByService(serviceID uint, startDate, endDate string) (float64, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	return s.paymentRepo.GetRevenueByService(serviceID, start, end)
}

type CreateMemberCardRequest struct {
	CustomerID uint    `json:"customer_id" binding:"required"`
	CardType   string  `json:"card_type"`
	Balance    float64 `json:"balance"`
	Discount   float64 `json:"discount"`
}

func (s *PaymentService) CreateMemberCard(req *CreateMemberCardRequest) (*model.MemberCard, error) {
	card := &model.MemberCard{
		CustomerID: req.CustomerID,
		CardNo:     fmt.Sprintf("MC%d%d", req.CustomerID, time.Now().Unix()),
		CardType:   req.CardType,
		Balance:    req.Balance,
		Discount:   req.Discount,
		Status:     1,
	}

	if card.Discount <= 0 || card.Discount > 1 {
		card.Discount = 1.0
	}

	if err := s.memberCardRepo.Create(card); err != nil {
		return nil, fmt.Errorf("create member card: %w", err)
	}

	return card, nil
}

func (s *PaymentService) GetMemberCards(customerID uint) ([]model.MemberCard, error) {
	return s.memberCardRepo.GetByCustomerID(customerID)
}

func (s *PaymentService) RechargeCard(cardID uint, amount float64) error {
	return s.memberCardRepo.AddBalance(cardID, amount)
}

type PurchasePackageRequest struct {
	CustomerID uint `json:"customer_id" binding:"required"`
	ServiceID  uint `json:"service_id" binding:"required"`
}

func (s *PaymentService) PurchasePackage(req *PurchasePackageRequest) (*model.CustomerPackage, error) {
	serviceItem, err := s.serviceItemService.GetByID(req.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	if !serviceItem.IsPackage {
		return nil, fmt.Errorf("this service is not a package")
	}

	customerPkg := &model.CustomerPackage{
		CustomerID:   req.CustomerID,
		ServiceID:    req.ServiceID,
		TotalCount:   serviceItem.PackageCount,
		UsedCount:    0,
		PurchaseDate: time.Now(),
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		Status:       1,
	}

	if err := s.customerPkgRepo.Create(customerPkg, nil); err != nil {
		return nil, fmt.Errorf("purchase package: %w", err)
	}

	return customerPkg, nil
}

func (s *PaymentService) GetCustomerPackages(customerID uint) ([]model.CustomerPackage, error) {
	return s.customerPkgRepo.GetByCustomerID(customerID)
}
