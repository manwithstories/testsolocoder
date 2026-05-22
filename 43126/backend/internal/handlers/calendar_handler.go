package handlers

import (
	"strconv"
	"time"

	"meeting-room/internal/services"
	"meeting-room/internal/utils"

	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	calendarService *services.CalendarService
}

func NewCalendarHandler() *CalendarHandler {
	return &CalendarHandler{
		calendarService: services.NewCalendarService(),
	}
}

func (h *CalendarHandler) GetWeekCalendar(c *gin.Context) {
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	roomID, _ := strconv.ParseUint(c.Query("room_id"), 10, 32)
	floor := c.Query("floor")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.BadRequest(c, "日期格式错误")
		return
	}

	bookings, err := h.calendarService.GetWeekCalendar(date, uint(roomID), floor)
	if err != nil {
		utils.InternalError(c, "获取日历数据失败")
		return
	}

	utils.Success(c, bookings)
}

func (h *CalendarHandler) GetMonthCalendar(c *gin.Context) {
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01"))
	roomID, _ := strconv.ParseUint(c.Query("room_id"), 10, 32)
	floor := c.Query("floor")

	date, err := time.Parse("2006-01", dateStr)
	if err != nil {
		utils.BadRequest(c, "日期格式错误")
		return
	}

	bookings, err := h.calendarService.GetMonthCalendar(date.Year(), int(date.Month()), uint(roomID), floor)
	if err != nil {
		utils.InternalError(c, "获取日历数据失败")
		return
	}

	utils.Success(c, bookings)
}

func (h *CalendarHandler) GetRoomAvailability(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "会议室ID错误")
		return
	}

	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.BadRequest(c, "日期格式错误")
		return
	}

	availability := h.calendarService.GetRoomAvailability(uint(id), date)
	utils.Success(c, availability)
}
