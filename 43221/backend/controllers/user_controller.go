package controllers

import (
	"net/http"
	"strconv"

	"consultation-platform/config"
	"consultation-platform/models"
	"consultation-platform/services"
	"consultation-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userService *services.UserService
	cfg         *config.Config
}

func NewUserController(cfg *config.Config) *UserController {
	return &UserController{
		userService: services.NewUserService(cfg),
		cfg:         cfg,
	}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	user, err := ctrl.userService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

func (ctrl *UserController) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	result, err := ctrl.userService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, 404, "User not found")
		return
	}

	utils.SuccessResponse(c, user)
}

func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	user, err := ctrl.userService.UpdateUser(userID, updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")

	users, total, err := ctrl.userService.GetUsers(page, pageSize, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    users,
	})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid user ID")
		return
	}

	user, err := ctrl.userService.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, 404, "User not found")
		return
	}

	utils.SuccessResponse(c, user)
}

func (ctrl *UserController) VerifyProfessional(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid user ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=approved rejected"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	err = ctrl.userService.VerifyProfessional(id, models.VerificationStatus(req.Status), req.Note)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *UserController) GetPendingVerifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := ctrl.userService.GetPendingVerifications(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    users,
	})
}

func (ctrl *UserController) RefreshToken(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	token, err := ctrl.userService.RefreshToken(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"access_token": token})
}

type UserControllerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	VerifyProfessional(c *gin.Context)
	GetPendingVerifications(c *gin.Context)
	RefreshToken(c *gin.Context)
}
