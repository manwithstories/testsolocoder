package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"museum-server/internal/dto"
	"museum-server/internal/services"
	"museum-server/pkg/response"
)

type ReservationHandler struct {
	reservationService *services.ReservationService
}

func NewReservationHandler(reservationService *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req dto.ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	reservation, err := h.reservationService.Create(uid, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, reservation)
}

func (h *ReservationHandler) Confirm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid reservation ID")
		return
	}

	if err := h.reservationService.Confirm(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ReservationHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid reservation ID")
		return
	}

	var req dto.ReservationCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.reservationService.Cancel(uint(id), req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ReservationHandler) Reschedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid reservation ID")
		return
	}

	var req struct {
		NewTimeSlotID uint `json:"new_time_slot_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.reservationService.Reschedule(uint(id), req.NewTimeSlotID); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ReservationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid reservation ID")
		return
	}

	reservation, err := h.reservationService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, reservation)
}

func (h *ReservationHandler) GetByQRCode(c *gin.Context) {
	qrCode := c.Param("qr_code")

	reservation, err := h.reservationService.GetByQRCode(qrCode)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, reservation)
}

func (h *ReservationHandler) ListByUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reservations, total, err := h.reservationService.ListByUser(uid, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, reservations, total, page, pageSize)
}

func (h *ReservationHandler) ListByExhibition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	reservations, total, err := h.reservationService.ListByExhibition(uint(id), page, pageSize, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, reservations, total, page, pageSize)
}

func (h *ReservationHandler) CheckIn(c *gin.Context) {
	var req struct {
		QRCode string `json:"qr_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.reservationService.CheckIn(req.QRCode); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ReservationHandler) CheckOut(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid reservation ID")
		return
	}

	if err := h.reservationService.CheckOut(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ReservationHandler) RateVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid reservation ID")
		return
	}

	var req struct {
		Rating   int    `json:"rating" binding:"required,min=1,max=5"`
		Comment  string `json:"comment"`
		Favorite bool   `json:"favorite"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.reservationService.RateVisit(uint(id), req.Rating, req.Comment, req.Favorite); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ReservationHandler) ListVisitRecords(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	records, total, err := h.reservationService.ListVisitRecords(uid, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, records, total, page, pageSize)
}

func (h *ReservationHandler) GetUserVisitStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	stats, err := h.reservationService.GetUserVisitStats(uid)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, stats)
}

func (h *ReservationHandler) GetReservationStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	status, err := h.reservationService.GetCachedReservationStatus(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, status)
}
