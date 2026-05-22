package services

import (
	"errors"
	"time"

	"meeting-room/internal/middleware"
	"meeting-room/internal/models"
	"meeting-room/internal/repositories"
	"meeting-room/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	RealName   string `json:"real_name"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
}

type UpdateUserRequest struct {
	RealName   string `json:"real_name"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱或密码错误")
		}
		return nil, err
	}

	if !s.userRepo.VerifyPassword(user, req.Password) {
		return nil, errors.New("邮箱或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	token, err := middleware.GenerateToken(user.ID, user.Username, string(user.Role), user.Email)
	if err != nil {
		return nil, err
	}

	utils.RedisSet("session:"+token, user.ID, 24*time.Hour)

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	_, err = s.userRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("用户名已被使用")
	}

	user := &models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   req.Password,
		RealName:   req.RealName,
		Phone:      req.Phone,
		Department: req.Department,
		Role:       models.RoleUser,
		Status:     1,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateProfile(userID uint, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Department != "" {
		user.Department = req.Department
	}

	err = s.userRepo.Update(user)
	return user, err
}

func (s *UserService) ListUsers(page, pageSize int, role string) ([]models.User, int64, error) {
	return s.userRepo.List(page, pageSize, role)
}

func (s *UserService) UpdateUserRole(userID uint, role string) error {
	validRoles := map[string]bool{
		"admin":       true,
		"space_admin": true,
		"user":        true,
	}
	if !validRoles[role] {
		return errors.New("无效的角色")
	}
	return s.userRepo.UpdateRole(userID, role)
}

func (s *UserService) DeleteUser(userID uint) error {
	return s.userRepo.Delete(userID)
}
