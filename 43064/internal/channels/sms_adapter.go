package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"go.uber.org/zap"
)

type SMSAdapter struct{}

func NewSMSAdapter() *SMSAdapter {
	return &SMSAdapter{}
}

func (s *SMSAdapter) Send(message *models.Message, config string) error {
	var smsConfig models.SMSConfig
	if err := json.Unmarshal([]byte(config), &smsConfig); err != nil {
		return fmt.Errorf("invalid sms config: %w", err)
	}

	switch smsConfig.Provider {
	case "aliyun":
		return s.sendAliyun(smsConfig, message)
	case "tencent":
		return s.sendTencent(smsConfig, message)
	default:
		return s.sendGeneric(smsConfig, message)
	}
}

func (s *SMSAdapter) sendAliyun(cfg models.SMSConfig, message *models.Message) error {
	params := url.Values{}
	params.Set("Action", "SendSms")
	params.Set("PhoneNumbers", message.Recipient)
	params.Set("SignName", cfg.SignName)
	params.Set("TemplateCode", cfg.TemplateCode)

	params.Set("TemplateParam", message.Content)

	reqURL := fmt.Sprintf("%s?%s", cfg.Endpoint, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Debug("aliyun sms response", zap.String("response", string(body)))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("aliyun sms failed: %s", string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if code, ok := result["Code"].(string); ok && code != "OK" {
		return fmt.Errorf("aliyun sms error: %s, %s", code, result["Message"])
	}

	return nil
}

func (s *SMSAdapter) sendTencent(cfg models.SMSConfig, message *models.Message) error {
	payload := map[string]interface{}{
		"PhoneNumberSet":   []string{message.Recipient},
		"SignName":         cfg.SignName,
		"TemplateCode":     cfg.TemplateCode,
		"TemplateParamSet": []string{message.Content},
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", cfg.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logger.Debug("tencent sms response", zap.String("response", string(respBody)))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tencent sms failed: %s", string(respBody))
	}

	return nil
}

func (s *SMSAdapter) sendGeneric(cfg models.SMSConfig, message *models.Message) error {
	formData := url.Values{}
	formData.Set("mobile", message.Recipient)
	formData.Set("content", message.Content)
	formData.Set("apikey", cfg.APIKey)

	req, err := http.NewRequest("POST", cfg.Endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logger.Debug("generic sms response", zap.String("response", string(respBody)))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sms failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *SMSAdapter) TestConnection(config string) error {
	var smsConfig models.SMSConfig
	if err := json.Unmarshal([]byte(config), &smsConfig); err != nil {
		return fmt.Errorf("invalid sms config: %w", err)
	}

	if smsConfig.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(smsConfig.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to connect to sms endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("sms endpoint returned server error: %d", resp.StatusCode)
	}

	logger.Info("sms connection test successful", zap.String("provider", smsConfig.Provider), zap.String("endpoint", smsConfig.Endpoint))
	return nil
}
