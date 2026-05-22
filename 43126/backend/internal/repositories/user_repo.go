package repositories

import (
	"meeting-room/internal/models"
	"meeting-room/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return utils.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := utils.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := utils.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := utils.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) VerifyPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (r *UserRepository) Update(user *models.User) error {
	return utils.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return utils.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) List(page, pageSize int, role string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := utils.DB.Model(&models.User{})
	if role != "" {
		db = db.Where("role = ?", role)
	}

	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) UpdateRole(id uint, role string) error {
	return utils.DB.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error
}
