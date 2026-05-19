package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/apitester/apitester/internal/logger"
)

type Client struct {
	httpClient *http.Client
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
	Duration   time.Duration
	Error      error
}

func NewClient(timeout time.Duration) *Client {
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *Client) Do(ctx context.Context, method, url string, headers map[string]string, body string, retries int) *Response {
	var lastErr error
	var resp *Response

	for attempt := 0; attempt <= retries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt*100) * time.Millisecond
			time.Sleep(backoff)
			logger.Debug("Retrying request (attempt %d/%d): %s %s", attempt, retries, method, url)
		}

		resp = c.doOnce(ctx, method, url, headers, body)
		if resp.Error == nil {
			return resp
		}

		lastErr = resp.Error

		if !isRetryableError(resp.Error) {
			break
		}
	}

	if lastErr != nil {
		resp.Error = fmt.Errorf("request failed after %d attempts: %w", retries+1, lastErr)
	}

	return resp
}

func (c *Client) doOnce(ctx context.Context, method, url string, headers map[string]string, body string) *Response {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return &Response{Error: fmt.Errorf("failed to create request: %w", err)}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return &Response{
			Error:    err,
			Duration: duration,
		}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			StatusCode: resp.StatusCode,
			Error:      fmt.Errorf("failed to read response body: %w", err),
			Duration:   duration,
		}
	}

	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		respHeaders[k] = strings.Join(v, ", ")
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       string(respBody),
		Duration:   duration,
	}
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	retryableErrors := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"no such host",
		"temporary failure",
		"deadline exceeded",
		"context deadline exceeded",
		"net/http: TLS handshake timeout",
		"EOF",
	}

	for _, re := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), re) {
			return true
		}
	}

	if ne, ok := err.(net.Error); ok {
		return ne.Temporary() || ne.Timeout()
	}

	return false
}
