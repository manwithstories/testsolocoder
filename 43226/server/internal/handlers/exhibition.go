package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"museum-server/internal/dto"
	"museum-server/internal/services"
	"museum-server/pkg/response"
)

type ExhibitionHandler struct {
	exhibitionService *services.ExhibitionService
}

func NewExhibitionHandler(exhibitionService *services.ExhibitionService) *ExhibitionHandler {
	return &ExhibitionHandler{exhibitionService: exhibitionService}
}

func (h *ExhibitionHandler) Create(c *gin.Context) {
	var req dto.ExhibitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	museumID := c.GetUint("museum_id")
	if museumID == 0 {
		museumID = 1
	}

	exhibition, err := h.exhibitionService.Create(museumID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, exhibition)
}

func (h *ExhibitionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	exhibition, err := h.exhibitionService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, exhibition)
}

func (h *ExhibitionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	var req dto.ExhibitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.exhibitionService.Update(uint(id), &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ExhibitionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	if err := h.exhibitionService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ExhibitionHandler) List(c *gin.Context) {
	var query dto.ExhibitionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	exhibitions, total, err := h.exhibitionService.List(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, exhibitions, total, query.Page, query.PageSize)
}

func (h *ExhibitionHandler) AddCollections(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	var req struct {
		CollectionIDs []uint `json:"collection_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.exhibitionService.AddCollections(uint(id), req.CollectionIDs); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ExhibitionHandler) RemoveCollections(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	var req struct {
		CollectionIDs []uint `json:"collection_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.exhibitionService.RemoveCollections(uint(id), req.CollectionIDs); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ExhibitionHandler) GetCollections(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	collections, err := h.exhibitionService.GetCollections(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, collections)
}

func (h *ExhibitionHandler) CreateTimeSlot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	var req dto.TimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	_ = id
	slot, err := h.exhibitionService.CreateTimeSlot(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, slot)
}

func (h *ExhibitionHandler) BatchCreateTimeSlots(c *gin.Context) {
	var req dto.BatchTimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.exhibitionService.BatchCreateTimeSlots(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ExhibitionHandler) ListTimeSlots(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid exhibition ID")
		return
	}

	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		date, _ = time.Parse("2006-01-02", dateStr)
	}

	slots, err := h.exhibitionService.ListTimeSlots(uint(id), date)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, slots)
}

func (h *ExhibitionHandler) GetHotExhibitions(c *gin.Context) {
	exhibitions, err := h.exhibitionService.GetHotExhibitions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, exhibitions)
}
