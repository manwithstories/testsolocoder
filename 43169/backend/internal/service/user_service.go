package service

import (
	"errors"
	"strings"
	"time"

	"matchmaking-platform/internal/dto"
	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/repository"
	"matchmaking-platform/internal/utils"
)

type UserService struct {
	userRepo    *repository.UserRepo
	profileRepo *repository.ProfileRepo
	memberRepo  *repository.MemberRepo
	logRepo     *repository.SystemLogRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:    repository.NewUserRepo(),
		profileRepo: repository.NewProfileRepo(),
		memberRepo:  repository.NewMemberRepo(),
		logRepo:     repository.NewSystemLogRepo(),
	}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	if !utils.ValidatePhone(req.Phone) {
		return nil, errors.New("手机号格式不正确")
	}

	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	if _, err := s.userRepo.FindByPhone(req.Phone); err == nil {
		return nil, errors.New("手机号已注册")
	}

	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		Password:     hashedPwd,
		Phone:        req.Phone,
		Email:        req.Email,
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		VerifyStatus: model.VerifyStatusPending,
		MemberLevel:  string(model.MemberFree),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	profile := &model.Profile{
		UserID: user.ID,
	}
	s.profileRepo.Create(profile)

	s.logRepo.Create(&model.SystemLog{
		UserID: user.ID,
		Module: "user",
		Action: "register",
	})

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *s.toUserInfo(user, profile),
	}, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *model.User
	var err error

	if strings.Contains(req.Account, "@") {
		user, err = s.userRepo.FindByEmail(req.Account)
	} else if utils.ValidatePhone(req.Account) {
		user, err = s.userRepo.FindByPhone(req.Account)
	} else {
		user, err = s.userRepo.FindByUsername(req.Account)
	}

	if err != nil {
		return nil, errors.New("账号或密码错误")
	}

	if user.Status != model.UserStatusActive {
		return nil, errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("账号或密码错误")
	}

	now := time.Now()
	s.userRepo.Update(user.ID, map[string]interface{}{
		"last_login_at": now,
	})

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}

	s.logRepo.Create(&model.SystemLog{
		UserID: user.ID,
		Module: "user",
		Action: "login",
	})

	return &dto.LoginResponse{
		Token: token,
		User:  *s.toUserInfo(user, user.Profile),
	}, nil
}

func (s *UserService) GetUserInfo(userID uint) (*dto.UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return s.toUserInfo(user, user.Profile), nil
}

func (s *UserService) Verify(userID uint, req *dto.VerifyRequest) error {
	if !utils.ValidateIDCard(req.IDCard) {
		return errors.New("身份证号格式不正确")
	}
	if !utils.ValidatePhone(req.Phone) {
		return errors.New("手机号格式不正确")
	}

	updates := map[string]interface{}{
		"real_name":    req.RealName,
		"id_card":      req.IDCard,
		"id_card_front": req.IDCardFront,
		"id_card_back":  req.IDCardBack,
		"phone":         req.Phone,
		"verify_status": model.VerifyStatusPending,
	}

	if err := s.userRepo.Update(userID, updates); err != nil {
		return err
	}

	s.logRepo.Create(&model.SystemLog{
		UserID: userID,
		Module: "user",
		Action: "verify_submit",
	})
	return nil
}

func (s *UserService) ApproveVerify(userID uint) error {
	return s.userRepo.UpdateVerifyStatus(userID, model.VerifyStatusVerified)
}

func (s *UserService) RejectVerify(userID uint) error {
	return s.userRepo.UpdateVerifyStatus(userID, model.VerifyStatusRejected)
}

func (s *UserService) UpdateProfile(userID uint, req *dto.ProfileUpdateRequest) error {
	updates := map[string]interface{}{
		"nickname":         req.Nickname,
		"gender":           req.Gender,
		"birthday":         req.Birthday,
		"age":              utils.CalcAge(req.Birthday),
		"height":           req.Height,
		"weight":           req.Weight,
		"education":        req.Education,
		"occupation":       req.Occupation,
		"income":           req.Income,
		"city":             req.City,
		"district":         req.District,
		"address":          req.Address,
		"intro":            req.Intro,
		"hobbies":          req.Hobbies,
		"tags":             req.Tags,
		"min_age":          req.MinAge,
		"max_age":          req.MaxAge,
		"min_height":       req.MinHeight,
		"max_height":       req.MaxHeight,
		"prefer_education": req.PreferEducation,
		"prefer_income":    req.PreferIncome,
		"prefer_city":      req.PreferCity,
	}
	return s.profileRepo.Update(userID, updates)
}

func (s *UserService) UploadAvatar(userID uint, avatarURL string) error {
	return s.userRepo.Update(userID, map[string]interface{}{"avatar": avatarURL})
}

func (s *UserService) UploadPhotos(userID uint, photos []string) error {
	return s.profileRepo.Update(userID, map[string]interface{}{"photos": strings.Join(photos, ",")})
}

func (s *UserService) ListUsers(page, pageSize int, keyword string) ([]dto.UserInfo, int64, error) {
	users, total, err := s.userRepo.List(page, pageSize, keyword)
	if err != nil {
		return nil, 0, err
	}
	var infos []dto.UserInfo
	for _, u := range users {
		infos = append(infos, *s.toUserInfo(&u, u.Profile))
	}
	return infos, total, nil
}

func (s *UserService) DisableUser(userID uint) error {
	return s.userRepo.UpdateStatus(userID, model.UserStatusDisabled)
}

func (s *UserService) EnableUser(userID uint) error {
	return s.userRepo.UpdateStatus(userID, model.UserStatusActive)
}

func (s *UserService) toUserInfo(user *model.User, profile *model.Profile) *dto.UserInfo {
	info := &dto.UserInfo{
		ID:           user.ID,
		Username:     user.Username,
		Phone:        user.Phone,
		Email:        user.Email,
		Role:         string(user.Role),
		Status:       string(user.Status),
		VerifyStatus: string(user.VerifyStatus),
		RealName:     user.RealName,
		Avatar:       user.Avatar,
		MemberLevel:  user.MemberLevel,
		MemberExpire: user.MemberExpire,
	}
	if profile != nil {
		info.Profile = &dto.ProfileInfo{
			ID:         profile.ID,
			UserID:     profile.UserID,
			Nickname:   profile.Nickname,
			Gender:     string(profile.Gender),
			Birthday:   profile.Birthday,
			Age:        profile.Age,
			Height:     profile.Height,
			Weight:     profile.Weight,
			Education:  string(profile.Education),
			Occupation: profile.Occupation,
			Income:     string(profile.Income),
			City:       profile.City,
			District:   profile.District,
			Intro:      profile.Intro,
			Hobbies:    profile.Hobbies,
			Tags:       profile.Tags,
			Photos:     strings.Split(profile.Photos, ","),
		}
	}
	return info
}
