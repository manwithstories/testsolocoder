package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/services"
)

type ChannelHandler struct {
	service *services.ChannelService
}

func NewChannelHandler(service *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{service: service}
}

func (h *ChannelHandler) Create(c *gin.Context) {
	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	created, err := h.service.Create(&channel)
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

func (h *ChannelHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid channel id", err))
		return
	}

	channel, err := h.service.GetByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    channel,
	})
}

func (h *ChannelHandler) List(c *gin.Context) {
	enabledOnly := c.Query("enabled") == "true"
	channelType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	channels, total, err := h.service.List(enabledOnly, models.ChannelType(channelType), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     channels,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *ChannelHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid channel id", err))
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

func (h *ChannelHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid channel id", err))
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

func (h *ChannelHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid channel id", err))
		return
	}

	if err := h.service.Enable(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *ChannelHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid channel id", err))
		return
	}

	if err := h.service.Disable(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *ChannelHandler) UpdatePriority(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid channel id", err))
		return
	}

	var req struct {
		Priority int `json:"priority" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	if err := h.service.UpdatePriority(uint(id), req.Priority); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *ChannelHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(errors.InvalidParams("invalid channel id", err))
		return
	}

	if err := h.service.TestConnection(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "connection test successful",
	})
}
