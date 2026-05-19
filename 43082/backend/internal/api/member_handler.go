package api

import (
	"gym-management/internal/middleware"
	"gym-management/internal/models"
	"gym-management/internal/pkg/utils"
	"gym-management/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService service.MemberService
}

func NewMemberHandler() *MemberHandler {
	return &MemberHandler{
		memberService: service.NewMemberService(),
	}
}

func (h *MemberHandler) RegisterRoutes(r *gin.RouterGroup) {
	member := r.Group("/members")
	{
		member.POST("/register", h.Register)
		member.POST("/login", h.Login)
		member.GET("/", middleware.JWTAuth(), h.List)
		member.GET("/:id", middleware.JWTAuth(), h.GetByID)
		member.PUT("/:id", middleware.JWTAuth(), h.Update)
		member.DELETE("/:id", middleware.JWTAuth(), middleware.AdminRequired(), h.Delete)
		member.PATCH("/:id/status", middleware.JWTAuth(), middleware.AdminRequired(), h.UpdateStatus)
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=6"`
	Gender   string `json:"gender"`
}

func (h *MemberHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	member := &models.Member{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Gender:   req.Gender,
	}

	if err := h.memberService.Register(member); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, member)
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *MemberHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	member, err := h.memberService.Login(req.Phone, req.Password)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	token, err := middleware.GenerateToken(member.ID, member.Name, "member")
	if err != nil {
		utils.InternalServerError(c, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token":  token,
		"member": member,
	})
}

func (h *MemberHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPageInfo(c)
	keyword := c.Query("keyword")

	members, total, err := h.memberService.List(page, pageSize, keyword)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, members, page, pageSize, total)
}

func (h *MemberHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	member, err := h.memberService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "会员不存在")
		return
	}

	utils.Success(c, member)
}

func (h *MemberHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	delete(updates, "password")
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	if err := h.memberService.Update(uint(id), updates); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *MemberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	if err := h.memberService.Delete(uint(id)); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *MemberHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := h.memberService.UpdateStatus(uint(id), req.Status); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}
