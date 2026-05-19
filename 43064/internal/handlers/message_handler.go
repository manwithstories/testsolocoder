package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/services"
)

type MessageHandler struct {
	queueSvc      *services.QueueService
	templateSvc   *services.TemplateService
	channelSvc    *services.ChannelService
	statisticsSvc *services.StatisticsService
}

func NewMessageHandler(queueSvc *services.QueueService, templateSvc *services.TemplateService,
	channelSvc *services.ChannelService, statisticsSvc *services.StatisticsService) *MessageHandler {
	return &MessageHandler{
		queueSvc:      queueSvc,
		templateSvc:   templateSvc,
		channelSvc:    channelSvc,
		statisticsSvc: statisticsSvc,
	}
}

type SendRequest struct {
	ChannelID    uint                   `json:"channel_id" binding:"required"`
	TemplateID   *uint                  `json:"template_id"`
	Recipient    string                 `json:"recipient" binding:"required"`
	RecipientID  *uint                  `json:"recipient_id"`
	Subject      string                 `json:"subject"`
	Content      string                 `json:"content"`
	Variables    map[string]interface{} `json:"variables"`
	Priority     models.MessagePriority `json:"priority"`
	ScheduledAt  *time.Time             `json:"scheduled_at"`
	WebhookURL   string                 `json:"webhook_url"`
}

func (h *MessageHandler) Send(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	channel, err := h.channelSvc.GetByID(req.ChannelID)
	if err != nil {
		c.Error(err)
		return
	}

	if !channel.Enabled {
		c.Error(errors.ChannelError("channel is disabled", nil))
		return
	}

	content := req.Content
	subject := req.Subject

	if req.TemplateID != nil {
		template, err := h.templateSvc.GetByID(*req.TemplateID)
		if err != nil {
			c.Error(err)
			return
		}

		if template.ChannelID != req.ChannelID {
			c.Error(errors.InvalidParams("template does not belong to this channel", nil))
			return
		}

		if req.Variables == nil {
			req.Variables = make(map[string]interface{})
		}

		subject, content, err = h.templateSvc.Render(template, req.Variables)
		if err != nil {
			c.Error(err)
			return
		}
	}

	if content == "" {
		c.Error(errors.MissingRequiredField("content"))
		return
	}

	message := &models.Message{
		MessageID:   models.GenerateMessageID(),
		ChannelID:   req.ChannelID,
		TemplateID:  req.TemplateID,
		Recipient:   req.Recipient,
		RecipientID: req.RecipientID,
		Subject:     subject,
		Content:     content,
		Variables:   req.Variables,
		Status:      models.MessageStatusPending,
		Priority:    req.Priority,
		ScheduledAt: req.ScheduledAt,
		WebhookURL:  req.WebhookURL,
		MaxRetries:  5,
	}

	if message.Priority == 0 {
		message.Priority = models.PriorityNormal
	}

	requestID, _ := c.Get("request_id")
	message.RequestID = requestID.(string)

	if err := h.queueSvc.Enqueue(message); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"code":      0,
		"message":   "message queued",
		"message_id": message.MessageID,
	})
}

func (h *MessageHandler) BatchSend(c *gin.Context) {
	var req struct {
		Messages []SendRequest `json:"messages" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidParams("invalid request body", err))
		return
	}

	requestID, _ := c.Get("request_id")
	results := make([]map[string]interface{}, 0, len(req.Messages))

	for _, msgReq := range req.Messages {
		result := make(map[string]interface{})
		result["recipient"] = msgReq.Recipient

		channel, err := h.channelSvc.GetByID(msgReq.ChannelID)
		if err != nil {
			result["success"] = false
			result["error"] = err.Error()
			results = append(results, result)
			continue
		}

		if !channel.Enabled {
			result["success"] = false
			result["error"] = "channel is disabled"
			results = append(results, result)
			continue
		}

		content := msgReq.Content
		subject := msgReq.Subject

		if msgReq.TemplateID != nil {
			template, err := h.templateSvc.GetByID(*msgReq.TemplateID)
			if err != nil {
				result["success"] = false
				result["error"] = err.Error()
				results = append(results, result)
				continue
			}

			if msgReq.Variables == nil {
				msgReq.Variables = make(map[string]interface{})
			}

			subject, content, err = h.templateSvc.Render(template, msgReq.Variables)
			if err != nil {
				result["success"] = false
				result["error"] = err.Error()
				results = append(results, result)
				continue
			}
		}

		message := &models.Message{
			MessageID:  models.GenerateMessageID(),
			ChannelID:  msgReq.ChannelID,
			TemplateID: msgReq.TemplateID,
			Recipient:  msgReq.Recipient,
			Subject:    subject,
			Content:    content,
			Variables:  msgReq.Variables,
			Status:     models.MessageStatusPending,
			Priority:   msgReq.Priority,
			ScheduledAt: msgReq.ScheduledAt,
			WebhookURL: msgReq.WebhookURL,
			MaxRetries: 5,
			RequestID:  requestID.(string),
		}

		if message.Priority == 0 {
			message.Priority = models.PriorityNormal
		}

		if err := h.queueSvc.Enqueue(message); err != nil {
			result["success"] = false
			result["error"] = err.Error()
		} else {
			result["success"] = true
			result["message_id"] = message.MessageID
		}
		results = append(results, result)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"code":    0,
		"message": "batch send completed",
		"data":    results,
	})
}

func (h *MessageHandler) GetStatus(c *gin.Context) {
	messageID := c.Param("message_id")

	message, err := h.statisticsSvc.GetMessageByID(messageID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    message,
	})
}

func (h *MessageHandler) List(c *gin.Context) {
	channelID, _ := strconv.ParseUint(c.Query("channel_id"), 10, 32)
	status := c.Query("status")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var start, end time.Time
	if startTime != "" {
		start, _ = time.Parse(time.RFC3339, startTime)
	}
	if endTime != "" {
		end, _ = time.Parse(time.RFC3339, endTime)
	}

	messages, total, err := h.statisticsSvc.GetMessageList(uint(channelID), models.MessageStatus(status), start, end, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     messages,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *MessageHandler) GetStats(c *gin.Context) {
	channelID, _ := strconv.ParseUint(c.Query("channel_id"), 10, 32)
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var start, end time.Time
	if startTime != "" {
		start, _ = time.Parse(time.RFC3339, startTime)
	}
	if endTime != "" {
		end, _ = time.Parse(time.RFC3339, endTime)
	}

	stats, err := h.statisticsSvc.GetMessageStats(uint(channelID), start, end)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

func (h *MessageHandler) GetChannelStats(c *gin.Context) {
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var start, end time.Time
	if startTime != "" {
		start, _ = time.Parse(time.RFC3339, startTime)
	}
	if endTime != "" {
		end, _ = time.Parse(time.RFC3339, endTime)
	}

	stats, err := h.statisticsSvc.GetChannelStats(start, end)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

func (h *MessageHandler) GetDailyStats(c *gin.Context) {
	channelID, _ := strconv.ParseUint(c.Query("channel_id"), 10, 32)
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	if days < 1 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	stats, err := h.statisticsSvc.GetDailyStats(uint(channelID), days)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

func (h *MessageHandler) GetQueueStats(c *gin.Context) {
	pending, processing, err := h.queueSvc.GetQueueStats()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"pending":    pending,
			"processing": processing,
		},
	})
}
