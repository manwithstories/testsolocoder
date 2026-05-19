package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/services"
)

type TemplateHandler struct {
	service *services.TemplateService
}

func NewTemplateHandler(service *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

func (h *TemplateHandler) Create(c *gin.Context) {
	var template models.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if err := h.service.ValidateTemplateSyntax(template.Content); err != nil {
		c.Error(err)
		return
	}

	created, err := h.service.Create(&template)
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

func (h *TemplateHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid template id", err))
		return
	}

	template, err := h.service.GetByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    template,
	})
}

func (h *TemplateHandler) List(c *gin.Context) {
	channelID, _ := strconv.ParseUint(c.Query("channel_id"), 10, 32)
	language := c.Query("language")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	templates, total, err := h.service.List(uint(channelID), language, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     templates,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *TemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid template id", err))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if content, ok := updates["content"].(string); ok {
		if err := h.service.ValidateTemplateSyntax(content); err != nil {
			c.Error(err)
			return
		}
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

func (h *TemplateHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid template id", err))
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

func (h *TemplateHandler) Render(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid template id", err))
		return
	}

	var req struct {
		Variables map[string]interface{} `json:"variables" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	template, err := h.service.GetByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	subject, content, err := h.service.Render(template, req.Variables)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"subject": subject,
			"content": content,
		},
	})
}

func (h *TemplateHandler) Validate(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if err := h.service.ValidateTemplateSyntax(req.Content); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "invalid",
			"data": gin.H{
				"valid":   false,
				"error":   err.Error(),
				"variables": []string{},
			},
		})
		return
	}

	variables := h.service.ExtractVariables(req.Content)
	varNames := make([]string, len(variables))
	for i, v := range variables {
		varNames[i] = v.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "valid",
		"data": gin.H{
			"valid":     true,
			"variables": varNames,
		},
	})
}
