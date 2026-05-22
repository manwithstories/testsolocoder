package handler

import (
	"hotel-system/internal/dto"
	"hotel-system/internal/middleware"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("登录参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("注册参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		RealName:  user.RealName,
		Phone:     user.Phone,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, userResponse)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.ErrorWithStatus(c, http.StatusUnauthorized, "未授权")
		return
	}

	user, err := h.userService.GetUserInfo(userID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		RealName:  user.RealName,
		Phone:     user.Phone,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, userResponse)
}
