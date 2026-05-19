package channels

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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

type DingTalkAdapter struct{}

func NewDingTalkAdapter() *DingTalkAdapter {
	return &DingTalkAdapter{}
}

func (d *DingTalkAdapter) Send(message *models.Message, config string) error {
	var dtConfig models.DingTalkConfig
	if err := json.Unmarshal([]byte(config), &dtConfig); err != nil {
		return fmt.Errorf("invalid dingtalk config: %w", err)
	}

	webhookURL := dtConfig.WebhookURL
	if dtConfig.Secret != "" {
		timestamp := time.Now().UnixMilli()
		sign := d.sign(dtConfig.Secret, timestamp)
		webhookURL = fmt.Sprintf("%s&timestamp=%d&sign=%s", dtConfig.WebhookURL, timestamp, sign)
	}

	if dtConfig.AccessToken != "" {
		if webhookURL == "" {
			return d.sendMessageByAccessToken(dtConfig.AccessToken, message)
		}
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message.Content,
		},
	}

	return d.sendRequest(webhookURL, payload)
}

func (d *DingTalkAdapter) sendMessageByAccessToken(accessToken string, message *models.Message) error {
	payload := map[string]interface{}{
		"userid_list": message.Recipient,
		"msg": map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": message.Content,
			},
		},
	}

	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=%s", accessToken)
	return d.sendRequest(url, payload)
}

func (d *DingTalkAdapter) sendRequest(url string, payload interface{}) error {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
	logger.Debug("dingtalk response", zap.String("response", string(respBody)))

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("dingtalk error: %v, %v", result["errcode"], result["errmsg"])
	}

	return nil
}

func (d *DingTalkAdapter) sign(secret string, timestamp int64) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	signData := h.Sum(nil)
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signData))
}

func (d *DingTalkAdapter) TestConnection(config string) error {
	var dtConfig models.DingTalkConfig
	if err := json.Unmarshal([]byte(config), &dtConfig); err != nil {
		return fmt.Errorf("invalid dingtalk config: %w", err)
	}

	if dtConfig.WebhookURL == "" && dtConfig.AccessToken == "" {
		return fmt.Errorf("webhook_url or access_token is required")
	}

	if dtConfig.WebhookURL != "" {
		testURL := dtConfig.WebhookURL
		if dtConfig.Secret != "" {
			timestamp := time.Now().UnixMilli()
			sign := d.sign(dtConfig.Secret, timestamp)
			testURL = fmt.Sprintf("%s&timestamp=%d&sign=%s", dtConfig.WebhookURL, timestamp, sign)
		}

		payload := map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": "【连接测试】通知中心渠道连通性测试",
			},
		}

		body, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", testURL, bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("create request failed: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("connect to dingtalk webhook failed: %w", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(respBody, &result)

		if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
			return fmt.Errorf("dingtalk webhook error: %v, %v", result["errcode"], result["errmsg"])
		}
	}

	if dtConfig.AccessToken != "" {
		testURL := "https://oapi.dingtalk.com/topapi/message/corpconversation/getsendprogress?access_token=" + dtConfig.AccessToken
		payload := map[string]interface{}{
			"agent_id": 0,
			"task_id":  0,
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(testURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("connect to dingtalk api failed: %w", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(respBody, &result)

		if errcode, ok := result["errcode"].(float64); ok && errcode != 0 && errcode != 40004 && errcode != 41043 {
			return fmt.Errorf("dingtalk api error: %v, %v", result["errcode"], result["errmsg"])
		}
	}

	logger.Info("dingtalk connection test successful")
	return nil
}
