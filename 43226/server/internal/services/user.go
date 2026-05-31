package services

import (
	"fmt"
	"time"

	"museum-server/internal/config"
	"museum-server/internal/dto"
	"museum-server/internal/models"
	"museum-server/internal/repository"
	"museum-server/pkg/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
	jwtCfg   *config.JWTConfig
}

func NewUserService(userRepo *repository.UserRepository, jwtCfg *config.JWTConfig) *UserService {
	return &UserService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return nil, fmt.Errorf("username already exists")
	}

	existing, _ = s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     models.UserRoleNormal,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, expiresAt, err := utils.GenerateToken(s.jwtCfg, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Avatar:   user.Avatar,
		},
	}, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, fmt.Errorf("invalid username or password")
	}

	if user.Status != 1 {
		return nil, fmt.Errorf("account is disabled")
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	token, expiresAt, err := utils.GenerateToken(s.jwtCfg, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Avatar:   user.Avatar,
		},
	}, nil
}

func (s *UserService) GetUser(id uint) (*dto.UserInfo, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Avatar:   user.Avatar,
	}, nil
}

func (s *UserService) UpdateUser(id uint, req *dto.UpdateUserRequest) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	return s.userRepo.Update(user)
}

func (s *UserService) ListUsers(query *dto.CollectionListQuery) ([]models.User, int64, error) {
	return s.userRepo.List(query)
}

func (s *UserService) GetGuides() ([]models.User, error) {
	return s.userRepo.FindByRole(models.UserRoleGuide)
}
