package handler

import (
	"net/http"
	"strconv"
	"venue-booking/internal/dto"
	"venue-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	logService  *service.OperationLogService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
		logService:  service.NewOperationLogService(),
	}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "User not found"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(user))
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	user, err := h.userService.UpdateProfile(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to update profile"))
		return
	}

	h.logService.Log(c, userID.(uint), "update_profile", "user", nil)
	c.JSON(http.StatusOK, dto.Success(user))
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	users, total, err := h.userService.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get users"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     users,
	}))
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid user ID"))
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err = h.userService.UpdateRole(uint(id), req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to update role"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "update_role", "user", map[string]interface{}{
		"target_user_id": id,
		"new_role":       req.Role,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}
