package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

func RegisterAuthRoutes(r *gin.RouterGroup) {
	handler := NewAuthHandler()

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.GET("/me", middleware.Auth(), handler.GetCurrentUser)
		auth.PUT("/password", middleware.Auth(), handler.ChangePassword)
		auth.POST("/logout", middleware.Auth(), handler.Logout)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"full_name": user.FullName,
		"role":     user.Role,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, token, err := h.authService.Login(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
			"avatar":    user.AvatarURL,
		},
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	claims := utils.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "请先登录")
		return
	}

	user, err := h.authService.GetUserByID(claims.UserID)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"phone":      user.Phone,
		"full_name":  user.FullName,
		"gender":     user.Gender,
		"birth_date": user.BirthDate,
		"avatar_url": user.AvatarURL,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	claims := utils.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "请先登录")
		return
	}

	var req services.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.authService.ChangePassword(claims.UserID, &req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "密码修改成功"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.Response{
		Code:    0,
		Message: "登出成功",
	})
}
