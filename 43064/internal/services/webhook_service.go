package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/notification-center/internal/config"
	"github.com/notification-center/internal/database"
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/utils"
	"go.uber.org/zap"
)

type WebhookService struct {
	cfg *config.WebhookConfig
}

func NewWebhookService(cfg *config.WebhookConfig) *WebhookService {
	return &WebhookService{cfg: cfg}
}

func (ws *WebhookService) Notify(url string, message *models.Message) error {
	payload := map[string]interface{}{
		"event":      "message.status_update",
		"message_id": message.MessageID,
		"status":     message.Status,
		"channel_id": message.ChannelID,
		"recipient":  message.Recipient,
		"sent_at":    message.SentAt,
		"duration_ms": message.DurationMs,
		"retry_count": message.RetryCount,
		"error":      message.LastError,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)

	var lastErr error
	for attempt := 0; attempt < ws.cfg.MaxRetries; attempt++ {
		err := ws.sendRequest(url, body, message)
		if err == nil {
			return nil
		}
		lastErr = err

		if attempt < ws.cfg.MaxRetries-1 {
			backoff := utils.CalculateBackoff(attempt+1, ws.cfg.Backoff, ws.cfg.Backoff*10, 2)
			time.Sleep(backoff)
		}
	}

	logger.Error("webhook notify failed after retries",
		zap.String("url", url),
		zap.String("message_id", message.MessageID),
		zap.Error(lastErr),
	)
	return lastErr
}

func (ws *WebhookService) NotifySubscribers(message *models.Message) {
	var webhooks []models.Webhook
	if err := database.DB.Where("enabled = ?", true).Find(&webhooks).Error; err != nil {
		logger.Error("find enabled webhooks failed", zap.Error(err))
		return
	}

	for _, webhook := range webhooks {
		if !ws.shouldNotify(&webhook, message) {
			continue
		}

		go func(wh models.Webhook) {
			if err := ws.Notify(wh.URL, message); err != nil {
				logger.Error("notify webhook subscriber failed",
					zap.String("webhook", wh.Name),
					zap.String("url", wh.URL),
					zap.String("message_id", message.MessageID),
					zap.Error(err),
				)
			} else {
				logger.Info("webhook subscriber notified",
					zap.String("webhook", wh.Name),
					zap.String("url", wh.URL),
					zap.String("message_id", message.MessageID),
					zap.String("status", string(message.Status)),
				)
			}
		}(webhook)
	}
}

func (ws *WebhookService) shouldNotify(webhook *models.Webhook, message *models.Message) bool {
	if len(webhook.Events) > 0 {
		eventMatched := false
		for _, event := range webhook.Events {
			if event == "*" || event == string(message.Status) {
				eventMatched = true
				break
			}
		}
		if !eventMatched {
			return false
		}
	}

	if len(webhook.ChannelIDs) > 0 {
		channelMatched := false
		for _, chID := range webhook.ChannelIDs {
			if chID == message.ChannelID {
				channelMatched = true
				break
			}
		}
		if !channelMatched {
			return false
		}
	}

	return true
}

func (ws *WebhookService) sendRequest(url string, body []byte, message *models.Message) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Duration(ws.cfg.Timeout) * time.Millisecond}
	startTime := time.Now()

	resp, err := client.Do(req)
	duration := time.Since(startTime).Milliseconds()
	if err != nil {
		ws.saveLog(0, url, message.MessageID, "message.status_update", "failed", 0, string(body), "", err.Error(), duration, 0)
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	status := "success"
	if resp.StatusCode >= 400 {
		status = "failed"
		err = fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	ws.saveLog(0, url, message.MessageID, "message.status_update", status, resp.StatusCode, string(body), string(respBody), "", duration, 0)

	if resp.StatusCode >= 400 {
		return err
	}

	return nil
}

func (ws *WebhookService) Create(webhook *models.Webhook) (*models.Webhook, error) {
	if err := database.DB.Create(webhook).Error; err != nil {
		logger.Error("create webhook failed", zap.Error(err))
		return nil, err
	}
	logger.Info("webhook created", zap.Uint("id", webhook.ID), zap.String("name", webhook.Name))
	return webhook, nil
}

func (ws *WebhookService) GetByID(id uint) (*models.Webhook, error) {
	var webhook models.Webhook
	if err := database.DB.First(&webhook, id).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (ws *WebhookService) List(page, pageSize int) ([]models.Webhook, int64, error) {
	var webhooks []models.Webhook
	var total int64

	if err := database.DB.Model(&models.Webhook{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := database.DB.Order("id DESC").Offset(offset).Limit(pageSize).Find(&webhooks).Error; err != nil {
		return nil, 0, err
	}

	return webhooks, total, nil
}

func (ws *WebhookService) Update(id uint, updates map[string]interface{}) (*models.Webhook, error) {
	webhook, err := ws.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := database.DB.Model(webhook).Updates(updates).Error; err != nil {
		return nil, err
	}

	return webhook, nil
}

func (ws *WebhookService) Delete(id uint) error {
	return database.DB.Delete(&models.Webhook{}, id).Error
}

func (ws *WebhookService) saveLog(webhookID uint, url, messageID, event, status string, statusCode int, request, response, errStr string, durationMs int64, retryCount int) {
	log := &models.WebhookLog{
		WebhookID:  webhookID,
		MessageID:  messageID,
		Event:      event,
		Status:     status,
		StatusCode: statusCode,
		Request:    request,
		Response:   response,
		Error:      errStr,
		DurationMs: durationMs,
		RetryCount: retryCount,
	}
	database.DB.Create(log)
}
