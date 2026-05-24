package service

import (
	"errors"
	"survey-platform/internal/dto"
	"survey-platform/internal/model"
	"survey-platform/internal/repository"
	"survey-platform/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		RoleID:   3,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Email, "viewer")
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Role:     "viewer",
	}, nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Status != 1 {
		return nil, errors.New("account is disabled")
	}

	now := time.Now()
	user.LastLogin = &now
	s.userRepo.Update(user)

	roleName := "viewer"
	if user.Role != nil {
		roleName = user.Role.Name
	}

	token, err := utils.GenerateToken(user.ID, user.Email, roleName)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Role:     roleName,
	}, nil
}

func (s *AuthService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	roleName := "viewer"
	if user.Role != nil {
		roleName = user.Role.Name
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Role:      roleName,
		Status:    user.Status,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AuthService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	return s.userRepo.Update(user)
}

func (s *AuthService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}
