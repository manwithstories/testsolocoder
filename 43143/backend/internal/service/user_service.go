package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"skillshare/internal/models"
	"skillshare/internal/repository"
	"skillshare/pkg/auth"
	"skillshare/pkg/encryption"
	"skillshare/pkg/validator"
)

type UserService struct {
	userRepo *repository.UserRepository
	secret   string
}

func NewUserService(userRepo *repository.UserRepository, secret string) *UserService {
	return &UserService{
		userRepo: userRepo,
		secret:   secret,
	}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	AuthType string `json:"auth_type"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (s *UserService) Register(input *RegisterInput, ip string) (*models.User, *auth.TokenPair, error) {
	if input.AuthType == string(models.AuthTypeEmail) {
		if err := validator.ValidateEmail(input.Email); err != nil {
			return nil, nil, err
		}
		existingUser, _ := s.userRepo.FindByEmail(input.Email)
		if existingUser != nil {
			return nil, nil, errors.New("邮箱已被注册")
		}
	} else if input.AuthType == string(models.AuthTypePhone) {
		if err := validator.ValidatePhone(input.Phone); err != nil {
			return nil, nil, err
		}
		existingUser, _ := s.userRepo.FindByPhone(input.Phone)
		if existingUser != nil {
			return nil, nil, errors.New("手机号已被注册")
		}
	}

	if err := validator.ValidatePassword(input.Password); err != nil {
		return nil, nil, err
	}

	if err := validator.ValidateNickname(input.Nickname); err != nil {
		return nil, nil, err
	}

	hashedPassword, err := encryption.HashPassword(input.Password)
	if err != nil {
		return nil, nil, errors.New("密码加密失败")
	}

	user := &models.User{
		Email:      input.Email,
		Phone:      input.Phone,
		Password:   hashedPassword,
		Nickname:   input.Nickname,
		Role:       models.RoleLearner,
		Status:     models.UserStatusActive,
		AuthType:   models.AuthType(input.AuthType),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, errors.New("注册失败")
	}

	s.userRepo.UpdateLastLogin(user.ID, ip)

	token, err := auth.GenerateToken(user.ID, string(user.Role), user.Nickname, s.secret, 24, 168)
	if err != nil {
		return nil, nil, errors.New("生成令牌失败")
	}

	return user, token, nil
}

func (s *UserService) Login(input *LoginInput, ip string) (*models.User, *auth.TokenPair, error) {
	var user *models.User
	var err error

	if input.Email != "" {
		user, err = s.userRepo.FindByEmail(input.Email)
	} else if input.Phone != "" {
		user, err = s.userRepo.FindByPhone(input.Phone)
	}

	if err != nil || user == nil {
		return nil, nil, errors.New("用户不存在")
	}

	if user.Status != models.UserStatusActive {
		return nil, nil, errors.New("账号已被禁用")
	}

	if !encryption.VerifyPassword(input.Password, user.Password) {
		return nil, nil, errors.New("密码错误")
	}

	s.userRepo.UpdateLastLogin(user.ID, ip)

	token, err := auth.GenerateToken(user.ID, string(user.Role), user.Nickname, s.secret, 24, 168)
	if err != nil {
		return nil, nil, errors.New("生成令牌失败")
	}

	return user, token, nil
}

func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if nickname, ok := updates["nickname"].(string); ok {
		if err := validator.ValidateNickname(nickname); err != nil {
			return nil, err
		}
		user.Nickname = nickname
	}

	if avatar, ok := updates["avatar"].(string); ok {
		user.Avatar = avatar
	}

	if bio, ok := updates["bio"].(string); ok {
		user.Bio = bio
	}

	if gender, ok := updates["gender"].(string); ok {
		user.Gender = gender
	}

	if location, ok := updates["location"].(string); ok {
		user.Location = location
	}

	if latitude, ok := updates["latitude"].(float64); ok {
		user.Latitude = latitude
	}

	if longitude, ok := updates["longitude"].(float64); ok {
		user.Longitude = longitude
	}

	if birthday, ok := updates["birthday"].(string); ok {
		t, err := time.Parse("2006-01-02", birthday)
		if err == nil {
			user.Birthday = &t
		}
	}

	if role, ok := updates["role"].(string); ok {
		user.Role = models.UserRole(role)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新资料失败")
	}

	return user, nil
}

func (s *UserService) RefreshToken(refreshToken string) (*auth.TokenPair, error) {
	return auth.RefreshToken(refreshToken, s.secret, 24, 168)
}

func (s *UserService) ListUsers(page, pageSize int, keyword string) ([]*models.User, int64, error) {
	return s.userRepo.List(page, pageSize, keyword)
}

func (s *UserService) AddSkillTags(userID uuid.UUID, tagIDs []uuid.UUID) error {
	return s.userRepo.AddSkillTags(userID, tagIDs)
}

func (s *UserService) RemoveSkillTag(userID, tagID uuid.UUID) error {
	return s.userRepo.RemoveSkillTag(userID, tagID)
}

func (s *UserService) SubmitCertification(userID uuid.UUID, cert *models.Certification) error {
	cert.UserID = userID
	return nil
}
