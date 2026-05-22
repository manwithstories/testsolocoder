package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"car-rental/internal/utils"
	"errors"
	"time"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	RealName   string `json:"real_name"`
	AuthStatus string `json:"auth_status"`
	Status     string `json:"status"`
	RoleID     uint   `json:"role_id"`
	RoleName   string `json:"role_name"`
	Avatar     string `json:"avatar"`
}

func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
	if s.userRepo.ExistsByUsername(req.Username) {
		return nil, errors.New("用户名已存在")
	}
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("邮箱已被注册")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:   req.Username,
		Password:   hashedPassword,
		Email:      req.Email,
		Phone:      req.Phone,
		RoleID:     2,
		AuthStatus: model.UserStatusPending,
		Status:     model.UserStatusActive,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *LoginRequest, accessExpire, refreshExpire int) (*LoginResponse, error) {
	var user *model.User
	var err error

	if utils.ValidateEmail(req.Account) {
		user, err = s.userRepo.FindByEmail(req.Account)
	} else {
		user, err = s.userRepo.FindByUsername(req.Account)
	}

	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status == model.UserStatusDisabled {
		return nil, errors.New("账户已被禁用")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleID, accessExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, refreshExpire)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		return nil, err
	}

	roleName := "用户"
	if user.Role != nil {
		roleName = user.Role.Name
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessExpire * 3600,
		User: &UserInfo{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Phone:      user.Phone,
			RealName:   user.RealName,
			AuthStatus: string(user.AuthStatus),
			Status:     string(user.Status),
			RoleID:     user.RoleID,
			RoleName:   roleName,
			Avatar:     user.Avatar,
		},
	}, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetAllUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.userRepo.FindAll(page, pageSize, keyword)
}

func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if realName, ok := updates["real_name"]; ok {
		user.RealName = realName.(string)
	}
	if phone, ok := updates["phone"]; ok {
		user.Phone = phone.(string)
	}
	if avatar, ok := updates["avatar"]; ok {
		user.Avatar = avatar.(string)
	}
	if licenseNumber, ok := updates["license_number"]; ok {
		user.LicenseNumber = licenseNumber.(string)
	}
	if licenseImage, ok := updates["license_image"]; ok {
		user.LicenseImage = licenseImage.(string)
	}
	if idCard, ok := updates["id_card"]; ok {
		user.IDCard = idCard.(string)
	}
	if idCardFront, ok := updates["id_card_front"]; ok {
		user.IDCardFront = idCardFront.(string)
	}
	if idCardBack, ok := updates["id_card_back"]; ok {
		user.IDCardBack = idCardBack.(string)
	}

	if user.IDCard != "" && user.LicenseNumber != "" {
		user.AuthStatus = model.UserStatusPending
	}

	return s.userRepo.Update(user)
}

func (s *UserService) UpdateAuthStatus(id uint, status model.UserStatus) error {
	return s.userRepo.UpdateAuthStatus(id, status)
}

func (s *UserService) UpdateStatus(id uint, status model.UserStatus) error {
	return s.userRepo.UpdateStatus(id, status)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(user.Password, oldPassword) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *UserService) RefreshToken(tokenString string, accessExpire, refreshExpire int) (*LoginResponse, error) {
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return nil, errors.New("Token无效")
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, errors.New("Refresh Token已过期")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleID, accessExpire)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, refreshExpire)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    accessExpire * 3600,
	}, nil
}