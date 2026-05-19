package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/services"
)

type WebhookHandler struct {
	service *services.WebhookService
}

func NewWebhookHandler(service *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{service: service}
}

func (h *WebhookHandler) Create(c *gin.Context) {
	var webhook models.Webhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	created, err := h.service.Create(&webhook)
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

func (h *WebhookHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid webhook id", err))
		return
	}

	webhook, err := h.service.GetByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    webhook,
	})
}

func (h *WebhookHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	webhooks, total, err := h.service.List(page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     webhooks,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *WebhookHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid webhook id", err))
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

func (h *WebhookHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid webhook id", err))
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
