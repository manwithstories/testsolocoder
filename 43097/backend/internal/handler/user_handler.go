package handler

import (
	"hotel-system/internal/dto"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的用户ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := h.userService.GetUserInfo(uint(id))
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

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("创建用户参数校验失败: %v", err)
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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的用户ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新用户参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
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

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的用户ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("获取用户列表参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	users, total, err := h.userService.ListUsers(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			RealName:  user.RealName,
			Phone:     user.Phone,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.PageResult(c, userResponses, total, req.GetPage(), req.GetPageSize())
}
