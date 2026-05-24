package repository

import (
	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/utils"
	"time"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(user *model.User) error {
	return utils.DB.Create(user).Error
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := utils.DB.Preload("Profile").First(&user, id).Error
	return &user, err
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := utils.DB.Where("username = ?", username).Preload("Profile").First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := utils.DB.Where("phone = ?", phone).Preload("Profile").First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := utils.DB.Where("email = ?", email).Preload("Profile").First(&user).Error
	return &user, err
}

func (r *UserRepo) Update(id uint, updates map[string]interface{}) error {
	return utils.DB.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepo) UpdateStatus(id uint, status model.UserStatus) error {
	return utils.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRepo) UpdateVerifyStatus(id uint, status model.VerifyStatus) error {
	return utils.DB.Model(&model.User{}).Where("id = ?", id).Update("verify_status", status).Error
}

func (r *UserRepo) UpdateMember(id uint, level string, expire *time.Time) error {
	return utils.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"member_level":  level,
		"member_expire": expire,
	}).Error
}

func (r *UserRepo) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := utils.DB.Model(&model.User{})
	if keyword != "" {
		db = db.Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepo) CountTotal() (int64, error) {
	var count int64
	err := utils.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepo) CountActiveToday() (int64, error) {
	var count int64
	err := utils.DB.Model(&model.User{}).Where("last_login_at > CURDATE()").Count(&count).Error
	return count, err
}

func (r *UserRepo) CountByRole(role string) (int64, error) {
	var count int64
	err := utils.DB.Model(&model.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

func (r *UserRepo) CountByVerifyStatus(status string) (int64, error) {
	var count int64
	err := utils.DB.Model(&model.User{}).Where("verify_status = ?", status).Count(&count).Error
	return count, err
}
