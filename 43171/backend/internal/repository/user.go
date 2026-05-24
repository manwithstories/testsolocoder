package repository

import (
	"drone-rental/internal/config"
	"drone-rental/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(user *model.User) error {
	return config.DB.Create(user).Error
}

func (r *UserRepo) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	err := config.DB.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *UserRepo) Update(user *model.User) error {
	return config.DB.Save(user).Error
}

func (r *UserRepo) UpdateRating(id uint, rating float64) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).
		Update("rating", rating).Error
}

func (r *UserRepo) UpdateBalance(id uint, amount float64) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("balance", config.DB.Raw("balance + ?", amount)).Error
}

func (r *UserRepo) UpdateDeposit(id uint, amount float64) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("deposit", config.DB.Raw("deposit + ?", amount)).Error
}

func (r *UserRepo) List(page, pageSize int, role model.Role, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := config.DB.Model(&model.User{})
	if role != "" {
		db = db.Where("role = ?", role)
	}
	if keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	db.Count(&total)
	err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
