package services

import (
	"errors"
	"qa-platform/models"
	"qa-platform/repository"

	"gorm.io/gorm"
)

type FavoriteService struct{}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

type FavoriteRequest struct {
	TargetType string `json:"targetType" binding:"required,oneof=question answer"`
	TargetID   uint   `json:"targetId" binding:"required"`
}

func (s *FavoriteService) AddFavorite(userID uint, req *FavoriteRequest) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var existingFav models.Favorite
		result := tx.Where("user_id = ? AND target_type = ? AND target_id = ?",
			userID, req.TargetType, req.TargetID).First(&existingFav)
		if result.Error == nil {
			return errors.New("已收藏")
		}

		fav := models.Favorite{
			UserID:     userID,
			TargetType: req.TargetType,
			TargetID:   req.TargetID,
		}

		if err := tx.Create(&fav).Error; err != nil {
			return err
		}

		switch req.TargetType {
		case "question":
			if err := tx.Model(&models.Question{}).Where("id = ?", req.TargetID).
				UpdateColumn("collect_count", gorm.Expr("collect_count + ?", 1)).Error; err != nil {
				return err
			}
		case "answer":
			if err := tx.Model(&models.Answer{}).Where("id = ?", req.TargetID).
				UpdateColumn("collect_count", gorm.Expr("collect_count + ?", 1)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *FavoriteService) RemoveFavorite(userID uint, req *FavoriteRequest) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ? AND target_type = ? AND target_id = ?",
			userID, req.TargetType, req.TargetID).Delete(&models.Favorite{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("未收藏")
		}

		switch req.TargetType {
		case "question":
			if err := tx.Model(&models.Question{}).Where("id = ?", req.TargetID).
				UpdateColumn("collect_count", gorm.Expr("GREATEST(collect_count - ?, 0)", 1)).Error; err != nil {
				return err
			}
		case "answer":
			if err := tx.Model(&models.Answer{}).Where("id = ?", req.TargetID).
				UpdateColumn("collect_count", gorm.Expr("GREATEST(collect_count - ?, 0)", 1)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *FavoriteService) GetUserFavorites(userID uint, targetType string, page, pageSize int) ([]models.Favorite, int64, error) {
	var favorites []models.Favorite
	var total int64

	dbQuery := repository.DB.Where("user_id = ?", userID)
	if targetType != "" {
		dbQuery = dbQuery.Where("target_type = ?", targetType)
	}

	dbQuery.Model(&models.Favorite{}).Count(&total)
	dbQuery.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&favorites)

	return favorites, total, nil
}

type FollowService struct{}

func NewFollowService() *FollowService {
	return &FollowService{}
}

type FollowRequest struct {
	FollowingType string `json:"followingType" binding:"required,oneof=user tag"`
	FollowingID   uint   `json:"followingId" binding:"required"`
}

func (s *FollowService) Follow(userID uint, req *FollowRequest) error {
	var existingFollow models.Follow
	result := repository.DB.Where("follower_id = ? AND following_type = ? AND following_id = ?",
		userID, req.FollowingType, req.FollowingID).First(&existingFollow)
	if result.Error == nil {
		return errors.New("已关注")
	}

	follow := models.Follow{
		FollowerID:    userID,
		FollowingType: req.FollowingType,
		FollowingID:   req.FollowingID,
	}

	return repository.DB.Create(&follow).Error
}

func (s *FollowService) Unfollow(userID uint, req *FollowRequest) error {
	result := repository.DB.Where("follower_id = ? AND following_type = ? AND following_id = ?",
		userID, req.FollowingType, req.FollowingID).Delete(&models.Follow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未关注")
	}
	return nil
}

func (s *FollowService) GetUserFollows(userID uint, followingType string, page, pageSize int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64

	dbQuery := repository.DB.Preload("Follower").Where("follower_id = ?", userID)
	if followingType != "" {
		dbQuery = dbQuery.Where("following_type = ?", followingType)
	}

	dbQuery.Model(&models.Follow{}).Count(&total)
	dbQuery.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&follows)

	return follows, total, nil
}

func (s *FollowService) GetUserFollowers(userID uint, page, pageSize int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64

	dbQuery := repository.DB.Preload("Follower").
		Where("following_type = ? AND following_id = ?", "user", userID)

	dbQuery.Model(&models.Follow{}).Count(&total)
	dbQuery.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&follows)

	return follows, total, nil
}

func (s *FollowService) IsFollowing(followerID uint, followingType string, followingID uint) bool {
	var follow models.Follow
	result := repository.DB.Where("follower_id = ? AND following_type = ? AND following_id = ?",
		followerID, followingType, followingID).First(&follow)
	return result.Error == nil
}

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) CreateNotification(userID uint, notifType, title, content, refType string, refID uint) error {
	notification := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Content: content,
		RefType: refType,
		RefID:   refID,
	}

	return repository.DB.Create(&notification).Error
}

func (s *NotificationService) GetUserNotifications(userID uint, isRead *bool, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	dbQuery := repository.DB.Where("user_id = ?", userID)
	if isRead != nil {
		dbQuery = dbQuery.Where("is_read = ?", *isRead)
	}

	dbQuery.Model(&models.Notification{}).Count(&total)
	dbQuery.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&notifications)

	return notifications, total, nil
}

func (s *NotificationService) MarkAsRead(userID, notificationID uint) error {
	result := repository.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)
	return result.Error
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	result := repository.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)
	return result.Error
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := repository.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}
