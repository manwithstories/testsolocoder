package handlers

import (
	"strconv"

	"meeting-room/internal/middleware"
	"meeting-room/internal/services"
	"meeting-room/internal/utils"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler() *BookingHandler {
	return &BookingHandler{
		bookingService: services.NewBookingService(),
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req services.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: " + err.Error())
		return
	}

	booking, err := h.bookingService.CreateBooking(&req, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "预订ID错误")
		return
	}

	booking, err := h.bookingService.GetBooking(uint(id))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) ListBookings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	roomID, _ := strconv.ParseUint(c.Query("room_id"), 10, 32)
	status, _ := strconv.Atoi(c.Query("status"))

	userID := uint(0)
	role := middleware.GetUserRole(c)
	if role == "user" {
		userID = middleware.GetUserID(c)
	}

	bookings, total, err := h.bookingService.ListBookings(page, pageSize, userID, uint(roomID), status)
	if err != nil {
		utils.InternalError(c, "获取预订列表失败")
		return
	}

	utils.Success(c, gin.H{
		"bookings": bookings,
		"total":    total,
	})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "预订ID错误")
		return
	}

	userID := middleware.GetUserID(c)

	var req services.CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	err = h.bookingService.CancelBooking(uint(id), userID, req.Reason)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *BookingHandler) RescheduleBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "预订ID错误")
		return
	}

	userID := middleware.GetUserID(c)

	var req services.RescheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	booking, err := h.bookingService.RescheduleBooking(uint(id), userID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) ApproveBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "预订ID错误")
		return
	}

	err = h.bookingService.ApproveBooking(uint(id))
	if err != nil {
		utils.InternalError(c, "审核预订失败")
		return
	}

	utils.Success(c, nil)
}

func (h *BookingHandler) CompleteBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "预订ID错误")
		return
	}

	err = h.bookingService.CompleteBooking(uint(id))
	if err != nil {
		utils.InternalError(c, "完成预订失败")
		return
	}

	utils.Success(c, nil)
}
