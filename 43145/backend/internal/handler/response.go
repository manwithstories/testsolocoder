package handler

import (
	"net/http"
	"strconv"
	"survey-platform/internal/dto"
	"survey-platform/internal/service"
	"survey-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type ResponseHandler struct {
	responseService *service.ResponseService
	distService     *service.DistributionService
}

func NewResponseHandler(
	responseService *service.ResponseService,
	distService *service.DistributionService,
) *ResponseHandler {
	return &ResponseHandler{
		responseService: responseService,
		distService:     distService,
	}
}

func (h *ResponseHandler) StartResponse(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}

	sessionID := h.responseService.GenerateSessionID()
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	token := c.Query("token")
	var distributionID *uint
	if token != "" {
		link, err := h.distService.GetByToken(token)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		distID := link.ID
		distributionID = &distID
	}

	response, err := h.responseService.StartResponse(uint(surveyID), userID, sessionID, ipAddress, userAgent, distributionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"response_id": response.ID,
		"session_id":  response.SessionID,
	})
}

func (h *ResponseHandler) SaveProgress(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.SubmitResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.responseService.SaveProgress(uint(surveyID), req.SessionID, req.Answers); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ResponseHandler) SubmitResponse(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.SubmitResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.responseService.SubmitResponse(uint(surveyID), &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ResponseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid response id")
		return
	}

	result, err := h.responseService.GetDetail(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *ResponseHandler) List(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var query dto.ResponseListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	responses, total, err := h.responseService.List(uint(surveyID), &query)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     responses,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

func (h *ResponseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid response id")
		return
	}

	if err := h.responseService.Delete(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ResponseHandler) ValidateAccess(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.AccessSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uint)
		userID = &id
	}

	if err := h.responseService.ValidateSurveyAccess(uint(surveyID), req.Password, userID); err != nil {
		utils.Error(c, http.StatusForbidden, err.Error())
		return
	}

	utils.Success(c, nil)
}
