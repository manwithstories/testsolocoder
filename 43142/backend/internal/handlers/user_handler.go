package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/middleware"
	"recruitment-platform/internal/models"
	"recruitment-platform/internal/services"
	"recruitment-platform/internal/utils"
	"recruitment-platform/pkg/logger"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	limiter := middleware.GetRateLimiter("register")
	if !limiter.Allow(c.ClientIP(), 5, time.Hour) {
		utils.BadRequest(c, "注册过于频繁，请稍后再试")
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		logger.Error("注册失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		logger.Error("生成Token失败: %v", err)
		utils.InternalError(c, "注册成功，但登录失败")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	limiter := middleware.GetRateLimiter("login")
	if !limiter.Allow(req.Email, 10, 30*time.Minute) {
		utils.BadRequest(c, "登录尝试过于频繁，请稍后再试")
		return
	}

	user, err := h.userService.Login(&req)
	if err != nil {
		logger.Error("登录失败: %v", err)
		utils.Unauthorized(c, err.Error())
		return
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		logger.Error("生成Token失败: %v", err)
		utils.InternalError(c, "生成登录凭证失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var profile models.ApplicantProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	profile.UserID = userID

	if err := h.userService.UpdateProfile(userID, &profile); err != nil {
		logger.Error("更新档案失败: %v", err)
		utils.InternalError(c, "更新档案失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", profile)
}

func (h *UserHandler) GetCompany(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if user.Company == nil {
		utils.NotFound(c, "公司信息不存在")
		return
	}

	utils.Success(c, user.Company)
}

func (h *UserHandler) UpdateCompany(c *gin.Context) {
	userID := c.GetUint("user_id")

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userService.UpdateCompany(userID, &company); err != nil {
		logger.Error("更新公司信息失败: %v", err)
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", company)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	role := c.Query("role")
	status := c.Query("status")

	users, total, err := h.userService.ListUsers(page, pageSize, role, status)
	if err != nil {
		logger.Error("获取用户列表失败: %v", err)
		utils.InternalError(c, "获取用户列表失败")
		return
	}

	utils.Paginated(c, users, page, pageSize, total)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.userService.UpdateUserStatus(id, models.UserStatus(req.Status)); err != nil {
		utils.InternalError(c, "更新用户状态失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}
