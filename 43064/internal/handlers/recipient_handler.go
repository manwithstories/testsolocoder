package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/services"
)

type RecipientHandler struct {
	service *services.RecipientService
}

func NewRecipientHandler(service *services.RecipientService) *RecipientHandler {
	return &RecipientHandler{service: service}
}

func (h *RecipientHandler) Create(c *gin.Context) {
	var recipient models.Recipient
	if err := c.ShouldBindJSON(&recipient); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	created, err := h.service.Create(&recipient)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "success",
		"data":    created,
	})
}

func (h *RecipientHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid recipient id", err))
		return
	}

	recipient, err := h.service.GetByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    recipient,
	})
}

func (h *RecipientHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	enabledOnly := c.Query("enabled") == "true"
	tagIDs := parseUintIDs(c.Query("tag_ids"))
	groupIDs := parseUintIDs(c.Query("group_ids"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	recipients, total, err := h.service.List(keyword, tagIDs, groupIDs, enabledOnly, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     recipients,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *RecipientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid recipient id", err))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	updated, err := h.service.Update(uint(id), updates)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    updated,
	})
}

func (h *RecipientHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid recipient id", err))
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *RecipientHandler) AddTags(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid recipient id", err))
		return
	}

	var req struct {
		TagIDs []uint `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if err := h.service.AddTags(uint(id), req.TagIDs); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *RecipientHandler) RemoveTags(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid recipient id", err))
		return
	}

	var req struct {
		TagIDs []uint `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if err := h.service.RemoveTags(uint(id), req.TagIDs); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *RecipientHandler) AddToGroups(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid recipient id", err))
		return
	}

	var req struct {
		GroupIDs []uint `json:"group_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if err := h.service.AddToGroups(uint(id), req.GroupIDs); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *RecipientHandler) RemoveFromGroups(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid recipient id", err))
		return
	}

	var req struct {
		GroupIDs []uint `json:"group_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if err := h.service.RemoveFromGroups(uint(id), req.GroupIDs); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *RecipientHandler) BatchImport(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(errors.InvalidParams("failed to get file", err))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.Error(errors.InvalidParams("failed to parse CSV", err))
		return
	}

	if len(records) > 0 {
		records = records[1:]
	}

	success, failed, errs := h.service.BatchImport(records)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "import completed",
		"data": gin.H{
			"success": success,
			"failed":  failed,
			"errors":  errs,
		},
	})
}

func (h *RecipientHandler) Export(c *gin.Context) {
	groupIDs := parseUintIDs(c.Query("group_ids"))
	tagIDs := parseUintIDs(c.Query("tag_ids"))

	data, err := h.service.Export(groupIDs, tagIDs)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=recipients.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	for _, row := range data {
		writer.Write(row)
	}
}

func (h *RecipientHandler) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	created, err := h.service.CreateTag(&tag)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "success",
		"data":    created,
	})
}

func (h *RecipientHandler) ListTags(c *gin.Context) {
	tags, err := h.service.ListTags()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tags,
	})
}

func (h *RecipientHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid tag id", err))
		return
	}

	if err := h.service.DeleteTag(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *RecipientHandler) CreateGroup(c *gin.Context) {
	var group models.RecipientGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	created, err := h.service.CreateGroup(&group)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "success",
		"data":    created,
	})
}

func (h *RecipientHandler) ListGroups(c *gin.Context) {
	groups, err := h.service.ListGroups()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    groups,
	})
}

func (h *RecipientHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid group id", err))
		return
	}

	group, err := h.service.GetGroup(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    group,
	})
}

func (h *RecipientHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid group id", err))
		return
	}

	if err := h.service.DeleteGroup(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func parseUintIDs(s string) []uint {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		if id, err := strconv.ParseUint(p, 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids
}
