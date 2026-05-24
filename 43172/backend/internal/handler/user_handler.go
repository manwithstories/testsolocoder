package handler

import (
	"net/http"
	"strconv"
	"strings"

	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/service"
	resp "luxury-trading-platform/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			resp.Conflict(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Created(c, user)
}

func (h *UserHandler) RegisterAuthenticator(c *gin.Context) {
	var req service.AuthenticatorRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	user, err := h.userService.RegisterAuthenticator(c.Request.Context(), &req)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Created(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	result, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		resp.Unauthorized(c, err.Error())
		return
	}

	resp.Success(c, result)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		resp.BadRequest(c, err)
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID.(uint), updates)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")
	status := c.Query("status")

	users, total, err := h.userService.ListUsers(page, pageSize, model.UserRole(role), model.UserStatus(status))
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.SuccessWithPage(c, users, total, page, pageSize)
}

func (h *UserHandler) ApproveAuthenticator(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	if err := h.userService.ApproveAuthenticator(uint(id)); err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, gin.H{"message": "authenticator approved successfully"})
}

func (h *UserHandler) RejectAuthenticator(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	if err := h.userService.RejectAuthenticator(uint(id), req.Reason); err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, gin.H{"message": "authenticator rejected successfully"})
}

func (h *UserHandler) ListAuthenticators(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	authenticators, total, err := h.userService.ListAuthenticators(page, pageSize, model.AuthenticatorStatus(status))
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.SuccessWithPage(c, authenticators, total, page, pageSize)
}

func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "luxury-trading-platform",
	})
}
