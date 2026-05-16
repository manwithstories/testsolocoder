package controllers

import (
	"finance-api/config"
	"finance-api/middleware"
	"finance-api/models"
	"finance-api/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	cfg *config.Config
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{cfg: cfg}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	var existingUser models.User
	if result := utils.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		utils.BadRequest(c, "Username already exists")
		return
	}

	if req.Email != "" {
		if result := utils.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
			utils.BadRequest(c, "Email already exists")
			return
		}
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := user.HashPassword(); err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	if result := utils.DB.Create(&user); result.Error != nil {
		utils.InternalError(c, "Failed to create user: "+result.Error.Error())
		return
	}

	token, err := utils.GenerateToken(ctrl.cfg, user.ID, user.Username)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	var user models.User
	if result := utils.DB.Where("username = ?", req.Username).First(&user); result.Error != nil {
		utils.BadRequest(c, "Invalid username or password")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.BadRequest(c, "Invalid username or password")
		return
	}

	token, err := utils.GenerateToken(ctrl.cfg, user.ID, user.Username)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.User
	if result := utils.DB.Select("id", "username", "email", "created_at").First(&user, userID); result.Error != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user)
}
