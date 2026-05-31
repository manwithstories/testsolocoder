package repository

import (
	"health-platform/models"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	*BaseRepository
}

func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *EmployeeRepository) FindByCompanyID(companyID uint, page, pageSize int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.DB.Model(&models.Employee{}).Where("company_id = ?", companyID)
	query.Count(&total)

	err := query.Preload("Department").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&employees).Error
	return employees, total, err
}

func (r *EmployeeRepository) FindByDepartmentID(departmentID uint) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.DB.Where("department_id = ?", departmentID).Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) FindByUserID(userID uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.DB.Where("user_id = ?", userID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) FindByEmployeeNo(companyID uint, employeeNo string) (*models.Employee, error) {
	var employee models.Employee
	err := r.DB.Where("company_id = ? AND employee_no = ?", companyID, employeeNo).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) UpdateQuota(employeeID uint, usedQuota int) error {
	return r.DB.Model(&models.Employee{}).Where("id = ?", employeeID).
		Update("used_quota", gorm.Expr("used_quota + ?", usedQuota)).Error
}

func (r *EmployeeRepository) BatchInsert(employees []models.Employee) error {
	return r.DB.Create(&employees).Error
}

func (r *EmployeeRepository) CountByCompanyID(companyID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Employee{}).Where("company_id = ?", companyID).Count(&count).Error
	return count, err
}
