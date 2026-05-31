package services

import (
	"errors"
	"time"

	"consultation-platform/config"
	"consultation-platform/models"
	"consultation-platform/repositories"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
		cfg:      cfg,
	}
}

type RegisterRequest struct {
	Username         string `json:"username" binding:"required,min=3,max=50"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6,max=100"`
	Role             string `json:"role" binding:"required,oneof=client professional"`
	FullName         string `json:"full_name" binding:"required"`
	Phone            string `json:"phone"`
	VerificationDocs string `json:"verification_docs"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *models.User `json:"user"`
}

type UserServiceInterface interface {
	Register(req *RegisterRequest) (*models.User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, updates map[string]interface{}) (*models.User, error)
	GetUsers(page, pageSize int, role string) ([]models.User, int64, error)
	VerifyProfessional(id uuid.UUID, status models.VerificationStatus, note string) error
	GetPendingVerifications(page, pageSize int) ([]models.User, int64, error)
}

func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	exists, err := s.userRepo.UsernameOrEmailExists(req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := models.Role(req.Role)
	user := &models.User{
		Username:         req.Username,
		Email:            req.Email,
		Password:         string(hashedPassword),
		Role:             role,
		FullName:         req.FullName,
		Phone:            req.Phone,
		VerificationDocs: req.VerificationDocs,
		VerificationStatus: models.VerificationPending,
		IsActive:         true,
	}

	if role == models.RoleClient {
		user.VerificationStatus = models.VerificationApproved
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.Role == models.RoleProfessional && user.VerificationStatus != models.VerificationApproved {
		return nil, errors.New("account is pending verification")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.AccessTokenExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	user.Password = ""
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(s.cfg.JWT.AccessTokenExpire) * time.Hour),
		User:         user,
	}, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if updates["password"] != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updates["password"].(string)), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	for key, value := range updates {
		switch key {
		case "full_name":
			user.FullName = value.(string)
		case "avatar":
			user.Avatar = value.(string)
		case "phone":
			user.Phone = value.(string)
		case "password":
			user.Password = value.(string)
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *UserService) GetUsers(page, pageSize int, role string) ([]models.User, int64, error) {
	if role != "" {
		return s.userRepo.FindByRole(models.Role(role), page, pageSize)
	}
	return s.userRepo.FindAll(page, pageSize)
}

func (s *UserService) VerifyProfessional(id uuid.UUID, status models.VerificationStatus, note string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if user.Role != models.RoleProfessional {
		return errors.New("user is not a professional")
	}

	return s.userRepo.UpdateVerificationStatus(id, status, note)
}

func (s *UserService) GetPendingVerifications(page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.FindPendingVerifications(page, pageSize)
}

func (s *UserService) RefreshToken(userID uuid.UUID) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	if !user.IsActive {
		return "", errors.New("account is disabled")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.AccessTokenExpire)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
