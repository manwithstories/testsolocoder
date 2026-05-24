package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	user, err := ctrl.authService.Register(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, user)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	response, err := ctrl.authService.Login(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, response)
}

func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := ctrl.authService.GetUserByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user)
}

func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.authService.UpdateProfile(userID.(uint), data); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.authService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "Logged out successfully",
	})
}
