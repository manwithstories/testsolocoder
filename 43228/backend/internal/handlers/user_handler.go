package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"tea-platform/internal/service"
	"tea-platform/internal/utils"
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
	var req service.RegisterRequest
	if !utils.Validate(c, &req) {
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if !utils.Validate(c, &req) {
		return
	}

	token, user, err := h.userService.Login(&req)
	if err != nil {
		utils.Fail(c, utils.CodeUnauthorized, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"avatar":   user.Avatar,
		},
	})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	user, err := h.userService.GetUserInfo(userID)
	if err != nil {
		utils.Fail(c, utils.CodeNotFound, err.Error())
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	var req service.UpdateProfileRequest
	if !utils.Validate(c, &req) {
		return
	}

	if err := h.userService.UpdateProfile(userID, &req); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	var req service.ChangePasswordRequest
	if !utils.Validate(c, &req) {
		return
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

type UserListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Keyword  string `form:"keyword"`
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	var query UserListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Fail(c, utils.CodeBadRequest, "参数错误")
		return
	}

	users, total, err := h.userService.GetUserList(query.Page, query.PageSize, query.Keyword)
	if err != nil {
		utils.Fail(c, utils.CodeInternalError, "查询用户列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":  users,
		"total": total,
	})
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive banned"`
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "用户ID格式错误")
		return
	}

	var req UpdateUserStatusRequest
	if !utils.Validate(c, &req) {
		return
	}

	if err := h.userService.UpdateUserStatus(uint(id), req.Status); err != nil {
		utils.Fail(c, utils.CodeInternalError, "更新用户状态失败")
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "请上传文件")
		return
	}

	filePath, err := utils.UploadFile(c, file, "avatar")
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	if err := h.userService.UpdateProfile(userID, &service.UpdateProfileRequest{
		Avatar: filePath,
	}); err != nil {
		utils.Fail(c, utils.CodeInternalError, "更新头像失败")
		return
	}

	utils.Success(c, gin.H{
		"avatar": filePath,
	})
}

func (h *UserHandler) getCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, utils.CodeUnauthorized, utils.MsgUnauthorized)
		return 0
	}
	uid, ok := userID.(uint)
	if !ok {
		utils.Fail(c, utils.CodeUnauthorized, "用户ID类型错误")
		return 0
	}
	return uid
}