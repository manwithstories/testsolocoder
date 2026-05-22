package handlers

import (
	"net/http"
	"podcast-manager/internal/database"
	"podcast-manager/internal/models"
	"podcast-manager/internal/services"
	"podcast-manager/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PodcastHandler struct {
	podcastService *services.PodcastService
}

func NewPodcastHandler() *PodcastHandler {
	return &PodcastHandler{
		podcastService: services.NewPodcastService(),
	}
}

type AddPodcastRequest struct {
	FeedURL string `json:"feed_url" binding:"required,url"`
}

func (h *PodcastHandler) AddPodcast(c *gin.Context) {
	var req AddPodcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	podcast, err := h.podcastService.AddPodcast(req.FeedURL)
	if err != nil {
		utils.InternalError(c, "Failed to add podcast: "+err.Error())
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (h *PodcastHandler) GetPodcasts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	podcasts, total, err := h.podcastService.GetPodcastList(page, perPage, search)
	if err != nil {
		utils.InternalError(c, "Failed to get podcasts: "+err.Error())
		return
	}

	meta := utils.Meta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	utils.SuccessResponseWithMeta(c, podcasts, meta)
}

func (h *PodcastHandler) GetPodcast(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid podcast ID")
		return
	}

	podcast, err := h.podcastService.GetPodcastByID(id)
	if err != nil {
		utils.NotFound(c, "Podcast not found")
		return
	}

	stats, err := h.podcastService.GetPodcastStats(id)
	if err != nil {
		utils.InternalError(c, "Failed to get podcast stats: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"podcast": podcast,
		"stats":   stats,
	})
}

type UpdatePodcastRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	CoverImage  string `json:"cover_image"`
}

func (h *PodcastHandler) UpdatePodcast(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid podcast ID")
		return
	}

	var req UpdatePodcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Author != "" {
		updates["author"] = req.Author
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}

	podcast, err := h.podcastService.UpdatePodcast(id, updates)
	if err != nil {
		utils.InternalError(c, "Failed to update podcast: "+err.Error())
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (h *PodcastHandler) DeletePodcast(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid podcast ID")
		return
	}

	if err := h.podcastService.DeletePodcast(id); err != nil {
		utils.InternalError(c, "Failed to delete podcast: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PodcastHandler) RefreshPodcast(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid podcast ID")
		return
	}

	podcast, newEpisodes, err := h.podcastService.RefreshPodcast(id)
	if err != nil {
		utils.InternalError(c, "Failed to refresh podcast: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"podcast":      podcast,
		"new_episodes": newEpisodes,
	})
}

func (h *PodcastHandler) GetNewEpisodesCount(c *gin.Context) {
	var count int64
	err := database.DB.Model(&models.Episode{}).Where("is_new = true").Count(&count).Error
	if err != nil {
		utils.InternalError(c, "Failed to get new episodes count: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"new_episodes_count": count,
	})
}
