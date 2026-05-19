package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"go.uber.org/zap"
)

type WeChatAdapter struct {
	accessTokenCache map[string]*tokenInfo
	mu               sync.RWMutex
}

type tokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

func NewWeChatAdapter() *WeChatAdapter {
	return &WeChatAdapter{
		accessTokenCache: make(map[string]*tokenInfo),
	}
}

func (w *WeChatAdapter) Send(message *models.Message, config string) error {
	var wxConfig models.WeChatConfig
	if err := json.Unmarshal([]byte(config), &wxConfig); err != nil {
		return fmt.Errorf("invalid wechat config: %w", err)
	}

	accessToken, err := w.getAccessToken(wxConfig)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"touser":  message.Recipient,
		"msgtype": "text",
		"text": map[string]string{
			"content": message.Content,
		},
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)

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
	logger.Debug("wechat send response", zap.String("response", string(respBody)))

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		if errcode == 40001 || errcode == 42001 {
			w.clearAccessToken(wxConfig.AppID)
			accessToken, err := w.getAccessToken(wxConfig)
			if err == nil {
				url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				respBody, _ := io.ReadAll(resp.Body)
				json.Unmarshal(respBody, &result)
				if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
					return fmt.Errorf("wechat error: %v, %v", result["errcode"], result["errmsg"])
				}
				return nil
			}
		}
		return fmt.Errorf("wechat error: %v, %v", result["errcode"], result["errmsg"])
	}

	return nil
}

func (w *WeChatAdapter) getAccessToken(cfg models.WeChatConfig) (string, error) {
	w.mu.RLock()
	cached, exists := w.accessTokenCache[cfg.AppID]
	w.mu.RUnlock()

	if exists && time.Now().Before(cached.ExpiresAt) {
		return cached.Token, nil
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", cfg.AppID, cfg.AppSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return "", fmt.Errorf("get access token failed: %v, %v", result["errcode"], result["errmsg"])
	}

	token, _ := result["access_token"].(string)
	expiresIn, _ := result["expires_in"].(float64)

	w.mu.Lock()
	w.accessTokenCache[cfg.AppID] = &tokenInfo{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(expiresIn-300) * time.Second),
	}
	w.mu.Unlock()

	return token, nil
}

func (w *WeChatAdapter) clearAccessToken(appID string) {
	w.mu.Lock()
	delete(w.accessTokenCache, appID)
	w.mu.Unlock()
}

func (w *WeChatAdapter) TestConnection(config string) error {
	var wxConfig models.WeChatConfig
	if err := json.Unmarshal([]byte(config), &wxConfig); err != nil {
		return fmt.Errorf("invalid wechat config: %w", err)
	}

	if wxConfig.AppID == "" || wxConfig.AppSecret == "" {
		return fmt.Errorf("app_id and app_secret are required")
	}

	_, err := w.getAccessToken(wxConfig)
	if err != nil {
		return err
	}

	logger.Info("wechat connection test successful", zap.String("app_id", wxConfig.AppID))
	return nil
}
