package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"tea-platform/internal/models"
	"tea-platform/internal/repository"
	"tea-platform/pkg/auth"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,min=11,max=20"`
	Role     string `json:"role" binding:"required,oneof=admin buyer seller appraiser"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Phone  string `json:"phone" binding:"omitempty,min=11,max=20"`
	Avatar string `json:"avatar" binding:"omitempty,max=512"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=128"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	existing, _ := s.repo.GetUserByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("用户名已存在")
	}

	existing, _ = s.repo.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("邮箱已被注册")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     models.UserRole(req.Role),
		Status:   models.UserStatusActive,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return user, nil
}

func (s *UserService) Login(req *LoginRequest) (string, *models.User, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	if user.Status != models.UserStatusActive {
		return "", nil, errors.New("账号已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	token, err := auth.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return "", nil, errors.New("生成令牌失败")
	}

	return token, user, nil
}

func (s *UserService) GetUserInfo(id uint) (*models.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(id uint, req *UpdateProfileRequest) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	if req.Email != "" {
		existing, _ := s.repo.GetUserByEmail(req.Email)
		if existing != nil && existing.ID != id {
			return errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	return s.repo.UpdateUser(user)
}

func (s *UserService) ChangePassword(id uint, req *ChangePasswordRequest) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = string(hashedPassword)
	return s.repo.UpdateUser(user)
}

func (s *UserService) GetUserList(page, pageSize int, keyword string) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.GetUserList(page, pageSize, keyword)
}

func (s *UserService) UpdateUserStatus(id uint, status string) error {
	return s.repo.UpdateUserStatus(id, status)
}