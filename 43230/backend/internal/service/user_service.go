package service

import (
	"context"
	"errors"
	"time"

	"print3d-platform/internal/config"
	"print3d-platform/internal/models"
	"print3d-platform/internal/repository"
	"print3d-platform/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type RegisterRequest struct {
	Email        string          `json:"email" binding:"required,email"`
	Username     string          `json:"username" binding:"required,min=3,max=50"`
	Password     string          `json:"password" binding:"required,min=6,max=128"`
	Role         models.UserRole `json:"role" binding:"required"`
	Phone        string          `json:"phone"`
	RealName     string          `json:"real_name"`
	CompanyName  string          `json:"company_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    int64       `json:"expires_at"`
	User         interface{} `json:"user"`
}

func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	emailExists, err := s.userRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("email already exists")
	}

	usernameExists, err := s.userRepo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("username already exists")
	}

	now := time.Now()
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: utils.HashPassword(req.Password),
		Role:         req.Role,
		Status:       models.UserStatusActive,
		Phone:        req.Phone,
		RealName:     req.RealName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if req.Role == models.RoleDesigner {
		user.DesignerProfile = &models.DesignerProfile{
			ID:              uuid.New(),
			UserID:          user.ID,
			Nickname:        req.Username,
			Rating:          5.0,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
	} else if req.Role == models.RolePrinter {
		user.PrinterProfile = &models.PrinterProfile{
			ID:              uuid.New(),
			UserID:          user.ID,
			CompanyName:     req.CompanyName,
			Rating:          5.0,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	utils.LogInfo("New user registered: %s (%s)", user.Username, user.Email)
	return user, nil
}

func (s *UserService) Login(ctx context.Context, req *LoginRequest, clientIP string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Status != models.UserStatusActive {
		return nil, errors.New("account is not active")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	tokenPair, err := utils.GenerateTokenPair(
		user.ID,
		user.Email,
		user.Username,
		user.Role,
		s.cfg.JWT.SecretKey,
		s.cfg.JWT.AccessTokenExpireHours,
		s.cfg.JWT.RefreshTokenExpireHours,
	)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateLoginInfo(ctx, user.ID, clientIP)
	if err != nil {
		utils.LogWarn("Failed to update login info for user %s: %v", user.Email, err)
	}

	utils.LogInfo("User logged in: %s from %s", user.Email, clientIP)

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         s.sanitizeUser(user),
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.Subject != "refresh_token" {
		return nil, errors.New("invalid token type")
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, errors.New("account is not active")
	}

	tokenPair, err := utils.GenerateTokenPair(
		user.ID,
		user.Email,
		user.Username,
		user.Role,
		s.cfg.JWT.SecretKey,
		s.cfg.JWT.AccessTokenExpireHours,
		s.cfg.JWT.RefreshTokenExpireHours,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         s.sanitizeUser(user),
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) UpdateDesignerProfile(ctx context.Context, profile *models.DesignerProfile) error {
	profile.UpdatedAt = time.Now()
	return s.userRepo.UpdateDesignerProfile(ctx, profile)
}

func (s *UserService) UpdatePrinterProfile(ctx context.Context, profile *models.PrinterProfile) error {
	profile.UpdatedAt = time.Now()
	return s.userRepo.UpdatePrinterProfile(ctx, profile)
}

func (s *UserService) ListDesigners(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.ListDesigners(ctx, page, pageSize)
}

func (s *UserService) ListPrinters(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.ListPrinters(ctx, page, pageSize)
}

func (s *UserService) GetUserStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	return s.userRepo.GetUserStats(ctx, userID)
}

func (s *UserService) GetNotifications(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Notification, int64, error) {
	return s.userRepo.GetNotifications(ctx, userID, page, pageSize)
}

func (s *UserService) MarkNotificationRead(ctx context.Context, id, userID uuid.UUID) error {
	return s.userRepo.MarkNotificationRead(ctx, id, userID)
}

func (s *UserService) GetTransactions(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Transaction, int64, error) {
	return s.userRepo.GetTransactions(ctx, userID, page, pageSize)
}

func (s *UserService) UpdateBalance(ctx context.Context, userID uuid.UUID, amount float64, tx *gorm.DB) error {
	if tx != nil {
		return tx.WithContext(ctx).Model(&models.User{}).
			Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error
	}
	return s.userRepo.UpdateBalance(ctx, userID, amount)
}

func (s *UserService) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	return s.userRepo.CreateTransaction(ctx, tx, transaction)
}

func (s *UserService) ExecTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return s.userRepo.ExecTx(ctx, fn)
}

func (s *UserService) sanitizeUser(user *models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":               user.ID,
		"email":            user.Email,
		"username":         user.Username,
		"role":             user.Role,
		"status":           user.Status,
		"phone":            user.Phone,
		"avatar":           user.Avatar,
		"credit_score":     user.CreditScore,
		"real_name":        user.RealName,
		"balance":          user.Balance,
		"email_verified":   user.EmailVerified,
		"designer_profile": user.DesignerProfile,
		"printer_profile":  user.PrinterProfile,
		"created_at":       user.CreatedAt,
	}
}
