package repository

import (
	"music-platform/internal/model"
	"music-platform/pkg/database"
	"music-platform/pkg/utils"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.Preload("ArtistInfo").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).Preload("ArtistInfo").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).Preload("ArtistInfo").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("phone = ?", phone).Preload("ArtistInfo").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) UpdateUserInfo(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepository) Delete(id uint) error {
	return database.DB.Delete(&model.User{}, id).Error
}

func (r *UserRepository) List(page, pageSize int, keyword string, role string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := database.DB.Model(&model.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("ArtistInfo").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) UpdateLoginInfo(id uint, ip string) error {
	now := utils.GetTodayStart()
	return database.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	}).Error
}

func (r *UserRepository) FindByIds(ids []uint) ([]model.User, error) {
	var users []model.User
	err := database.DB.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (r *UserRepository) UpdateArtistInfo(info *model.ArtistInfo) error {
	return database.DB.Save(info).Error
}

func (r *UserRepository) CreateArtistInfo(info *model.ArtistInfo) error {
	return database.DB.Create(info).Error
}

func (r *UserRepository) FindArtistInfoByUserID(userID uint) (*model.ArtistInfo, error) {
	var info model.ArtistInfo
	err := database.DB.Where("user_id = ?", userID).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *UserRepository) UpdateArtistBalance(artistID uint, amount float64) error {
	return database.DB.Model(&model.ArtistInfo{}).Where("id = ?", artistID).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *UserRepository) UpdateFrozenBalance(artistID uint, amount float64) error {
	return database.DB.Model(&model.ArtistInfo{}).Where("id = ?", artistID).
		UpdateColumn("frozen_balance", gorm.Expr("frozen_balance + ?", amount)).Error
}
