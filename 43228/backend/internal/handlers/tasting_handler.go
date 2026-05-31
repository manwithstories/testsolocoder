package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"tea-platform/internal/service"
	"tea-platform/internal/utils"
)

type TastingHandler struct {
	tastingService *service.TastingService
}

func NewTastingHandler() *TastingHandler {
	return &TastingHandler{
		tastingService: service.NewTastingService(),
	}
}

func (h *TastingHandler) CreateTastingRecord(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	var req service.CreateTastingRequest
	if !utils.Validate(c, &req) {
		return
	}

	record, err := h.tastingService.CreateTastingRecord(userID, &req)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id": record.ID,
	})
}

func (h *TastingHandler) UpdateTastingRecord(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "品鉴记录ID格式错误")
		return
	}

	var req service.UpdateTastingRequest
	if !utils.Validate(c, &req) {
		return
	}

	if err := h.tastingService.UpdateTastingRecord(userID, uint(id), &req); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TastingHandler) DeleteTastingRecord(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "品鉴记录ID格式错误")
		return
	}

	if err := h.tastingService.DeleteTastingRecord(userID, uint(id)); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TastingHandler) GetTastingRecordDetail(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "品鉴记录ID格式错误")
		return
	}

	record, err := h.tastingService.GetTastingRecordDetail(uint(id))
	if err != nil {
		utils.Fail(c, utils.CodeNotFound, err.Error())
		return
	}

	utils.Success(c, record)
}

type TastingListQuery struct {
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=10"`
	UserID          uint   `form:"user_id"`
	TeaID           uint   `form:"tea_id"`
	MinOverallScore float64 `form:"min_overall_score"`
	MaxOverallScore float64 `form:"max_overall_score"`
	StartDate       string `form:"start_date"`
	EndDate         string `form:"end_date"`
	BrewMethod      string `form:"brew_method"`
	Keyword         string `form:"keyword"`
}

func (h *TastingHandler) GetTastingRecordList(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	var query TastingListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Fail(c, utils.CodeBadRequest, "参数错误")
		return
	}

	req := &service.TastingFilterRequest{
		UserID:          query.UserID,
		TeaID:           query.TeaID,
		MinOverallScore: query.MinOverallScore,
		MaxOverallScore: query.MaxOverallScore,
		StartDate:       query.StartDate,
		EndDate:         query.EndDate,
		BrewMethod:      query.BrewMethod,
		Keyword:         query.Keyword,
	}

	records, total, err := h.tastingService.GetTastingRecordList(query.Page, query.PageSize, req)
	if err != nil {
		utils.Fail(c, utils.CodeInternalError, "查询品鉴记录列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":  records,
		"total": total,
	})
}

func (h *TastingHandler) UploadTastingImage(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "品鉴记录ID格式错误")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, utils.CodeBadRequest, "请上传文件")
		return
	}

	h.tastingService.SetContext(c)
	if err := h.tastingService.UploadTastingImage(uint(id), userID, file); err != nil {
		utils.Fail(c, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *TastingHandler) GetUserTastingStats(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return
	}

	stats, err := h.tastingService.GetUserTastingStats(userID)
	if err != nil {
		utils.Fail(c, utils.CodeInternalError, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *TastingHandler) getCurrentUserID(c *gin.Context) uint {
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
