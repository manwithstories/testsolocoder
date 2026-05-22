package services

import (
	"errors"
	"qa-platform/models"
	"qa-platform/repository"
	"qa-platform/utils"
	"time"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	Level        int       `json:"level"`
	Points       int       `json:"points"`
	IsExpert     bool      `json:"isExpert"`
	ExpertStatus string    `json:"expertStatus"`
	Role         string    `json:"role"`
	Bio          string    `json:"bio"`
	LastLoginAt  time.Time `json:"lastLoginAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func (s *UserService) Register(req *RegisterRequest) (*LoginResponse, error) {
	var existingUser models.User
	result := repository.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("用户已存在")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     "user",
		Status:   "active",
		Level:    1,
		Points:   100,
	}

	if err := repository.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	pointLog := models.PointLog{
		UserID:      user.ID,
		Type:        "register",
		Points:      100,
		Balance:     100,
		Description: "注册奖励",
	}
	repository.DB.Create(&pointLog)

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, user.Level, user.IsExpert, 24)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user models.User
	result := repository.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != "active" {
		return nil, errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	now := time.Now()
	user.LastLoginAt = now
	repository.DB.Save(&user)

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, user.Level, user.IsExpert, 24)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (s *UserService) GetUserByID(id uint) (*UserResponse, error) {
	var user models.User
	if err := repository.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	resp := toUserResponse(user)
	return &resp, nil
}

func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*UserResponse, error) {
	var user models.User
	if err := repository.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if err := repository.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	repository.DB.First(&user, id)
	resp := toUserResponse(user)
	return &resp, nil
}

func (s *UserService) GetUserList(page, pageSize int, keyword string) ([]UserResponse, int64, error) {
	var users []models.User
	var total int64

	query := repository.DB.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&users)

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, toUserResponse(user))
	}

	return userResponses, total, nil
}

func (s *UserService) UpdateUserStatus(id uint, status string) error {
	return repository.DB.Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}

func (s *UserService) GetPointLogs(userID uint, page, pageSize int) ([]models.PointLog, int64, error) {
	var logs []models.PointLog
	var total int64

	query := repository.DB.Where("user_id = ?", userID)
	query.Model(&models.PointLog{}).Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&logs)

	return logs, total, nil
}

func (s *UserService) ApplyExpert(userID uint, field, description string) error {
	var existingApp models.ExpertApplication
	result := repository.DB.Where("user_id = ? AND status = ?", userID, "pending").First(&existingApp)
	if result.Error == nil {
		return errors.New("已有待审核的申请")
	}

	app := models.ExpertApplication{
		UserID:      userID,
		Field:       field,
		Description: description,
		Status:      "pending",
	}

	return repository.DB.Create(&app).Error
}

func (s *UserService) ReviewExpertApplication(id uint, reviewerID uint, status, remark string) error {
	var app models.ExpertApplication
	if err := repository.DB.First(&app, id).Error; err != nil {
		return errors.New("申请不存在")
	}

	now := time.Now()
	app.Status = status
	app.ReviewerID = &reviewerID
	app.ReviewRemark = remark
	app.ReviewedAt = &now

	if err := repository.DB.Save(&app).Error; err != nil {
		return err
	}

	if status == "approved" {
		repository.DB.Model(&models.User{}).Where("id = ?", app.UserID).Updates(map[string]interface{}{
			"is_expert":    true,
			"expert_status": "approved",
		})
	}

	return nil
}

func (s *UserService) GetExpertApplications(page, pageSize int, status string) ([]models.ExpertApplication, int64, error) {
	var apps []models.ExpertApplication
	var total int64

	query := repository.DB.Preload("User").Preload("Reviewer")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.ExpertApplication{}).Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&apps)

	return apps, total, nil
}

func toUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Level:        user.Level,
		Points:       user.Points,
		IsExpert:     user.IsExpert,
		ExpertStatus: user.ExpertStatus,
		Role:         user.Role,
		Bio:          user.Bio,
		LastLoginAt:  user.LastLoginAt,
		CreatedAt:    user.CreatedAt,
	}
}

func CalculateLevel(points int) int {
	levels := []int{0, 100, 500, 1000, 2000, 5000, 10000, 20000, 50000, 100000}
	for i := len(levels) - 1; i >= 0; i-- {
		if points >= levels[i] {
			return i + 1
		}
	}
	return 1
}

func AddPoints(userID uint, points int, logType, description string, refType string, refID uint) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		newPoints := user.Points + points
		if newPoints < 0 {
			return errors.New("积分不足")
		}

		newLevel := CalculateLevel(newPoints)

		if err := tx.Model(&user).Updates(map[string]interface{}{
			"points": newPoints,
			"level":  newLevel,
		}).Error; err != nil {
			return err
		}

		log := models.PointLog{
			UserID:      userID,
			Type:        logType,
			Points:      points,
			Balance:     newPoints,
			Description: description,
			RefType:     refType,
			RefID:       refID,
		}

		return tx.Create(&log).Error
	})
}
