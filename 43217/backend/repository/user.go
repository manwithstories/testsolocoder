package repository

import (
	"health-platform/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(userID uint) error {
	return r.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

func (r *UserRepository) FindByCompanyID(companyID uint) ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("company_id = ?", companyID).Find(&users).Error
	return users, err
}

func (r *UserRepository) FindByAgencyID(agencyID uint) ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("agency_id = ?", agencyID).Find(&users).Error
	return users, err
}
