package repositories

import (
	"context"
	"errors"
	"time"

	"consultation-platform/models"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: utils.GetDB()}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

func (r *UserRepository) FindAll(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	err := r.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FindByRole(role models.Role, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.User{}).Where("role = ?", role)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) UpdateVerificationStatus(id uuid.UUID, status models.VerificationStatus, note string) error {
	updates := map[string]interface{}{
		"verification_status": status,
		"verification_note":   note,
	}
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepository) UsernameOrEmailExists(username, email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) FindPendingVerifications(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.User{}).Where("verification_status = ?", models.VerificationPending)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) Transaction(ctx context.Context, fn func(*UserRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := &UserRepository{db: tx}
		return fn(repo)
	})
}

func (r *UserRepository) FindActiveProfessionalByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND role = ? AND verification_status = ? AND is_active = ?",
		id, models.RoleProfessional, models.VerificationApproved, true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("professional not found or not active")
		}
		return nil, err
	}
	return &user, nil
}
