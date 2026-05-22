package handler

import (
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req dto.BookingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("创建预订参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	booking, err := h.bookingService.CreateBooking(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	bookingResponse := convertToBookingResponse(booking)
	utils.Success(c, bookingResponse)
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的预订ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的预订ID")
		return
	}

	booking, err := h.bookingService.GetBookingDetail(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	bookingResponse := convertToBookingResponse(booking)
	utils.Success(c, bookingResponse)
}

func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的预订ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的预订ID")
		return
	}

	var req dto.BookingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新预订参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	booking, err := h.bookingService.UpdateBooking(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	bookingResponse := convertToBookingResponse(booking)
	utils.Success(c, bookingResponse)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的预订ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的预订ID")
		return
	}

	var req dto.BookingCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("取消预订参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	booking, err := h.bookingService.CancelBooking(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	bookingResponse := convertToBookingResponse(booking)
	utils.Success(c, bookingResponse)
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的预订ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的预订ID")
		return
	}

	var req dto.BookingConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("确认预订参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	booking, err := h.bookingService.ConfirmBooking(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	bookingResponse := convertToBookingResponse(booking)
	utils.Success(c, bookingResponse)
}

func (h *BookingHandler) ListBookings(c *gin.Context) {
	var req dto.BookingListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("获取预订列表参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	bookings, total, err := h.bookingService.ListBookings(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var bookingResponses []dto.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, convertToBookingResponse(&booking))
	}

	utils.PageResult(c, bookingResponses, total, req.GetPage(), req.GetPageSize())
}

func (h *BookingHandler) CalculatePrice(c *gin.Context) {
	var req dto.BookingPriceCalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("计算价格参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := h.bookingService.CalculatePrice(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, result)
}

func convertToBookingResponse(booking *model.Booking) dto.BookingResponse {
	resp := dto.BookingResponse{
		ID:          booking.ID,
		BookingNo:   booking.BookingNo,
		RoomID:      booking.RoomID,
		MemberID:    booking.MemberID,
		GuestName:   booking.GuestName,
		GuestPhone:  booking.GuestPhone,
		GuestIDCard: booking.GuestIDCard,
		CheckInDate:  booking.CheckInDate.Format("2006-01-02"),
		CheckOutDate: booking.CheckOutDate.Format("2006-01-02"),
		Days:        booking.Days,
		TotalPrice:  booking.TotalPrice,
		Status:      booking.Status,
		PaidAmount:  booking.PaidAmount,
		Remarks:     booking.Remarks,
		CreatedAt:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   booking.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if booking.CancelDeadline != nil {
		cancelDeadlineStr := booking.CancelDeadline.Format("2006-01-02 15:04:05")
		resp.CancelDeadline = &cancelDeadlineStr
	}

	if booking.Room != nil {
		roomTypeResp := &dto.RoomTypeResponse{
			ID:          booking.Room.RoomType.ID,
			Name:        booking.Room.RoomType.Name,
			Description: booking.Room.RoomType.Description,
			BasePrice:   booking.Room.RoomType.BasePrice,
			BedCount:    booking.Room.RoomType.BedCount,
			MaxGuests:   booking.Room.RoomType.MaxGuests,
			Facilities:  []string(booking.Room.RoomType.Facilities),
			CreatedAt:   booking.Room.RoomType.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   booking.Room.RoomType.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		resp.Room = &dto.RoomResponse{
			ID:         booking.Room.ID,
			RoomNo:     booking.Room.RoomNo,
			Floor:      booking.Room.Floor,
			RoomTypeID: booking.Room.RoomTypeID,
			RoomType:   roomTypeResp,
			Status:     booking.Room.Status,
			Price:      booking.Room.Price,
			Facilities: []string(booking.Room.Facilities),
			CreatedAt:  booking.Room.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  booking.Room.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	if booking.Member != nil {
		levelResp := &dto.MemberLevelResponse{
			ID:           booking.Member.Level.ID,
			Name:         booking.Member.Level.Name,
			DiscountRate: booking.Member.Level.DiscountRate,
		}
		resp.Member = &dto.MemberResponse{
			ID:       booking.Member.ID,
			MemberNo: booking.Member.MemberNo,
			Name:     booking.Member.Name,
			Phone:    booking.Member.Phone,
			LevelID:  booking.Member.LevelID,
			Level:    levelResp,
		}
	}

	return resp
}
