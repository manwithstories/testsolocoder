package service

import (
	"drone-rental/internal/config"
	"drone-rental/internal/dto"
	"drone-rental/internal/middleware"
	"drone-rental/internal/model"
	"drone-rental/internal/repository"
	"errors"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepo(),
	}
}

func (s *UserService) Register(req *dto.RegisterReq) (*model.User, error) {
	existing, _ := s.userRepo.GetByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("用户名已存在")
	}
	if req.Phone != "" {
		existingPhone, _ := s.userRepo.GetByPhone(req.Phone)
		if existingPhone != nil {
			return nil, errors.New("手机号已注册")
		}
	}
	hashedPwd, err := repository.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username:     req.Username,
		Password:     hashedPwd,
		Nickname:     req.Nickname,
		Phone:        req.Phone,
		Email:        req.Email,
		Role:         req.Role,
		VerifyStatus: model.VerifyPending,
		Rating:       5.0,
		Status:       1,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(req *dto.LoginReq) (*dto.LoginResp, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil || user == nil {
		return nil, errors.New("用户名或密码错误")
	}
	if !repository.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}
	token, err := middleware.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResp{
		Token:        token,
		UserID:       user.ID,
		Username:     user.Username,
		Role:         user.Role,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		VerifyStatus: user.VerifyStatus,
	}, nil
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Update(id uint, req *dto.UpdateUserReq) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	return s.userRepo.Update(user)
}

func (s *UserService) VerifyPilot(id uint, req *dto.VerifyPilotReq, licenseImage string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	user.RealName = req.RealName
	user.IDCardNo = req.IDCardNo
	user.LicenseNo = req.LicenseNo
	user.LicenseImage = licenseImage
	user.VerifyStatus = model.VerifyPending
	return s.userRepo.Update(user)
}

func (s *UserService) VerifyOwner(id uint, req *dto.VerifyOwnerReq) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	user.RealName = req.RealName
	user.IDCardNo = req.IDCardNo
	user.VerifyStatus = model.VerifyPending
	return s.userRepo.Update(user)
}

func (s *UserService) List(page, pageSize int, role model.Role, keyword string) ([]model.User, int64, error) {
	return s.userRepo.List(page, pageSize, role, keyword)
}

func (s *UserService) AuditVerify(id uint, status model.VerifyStatus, remark string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	user.VerifyStatus = status
	user.VerifyRemark = remark
	return s.userRepo.Update(user)
}

func (s *UserService) UpdateBalance(id uint, amount float64) error {
	return s.userRepo.UpdateBalance(id, amount)
}

func (s *UserService) UpdateDeposit(id uint, amount float64) error {
	return s.userRepo.UpdateDeposit(id, amount)
}

func (s *UserService) UpdateRating(id uint) error {
	reviewRepo := repository.NewReviewRepo()
	avg, count, err := reviewRepo.GetAvgRatingByUser(id)
	if err != nil {
		return err
	}
	return config.DB.Model(&model.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"rating":       avg,
			"rating_count": count,
		}).Error
}
