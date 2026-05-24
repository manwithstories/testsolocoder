package repository

import (
	"music-platform/internal/model"
	"music-platform/pkg/database"
	"music-platform/pkg/utils"

	"gorm.io/gorm"
)

type CommunityRepository struct{}

func NewCommunityRepository() *CommunityRepository {
	return &CommunityRepository{}
}

func (r *CommunityRepository) Follow(followerID, followingID uint) error {
	follow := model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return database.DB.Create(&follow).Error
}

func (r *CommunityRepository) Unfollow(followerID, followingID uint) error {
	return database.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&model.Follow{}).Error
}

func (r *CommunityRepository) IsFollowing(followerID, followingID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *CommunityRepository) GetFollowers(userID uint, page, pageSize int) ([]model.Follow, int64, error) {
	var follows []model.Follow
	var total int64

	query := database.DB.Model(&model.Follow{}).Where("following_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Follower").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&follows).Error
	return follows, total, err
}

func (r *CommunityRepository) GetFollowings(userID uint, page, pageSize int) ([]model.Follow, int64, error) {
	var follows []model.Follow
	var total int64

	query := database.DB.Model(&model.Follow{}).Where("follower_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Following").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&follows).Error
	return follows, total, err
}

func (r *CommunityRepository) GetFollowerCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *CommunityRepository) GetFollowingCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *CommunityRepository) CreateComment(comment *model.Comment) error {
	return database.DB.Create(comment).Error
}

func (r *CommunityRepository) DeleteComment(id uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Comment{}).Error
}

func (r *CommunityRepository) GetComments(workID uint, albumID *uint, playlistID *uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := database.DB.Model(&model.Comment{}).Where("parent_id IS NULL")

	if workID > 0 {
		query = query.Where("work_id = ?", workID)
	}
	if albumID != nil {
		query = query.Where("album_id = ?", *albumID)
	}
	if playlistID != nil {
		query = query.Where("playlist_id = ?", *playlistID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("User").Preload("Replies.User").
		Offset(offset).Limit(pageSize).Order("is_pinned DESC, created_at DESC").Find(&comments).Error
	return comments, total, err
}

func (r *CommunityRepository) GetCommentByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := database.DB.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommunityRepository) UpdateCommentLikeCount(id uint, count int64) error {
	return database.DB.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", count)).Error
}

func (r *CommunityRepository) UpdateCommentReplyCount(id uint, count int64) error {
	return database.DB.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("reply_count", gorm.Expr("reply_count + ?", count)).Error
}

func (r *CommunityRepository) CreatePlaylist(playlist *model.Playlist) error {
	return database.DB.Create(playlist).Error
}

func (r *CommunityRepository) UpdatePlaylist(playlist *model.Playlist) error {
	return database.DB.Save(playlist).Error
}

func (r *CommunityRepository) DeletePlaylist(id uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Playlist{}).Error
}

func (r *CommunityRepository) GetPlaylistByID(id uint) (*model.Playlist, error) {
	var playlist model.Playlist
	err := database.DB.Preload("User").Preload("Works.Tags").First(&playlist, id).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *CommunityRepository) ListPlaylists(userID uint, page, pageSize int, keyword string) ([]model.Playlist, int64, error) {
	var playlists []model.Playlist
	var total int64

	query := database.DB.Model(&model.Playlist{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&playlists).Error
	return playlists, total, err
}

func (r *CommunityRepository) AddWorkToPlaylist(playlistID, workID uint) error {
	pw := model.PlaylistWork{
		PlaylistID: playlistID,
		WorkID:     workID,
	}
	return database.DB.Create(&pw).Error
}

func (r *CommunityRepository) RemoveWorkFromPlaylist(playlistID, workID uint) error {
	return database.DB.Where("playlist_id = ? AND work_id = ?", playlistID, workID).
		Delete(&model.PlaylistWork{}).Error
}

func (r *CommunityRepository) IsWorkInPlaylist(playlistID, workID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&model.PlaylistWork{}).
		Where("playlist_id = ? AND work_id = ?", playlistID, workID).
		Count(&count).Error
	return count > 0, err
}

func (r *CommunityRepository) UpdatePlaylistWorkCount(playlistID uint, count int) error {
	return database.DB.Model(&model.Playlist{}).Where("id = ?", playlistID).
		UpdateColumn("work_count", gorm.Expr("work_count + ?", count)).Error
}

func (r *CommunityRepository) CreatePlayRecord(record *model.PlayRecord) error {
	return database.DB.Create(record).Error
}

func (r *CommunityRepository) GetPlayRecords(userID uint, page, pageSize int) ([]model.PlayRecord, int64, error) {
	var records []model.PlayRecord
	var total int64

	query := database.DB.Model(&model.PlayRecord{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Work").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error
	return records, total, err
}

func (r *CommunityRepository) CreateNotification(notification *model.Notification) error {
	return database.DB.Create(notification).Error
}

func (r *CommunityRepository) GetNotifications(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := database.DB.Model(&model.Notification{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notifications).Error
	return notifications, total, err
}

func (r *CommunityRepository) MarkNotificationAsRead(id uint, userID uint) error {
	return database.DB.Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *CommunityRepository) MarkAllNotificationsAsRead(userID uint) error {
	return database.DB.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

func (r *CommunityRepository) GetUnreadNotificationCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}

func (r *CommunityRepository) GetPlaylistsByUser(userID uint) ([]model.Playlist, error) {
	var playlists []model.Playlist
	err := database.DB.Where("user_id = ?", userID).
		Preload("Works").
		Order("created_at DESC").
		Find(&playlists).Error
	return playlists, err
}

func (r *CommunityRepository) GetFollowingUsersByRole(userID uint, role string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	subQuery := database.DB.Model(&model.Follow{}).Select("following_id").Where("follower_id = ?", userID)

	query := database.DB.Model(&model.User{}).Where("id IN (?) AND role = ?", subQuery, role)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}
