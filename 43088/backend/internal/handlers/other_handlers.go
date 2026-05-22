package handlers

import (
	"encoding/csv"
	"net/http"
	"podcast-manager/internal/services"
	"podcast-manager/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NoteHandler struct {
	noteService *services.NoteService
}

func NewNoteHandler() *NoteHandler {
	return &NoteHandler{
		noteService: services.NewNoteService(),
	}
}

type AddNoteRequest struct {
	Timestamp float64  `json:"timestamp" binding:"required,min=0"`
	Content   string   `json:"content" binding:"required"`
	Tags      []string `json:"tags"`
}

func (h *NoteHandler) AddNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("episode_id"))
	if err != nil {
		utils.BadRequest(c, "Invalid episode ID")
		return
	}

	var req AddNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	note, err := h.noteService.AddNote(id, req.Timestamp, req.Content, req.Tags)
	if err != nil {
		utils.InternalError(c, "Failed to add note: "+err.Error())
		return
	}

	utils.SuccessResponse(c, note)
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	id, err := uuid.Parse(c.Param("episode_id"))
	if err != nil {
		utils.BadRequest(c, "Invalid episode ID")
		return
	}

	search := c.Query("search")
	tag := c.Query("tag")

	notes, err := h.noteService.GetNotes(id, search, tag)
	if err != nil {
		utils.InternalError(c, "Failed to get notes: "+err.Error())
		return
	}

	utils.SuccessResponse(c, notes)
}

type UpdateNoteRequest struct {
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid note ID")
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	note, err := h.noteService.UpdateNote(id, req.Content, req.Tags)
	if err != nil {
		utils.InternalError(c, "Failed to update note: "+err.Error())
		return
	}

	utils.SuccessResponse(c, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid note ID")
		return
	}

	if err := h.noteService.DeleteNote(id); err != nil {
		utils.InternalError(c, "Failed to delete note: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NoteHandler) SearchNotes(c *gin.Context) {
	search := c.Query("q")
	if search == "" {
		utils.BadRequest(c, "Search query is required")
		return
	}

	notes, err := h.noteService.SearchNotes(search)
	if err != nil {
		utils.InternalError(c, "Failed to search notes: "+err.Error())
		return
	}

	utils.SuccessResponse(c, notes)
}

type PlaylistHandler struct {
	playlistService *services.PlaylistService
}

func NewPlaylistHandler() *PlaylistHandler {
	return &PlaylistHandler{
		playlistService: services.NewPlaylistService(),
	}
}

type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

func (h *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	var req CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	playlist, err := h.playlistService.CreatePlaylist(req.Name, req.Description, req.CoverImage)
	if err != nil {
		utils.InternalError(c, "Failed to create playlist: "+err.Error())
		return
	}

	utils.SuccessResponse(c, playlist)
}

func (h *PlaylistHandler) GetPlaylists(c *gin.Context) {
	playlists, err := h.playlistService.GetPlaylistList()
	if err != nil {
		utils.InternalError(c, "Failed to get playlists: "+err.Error())
		return
	}

	utils.SuccessResponse(c, playlists)
}

func (h *PlaylistHandler) GetPlaylist(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid playlist ID")
		return
	}

	playlist, err := h.playlistService.GetPlaylistByID(id)
	if err != nil {
		utils.NotFound(c, "Playlist not found")
		return
	}

	utils.SuccessResponse(c, playlist)
}

type UpdatePlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

func (h *PlaylistHandler) UpdatePlaylist(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid playlist ID")
		return
	}

	var req UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	playlist, err := h.playlistService.UpdatePlaylist(id, req.Name, req.Description, req.CoverImage)
	if err != nil {
		utils.InternalError(c, "Failed to update playlist: "+err.Error())
		return
	}

	utils.SuccessResponse(c, playlist)
}

func (h *PlaylistHandler) DeletePlaylist(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid playlist ID")
		return
	}

	if err := h.playlistService.DeletePlaylist(id); err != nil {
		utils.InternalError(c, "Failed to delete playlist: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

type AddToPlaylistRequest struct {
	EpisodeID uuid.UUID `json:"episode_id" binding:"required"`
}

func (h *PlaylistHandler) AddEpisodeToPlaylist(c *gin.Context) {
	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid playlist ID")
		return
	}

	var req AddToPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	item, err := h.playlistService.AddEpisodeToPlaylist(playlistID, req.EpisodeID)
	if err != nil {
		utils.InternalError(c, "Failed to add episode to playlist: "+err.Error())
		return
	}

	utils.SuccessResponse(c, item)
}

func (h *PlaylistHandler) RemoveEpisodeFromPlaylist(c *gin.Context) {
	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid playlist ID")
		return
	}

	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		utils.BadRequest(c, "Invalid item ID")
		return
	}

	if err := h.playlistService.RemoveEpisodeFromPlaylist(playlistID, itemID); err != nil {
		utils.InternalError(c, "Failed to remove episode from playlist: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

type ReorderPlaylistRequest struct {
	ItemIDs []uuid.UUID `json:"item_ids" binding:"required"`
}

func (h *PlaylistHandler) ReorderPlaylistItems(c *gin.Context) {
	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid playlist ID")
		return
	}

	var req ReorderPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.playlistService.ReorderPlaylistItems(playlistID, req.ItemIDs); err != nil {
		utils.InternalError(c, "Failed to reorder playlist items: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Playlist reordered successfully"})
}

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: services.NewStatsService(),
	}
}

func (h *StatsHandler) GetListeningStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 || days > 365 {
		days = 30
	}

	stats, err := h.statsService.GetListeningStats(days)
	if err != nil {
		utils.InternalError(c, "Failed to get listening stats: "+err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}

func (h *StatsHandler) GetPodcastDistribution(c *gin.Context) {
	distribution, err := h.statsService.GetPodcastDistribution()
	if err != nil {
		utils.InternalError(c, "Failed to get podcast distribution: "+err.Error())
		return
	}

	utils.SuccessResponse(c, distribution)
}

func (h *StatsHandler) GetCompletionStats(c *gin.Context) {
	stats, err := h.statsService.GetCompletionStats()
	if err != nil {
		utils.InternalError(c, "Failed to get completion stats: "+err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}

func (h *StatsHandler) GetListeningHabits(c *gin.Context) {
	habits, err := h.statsService.GetListeningHabits()
	if err != nil {
		utils.InternalError(c, "Failed to get listening habits: "+err.Error())
		return
	}

	utils.SuccessResponse(c, habits)
}

type ImportExportHandler struct {
	importExportService *services.ImportExportService
}

func NewImportExportHandler() *ImportExportHandler {
	return &ImportExportHandler{
		importExportService: services.NewImportExportService(),
	}
}

func (h *ImportExportHandler) ExportOPML(c *gin.Context) {
	opml, err := h.importExportService.ExportOPML()
	if err != nil {
		utils.InternalError(c, "Failed to export OPML: "+err.Error())
		return
	}

	filename := "podcasts_" + time.Now().Format("20060102_150405") + ".opml"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, opml)
}

func (h *ImportExportHandler) ImportOPML(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "Failed to get file: "+err.Error())
		return
	}
	defer file.Close()

	count, err := h.importExportService.ImportOPML(file)
	if err != nil {
		utils.InternalError(c, "Failed to import OPML: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"imported_count": count,
		"message":        "OPML imported successfully",
	})
}

func (h *ImportExportHandler) ExportHistoryCSV(c *gin.Context) {
	records, err := h.importExportService.ExportHistoryCSV()
	if err != nil {
		utils.InternalError(c, "Failed to export history: "+err.Error())
		return
	}

	filename := "listening_history_" + time.Now().Format("20060102_150405") + ".csv"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv; charset=utf-8")

	writer := csv.NewWriter(c.Writer)
	writer.UseCRLF = true
	writer.WriteAll(records)
	writer.Flush()
}

func (h *ImportExportHandler) ExportNotesCSV(c *gin.Context) {
	records, err := h.importExportService.ExportNotesCSV()
	if err != nil {
		utils.InternalError(c, "Failed to export notes: "+err.Error())
		return
	}

	filename := "notes_" + time.Now().Format("20060102_150405") + ".csv"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv; charset=utf-8")

	writer := csv.NewWriter(c.Writer)
	writer.UseCRLF = true
	writer.WriteAll(records)
	writer.Flush()
}

func (h *ImportExportHandler) ExportNotesMarkdown(c *gin.Context) {
	episodeIDStr := c.Query("episode_id")
	var episodeID uuid.UUID
	if episodeIDStr != "" {
		id, err := uuid.Parse(episodeIDStr)
		if err != nil {
			utils.BadRequest(c, "Invalid episode ID")
			return
		}
		episodeID = id
	}

	content, err := h.importExportService.ExportNotesMarkdown(episodeID)
	if err != nil {
		utils.InternalError(c, "Failed to export notes: "+err.Error())
		return
	}

	filename := "podcast_notes_" + time.Now().Format("20060102_150405") + ".md"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/markdown; charset=utf-8")
	c.String(http.StatusOK, content)
}
