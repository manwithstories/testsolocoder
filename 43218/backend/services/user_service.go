package services

import (
	"context"
	"errors"
	"time"

	"secondhand-platform/cache"
	"secondhand-platform/database"
	"secondhand-platform/models"
	"secondhand-platform/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(username, password, email, phone, role, nickname string) (*models.User, error) {
	if database.DB.Where("username = ?", username).First(&models.User{}).Error == nil {
		return nil, errors.New("用户名已存在")
	}

	if email != "" && database.DB.Where("email = ?", email).First(&models.User{}).Error == nil {
		return nil, errors.New("邮箱已被注册")
	}

	if phone != "" && database.DB.Where("phone = ?", phone).First(&models.User{}).Error == nil {
		return nil, errors.New("手机号已被注册")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:    username,
		Password:    hashedPassword,
		Email:       email,
		Phone:       phone,
		Role:        role,
		Nickname:    nickname,
		CreditScore: 100,
		Status:      models.UserStatusNormal,
	}

	result := database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *UserService) Login(username, password string) (*models.User, string, string, error) {
	var user models.User
	if err := database.DB.Where("username = ? OR email = ? OR phone = ?", username, username, username).First(&user).Error; err != nil {
		return nil, "", "", errors.New("用户不存在")
	}

	if user.Status != models.UserStatusNormal {
		return nil, "", "", errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, "", "", errors.New("密码错误")
	}

	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	accessToken, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	cache.Set(context.Background(), "user:online:"+utils.GenerateUUID(), user.ID, 24*time.Hour)

	return &user, accessToken, refreshToken, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateProfile(id uint, updates map[string]interface{}) error {
	result := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return database.DB.Model(&user).Update("password", hashedPassword).Error
}

func (s *UserService) SubmitRealNameAuth(userID uint, realName, idCard string) error {
	updates := map[string]interface{}{
		"real_name": realName,
		"id_card":   idCard,
	}
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserService) ReviewRealNameAuth(userID uint, approved bool, rejectReason string) error {
	updates := map[string]interface{}{
		"is_authenticated": approved,
	}
	if !approved && rejectReason != "" {
		updates["reject_reason"] = rejectReason
	}
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserService) UpdateCreditScore(userID uint, delta int) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	newScore := user.CreditScore + delta
	if newScore < 0 {
		newScore = 0
	}
	if newScore > 100 {
		newScore = 100
	}

	return database.DB.Model(&user).Update("credit_score", newScore).Error
}

func (s *UserService) ListUsers(page, pageSize int, role string, status int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := database.DB.Model(&models.User{})
	if role != "" {
		db = db.Where("role = ?", role)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) UpdateUserStatus(userID uint, status int) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (s *UserService) SubmitTechnicianCert(userID uint, certType, certNumber, certImage string) (*models.TechnicianCert, error) {
	var existingCert models.TechnicianCert
	if database.DB.Where("user_id = ? AND status = ?", userID, models.CertStatusPending).First(&existingCert).Error == nil {
		return nil, errors.New("已有审核中的认证申请")
	}

	cert := &models.TechnicianCert{
		UserID:     userID,
		CertType:   certType,
		CertNumber: certNumber,
		CertImage:  certImage,
		Status:     models.CertStatusPending,
	}

	result := database.DB.Create(cert)
	if result.Error != nil {
		return nil, result.Error
	}

	return cert, nil
}

func (s *UserService) ReviewTechnicianCert(certID uint, approved bool, rejectReason string, reviewerID uint) error {
	var cert models.TechnicianCert
	if err := database.DB.First(&cert, certID).Error; err != nil {
		return errors.New("认证申请不存在")
	}

	now := time.Now()
	cert.ReviewedAt = &now
	cert.ReviewedBy = &reviewerID

	if approved {
		cert.Status = models.CertStatusApproved
		database.DB.Model(&models.User{}).Where("id = ?", cert.UserID).Update("is_authenticated", true)
	} else {
		cert.Status = models.CertStatusRejected
		cert.RejectReason = rejectReason
	}

	return database.DB.Save(&cert).Error
}

func (s *UserService) GetWalletBalance(userID uint) (float64, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return 0, err
	}
	return user.WalletBalance, nil
}

func (s *UserService) Recharge(userID uint, amount float64, paymentMethod string) (*models.Transaction, error) {
	var transaction models.Transaction

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		newBalance := user.WalletBalance + amount
		if err := tx.Model(&user).Update("wallet_balance", newBalance).Error; err != nil {
			return err
		}

		transaction = models.Transaction{
			UserID:        userID,
			Type:          models.TransactionTypeRecharge,
			Amount:        amount,
			PaymentMethod: paymentMethod,
			TransactionNo: utils.GenerateOrderNo(),
			Status:        models.TransactionStatusSuccess,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		walletLog := models.WalletLog{
			UserID:      userID,
			Type:        models.WalletTypeRecharge,
			Amount:      amount,
			Balance:     newBalance,
			OrderNo:     transaction.TransactionNo,
			Description: "账户充值",
		}
		return tx.Create(&walletLog).Error
	})

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *UserService) Withdraw(userID uint, amount float64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		if user.WalletBalance < amount {
			return errors.New("余额不足")
		}

		newBalance := user.WalletBalance - amount
		if err := tx.Model(&user).Update("wallet_balance", newBalance).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			UserID:        userID,
			Type:          models.TransactionTypeWithdraw,
			Amount:        -amount,
			TransactionNo: utils.GenerateOrderNo(),
			Status:        models.TransactionStatusPending,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		walletLog := models.WalletLog{
			UserID:      userID,
			Type:        models.WalletTypeWithdraw,
			Amount:      -amount,
			Balance:     newBalance,
			OrderNo:     transaction.TransactionNo,
			Description: "账户提现",
		}
		return tx.Create(&walletLog).Error
	})
}

func (s *UserService) ListWalletLogs(userID uint, page, pageSize int) ([]models.WalletLog, int64, error) {
	var logs []models.WalletLog
	var total int64

	db := database.DB.Model(&models.WalletLog{}).Where("user_id = ?", userID)
	db.Count(&total)
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (s *UserService) Logout(token string) error {
	claims, err := utils.ParseToken(token)
	if err != nil {
		return err
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		return cache.Set(context.Background(), "token:blacklist:"+token, 1, ttl)
	}

	return nil
}

func (s *UserService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	newToken, err := utils.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (s *UserService) GetUserStats(userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var sellOrders int64
	database.DB.Model(&models.Order{}).Where("seller_id = ? AND status IN ?", userID, []int{
		models.OrderStatusPaid, models.OrderStatusShipped,
		models.OrderStatusDelivered, models.OrderStatusCompleted,
	}).Count(&sellOrders)
	stats["sell_orders"] = sellOrders

	var buyOrders int64
	database.DB.Model(&models.Order{}).Where("buyer_id = ? AND status IN ?", userID, []int{
		models.OrderStatusPaid, models.OrderStatusShipped,
		models.OrderStatusDelivered, models.OrderStatusCompleted,
	}).Count(&buyOrders)
	stats["buy_orders"] = buyOrders

	var products int64
	database.DB.Model(&models.Product{}).Where("seller_id = ?", userID).Count(&products)
	stats["products"] = products

	var reviews int64
	database.DB.Model(&models.Review{}).Where("reviewer_id = ?", userID).Count(&reviews)
	stats["reviews"] = reviews

	return stats, nil
}

func init() {
	logrus.Info("User service initialized")
}
