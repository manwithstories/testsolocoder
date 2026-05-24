package repository

import (
	"errors"
	"time"

	"luxury-trading-platform/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("AuthenticatorProfile").
		Preload("SellerProfile").
		Preload("BuyerProfile").
		First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) List(page, pageSize int, role model.UserRole, status model.UserStatus) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("AuthenticatorProfile").
		Preload("SellerProfile").
		Preload("BuyerProfile").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) UpdateStatus(id uint, status model.UserStatus) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRepository) UpdateCreditScore(id uint, score int) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("credit_score", score).Error
}

func (r *UserRepository) CreateAuthenticatorProfile(profile *model.AuthenticatorProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserRepository) UpdateAuthenticatorProfile(profile *model.AuthenticatorProfile) error {
	return r.db.Save(profile).Error
}

func (r *UserRepository) FindAuthenticatorProfileByUserID(userID uint) (*model.AuthenticatorProfile, error) {
	var profile model.AuthenticatorProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) ListAuthenticators(page, pageSize int, status model.AuthenticatorStatus) ([]model.AuthenticatorProfile, int64, error) {
	var profiles []model.AuthenticatorProfile
	var total int64

	query := r.db.Model(&model.AuthenticatorProfile{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&profiles).Error

	return profiles, total, err
}

func (r *UserRepository) ApproveAuthenticator(id uint) error {
	now := time.Now()
	return r.db.Model(&model.AuthenticatorProfile{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      model.AuthenticatorStatusApproved,
			"verified_at": now,
		}).Error
}

func (r *UserRepository) RejectAuthenticator(id uint, reason string) error {
	return r.db.Model(&model.AuthenticatorProfile{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           model.AuthenticatorStatusRejected,
			"rejection_reason": reason,
		}).Error
}

func (r *UserRepository) CreateSellerProfile(profile *model.SellerProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserRepository) FindSellerProfileByUserID(userID uint) (*model.SellerProfile, error) {
	var profile model.SellerProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) CreateBuyerProfile(profile *model.BuyerProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserRepository) FindBuyerProfileByUserID(userID uint) (*model.BuyerProfile, error) {
	var profile model.BuyerProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}
