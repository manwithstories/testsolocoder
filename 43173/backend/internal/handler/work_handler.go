package handler

import (
	"io"
	"strconv"

	"music-platform/internal/service"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/jwt"
	"music-platform/pkg/response"
	"music-platform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WorkHandler struct {
	workService *service.WorkService
}

func NewWorkHandler() *WorkHandler {
	return &WorkHandler{
		workService: service.NewWorkService(),
	}
}

func (h *WorkHandler) UploadWork(c *gin.Context) {
	userID := jwt.GetUserID(c)

	title := c.PostForm("title")
	artistName := c.PostForm("artist_name")
	genre := c.PostForm("genre")
	description := c.PostForm("description")

	if title == "" || artistName == "" {
		response.BadRequest(c, "标题和艺术家名称不能为空")
		return
	}

	audioFile, audioHeader, err := c.Request.FormFile("audio")
	if err != nil {
		response.BadRequest(c, "请上传音频文件")
		return
	}
	defer audioFile.Close()

	audioExt := utils.GetFileExt(audioHeader.Filename)
	if !utils.IsValidAudioExt(audioExt) {
		response.BadRequest(c, "不支持的音频格式")
		return
	}

	if audioHeader.Size > 50<<20 {
		response.BadRequest(c, "音频文件大小不能超过50MB")
		return
	}

	var coverFile io.Reader
	var coverFilename string
	cover, coverHeader, err := c.Request.FormFile("cover")
	if err == nil {
		coverFile = cover
		coverFilename = coverHeader.Filename
		defer cover.Close()
	}

	req := &service.UploadWorkRequest{
		Title:       title,
		ArtistName:  artistName,
		Genre:       genre,
		Description: description,
	}

	work, err := h.workService.UploadWork(userID, req, audioFile, coverFile, audioHeader.Filename, coverFilename)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "上传失败")
		return
	}

	response.Success(c, work)
}

func (h *WorkHandler) GetWorkByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	work, err := h.workService.GetWorkByID(uint(id))
	if err != nil {
		response.NotFound(c, "作品不存在")
		return
	}

	response.Success(c, work)
}

func (h *WorkHandler) ListWorks(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)
	keyword := c.Query("keyword")
	genre := c.Query("genre")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	artistID, _ := strconv.ParseUint(c.DefaultQuery("artist_id", "0"), 10, 64)

	works, total, err := h.workService.ListWorks(page, pageSize, keyword, uint(artistID), genre, status)
	if err != nil {
		response.InternalError(c, "获取作品列表失败")
		return
	}

	response.Page(c, works, total, page, pageSize)
}

func (h *WorkHandler) GetArtistWorks(c *gin.Context) {
	artistID, err := strconv.ParseUint(c.Param("artist_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	page, pageSize := utils.GetPageAndPageSize(c)

	works, total, err := h.workService.GetArtistWorks(uint(artistID), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取作品列表失败")
		return
	}

	response.Page(c, works, total, page, pageSize)
}

func (h *WorkHandler) UpdateWork(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.workService.UpdateWork(uint(id), userID, req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) DeleteWork(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.workService.DeleteWork(uint(id), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) BatchPublish(c *gin.Context) {
	var req service.BatchPublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.workService.BatchPublish(&req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "批量发布失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) CreateAlbum(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.CreateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	album, err := h.workService.CreateAlbum(userID, &req)
	if err != nil {
		response.InternalError(c, "创建专辑失败")
		return
	}

	response.Success(c, album)
}

func (h *WorkHandler) GetAlbumByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	album, err := h.workService.GetAlbumByID(uint(id))
	if err != nil {
		response.NotFound(c, "专辑不存在")
		return
	}

	response.Success(c, album)
}

func (h *WorkHandler) ListAlbums(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)
	keyword := c.Query("keyword")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	artistID, _ := strconv.ParseUint(c.DefaultQuery("artist_id", "0"), 10, 64)

	albums, total, err := h.workService.ListAlbums(page, pageSize, keyword, uint(artistID), status)
	if err != nil {
		response.InternalError(c, "获取专辑列表失败")
		return
	}

	response.Page(c, albums, total, page, pageSize)
}

func (h *WorkHandler) UpdateAlbum(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.workService.UpdateAlbum(uint(id), userID, req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) DeleteAlbum(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.workService.DeleteAlbum(uint(id), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) AddWorkToAlbum(c *gin.Context) {
	userID := jwt.GetUserID(c)

	albumID, err := strconv.ParseUint(c.Param("album_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		WorkID uint `json:"work_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.workService.AddWorkToAlbum(uint(albumID), req.WorkID, userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "添加失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) ListTags(c *gin.Context) {
	keyword := c.Query("keyword")

	tags, err := h.workService.ListTags(keyword)
	if err != nil {
		response.InternalError(c, "获取标签失败")
		return
	}

	response.Success(c, tags)
}

func (h *WorkHandler) RecordPlay(c *gin.Context) {
	userID := jwt.GetUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Duration int `json:"duration"`
	}
	c.ShouldBindJSON(&req)

	ip := c.ClientIP()

	err = h.workService.RecordPlay(uint(id), userID, req.Duration, ip)
	if err != nil {
		response.InternalError(c, "记录播放失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) ApproveWork(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.workService.UpdateWorkStatus(uint(id), 2)
	if err != nil {
		response.InternalError(c, "审核失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) RejectWork(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	err = h.workService.UpdateWorkStatus(uint(id), 3)
	if err != nil {
		response.InternalError(c, "拒绝失败")
		return
	}

	response.Success(c, nil)
}

func (h *WorkHandler) UpdateWorkStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.workService.UpdateWorkStatus(uint(id), req.Status)
	if err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}
