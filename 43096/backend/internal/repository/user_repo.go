package repository

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"ticket-system/internal/database"
	"ticket-system/internal/models"
	"ticket-system/internal/redis"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetByID(id uint64) (*models.User, error) {
	cached, err := redis.GetCachedUserInfo(id)
	if err == nil && cached != "" {
		var user models.User
		if json.Unmarshal([]byte(cached), &user) == nil {
			return &user, nil
		}
	}

	var user models.User
	err = database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	userJSON, _ := json.Marshal(user)
	_ = redis.CacheUserInfo(id, string(userJSON))

	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
	err := database.DB.Save(user).Error
	if err == nil {
		_ = redis.DeleteCachedUserInfo(user.ID)
	}
	return err
}

func (r *UserRepository) UpdatePoints(userID uint64, points int) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("points", gorm.Expr("points + ?", points)).Error
}

func (r *UserRepository) UpdateMemberLevel(userID uint64, level int) error {
	err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("member_level", level).Error
	if err == nil {
		_ = redis.DeleteCachedUserInfo(userID)
	}
	return err
}

func (r *UserRepository) CheckMemberLevelUpdate(userID uint64) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var levels []models.MemberLevel
	database.DB.Order("level desc").Find(&levels)

	for _, level := range levels {
		if user.Points >= level.MinPoints && user.MemberLevel != level.Level {
			return r.UpdateMemberLevel(userID, level.Level)
		}
	}

	return nil
}

func (r *UserRepository) GetMemberLevels() ([]models.MemberLevel, error) {
	var levels []models.MemberLevel
	err := database.DB.Order("level asc").Find(&levels).Error
	return levels, err
}

func (r *UserRepository) GetMemberLevelByLevel(level int) (*models.MemberLevel, error) {
	var memberLevel models.MemberLevel
	err := database.DB.Where("level = ?", level).First(&memberLevel).Error
	return &memberLevel, err
}

func (r *UserRepository) CreateCoupon(coupon *models.Coupon) error {
	return database.DB.Create(coupon).Error
}

func (r *UserRepository) GetCouponByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := database.DB.Where("code = ?", code).First(&coupon).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &coupon, err
}

func (r *UserRepository) UpdateCoupon(coupon *models.Coupon) error {
	return database.DB.Save(coupon).Error
}

func (r *UserRepository) GetUserCoupons(userID uint64) ([]models.Coupon, error) {
	var coupons []models.Coupon
	err := database.DB.Where("user_id = ? AND status = 1", userID).Find(&coupons).Error
	return coupons, err
}
