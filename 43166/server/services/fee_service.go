package services

import (
	"errors"
	"fmt"
	"time"

	"business-registration-platform/database"
	"business-registration-platform/models"

	"github.com/google/uuid"
)

type FeeService struct{}

func NewFeeService() *FeeService {
	return &FeeService{}
}

type CalculateFeeRequest struct {
	ApplicationID uint              `json:"applicationId"`
	CompanyType   models.CompanyType `json:"companyType"`
	Capital       float64           `json:"capital"`
	DiscountCode  string            `json:"discountCode"`
}

type PayFeeRequest struct {
	ApplicationID uint   `json:"applicationId"`
	PaymentMethod string `json:"paymentMethod"`
}

func (s *FeeService) CalculateFee(req *CalculateFeeRequest) (*models.ApplicationFee, error) {
	var standard models.FeeStandard
	if err := database.DB.Where("company_type = ?", req.CompanyType).First(&standard).Error; err != nil {
		return nil, errors.New("fee standard not found for this company type")
	}

	feeItems := []models.FeeItem{
		{ItemName: "核名费", Amount: standard.NamingFee, Description: "公司名称核准费用"},
		{ItemName: "工商登记费", Amount: standard.RegistrationFee, Description: "工商注册登记费用"},
		{ItemName: "税务登记费", Amount: standard.TaxFee, Description: "税务登记相关费用"},
		{ItemName: "银行开户费", Amount: standard.BankFee, Description: "银行账户开立费用"},
		{ItemName: "刻章费", Amount: standard.SealFee, Description: "公章刻制备案费用"},
		{ItemName: "服务费", Amount: standard.ServiceFee, Description: "代办服务费用"},
	}

	if standard.CapitalRate > 0 && req.Capital > 0 {
		capitalFee := req.Capital * standard.CapitalRate / 10000
		feeItems = append(feeItems, models.FeeItem{
			ItemName:    "注册资本印花税",
			Amount:      capitalFee,
			Description: "按注册资本比例计算",
		})
	}

	var totalAmount float64
	for _, item := range feeItems {
		totalAmount += item.Amount
	}

	discountAmount := float64(0)
	if req.DiscountCode != "" {
		var policy models.DiscountPolicy
		if err := database.DB.Where("code = ? AND is_active = ?", req.DiscountCode, true).First(&policy).Error; err == nil {
			now := time.Now()
			if (policy.StartDate == nil || now.After(*policy.StartDate)) &&
				(policy.EndDate == nil || now.Before(*policy.EndDate)) {
				if totalAmount >= policy.MinAmount {
					if policy.Type == "fixed" {
						discountAmount = policy.Value
					} else if policy.Type == "percent" {
						discountAmount = totalAmount * policy.Value / 100
					}
					if policy.MaxDiscount > 0 && discountAmount > policy.MaxDiscount {
						discountAmount = policy.MaxDiscount
					}
				}
			}
		}
	}

	paidAmount := totalAmount - discountAmount
	if paidAmount < 0 {
		paidAmount = 0
	}

	return &models.ApplicationFee{
		ApplicationID:  req.ApplicationID,
		TotalAmount:    totalAmount,
		DiscountAmount: discountAmount,
		PaidAmount:     paidAmount,
		Status:         models.FeeStatusPending,
		FeeItems:       feeItems,
	}, nil
}

func (s *FeeService) CreateApplicationFee(applicationID uint, companyType models.CompanyType, capital float64, discountCode string) (*models.ApplicationFee, error) {
	calcReq := &CalculateFeeRequest{
		ApplicationID: applicationID,
		CompanyType:   companyType,
		Capital:       capital,
		DiscountCode:  discountCode,
	}

	fee, err := s.CalculateFee(calcReq)
	if err != nil {
		return nil, err
	}

	tx := database.DB.Begin()
	if err := tx.Create(fee).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range fee.FeeItems {
		fee.FeeItems[i].ApplicationFeeID = fee.ID
	}
	if err := tx.Create(&fee.FeeItems).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return fee, nil
}

func (s *FeeService) PayFee(req *PayFeeRequest) (*models.ApplicationFee, error) {
	var fee models.ApplicationFee
	if err := database.DB.Where("application_id = ?", req.ApplicationID).First(&fee).Error; err != nil {
		return nil, errors.New("fee record not found")
	}

	if fee.Status == models.FeeStatusPaid {
		return nil, errors.New("fee already paid")
	}

	transactionNo := fmt.Sprintf("PAY%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8])
	now := time.Now()

	fee.Status = models.FeeStatusPaid
	fee.PaymentMethod = req.PaymentMethod
	fee.PaymentTime = &now
	fee.TransactionNo = transactionNo

	if err := database.DB.Save(&fee).Error; err != nil {
		return nil, err
	}

	var application models.Application
	if err := database.DB.First(&application, req.ApplicationID).Error; err == nil {
		application.Status = models.AppStatusPendingReview
		database.DB.Save(&application)
	}

	return &fee, nil
}

func (s *FeeService) GetApplicationFee(applicationID uint) (*models.ApplicationFee, error) {
	var fee models.ApplicationFee
	if err := database.DB.Where("application_id = ?", applicationID).
		Preload("FeeItems").
		First(&fee).Error; err != nil {
		return nil, err
	}
	return &fee, nil
}

func (s *FeeService) GetFeeList(page, pageSize int, status string) ([]models.ApplicationFee, int64, error) {
	var fees []models.ApplicationFee
	var total int64

	db := database.DB.Model(&models.ApplicationFee{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Application").Preload("FeeItems").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&fees).Error; err != nil {
		return nil, 0, err
	}

	return fees, total, nil
}

func (s *FeeService) GetFeeStandards() ([]models.FeeStandard, error) {
	var standards []models.FeeStandard
	if err := database.DB.Find(&standards).Error; err != nil {
		return nil, err
	}
	return standards, nil
}

func (s *FeeService) UpdateFeeStandard(id uint, data map[string]interface{}) error {
	return database.DB.Model(&models.FeeStandard{}).Where("id = ?", id).Updates(data).Error
}

func (s *FeeService) CreateFeeStandard(standard *models.FeeStandard) error {
	return database.DB.Create(standard).Error
}

func (s *FeeService) GetDiscountPolicies() ([]models.DiscountPolicy, error) {
	var policies []models.DiscountPolicy
	if err := database.DB.Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}

func (s *FeeService) CreateDiscountPolicy(policy *models.DiscountPolicy) error {
	return database.DB.Create(policy).Error
}

func (s *FeeService) UpdateDiscountPolicy(id uint, data map[string]interface{}) error {
	return database.DB.Model(&models.DiscountPolicy{}).Where("id = ?", id).Updates(data).Error
}

func (s *FeeService) DeleteDiscountPolicy(id uint) error {
	return database.DB.Delete(&models.DiscountPolicy{}, id).Error
}
