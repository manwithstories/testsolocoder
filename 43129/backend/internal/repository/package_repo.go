package repository

import (
	"beauty-salon-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type CustomerPackageRepository struct {
	db *gorm.DB
}

func NewCustomerPackageRepository(db *gorm.DB) *CustomerPackageRepository {
	return &CustomerPackageRepository{db: db}
}

func (r *CustomerPackageRepository) Create(pkg *model.CustomerPackage, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(pkg).Error
	}
	return r.db.Create(pkg).Error
}

func (r *CustomerPackageRepository) GetByID(id uint) (*model.CustomerPackage, error) {
	var pkg model.CustomerPackage
	err := r.db.Preload("Customer.User").Preload("Service").First(&pkg, id).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *CustomerPackageRepository) GetByCustomerID(customerID uint) ([]model.CustomerPackage, error) {
	var packages []model.CustomerPackage
	err := r.db.Where("customer_id = ? AND status = 1", customerID).
		Preload("Service").Find(&packages).Error
	return packages, err
}

func (r *CustomerPackageRepository) GetAvailableByCustomerAndService(customerID, serviceID uint) (*model.CustomerPackage, error) {
	var pkg model.CustomerPackage
	err := r.db.Where("customer_id = ? AND service_id = ? AND used_count < total_count AND status = 1 AND expire_date >= ?",
		customerID, serviceID, time.Now()).
		First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *CustomerPackageRepository) IncrementUsedCount(id uint, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&model.CustomerPackage{}).Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *CustomerPackageRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.CustomerPackage{}).Where("id = ?", id).
		Update("status", status).Error
}

type MemberCardRepository struct {
	db *gorm.DB
}

func NewMemberCardRepository(db *gorm.DB) *MemberCardRepository {
	return &MemberCardRepository{db: db}
}

func (r *MemberCardRepository) Create(card *model.MemberCard) error {
	return r.db.Create(card).Error
}

func (r *MemberCardRepository) GetByID(id uint) (*model.MemberCard, error) {
	var card model.MemberCard
	err := r.db.Preload("Customer.User").First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *MemberCardRepository) GetByCustomerID(customerID uint) ([]model.MemberCard, error) {
	var cards []model.MemberCard
	err := r.db.Where("customer_id = ? AND status = 1", customerID).Find(&cards).Error
	return cards, err
}

func (r *MemberCardRepository) DeductBalance(id uint, amount float64, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&model.MemberCard{}).Where("id = ? AND balance >= ?", id, amount).
		UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error
}

func (r *MemberCardRepository) AddBalance(id uint, amount float64) error {
	return r.db.Model(&model.MemberCard{}).Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *MemberCardRepository) Update(card *model.MemberCard) error {
	return r.db.Save(card).Error
}
