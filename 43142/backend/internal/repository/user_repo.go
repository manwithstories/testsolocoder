package repository

import (
	"time"

	"gorm.io/gorm"

	"recruitment-platform/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Company").Preload("Profile").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Company").Preload("Profile").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) List(page, pageSize int, role models.UserRole, status models.UserStatus) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Company").Preload("Profile").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) UpdateLoginAttempts(id uint, attempts int) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"login_attempts": attempts,
			"last_login_at":  time.Now(),
		}).Error
}

func (r *UserRepository) LockAccount(id uint, until time.Time) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"locked_until": until,
			"status":       models.UserStatusInactive,
		}).Error
}

func (r *UserRepository) CreateProfile(profile *models.ApplicantProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserRepository) UpdateProfile(profile *models.ApplicantProfile) error {
	return r.db.Save(profile).Error
}

func (r *UserRepository) FindProfileByUserID(userID uint) (*models.ApplicantProfile, error) {
	var profile models.ApplicantProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) CreateCompany(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *UserRepository) UpdateCompany(company *models.Company) error {
	return r.db.Save(company).Error
}

func (r *UserRepository) FindCompanyByID(id uint) (*models.Company, error) {
	var company models.Company
	err := r.db.First(&company, id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *UserRepository) ListCompanies(page, pageSize int) ([]models.Company, int64, error) {
	var companies []models.Company
	var total int64

	r.db.Model(&models.Company{}).Count(&total)
	err := r.db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&companies).Error

	return companies, total, err
}
