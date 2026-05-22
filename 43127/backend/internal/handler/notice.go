package handler

import (
	"property-management/internal/database"
	"property-management/internal/model"
	"property-management/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoticeHandler struct{}

func NewNoticeHandler() *NoticeHandler {
	return &NoticeHandler{}
}

type NoticeRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Type     int    `json:"type"`
	Building string `json:"building"`
	IsTop    int    `json:"isTop"`
}

func (h *NoticeHandler) Create(c *gin.Context) {
	var req NoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	notice := model.Notice{
		Title:       req.Title,
		Content:     req.Content,
		Type:        req.Type,
		Building:    req.Building,
		PublisherID: userID,
		IsTop:       req.IsTop,
		Status:      1,
	}

	if err := database.DB.Create(&notice).Error; err != nil {
		utils.Error(c, 500, "Failed to create notice")
		return
	}

	utils.Success(c, notice)
}

func (h *NoticeHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var notice model.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		utils.Error(c, 404, "Notice not found")
		return
	}

	var req NoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&notice).Updates(map[string]interface{}{
		"title":    req.Title,
		"content":  req.Content,
		"type":     req.Type,
		"building": req.Building,
		"is_top":   req.IsTop,
	})

	utils.Success(c, nil)
}

type NoticeUpdateFieldsRequest struct {
	IsTop  *int    `json:"isTop"`
	Status *int    `json:"status"`
	Title  *string `json:"title"`
}

func (h *NoticeHandler) UpdateFields(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var notice model.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		utils.Error(c, 404, "Notice not found")
		return
	}

	var req NoticeUpdateFieldsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	updates := make(map[string]interface{})
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}

	if len(updates) > 0 {
		database.DB.Model(&notice).UpdateColumns(updates)
	}

	utils.Success(c, nil)
}

func (h *NoticeHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&model.Notice{}, id)
	utils.Success(c, nil)
}

func (h *NoticeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.Notice{})

	if noticeType := c.Query("type"); noticeType != "" {
		t, _ := strconv.Atoi(noticeType)
		query = query.Where("type = ?", t)
	}
	if building := c.Query("building"); building != "" {
		query = query.Where("building = ? OR building = ''", building)
	}

	var total int64
	query.Count(&total)

	var notices []model.Notice
	offset := (page - 1) * pageSize
	query.Preload("Publisher").
		Order("is_top DESC, id DESC").
		Offset(offset).Limit(pageSize).Find(&notices)

	utils.Success(c, gin.H{
		"list":     notices,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *NoticeHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var notice model.Notice
	if err := database.DB.Preload("Publisher").First(&notice, id).Error; err != nil {
		utils.Error(c, 404, "Notice not found")
		return
	}
	utils.Success(c, notice)
}

func (h *NoticeHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&model.Notice{}).Where("id = ?", id).Update("status", req.Status)
	utils.Success(c, nil)
}
