package controller

import (
	"strconv"
	"ticket-system/internal/common/response"
	"ticket-system/internal/dto"
	"ticket-system/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := c.userService.Register(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, user)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	resp, err := c.userService.Login(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}

func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		response.NotFound(ctx, "用户不存在")
		return
	}
	response.Success(ctx, user)
}

func (c *UserController) GetUserList(ctx *gin.Context) {
	var req dto.UserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	users, total, err := c.userService.GetUserList(&req)
	if err != nil {
		response.ServerError(ctx, "获取用户列表失败")
		return
	}

	response.Page(ctx, users, total, req.Page, req.PageSize)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的用户ID")
		return
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		response.NotFound(ctx, "用户不存在")
		return
	}

	response.Success(ctx, user)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的用户ID")
		return
	}

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := c.userService.UpdateUser(uint(id), &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, user)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的用户ID")
		return
	}

	if err := c.userService.DeleteUser(uint(id)); err != nil {
		response.ServerError(ctx, "删除用户失败")
		return
	}

	response.Success(ctx, nil)
}
