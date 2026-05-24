package handler

import (
	"net/http"
	"strconv"
	"survey-platform/internal/dto"
	"survey-platform/internal/service"
	"survey-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionService *service.QuestionService
}

func NewQuestionHandler(questionService *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

func (h *QuestionHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	question, err := h.questionService.Create(uint(surveyID), userID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, question)
}

func (h *QuestionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid question id")
		return
	}

	question, err := h.questionService.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, question)
}

func (h *QuestionHandler) GetBySurveyID(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	questions, err := h.questionService.GetBySurveyID(uint(surveyID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, questions)
}

func (h *QuestionHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid question id")
		return
	}

	var req dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.questionService.Update(uint(id), userID, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *QuestionHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid question id")
		return
	}

	if err := h.questionService.Delete(uint(id), userID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *QuestionHandler) Reorder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid question id")
		return
	}

	var req dto.ReorderQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.questionService.Reorder(uint(id), userID, req.OrderIndex); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *QuestionHandler) BatchCreate(c *gin.Context) {
	userID := c.GetUint("user_id")
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	var req dto.BatchUpdateQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.questionService.BatchCreate(uint(surveyID), userID, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}
