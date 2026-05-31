package services

import (
	"errors"
	"fmt"
	"time"

	"health-platform/models"
	"health-platform/repository"
	"health-platform/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

type RegisterRequest struct {
	Username string      `json:"username" binding:"required,min=3,max=50"`
	Password string      `json:"password" binding:"required,min=6,max=50"`
	RealName string      `json:"real_name" binding:"required,max=50"`
	Phone    string      `json:"phone" binding:"required,max=20"`
	Email    string      `json:"email" binding:"omitempty,email,max=100"`
	Role     models.UserRole `json:"role" binding:"required"`
	CompanyID *uint     `json:"company_id"`
	AgencyID  *uint     `json:"agency_id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	User     *models.User `json:"user"`
	ExpiresAt int64      `json:"expires_at"`
}

func (s *AuthService) Register(req *RegisterRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	existingPhone, _ := s.userRepo.FindByPhone(req.Phone)
	if existingPhone != nil {
		return nil, errors.New("手机号已注册")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		RealName:  req.RealName,
		Phone:     req.Phone,
		Email:     req.Email,
		Role:      req.Role,
		CompanyID: req.CompanyID,
		AgencyID:  req.AgencyID,
		Status:    1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("注册失败: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), user.CompanyID, user.AgencyID)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	s.userRepo.UpdateLastLogin(user.ID)

	return &LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.userRepo.FindByID(&user, userID); err != nil {
		return nil, err
	}
	return &user, nil
}
