package services

import (
	"errors"
	"time"

	"business-registration-platform/database"
	"business-registration-platform/models"
	"business-registration-platform/utils"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	RealName string `json:"realName" binding:"required,max=100"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Role     string `json:"role" binding:"required,oneof=entrepreneur agent"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string     `json:"token"`
	UserInfo *models.User `json:"userInfo"`
}

func (s *AuthService) Register(req *RegisterRequest) (*models.User, error) {
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     models.UserRole(req.Role),
		Status:   models.UserStatusActive,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	if user.Role == models.RoleAgent {
		agentProfile := &models.AgentProfile{
			UserID:          user.ID,
			EmployeeNo:      "AGT" + time.Now().Format("20060102150405"),
			SpecialtyTags:   "general",
			MaxApplications: 5,
			CurrentApps:     0,
			WorkStartTime:   "09:00",
			WorkEndTime:     "18:00",
			Status:          "available",
		}
		database.DB.Create(agentProfile)
	}

	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.Status != models.UserStatusActive {
		return nil, errors.New("account is not active")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	database.DB.Model(&user).Update("last_login_at", now)

	if user.Role == models.RoleAgent {
		database.DB.Preload("AgentProfile").First(&user, user.ID)
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: &user,
	}, nil
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateProfile(userID uint, data map[string]interface{}) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(data).Error
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return database.DB.Model(&user).Update("password", hashedPassword).Error
}
