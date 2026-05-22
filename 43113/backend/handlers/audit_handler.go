package handlers

import (
	"qa-platform/services"
	"qa-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *services.AuditService
}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{
		auditService: services.NewAuditService(),
	}
}

func (h *AuditHandler) GetAuditList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	targetType := c.Query("targetType")
	status := c.Query("status")

	query := services.AuditQuery{
		Page:       page,
		PageSize:   pageSize,
		TargetType: targetType,
		Status:     status,
	}

	records, total, err := h.auditService.GetAuditList(query)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     records,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AuditHandler) AuditContent(c *gin.Context) {
	adminID := c.GetUint("userId")

	var req services.AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.auditService.AuditContent(adminID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AuditHandler) GetPendingAuditCount(c *gin.Context) {
	counts, err := h.auditService.GetPendingAuditCount()
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, counts)
}

func (h *AuditHandler) GetReportList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	query := services.ReportQuery{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
	}

	reports, total, err := h.auditService.GetReportList(query)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     reports,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AuditHandler) HandleReport(c *gin.Context) {
	handlerID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	var req services.HandleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.auditService.HandleReport(uint(id), handlerID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AuditHandler) CreateReport(c *gin.Context) {
	reporterID := c.GetUint("userId")

	var req struct {
		TargetType  string `json:"targetType" binding:"required,oneof=question answer comment"`
		TargetID    uint   `json:"targetId" binding:"required"`
		Reason      string `json:"reason" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.auditService.CreateReport(reporterID, req.TargetType, req.TargetID, req.Reason, req.Description); err != nil {
		utils.ErrorResponseWithMessage(c, utils.InternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AuditHandler) GetSensitiveWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	query := services.SensitiveWordQuery{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Keyword:  keyword,
	}

	words, total, err := h.auditService.GetSensitiveWords(query)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     words,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AuditHandler) CreateSensitiveWord(c *gin.Context) {
	var req struct {
		Word      string `json:"word" binding:"required"`
		Category  string `json:"category"`
		ReplaceTo string `json:"replaceTo"`
		Level     int    `json:"level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.auditService.CreateSensitiveWord(req.Word, req.Category, req.ReplaceTo, req.Level); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AuditHandler) DeleteSensitiveWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.auditService.DeleteSensitiveWord(uint(id)); err != nil {
		utils.ErrorResponseWithMessage(c, utils.InternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AuditHandler) CheckContent(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	filteredContent, sensitiveWords := h.auditService.BatchCheckContent(req.Content)

	utils.SuccessResponse(c, gin.H{
		"filteredContent": filteredContent,
		"sensitiveWords":  sensitiveWords,
		"hasSensitive":    len(sensitiveWords) > 0,
	})
}

type FavoriteHandler struct {
	favoriteService *services.FavoriteService
}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: services.NewFavoriteService(),
	}
}

func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := c.GetUint("userId")

	var req services.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.favoriteService.AddFavorite(userID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.AlreadyFavorited, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("userId")

	var req services.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.favoriteService.RemoveFavorite(userID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	userID := c.GetUint("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	targetType := c.Query("targetType")

	favorites, total, err := h.favoriteService.GetUserFavorites(userID, targetType, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     favorites,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

type FollowHandler struct {
	followService *services.FollowService
}

func NewFollowHandler() *FollowHandler {
	return &FollowHandler{
		followService: services.NewFollowService(),
	}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	userID := c.GetUint("userId")

	var req services.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.followService.Follow(userID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.AlreadyFollowed, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID := c.GetUint("userId")

	var req services.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.followService.Unfollow(userID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *FollowHandler) GetUserFollows(c *gin.Context) {
	userID := c.GetUint("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	followingType := c.Query("followingType")

	follows, total, err := h.followService.GetUserFollows(userID, followingType, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     follows,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *FollowHandler) GetUserFollowers(c *gin.Context) {
	userID := c.GetUint("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	followers, total, err := h.followService.GetUserFollowers(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     followers,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *FollowHandler) IsFollowing(c *gin.Context) {
	followerID := c.GetUint("userId")
	followingType := c.Query("followingType")
	followingID, _ := strconv.ParseUint(c.Query("followingId"), 10, 64)

	isFollowing := h.followService.IsFollowing(followerID, followingType, uint(followingID))
	utils.SuccessResponse(c, gin.H{"isFollowing": isFollowing})
}

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notificationService: services.NewNotificationService(),
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetUint("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	isReadStr := c.Query("isRead")

	var isRead *bool
	if isReadStr != "" {
		read := isReadStr == "true"
		isRead = &read
	}

	notifications, total, err := h.notificationService.GetUserNotifications(userID, isRead, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     notifications,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.notificationService.MarkAsRead(userID, uint(id)); err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("userId")

	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userId")

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{"unreadCount": count})
}
