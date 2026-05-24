package handler

import (
	"net/http"
	"strconv"
	"survey-platform/internal/dto"
	"survey-platform/internal/service"
	"survey-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type DistributionHandler struct {
	distService *service.DistributionService
}

func NewDistributionHandler(distService *service.DistributionService) *DistributionHandler {
	return &DistributionHandler{distService: distService}
}

func (h *DistributionHandler) CreateLink(c *gin.Context) {
	userID := c.GetUint("user_id")
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.CreateDistributionLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.distService.CreateLink(uint(surveyID), userID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *DistributionHandler) GetByToken(c *gin.Context) {
	token := c.Param("token")

	link, err := h.distService.GetByToken(token)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, link)
}

func (h *DistributionHandler) ListBySurveyID(c *gin.Context) {
	userID := c.GetUint("user_id")
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	links, err := h.distService.ListBySurveyID(uint(surveyID), userID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, links)
}

func (h *DistributionHandler) SendInvitations(c *gin.Context) {
	userID := c.GetUint("user_id")
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.SendInvitationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.distService.SendInvitations(uint(surveyID), userID, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *DistributionHandler) ListInvitations(c *gin.Context) {
	userID := c.GetUint("user_id")
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invitations, total, err := h.distService.ListInvitations(uint(surveyID), userID, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     invitations,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *DistributionHandler) DeleteLink(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid link id")
		return
	}

	if err := h.distService.DeleteLink(uint(id), userID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}
