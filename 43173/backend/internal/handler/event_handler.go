package handler

import (
	"strconv"

	"music-platform/internal/service"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/jwt"
	"music-platform/pkg/response"
	"music-platform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		eventService: service.NewEventService(),
	}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	event, err := h.eventService.CreateEvent(userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "创建演出失败")
		return
	}

	response.Success(c, event)
}

func (h *EventHandler) GetEventByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	event, err := h.eventService.GetEventByID(uint(id))
	if err != nil {
		response.NotFound(c, "演出不存在")
		return
	}

	response.Success(c, event)
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)
	keyword := c.Query("keyword")
	city := c.Query("city")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	artistID, _ := strconv.ParseUint(c.DefaultQuery("artist_id", "0"), 10, 64)

	events, total, err := h.eventService.ListEvents(page, pageSize, keyword, uint(artistID), city, status)
	if err != nil {
		response.InternalError(c, "获取演出列表失败")
		return
	}

	response.Page(c, events, total, page, pageSize)
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req service.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err = h.eventService.UpdateEvent(uint(id), userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "更新演出失败")
		return
	}

	response.Success(c, nil)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.eventService.DeleteEvent(uint(id), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "删除演出失败")
		return
	}

	response.Success(c, nil)
}

func (h *EventHandler) PublishEvent(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.eventService.PublishEvent(uint(id), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "发布演出失败")
		return
	}

	response.Success(c, nil)
}

func (h *EventHandler) PurchaseTicket(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.PurchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	order, err := h.eventService.PurchaseTicket(userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "购票失败")
		return
	}

	response.Success(c, order)
}

func (h *EventHandler) GetOrderByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	order, err := h.eventService.GetOrderByID(uint(id))
	if err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	response.Success(c, order)
}

func (h *EventHandler) GetOrdersByUser(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	orders, total, err := h.eventService.GetOrdersByUser(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取订单列表失败")
		return
	}

	response.Page(c, orders, total, page, pageSize)
}

func (h *EventHandler) GetOrdersByArtist(c *gin.Context) {
	artistID, err := strconv.ParseUint(c.Param("artist_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	page, pageSize := utils.GetPageAndPageSize(c)

	orders, total, err := h.eventService.GetOrdersByArtist(uint(artistID), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取订单列表失败")
		return
	}

	response.Page(c, orders, total, page, pageSize)
}

func (h *EventHandler) GetTicketsByUser(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	tickets, total, err := h.eventService.GetTicketsByUser(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取票列表失败")
		return
	}

	response.Page(c, tickets, total, page, pageSize)
}

func (h *EventHandler) UseTicket(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.eventService.UseTicket(uint(id), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "核销失败")
		return
	}

	response.Success(c, nil)
}

func (h *EventHandler) GetEventStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	stats, err := h.eventService.GetEventStats(uint(id))
	if err != nil {
		response.InternalError(c, "获取统计失败")
		return
	}

	response.Success(c, stats)
}

func (h *EventHandler) GetSeatAvailability(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	availability, err := h.eventService.GetSeatAvailability(uint(id))
	if err != nil {
		response.InternalError(c, "获取座位信息失败")
		return
	}

	response.Success(c, availability)
}
