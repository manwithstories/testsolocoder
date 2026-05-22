package handlers

import (
	"qa-platform/models"
	"qa-platform/repository"
	"qa-platform/services"
	"qa-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionService *services.QuestionService
}

func NewQuestionHandler() *QuestionHandler {
	return &QuestionHandler{
		questionService: services.NewQuestionService(),
	}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	userID := c.GetUint("userId")

	var req services.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	question, err := h.questionService.CreateQuestion(userID, &req)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, question)
}

func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	currentUserID := c.GetUint("userId")

	question, err := h.questionService.GetQuestionByID(uint(id), currentUserID)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.QuestionNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, question)
}

func (h *QuestionHandler) GetQuestionList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("categoryId", "0"), 10, 64)
	tagID, _ := strconv.ParseUint(c.DefaultQuery("tagId", "0"), 10, 64)
	userID, _ := strconv.ParseUint(c.DefaultQuery("userId", "0"), 10, 64)
	keyword := c.Query("keyword")
	sort := c.Query("sort")

	query := services.QuestionListQuery{
		Page:       page,
		PageSize:   pageSize,
		CategoryID: uint(categoryID),
		TagID:      uint(tagID),
		UserID:     uint(userID),
		Keyword:    keyword,
		Sort:       sort,
	}

	questions, total, err := h.questionService.GetQuestionList(query)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     questions,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	var req services.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.questionService.UpdateQuestion(uint(id), userID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	userID := c.GetUint("userId")
	role := c.GetString("role")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.questionService.DeleteQuestion(uint(id), userID, role); err != nil {
		utils.ErrorResponseWithMessage(c, utils.Forbidden, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *QuestionHandler) AcceptAnswer(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	var req struct {
		AnswerID uint `json:"answerId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.questionService.AcceptAnswer(uint(id), req.AnswerID, userID); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *QuestionHandler) LikeQuestion(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.questionService.LikeQuestion(uint(id), userID); err != nil {
		utils.ErrorResponseWithMessage(c, utils.QuestionNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *QuestionHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	repository.DB.Where("status = ?", "active").Order("sort_order ASC").Find(&categories)
	utils.SuccessResponse(c, categories)
}

func (h *QuestionHandler) GetTags(c *gin.Context) {
	var tags []models.Tag
	repository.DB.Order("usage_count DESC").Limit(50).Find(&tags)
	utils.SuccessResponse(c, tags)
}

func (h *QuestionHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		SortOrder   int    `json:"sortOrder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	if err := repository.DB.Create(&category).Error; err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, category)
}

func (h *QuestionHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	var existingTag models.Tag
	result := repository.DB.Where("name = ?", req.Name).First(&existingTag)
	if result.Error == nil {
		utils.SuccessResponse(c, existingTag)
		return
	}

	tag := models.Tag{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := repository.DB.Create(&tag).Error; err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, tag)
}

type AnswerHandler struct {
	answerService *services.AnswerService
}

func NewAnswerHandler() *AnswerHandler {
	return &AnswerHandler{
		answerService: services.NewAnswerService(),
	}
}

func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	userID := c.GetUint("userId")

	var req services.CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	answer, err := h.answerService.CreateAnswer(userID, &req)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, answer)
}

func (h *AnswerHandler) GetAnswer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	answer, err := h.answerService.GetAnswerByID(uint(id))
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.AnswerNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, answer)
}

func (h *AnswerHandler) UpdateAnswer(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	var req services.UpdateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.answerService.UpdateAnswer(uint(id), userID, &req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AnswerHandler) DeleteAnswer(c *gin.Context) {
	userID := c.GetUint("userId")
	role := c.GetString("role")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.answerService.DeleteAnswer(uint(id), userID, role); err != nil {
		utils.ErrorResponseWithMessage(c, utils.Forbidden, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AnswerHandler) LikeAnswer(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.answerService.LikeAnswer(uint(id), userID); err != nil {
		utils.ErrorResponseWithMessage(c, utils.AnswerNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *AnswerHandler) DislikeAnswer(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.answerService.DislikeAnswer(uint(id), userID); err != nil {
		utils.ErrorResponseWithMessage(c, utils.AnswerNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentService: services.NewCommentService(),
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID := c.GetUint("userId")

	var req services.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	comment, err := h.commentService.CreateComment(userID, &req)
	if err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID := c.GetUint("userId")
	role := c.GetString("role")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.commentService.DeleteComment(uint(id), userID, role); err != nil {
		utils.ErrorResponseWithMessage(c, utils.Forbidden, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *CommentHandler) LikeComment(c *gin.Context) {
	userID := c.GetUint("userId")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.commentService.LikeComment(uint(id), userID); err != nil {
		utils.ErrorResponseWithMessage(c, utils.CommentNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}
