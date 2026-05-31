package repository

import (
	"tea-platform/internal/models"
	"tea-platform/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return database.GetDB().Create(user).Error
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.GetDB().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return database.GetDB().Save(user).Error
}

func (r *UserRepository) DeleteUser(id uint) error {
	return database.GetDB().Delete(&models.User{}, id).Error
}

func (r *UserRepository) GetUserList(page, pageSize int, keyword string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := database.GetDB().Model(&models.User{})

	if keyword != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) UpdateUserStatus(id uint, status string) error {
	return database.GetDB().Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}