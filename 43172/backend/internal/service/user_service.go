package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/middleware"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/repository"
	"luxury-trading-platform/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo    *repository.UserRepository
	redisClient *cache.RedisClient
	db          *gorm.DB
}

func NewUserService(userRepo *repository.UserRepository, redisClient *cache.RedisClient, db *gorm.DB) *UserService {
	return &UserService{
		userRepo:    userRepo,
		redisClient: redisClient,
		db:          db,
	}
}

type RegisterRequest struct {
	Username string         `json:"username" binding:"required,min=3,max=100"`
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Phone    string         `json:"phone"`
	RealName string         `json:"real_name"`
	Role     model.UserRole `json:"role" binding:"required,oneof=buyer seller authenticator"`
}

type AuthenticatorRegisterRequest struct {
	RegisterRequest
	LicenseNumber  string `json:"license_number" binding:"required"`
	LicenseFile    string `json:"license_file" binding:"required"`
	Certifications string `json:"certifications"`
	Specialties    string `json:"specialties"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token  string     `json:"token"`
	User   *model.User `json:"user"`
	Expire time.Time  `json:"expire"`
}

func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	existingUser, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingEmail, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Username:    req.Username,
		Email:       req.Email,
		Phone:       req.Phone,
		Password:    hashedPassword,
		RealName:    req.RealName,
		Role:        req.Role,
		Status:      model.UserStatusActive,
		CreditScore: 100,
	}

	if req.Role == model.RoleAuthenticator {
		user.Status = model.UserStatusPending
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	switch req.Role {
	case model.RoleBuyer:
		buyerProfile := &model.BuyerProfile{
			UserID: user.ID,
		}
		if err := s.userRepo.CreateBuyerProfile(buyerProfile); err != nil {
			return nil, fmt.Errorf("failed to create buyer profile: %w", err)
		}
	case model.RoleSeller:
		sellerProfile := &model.SellerProfile{
			UserID: user.ID,
		}
		if err := s.userRepo.CreateSellerProfile(sellerProfile); err != nil {
			return nil, fmt.Errorf("failed to create seller profile: %w", err)
		}
	}

	cacheKey := fmt.Sprintf("user:%d", user.ID)
	if s.redisClient != nil {
		_ = s.redisClient.Set(ctx, cacheKey, user, 24*time.Hour)
	}

	return user, nil
}

func (s *UserService) RegisterAuthenticator(ctx context.Context, req *AuthenticatorRegisterRequest) (*model.User, error) {
	user, err := s.Register(ctx, &req.RegisterRequest)
	if err != nil {
		return nil, err
	}

	authProfile := &model.AuthenticatorProfile{
		UserID:         user.ID,
		LicenseNumber:  req.LicenseNumber,
		LicenseFile:    req.LicenseFile,
		Certifications: req.Certifications,
		Specialties:    req.Specialties,
		Status:         model.AuthenticatorStatusPending,
	}

	if err := s.userRepo.CreateAuthenticatorProfile(authProfile); err != nil {
		return nil, fmt.Errorf("failed to create authenticator profile: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("invalid username or password")
	}

	if user.Status == model.UserStatusBanned {
		return nil, errors.New("account is banned")
	}

	if user.Status == model.UserStatusInactive {
		return nil, errors.New("account is inactive")
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenKey := fmt.Sprintf("token:%d", user.ID)
	if s.redisClient != nil {
		_ = s.redisClient.Set(ctx, tokenKey, token, 24*time.Hour)
	}

	return &LoginResponse{
		Token:  token,
		User:   user,
		Expire: time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	if s.redisClient != nil {
		cached, err := s.redisClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			return s.userRepo.FindByID(id)
		}
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if s.redisClient != nil {
		_ = s.redisClient.Set(ctx, cacheKey, user, 24*time.Hour)
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if avatar, ok := updates["avatar"].(string); ok {
		user.Avatar = avatar
	}
	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}
	if address, ok := updates["address"].(string); ok {
		user.Address = address
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	cacheKey := fmt.Sprintf("user:%d", id)
	if s.redisClient != nil {
		_ = s.redisClient.Del(ctx, cacheKey)
	}

	return user, nil
}

func (s *UserService) ListUsers(page, pageSize int, role model.UserRole, status model.UserStatus) ([]model.User, int64, error) {
	page, pageSize = utils.ValidatePage(page, pageSize)
	return s.userRepo.List(page, pageSize, role, status)
}

func (s *UserService) ApproveAuthenticator(id uint) error {
	profile, err := s.userRepo.FindAuthenticatorProfileByUserID(id)
	if err != nil {
		return fmt.Errorf("failed to find authenticator profile: %w", err)
	}
	if profile == nil {
		return errors.New("authenticator profile not found")
	}

	if profile.Status != model.AuthenticatorStatusPending {
		return errors.New("authenticator is not in pending status")
	}

	err = s.userRepo.ApproveAuthenticator(profile.ID)
	if err != nil {
		return fmt.Errorf("failed to approve authenticator: %w", err)
	}

	err = s.userRepo.UpdateStatus(id, model.UserStatusActive)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

func (s *UserService) RejectAuthenticator(id uint, reason string) error {
	profile, err := s.userRepo.FindAuthenticatorProfileByUserID(id)
	if err != nil {
		return fmt.Errorf("failed to find authenticator profile: %w", err)
	}
	if profile == nil {
		return errors.New("authenticator profile not found")
	}

	if profile.Status != model.AuthenticatorStatusPending {
		return errors.New("authenticator is not in pending status")
	}

	err = s.userRepo.RejectAuthenticator(profile.ID, reason)
	if err != nil {
		return fmt.Errorf("failed to reject authenticator: %w", err)
	}

	return nil
}

func (s *UserService) ListAuthenticators(page, pageSize int, status model.AuthenticatorStatus) ([]model.AuthenticatorProfile, int64, error) {
	page, pageSize = utils.ValidatePage(page, pageSize)
	return s.userRepo.ListAuthenticators(page, pageSize, status)
}
