package service

import (
	"music-platform/internal/model"
	"music-platform/internal/repository"
	apperrors "music-platform/pkg/errors"
)

type CommunityService struct {
	communityRepo *repository.CommunityRepository
	userRepo      *repository.UserRepository
	workRepo      *repository.WorkRepository
}

func NewCommunityService() *CommunityService {
	return &CommunityService{
		communityRepo: repository.NewCommunityRepository(),
		userRepo:      repository.NewUserRepository(),
		workRepo:      repository.NewWorkRepository(),
	}
}

type FollowRequest struct {
	FollowingID uint `json:"following_id" binding:"required"`
}

type CreateCommentRequest struct {
	WorkID     *uint  `json:"work_id"`
	AlbumID    *uint  `json:"album_id"`
	PlaylistID *uint  `json:"playlist_id"`
	ParentID   *uint  `json:"parent_id"`
	Content    string `json:"content" binding:"required,max=500"`
}

type CreatePlaylistRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	IsPublic    bool   `json:"is_public"`
}

type UpdatePlaylistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	IsPublic    *bool  `json:"is_public"`
}

type PlaylistWorkRequest struct {
	PlaylistID uint `json:"playlist_id" binding:"required"`
	WorkID     uint `json:"work_id" binding:"required"`
}

func (s *CommunityService) Follow(userID uint, req *FollowRequest) error {
	if userID == req.FollowingID {
		return apperrors.NewAppError(3001, "不能关注自己")
	}

	user, err := s.userRepo.FindByID(req.FollowingID)
	if err != nil || user == nil {
		return apperrors.ErrUserNotFound
	}

	isFollowing, _ := s.communityRepo.IsFollowing(userID, req.FollowingID)
	if isFollowing {
		return apperrors.NewAppError(3002, "已关注该用户")
	}

	err = s.communityRepo.Follow(userID, req.FollowingID)
	if err != nil {
		return err
	}

	notification := &model.Notification{
		UserID:  req.FollowingID,
		Type:    "follow",
		Title:   "新的关注者",
		Content: "有人关注了你",
	}
	_ = s.communityRepo.CreateNotification(notification)

	return nil
}

func (s *CommunityService) Unfollow(userID uint, req *FollowRequest) error {
	return s.communityRepo.Unfollow(userID, req.FollowingID)
}

func (s *CommunityService) IsFollowing(followerID, followingID uint) (bool, error) {
	return s.communityRepo.IsFollowing(followerID, followingID)
}

func (s *CommunityService) GetFollowers(userID uint, page, pageSize int) ([]model.Follow, int64, error) {
	return s.communityRepo.GetFollowers(userID, page, pageSize)
}

func (s *CommunityService) GetFollowings(userID uint, page, pageSize int) ([]model.Follow, int64, error) {
	return s.communityRepo.GetFollowings(userID, page, pageSize)
}

func (s *CommunityService) GetFollowerCount(userID uint) (int64, error) {
	return s.communityRepo.GetFollowerCount(userID)
}

func (s *CommunityService) GetFollowingCount(userID uint) (int64, error) {
	return s.communityRepo.GetFollowingCount(userID)
}

func (s *CommunityService) CreateComment(userID uint, req *CreateCommentRequest) (*model.Comment, error) {
	comment := &model.Comment{
		UserID:     userID,
		WorkID:     *req.WorkID,
		AlbumID:    req.AlbumID,
		PlaylistID: req.PlaylistID,
		ParentID:   req.ParentID,
		Content:    req.Content,
	}

	err := s.communityRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	if req.WorkID != nil && *req.WorkID > 0 {
		_ = s.workRepo.UpdateCommentCount(*req.WorkID, 1)
	}

	if req.ParentID != nil && *req.ParentID > 0 {
		_ = s.communityRepo.UpdateCommentReplyCount(*req.ParentID, 1)
	}

	if req.WorkID != nil && *req.WorkID > 0 {
		work, _ := s.workRepo.FindByID(*req.WorkID)
		if work != nil {
			notification := &model.Notification{
				UserID:  work.UserID,
				Type:    "comment",
				Title:   "新的评论",
				Content: "你的作品有新评论",
			}
			_ = s.communityRepo.CreateNotification(notification)
		}
	}

	return comment, nil
}

func (s *CommunityService) DeleteComment(commentID uint, userID uint) error {
	comment, err := s.communityRepo.GetCommentByID(commentID)
	if err != nil {
		return apperrors.NewAppError(3003, "评论不存在")
	}

	if comment.UserID != userID {
		return apperrors.ErrForbidden
	}

	return s.communityRepo.DeleteComment(commentID, userID)
}

func (s *CommunityService) GetComments(workID uint, albumID *uint, playlistID *uint, page, pageSize int) ([]model.Comment, int64, error) {
	return s.communityRepo.GetComments(workID, albumID, playlistID, page, pageSize)
}

func (s *CommunityService) CreatePlaylist(userID uint, req *CreatePlaylistRequest) (*model.Playlist, error) {
	playlist := &model.Playlist{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		IsPublic:    req.IsPublic,
	}

	err := s.communityRepo.CreatePlaylist(playlist)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func (s *CommunityService) UpdatePlaylist(playlistID uint, userID uint, req *UpdatePlaylistRequest) error {
	playlist, err := s.communityRepo.GetPlaylistByID(playlistID)
	if err != nil {
		return apperrors.ErrPlaylistNotFound
	}

	if playlist.UserID != userID {
		return apperrors.ErrForbidden
	}

	if req.Title != "" {
		playlist.Title = req.Title
	}
	if req.Description != "" {
		playlist.Description = req.Description
	}
	if req.CoverURL != "" {
		playlist.CoverURL = req.CoverURL
	}
	if req.IsPublic != nil {
		playlist.IsPublic = *req.IsPublic
	}

	return s.communityRepo.UpdatePlaylist(playlist)
}

func (s *CommunityService) DeletePlaylist(playlistID uint, userID uint) error {
	playlist, err := s.communityRepo.GetPlaylistByID(playlistID)
	if err != nil {
		return apperrors.ErrPlaylistNotFound
	}

	if playlist.UserID != userID {
		return apperrors.ErrForbidden
	}

	return s.communityRepo.DeletePlaylist(playlistID, userID)
}

func (s *CommunityService) GetPlaylistByID(id uint) (*model.Playlist, error) {
	playlist, err := s.communityRepo.GetPlaylistByID(id)
	if err != nil {
		return nil, apperrors.ErrPlaylistNotFound
	}
	return playlist, nil
}

func (s *CommunityService) ListPlaylists(userID uint, page, pageSize int, keyword string) ([]model.Playlist, int64, error) {
	return s.communityRepo.ListPlaylists(userID, page, pageSize, keyword)
}

func (s *CommunityService) AddWorkToPlaylist(req *PlaylistWorkRequest, userID uint) error {
	playlist, err := s.communityRepo.GetPlaylistByID(req.PlaylistID)
	if err != nil {
		return apperrors.ErrPlaylistNotFound
	}

	if playlist.UserID != userID {
		return apperrors.ErrForbidden
	}

	_, err = s.workRepo.FindByID(req.WorkID)
	if err != nil {
		return apperrors.ErrWorkNotFound
	}

	isInPlaylist, _ := s.communityRepo.IsWorkInPlaylist(req.PlaylistID, req.WorkID)
	if isInPlaylist {
		return apperrors.NewAppError(3004, "作品已在歌单中")
	}

	err = s.communityRepo.AddWorkToPlaylist(req.PlaylistID, req.WorkID)
	if err != nil {
		return err
	}

	_ = s.communityRepo.UpdatePlaylistWorkCount(req.PlaylistID, 1)

	return nil
}

func (s *CommunityService) RemoveWorkFromPlaylist(req *PlaylistWorkRequest, userID uint) error {
	playlist, err := s.communityRepo.GetPlaylistByID(req.PlaylistID)
	if err != nil {
		return apperrors.ErrPlaylistNotFound
	}

	if playlist.UserID != userID {
		return apperrors.ErrForbidden
	}

	err = s.communityRepo.RemoveWorkFromPlaylist(req.PlaylistID, req.WorkID)
	if err != nil {
		return err
	}

	_ = s.communityRepo.UpdatePlaylistWorkCount(req.PlaylistID, -1)

	return nil
}

func (s *CommunityService) GetNotifications(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	return s.communityRepo.GetNotifications(userID, page, pageSize)
}

func (s *CommunityService) MarkNotificationAsRead(notificationID uint, userID uint) error {
	return s.communityRepo.MarkNotificationAsRead(notificationID, userID)
}

func (s *CommunityService) MarkAllNotificationsAsRead(userID uint) error {
	return s.communityRepo.MarkAllNotificationsAsRead(userID)
}

func (s *CommunityService) GetUnreadNotificationCount(userID uint) (int64, error) {
	return s.communityRepo.GetUnreadNotificationCount(userID)
}

func (s *CommunityService) GetPlayRecords(userID uint, page, pageSize int) ([]model.PlayRecord, int64, error) {
	return s.communityRepo.GetPlayRecords(userID, page, pageSize)
}

func (s *CommunityService) GetPlaylistsByUser(userID uint) ([]model.Playlist, error) {
	return s.communityRepo.GetPlaylistsByUser(userID)
}

func (s *CommunityService) GetFollowingArtists(userID uint, page, pageSize int) ([]model.User, int64, error) {
	return s.communityRepo.GetFollowingUsersByRole(userID, "artist", page, pageSize)
}

func (s *CommunityService) GetFollowingUsers(userID uint, page, pageSize int) ([]model.User, int64, error) {
	return s.communityRepo.GetFollowingUsersByRole(userID, "fan", page, pageSize)
}
