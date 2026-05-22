package handlers

import (
	"qa-platform/services"
	"qa-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	rewardService *services.RewardService
}

func NewRewardHandler() *RewardHandler {
	return &RewardHandler{
		rewardService: services.NewRewardService(),
	}
}

func (h *RewardHandler) GetRewardList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	rewards, total, err := h.rewardService.GetRewardList(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     rewards,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Image       string `json:"image"`
		PointsCost  int    `json:"pointsCost" binding:"required"`
		Stock       int    `json:"stock"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	reward, err := h.rewardService.CreateReward(req.Name, req.Description, req.Image, req.PointsCost, req.Stock)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, reward)
}

func (h *RewardHandler) UpdateReward(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponseWithMessage(c, utils.BadRequest, err.Error())
		return
	}

	if err := h.rewardService.UpdateReward(uint(id), updates); err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *RewardHandler) DeleteReward(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.rewardService.DeleteReward(uint(id)); err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *RewardHandler) ExchangeReward(c *gin.Context) {
	userID := c.GetUint("userId")
	rewardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	if err := h.rewardService.ExchangeReward(userID, uint(rewardID)); err != nil {
		utils.ErrorResponseWithMessage(c, utils.InsufficientPoints, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *RewardHandler) GetExchangeList(c *gin.Context) {
	userID := c.GetUint("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	exchanges, total, err := h.rewardService.GetExchangeList(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     exchanges,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *RewardHandler) GetPointLogs(c *gin.Context) {
	userID := c.GetUint("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	logs, total, err := h.rewardService.GetPointLogs(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		searchService: services.NewSearchService(),
	}
}

func (h *SearchHandler) SearchQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("categoryId", "0"), 10, 64)
	tagID, _ := strconv.ParseUint(c.DefaultQuery("tagId", "0"), 10, 64)
	keyword := c.Query("keyword")
	sort := c.Query("sort")

	query := services.SearchQuery{
		Keyword:    keyword,
		CategoryID: uint(categoryID),
		TagID:      uint(tagID),
		Page:       page,
		PageSize:   pageSize,
		Sort:       sort,
	}

	questions, total, err := h.searchService.SearchQuestions(query)
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

func (h *SearchHandler) GetRecommendations(c *gin.Context) {
	userID := c.GetUint("userId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	questions, err := h.searchService.GetRecommendations(userID, limit)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, questions)
}

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: services.NewStatsService(),
	}
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, stats)
}

func (h *StatsHandler) GetActivityReport(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	reports, err := h.statsService.GetActivityReport(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, reports)
}

func (h *StatsHandler) GetAuditReport(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	reports, err := h.statsService.GetAuditReport(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	utils.SuccessResponse(c, reports)
}

func (h *StatsHandler) ExportActivityReport(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	reports, err := h.statsService.GetActivityReport(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	csvContent := "date,newQuestions,newAnswers,newUsers,newComments\n"
	for _, report := range reports {
		csvContent += report.Date + "," +
			strconv.FormatInt(report.NewQuestions, 10) + "," +
			strconv.FormatInt(report.NewAnswers, 10) + "," +
			strconv.FormatInt(report.NewUsers, 10) + "," +
			strconv.FormatInt(report.NewComments, 10) + "\n"
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=activity_report.csv")
	c.String(200, csvContent)
}

func (h *StatsHandler) ExportAuditReport(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		utils.ErrorResponse(c, utils.BadRequest)
		return
	}

	reports, err := h.statsService.GetAuditReport(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, utils.InternalServerError)
		return
	}

	csvContent := "date,reviewedCount,approvedCount,rejectedCount\n"
	for _, report := range reports {
		csvContent += report.Date + "," +
			strconv.FormatInt(report.ReviewedCount, 10) + "," +
			strconv.FormatInt(report.ApprovedCount, 10) + "," +
			strconv.FormatInt(report.RejectedCount, 10) + "\n"
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=audit_report.csv")
	c.String(200, csvContent)
}
