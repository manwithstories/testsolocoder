package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

// UserHandler 用户 HTTP 处理器
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// toUserResponse 模型转响应 DTO
func toUserResponse(u *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Role:      u.Role,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// currentUserID 从上下文中获取当前用户 ID
func currentUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := val.(uint)
	if !ok {
		return 0, false
	}
	return id, true
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toUserResponse(user))
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, token, err := h.service.Login(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, dto.LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	})
}

// Profile 获取当前登录用户资料
func (h *UserHandler) Profile(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, toUserResponse(user))
}

// UpdateProfile 更新当前登录用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.service.UpdateProfile(id, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toUserResponse(user))
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.service.ChangePassword(id, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListUsers 分页查询用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	users, total, err := h.service.List(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	list := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		list = append(list, toUserResponse(u))
	}

	response.Success(c, dto.UserListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	})
}

// GetUser 根据 ID 获取用户
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "用户 ID 格式错误")
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toUserResponse(user))
}
