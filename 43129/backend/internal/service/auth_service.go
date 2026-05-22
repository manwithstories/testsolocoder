package service

import (
	"errors"
	"fmt"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
	"beauty-salon-system/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtSecret string
	jwtExpire int
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpire int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  *model.User `json:"user"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	if !utils.ValidatePhone(req.Phone) {
		return nil, errors.New("invalid phone number")
	}

	user, err := s.userRepo.GetByPhone(req.Phone)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != 1 {
		return nil, errors.New("account is disabled")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.ID, user.Phone, user.Role, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Register(req *RegisterRequest) (*LoginResponse, error) {
	if !utils.ValidatePhone(req.Phone) {
		return nil, errors.New("invalid phone number")
	}

	existingUser, _ := s.userRepo.GetByPhone(req.Phone)
	if existingUser != nil {
		return nil, errors.New("phone already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "customer"
	}

	user := &model.User{
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     role,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := utils.GenerateToken(user.ID, user.Phone, user.Role, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GetUser(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}
