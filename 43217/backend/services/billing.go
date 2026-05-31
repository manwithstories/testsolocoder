package services

import (
	"errors"
	"fmt"
	"time"

	"health-platform/models"
	"health-platform/repository"
	"health-platform/utils"
)

type BillingService struct {
	billingRepo     *repository.BillingRepository
	billingItemRepo *repository.BillingItemRepository
	transactionRepo *repository.TransactionRepository
	companyRepo     *repository.CompanyRepository
	appointmentRepo *repository.AppointmentRepository
	packageRepo     *repository.PackageRepository
}

func NewBillingService() *BillingService {
	return &BillingService{
		billingRepo:     repository.NewBillingRepository(),
		billingItemRepo: repository.NewBillingItemRepository(),
		transactionRepo: repository.NewTransactionRepository(),
		companyRepo:     repository.NewCompanyRepository(),
		appointmentRepo: repository.NewAppointmentRepository(),
		packageRepo:     repository.NewPackageRepository(),
	}
}

type GenerateBillingRequest struct {
	CompanyID uint   `json:"company_id" binding:"required"`
	AgencyID  uint   `json:"agency_id" binding:"required"`
	Period    string `json:"period" binding:"required"`
}

type PayBillingRequest struct {
	BillingID     uint   `json:"billing_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type RechargeRequest struct {
	CompanyID     uint   `json:"company_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,min=0"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

func (s *BillingService) GenerateMonthlyBilling(req *GenerateBillingRequest) (*models.Billing, error) {
	existing, _ := s.billingRepo.GetByPeriod(req.CompanyID, req.AgencyID, req.Period)
	if existing != nil {
		return nil, errors.New("该期账单已生成")
	}

	appointments, _, _ := s.appointmentRepo.FindByCompanyIDAndAgency(req.CompanyID, req.AgencyID, 1, 10000)
	_ = appointments

	billingNo := utils.GenerateOrderNo("BL")

	billing := &models.Billing{
		BillingNo: billingNo,
		CompanyID: req.CompanyID,
		AgencyID:  req.AgencyID,
		Period:    req.Period,
		Status:    models.BillingStatusPending,
	}

	var totalAmount float64
	var billingItems []models.BillingItem

	for _, appt := range appointments {
		if appt.Status != models.AppointmentStatusCompleted {
			continue
		}

		var pkg models.Package
		if err := s.packageRepo.FindByID(&pkg, appt.PackageID); err != nil {
			continue
		}

		item := models.BillingItem{
			AppointmentID:   appt.ID,
			EmployeeID:      appt.EmployeeID,
			PackageID:       appt.PackageID,
			PackageName:     pkg.Name,
			UnitPrice:       pkg.Price,
			Quantity:        1,
			Amount:          pkg.Price,
			AppointmentDate: appt.AppointmentDate,
		}
		billingItems = append(billingItems, item)
		totalAmount += pkg.Price
	}

	billing.TotalAmount = totalAmount

	if err := s.billingRepo.Create(billing); err != nil {
		return nil, fmt.Errorf("生成账单失败: %w", err)
	}

	for i := range billingItems {
		billingItems[i].BillingID = billing.ID
		if err := s.billingItemRepo.Create(&billingItems[i]); err != nil {
			return nil, fmt.Errorf("添加账单明细失败: %w", err)
		}
	}

	return billing, nil
}

func (s *BillingService) GetBilling(billingID uint) (*models.Billing, error) {
	return s.billingRepo.GetWithItems(billingID)
}

func (s *BillingService) GetCompanyBillings(companyID uint, page, pageSize int) ([]models.Billing, int64, error) {
	return s.billingRepo.FindByCompanyID(companyID, page, pageSize)
}

func (s *BillingService) GetAgencyBillings(agencyID uint, page, pageSize int) ([]models.Billing, int64, error) {
	return s.billingRepo.FindByAgencyID(agencyID, page, pageSize)
}

func (s *BillingService) PayBilling(req *PayBillingRequest) error {
	var billing models.Billing
	if err := s.billingRepo.FindByID(&billing, req.BillingID); err != nil {
		return errors.New("账单不存在")
	}

	if billing.Status == models.BillingStatusPaid {
		return errors.New("账单已支付")
	}

	var company models.Company
	if err := s.companyRepo.FindByID(&company, billing.CompanyID); err != nil {
		return errors.New("企业不存在")
	}

	if company.PaymentType == 1 {
		if company.Balance < billing.TotalAmount {
			return errors.New("余额不足，请先充值")
		}
		s.companyRepo.UpdateBalance(billing.CompanyID, -billing.TotalAmount)
	}

	s.billingRepo.UpdateStatus(billing.ID, models.BillingStatusPaid)
	s.billingRepo.UpdatePaidAmount(billing.ID, billing.TotalAmount)

	transactionNo := utils.GenerateOrderNo("TX")
	transaction := &models.Transaction{
		TransactionNo: transactionNo,
		CompanyID:     billing.CompanyID,
		Type:          "payment",
		Amount:        -billing.TotalAmount,
		Balance:       company.Balance,
		PaymentMethod: req.PaymentMethod,
		Status:        1,
		Remark:        fmt.Sprintf("支付账单%s", billing.BillingNo),
	}
	return s.transactionRepo.Create(transaction)
}

func (s *BillingService) Recharge(req *RechargeRequest) error {
	var company models.Company
	if err := s.companyRepo.FindByID(&company, req.CompanyID); err != nil {
		return errors.New("企业不存在")
	}

	s.companyRepo.UpdateBalance(req.CompanyID, req.Amount)

	transactionNo := utils.GenerateOrderNo("TX")
	transaction := &models.Transaction{
		TransactionNo: transactionNo,
		CompanyID:     req.CompanyID,
		Type:          "recharge",
		Amount:        req.Amount,
		Balance:       company.Balance + req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        1,
		Remark:        "账户充值",
	}
	return s.transactionRepo.Create(transaction)
}

func (s *BillingService) GetTransactions(companyID uint, page, pageSize int) ([]models.Transaction, int64, error) {
	return s.transactionRepo.FindByCompanyID(companyID, page, pageSize)
}

func (s *BillingService) GetCompanyBalance(companyID uint) (map[string]interface{}, error) {
	var company models.Company
	if err := s.companyRepo.FindByID(&company, companyID); err != nil {
		return nil, errors.New("企业不存在")
	}

	return map[string]interface{}{
		"company_id":    companyID,
		"balance":       company.Balance,
		"credit_limit":  company.CreditLimit,
		"payment_type":  company.PaymentType,
		"annual_budget": company.AnnualBudget,
		"used_budget":   company.UsedBudget,
	}, nil
}
