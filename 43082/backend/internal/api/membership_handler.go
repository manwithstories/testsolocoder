package api

import (
	"gym-management/internal/middleware"
	"gym-management/internal/models"
	"gym-management/internal/pkg/utils"
	"gym-management/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipService service.MembershipService
}

func NewMembershipHandler() *MembershipHandler {
	return &MembershipHandler{
		membershipService: service.NewMembershipService(),
	}
}

func (h *MembershipHandler) RegisterRoutes(r *gin.RouterGroup) {
	membership := r.Group("/memberships")
	membership.Use(middleware.JWTAuth())
	{
		membership.POST("/", h.Create)
		membership.GET("/", h.List)
		membership.GET("/:id", h.GetByID)
		membership.GET("/member/:memberId", h.GetByMemberID)
		membership.POST("/:id/renew", h.Renew)
		membership.POST("/:memberId/upgrade", h.Upgrade)
		membership.PATCH("/:id/status", h.UpdateStatus)
		membership.GET("/:memberId/validity", h.CheckValidity)
	}
}

type CreateMembershipRequest struct {
	MemberID   uint                    `json:"member_id" binding:"required"`
	Type       models.MembershipType   `json:"type" binding:"required,oneof=monthly quarter yearly"`
	StartDate  string                  `json:"start_date"`
	Price      float64                 `json:"price"`
	AutoRenew  bool                    `json:"auto_renew"`
}

func (h *MembershipHandler) Create(c *gin.Context) {
	var req CreateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	membership := &models.Membership{
		MemberID:  req.MemberID,
		Type:      req.Type,
		Price:     req.Price,
		AutoRenew: req.AutoRenew,
	}

	if req.StartDate != "" {
		startDate, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
		if err != nil {
			utils.BadRequest(c, "日期格式错误，请使用 YYYY-MM-DD", nil)
			return
		}
		membership.StartDate = startDate
	}

	if err := h.membershipService.Create(membership); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, membership)
}

func (h *MembershipHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPageInfo(c)
	memberID, _ := strconv.ParseUint(c.Query("member_id"), 10, 32)

	memberships, total, err := h.membershipService.List(page, pageSize, uint(memberID))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, memberships, page, pageSize, total)
}

func (h *MembershipHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	membership, err := h.membershipService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "会员卡不存在")
		return
	}

	utils.Success(c, membership)
}

func (h *MembershipHandler) GetByMemberID(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会员ID", nil)
		return
	}

	membership, err := h.membershipService.GetByMemberID(uint(memberID))
	if err != nil {
		utils.NotFound(c, "会员没有有效会员卡")
		return
	}

	utils.Success(c, membership)
}

type RenewRequest struct {
	NewType models.MembershipType `json:"new_type" binding:"required,oneof=monthly quarter yearly"`
}

func (h *MembershipHandler) Renew(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var req RenewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	newMembership, err := h.membershipService.Renew(uint(id), req.NewType)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, newMembership)
}

type UpgradeRequest struct {
	NewType models.MembershipType `json:"new_type" binding:"required,oneof=monthly quarter yearly"`
	Price   float64               `json:"price" binding:"required,min=0"`
}

func (h *MembershipHandler) Upgrade(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会员ID", nil)
		return
	}

	var req UpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	newMembership, err := h.membershipService.Upgrade(uint(memberID), req.NewType, req.Price)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, newMembership)
}

func (h *MembershipHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2 3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := h.membershipService.UpdateStatus(uint(id), req.Status); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *MembershipHandler) CheckValidity(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会员ID", nil)
		return
	}

	valid, err := h.membershipService.CheckValidity(uint(memberID))
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, gin.H{"valid": valid})
}
