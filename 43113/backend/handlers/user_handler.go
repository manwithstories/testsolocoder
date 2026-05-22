package handlers

import (
	"net/http"
	"qa-platform/services"
	"qa-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	result, err := h.userService.Register(&req)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.UserAlreadyExists, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	result, err := h.userService.Login(&req)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.InvalidCredentials, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userId")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.UserNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userId")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(userID, updates)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.InternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.UserNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	users, total, err := h.userService.GetUserList(page, pageSize, keyword)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active disabled banned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.userService.UpdateUserStatus(uint(id), req.Status); err != nil {
		utils.ErrorResponseWithMessage(c, utils.InternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *UserHandler) GetPointLogs(c *gin.Context) {
	userID := c.GetUint("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	logs, total, err := h.userService.GetPointLogs(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *UserHandler) ApplyExpert(c *gin.Context) {
	userID := c.GetUint("userId")

	var req struct {
		Field       string `json:"field" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.userService.ApplyExpert(userID, req.Field, req.Description); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *UserHandler) ReviewExpertApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	reviewerID := c.GetUint("userId")

	var req struct {
		Status string `json:"status" binding:"required,oneof=approved rejected"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.userService.ReviewExpertApplication(uint(id), reviewerID, req.Status, req.Remark); err != nil {
		utils.ErrorResponseWithMessage(c, utils.InternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *UserHandler) GetExpertApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	apps, total, err := h.userService.GetExpertApplications(page, pageSize, status)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     apps,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("userId")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.UserNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}
