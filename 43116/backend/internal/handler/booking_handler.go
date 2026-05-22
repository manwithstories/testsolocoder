package handler

import (
	"car-rental/internal/config"
	"car-rental/internal/middleware"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(cfg *config.EmailConfig) *BookingHandler {
	messageService := service.NewMessageService(cfg)
	return &BookingHandler{
		bookingService: service.NewBookingService(messageService),
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	booking, err := h.bookingService.CreateBooking(user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	booking, err := h.bookingService.GetBookingByID(uint(id))
	if err != nil {
		utils.NotFound(c, "预订不存在")
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) GetBookingByNo(c *gin.Context) {
	bookingNo := c.Param("no")

	booking, err := h.bookingService.GetBookingByNo(bookingNo)
	if err != nil {
		utils.NotFound(c, "预订不存在")
		return
	}

	utils.Success(c, booking)
}

func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	status := c.Query("status")
	carID, _ := strconv.ParseUint(c.Query("car_id"), 10, 64)

	var startDate, endDate *time.Time
	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			startDate = &t
		}
	}
	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			endDate = &t
		}
	}

	bookings, total, err := h.bookingService.GetAllBookings(page, pageSize, uint(userID), status, uint(carID), startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, bookings, total, page, pageSize)
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	page, pageSize, _, _ := utils.ParsePageParams(c)

	bookings, total, err := h.bookingService.GetUserBookings(user.UserID, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, bookings, total, page, pageSize)
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.bookingService.ConfirmBooking(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	err := h.bookingService.CancelBooking(uint(id), req.Reason)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *BookingHandler) CompleteBooking(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.bookingService.CompleteBooking(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *BookingHandler) CheckCarAvailability(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Query("car_id"), 10, 64)
	pickupTime, _ := time.Parse(time.RFC3339, c.Query("pickup_time"))
	returnTime, _ := time.Parse(time.RFC3339, c.Query("return_time"))

	available := h.bookingService.CheckCarAvailability(uint(carID), pickupTime, returnTime)

	utils.Success(c, gin.H{"available": available})
}

func (h *BookingHandler) CalculatePrice(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Query("car_id"), 10, 64)
	pickupTime, _ := time.Parse(time.RFC3339, c.Query("pickup_time"))
	returnTime, _ := time.Parse(time.RFC3339, c.Query("return_time"))
	promoCode := c.Query("promo_code")

	price, err := h.bookingService.CalculatePrice(uint(carID), pickupTime, returnTime, promoCode)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, price)
}
