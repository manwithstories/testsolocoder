package service

import (
	"errors"
	"time"

	"ticket-system/internal/dto"
	"ticket-system/internal/logging"
	"ticket-system/internal/models"
	"ticket-system/internal/repository"
	"ticket-system/internal/util"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*models.User, error) {
	_, err := s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	_, err = s.userRepo.GetByPhone(req.Phone)
	if err == nil {
		return nil, errors.New("手机号已注册")
	}

	if req.Email != "" {
		_, err = s.userRepo.GetByEmail(req.Email)
		if err == nil {
			return nil, errors.New("邮箱已注册")
		}
	}

	if !util.ValidatePhone(req.Phone) {
		return nil, errors.New("手机号格式不正确")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:    req.Username,
		Password:    hashedPassword,
		Phone:       req.Phone,
		Email:       req.Email,
		MemberLevel: 1,
		Points:      0,
		Status:      1,
		Role:        "user",
	}

	err = s.userRepo.Create(user)
	if err != nil {
		logging.Errorf("Failed to create user: %v", err)
		return nil, errors.New("注册失败")
	}

	logging.Infof("User registered successfully: %s", user.Username)
	return user, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := util.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	logging.Infof("User logged in: %s", user.Username)
	return &dto.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetUserByID(id uint64) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateUser(id uint64, req *dto.UserUpdateRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.RealName != "" {
		if !util.ValidateRealName(req.RealName) {
			return nil, errors.New("真实姓名格式不正确")
		}
		user.RealName = req.RealName
	}

	if req.IDCard != "" {
		if !util.ValidateIDCard(req.IDCard) {
			return nil, errors.New("身份证号格式不正确")
		}
		user.IDCard = req.IDCard
	}

	if req.Phone != "" {
		if !util.ValidatePhone(req.Phone) {
			return nil, errors.New("手机号格式不正确")
		}
		user.Phone = req.Phone
	}

	if req.Email != "" {
		if !util.ValidateEmail(req.Email) {
			return nil, errors.New("邮箱格式不正确")
		}
		user.Email = req.Email
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	logging.Infof("User updated: %d", id)
	return user, nil
}

func (s *UserService) GetMemberLevels() ([]models.MemberLevel, error) {
	return s.userRepo.GetMemberLevels()
}

func (s *UserService) ExchangeCoupon(userID uint64, points int) (*models.Coupon, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Points < points {
		return nil, errors.New("积分不足")
	}

	couponValue := float64(points) / 100

	coupon := &models.Coupon{
		Code:      util.GenerateCouponCode(),
		Name:      "积分兑换优惠券",
		Type:      1,
		Value:     couponValue,
		MinAmount: 0,
		UserID:    userID,
		Status:    1,
		ExpireAt:  time.Now().AddDate(0, 3, 0),
	}

	err = s.userRepo.CreateCoupon(coupon)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdatePoints(userID, -points)
	if err != nil {
		return nil, err
	}

	_ = s.userRepo.CheckMemberLevelUpdate(userID)

	logging.Infof("Coupon exchanged: user=%d, points=%d", userID, points)
	return coupon, nil
}

func (s *UserService) GetUserCoupons(userID uint64) ([]models.Coupon, error) {
	return s.userRepo.GetUserCoupons(userID)
}

func (s *UserService) AddPoints(userID uint64, points int) error {
	err := s.userRepo.UpdatePoints(userID, points)
	if err != nil {
		return err
	}
	_ = s.userRepo.CheckMemberLevelUpdate(userID)
	return nil
}

func (s *UserService) GetMemberDiscount(userID uint64) (float64, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 1.0, err
	}

	level, err := s.userRepo.GetMemberLevelByLevel(user.MemberLevel)
	if err != nil {
		return 1.0, nil
	}

	return level.Discount, nil
}

func (s *UserService) ChangePassword(userID uint64, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if !util.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	if len(req.NewPassword) < 6 || len(req.NewPassword) > 20 {
		return errors.New("新密码长度必须在6-20位之间")
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = hashedPassword
	err = s.userRepo.Update(user)
	if err != nil {
		return errors.New("修改密码失败")
	}

	logging.Infof("User changed password: %d", userID)
	return nil
}
