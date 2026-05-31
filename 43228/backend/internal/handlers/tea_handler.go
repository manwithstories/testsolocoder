package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"tea-platform/internal/repository"
	"tea-platform/internal/service"
	"tea-platform/internal/utils"
)

type TeaHandler struct {
	teaService *service.TeaService
}

func NewTeaHandler() *TeaHandler {
	return &TeaHandler{
		teaService: service.NewTeaService(),
	}
}

func (h *TeaHandler) CreateTea(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	var req service.CreateTeaRequest
	if !utils.Validate(c, &req) {
		return
	}

	tea, err := h.teaService.CreateTea(userID, &req)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":   tea.ID,
		"name": tea.Name,
	})
}

func (h *TeaHandler) UpdateTea(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "茶叶ID格式错误")
		return
	}

	var req service.UpdateTeaRequest
	if !utils.Validate(c, &req) {
		return
	}

	if err := h.teaService.UpdateTea(userID, uint(id), &req); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TeaHandler) DeleteTea(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "茶叶ID格式错误")
		return
	}

	if err := h.teaService.DeleteTea(userID, uint(id)); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TeaHandler) GetTeaDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "茶叶ID格式错误")
		return
	}

	tea, err := h.teaService.GetTeaByID(uint(id))
	if err != nil {
		utils.Fail(c, utils.CodeNotFound, err.Error())
		return
	}

	utils.Success(c, tea)
}

type TeaListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Type     string `form:"type"`
	Origin   string `form:"origin"`
	Year     int    `form:"year"`
	Grade    string `form:"grade"`
	Keyword  string `form:"keyword"`
}

func (h *TeaHandler) GetTeaList(c *gin.Context) {
	var query TeaListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Fail(c, utils.CodeBadRequest, "参数错误")
		return
	}

	filters := repository.TeaFilters{
		Type:    query.Type,
		Origin:  query.Origin,
		Year:    query.Year,
		Grade:   query.Grade,
		Keyword: query.Keyword,
	}

	teas, total, err := h.teaService.GetTeaList(query.Page, query.PageSize, filters)
	if err != nil {
		utils.Fail(c, utils.CodeInternalError, "查询茶叶列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":  teas,
		"total": total,
	})
}

func (h *TeaHandler) UploadTeaImage(c *gin.Context) {
	h.getCurrentUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "茶叶ID格式错误")
		return
	}

	imageType := c.Query("image_type")
	if imageType == "" {
		imageType = "detail"
	}
	if imageType != "main" && imageType != "detail" && imageType != "packaging" {
		utils.Fail(c, utils.CodeBadRequest, "图片类型错误，支持 main/detail/packaging")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "请上传文件")
		return
	}

	if err := h.teaService.UploadTeaImage(uint(id), imageType, file); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TeaHandler) DeleteTeaImage(c *gin.Context) {
	h.getCurrentUserID(c)

	imageIDStr := c.Param("image_id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "图片ID格式错误")
		return
	}

	if err := h.teaService.DeleteTeaImage(uint(imageID)); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TeaHandler) AddTraceability(c *gin.Context) {
	h.getCurrentUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "茶叶ID格式错误")
		return
	}

	var req service.AddTraceabilityRequest
	if !utils.Validate(c, &req) {
		return
	}

	if err := h.teaService.AddTraceability(uint(id), &req); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TeaHandler) GetTraceability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "茶叶ID格式错误")
		return
	}

	traces, err := h.teaService.GetTraceability(uint(id))
	if err != nil {
		utils.Fail(c, utils.CodeNotFound, err.Error())
		return
	}

	utils.Success(c, traces)
}

func (h *TeaHandler) getCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, utils.CodeUnauthorized, utils.MsgUnauthorized)
		return 0
	}
	uid, ok := userID.(uint)
	if !ok {
		utils.Fail(c, utils.CodeUnauthorized, "用户ID类型错误")
		return 0
	}
	return uid
}
