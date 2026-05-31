package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"museum-server/internal/dto"
	"museum-server/internal/services"
	"museum-server/pkg/response"
)

type GuideHandler struct {
	guideService *services.GuideService
}

func NewGuideHandler(guideService *services.GuideService) *GuideHandler {
	return &GuideHandler{guideService: guideService}
}

func (h *GuideHandler) CreateSchedule(c *gin.Context) {
	guideID, _ := c.Get("user_id")
	gid := guideID.(uint)

	var req dto.GuideScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	schedule, err := h.guideService.CreateSchedule(gid, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, schedule)
}

func (h *GuideHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid schedule ID")
		return
	}

	var req dto.GuideScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.guideService.UpdateSchedule(uint(id), &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GuideHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid schedule ID")
		return
	}

	if err := h.guideService.DeleteSchedule(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GuideHandler) ListSchedules(c *gin.Context) {
	guideID, _ := c.Get("user_id")
	gid := guideID.(uint)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	schedules, err := h.guideService.ListSchedules(gid, startDate, endDate)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, schedules)
}

func (h *GuideHandler) CreateContent(c *gin.Context) {
	var req dto.GuideContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	content, err := h.guideService.CreateContent(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, content)
}

func (h *GuideHandler) UpdateContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid content ID")
		return
	}

	var req dto.GuideContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.guideService.UpdateContent(uint(id), &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GuideHandler) DeleteContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid content ID")
		return
	}

	if err := h.guideService.DeleteContent(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GuideHandler) ListContents(c *gin.Context) {
	collectionID, _ := strconv.ParseUint(c.Query("collection_id"), 10, 32)
	exhibitionID, _ := strconv.ParseUint(c.Query("exhibition_id"), 10, 32)
	language := c.Query("language")

	contents, err := h.guideService.ListContents(uint(collectionID), uint(exhibitionID), language)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, contents)
}

type ResearchHandler struct {
	researchService *services.ResearchService
}

func NewResearchHandler(researchService *services.ResearchService) *ResearchHandler {
	return &ResearchHandler{researchService: researchService}
}

func (h *ResearchHandler) CreateApplication(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req dto.ResearchApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	app, err := h.researchService.CreateApplication(uid, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, app)
}

func (h *ResearchHandler) ReviewApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid application ID")
		return
	}

	reviewerID, _ := c.Get("user_id")
	rid := reviewerID.(uint)

	var req dto.ApplicationReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.researchService.ReviewApplication(uint(id), rid, &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ResearchHandler) ListApplications(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	apps, total, err := h.researchService.ListApplications(status, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, apps, total, page, pageSize)
}

func (h *ResearchHandler) ListMyApplications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	apps, total, err := h.researchService.ListApplicationsByUser(uid, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, apps, total, page, pageSize)
}

func (h *ResearchHandler) GetApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid application ID")
		return
	}

	app, err := h.researchService.GetApplication(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, app)
}

type StatisticsHandler struct {
	statisticsService *services.StatisticsService
}

func NewStatisticsHandler(statisticsService *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statisticsService: statisticsService}
}

func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	stats, err := h.statisticsService.GetStatistics(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, stats)
}

func (h *StatisticsHandler) ExportExcel(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	data, err := h.statisticsService.ExportExcel(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=museum_statistics.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *StatisticsHandler) ExportPDF(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	data, err := h.statisticsService.ExportPDF(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, data)
}

type MuseumHandler struct {
	museumService *services.MuseumService
}

func NewMuseumHandler(museumService *services.MuseumService) *MuseumHandler {
	return &MuseumHandler{museumService: museumService}
}

func (h *MuseumHandler) Create(c *gin.Context) {
	var req dto.MuseumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	museum, err := h.museumService.Create(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, museum)
}

func (h *MuseumHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid museum ID")
		return
	}

	var req dto.MuseumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.museumService.Update(uint(id), &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *MuseumHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid museum ID")
		return
	}

	if err := h.museumService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *MuseumHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid museum ID")
		return
	}

	museum, err := h.museumService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, museum)
}

func (h *MuseumHandler) List(c *gin.Context) {
	museums, err := h.museumService.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, museums)
}

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "No file uploaded")
		return
	}

	filename := file.Filename
	filePath := h.uploadDir + "/" + filename

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to save file")
		return
	}

	response.Success(c, gin.H{
		"url":      "/uploads/" + filename,
		"filename": filename,
		"size":     file.Size,
	})
}
