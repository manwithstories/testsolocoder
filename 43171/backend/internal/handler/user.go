package handler

import (
	"drone-rental/internal/dto"
	"drone-rental/internal/middleware"
	"drone-rental/internal/model"
	"drone-rental/internal/pkg/response"
	"drone-rental/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	user, err := h.userService.Register(&req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	resp, err := h.userService.Login(&req)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	response.Success(c, resp)
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.ErrNotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.userService.Update(userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) VerifyPilot(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.VerifyPilotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	licenseImage := c.PostForm("license_image")
	if err := h.userService.VerifyPilot(userID, &req, licenseImage); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) VerifyOwner(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.VerifyOwnerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.userService.VerifyOwner(userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	users, total, err := h.userService.List(page, pageSize, "", keyword)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, users, total, page, pageSize)
}

func (h *UserHandler) AuditVerify(c *gin.Context) {
	var req struct {
		UserID uint                `json:"user_id" binding:"required"`
		Status model.VerifyStatus `json:"status" binding:"required,oneof=approved rejected"`
		Remark string            `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.userService.AuditVerify(req.UserID, req.Status, req.Remark); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}
