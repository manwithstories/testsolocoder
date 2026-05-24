package handler

import (
	"strconv"

	"music-platform/internal/model"
	"music-platform/internal/service"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/jwt"
	"music-platform/pkg/response"
	"music-platform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "注册失败")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ip := c.ClientIP()
	user, token, err := h.userService.Login(&req, ip)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 401, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "登录失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := jwt.GetUserID(c)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)
	keyword := c.Query("keyword")
	role := c.Query("role")

	users, total, err := h.userService.ListUsers(page, pageSize, keyword, role)
	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}

	response.Page(c, users, total, page, pageSize)
}

func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.userService.UpdateUserRole(uint(id), req.Role)
	if err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.userService.UpdateUserStatus(uint(id), model.UserStatus(req.Status))
	if err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) UpdateArtistInfo(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.UpdateArtistInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.userService.UpdateArtistInfo(userID, &req)
	if err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetArtistInfo(c *gin.Context) {
	userID := jwt.GetUserID(c)

	info, err := h.userService.GetArtistInfo(userID)
	if err != nil {
		response.NotFound(c, "音乐人信息不存在")
		return
	}

	response.Success(c, info)
}

func (h *UserHandler) GetBalance(c *gin.Context) {
	userID := jwt.GetUserID(c)

	balance, frozenBalance, err := h.userService.GetBalance(userID)
	if err != nil {
		response.InternalError(c, "获取余额失败")
		return
	}

	response.Success(c, gin.H{
		"balance":        balance,
		"frozen_balance": frozenBalance,
	})
}

func (h *UserHandler) VerifyArtist(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.userService.VerifyArtist(uint(id))
	if err != nil {
		response.InternalError(c, "认证失败")
		return
	}

	response.Success(c, nil)
}
