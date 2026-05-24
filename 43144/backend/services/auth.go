package services

import (
	"errors"
	"fmt"

	"pet-adoption-platform/database"
	"pet-adoption-platform/models"
	"pet-adoption-platform/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(req *models.RegisterRequest) (*models.User, error) {
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     models.UserRole(req.Role),
	}

	if req.Role == "rescue" && req.RescueName != "" {
		rescue := &models.RescueStation{
			Name:   req.RescueName,
			Status: models.RescueStatusPending,
		}
		if err := database.DB.Create(rescue).Error; err != nil {
			return nil, fmt.Errorf("failed to create rescue station: %w", err)
		}
		user.RescueID = &rescue.ID
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User
	if err := database.DB.Preload("Rescue").Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	var rescueID *uint
	if user.RescueID != nil {
		rescueID = user.RescueID
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), rescueID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     string(user.Role),
		RescueID: rescueID,
	}, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Rescue").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id uint, updates map[string]interface{}) (*models.User, error) {
	result := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	return GetUserByID(id)
}

func ListUsers(page, pageSize int, role string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Rescue").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func WithTransaction(fn func(tx *gorm.DB) error) error {
	return database.DB.Transaction(fn)
}
