package handlers

import (
	"podcast-manager/internal/services"
	"podcast-manager/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EpisodeHandler struct {
	episodeService  *services.EpisodeService
	playbackService *services.PlaybackService
	historyService  *services.HistoryService
}

func NewEpisodeHandler() *EpisodeHandler {
	return &EpisodeHandler{
		episodeService:  services.NewEpisodeService(),
		playbackService: services.NewPlaybackService(),
		historyService:  services.NewHistoryService(),
	}
}

func (h *EpisodeHandler) GetEpisodes(c *gin.Context) {
	podcastIDStr := c.Query("podcast_id")
	var podcastID uuid.UUID
	if podcastIDStr != "" {
		id, err := uuid.Parse(podcastIDStr)
		if err != nil {
			utils.BadRequest(c, "Invalid podcast ID")
			return
		}
		podcastID = id
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	episodes, total, err := h.episodeService.GetEpisodeList(podcastID, page, perPage, search)
	if err != nil {
		utils.InternalError(c, "Failed to get episodes: "+err.Error())
		return
	}

	meta := utils.Meta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	utils.SuccessResponseWithMeta(c, episodes, meta)
}

func (h *EpisodeHandler) GetEpisode(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid episode ID")
		return
	}

	episode, err := h.episodeService.GetEpisodeByID(id)
	if err != nil {
		utils.NotFound(c, "Episode not found")
		return
	}

	progress, err := h.playbackService.GetPlaybackProgress(id)
	if err != nil {
		utils.InternalError(c, "Failed to get playback progress: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"episode":  episode,
		"progress": progress,
	})
}

type UpdateProgressRequest struct {
	CurrentTime float64 `json:"current_time" binding:"required,min=0"`
	Duration    int     `json:"duration" binding:"min=0"`
}

func (h *EpisodeHandler) UpdatePlaybackProgress(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid episode ID")
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	progress, err := h.playbackService.UpdatePlaybackProgress(id, req.CurrentTime, req.Duration)
	if err != nil {
		utils.InternalError(c, "Failed to update playback progress: "+err.Error())
		return
	}

	utils.SuccessResponse(c, progress)
}

func (h *EpisodeHandler) GetPlaybackProgress(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid episode ID")
		return
	}

	progress, err := h.playbackService.GetPlaybackProgress(id)
	if err != nil {
		utils.InternalError(c, "Failed to get playback progress: "+err.Error())
		return
	}

	utils.SuccessResponse(c, progress)
}

func (h *EpisodeHandler) MarkAsCompleted(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid episode ID")
		return
	}

	if err := h.playbackService.MarkAsCompleted(id); err != nil {
		utils.InternalError(c, "Failed to mark as completed: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Episode marked as completed"})
}

type AddHistoryRequest struct {
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
	Duration   int       `json:"duration" binding:"required,min=0"`
	Completion float64   `json:"completion" binding:"min=0,max=1"`
}

func (h *EpisodeHandler) AddListeningHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid episode ID")
		return
	}

	var req AddHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	history, err := h.historyService.AddListeningHistory(id, req.StartTime, req.EndTime, req.Duration, req.Completion)
	if err != nil {
		utils.InternalError(c, "Failed to add listening history: "+err.Error())
		return
	}

	utils.SuccessResponse(c, history)
}

func (h *EpisodeHandler) GetListeningHistory(c *gin.Context) {
	podcastIDStr := c.Query("podcast_id")
	var podcastID uuid.UUID
	if podcastIDStr != "" {
		id, err := uuid.Parse(podcastIDStr)
		if err != nil {
			utils.BadRequest(c, "Invalid podcast ID")
			return
		}
		podcastID = id
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	completedStr := c.Query("completed")

	var completed *bool
	if completedStr != "" {
		b, err := strconv.ParseBool(completedStr)
		if err == nil {
			completed = &b
		}
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	histories, total, err := h.historyService.GetListeningHistory(page, perPage, podcastID, startDate, endDate, completed)
	if err != nil {
		utils.InternalError(c, "Failed to get listening history: "+err.Error())
		return
	}

	meta := utils.Meta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	utils.SuccessResponseWithMeta(c, histories, meta)
}
