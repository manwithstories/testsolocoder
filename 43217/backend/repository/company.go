package repository

import (
	"health-platform/models"

	"gorm.io/gorm"
)

type CompanyRepository struct {
	*BaseRepository
}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *CompanyRepository) FindByName(name string) (*models.Company, error) {
	var company models.Company
	err := r.DB.Where("name = ?", name).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) FindByUnifiedCode(code string) (*models.Company, error) {
	var company models.Company
	err := r.DB.Where("unified_code = ?", code).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) GetWithDepartments(companyID uint) (*models.Company, error) {
	var company models.Company
	err := r.DB.Preload("Departments").First(&company, companyID).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) UpdateBudget(companyID uint, usedBudget float64) error {
	return r.DB.Model(&models.Company{}).Where("id = ?", companyID).
		Update("used_budget", gorm.Expr("used_budget + ?", usedBudget)).Error
}

func (r *CompanyRepository) UpdateBalance(companyID uint, amount float64) error {
	return r.DB.Model(&models.Company{}).Where("id = ?", companyID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}
