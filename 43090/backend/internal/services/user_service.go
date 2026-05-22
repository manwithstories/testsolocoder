package services

import (
	"errors"
	"time"

	"auction-system/internal/dto"
	"auction-system/internal/models"
	"auction-system/internal/utils"
	"auction-system/pkg/logger"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*models.User, error) {
	var count int64
	models.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名或邮箱已存在")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     "user",
		Status:   1,
	}

	if err := models.DB.Create(user).Error; err != nil {
		logger.Error("Failed to create user: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user models.User
	if err := models.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("密码错误")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		UserInfo: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Balance:  user.Balance,
		},
	}, nil
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(userID uint, req *dto.UpdateUserRequest) error {
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	updates["updated_at"] = time.Now()

	return models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return models.DB.Model(&user).Update("password", hashedPassword).Error
}

func (s *UserService) GetUserList(page, pageSize int, keyword string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := models.DB.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (s *UserService) UpdateUserStatus(userID uint, status int) error {
	return models.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error
}
