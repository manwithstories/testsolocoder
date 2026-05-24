package services

import (
	"errors"
	"time"

	"recruitment-platform/internal/models"
	"recruitment-platform/internal/repository"
	"recruitment-platform/internal/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required,oneof=company applicant"`
	Company   *CompanyRegister `json:"company,omitempty"`
	FullName  string `json:"full_name" binding:"required"`
	Phone     string `json:"phone"`
}

type CompanyRegister struct {
	Name     string `json:"name" binding:"required"`
	Industry string `json:"industry"`
	Size     string `json:"size"`
	Address  string `json:"address"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("该邮箱已注册")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	var companyID *uint
	if req.Role == "company" && req.Company != nil {
		company := &models.Company{
			Name:     req.Company.Name,
			Industry: req.Company.Industry,
			Size:     req.Company.Size,
			Address:  req.Company.Address,
		}
		if err := s.userRepo.CreateCompany(company); err != nil {
			return nil, errors.New("创建公司信息失败")
		}
		id := company.ID
		companyID = &id
	}

	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      models.UserRole(req.Role),
		Status:    models.UserStatusActive,
		CompanyID: companyID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	if req.Role == "applicant" {
		profile := &models.ApplicantProfile{
			UserID:   user.ID,
			FullName: req.FullName,
			Phone:    req.Phone,
		}
		if err := s.userRepo.CreateProfile(profile); err != nil {
			return nil, errors.New("创建用户档案失败")
		}
		user.Profile = profile
	}

	return user, nil
}

func (s *UserService) Login(req *LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != models.UserStatusActive {
		if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
			return nil, errors.New("账户已被锁定，请稍后再试")
		}
		return nil, errors.New("账户状态异常")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		user.LoginAttempts++
		if user.LoginAttempts >= 5 {
			s.userRepo.LockAccount(user.ID, time.Now().Add(30*time.Minute))
			return nil, errors.New("登录失败次数过多，账户已锁定30分钟")
		}
		s.userRepo.UpdateLoginAttempts(user.ID, user.LoginAttempts)
		return nil, errors.New("密码错误")
	}

	user.LoginAttempts = 0
	user.LockedUntil = nil
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateProfile(userID uint, profile *models.ApplicantProfile) error {
	existing, err := s.userRepo.FindProfileByUserID(userID)
	if err != nil {
		return s.userRepo.CreateProfile(profile)
	}
	profile.ID = existing.ID
	return s.userRepo.UpdateProfile(profile)
}

func (s *UserService) ListUsers(page, pageSize int, role, status string) ([]models.User, int64, error) {
	return s.userRepo.List(page, pageSize, models.UserRole(role), models.UserStatus(status))
}

func (s *UserService) UpdateUserStatus(id uint, status models.UserStatus) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	user.Status = status
	return s.userRepo.Update(user)
}

func (s *UserService) UpdateCompany(userID uint, company *models.Company) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user.CompanyID == nil {
		return errors.New("用户未关联公司")
	}
	company.ID = *user.CompanyID
	return s.userRepo.UpdateCompany(company)
}
