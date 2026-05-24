package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
	jwtpkg "furniture-platform/pkg/jwt"
	"furniture-platform/pkg/password"
)

// UserService 用户业务逻辑层
type UserService struct {
	repo      *repository.UserRepository
	db        *gorm.DB
	jwtSecret string
	jwtTTL    time.Duration
}

// NewUserService 创建用户服务
func NewUserService(repo *repository.UserRepository, db *gorm.DB, jwtSecret string, jwtTTLSeconds int) *UserService {
	return &UserService{
		repo:      repo,
		db:        db,
		jwtSecret: jwtSecret,
		jwtTTL:    time.Duration(jwtTTLSeconds) * time.Second,
	}
}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterRequest) (*model.User, error) {
	if !model.ValidRole(req.Role) {
		return nil, errors.New("角色不合法")
	}

	existing, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("用户名已存在")
	}

	hashed, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashed,
		Role:     req.Role,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   model.StatusActive,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *dto.LoginRequest) (*model.User, string, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("用户名或密码错误")
	}
	if !user.IsActive() {
		return nil, "", errors.New("账号已被禁用")
	}
	if !password.Verify(req.Password, user.Password) {
		return nil, "", errors.New("用户名或密码错误")
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// GetByID 根据 ID 获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(id uint, req *dto.UpdateProfileRequest) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}
	if !password.Verify(req.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	hashed, err := password.Hash(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.repo.Update(user)
}

// List 分页查询用户列表
func (s *UserService) List(params *dto.UserListRequest) ([]*model.User, int64, error) {
	p := &repository.ListParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.Keyword,
		Role:     params.Role,
		Status:   params.Status,
	}
	return s.repo.List(p)
}
