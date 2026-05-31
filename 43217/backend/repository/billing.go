package repository

import (
	"health-platform/models"
	"time"

	"gorm.io/gorm"
)

type BillingRepository struct {
	*BaseRepository
}

func NewBillingRepository() *BillingRepository {
	return &BillingRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *BillingRepository) FindByCompanyID(companyID uint, page, pageSize int) ([]models.Billing, int64, error) {
	var billings []models.Billing
	var total int64

	query := r.DB.Model(&models.Billing{}).Where("company_id = ?", companyID)
	query.Count(&total)

	err := query.Preload("Agency").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&billings).Error
	return billings, total, err
}

func (r *BillingRepository) FindByAgencyID(agencyID uint, page, pageSize int) ([]models.Billing, int64, error) {
	var billings []models.Billing
	var total int64

	query := r.DB.Model(&models.Billing{}).Where("agency_id = ?", agencyID)
	query.Count(&total)

	err := query.Preload("Company").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&billings).Error
	return billings, total, err
}

func (r *BillingRepository) GetWithItems(billingID uint) (*models.Billing, error) {
	var billing models.Billing
	err := r.DB.Preload("Items").Preload("Company").Preload("Agency").
		First(&billing, billingID).Error
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *BillingRepository) GetByPeriod(companyID, agencyID uint, period string) (*models.Billing, error) {
	var billing models.Billing
	err := r.DB.Where("company_id = ? AND agency_id = ? AND period = ?",
		companyID, agencyID, period).
		First(&billing).Error
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *BillingRepository) UpdateStatus(billingID uint, status models.BillingStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == models.BillingStatusPaid {
		updates["paid_at"] = time.Now()
	}
	return r.DB.Model(&models.Billing{}).Where("id = ?", billingID).Updates(updates).Error
}

func (r *BillingRepository) UpdatePaidAmount(billingID uint, amount float64) error {
	return r.DB.Model(&models.Billing{}).Where("id = ?", billingID).
		Update("paid_amount", gorm.Expr("paid_amount + ?", amount)).Error
}

func (r *BillingRepository) UpdateInvoice(billingID uint, invoiceNo, invoiceFile string) error {
	return r.DB.Model(&models.Billing{}).Where("id = ?", billingID).
		Updates(map[string]interface{}{
			"invoice_no":   invoiceNo,
			"invoice_file": invoiceFile,
			"invoice_status": 1,
		}).Error
}

func (r *BillingRepository) GetBillingByNo(billingNo string) (*models.Billing, error) {
	var billing models.Billing
	err := r.DB.Where("billing_no = ?", billingNo).First(&billing).Error
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

type BillingItemRepository struct {
	*BaseRepository
}

func NewBillingItemRepository() *BillingItemRepository {
	return &BillingItemRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *BillingItemRepository) FindByBillingID(billingID uint) ([]models.BillingItem, error) {
	var items []models.BillingItem
	err := r.DB.Where("billing_id = ?", billingID).Order("id ASC").Find(&items).Error
	return items, err
}

type TransactionRepository struct {
	*BaseRepository
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *TransactionRepository) FindByCompanyID(companyID uint, page, pageSize int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.DB.Model(&models.Transaction{}).Where("company_id = ?", companyID)
	query.Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&transactions).Error
	return transactions, total, err
}

func (r *TransactionRepository) GetTransactionByNo(transactionNo string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.DB.Where("transaction_no = ?", transactionNo).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

type CompanyBudgetRepository struct {
	*BaseRepository
}

func NewCompanyBudgetRepository() *CompanyBudgetRepository {
	return &CompanyBudgetRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *CompanyBudgetRepository) FindByCompanyIDAndYear(companyID uint, year int) (*models.CompanyBudget, error) {
	var budget models.CompanyBudget
	err := r.DB.Where("company_id = ? AND year = ?", companyID, year).First(&budget).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *CompanyBudgetRepository) UpdateUsedBudget(id uint, amount float64) error {
	return r.DB.Model(&models.CompanyBudget{}).Where("id = ?", id).
		Update("used_budget", gorm.Expr("used_budget + ?", amount)).Error
}

type DepartmentAppointmentRepository struct {
	*BaseRepository
}

func NewDepartmentAppointmentRepository() *DepartmentAppointmentRepository {
	return &DepartmentAppointmentRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *DepartmentAppointmentRepository) FindByDepartmentID(departmentID uint, year int) ([]models.DepartmentAppointment, error) {
	var appointments []models.DepartmentAppointment
	err := r.DB.Where("department_id = ? AND year = ?", departmentID, year).
		Preload("Agency").Preload("Package").
		Find(&appointments).Error
	return appointments, err
}

func (r *DepartmentAppointmentRepository) UpdateUsedQuota(id uint, count int) error {
	return r.DB.Model(&models.DepartmentAppointment{}).Where("id = ?", id).
		Update("used_quota", gorm.Expr("used_quota + ?", count)).Error
}
