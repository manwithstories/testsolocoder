package service

import (
	"errors"
	"time"

	"music-platform/internal/model"
	"music-platform/internal/repository"
	"music-platform/pkg/jwt"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/utils"
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
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UpdateArtistInfoRequest struct {
	RealName       string `json:"real_name"`
	IDCard         string `json:"id_card"`
	ArtistName     string `json:"artist_name"`
	Genre          string `json:"genre"`
	Label          string `json:"label"`
	Website        string `json:"website"`
	Facebook       string `json:"facebook"`
	Instagram      string `json:"instagram"`
	Twitter        string `json:"twitter"`
	YouTube        string `json:"youtube"`
	Spotify        string `json:"spotify"`
	AppleMusic     string `json:"apple_music"`
	BankAccount    string `json:"bank_account"`
	BankName       string `json:"bank_name"`
	BankHolder     string `json:"bank_holder"`
	AlipayAccount  string `json:"alipay_account"`
	WechatAccount  string `json:"wechat_account"`
}

func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	if req.Email != "" {
		existingUser, _ = s.userRepo.FindByEmail(req.Email)
		if existingUser != nil {
			return nil, apperrors.ErrUserAlreadyExists
		}
	}

	if req.Phone != "" {
		existingUser, _ = s.userRepo.FindByPhone(req.Phone)
		if existingUser != nil {
			return nil, apperrors.ErrUserAlreadyExists
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.ErrSystemError
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     model.UserRole(req.Role),
		Status:   model.UserStatusActive,
	}

	if user.Role == "" {
		user.Role = model.RoleFan
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	if user.Role == model.RoleArtist || user.Role == model.RoleLabel {
		artistInfo := &model.ArtistInfo{
			UserID:      user.ID,
			ArtistName:  req.Nickname,
			IsVerified:  false,
		}
		_ = s.userRepo.CreateArtistInfo(artistInfo)
	}

	return user, nil
}

func (s *UserService) Login(req *LoginRequest, ip string) (*model.User, string, error) {
	var user *model.User
	var err error

	if utils.IsValidEmail(req.Account) {
		user, err = s.userRepo.FindByEmail(req.Account)
	} else if utils.IsValidPhone(req.Account) {
		user, err = s.userRepo.FindByPhone(req.Account)
	} else {
		user, err = s.userRepo.FindByUsername(req.Account)
	}

	if err != nil || user == nil {
		return nil, "", apperrors.ErrUserNotFound
	}

	if user.Status != model.UserStatusActive {
		return nil, "", errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, "", apperrors.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, "", err
	}

	_ = s.userRepo.UpdateLoginInfo(user.ID, ip)

	return user, token, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) error {
	updates := map[string]interface{}{}

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	return s.userRepo.UpdateUserInfo(userID, updates)
}

func (s *UserService) UpdateArtistInfo(userID uint, req *UpdateArtistInfoRequest) error {
	info, err := s.userRepo.FindArtistInfoByUserID(userID)
	if err != nil {
		info = &model.ArtistInfo{
			UserID: userID,
		}
	}

	if req.RealName != "" {
		info.RealName = req.RealName
	}
	if req.IDCard != "" {
		info.IDCard = req.IDCard
	}
	if req.ArtistName != "" {
		info.ArtistName = req.ArtistName
	}
	if req.Genre != "" {
		info.Genre = req.Genre
	}
	if req.Label != "" {
		info.Label = req.Label
	}
	if req.Website != "" {
		info.Website = req.Website
	}
	if req.Facebook != "" {
		info.Facebook = req.Facebook
	}
	if req.Instagram != "" {
		info.Instagram = req.Instagram
	}
	if req.Twitter != "" {
		info.Twitter = req.Twitter
	}
	if req.YouTube != "" {
		info.YouTube = req.YouTube
	}
	if req.Spotify != "" {
		info.Spotify = req.Spotify
	}
	if req.AppleMusic != "" {
		info.AppleMusic = req.AppleMusic
	}
	if req.BankAccount != "" {
		info.BankAccount = req.BankAccount
	}
	if req.BankName != "" {
		info.BankName = req.BankName
	}
	if req.BankHolder != "" {
		info.BankHolder = req.BankHolder
	}
	if req.AlipayAccount != "" {
		info.AlipayAccount = req.AlipayAccount
	}
	if req.WechatAccount != "" {
		info.WechatAccount = req.WechatAccount
	}

	if info.ID == 0 {
		return s.userRepo.CreateArtistInfo(info)
	}
	return s.userRepo.UpdateArtistInfo(info)
}

func (s *UserService) ListUsers(page, pageSize int, keyword string, role string) ([]model.User, int64, error) {
	return s.userRepo.List(page, pageSize, keyword, role)
}

func (s *UserService) UpdateUserStatus(userID uint, status model.UserStatus) error {
	return s.userRepo.UpdateUserInfo(userID, map[string]interface{}{
		"status": status,
	})
}

func (s *UserService) UpdateUserRole(userID uint, role string) error {
	return s.userRepo.UpdateUserInfo(userID, map[string]interface{}{
		"role": model.UserRole(role),
	})
}

func (s *UserService) VerifyArtist(userID uint) error {
	info, err := s.userRepo.FindArtistInfoByUserID(userID)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	now := time.Now()
	info.IsVerified = true
	info.VerifiedAt = &now

	return s.userRepo.UpdateArtistInfo(info)
}

func (s *UserService) GetArtistInfo(userID uint) (*model.ArtistInfo, error) {
	return s.userRepo.FindArtistInfoByUserID(userID)
}

func (s *UserService) GetBalance(userID uint) (float64, float64, error) {
	info, err := s.userRepo.FindArtistInfoByUserID(userID)
	if err != nil {
		return 0, 0, apperrors.ErrUserNotFound
	}
	return info.Balance, info.FrozenBalance, nil
}
