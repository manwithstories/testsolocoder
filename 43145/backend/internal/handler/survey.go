package handler

import (
	"net/http"
	"strconv"
	"survey-platform/internal/dto"
	"survey-platform/internal/service"
	"survey-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type SurveyHandler struct {
	surveyService *service.SurveyService
}

func NewSurveyHandler(surveyService *service.SurveyService) *SurveyHandler {
	return &SurveyHandler{surveyService: surveyService}
}

func (h *SurveyHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	survey, err := h.surveyService.Create(userID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, survey)
}

func (h *SurveyHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	survey, err := h.surveyService.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, survey)
}

func (h *SurveyHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.surveyService.Update(uint(id), userID, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *SurveyHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	if err := h.surveyService.Delete(uint(id), userID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *SurveyHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.SurveyListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	surveys, total, err := h.surveyService.List(userID, &query)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     surveys,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

func (h *SurveyHandler) ListAll(c *gin.Context) {
	var query dto.SurveyListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	surveys, total, err := h.surveyService.ListAll(&query)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     surveys,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

func (h *SurveyHandler) Publish(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.PublishSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.surveyService.Publish(uint(id), userID, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *SurveyHandler) Close(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	if err := h.surveyService.Close(uint(id), userID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *SurveyHandler) Copy(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.CopySurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	survey, err := h.surveyService.Copy(uint(id), userID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, survey)
}
