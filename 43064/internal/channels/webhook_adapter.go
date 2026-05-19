package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"go.uber.org/zap"
)

type WebhookChannelAdapter struct{}

func NewWebhookChannelAdapter() *WebhookChannelAdapter {
	return &WebhookChannelAdapter{}
}

func (w *WebhookChannelAdapter) Send(message *models.Message, config string) error {
	var webhookConfig models.WebhookConfig
	if err := json.Unmarshal([]byte(config), &webhookConfig); err != nil {
		return fmt.Errorf("invalid webhook config: %w", err)
	}

	method := webhookConfig.Method
	if method == "" {
		method = "POST"
	}

	timeout := webhookConfig.Timeout
	if timeout == 0 {
		timeout = 5000
	}

	payload := map[string]interface{}{
		"recipient": message.Recipient,
		"subject":   message.Subject,
		"content":   message.Content,
		"variables": message.Variables,
		"message_id": message.MessageID,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, webhookConfig.URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range webhookConfig.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logger.Debug("webhook channel response",
		zap.String("url", webhookConfig.URL),
		zap.Int("status", resp.StatusCode),
		zap.String("response", string(respBody)),
	)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook channel failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (w *WebhookChannelAdapter) TestConnection(config string) error {
	var webhookConfig models.WebhookConfig
	if err := json.Unmarshal([]byte(config), &webhookConfig); err != nil {
		return fmt.Errorf("invalid webhook config: %w", err)
	}

	if webhookConfig.URL == "" {
		return fmt.Errorf("url is required")
	}

	method := webhookConfig.Method
	if method == "" {
		method = "GET"
	}

	timeout := webhookConfig.Timeout
	if timeout == 0 {
		timeout = 5000
	}

	payload := map[string]interface{}{
		"test": true,
		"message": "Notification center connection test",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, webhookConfig.URL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range webhookConfig.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connect to webhook failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	logger.Info("webhook channel test successful", zap.String("url", webhookConfig.URL), zap.Int("status", resp.StatusCode))
	return nil
}
