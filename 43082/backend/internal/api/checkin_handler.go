package api

import (
	"gym-management/internal/middleware"
	"gym-management/internal/pkg/utils"
	"gym-management/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CheckInHandler struct {
	checkInService service.CheckInService
}

func NewCheckInHandler() *CheckInHandler {
	return &CheckInHandler{
		checkInService: service.NewCheckInService(),
	}
}

func (h *CheckInHandler) RegisterRoutes(r *gin.RouterGroup) {
	checkIn := r.Group("/check-ins")
	checkIn.Use(middleware.JWTAuth())
	{
		checkIn.POST("/", h.CheckIn)
		checkIn.GET("/", h.List)
		checkIn.GET("/:id", h.GetByID)
		checkIn.GET("/member/:memberId", h.ListByMember)
		checkIn.GET("/date/:date", h.ListByDate)
	}
}

type CheckInRequest struct {
	MemberID   uint  `json:"member_id" binding:"required"`
	ScheduleID *uint `json:"schedule_id"`
}

func (h *CheckInHandler) CheckIn(c *gin.Context) {
	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	checkIn, err := h.checkInService.ValidateAndCheckIn(req.MemberID, req.ScheduleID)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, checkIn)
}

func (h *CheckInHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPageInfo(c)
	memberID, _ := strconv.ParseUint(c.Query("member_id"), 10, 32)
	dateStr := c.Query("date")

	if memberID > 0 {
		checkIns, total, err := h.checkInService.ListByMember(uint(memberID), page, pageSize)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}
		utils.SuccessWithPagination(c, checkIns, page, pageSize, total)
	} else if dateStr != "" {
		date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			utils.BadRequest(c, "日期格式错误，请使用 YYYY-MM-DD", nil)
			return
		}
		checkIns, err := h.checkInService.ListByDate(date)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}
		utils.Success(c, checkIns)
	} else {
		utils.BadRequest(c, "请提供 member_id 或 date", nil)
	}
}

func (h *CheckInHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	checkIn, err := h.checkInService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "签到记录不存在")
		return
	}

	utils.Success(c, checkIn)
}

func (h *CheckInHandler) ListByMember(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会员ID", nil)
		return
	}

	page, pageSize := utils.GetPageInfo(c)
	checkIns, total, err := h.checkInService.ListByMember(uint(memberID), page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, checkIns, page, pageSize, total)
}

func (h *CheckInHandler) ListByDate(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		utils.BadRequest(c, "日期格式错误，请使用 YYYY-MM-DD", nil)
		return
	}

	checkIns, err := h.checkInService.ListByDate(date)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, checkIns)
}
