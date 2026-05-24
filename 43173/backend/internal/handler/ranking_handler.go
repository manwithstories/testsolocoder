package handler

import (
	"strconv"

	"music-platform/internal/service"
	"music-platform/pkg/response"

	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	rankingService *service.RankingService
}

func NewRankingHandler() *RankingHandler {
	return &RankingHandler{
		rankingService: service.NewRankingService(),
	}
}

func (h *RankingHandler) GetRanking(c *gin.Context) {
	rankingType := c.DefaultQuery("type", "daily")
	category := c.DefaultQuery("category", "plays")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	items, err := h.rankingService.GetRanking(rankingType, category, limit)
	if err != nil {
		response.InternalError(c, "获取排行榜失败")
		return
	}

	response.Success(c, items)
}

func (h *RankingHandler) GetWorkRanking(c *gin.Context) {
	workID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	rankingType := c.DefaultQuery("type", "daily")
	category := c.DefaultQuery("category", "plays")

	rank, score, err := h.rankingService.GetWorkRanking(uint(workID), rankingType, category)
	if err != nil {
		response.InternalError(c, "获取排名失败")
		return
	}

	response.Success(c, gin.H{
		"rank":  rank,
		"score": score,
	})
}

func (h *RankingHandler) GetDailyRanking(c *gin.Context) {
	category := c.DefaultQuery("category", "plays")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	items, err := h.rankingService.GetRanking("daily", category, limit)
	if err != nil {
		response.InternalError(c, "获取日榜失败")
		return
	}

	response.Success(c, items)
}

func (h *RankingHandler) GetWeeklyRanking(c *gin.Context) {
	category := c.DefaultQuery("category", "plays")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	items, err := h.rankingService.GetRanking("weekly", category, limit)
	if err != nil {
		response.InternalError(c, "获取周榜失败")
		return
	}

	response.Success(c, items)
}

func (h *RankingHandler) GetMonthlyRanking(c *gin.Context) {
	category := c.DefaultQuery("category", "plays")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	items, err := h.rankingService.GetRanking("monthly", category, limit)
	if err != nil {
		response.InternalError(c, "获取月榜失败")
		return
	}

	response.Success(c, items)
}

func (h *RankingHandler) GetHotRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	items, err := h.rankingService.GetRanking("daily", "hot", limit)
	if err != nil {
		response.InternalError(c, "获取热歌榜失败")
		return
	}

	response.Success(c, items)
}
