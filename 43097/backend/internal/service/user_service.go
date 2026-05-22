package service

import (
	"errors"
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/jwt"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/repository"
	"hotel-system/internal/utils"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*model.User, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetUserInfo(id uint) (*model.User, error)
	UpdateUser(id uint, req *dto.UserUpdateRequest) (*model.User, error)
	DeleteUser(id uint) error
	ListUsers(req *dto.UserListRequest) ([]model.User, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(req *dto.RegisterRequest) (*model.User, error) {
	existingUser, err := s.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		logger.Warnf("用户名已存在: %s", req.Username)
		return nil, errors.New("用户名已存在")
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		logger.Errorf("创建用户失败: %v", err)
		return nil, err
	}

	logger.Infof("用户注册成功: %s", user.Username)
	return user, nil
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		logger.Warnf("登录失败，用户不存在: %s", req.Username)
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != model.UserStatusActive {
		logger.Warnf("登录失败，用户状态非活跃: %s, status: %s", req.Username, user.Status)
		return nil, errors.New("用户已被禁用")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		logger.Warnf("登录失败，密码错误: %s", req.Username)
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		logger.Errorf("生成token失败: %v", err)
		return nil, errors.New("生成token失败")
	}

	logger.Infof("用户登录成功: %s", user.Username)

	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		RealName:  user.RealName,
		Phone:     user.Phone,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return &dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}, nil
}

func (s *userService) GetUserInfo(id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取用户信息失败: id=%d, err=%v", id, err)
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (s *userService) UpdateUser(id uint, req *dto.UserUpdateRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Errorf("更新用户失败，用户不存在: id=%d, err=%v", id, err)
		return nil, errors.New("用户不存在")
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logger.Errorf("密码哈希失败: %v", err)
			return nil, errors.New("密码更新失败")
		}
		user.Password = hashedPassword
	}

	err = s.userRepo.Update(user)
	if err != nil {
		logger.Errorf("更新用户失败: %v", err)
		return nil, errors.New("更新用户失败")
	}

	logger.Infof("用户更新成功: id=%d, username=%s", user.ID, user.Username)
	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Errorf("删除用户失败，用户不存在: id=%d, err=%v", id, err)
		return errors.New("用户不存在")
	}

	err = s.userRepo.Delete(id)
	if err != nil {
		logger.Errorf("删除用户失败: %v", err)
		return errors.New("删除用户失败")
	}

	logger.Infof("用户删除成功: id=%d", id)
	return nil
}

func (s *userService) ListUsers(req *dto.UserListRequest) ([]model.User, int64, error) {
	page := req.GetPage()
	pageSize := req.GetPageSize()

	users, total, err := s.userRepo.List(page, pageSize, req.Username, req.Role, req.Status)
	if err != nil {
		logger.Errorf("获取用户列表失败: %v", err)
		return nil, 0, errors.New("获取用户列表失败")
	}

	return users, total, nil
}
