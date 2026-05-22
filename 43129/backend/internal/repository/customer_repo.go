package repository

import (
	"beauty-salon-system/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) GetByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Preload("User").First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) GetByUserID(userID uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Update(customer *model.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) List(page, pageSize int, keyword string) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	query := r.db.Model(&model.Customer{}).Preload("User")
	if keyword != "" {
		query = query.Joins("JOIN users ON users.id = customers.user_id").
			Where("customers.name LIKE ? OR users.phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("customers.created_at DESC").Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepository) AddPoints(id uint, points int) error {
	return r.db.Model(&model.Customer{}).Where("id = ?", id).
		UpdateColumn("points", gorm.Expr("points + ?", points)).Error
}

func (r *CustomerRepository) DeductPoints(id uint, points int) error {
	return r.db.Model(&model.Customer{}).Where("id = ? AND points >= ?", id, points).
		UpdateColumn("points", gorm.Expr("points - ?", points)).Error
}

func (r *CustomerRepository) UpdateLevel(id uint, level int) error {
	return r.db.Model(&model.Customer{}).Where("id = ?", id).
		Update("member_level", level).Error
}

func (r *CustomerRepository) AddVisit(id uint) error {
	return r.db.Model(&model.Customer{}).Where("id = ?", id).
		UpdateColumns(map[string]interface{}{
			"visit_count":   gorm.Expr("visit_count + 1"),
			"last_visit_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *CustomerRepository) AddTotalSpent(id uint, amount float64) error {
	return r.db.Model(&model.Customer{}).Where("id = ?", id).
		UpdateColumn("total_spent", gorm.Expr("total_spent + ?", amount)).Error
}
