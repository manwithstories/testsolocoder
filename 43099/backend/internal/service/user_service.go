package service

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/internal/repository"
	"venue-booking/pkg/auth"
	"venue-booking/pkg/email"
)

type UserService struct {
	userRepo *repository.UserRepository
	logService *OperationLogService
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:   repository.NewUserRepository(),
		logService: NewOperationLogService(),
	}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*model.User, string, error) {
	_, err := s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return nil, "", errors.New("username already exists")
	}

	_, err = s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, "", errors.New("email already exists")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		Username:     req.Username,
		Email:          req.Email,
		PasswordHash:     hashedPassword,
		RealName:       req.RealName,
		Phone:          req.Phone,
		Role:           model.RoleUser,
		EmailVerified: false,
		Status:         model.UserStatusActive,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	return user, token, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*model.User, string, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		user, err = s.userRepo.GetByEmail(req.Username)
		if err != nil {
			return nil, "", errors.New("invalid credentials")
		}
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	if user.Status != model.UserStatusActive {
		return nil, "", errors.New("account is inactive")
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	return user, token, err
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateProfile(id uint, req *dto.UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	err = s.userRepo.Update(user)
	return user, err
}

func (s *UserService) List(req *dto.PaginationRequest) ([]model.User, int64, error) {
	return s.userRepo.List(req)
}

func (s *UserService) UpdateRole(id uint, role string) error {
	return s.userRepo.UpdateRole(id, model.UserRole(role))
}

func (s *UserService) SendVerificationCode(email string) error {
	code := generateVerificationCode()
	store := GetVerificationStore()
	store.Store(email, code, 10*time.Minute)

	emailService := email.NewEmailService()
	return emailService.SendVerificationEmail(email, code)
}

func (s *UserService) VerifyEmail(email, code string) error {
	store := GetVerificationStore()
	if err := store.Verify(email, code); err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	user.EmailVerified = true
	return s.userRepo.Update(user)
}

func (s *UserService) SendPasswordReset(email string) error {
	code := generateVerificationCode()
	store := GetVerificationStore()
	store.Store(email, code, 10*time.Minute)

	emailService := email.NewEmailService()
	return emailService.SendPasswordResetEmail(email, code)
}

func (s *UserService) ResetPassword(email, code, newPassword string) error {
	store := GetVerificationStore()
	if err := store.Verify(email, code); err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

func generateVerificationCode() string {
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}
