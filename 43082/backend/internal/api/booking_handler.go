package api

import (
	"gym-management/internal/middleware"
	"gym-management/internal/pkg/utils"
	"gym-management/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler() *BookingHandler {
	return &BookingHandler{
		bookingService: service.NewBookingService(),
	}
}

func (h *BookingHandler) RegisterRoutes(r *gin.RouterGroup) {
	booking := r.Group("/bookings")
	booking.Use(middleware.JWTAuth())
	{
		booking.POST("/", h.Book)
		booking.GET("/", h.List)
		booking.GET("/:id", h.GetByID)
		booking.DELETE("/:id", h.Cancel)
		booking.GET("/member/:memberId", h.ListByMember)
		booking.GET("/schedule/:scheduleId", h.ListBySchedule)

		waitlist := booking.Group("/waitlist")
		{
			waitlist.POST("/", h.AddToWaitlist)
			waitlist.DELETE("/:id", h.RemoveFromWaitlist)
		}
	}
}

type BookRequest struct {
	ScheduleID uint `json:"schedule_id" binding:"required"`
	MemberID   uint `json:"member_id" binding:"required"`
}

func (h *BookingHandler) Book(c *gin.Context) {
	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	booking, err := h.bookingService.Book(req.MemberID, req.ScheduleID)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPageInfo(c)
	memberID, _ := strconv.ParseUint(c.Query("member_id"), 10, 32)
	scheduleID, _ := strconv.ParseUint(c.Query("schedule_id"), 10, 32)

	if memberID > 0 {
		bookings, total, err := h.bookingService.ListByMember(uint(memberID), page, pageSize)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}
		utils.SuccessWithPagination(c, bookings, page, pageSize, total)
	} else if scheduleID > 0 {
		bookings, total, err := h.bookingService.ListBySchedule(uint(scheduleID), page, pageSize)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}
		utils.SuccessWithPagination(c, bookings, page, pageSize, total)
	} else {
		utils.BadRequest(c, "请提供 member_id 或 schedule_id", nil)
	}
}

func (h *BookingHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	booking, err := h.bookingService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "预约不存在")
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	userID, _ := c.Get("userID")
	memberID, ok := userID.(uint)
	if !ok {
		utils.Unauthorized(c, "用户信息错误")
		return
	}

	if err := h.bookingService.Cancel(uint(id), memberID); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *BookingHandler) ListByMember(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的会员ID", nil)
		return
	}

	page, pageSize := utils.GetPageInfo(c)
	bookings, total, err := h.bookingService.ListByMember(uint(memberID), page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, bookings, page, pageSize, total)
}

func (h *BookingHandler) ListBySchedule(c *gin.Context) {
	scheduleID, err := strconv.ParseUint(c.Param("scheduleId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的排期ID", nil)
		return
	}

	page, pageSize := utils.GetPageInfo(c)
	bookings, total, err := h.bookingService.ListBySchedule(uint(scheduleID), page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, bookings, page, pageSize, total)
}

type WaitlistRequest struct {
	ScheduleID uint `json:"schedule_id" binding:"required"`
	MemberID   uint `json:"member_id" binding:"required"`
}

func (h *BookingHandler) AddToWaitlist(c *gin.Context) {
	var req WaitlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	waitlist, err := h.bookingService.AddToWaitlist(req.MemberID, req.ScheduleID)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, waitlist)
}

func (h *BookingHandler) RemoveFromWaitlist(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	userID, _ := c.Get("userID")
	memberID, ok := userID.(uint)
	if !ok {
		utils.Unauthorized(c, "用户信息错误")
		return
	}

	if err := h.bookingService.RemoveFromWaitlist(uint(id), memberID); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}
