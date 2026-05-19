package handler

import (
	"net/http"
	"strconv"
	"venue-booking/internal/dto"
	"venue-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type VenueHandler struct {
	venueService *service.VenueService
	logService   *service.OperationLogService
}

func NewVenueHandler() *VenueHandler {
	return &VenueHandler{
		venueService: service.NewVenueService(),
		logService:   service.NewOperationLogService(),
	}
}

func (h *VenueHandler) Create(c *gin.Context) {
	var req dto.VenueCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	userID, _ := c.Get("userID")
	venue, err := h.venueService.Create(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to create venue"))
		return
	}

	h.logService.Log(c, userID.(uint), "create_venue", "venue", map[string]interface{}{
		"venue_id":   venue.ID,
		"venue_name": venue.Name,
	})

	c.JSON(http.StatusOK, dto.Success(venue))
}

func (h *VenueHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid venue ID"))
		return
	}

	venue, err := h.venueService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "Venue not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(venue))
}

func (h *VenueHandler) List(c *gin.Context) {
	var req dto.VenueListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	venues, total, err := h.venueService.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get venues"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     venues,
	}))
}

func (h *VenueHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid venue ID"))
		return
	}

	var req dto.VenueUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	venue, err := h.venueService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to update venue"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "update_venue", "venue", map[string]interface{}{
		"venue_id": id,
	})

	c.JSON(http.StatusOK, dto.Success(venue))
}

func (h *VenueHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid venue ID"))
		return
	}

	err = h.venueService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to delete venue"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "delete_venue", "venue", map[string]interface{}{
		"venue_id": id,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *VenueHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid venue ID"))
		return
	}

	var req dto.VenueStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err = h.venueService.UpdateStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to update venue status"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "update_venue_status", "venue", map[string]interface{}{
		"venue_id": id,
		"status":   req.Status,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *VenueHandler) SetPrice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid venue ID"))
		return
	}

	var req dto.VenuePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err = h.venueService.SetPrice(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to set venue price"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "set_venue_price", "venue", map[string]interface{}{
		"venue_id":    id,
		"day_of_week": req.DayOfWeek,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *VenueHandler) GetAvailability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid venue ID"))
		return
	}

	var req dto.VenueAvailabilityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	availability, err := h.venueService.GetAvailability(uint(id), req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(availability))
}
