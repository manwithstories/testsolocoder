package handler

import (
	"strconv"

	"music-platform/internal/service"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/jwt"
	"music-platform/pkg/response"
	"music-platform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CommunityHandler struct {
	communityService *service.CommunityService
}

func NewCommunityHandler() *CommunityHandler {
	return &CommunityHandler{
		communityService: service.NewCommunityService(),
	}
}

func (h *CommunityHandler) Follow(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.communityService.Follow(userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "关注失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) Unfollow(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.communityService.Unfollow(userID, &req)
	if err != nil {
		response.InternalError(c, "取消关注失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) IsFollowing(c *gin.Context) {
	userID := jwt.GetUserID(c)

	followingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	isFollowing, err := h.communityService.IsFollowing(userID, uint(followingID))
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, gin.H{"is_following": isFollowing})
}

func (h *CommunityHandler) GetFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	page, pageSize := utils.GetPageAndPageSize(c)

	followers, total, err := h.communityService.GetFollowers(uint(userID), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取粉丝列表失败")
		return
	}

	response.Page(c, followers, total, page, pageSize)
}

func (h *CommunityHandler) GetFollowings(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	page, pageSize := utils.GetPageAndPageSize(c)

	followings, total, err := h.communityService.GetFollowings(uint(userID), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取关注列表失败")
		return
	}

	response.Page(c, followings, total, page, pageSize)
}

func (h *CommunityHandler) GetFollowerCount(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	count, err := h.communityService.GetFollowerCount(uint(userID))
	if err != nil {
		response.InternalError(c, "获取粉丝数失败")
		return
	}

	response.Success(c, gin.H{"count": count})
}

func (h *CommunityHandler) GetFollowingCount(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	count, err := h.communityService.GetFollowingCount(uint(userID))
	if err != nil {
		response.InternalError(c, "获取关注数失败")
		return
	}

	response.Success(c, gin.H{"count": count})
}

func (h *CommunityHandler) CreateComment(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	comment, err := h.communityService.CreateComment(userID, &req)
	if err != nil {
		response.InternalError(c, "评论失败")
		return
	}

	response.Success(c, comment)
}

func (h *CommunityHandler) DeleteComment(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.communityService.DeleteComment(uint(id), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "删除评论失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) GetComments(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)

	workID, _ := strconv.ParseUint(c.DefaultQuery("work_id", "0"), 10, 64)
	albumIDStr := c.Query("album_id")
	playlistIDStr := c.Query("playlist_id")

	var albumID *uint
	if albumIDStr != "" {
		id, _ := strconv.ParseUint(albumIDStr, 10, 64)
		uid := uint(id)
		albumID = &uid
	}

	var playlistID *uint
	if playlistIDStr != "" {
		id, _ := strconv.ParseUint(playlistIDStr, 10, 64)
		uid := uint(id)
		playlistID = &uid
	}

	comments, total, err := h.communityService.GetComments(uint(workID), albumID, playlistID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取评论列表失败")
		return
	}

	response.Page(c, comments, total, page, pageSize)
}

func (h *CommunityHandler) CreatePlaylist(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	playlist, err := h.communityService.CreatePlaylist(userID, &req)
	if err != nil {
		response.InternalError(c, "创建歌单失败")
		return
	}

	response.Success(c, playlist)
}

func (h *CommunityHandler) UpdatePlaylist(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req service.UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.communityService.UpdatePlaylist(uint(id), userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "更新歌单失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) DeletePlaylist(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.communityService.DeletePlaylist(uint(id), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "删除歌单失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) GetPlaylistByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	playlist, err := h.communityService.GetPlaylistByID(uint(id))
	if err != nil {
		response.NotFound(c, "歌单不存在")
		return
	}

	response.Success(c, playlist)
}

func (h *CommunityHandler) ListPlaylists(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.DefaultQuery("user_id", "0"), 10, 64)
	keyword := c.Query("keyword")
	page, pageSize := utils.GetPageAndPageSize(c)

	playlists, total, err := h.communityService.ListPlaylists(uint(userID), page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, "获取歌单列表失败")
		return
	}

	response.Page(c, playlists, total, page, pageSize)
}

func (h *CommunityHandler) AddWorkToPlaylist(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.PlaylistWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.communityService.AddWorkToPlaylist(&req, userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "添加失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) RemoveWorkFromPlaylist(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.PlaylistWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.communityService.RemoveWorkFromPlaylist(&req, userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "移除失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) GetNotifications(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	notifications, total, err := h.communityService.GetNotifications(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取通知列表失败")
		return
	}

	response.Page(c, notifications, total, page, pageSize)
}

func (h *CommunityHandler) MarkNotificationAsRead(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.communityService.MarkNotificationAsRead(uint(id), userID)
	if err != nil {
		response.InternalError(c, "标记失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	userID := jwt.GetUserID(c)

	err := h.communityService.MarkAllNotificationsAsRead(userID)
	if err != nil {
		response.InternalError(c, "标记失败")
		return
	}

	response.Success(c, nil)
}

func (h *CommunityHandler) GetUnreadNotificationCount(c *gin.Context) {
	userID := jwt.GetUserID(c)

	count, err := h.communityService.GetUnreadNotificationCount(userID)
	if err != nil {
		response.InternalError(c, "获取未读数失败")
		return
	}

	response.Success(c, gin.H{"count": count})
}

func (h *CommunityHandler) GetPlayRecords(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	records, total, err := h.communityService.GetPlayRecords(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取播放记录失败")
		return
	}

	response.Page(c, records, total, page, pageSize)
}

func (h *CommunityHandler) GetMyPlaylists(c *gin.Context) {
	userID := jwt.GetUserID(c)

	playlists, err := h.communityService.GetPlaylistsByUser(userID)
	if err != nil {
		response.InternalError(c, "获取歌单列表失败")
		return
	}

	response.Success(c, playlists)
}

func (h *CommunityHandler) GetMyFollowingArtists(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	artists, total, err := h.communityService.GetFollowingArtists(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取关注列表失败")
		return
	}

	response.Page(c, artists, total, page, pageSize)
}

func (h *CommunityHandler) GetMyFollowingUsers(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	users, total, err := h.communityService.GetFollowingUsers(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取关注列表失败")
		return
	}

	response.Page(c, users, total, page, pageSize)
}

func (h *CommunityHandler) GetMyFollowers(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	followers, total, err := h.communityService.GetFollowers(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取粉丝列表失败")
		return
	}

	response.Page(c, followers, total, page, pageSize)
}
