package handler

import (
	"net/http"
	"venue-booking/internal/dto"
	"venue-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *service.UserService
	logService  *service.OperationLogService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: service.NewUserService(),
		logService:  service.NewOperationLogService(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	user, token, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, user.ID, "register", "auth", map[string]interface{}{
		"username": req.Username,
		"email": req.Email,
	})

	c.JSON(http.StatusOK, dto.Success(dto.LoginResponse{
		Token: token,
		User:  user,
	}))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	user, token, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Error(401, err.Error()))
		return
	}

	h.logService.Log(c, user.ID, "login", "auth", map[string]interface{}{
		"username": req.Username,
	})

	c.JSON(http.StatusOK, dto.Success(dto.LoginResponse{
		Token: token,
		User:  user,
	}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "logout", "auth", nil)
	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *AuthHandler) SendVerifyEmail(c *gin.Context) {
	var req dto.SendVerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err := h.userService.SendVerificationCode(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to send verification email"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err := h.userService.VerifyEmail(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid verification code"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err := h.userService.SendPasswordReset(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to send reset email"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err := h.userService.ResetPassword(req.Email, req.Code, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Failed to reset password"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessNoData())
}
