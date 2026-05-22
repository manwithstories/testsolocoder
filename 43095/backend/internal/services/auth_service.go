package services

import (
	"errors"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"
	"medical-platform/pkg/utils"

	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	FullName string `json:"full_name" binding:"required,max=50"`
	Gender   string `json:"gender" binding:"omitempty,oneof=男 女 其他"`
	Role     string `json:"role" binding:"required,oneof=patient doctor admin"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

func (s *AuthService) Register(req *RegisterRequest) (*models.User, error) {
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名或邮箱已存在")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Gender:       req.Gender,
		Role:         models.UserRole(req.Role),
		IsActive:     true,
	}

	err = database.WithTransaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return errors.New("注册失败")
		}

		switch user.Role {
		case models.RolePatient:
			patient := &models.Patient{
				UserID: user.ID,
			}
			if err := tx.Create(patient).Error; err != nil {
				return errors.New("创建患者信息失败")
			}

			healthRecord := &models.HealthRecord{
				PatientID: patient.ID,
			}
			if err := tx.Create(healthRecord).Error; err != nil {
				return errors.New("创建健康档案失败")
			}
		case models.RoleDoctor:
			doctor := &models.Doctor{
				UserID: user.ID,
				Title:  models.TitleResident,
			}
			if err := tx.Create(doctor).Error; err != nil {
				return errors.New("创建医生信息失败")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户名或密码错误")
		}
		return nil, "", errors.New("登录失败")
	}

	if !user.IsActive {
		return nil, "", errors.New("账号已被禁用")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, "", errors.New("用户名或密码错误")
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		return nil, "", errors.New("生成令牌失败")
	}

	return &user, token, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("获取用户信息失败")
	}
	return &user, nil
}

func (s *AuthService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	if err := s.db.Model(&user).Update("password_hash", hashedPassword).Error; err != nil {
		return errors.New("修改密码失败")
	}

	return nil
}
