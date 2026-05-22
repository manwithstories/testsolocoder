package handler

import (
	"car-rental/internal/config"
	"car-rental/internal/middleware"
	"car-rental/internal/model"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	cfg         *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
		cfg:         cfg,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if !utils.ValidateEmail(req.Email) {
		utils.BadRequest(c, "邮箱格式不正确")
		return
	}

	if !utils.ValidatePhone(req.Phone) {
		utils.BadRequest(c, "手机号格式不正确")
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	resp, err := h.userService.Login(&req, h.cfg.JWT.AccessExpire, h.cfg.JWT.RefreshExpire)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	resp, err := h.userService.RefreshToken(req.RefreshToken, h.cfg.JWT.AccessExpire, h.cfg.JWT.RefreshExpire)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	userData, err := h.userService.GetUserByID(user.UserID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, userData)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.userService.UpdateUser(user.UserID, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.userService.ChangePassword(user.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	keyword := c.Query("keyword")

	users, total, err := h.userService.GetAllUsers(page, pageSize, keyword)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, users, total, page, pageSize)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) UpdateUserAuthStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Status model.UserStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.userService.UpdateAuthStatus(uint(id), req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Status model.UserStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.userService.UpdateStatus(uint(id), req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.userService.DeleteUser(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) UploadLicense(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	field := c.Query("field")
	if field != "license_image" && field != "id_card_front" && field != "id_card_back" {
		utils.BadRequest(c, "无效的字段类型")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}

	if file.Size > h.cfg.Upload.MaxSize {
		utils.BadRequest(c, "文件大小超出限制")
		return
	}

	fileName := utils.GenerateFileName(file.Filename)
	savePath := h.cfg.Upload.Path + "/licenses/" + fileName

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.InternalServerError(c, "文件保存失败")
		return
	}

	updates := map[string]interface{}{
		field: "/uploads/licenses/" + fileName,
	}

	err = h.userService.UpdateUser(user.UserID, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"url": "/uploads/licenses/" + fileName})
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}

	if file.Size > h.cfg.Upload.MaxSize {
		utils.BadRequest(c, "文件大小超出限制")
		return
	}

	fileName := utils.GenerateFileName(file.Filename)
	savePath := h.cfg.Upload.Path + "/avatars/" + fileName

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.InternalServerError(c, "文件保存失败")
		return
	}

	updates := map[string]interface{}{
		"avatar": "/uploads/avatars/" + fileName,
	}

	err = h.userService.UpdateUser(user.UserID, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"url": "/uploads/avatars/" + fileName})
}
