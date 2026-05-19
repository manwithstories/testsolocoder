package repository

import (
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) List(req *dto.PaginationRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := database.DB.Model(&model.User{})
	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) UpdateRole(id uint, role model.UserRole) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).Update("role", role).Error
}
