package service

import (
	"errors"
	"ticket-system/internal/dto"
	"ticket-system/internal/jwt"
	"ticket-system/internal/logger"
	"ticket-system/internal/models"
	"ticket-system/internal/util"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*models.User, error) {
	var count int64
	models.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     models.RoleUser,
	}

	if err := models.DB.Create(user).Error; err != nil {
		logger.Log.Errorf("Create user failed: %v", err)
		return nil, errors.New("创建用户失败")
	}

	logger.Log.Infof("User registered: %s", user.Username)
	return user, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	if !util.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	logger.Log.Infof("User logged in: %s", user.Username)
	return &dto.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

func (s *UserService) GetUserList(req *dto.UserListRequest) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := models.DB.Model(&models.User{})

	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uint, req *dto.UserUpdateRequest) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := models.DB.Save(user).Error; err != nil {
		return nil, err
	}

	logger.Log.Infof("User updated: %d", id)
	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	if err := models.DB.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	logger.Log.Infof("User deleted: %d", id)
	return nil
}
